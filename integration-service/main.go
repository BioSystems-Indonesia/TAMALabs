package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/BioSystems-Indonesia/integration-services-lis/internal/database"
	"github.com/BioSystems-Indonesia/integration-services-lis/internal/helper"
	repositories "github.com/BioSystems-Indonesia/integration-services-lis/internal/repository"
	services "github.com/BioSystems-Indonesia/integration-services-lis/internal/service"
)

const (
	BaseURL                       = "tamalabs.biosystems.id"
	ApiKey                        = "KJKDANCJSANIUWYR6243UJFOISJFJKVOMV72487YEHFHFHSDVOHF9AMDC9AN9SDN98YE98YEHDIU2Y897873YYY68686487WGDUDUAGYTE8QTEYADIUHADUYW8E8BWTNC8N8NAMDOAIMDAUDUWYAD87NYW7Y7CBT87EY8142164B36248732M87MCIFH8NYRWCM8MYCMUOIDOIADOIDOIUR83YR983Y98328N32C83NYC8732NYC8732Y87Y32NCNSAIHJAOJFOIJFOIQFIUIUNCNHCIUHWV8NRYNV8Y989N9198298YOIJOI090103021313JKJDHAHDJAJASHHAH"
	HTTPPort                      = ":8214"
	UserSyncInterval              = 1 * time.Minute
	TestTypeSyncInterval          = 30 * time.Second
	ObservationResultSyncInterval = 10 * time.Second
	SyncContextTimeout            = 2 * time.Minute
	SyncQueueSize                 = 10
)

var (
	KeysDir    string
	LabKeyPath string
	LogDir     string
)

func init() {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		localAppData = "."
	}

	baseDir := filepath.Join(localAppData, "TAMALabs")
	KeysDir = filepath.Join(baseDir, "keys")
	LabKeyPath = filepath.Join(KeysDir, "lab_key.json")
	LogDir = "./logs"
}

type Config struct {
	BaseURL string
	ApiKey  string
	KeysDir string
}

type Services struct {
	User              *services.UserSyncService
	TestType          *services.SyncTestType
	ObservationResult *services.ObservationResultSyncService
}

type DailyLogWriter struct {
	mu       sync.Mutex
	file     *os.File
	lastDate string
	logDir   string
}

func NewDailyLogWriter(logDir string) (*DailyLogWriter, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	dlw := &DailyLogWriter{
		logDir:   logDir,
		lastDate: "",
	}

	if err := dlw.rotateIfNeeded(); err != nil {
		return nil, err
	}

	return dlw, nil
}

func (dlw *DailyLogWriter) rotateIfNeeded() error {
	currentDate := time.Now().Format("2006-01-02")

	if dlw.lastDate != currentDate {
		if dlw.file != nil {
			dlw.file.Close()
		}

		logFile, err := os.OpenFile(
			fmt.Sprintf("%s/Integration_Service_%s.log", dlw.logDir, currentDate),
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			0666,
		)
		if err != nil {
			return err
		}

		dlw.file = logFile
		dlw.lastDate = currentDate
	}

	return nil
}

func (dlw *DailyLogWriter) Write(p []byte) (n int, err error) {
	dlw.mu.Lock()
	defer dlw.mu.Unlock()

	if err := dlw.rotateIfNeeded(); err != nil {
		return 0, err
	}

	return dlw.file.Write(p)
}

func (dlw *DailyLogWriter) Close() error {
	dlw.mu.Lock()
	defer dlw.mu.Unlock()

	if dlw.file != nil {
		return dlw.file.Close()
	}

	return nil
}

type SyncScheduler struct {
	queue                     chan func()
	wg                        sync.WaitGroup
	stop                      chan os.Signal
	tickers                   map[string]*time.Ticker
	enqueue                   func(string, func())
	services                  *Services
	baseURL                   string
	apiKey                    string
	tickersMutex              sync.Mutex
	servicesActive            bool
	consecutiveHealthFailures int
}

