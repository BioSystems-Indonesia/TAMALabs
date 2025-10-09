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
		ISNULL(last_name, '') as last_name,
		birthdate,
		gender,
		ISNULL(address, '') as address,
		ISNULL(phone, '') as phone,
		created_at,
		updated_at
	FROM patient 
	WHERE updated_at >= DATEADD(HOUR, -24, GETDATE())
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
		ISNULL(last_name, '') as last_name,
		birthdate,
		gender,
		ISNULL(address, '') as address,
		ISNULL(phone, '') as phone,
		created_at,
		updated_at
	FROM patient 
	WHERE patient_id = @patientID`

	row := conn.QueryRowContext(ctx, query, sql.Named("patientID", patientID))

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
	WHERE created_at >= DATEADD(HOUR, -24, GETDATE())
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
	WHERE no_order = @noOrder`

	row := conn.QueryRowContext(ctx, query, sql.Named("noOrder", noOrder))

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

	// First, check which records already exist
	checkQuery := `
	SELECT no_order, param_code 
	FROM lab_results 
	WHERE no_order = @noOrder AND param_code = @paramCode
	`

	insertQuery := `
	INSERT INTO lab_results (
		no_order, param_code, result_value, unit, ref_range, flag, created_at
	) VALUES (
		@noOrder, @paramCode, @resultValue, @unit, @refRange, @flag, @createdAt
	)
	`

	now := time.Now()
	var insertedCount int
	var skippedCount int

	checkStmt, err := tx.PrepareContext(ctx, checkQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare check statement: %w", err)
	}
	defer checkStmt.Close()

	insertStmt, err := tx.PrepareContext(ctx, insertQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer insertStmt.Close()

	for _, result := range results {
		// Check if record already exists
		row := checkStmt.QueryRowContext(ctx,
			sql.Named("noOrder", result.NoOrder),
			sql.Named("paramCode", result.ParamCode),
		)

		var existingNoOrder, existingParamCode string
		err := row.Scan(&existingNoOrder, &existingParamCode)

		if err == sql.ErrNoRows {
			// Record doesn't exist, insert it
			_, err = insertStmt.ExecContext(ctx,
				sql.Named("noOrder", result.NoOrder),
				sql.Named("paramCode", result.ParamCode),
				sql.Named("resultValue", result.ResultValue),
				sql.Named("unit", result.Unit),
				sql.Named("refRange", result.RefRange),
				sql.Named("flag", result.Flag),
				sql.Named("createdAt", now),
			)
			if err != nil {
				return fmt.Errorf("failed to execute insert for no_order %s, param_code %s: %w", result.NoOrder, result.ParamCode, err)
			}
			insertedCount++
			slog.Info("Inserted new lab result", "no_order", result.NoOrder, "param_code", result.ParamCode)
		} else if err != nil {
			return fmt.Errorf("failed to check existing record for no_order %s, param_code %s: %w", result.NoOrder, result.ParamCode, err)
		} else {
			// Record already exists, skip it
			skippedCount++
			slog.Info("Skipped existing lab result", "no_order", result.NoOrder, "param_code", result.ParamCode)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	slog.Info("Batch insert lab results completed",
		"total_processed", len(results),
		"inserted", insertedCount,
		"skipped", skippedCount)

	return nil
}

// CheckIfLabResultExists checks if a lab result already exists in the database
func (s *Repository) CheckIfLabResultExists(ctx context.Context, noOrder string, paramCode string) (bool, error) {
	if s.simrsDB == nil {
		return false, fmt.Errorf("SIMRS database not initialized")
	}

	conn := s.simrsDB.GetConnection()
	if conn == nil {
		return false, fmt.Errorf("SIMRS database connection is nil")
	}

	query := `
	SELECT COUNT(1) 
	FROM lab_results 
	WHERE no_order = @noOrder AND param_code = @paramCode
	`

	var count int
	err := conn.QueryRowContext(ctx, query,
		sql.Named("noOrder", noOrder),
		sql.Named("paramCode", paramCode),
	).Scan(&count)

	if err != nil {
		return false, fmt.Errorf("failed to check if lab result exists: %w", err)
	}

	return count > 0, nil
}

// FilterNewLabResults filters out lab results that already exist in the database
func (s *Repository) FilterNewLabResults(ctx context.Context, results []entity.SimrsLabResult) ([]entity.SimrsLabResult, error) {
	if len(results) == 0 {
		return results, nil
	}

	var newResults []entity.SimrsLabResult

	for _, result := range results {
		exists, err := s.CheckIfLabResultExists(ctx, result.NoOrder, result.ParamCode)
		if err != nil {
			return nil, fmt.Errorf("failed to check existence for no_order %s, param_code %s: %w",
				result.NoOrder, result.ParamCode, err)
		}

		if !exists {
			newResults = append(newResults, result)
			slog.Info("Lab result is new, will be inserted",
				"no_order", result.NoOrder,
				"param_code", result.ParamCode)
		} else {
			slog.Info("Lab result already exists, skipping",
				"no_order", result.NoOrder,
				"param_code", result.ParamCode)
		}
	}

	slog.Info("Filtered lab results",
		"total_input", len(results),
		"new_results", len(newResults),
		"existing_results", len(results)-len(newResults))

	return newResults, nil
}

// BatchInsertNewLabResults inserts only new lab results (that don't already exist)
func (s *Repository) BatchInsertNewLabResults(ctx context.Context, results []entity.SimrsLabResult) error {
	// First filter out existing results
	newResults, err := s.FilterNewLabResults(ctx, results)
	if err != nil {
		return fmt.Errorf("failed to filter new lab results: %w", err)
	}

	if len(newResults) == 0 {
		slog.Info("No new lab results to insert")
		return nil
	}

	// Insert only new results using simple insert
	return s.insertLabResultsSimple(ctx, newResults)
}

// insertLabResultsSimple performs simple insert without existence check (assumes data is already filtered)
func (s *Repository) insertLabResultsSimple(ctx context.Context, results []entity.SimrsLabResult) error {
	if len(results) == 0 {
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

	insertQuery := `
	INSERT INTO lab_results (
		no_order, param_code, result_value, unit, ref_range, flag, created_at
	) VALUES (
		@noOrder, @paramCode, @resultValue, @unit, @refRange, @flag, @createdAt
	)
	`

	now := time.Now()

	stmt, err := tx.PrepareContext(ctx, insertQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, result := range results {
		_, err = stmt.ExecContext(ctx,
			sql.Named("noOrder", result.NoOrder),
			sql.Named("paramCode", result.ParamCode),
			sql.Named("resultValue", result.ResultValue),
			sql.Named("unit", result.Unit),
			sql.Named("refRange", result.RefRange),
			sql.Named("flag", result.Flag),
			sql.Named("createdAt", now),
		)
		if err != nil {
			return fmt.Errorf("failed to execute insert for no_order %s, param_code %s: %w", result.NoOrder, result.ParamCode, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	slog.Info("Successfully inserted new lab results", "count", len(results))
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
	WHERE no_order = @noOrder
	ORDER BY created_at DESC`

	rows, err := conn.QueryContext(ctx, query, sql.Named("noOrder", noOrder))
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

	query := `DELETE FROM lab_requests WHERE no_order = @noOrder`

	result, err := conn.ExecContext(ctx, query, sql.Named("noOrder", noOrder))
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

	deleteQuery := `DELETE FROM lab_requests WHERE no_order = @noOrder`
	stmt, err := tx.PrepareContext(ctx, deleteQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare delete statement: %w", err)
	}
	defer stmt.Close()

	var deletedCount int64
	for _, noOrder := range noOrders {
		result, err := stmt.ExecContext(ctx, sql.Named("noOrder", noOrder))
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

	query := `DELETE FROM patient WHERE patient_id = @patientID`

	result, err := conn.ExecContext(ctx, query, sql.Named("patientID", patientID))
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

	deleteQuery := `DELETE FROM patient WHERE patient_id = @patientID`
	stmt, err := tx.PrepareContext(ctx, deleteQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare delete statement: %w", err)
	}
	defer stmt.Close()

	var deletedCount int64
	for _, patientID := range patientIDs {
		result, err := stmt.ExecContext(ctx, sql.Named("patientID", patientID))
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
