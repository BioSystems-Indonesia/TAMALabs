package test_type

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/BioSystems-Indonesia/TAMALabs/internal/entity"
)

func (u *Usecase) BulkCreate(ctx context.Context, r io.Reader) error {

	var allErr []error
	csvReader := csv.NewReader(r)

	_, err := csvReader.Read()
	if err != nil {
		err = fmt.Errorf("failed to read csv header: %w", err)
		return err
	}

	rows, err := csvReader.ReadAll()
	if err != nil {
		err = fmt.Errorf("failed to read csv: %w", err)
		return err
	}

	for i, row := range rows {
		if len(row) < 18 {
			err := fmt.Errorf("data in row is only %d, need 18", len(row))
			allErr = append(allErr, err)
			continue
		}

		testType, err, skip := u.generateTestTypeFromRowCSV(ctx, row)
		if err != nil {
			err := fmt.Errorf("error in generate test type row %d: %w", i, err)
			allErr = append(allErr, err)
			if skip {
				continue
			}
		}

		_, err = u.repository.Create(ctx, &testType)
		if err != nil {
			err = fmt.Errorf("failed to insert %s in row %d: %w", testType.Code, i, err)
			allErr = append(allErr, err)
		}
	}

	if len(allErr) > 0 {
		return errors.Join(allErr...)
	}

	return nil
}

func (u *Usecase) generateTestTypeFromRowCSV(ctx context.Context, row []string) (entity.TestType, error, bool) {
	var allErr []error
	var skip bool

	t := entity.TestType{}

	// id (row[0])
	if strings.TrimSpace(row[0]) != "" {
		id, err := strconv.Atoi(strings.TrimSpace(row[0]))
		if err != nil {
			allErr = append(allErr, fmt.Errorf("cannot parse id %s, %w", row[0], err))
		} else {
			t.ID = id
		}
	}

	// name (row[1])
	t.Name = strings.TrimSpace(row[1])

	// code (row[2])
	t.Code = strings.TrimSpace(row[2])
	if len(t.Code) == 0 {
		return t, fmt.Errorf("code is empty"), true
	}

	// alias_code (row[3])
	if strings.TrimSpace(row[3]) != "" {
		t.AliasCode = strings.TrimSpace(row[3])
	}

	// unit (row[4])
	t.Unit = strings.TrimSpace(row[4])

	// low_ref_range (row[5])
	if strings.TrimSpace(row[5]) != "" {
		low, err := strconv.ParseFloat(row[5], 64)
		if err != nil {
			allErr = append(allErr, fmt.Errorf("cannot parse low_ref_range %s, %w", row[5], err))
		} else {
			t.LowRefRange = low
		}
	}

	// high_ref_range (row[6])
	if strings.TrimSpace(row[6]) != "" {
		high, err := strconv.ParseFloat(row[6], 64)
		if err != nil {
			allErr = append(allErr, fmt.Errorf("cannot parse high_ref_range %s, %w", row[6], err))
		} else {
			t.HighRefRange = high
		}
	}

	// normal_ref_string (row[7])
	if strings.TrimSpace(row[7]) != "" {
		t.NormalRefString = strings.TrimSpace(row[7])
	}

	// decimal (row[8])
	if strings.TrimSpace(row[8]) != "" {
		decimal, err := strconv.Atoi(strings.TrimSpace(row[8]))
		if err != nil {
			allErr = append(allErr, fmt.Errorf("cannot parse decimal %s, %w", row[8], err))
		} else {
			t.Decimal = decimal
		}
	}

	// category (row[9])
	if strings.TrimSpace(row[9]) != "" {
		t.Category = strings.TrimSpace(row[9])
	}

	// sub_category (row[10])
	if strings.TrimSpace(row[10]) != "" {
		t.SubCategory = strings.TrimSpace(row[10])
	}

	// description (row[11])
	if strings.TrimSpace(row[11]) != "" {
		t.Description = strings.TrimSpace(row[11])
	}

	// is_calculated_test (row[12])
	if strings.TrimSpace(row[12]) != "" {
		isCalculated, err := strconv.ParseBool(strings.TrimSpace(row[12]))
		if err != nil {
			// Try parsing as int (0 or 1)
			isCalculatedInt, errInt := strconv.Atoi(strings.TrimSpace(row[12]))
			if errInt != nil {
				allErr = append(allErr, fmt.Errorf("cannot parse is_calculated_test %s, %w", row[12], err))
			} else {
				t.IsCalculatedTest = isCalculatedInt == 1
			}
		} else {
			t.IsCalculatedTest = isCalculated
		}
	}

	// device_id (row[13])
	if strings.TrimSpace(row[13]) != "" {
		deviceID, err := strconv.Atoi(strings.TrimSpace(row[13]))
		if err != nil {
			allErr = append(allErr, fmt.Errorf("cannot parse device_id %s, %w", row[13], err))
		} else {
			t.DeviceID = &deviceID
		}
	}

	// type (row[14])
	if strings.TrimSpace(row[14]) != "" {
		typeStr := strings.TrimSpace(row[14])
		t.Type = []entity.TestTypeSpecimenType{{Type: typeStr}}
	}

	// specific_ref_ranges (row[15]) - JSON format
	if strings.TrimSpace(row[15]) != "" {
		specificRangesStr := strings.TrimSpace(row[15])
		var ranges []entity.SpecificReferenceRange
		if err := json.Unmarshal([]byte(specificRangesStr), &ranges); err != nil {
			allErr = append(allErr, fmt.Errorf("cannot parse specific_ref_ranges JSON %s, %w", row[15], err))
		} else {
			t.SpecificRefRanges = ranges
		}
	}

	// loinc_code (row[16])
	if strings.TrimSpace(row[16]) != "" {
		t.LoincCode = strings.TrimSpace(row[16])
	}

	return t, errors.Join(allErr...), skip
}
