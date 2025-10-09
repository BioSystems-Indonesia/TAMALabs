package test_type

import (
	"context"
	"encoding/csv"
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
		if len(row) < 5 {
			err := fmt.Errorf("data in row is only %d, need 10", len(row))
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
	t.Name = strings.TrimSpace(row[0])
	t.Code = strings.TrimSpace(row[1])
	if len(t.Code) == 0 {
		return t, fmt.Errorf("code is empty"), true
	}

	t.Unit = strings.TrimSpace(row[2])
	low, err := strconv.ParseFloat(row[3], 64)
	if err != nil {
		allErr = append(allErr, fmt.Errorf("cannot parse float %s, %w", row[3], err))
	}
	t.LowRefRange = low

	high, err := strconv.ParseFloat(row[4], 64)
	if err != nil {
		allErr = append(allErr, fmt.Errorf("cannot parse float %s, %w", row[4], err))
	}
	t.HighRefRange = high

	return t, errors.Join(allErr...), skip
}
