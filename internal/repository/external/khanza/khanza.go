package khanza

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

// Repository is a repository for external applications using SIM RS Repository
type Repository struct {
	db *DB
}

// NewRepository creates a new Khanza instance
func NewRepository(db *DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (s *Repository) RollbackTransaction(ctx context.Context, tx *sql.Tx) {
	if err := tx.Rollback(); err != nil {
		slog.Error("failed to rollback transaction", "error", err)
	}
}

// BatchUpsertRESDTO performs batch upsert operations for multiple RESDT records in a single transaction
func (s *Repository) BatchUpsertRESDTO(ctx context.Context, reqs []entity.KhanzaResDT) (err error) {
	tx, err := s.db.GetConnection().BeginTx(ctx, nil)
	if err != nil {
		s.RollbackTransaction(ctx, tx)

		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	var (
		values []string
		args   []interface{}
	)

	now := time.Now()
	values = make([]string, 0, len(reqs))
	for _ = range reqs {
		values = append(values, "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	}

	for _, resdt := range reqs {
		args = append(args, resdt.ONO, resdt.OrderTestID, resdt.TESTCD, resdt.TestNM, resdt.DataTyp,
			resdt.ResultValue, resdt.ResultFT, resdt.Unit, resdt.Flag, resdt.RefRange, now)
	}

	insertQuery := fmt.Sprintf(`INSERT INTO RESDT (
			ONO, ORDER_TESTID, TEST_CD, TEST_NM, DATA_TYP, RESULT_VALUE, RESULT_FT, UNIT, FLAG, REF_RANGE, VALIDATE_ON
		) VALUES %s
		ON DUPLICATE KEY UPDATE
			ORDER_TESTID = VALUES(ORDER_TESTID),
			DATA_TYP = VALUES(DATA_TYP),
			RESULT_VALUE = VALUES(RESULT_VALUE),
			RESULT_FT = VALUES(RESULT_FT),
			UNIT = VALUES(UNIT),
			FLAG = VALUES(FLAG),
			REF_RANGE = VALUES(REF_RANGE),
			VALIDATE_ON = VALUES(VALIDATE_ON)`, strings.Join(values, ", "))

	slog.Info("insert query", "query", insertQuery)

	_, err = tx.ExecContext(ctx, insertQuery, args...)
	if err != nil {
		s.RollbackTransaction(ctx, tx)
		return fmt.Errorf("failed to execute batch insert: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		s.RollbackTransaction(ctx, tx)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetLisOrderByONO retrieves a LisOrder record by ONO
func (k *Repository) GetLisOrderByONO(ono string) (*entity.KhanzaLisOrder, error) {
	query := `SELECT 
		ID,
		MESSAGE_DT,
		ORDER_CONTROL,
		PID,
		PNAME,
		ADDRESS1,
		ADDRESS2,
		ADDRESS3,
		ADDRESS4,
		PTYPE,
		BIRTH_DT,
		SEX,
		ONO,
		REQUEST_DT,
		SOURCE,
		CLINICIAN,
		ROOM_NO,
		PRIORITY,
		COMMENT,
		VISITNO,
		ORDER_TESTID,
		FLAG
	FROM lis_order 
	WHERE ONO = ?`

	conn := k.db.GetConnection()
	row := conn.QueryRow(query, ono)

	var order entity.KhanzaLisOrder
	err := row.Scan(
		&order.ID,
		&order.MessageDT,
		&order.OrderControl,
		&order.PID,
		&order.PName,
		&order.Address1,
		&order.Address2,
		&order.Address3,
		&order.Address4,
		&order.PType,
		&order.BirthDT,
		&order.Sex,
		&order.ONO,
		&order.RequestDT,
		&order.Source,
		&order.Clinician,
		&order.RoomNo,
		&order.Priority,
		&order.Comment,
		&order.VisitNo,
		&order.OrderTestID,
		&order.Flag,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get LisOrder by ONO: %w", err)
	}

	return &order, nil
}

// GetAllLisOrders retrieves all LisOrder records
func (k *Repository) GetAllLisOrders() ([]entity.KhanzaLisOrder, error) {
	query := `SELECT 
		ID,
		MESSAGE_DT,
		ORDER_CONTROL,
		PID,
		PNAME,
		ADDRESS1,
		ADDRESS2,
		ADDRESS3,
		ADDRESS4,
		PTYPE,
		BIRTH_DT,
		SEX,
		ONO,
		REQUEST_DT,
		SOURCE,
		CLINICIAN,
		ROOM_NO,
		PRIORITY,
		COMMENT,
		VISITNO,
		ORDER_TESTID,
		FLAG
	FROM lis_order`

	conn := k.db.GetConnection()
	rows, err := conn.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query LisOrders: %w", err)
	}
	defer rows.Close()

	var orders []entity.KhanzaLisOrder
	for rows.Next() {
		var order entity.KhanzaLisOrder
		err := rows.Scan(
			&order.ID,
			&order.MessageDT,
			&order.OrderControl,
			&order.PID,
			&order.PName,
			&order.Address1,
			&order.Address2,
			&order.Address3,
			&order.Address4,
			&order.PType,
			&order.BirthDT,
			&order.Sex,
			&order.ONO,
			&order.RequestDT,
			&order.Source,
			&order.Clinician,
			&order.RoomNo,
			&order.Priority,
			&order.Comment,
			&order.VisitNo,
			&order.OrderTestID,
			&order.Flag,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan LisOrder row: %w", err)
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating LisOrder rows: %w", err)
	}

	return orders, nil
}