func newConfig() *Config {
	return &Config{
		BaseURL: BaseURL,
		ApiKey:  ApiKey,
		KeysDir: KeysDir,
	}
}

func main() {
	dlw, err := NewDailyLogWriter(LogDir)
	if err != nil {
		panic(fmt.Sprintf("Failed to create log writer: %v", err))
	}
	defer dlw.Close()

	writers := []io.Writer{dlw}
	if stdoutIsAttached() {
		writers = append(writers, os.Stdout)
	}
	log.SetOutput(io.MultiWriter(writers...))
	log.SetFlags(log.LstdFlags)

	log.Println("🚀 LIS Integration Service started")

	cfg := newConfig()

	log.Println("🔌 Checking server connectivity...")
	if err := healthCheckWithRetry(cfg.BaseURL, cfg.ApiKey); err != nil {
		log.Fatal(err)
	}

	initializeKeyManagement(cfg)
	database.Connect()

	svcs := initializeServices(cfg)
	scheduler := runScheduler(svcs, cfg)

	go scheduler.startHTTPServer()
	scheduler.start()
}

func initializeKeyManagement(cfg *Config) {
	repoKey := repositories.NewFileSystemKeyPairRepository(cfg.KeysDir)
	serviceKey := services.NewKeyPairService(repoKey)

	labService := services.NewLabSyncService(
		fmt.Sprintf("https://%s/unauthenticated/lab-public/", cfg.BaseURL),
		cfg.ApiKey,
		LabKeyPath,
	)

	if !fileExists(LabKeyPath) {
		log.Println("🔑 Registering lab...")
		if err := retryWithBackoff(3, time.Second, func() error {
			return labService.CreateLab()
		}); err != nil {
			log.Printf("⚠️ Lab registration failed: %v, continuing with startup\n", err)
		}
	}

	log.Printf("🔍 Loading lab key from: %s\n", LabKeyPath)
	labKey, err := helper.LoadLabKey(LabKeyPath)
	if err != nil {
		log.Fatalf("❌ Failed to load lab key: %v\n", err)
	}
	log.Printf("✅ Lab key loaded: KeyId=%s, LabId=%s\n", labKey.KeyId, labKey.LabId)

	log.Printf("🔑 Ensuring key pair for KeyId: %s\n", labKey.KeyId)
	if _, err := serviceKey.EnsureKeyPair(labKey.KeyId); err != nil {
		log.Fatalf("❌ Failed to ensure key pair: %v\n", err)
	}
	log.Println("✅ Key pair ensured")

	keyDir := filepath.Join(cfg.KeysDir, labKey.KeyId)
	log.Printf("📤 Sending public key from: %s\n", keyDir)
	if err := retryWithBackoff(3, time.Second, func() error {
		err := serviceKey.SendPublicKey(
			cfg.BaseURL,
			keyDir,
			labKey.KeyId,
			cfg.ApiKey,
		)
		if err != nil {
			log.Printf("❌ SendPublicKey attempt failed: %v\n", err)
		}
		return err
	}); err != nil {
		log.Printf("⚠️ Sending public key failed after retries: %v, continuing with startup\n", err)
	} else {
		log.Println("✅ Public key sent successfully")
	}
}

func initializeServices(cfg *Config) *Services {
	userRepo := repositories.NewUserRepository(database.DB)
	userService := services.NewUserSyncService(
		userRepo,
		fmt.Sprintf("https://%s/protected/user/register/", cfg.BaseURL),
		LabKeyPath,
		cfg.KeysDir,
	)

	testTypeRepo := repositories.NewTestTypeRepository(database.DB)
	testTypeService := services.NewSyncTestType(
		testTypeRepo,
		fmt.Sprintf("https://%s/protected/test-type/", cfg.BaseURL),
		LabKeyPath,
		cfg.KeysDir,
	)

	observationResultRepo := repositories.NewObservationResultRepository(database.DB)
	observationResultService := services.NewObservationResultSyncService(
		observationResultRepo,
		userRepo,
		cfg.BaseURL,
		LabKeyPath,
		cfg.KeysDir,
	)

	return &Services{
		User:              userService,
		TestType:          testTypeService,
		ObservationResult: observationResultService,
	}
}

