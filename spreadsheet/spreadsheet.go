// Package spreadsheet provides a unified way to interact with spreadsheets.
package spreadsheet

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type XLSXFile struct {
	baseFile *excelize.File
}

func NewXLSXFile() *XLSXFile {
	return &XLSXFile{
		baseFile: excelize.NewFile(),
	}
}

func (f *XLSXFile) AddSheet(name string) (*XLSXSheet, error) {
	if _, err := f.baseFile.NewSheet(name); err != nil {
		return nil, fmt.Errorf("could not add sheet to file: %w", err)
	}

	return &XLSXSheet{
		baseFile:        f.baseFile,
		currentRowIndex: 1,
		name:            name,
	}, nil
}

func (f *XLSXFile) Bytes() ([]byte, error) {
	if err := f.deleteDefaultSheet(); err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	if err := f.baseFile.Write(&buffer); err != nil {
		return nil, fmt.Errorf("could not write file to buffer: %w", err)
	}

	return buffer.Bytes(), nil
}

func (f *XLSXFile) Close() error { return f.baseFile.Close() }

func (f *XLSXFile) deleteDefaultSheet() error {
	if err := f.baseFile.DeleteSheet("Sheet1"); err != nil {
		return fmt.Errorf("could not delete default sheet from file: %w", err)
	}

	return nil
}

type XLSXSheet struct {
	baseFile        *excelize.File
	currentRowIndex int
	name            string
}

func (s *XLSXSheet) AddRows(rows ...[]string) error {
	for _, row := range rows {
		if err := s.addRow(row); err != nil {
			return err
		}
	}

	return nil
}

func (s *XLSXSheet) addRow(row []string) error {
	if err := s.baseFile.SetSheetRow(s.name, "A"+strconv.Itoa(s.currentRowIndex), &row); err != nil {
		return fmt.Errorf("could not add row to sheet: %w", err)
	}

	s.currentRowIndex++

	return nil
}
