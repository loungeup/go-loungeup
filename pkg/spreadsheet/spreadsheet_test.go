package spreadsheet_test

import (
	"bytes"
	"testing"

	"github.com/loungeup/go-loungeup/pkg/spreadsheet"
	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
)

func TestXLSXFile(t *testing.T) {
	rows := [][]string{
		{"First Name", "Last Name"},
		{"John", "Doe"},
		{"Jane", "Doe"},
	}

	file := spreadsheet.NewXLSXFile()

	sheet, err := file.AddSheet("A Simple Sheet")
	assert.NoError(t, err)
	assert.NoError(t, sheet.AddRows(rows...))

	data, err := file.Bytes()
	assert.NoError(t, err)

	assert.NoError(t, file.Close())

	t.Run("validate file content", func(t *testing.T) {
		rawFile, err := excelize.OpenReader(bytes.NewReader(data))
		assert.NoError(t, err)

		assert.Equal(t, 1, rawFile.SheetCount)

		rawRows, err := rawFile.GetRows("A Simple Sheet")
		assert.NoError(t, err)
		assert.Equal(t, rows, rawRows)
	})
}
