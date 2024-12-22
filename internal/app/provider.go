package app

import (
	"fmt"
	"os"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/gommon/log"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/rest"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/tcp"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/tcp/ba400"
	"github.com/oibacidem/lims-hl-seven/pkg/server"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func provideTCP(config *config.Schema) *ba400.TCP {
	tcpEr := ba400.NewTCP(config)

	return tcpEr
}

func provideTCPServer(config *config.Schema, handler *tcp.Handler) server.TCPServer {
	serv := server.NewTCP("5678")
	tcp.RegisterRoutes(serv.GetClient(), handler)
	return serv
}

func provideRestServer(config *config.Schema, handlers *rest.Handler, validate *validator.Validate) server.RestServer {
	serv := server.NewRest(config.Port, validate)
	rest.RegisterMiddleware(serv.GetClient())
	rest.RegisterRoutes(serv.GetClient(), handlers)
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
) *rest.Handler {
	return &rest.Handler{
		hlSevenHandler,
		healthCheck,
		patientHandler,
		specimenHandler,
		workOrder,
		featureListHandler,
		observationRequest,
	}
}

func provideTCPHandler(
	HlSevenHHandler *tcp.HlSevenHandler,
) *tcp.Handler {
	return &tcp.Handler{
		HlSevenHandler: HlSevenHHandler,
	}
}

const dbFileName = "./tmp/biosystem-lims.db"

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
func InitSQLiteDB() (*gorm.DB, error) {
	if fileExists(dbFileName) {
		log.Info("db is existed already")
	} else {
		log.Info("db is not exists, start create and migrate db")
		err := os.MkdirAll("./tmp", 0755)
		if err != nil {
			return nil, err
		}

		_, err = os.Create(dbFileName)
		if err != nil {
			return nil, err
		}
	}

	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	db.Logger = db.Logger.LogMode(logger.Silent)

	return db, nil
}

func InitDatabase() (*gorm.DB, error) {
	db, err := InitSQLiteDB()
	if err != nil {
		return nil, err
	}

	autoMigrate := []interface{}{
		&entity.Observation{},
		&entity.ObservationRequest{},
		&entity.ObservationResult{},
		&entity.WorkOrderPatient{},
		&entity.Patient{},
		&entity.Specimen{},
		&entity.WorkOrder{},
		&entity.WorkOrderSpeciment{},
	}

	for _, model := range autoMigrate {
		log.Infof("AutoMigrate: %T", model)

		err = db.AutoMigrate(model)
		if err != nil {
			return nil, err
		}
	}

	err = seedTestData(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func seedTestData(db *gorm.DB) error {
	patient := []entity.Patient{
		{
			FirstName:   "Pasien",
			LastName:    "Pertama",
			Birthdate:   time.Date(1995, time.January, 1, 0, 0, 0, 0, time.UTC),
			Sex:         "M",
			PhoneNumber: "",
			Location:    "",
			Address:     "",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			FirstName:   "Pasien",
			LastName:    "Kedua",
			Birthdate:   time.Date(2002, time.October, 23, 0, 0, 0, 0, time.UTC),
			Sex:         "F",
			PhoneNumber: "",
			Location:    "",
			Address:     "",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			FirstName:   "Pasien",
			LastName:    "Ketiga",
			Birthdate:   time.Date(1998, time.February, 20, 0, 0, 0, 0, time.UTC),
			Sex:         "F",
			PhoneNumber: "",
			Location:    "",
			Address:     "",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	for _, p := range patient {
		err := db.Create(&p).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func provideDB(config *config.Schema) *gorm.DB {
	db, err := InitDatabase()
	if err != nil {
		panic(err)
	}

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
