package app

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/config"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/delivery/cron"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/delivery/rest"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/middleware"
	khanza "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/external/khanza"
	simrs "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/external/simrs"
	licenserepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/license"
	devicerepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/device"
	patientrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/patient"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/test_type"
	workOrderrepo "github.com/BioSystems-Indonesia/TAMALabs/internal/repository/sql/work_order"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/usecase"
	khanzauc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/external/khanza"
	simrsuc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/external/simrs"
	licenseuc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/license"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/result"
	workOrderuc "github.com/BioSystems-Indonesia/TAMALabs/internal/usecase/work_order"
	"github.com/BioSystems-Indonesia/TAMALabs/internal/util"
	"github.com/BioSystems-Indonesia/TAMALabs/migrations"
	"github.com/BioSystems-Indonesia/TAMALabs/pkg/server"
	"github.com/go-playground/validator/v10"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	sqliteMigrate "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	gormSqlite "gorm.io/driver/sqlite"

	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"

	_ "modernc.org/sqlite"
)

func provideAllDevices(deviceRepo *devicerepo.DeviceRepository) []entity.Device {
	allDevices, err := deviceRepo.FindAll(context.Background(), &entity.GetManyRequestDevice{})
	if err != nil {
		slog.Error("failed to get all devices", "error", err)
		panic(fmt.Sprintf("failed to get all devices: %v", err))
	}

	return allDevices.Data
}

func provideRestServer(
	config *config.Schema,
	handlers *rest.Handler,
	validate *validator.Validate,
	deviceHandler *rest.DeviceHandler,
	serverControllerHandler *rest.ServerControllerHandler,
	testTemplateHandler *rest.TestTemplateHandler,
	authHandler *rest.AuthHandler,
	adminHandler *rest.AdminHandler,
	roleHandler *rest.RoleHandler,
	hrisExternal *rest.KhanzaExternalHandler,
	khanzaHandler *rest.ExternalHandler,
	authMiddleware *middleware.JWTMiddleware,
	cronManager *cron.CronManager,
	summaryHandler *rest.SummaryHandler,
) server.RestServer {
	serv := server.NewRest(config.Port, validate, cronManager)
	rest.RegisterMiddleware(serv.GetClient())
	rest.RegisterRoutes(serv.GetClient(), handlers,
		deviceHandler,
		serverControllerHandler,
		testTemplateHandler,
		adminHandler,
		authHandler,
		roleHandler,
		hrisExternal,
		khanzaHandler,
		authMiddleware,
		summaryHandler,
	)

	return serv
}

func provideRestHandler(
	hlSevenHandler *rest.HlSevenHandler,
	healthCheck *rest.HealthCheckHandler,
	healthHandler *rest.HealthHandler,
	patientHandler *rest.PatientHandler,
	specimenHandler *rest.SpecimenHandler,
	workOrder *rest.WorkOrderHandler,
	featureListHandler *rest.FeatureListHandler,
	observationRequest *rest.ObservationRequestHandler,
	testTypeHandler *rest.TestTypeHandler,
	resultHandler *rest.ResultHandler,
	configHandler *rest.ConfigHandler,
	unitHandler *rest.UnitHandler,
	logHandler *rest.LogHandler,
	licenseHandler *rest.LicenseHandler,
) *rest.Handler {
	return &rest.Handler{
		HlSevenHandler:            hlSevenHandler,
		HealthCheckHandler:        healthCheck,
		HealthHandler:             healthHandler,
		PatientHandler:            patientHandler,
		SpecimenHandler:           specimenHandler,
		WorkOrderHandler:          workOrder,
		FeatureListHandler:        featureListHandler,
		ObservationRequestHandler: observationRequest,
		TestTypeHandler:           testTypeHandler,
		ResultHandler:             resultHandler,
		ConfigHandler:             configHandler,
		UnitHandler:               unitHandler,
		LogHandler:                logHandler,
		LicenseHandler:            licenseHandler,
	}
}

// dbFileName holds the path to the SQLite database file. On some environments
// (for example when running as a service) the ProgramData environment
// variable may not be set. In that case we fall back to the conventional
// Windows ProgramData location: C:\\ProgramData
var dbFileName = filepath.Join(os.Getenv("ProgramData"), "TAMALabs", "database", "TAMALabs.db")

