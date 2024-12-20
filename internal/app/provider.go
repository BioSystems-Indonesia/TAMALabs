package app

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/gommon/log"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/rest"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/tcp"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/tcp/ba400"
	"github.com/oibacidem/lims-hl-seven/pkg/server"
	"gorm.io/driver/sqlite"
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
) *rest.Handler {
	return &rest.Handler{
		hlSevenHandler,
		healthCheck,
		patientHandler,
		specimenHandler,
		workOrder,
		featureListHandler,
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
	db.Logger = db.Logger.LogMode(logger.Info)

	return db, nil
}

func InitDatabase() (*gorm.DB, error) {
	db, err := InitSQLiteDB()
	if err != nil {
		return nil, err
	}

	autoMigrate := []interface{}{
		&entity.WorkOrderPatient{},
		&entity.WorkOrder{},
		&entity.Patient{},
		&entity.Specimen{},
		&entity.ObservationRequest{},
	}

	for _, model := range autoMigrate {
		log.Infof("AutoMigrate: %T", model)

		err = db.AutoMigrate(model)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
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
