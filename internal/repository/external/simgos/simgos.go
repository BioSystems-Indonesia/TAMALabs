package simgos

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
)

type Repository struct {
	simgosDB *DB
}

func NewRepository(simgosDB *DB) *Repository {
	return &Repository{
		simgosDB: simgosDB,
	}
}

// GetNewLabOrders fetches all lab orders with status 'NEW' from Database Sharing
func (r *Repository) GetNewLabOrders(ctx context.Context) ([]entity.SimrsLabOrder, error) {
	if r.simgosDB == nil {
		return nil, fmt.Errorf("Database Sharing database not initialized")
	}

	conn := r.simgosDB.GetConnection()
	if conn == nil {
		return nil, fmt.Errorf("Database Sharing database connection is nil")
	}

	query := `
		SELECT 
			id,
			no_lab_order,
			no_rm,
			patient_name,
			birth_date,
			sex,
			IFNULL(doctor, '') as doctor,
			IFNULL(analyst, '') as analyst,
			status,
			created_at
		FROM lab_order
		WHERE status = ?
		  AND DATE(created_at) = CURDATE()
		ORDER BY created_at ASC
	`

	rows, err := conn.QueryContext(ctx, query, entity.SimrsStatusNew)
	if err != nil {
		return nil, fmt.Errorf("failed to query new lab orders: %w", err)
	}
	defer rows.Close()

	var orders []entity.SimrsLabOrder
	for rows.Next() {
		var order entity.SimrsLabOrder
		err := rows.Scan(
			&order.ID,
			&order.NoLabOrder,
			&order.NoRM,
			&order.PatientName,
			&order.BirthDate,
			&order.Sex,
			&order.Doctor,
			&order.Analyst,
			&order.Status,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lab order row: %w", err)
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating lab order rows: %w", err)
	}

	slog.Info("Fetched new lab orders from Database Sharing", "count", len(orders))
	return orders, nil
}

// GetOrderDetailsByNoLabOrder fetches all order details for a specific lab order
func (r *Repository) GetOrderDetailsByNoLabOrder(ctx context.Context, noLabOrder string) ([]entity.SimrsOrderDetail, error) {
	if r.simgosDB == nil {
		return nil, fmt.Errorf("Database Sharing database not initialized")
	}

	conn := r.simgosDB.GetConnection()
	if conn == nil {
		return nil, fmt.Errorf("Database Sharing database connection is nil")
	}

	query := `
		SELECT 
			id,
			no_lab_order,
			parameter_code,
			parameter_name,
			IFNULL(result_value, '') as result_value,
			IFNULL(unit, '') as unit,
			IFNULL(reference_range, '') as reference_range,
			IFNULL(flag, '') as flag,
			created_at
		FROM order_detail
		WHERE no_lab_order = ?
		ORDER BY id ASC
	`

	rows, err := conn.QueryContext(ctx, query, noLabOrder)
	if err != nil {
		return nil, fmt.Errorf("failed to query order details: %w", err)
	}
	defer rows.Close()

	var details []entity.SimrsOrderDetail
	for rows.Next() {
		var detail entity.SimrsOrderDetail
		err := rows.Scan(
			&detail.ID,
			&detail.NoLabOrder,
			&detail.ParameterCode,
			&detail.ParameterName,
			&detail.ResultValue,
			&detail.Unit,
			&detail.ReferenceRange,
			&detail.Flag,
			&detail.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order detail row: %w", err)
		}
		details = append(details, detail)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating order detail rows: %w", err)
	}

	return details, nil
}

// UpdateOrderStatus updates the status of a lab order
func (r *Repository) UpdateOrderStatus(ctx context.Context, noLabOrder string, status entity.SimrsOrderStatus) error {
	if r.simgosDB == nil {
		return fmt.Errorf("Database Sharing database not initialized")
	}

	conn := r.simgosDB.GetConnection()
	if conn == nil {
		return fmt.Errorf("Database Sharing database connection is nil")
	}

	query := `
		UPDATE lab_order 
		SET status = ?
		WHERE no_lab_order = ?
	`

	result, err := conn.ExecContext(ctx, query, status, noLabOrder)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no order found with no_lab_order: %s", noLabOrder)
	}

	slog.Info("Updated order status", "no_lab_order", noLabOrder, "status", status)
	return nil
}