func init() {
	// If ProgramData is empty, attempt to use the default Windows location.
	if os.Getenv("ProgramData") == "" && runtime.GOOS == "windows" {
		fallback := `C:\\ProgramData`
		dbFileName = filepath.Join(fallback, "TAMALabs", "database", "TAMALabs.db")
		slog.Info("ProgramData env not set, using fallback path", "dbFile", dbFileName)
	}
}

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
	// Ensure directory exists. Do not use initDBOnce here because
	// provideDB is already guarding initialization with initDBOnce.Do.
	// Calling initDBOnce.Do recursively would deadlock.
	if !fileExists(dbFileName) {
		slog.Info("DB not exists, creating folder...")
		if err := os.MkdirAll(filepath.Dir(dbFileName), 0755); err != nil {
			slog.Error("Failed to create folder", "error", err)
			return nil, err
		}
		slog.Info("Folder created, DB will be created automatically by GORM")
	} else {
		slog.Info("DB already exists")
	}

	dialec, err := sql.Open("sqlite", dbFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	dialec.SetMaxOpenConns(1)

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

	// for _, device := range seedDevice {
	// 	err := db.Clauses(clause.OnConflict{
	// 		DoNothing: true,
	// 	}).Create(&device).Error
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	for _, role := range seedRole {
		err := db.Clauses(clause.OnConflict{UpdateAll: true}).
			Create(&role).Error
		if err != nil {
			return err
		}
	}

	for _, admin := range seedAdmin {
		err := db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&admin).Error
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
	c := cache.New(time.Hour, 5*time.Minute)

	c.OnEvicted(func(key string, value interface{}) {
		slog.Debug("Cache item evicted", "key", key, "value_type", fmt.Sprintf("%T", value))
	})

	return c
}

func provideConfig(db *gorm.DB) *config.Schema {
	cfg, err := config.New(db)
	if err != nil {
		panic(err)
	}
	return &cfg
}

func provideKhanzaRepository(cfg *config.Schema) *khanza.Repository {
	if cfg.KhanzaIntegrationEnabled != "true" {
		return nil
	}

	bridgeDB, err := khanza.NewBridgeDB(cfg)
	if err != nil {
		slog.Error("Error on create khanza db connection. If you want to disable khanza integration, set KhanzaIntegrationEnabled to false on config", "error", err)
		log.Fatalf("failed to create khanza db connection: %v", err)
	}

	mainDB, err := khanza.NewMainDB(cfg)
	if err != nil {
		slog.Error("Error on create khanza db connection. If you want to disable khanza integration, set KhanzaIntegrationEnabled to false on config", "error", err)
		log.Fatalf("failed to create khanza db connection: %v", err)
	}

	return khanza.NewRepository(bridgeDB, mainDB)
}

func provideSimrsRepository(cfg *config.Schema) *simrs.Repository {
	if cfg.SimrsIntegrationEnabled != "true" {
		return nil
	}

	simrsDB, err := simrs.NewDB(cfg.SimrsDatabaseDSN)
	if err != nil {
		slog.Error("Error on create SIMRS db connection. If you want to disable SIMRS integration, set SimrsIntegrationEnabled to false on config", "error", err)
		slog.Info("failed to create SIMRS db connection, SIMRS integration will be disabled", "error", err)
		return nil
	}

	return simrs.NewRepository(simrsDB)
}

func provideSimrsUsecase(
	cfg *config.Schema,
	simrsRepo *simrs.Repository,
	workOrderRepo *workOrderrepo.WorkOrderRepository,
	workOrderUC *workOrderuc.WorkOrderUseCase,
	patientRepo *patientrepo.PatientRepository,
	testTypeRepo *test_type.Repository,
	resultUC *result.Usecase,
) *simrsuc.Usecase {
	if cfg.SimrsIntegrationEnabled != "true" {
		slog.Info("SIMRS integration is disabled, SIMRS Usecase will not be created")
		return nil
	}

	if simrsRepo == nil {
		slog.Info("SIMRS repository is nil (connection failed), SIMRS Usecase will not be created")
		return nil
	}

	return simrsuc.NewUsecase(
		simrsRepo,
		workOrderRepo,
		workOrderUC,
		patientRepo,
		testTypeRepo,
		cfg,
		resultUC,
	)
}

func provideCanalHandler(
	cfg *config.Schema,
	khanzaRepo *khanza.Repository,
	workOrderRepo *workOrderrepo.WorkOrderRepository,
	patientRepo *patientrepo.PatientRepository,
	testTypeRepo *test_type.Repository,
	barcodeUC usecase.BarcodeGenerator,
	resultUC *result.Usecase,
) *khanzauc.CanalHandler {
	if cfg.KhanzaIntegrationEnabled != "true" {
		slog.Info("Khanza integration is disabled, Canal Handler will not be created")
		return nil
	}

	khanzaUC := khanzauc.NewUsecase(
		khanzaRepo,
		workOrderRepo,
		patientRepo,
		testTypeRepo,
		barcodeUC,
		resultUC,
	)

	slog.Info("Creating Canal Handler with fully configured dependencies")
	return khanzauc.NewCanalHandler(khanzaUC, cfg)
}

