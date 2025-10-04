package simrs

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

type Repository struct {
	simrsDB *DB
}

func NewRepository(simrsDB *DB) *Repository {
	return &Repository{
		simrsDB: simrsDB,
	}
}

func (s *Repository) RollbackTransaction(ctx context.Context, tx *sql.Tx) {
	if err := tx.Rollback(); err != nil {
		slog.Error("failed to rollback transaction", "error", err)
	}
}

func (s *Repository) GetAllPatients(ctx context.Context) ([]entity.SimrsPatient, error) {
	if s.simrsDB == nil {
		return nil, fmt.Errorf("SIMRS database not initialized")
	}

	conn := s.simrsDB.GetConnection()
	if conn == nil {
		return nil, fmt.Errorf("SIMRS database connection is nil")
	}

	query := `SELECT 
		id,
		patient_id,
		first_name,
		COALESCE(last_name, '') as last_name,
		birthdate,
		gender,
		COALESCE(address, '') as address,
		COALESCE(phone, '') as phone,
		created_at,
		updated_at
	FROM patients 
	WHERE updated_at >= DATE_SUB(NOW(), INTERVAL 24 HOUR)
	ORDER BY updated_at DESC`

	rows, err := conn.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query patients: %w", err)
	}
	defer rows.Close()

	var patients []entity.SimrsPatient
	for rows.Next() {
		var patient entity.SimrsPatient
		err := rows.Scan(
			&patient.ID,
			&patient.PatientID,
			&patient.FirstName,
			&patient.LastName,
			&patient.Birthdate,
			&patient.Gender,
			&patient.Address,
			&patient.Phone,
			&patient.CreatedAt,
			&patient.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan patient row: %w", err)
		}
		patients = append(patients, patient)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating patient rows: %w", err)
	}

	return patients, nil
}

func (s *Repository) GetPatientByID(ctx context.Context, patientID string) (*entity.SimrsPatient, error) {
	if s.simrsDB == nil {
		return nil, fmt.Errorf("SIMRS database not initialized")
	}

	conn := s.simrsDB.GetConnection()
	if conn == nil {
		return nil, fmt.Errorf("SIMRS database connection is nil")
	}

	query := `SELECT 
		id,
		patient_id,
		first_name,
		COALESCE(last_name, '') as last_name,
		birthdate,
		gender,
		COALESCE(address, '') as address,
		COALESCE(phone, '') as phone,
		created_at,
		updated_at
	FROM patients 
	WHERE patient_id = ?`

	row := conn.QueryRowContext(ctx, query, patientID)

	var patient entity.SimrsPatient
	err := row.Scan(
		&patient.ID,
		&patient.PatientID,
		&patient.FirstName,
		&patient.LastName,
		&patient.Birthdate,
		&patient.Gender,
		&patient.Address,
		&patient.Phone,
		&patient.CreatedAt,
		&patient.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("patient with ID %s not found", patientID)
		}
		return nil, fmt.Errorf("failed to get patient by ID: %w", err)
	}

	return &patient, nil
}

func (s *Repository) GetAllLabRequests(ctx context.Context) ([]entity.SimrsLabRequest, error) {
	if s.simrsDB == nil {
		return nil, fmt.Errorf("SIMRS database not initialized")
	}

	conn := s.simrsDB.GetConnection()
	if conn == nil {
		return nil, fmt.Errorf("SIMRS database connection is nil")
	}

	query := `SELECT 
		id,
		no_order,
		patient_id,
		param_request,
		requested_by,
		requested_at,
		created_at
	FROM lab_requests 
	WHERE created_at >= DATE_SUB(NOW(), INTERVAL 24 HOUR)
	ORDER BY created_at DESC`

	rows, err := conn.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query lab requests: %w", err)
	}
	defer rows.Close()

	var labRequests []entity.SimrsLabRequest
	for rows.Next() {
		var labRequest entity.SimrsLabRequest
		err := rows.Scan(
			&labRequest.ID,
			&labRequest.NoOrder,
			&labRequest.PatientID,
			&labRequest.ParamRequest,
			&labRequest.RequestedBy,
			&labRequest.RequestedAt,
			&labRequest.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lab request row: %w", err)
		}
		labRequests = append(labRequests, labRequest)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating lab request rows: %w", err)
	}

	return labRequests, nil
}

