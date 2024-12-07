package app

import (
	"fmt"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/oibacidem/lims-hl-seven/config"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/rest"
	"github.com/oibacidem/lims-hl-seven/internal/delivery/tcp"
	"github.com/oibacidem/lims-hl-seven/internal/entity"
	"github.com/oibacidem/lims-hl-seven/internal/repository/tcp/hl_seven"
	"github.com/oibacidem/lims-hl-seven/pkg/server"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func provideTCP(config *config.Schema) *hl_seven.TCP {
	tcpEr := hl_seven.NewTCP(config)

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
) *rest.Handler {
	return &rest.Handler{
		hlSevenHandler,
		healthCheck,
		patientHandler,
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
		fmt.Println("db is existed already")
	} else {
		fmt.Println("db is not exists, start create and migrate db")
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
	return db, nil
}

func InitDatabase() (*gorm.DB, error) {
	db, err := InitSQLiteDB()
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&entity.Patient{})
	if err != nil {
		return nil, err
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

func provideValidator() *validator.Validate {
	return validator.New()
}