// BatchUpdateOrderDetails updates multiple order detail results
func (r *Repository) BatchUpdateOrderDetails(ctx context.Context, details []entity.SimrsOrderDetail) error {
	if len(details) == 0 {
		return nil
	}

	if r.simgosDB == nil {
		return fmt.Errorf("Database Sharing database not initialized")
	}

	conn := r.simgosDB.GetConnection()
	if conn == nil {
		return fmt.Errorf("Database Sharing database connection is nil")
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

	updateQuery := `
		UPDATE order_detail
		SET result_value = ?,
		    unit = ?,
		    reference_range = ?,
		    flag = ?
		WHERE no_lab_order = ? AND parameter_code = ?
	`

	stmt, err := tx.PrepareContext(ctx, updateQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer stmt.Close()

	var updatedCount int
	for _, detail := range details {
		result, err := stmt.ExecContext(ctx,
			detail.ResultValue,
			detail.Unit,
			detail.ReferenceRange,
			detail.Flag,
			detail.NoLabOrder,
			detail.ParameterCode,
		)
		if err != nil {
			return fmt.Errorf("failed to update order detail for no_lab_order %s, parameter_code %s: %w",
				detail.NoLabOrder, detail.ParameterCode, err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("failed to get rows affected: %w", err)
		}

		if rowsAffected > 0 {
			updatedCount++
			slog.Info("Updated order detail",
				"no_lab_order", detail.NoLabOrder,
				"parameter_code", detail.ParameterCode,
				"result_value", detail.ResultValue)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	slog.Info("Batch update order details completed",
		"total_processed", len(details),
		"updated", updatedCount)

	return nil
}

// GetCompletedLabOrders fetches all lab orders with status 'LIS_SUCCESS'
func (r *Repository) GetCompletedLabOrders(ctx context.Context) ([]entity.SimrsLabOrder, error) {
	if r.simgosDB == nil {
		return nil, fmt.Errorf("Database Sharing database not initialized")
	}

	conn := r.simgosDB.GetConnection()
	if conn == nil {
		return nil, fmt.Errorf("Database Sharing database connection is nil")
	}

	query := `
		SELECT 
			id,
			no_lab_order,
			no_rm,
			patient_name,
			birth_date,
			sex,
			IFNULL(doctor, '') as doctor,
			IFNULL(analyst, '') as analyst,
			status,
			created_at
		FROM lab_order
		WHERE status = ?
		  AND DATE(created_at) = CURDATE()
		ORDER BY created_at ASC
	`

	rows, err := conn.QueryContext(ctx, query, entity.SimrsStatusLISSuccess)
	if err != nil {
		return nil, fmt.Errorf("failed to query completed lab orders: %w", err)
	}
	defer rows.Close()

	var orders []entity.SimrsLabOrder
	for rows.Next() {
		var order entity.SimrsLabOrder
		err := rows.Scan(
			&order.ID,
			&order.NoLabOrder,
			&order.NoRM,
			&order.PatientName,
			&order.BirthDate,
			&order.Sex,
			&order.Doctor,
			&order.Analyst,
			&order.Status,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lab order row: %w", err)
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating lab order rows: %w", err)
	}

	return orders, nil
}

// CheckOrderExists checks if a lab order exists by no_lab_order
func (r *Repository) CheckOrderExists(ctx context.Context, noLabOrder string) (bool, error) {
	if r.simgosDB == nil {
		return false, fmt.Errorf("Database Sharing database not initialized")
	}

	conn := r.simgosDB.GetConnection()
	if conn == nil {
		return false, fmt.Errorf("Database Sharing database connection is nil")
	}

	query := `SELECT COUNT(1) FROM lab_order WHERE no_lab_order = ?`

	var count int
	err := conn.QueryRowContext(ctx, query, noLabOrder).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check order existence: %w", err)
	}

	return count > 0, nil
}

// GetLabOrderByNoLabOrder fetches a specific lab order by no_lab_order
func (r *Repository) GetLabOrderByNoLabOrder(ctx context.Context, noLabOrder string) (*entity.SimrsLabOrder, error) {
	if r.simgosDB == nil {
		return nil, fmt.Errorf("Database Sharing database not initialized")
	}

	conn := r.simgosDB.GetConnection()
	if conn == nil {
		return nil, fmt.Errorf("Database Sharing database connection is nil")
	}

	query := `
		SELECT 
			id,
			no_lab_order,
			no_rm,
			patient_name,
			birth_date,
			sex,
			IFNULL(doctor, '') as doctor,
			IFNULL(analyst, '') as analyst,
			status,
			created_at
		FROM lab_order
		WHERE no_lab_order = ?
	`

	var order entity.SimrsLabOrder
	err := conn.QueryRowContext(ctx, query, noLabOrder).Scan(
		&order.ID,
		&order.NoLabOrder,
		&order.NoRM,
		&order.PatientName,
		&order.BirthDate,
		&order.Sex,
		&order.Doctor,
		&order.Analyst,
		&order.Status,
		&order.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("lab order with no_lab_order %s not found", noLabOrder)
		}
		return nil, fmt.Errorf("failed to get lab order: %w", err)
	}

	return &order, nil
}

// CountOrderDetailsByNoLabOrder counts the number of order details for a lab order
func (r *Repository) CountOrderDetailsByNoLabOrder(ctx context.Context, noLabOrder string) (int, error) {
	if r.simgosDB == nil {
		return 0, fmt.Errorf("Database Sharing database not initialized")
	}

	conn := r.simgosDB.GetConnection()
	if conn == nil {
		return 0, fmt.Errorf("Database Sharing database connection is nil")
	}

	query := `SELECT COUNT(1) FROM order_detail WHERE no_lab_order = ?`

	var count int
	err := conn.QueryRowContext(ctx, query, noLabOrder).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count order details: %w", err)
	}

	return count, nil
}

// CountCompletedOrderDetailsByNoLabOrder counts the number of completed order details (with result_value)
func (r *Repository) CountCompletedOrderDetailsByNoLabOrder(ctx context.Context, noLabOrder string) (int, error) {
	if r.simgosDB == nil {
		return 0, fmt.Errorf("Database Sharing database not initialized")
	}

	conn := r.simgosDB.GetConnection()
	if conn == nil {
		return 0, fmt.Errorf("Database Sharing database connection is nil")
	}

	query := `
		SELECT COUNT(1) 
		FROM order_detail 
		WHERE no_lab_order = ? 
		  AND result_value IS NOT NULL 
		  AND result_value != ''
	`

	var count int
	err := conn.QueryRowContext(ctx, query, noLabOrder).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count completed order details: %w", err)
	}

	return count, nil
}

// TestConnection tests the Database Sharing database connection
func (r *Repository) TestConnection(ctx context.Context) error {
	if r.simgosDB == nil {
		return fmt.Errorf("Database Sharing database not initialized")
	}

	return r.simgosDB.Ping()
}

// GetOrdersForSync fetches lab orders that need to be synced (within last 14 days)
func (r *Repository) GetOrdersForSync(ctx context.Context, startDate, endDate time.Time) ([]entity.SimrsLabOrder, error) {
	if r.simgosDB == nil {
		return nil, fmt.Errorf("Database Sharing database not initialized")
	}

	conn := r.simgosDB.GetConnection()
	if conn == nil {
		return nil, fmt.Errorf("Database Sharing database connection is nil")
	}

	query := `
		SELECT 
			id,
			no_lab_order,
			no_rm,
			patient_name,
			birth_date,
			sex,
			IFNULL(doctor, '') as doctor,
			IFNULL(analyst, '') as analyst,
			status,
			created_at
		FROM lab_order
		WHERE DATE(created_at) = CURDATE()
		  AND status IN (?, ?)
		ORDER BY created_at DESC
	`

	rows, err := conn.QueryContext(ctx, query, entity.SimrsStatusPending, entity.SimrsStatusLISSuccess)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders for sync: %w", err)
	}
	defer rows.Close()

	var orders []entity.SimrsLabOrder
	for rows.Next() {
		var order entity.SimrsLabOrder
		err := rows.Scan(
			&order.ID,
			&order.NoLabOrder,
			&order.NoRM,
			&order.PatientName,
			&order.BirthDate,
			&order.Sex,
			&order.Doctor,
			&order.Analyst,
			&order.Status,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan lab order row: %w", err)
		}
		orders = append(orders, order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating lab order rows: %w", err)
	}

	return orders, nil
}
