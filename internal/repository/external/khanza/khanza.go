package khanza

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/oibacidem/lims-hl-seven/internal/entity"
)

// Repository is a repository for external applications using SIM RS Repository
type Repository struct {
	bridgeDB *DB
	mainDB   *DB
}

// NewRepository creates a new Khanza instance
func NewRepository(bridgeDB *DB, mainDB *DB) *Repository {
	return &Repository{
		bridgeDB: bridgeDB,
		mainDB:   mainDB,
	}
}

func (s *Repository) RollbackTransaction(ctx context.Context, tx *sql.Tx) {
	if err := tx.Rollback(); err != nil {
		slog.Error("failed to rollback transaction", "error", err)
	}
}

// BatchUpsertRESDTO performs batch upsert operations for multiple RESDT records in a single transaction
func (s *Repository) BatchUpsertRESDTO(ctx context.Context, reqs []entity.KhanzaResDT) (err error) {
	if len(reqs) == 0 {
		return nil // No records to process
	}

	tx, err := s.bridgeDB.GetConnection().BeginTx(ctx, &sql.TxOptions{})
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

	// Use individual INSERT ... ON DUPLICATE KEY UPDATE for better compatibility
	insertQuery := `INSERT INTO RESDT (
		ONO, ORDER_TESTID, TEST_CD, TEST_NM, DATA_TYP, RESULT_VALUE, RESULT_FT, UNIT, FLAG, REF_RANGE, VALIDATE_ON
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
		ORDER_TESTID = VALUES(ORDER_TESTID),
		DATA_TYP = VALUES(DATA_TYP),
		RESULT_VALUE = VALUES(RESULT_VALUE),
		RESULT_FT = VALUES(RESULT_FT),
		UNIT = VALUES(UNIT),
		FLAG = VALUES(FLAG),
		REF_RANGE = VALUES(REF_RANGE),
		VALIDATE_ON = VALUES(VALIDATE_ON)`

	now := time.Now()

	// Prepare statement for better performance
	stmt, err := tx.PrepareContext(ctx, insertQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	// Execute for each record
	for _, resdt := range reqs {
		_, err = stmt.ExecContext(ctx, resdt.ONO, resdt.OrderTestID, resdt.TESTCD, resdt.TestNM, resdt.DataTyp,
			resdt.ResultValue, resdt.ResultFT, resdt.Unit, resdt.Flag, resdt.RefRange, now.Format(time.RFC3339))
		if err != nil {
			return fmt.Errorf("failed to execute insert for ONO %s: %w", resdt.ONO, err)
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetLisOrderByVisitNo retrieves a LisOrder record by VISITNO
func (k *Repository) GetLisOrderByVisitNo(visitNo string) (*entity.KhanzaLisOrder, error) {
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
	FROM LIS_ORDER 
	WHERE VISITNO = ?`

	conn := k.bridgeDB.GetConnection()
	row := conn.QueryRow(query, visitNo)

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
		return nil, fmt.Errorf("failed to get LisOrder by visitNo: %w", err)
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
	FROM LIS_ORDER
	WHERE MESSAGE_DT >= DATE_SUB(NOW(), INTERVAL 72 HOUR)
	ORDER BY ID DESC 
	LIMIT 1000`

	conn := k.bridgeDB.GetConnection()
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

// GetLabRequestByNoOrder retrieves lab request details by noorder from main DB
func (k *Repository) GetLabRequestByNoOrder(ctx context.Context, noorder string) ([]entity.KhanzaLabRequest, error) {
	query := `SELECT 
		permintaan_lab.noorder,
		permintaan_lab.no_rawat,
		reg_periksa.no_rkm_medis,
		pasien.nm_pasien,
		jns_perawatan_lab.nm_perawatan,
		permintaan_detail_permintaan_lab.id_template,
		template_laboratorium.Pemeriksaan,
		template_laboratorium.satuan,
		template_laboratorium.nilai_rujukan_ld,
		reg_periksa.kd_pj,
		template_laboratorium.nilai_rujukan_la,
		template_laboratorium.nilai_rujukan_pd,
		template_laboratorium.nilai_rujukan_pa,
		permintaan_lab.tgl_permintaan,
		penjab.png_jawab,
		IF(permintaan_lab.jam_permintaan='00:00:00','',permintaan_lab.jam_permintaan) as jam_permintaan,
		permintaan_lab.tgl_sampel,
		IF(permintaan_lab.jam_sampel='00:00:00','',permintaan_lab.jam_sampel) as jam_sampel,
		permintaan_lab.tgl_hasil,
		IF(permintaan_lab.jam_hasil='00:00:00','',permintaan_lab.jam_hasil) as jam_hasil,
		permintaan_lab.dokter_perujuk,
		dokter.nm_dokter,
		poliklinik.nm_poli,
		permintaan_lab.informasi_tambahan,
		permintaan_lab.diagnosa_klinis 
	FROM permintaan_lab 
	INNER JOIN reg_periksa ON permintaan_lab.no_rawat=reg_periksa.no_rawat 
	INNER JOIN pasien ON reg_periksa.no_rkm_medis=pasien.no_rkm_medis 
	INNER JOIN permintaan_pemeriksaan_lab ON permintaan_lab.noorder=permintaan_pemeriksaan_lab.noorder 
	INNER JOIN poliklinik ON reg_periksa.kd_poli=poliklinik.kd_poli 
	INNER JOIN jns_perawatan_lab ON jns_perawatan_lab.kd_jenis_prw=permintaan_pemeriksaan_lab.kd_jenis_prw 
	INNER JOIN permintaan_detail_permintaan_lab ON permintaan_lab.noorder=permintaan_detail_permintaan_lab.noorder AND permintaan_detail_permintaan_lab.kd_jenis_prw=permintaan_pemeriksaan_lab.kd_jenis_prw 
	INNER JOIN template_laboratorium ON template_laboratorium.id_template=permintaan_detail_permintaan_lab.id_template 
	INNER JOIN dokter ON permintaan_lab.dokter_perujuk=dokter.kd_dokter 
	INNER JOIN penjab ON reg_periksa.kd_pj=penjab.kd_pj 
	WHERE permintaan_lab.noorder=?
	ORDER BY permintaan_lab.tgl_permintaan DESC, permintaan_lab.jam_permintaan DESC`

	conn := k.mainDB.GetConnection()
	rows, err := conn.QueryContext(ctx, query, noorder)
	if err != nil {
		return nil, fmt.Errorf("failed to query lab requests: %w", err)
	}
	defer rows.Close()

	var labRequests []entity.KhanzaLabRequest
	for rows.Next() {
		var labRequest entity.KhanzaLabRequest
		err := rows.Scan(
			&labRequest.NoOrder,
			&labRequest.NoRawat,
			&labRequest.NoRkmMedis,
			&labRequest.NmPasien,
			&labRequest.NmPerawatan,
			&labRequest.IDTemplate,
			&labRequest.Pemeriksaan,
			&labRequest.Satuan,
			&labRequest.NilaiRujukanLD,
			&labRequest.KdPj,
			&labRequest.NilaiRujukanLA,
			&labRequest.NilaiRujukanPD,
			&labRequest.NilaiRujukanPA,
			&labRequest.TglPermintaan,
			&labRequest.PngJawab,
			&labRequest.JamPermintaan,
			&labRequest.TglSampel,
			&labRequest.JamSampel,
			&labRequest.TglHasil,
			&labRequest.JamHasil,
			&labRequest.DokterPerujuk,
			&labRequest.NmDokter,
			&labRequest.NmPoli,
			&labRequest.InformasiTambahan,
			&labRequest.DiagnosaKlinis,
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