func runScheduler(svcs *Services, cfg *Config) *SyncScheduler {
	queue := make(chan func(), SyncQueueSize)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	scheduler := &SyncScheduler{
		queue:                     queue,
		stop:                      stop,
		tickers:                   make(map[string]*time.Ticker),
		services:                  svcs,
		baseURL:                   cfg.BaseURL,
		apiKey:                    cfg.ApiKey,
		servicesActive:            true,
		consecutiveHealthFailures: 0,
	}

	scheduler.enqueue = func(name string, job func()) {
		scheduler.tickersMutex.Lock()
		isActive := scheduler.servicesActive
		scheduler.tickersMutex.Unlock()

		if !isActive {
			return
		}

		scheduler.wg.Add(1)
		log.Println("📥 Enqueue sync:", name)
		queue <- job
	}

	go scheduler.runWorker()
	go scheduler.runVerification()
	go scheduler.healthMonitor()

	return scheduler
}

func (s *SyncScheduler) runWorker() {
	for job := range s.queue {
		job()
		s.wg.Done()
	}
	log.Println("🛑 Sync worker stopped")
}

func (s *SyncScheduler) runVerification() {
	for {
		s.tickersMutex.Lock()
		isActive := s.servicesActive
		s.tickersMutex.Unlock()

		if isActive {
			err := s.services.ObservationResult.VerifyObservationResult(context.Background(), ApiKey)
			if err != nil {
				log.Printf("❌ Observation result verification failed: %v\n", err)
			}
		}

		select {
		case <-time.After(5 * time.Second):
		case <-s.stop:
			return
		}
	}
}

func (s *SyncScheduler) start() {
	s.initializeTickers()

	go s.handleShutdown()
	log.Println("⏰ Schedulers started")

	s.runScheduleLoop()
	s.cleanup()
}

func (s *SyncScheduler) initializeTickers() {
	s.tickers["user"] = time.NewTicker(UserSyncInterval)
	s.tickers["testType"] = time.NewTicker(TestTypeSyncInterval)
	s.tickers["observationResult"] = time.NewTicker(ObservationResultSyncInterval)
}

func (s *SyncScheduler) handleShutdown() {
	<-s.stop
	log.Println("⚠️ Shutdown signal received")

	for _, ticker := range s.tickers {
		ticker.Stop()
	}

	close(s.queue)
}

