package khanzauc

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/config"
	"github.com/go-mysql-org/go-mysql/canal"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
)

type CanalHandler struct {
	usecase *Usecase
	config  *config.Schema
}

func NewCanalHandler(usecase *Usecase, config *config.Schema) *CanalHandler {
	return &CanalHandler{
		usecase: usecase,
		config:  config,
	}
}

func (h *CanalHandler) OnRow(e *canal.RowsEvent) error {
	if e.Action == canal.InsertAction {
		if e.Table.Name == "lis_order" {
			slog.Info("NEW LIS ORDER DETECTED!", "table", e.Table.Name, "action", e.Action)

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()

			maxAttempts := 3
			for attempt := 1; attempt <= maxAttempts; attempt++ {
				err := h.usecase.SyncAllRequest(ctx)
				if err == nil {
					slog.Info("Successfully synced all requests",
						"table", e.Table.Name,
						"attempt", attempt)
					return nil
				}

				if attempt < maxAttempts {
					slog.Warn("Failed to sync requests, retrying",
						"error", err,
						"attempt", attempt,
						"max_attempts", maxAttempts,
					)
					// Short delay before retry
					time.Sleep(time.Duration(attempt) * time.Second)
					continue
				}

				slog.Error("Failed to sync all requests after all attempts",
					"error", err,
					"table", e.Table.Name,
					"action", e.Action,
					"attempts", maxAttempts,
					"note", "Canal Handler will continue monitoring despite this error",
				)
			}
		}
	}
	return nil
}

// OnTableChanged handles table structure change events
func (h *CanalHandler) OnTableChanged(header *replication.EventHeader, schema string, table string) error {
	return nil
}

// OnDDL handles DDL (Data Definition Language) events
func (h *CanalHandler) OnDDL(header *replication.EventHeader, nextPos mysql.Position, queryEvent *replication.QueryEvent) error {
	return nil
}

// OnGTID handles GTID events
func (h *CanalHandler) OnGTID(header *replication.EventHeader, gtidEvent mysql.BinlogGTIDEvent) error {
	return nil
}

// OnPosSynced handles position sync events
func (h *CanalHandler) OnPosSynced(header *replication.EventHeader, pos mysql.Position, set mysql.GTIDSet, force bool) error {
	return nil
}

// OnRotate handles log rotation events
func (h *CanalHandler) OnRotate(header *replication.EventHeader, rotateEvent *replication.RotateEvent) error {
	return nil
}

// OnRowsQueryEvent handles rows query events
func (h *CanalHandler) OnRowsQueryEvent(event *replication.RowsQueryEvent) error {
	return nil
}

// OnXID handles transaction commit events
func (h *CanalHandler) OnXID(header *replication.EventHeader, nextPos mysql.Position) error {
	return nil
}

// String returns the handler name
func (h *CanalHandler) String() string {
	return "KhanzaCanalHandler"
}

// parseDSNForCanal extracts host, port, user, password, and database from MySQL DSN
func (h *CanalHandler) parseDSNForCanal(dsn string) (host, port, user, password, database string, err error) {
	// Parse DSN format: user:password@tcp(host:port)/database?params
	// Example: "admin:AdminLIS@HL7@tcp(192.168.0.100:3306)/lis_elga_tama?parseTime=true"

	// Find @tcp( part first to correctly identify where credentials end
	tcpIndex := strings.Index(dsn, "@tcp(")
	if tcpIndex == -1 {
		return "", "", "", "", "", fmt.Errorf("invalid DSN format: missing @tcp(")
	}

	userPass := dsn[:tcpIndex]
	remaining := dsn[tcpIndex+5:]

	colonIndex := strings.Index(userPass, ":")
	if colonIndex == -1 {
		return "", "", "", "", "", fmt.Errorf("invalid DSN format: missing user:password separator")
	}
	user = userPass[:colonIndex]
	password = userPass[colonIndex+1:]

	parenIndex := strings.Index(remaining, ")")
	if parenIndex == -1 {
		return "", "", "", "", "", fmt.Errorf("invalid DSN format: missing closing parenthesis")
	}

	hostPort := remaining[:parenIndex]
	remaining = remaining[parenIndex+1:]

	colonIndex = strings.LastIndex(hostPort, ":")
	if colonIndex == -1 {
		host = hostPort
		port = "3306"
	} else {
		host = hostPort[:colonIndex]
		port = hostPort[colonIndex+1:]
	}

	if strings.HasPrefix(remaining, "/") {
		remaining = remaining[1:]
		questionIndex := strings.Index(remaining, "?")
		if questionIndex == -1 {
			database = remaining
		} else {
			database = remaining[:questionIndex]
		}
	}

	slog.Debug("DSN parsing result",
		"original_dsn", dsn,
		"user", user,
		"password_length", len(password),
		"host", host,
		"port", port,
		"database", database,
	)

	return host, port, user, password, database, nil
}

// StartCanalHandler starts the Canal handler for MySQL binlog monitoring
func (h *CanalHandler) StartCanalHandler() {
	if h.config.KhanzaIntegrationEnabled != "true" {
		slog.Info("Khanza integration is disabled, skipping Canal Handler")
		return
	}

	host, port, user, password, database, err := h.parseDSNForCanal(h.config.KhanzaBridgeDatabaseDSN)
	if err != nil {
		slog.Error("Failed to parse Khanza Bridge DSN", "error", err, "dsn", h.config.KhanzaBridgeDatabaseDSN)
		return
	}

	slog.Info("Canal Handler configuration",
		"host", host,
		"port", port,
		"user", user,
		"password_masked", strings.Repeat("*", len(password)),
		"database", database,
	)

	cfg := canal.NewDefaultConfig()
	cfg.Addr = host + ":" + port
	cfg.User = user
	cfg.Password = password
	cfg.Flavor = "mysql"

	slog.Info("Canal configuration details",
		"skip_dump", true,
		"include_table_regex", cfg.IncludeTableRegex,
		"server_id", cfg.ServerID,
	)

	c, err := canal.NewCanal(cfg)
	if err != nil {
		slog.Error("Failed to create Canal instance", "error", err)
		return
	}

	c.SetEventHandler(h)

	slog.Info("Starting Canal Handler for MySQL binlog monitoring",
		"target_database", database,
		"target_table", "lis_order",
	)

	pos, err := c.GetMasterPos()
	if err != nil {
		slog.Error("Failed to get master position", "error", err)
		return
	}

	slog.Info("Starting from current binlog position", "position", pos)

	retryCount := 0
	maxRetries := 10
	baseDelay := 10 * time.Second
	maxDelay := 5 * time.Minute

	for {
		err = c.RunFrom(pos)
		if err != nil {
			retryCount++
			if retryCount >= maxRetries {
				slog.Error("Canal Handler max retries reached, giving up",
					"error", err,
					"retries", retryCount,
					"position", pos,
				)
				return
			}

			delay := time.Duration(retryCount) * baseDelay
			if delay > maxDelay {
				delay = maxDelay
			}

			slog.Error("Canal Handler encountered error, will retry",
				"error", err,
				"position", pos,
				"retry_count", retryCount,
				"delay_seconds", delay.Seconds(),
			)

			time.Sleep(delay)

			newPos, posErr := c.GetMasterPos()
			if posErr == nil {
				pos = newPos
				slog.Info("Updated binlog position for retry", "new_position", pos)
			}
			continue
		}
		retryCount = 0
		break
	}
}