func (s *Repository) GetLabRequestByNoOrder(ctx context.Context, noOrder string) (*entity.SimrsLabRequest, error) {
	if s.simrsDB == nil {
		return nil, fmt.Errorf("SIMRS database not initialized")
	}

	conn := s.simrsDB.GetConnection()
	if conn == nil {
		return nil, fmt.Errorf("SIMRS database connection is nil")
	}

	query := `SELECT 
		id,
		no_order,
		patient_id,
		param_request,
		requested_by,
		requested_at,
		created_at
	FROM lab_requests 
	WHERE no_order = ?`

	row := conn.QueryRowContext(ctx, query, noOrder)

	var labRequest entity.SimrsLabRequest
	err := row.Scan(
		&labRequest.ID,
		&labRequest.NoOrder,
		&labRequest.PatientID,
		&labRequest.ParamRequest,
		&labRequest.RequestedBy,
		&labRequest.RequestedAt,
		&labRequest.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("lab request with no_order %s not found", noOrder)
		}
		return nil, fmt.Errorf("failed to get lab request by no_order: %w", err)
	}

	return &labRequest, nil
}

func (s *Repository) BatchInsertLabResults(ctx context.Context, results []entity.SimrsLabResult) error {
	if len(results) == 0 {
		return nil // No records to process
	}

	if s.simrsDB == nil {
		return fmt.Errorf("SIMRS database not initialized")
	}

	conn := s.simrsDB.GetConnection()
	if conn == nil {
		return fmt.Errorf("SIMRS database connection is nil")
	}

	tx, err := conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				slog.Error("failed to rollback transaction", "error", rollbackErr)
			}
		}
	}()

	insertQuery := `
INSERT IGNORE INTO lab_results (
	no_order, param_code, result_value, unit, ref_range, flag, created_at
) VALUES (?, ?, ?, ?, ?, ?, ?);
`

	now := time.Now()

	stmt, err := tx.PrepareContext(ctx, insertQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, result := range results {
		_, err = stmt.ExecContext(ctx,
			result.NoOrder,
			result.ParamCode,
			result.ResultValue,
			result.Unit,
			result.RefRange,
			result.Flag,
			now,
		)
		if err != nil {
			return fmt.Errorf("failed to execute insert for no_order %s, param_code %s: %w", result.NoOrder, result.ParamCode, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *Repository) GetLabResultsByNoOrder(ctx context.Context, noOrder string) ([]entity.SimrsLabResult, error) {
	if s.simrsDB == nil {
		return nil, fmt.Errorf("SIMRS database not initialized")
	}

	conn := s.simrsDB.GetConnection()
	if conn == nil {
		return nil, fmt.Errorf("SIMRS database connection is nil")
	}

	query := `SELECT 
		id,
		no_order,
		param_code,
		result_value,
		unit,
		ref_range,
		flag,
		created_at
	FROM lab_results 
	WHERE no_order = ?
	ORDER BY created_at DESC`

	rows, err := conn.QueryContext(ctx, query, noOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to query lab results: %w", err)
	}
	defer rows.Close()

	var labResults []entity.SimrsLabResult
	for rows.Next() {
		var labResult entity.SimrsLabResult
		err := rows.Scan(
			&labResult.ID,
			&labResult.NoOrder,
			&labResult.ParamCode,
			&labResult.ResultValue,
			&labResult.Unit,
			&labResult.RefRange,
			&labResult.Flag,
			&labResult.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lab result row: %w", err)
		}
		labResults = append(labResults, labResult)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating lab result rows: %w", err)
	}

	return labResults, nil
}

// DeleteLabRequestByNoOrder deletes a lab request by no_order after processing
func (s *Repository) DeleteLabRequestByNoOrder(ctx context.Context, noOrder string) error {
	// Check if SIMRS DB is properly initialized
	if s.simrsDB == nil {
		return fmt.Errorf("SIMRS database not initialized")
	}

	conn := s.simrsDB.GetConnection()
	if conn == nil {
		return fmt.Errorf("SIMRS database connection is nil")
	}

	query := `DELETE FROM lab_requests WHERE no_order = ?`

	result, err := conn.ExecContext(ctx, query, noOrder)
	if err != nil {
		return fmt.Errorf("failed to delete lab request with no_order %s: %w", noOrder, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		slog.Warn("No lab request found to delete", "no_order", noOrder)
	} else {
		slog.Info("Successfully deleted lab request", "no_order", noOrder)
	}

	return nil
}

func (s *Repository) DeleteProcessedLabRequests(ctx context.Context, noOrders []string) error {
	if len(noOrders) == 0 {
		return nil
	}

	if s.simrsDB == nil {
		return fmt.Errorf("SIMRS database not initialized")
	}

	conn := s.simrsDB.GetConnection()
	if conn == nil {
		return fmt.Errorf("SIMRS database connection is nil")
	}

	tx, err := conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				slog.Error("failed to rollback transaction", "error", rollbackErr)
			}
		}
	}()

	deleteQuery := `DELETE FROM lab_requests WHERE no_order = ?`
	stmt, err := tx.PrepareContext(ctx, deleteQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare delete statement: %w", err)
	}
	defer stmt.Close()

	var deletedCount int64
	for _, noOrder := range noOrders {
		result, err := stmt.ExecContext(ctx, noOrder)
		if err != nil {
			return fmt.Errorf("failed to delete lab request with no_order %s: %w", noOrder, err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected for no_order %s: %w", noOrder, err)
		}

		deletedCount += rowsAffected
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit delete transaction: %w", err)
	}

	slog.Info("Successfully deleted processed lab requests", "count", deletedCount, "total_requested", len(noOrders))
	return nil
}

func (s *Repository) DeletePatientByID(ctx context.Context, patientID string) error {
	if s.simrsDB == nil {
		return fmt.Errorf("SIMRS database not initialized")
	}

	conn := s.simrsDB.GetConnection()
	if conn == nil {
		return fmt.Errorf("SIMRS database connection is nil")
	}

	query := `DELETE FROM patients WHERE patient_id = ?`

	result, err := conn.ExecContext(ctx, query, patientID)
	if err != nil {
		return fmt.Errorf("failed to delete patient with ID %s: %w", patientID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		slog.Warn("No patient found to delete", "patient_id", patientID)
	} else {
		slog.Info("Successfully deleted patient", "patient_id", patientID)
	}

	return nil
}

func (s *Repository) DeleteProcessedPatients(ctx context.Context, patientIDs []string) error {
	if len(patientIDs) == 0 {
		return nil // No records to delete
	}

	if s.simrsDB == nil {
		return fmt.Errorf("SIMRS database not initialized")
	}

	conn := s.simrsDB.GetConnection()
	if conn == nil {
		return fmt.Errorf("SIMRS database connection is nil")
	}

	tx, err := conn.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				slog.Error("failed to rollback transaction", "error", rollbackErr)
			}
		}
	}()

	deleteQuery := `DELETE FROM patients WHERE patient_id = ?`
	stmt, err := tx.PrepareContext(ctx, deleteQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare delete statement: %w", err)
	}
	defer stmt.Close()

	var deletedCount int64
	for _, patientID := range patientIDs {
		result, err := stmt.ExecContext(ctx, patientID)
		if err != nil {
			return fmt.Errorf("failed to delete patient with ID %s: %w", patientID, err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected for patient ID %s: %w", patientID, err)
		}

		deletedCount += rowsAffected
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit delete transaction: %w", err)
	}

	slog.Info("Successfully deleted processed patients", "count", deletedCount, "total_requested", len(patientIDs))
	return nil
}