func (s *SyncScheduler) runScheduleLoop() {
	for {
		s.tickersMutex.Lock()
		isActive := s.servicesActive
		s.tickersMutex.Unlock()

		if !isActive {
			time.Sleep(time.Second)
			continue
		}

		select {
		case <-s.tickers["user"].C:
			s.enqueue("user", func() {
				executeSync("user", func(ctx context.Context) error {
					return s.services.User.SyncUser(ctx)
				})
			})

		case <-s.tickers["testType"].C:
			s.enqueue("test-type", func() {
				executeSync("test-type", func(ctx context.Context) error {
					return s.services.TestType.SyncTestType(ctx)
				})
			})

		case <-s.tickers["observationResult"].C:
			s.enqueue("observation-result", func() {
				executeSync("observation-result", func(ctx context.Context) error {
					return s.services.ObservationResult.SyncObservationResult(ctx)
				})
			})

		case <-s.stop:
			log.Println("🛑 Stop scheduler loop")
			return

		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (s *SyncScheduler) cleanup() {
	log.Println("⏳ Waiting for running sync to finish...")
	s.wg.Wait()

	if sqlDB, err := database.DB.DB(); err == nil {
		sqlDB.Close()
	}

	log.Println("✅ Shutdown complete")
}

func executeSync(name string, fn func(ctx context.Context) error) {
	ctx, cancel := context.WithTimeout(context.Background(), SyncContextTimeout)
	defer cancel()

	log.Println("▶️ Start sync:", name)

	if err := executeWithRetry(ctx, name, fn); err != nil {
		log.Printf("❌ Sync %s failed: %v\n", name, err)
		return
	}

	log.Println("✅ Sync success:", name)
}

func executeWithRetry(ctx context.Context, name string, fn func(ctx context.Context) error) error {
	maxRetries := 10
	baseDelay := time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		err := fn(ctx)
		if err == nil {
			return nil
		}

		if attempt < maxRetries-1 {
			log.Printf("⚠️ Sync %s failed: %v, retrying in %v (attempt %d/%d)\n", name, err, baseDelay, attempt+1, maxRetries)
			select {
			case <-time.After(baseDelay):
			case <-ctx.Done():
				return ctx.Err()
			}
		} else {
			log.Printf("❌ Sync %s final attempt failed: %v\n", name, err)
		}
	}

	return fmt.Errorf("sync %s failed after %d retries", name, maxRetries)
}

func (s *SyncScheduler) startHTTPServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		s.tickersMutex.Lock()
		isActive := s.servicesActive
		s.tickersMutex.Unlock()

		if !isActive {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte("server connection lost"))
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})
	mux.HandleFunc("/generate_result_public", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var request struct {
			Barcode string `json:"barcode"`
		}

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			log.Printf("❌ Failed to parse request body: %v\n", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if request.Barcode == "" {
			http.Error(w, "Barcode is required", http.StatusBadRequest)
			return
		}

		result, err := s.services.ObservationResult.GeneratePublicLink(r.Context(), s.apiKey, request.Barcode)
		if err != nil {
			log.Printf("❌ Failed to generate public link: %v\n", err)
			http.Error(w, fmt.Sprintf("Failed to generate public link: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(result.Data))
	})

	handler := corsMiddleware(mux)

	log.Printf("🌐 HTTP server listening on %s\n", HTTPPort)
	if err := http.ListenAndServe(HTTPPort, handler); err != nil {
		log.Printf("❌ HTTP server error: %v\n", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func stdoutIsAttached() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func retryWithBackoff(maxRetries int, baseDelay time.Duration, fn func() error) error {
	for attempt := 0; attempt < maxRetries; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		if attempt < maxRetries-1 {
			log.Printf("⚠️ Retrying in %v (attempt %d/%d): %v\n", baseDelay, attempt+1, maxRetries, err)
			time.Sleep(baseDelay)
		} else {
			log.Printf("❌ Final attempt failed: %v\n", err)
		}
	}

	return fmt.Errorf("operation failed after %d retries", maxRetries)
}

func healthCheck(baseURL, apiKey string) error {
	url := fmt.Sprintf("https://%s/ping", baseURL)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-API-Key", apiKey)
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("server returned status %d", resp.StatusCode)
	}

	return nil
}

func healthCheckWithRetry(baseURL, apiKey string) error {
	log.Println("🔍 Attempting to connect to server...")
	retryDelay := 5 * time.Second
	for attempt := 0; ; attempt++ {
		if err := healthCheck(baseURL, apiKey); err == nil {
			log.Println("✅ Server is reachable")
			return nil
		}

		log.Printf("⚠️  Server unreachable, retrying in %v (attempt %d)\n", retryDelay, attempt+1)
		time.Sleep(retryDelay)
	}
}

func (s *SyncScheduler) healthMonitor() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := healthCheck(s.baseURL, s.apiKey); err != nil {
				s.tickersMutex.Lock()
				s.consecutiveHealthFailures++

				if s.servicesActive {
					s.servicesActive = false
					log.Printf("⚠️  Server connection lost, pausing services: %v\n", err)
				}
				s.tickersMutex.Unlock()
			} else {
				s.tickersMutex.Lock()
				s.consecutiveHealthFailures = 0
				if !s.servicesActive {
					s.servicesActive = true
					log.Println("✅ Server connection restored, resuming services")
				}
				s.tickersMutex.Unlock()
			}

		case <-s.stop:
			return
		}
	}
}
