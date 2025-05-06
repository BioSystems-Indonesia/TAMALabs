package app

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	sqliteMigrate "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/rest"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/tcp"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/tcp/ba400"
	"github.com/oibacidem/lims-hl-seven/migrations"
	"github.com/oibacidem/lims-hl-seven/pkg/server"
	gormSqlite "gorm.io/driver/sqlite"

	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"

	_ "modernc.org/sqlite"
)

func provideTCP(config *config.Schema) *ba400.TCP {
	tcpEr := ba400.NewTCP(config)

	return tcpEr
}

func provideTCPServer(config *config.Schema, handler *tcp.HlSevenHandler) *server.TCP {
	s := server.NewTCP("1024")
	s.SetHandler(handler)
	err := s.Start()
	if err != nil {
		log.Println(err)
	}
	go s.Serve()
	return s
}

func provideRestServer(
	config *config.Schema,
	handlers *rest.Handler,
	validate *validator.Validate,
	deviceHandler *rest.DeviceHandler,
	serverControllerHandler *rest.ServerControllerHandler,
	testTemplateHandler *rest.TestTemplateHandler,
) server.RestServer {
	serv := server.NewRest(config.Port, validate)
	rest.RegisterMiddleware(serv.GetClient())
	rest.RegisterRoutes(serv.GetClient(), handlers, deviceHandler, serverControllerHandler, testTemplateHandler)
	return serv
}

func provideRestHandler(
	hlSevenHandler *rest.HlSevenHandler,
	healthCheck *rest.HealthCheckHandler,
	patientHandler *rest.PatientHandler,
	specimenHandler *rest.SpecimenHandler,
	workOrder *rest.WorkOrderHandler,
	featureListHandler *rest.FeatureListHandler,
	observationRequest *rest.ObservationRequestHandler,
	testTypeHandler *rest.TestTypeHandler,
	resultHandler *rest.ResultHandler,
	configHandler *rest.ConfigHandler,
	unitHandler *rest.UnitHandler,
) *rest.Handler {
	return &rest.Handler{
		HlSevenHandler:            hlSevenHandler,
		HealthCheckHandler:        healthCheck,
		PatientHandler:            patientHandler,
		SpecimenHandler:           specimenHandler,
		WorkOrderHandler:          workOrder,
		FeatureListHandler:        featureListHandler,
		ObservationRequestHandler: observationRequest,
		TestTypeHandler:           testTypeHandler,
		ResultHandler:             resultHandler,
		ConfigHandler:             configHandler,
		UnitHandler:               unitHandler,
	}
}

const dbFileName = "./tmp/biosystem-lims.db"

// We init once here so it will not create DB twice
var initDBOnce sync.Once
var db *gorm.DB

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func InitSQLiteDB() (*gorm.DB, error) {
	if fileExists(dbFileName) {
		slog.Info("db is existed already")
	} else {
		slog.Info("db is not exists, start create and migrate db")
		err := os.MkdirAll("./tmp", 0755)
		if err != nil {
			return nil, err
		}

		_, err = os.Create(dbFileName)
		if err != nil {
			return nil, err
		}
	}

	dialec, err := sql.Open("sqlite", dbFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db, err := gorm.Open(gormSqlite.Dialector{
		Conn: dialec,
	}, &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.Logger = db.Logger.LogMode(logger.Error)

	return db, nil
}

func InitDatabase() (*gorm.DB, error) {
	db, err := InitSQLiteDB()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	sql, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	err = runMigrations(sql, db.Dialector.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	err = seedTestData(db)
	if err != nil {
		return nil, fmt.Errorf("failed to seed test data: %w", err)
	}

	err = sql.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

// runMigrations performs the database migration.
func runMigrations(db *sql.DB, databaseDriverName string) error {
	slog.Info("🏁 Starting database migration process...")

	sourceDriver, err := iofs.New(migrations.Files, ".")
	if err != nil {
		return fmt.Errorf("failed to create source driver (iofs): %w", err)
	}
	slog.Info("Source driver created (iofs)")

	var dbDriver database.Driver
	switch databaseDriverName {
	case "sqlite":
		dbDriver, err = sqliteMigrate.WithInstance(db, &sqliteMigrate.Config{})
		if err != nil {
			return err
		}
		slog.Info("Database driver created (sqlite)")
	default:
		return errors.New("unsupported database driver for migrations: " + databaseDriverName)
	}

	m, err := migrate.NewWithInstance(
		"iofs",
		sourceDriver,
		databaseDriverName,
		dbDriver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration instance: %w", err)
	}
	slog.Info("Migration instance created")

	slog.Info("Applying migrations...")
	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("No new migrations to apply.")
			return nil
		}
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	slog.Info("Successfully applied new migrations!")
	return nil
}

func seedTestData(db *gorm.DB) error {
	for _, p := range seedPatient {
		err := db.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&p).Error
		if err != nil {
			return err
		}
	}

	for _, config := range seedConfig {
		err := db.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&config).Error
		if err != nil {
			return err
		}
	}

	for _, testType := range seedDataTestType {
		err := db.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&testType).Error
		if err != nil {
			return err
		}
	}

	for _, device := range seedDevice {
		err := db.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(&device).Error
		if err != nil {
			return err
		}
	}

	for _, unit := range seedUnits {
		err := db.Clauses(clause.OnConflict{DoNothing: true}).
			Create(&unit).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func provideDB() *gorm.DB {
	var err error

	initDBOnce.Do(func() {
		db, err = InitDatabase()
		if err != nil {
			panic(err)
		}
	})

	return db
}
func TableKeyValidation(tables entity.Tables) validator.Func {
	return func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if value == "" {
			return true
		}

		_, ok := tables.Find(value)
		if !ok {
			return false
		}

		return true
	}
}

func registerTableValidation(v *validator.Validate) error {
	for key, table := range entity.TableList {
		err := v.RegisterValidation(key, TableKeyValidation(table))
		if err != nil {
			return err
		}
	}

	return nil
}

func provideValidator() *validator.Validate {
	v := validator.New()

	err := registerTableValidation(v)
	if err != nil {
		panic(err)
	}

	return v
}

func provideCache() *cache.Cache {
	return cache.New(time.Hour, 5*time.Minute)
}

func provideConfig(db *gorm.DB) *config.Schema {
	cfg, err := config.New(db)
	if err != nil {
		panic(err)
	}
	return &cfg
}