func provideLicenseService() *licenseuc.License {
	pubLoader := licenserepo.NewFSKeyLoader()
	fileLoader := licenserepo.NewFSFileLoader()

	// Ensure license directory exists
	licenseDir := "license"
	if err := os.MkdirAll(licenseDir, 0755); err != nil {
		slog.Warn("Failed to create license directory", "error", err)
	}

	// default paths (relative to working directory)
	pubKeyPath := "license/server_public.pem"
	licensePath := "license/license.json"

	lic := licenseuc.NewLicense(pubLoader, fileLoader, pubKeyPath, licensePath)

	// Start a background heartbeat goroutine to check license with license server.
	go func() {
		licenseServerURL := os.Getenv("LICENSE_SERVER_URL")
		if licenseServerURL == "" {
			licenseServerURL = "http://localhost:8080"
		}

		machineID, err := util.GenerateMachineID()
		if err != nil {
			slog.Error("Failed to generate machine ID for heartbeat", "error", err)
			machineID = "unknown"
		}

		client := &http.Client{Timeout: 30 * time.Second}
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		// Run immediately once
		for {
			// Read license file
			data, err := os.ReadFile(licensePath)
			if err != nil {
				slog.Debug("License file not found, skipping heartbeat", "error", err)
			} else {
				var m map[string]interface{}
				if err := json.Unmarshal(data, &m); err != nil {
					slog.Warn("Failed to decode license file, skipping heartbeat", "error", err)
				} else {
					codeRaw, ok := m["license_code"]
					if ok {
						codeStr, ok := codeRaw.(string)
						if ok {
							// prepare request
							hb := map[string]string{"machine_id": machineID, "license_code": codeStr}
							jb, _ := json.Marshal(hb)
							req, err := http.NewRequest("POST", fmt.Sprintf("%s/heartbeat", licenseServerURL), bytes.NewBuffer(jb))
							if err != nil {
								slog.Warn("License heartbeat request create failed", "error", err)
							} else {
								req.Header.Set("Content-Type", "application/json")
								apiKey := "KJKDANCJSANIUWYR6243UJFOISJFJKVOMV72487YEHFHFHSDVOHF9AMDC9AN9SDN98YE98YEHDIU2Y897873YYY68686487WGDUDUAGYTE8QTEYADIUHADUYW8E8BWTNC8N8NAMDOAIMDAUDUWYAD87NYW7Y7CBT87EY8142164B36248732M87MCIFH8NYRWCM8MYCMUOIDOIADOIDOIUR83YR983Y98328N32C83NYC8732NYC8732Y87Y32NCNSAIHJAOJFOIJFOIQFIUIUNCNHCIUHWV8NRYNV8Y989N9198298YOIJOI090103021313JKJDHAHDJAJASHHAH"
								if apiKey != "" {
									req.Header.Set("X-API-Key", apiKey)
								}
								resp, err := client.Do(req)
								if err != nil {
									slog.Warn("License heartbeat failed - server unreachable", "error", err, "server", licenseServerURL)
								} else {
									body, _ := io.ReadAll(resp.Body)
									_ = resp.Body.Close()
									bodyStr := strings.TrimSpace(string(body))
									bodyStr = strings.Trim(bodyStr, `"`)
									slog.Info("License heartbeat response", "msg", bodyStr)

									// helper to revoke license locally
									revoke := func(reason string) {
										slog.Error("License revoked by server, removing local license files", "reason", reason)
										runtime.GC()
										time.Sleep(100 * time.Millisecond)
										_ = os.Remove(licensePath)
										_ = os.Remove(pubKeyPath)
										rev := map[string]interface{}{"revoked_at": time.Now().Unix(), "reason": reason}
										if rb, err := json.MarshalIndent(rev, "", "  "); err == nil {
											_ = os.WriteFile("license/revoked.json", rb, 0644)
										}
									}

									// Try to parse as structured JSON: {code, status, message}
									var hr struct {
										Code    int    `json:"code"`
										Status  string `json:"status"`
										Message string `json:"message"`
									}

									parsed := false
									if err := json.Unmarshal(body, &hr); err == nil {
										parsed = true
									} else {
										// maybe the server returned a JSON string that contains JSON (quoted JSON)
										var inner string
										if err2 := json.Unmarshal(body, &inner); err2 == nil {
											// try to parse inner JSON
											if err3 := json.Unmarshal([]byte(inner), &hr); err3 == nil {
												parsed = true
											} else {
												// inner not structured JSON we expect
											}
										}
									}

									if parsed {
										lowerMsg := strings.ToLower(hr.Message)
										if hr.Code >= 400 && (strings.Contains(lowerMsg, "device not found") || strings.Contains(lowerMsg, "revoked") || strings.Contains(lowerMsg, "mismatch")) {
											revoke(hr.Message)
										}
									} else {
										// fallback to older string matching on bodyStr
										switch strings.ToLower(bodyStr) {
										case "device not found", "license mismatch", "device revoked":
											revoke(bodyStr)
										default:
											// OK or other non-critical responses
										}
									}
								}
							}
						}
					}
				}
			}

			// Wait for next tick
			<-ticker.C
		}
	}()

	return lic
}
