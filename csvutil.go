package csvutil

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
)

// Dataset contains the CSV content and eases modificating it.
type Dataset struct {
	data [][]string
}

var (
	// ErrColNotFound is returned when a named column is not found.
	ErrColNotFound = errors.New("column not found")
	// ErrSameColName is returned when a named column exists multiple times.
	ErrSameColName = errors.New("multiple columns with same name")
	// ErrColLen is returned when the length of a column to add differs from the already existing dataset column length.
	ErrColLen = errors.New("length of columns differ")
)

// Equals checks if both datasets are the same.
func Equals(datasetA, datasetB [][]string) bool {
	if len(datasetA) != len(datasetB) {
		return false
	}
	for rowIdx, rowA := range datasetA {
		rowB := datasetB[rowIdx]
		if len(rowA) != len(rowB) {
			return false
		}
		for colIdx := range rowA {
			if rowA[colIdx] != rowB[colIdx] {
				return false
			}
		}
	}
	return true
}

// New creates a new CSV dataset.
func New(dataset [][]string) Dataset {
	var ds Dataset
	ds.data = dataset
	return ds
}

// Raw returns the raw data usable with Go's stdlib csv package.
func (ds *Dataset) Raw() [][]string {
	return ds.data
}

func (ds *Dataset) indexOfCol(name string) (int, error) {
	index := -1
	for idxCol, col := range ds.data[0] {
		if col == name {
			if index != -1 {
				return -1, fmt.Errorf("%w (%v)", ErrSameColName, name)
			}
			index = idxCol
		}
	}
	if index == -1 {
		return -1, fmt.Errorf("%w (%v)", ErrColNotFound, name)
	}

	return index, nil
}

// DeleteCol deletes the column with the given name.
func (ds *Dataset) DeleteCol(name string) error {
	index, err := ds.indexOfCol(name)
	if err != nil {
		return err
	}

	for idxRow := range ds.data {
		ds.data[idxRow] = append(ds.data[idxRow][:index], ds.data[idxRow][index+1:]...)
	}
	return nil
}

// RenameCol renames the header of column 'old' to 'new'.
func (ds *Dataset) RenameCol(old, new string) error {
	index, err := ds.indexOfCol(old)
	if err != nil {
		return err
	}
	ds.data[0][index] = new

	return nil
}

// AddRow appends the given row to the dataset.
func (ds *Dataset) AddRow(row []string) error {
	return ds.AddRows([][]string{row})
}

// AddRows appends the given rows to the dataset.
func (ds *Dataset) AddRows(rows [][]string) error {
	if len(rows) == 0 {
		return nil
	}
	if len(ds.data) > 0 && len(ds.data[0]) != len(rows[0]) {
		return ErrColLen
	}

	ds.data = append(ds.data, rows...)
	return nil
}

// AddCol inserts the given column at the position of the index.
// -1 adds the column at the last column
// -2 adds the column as the second last column, and so on...
func (ds *Dataset) AddCol(column []string, index int) error {
	if index < 0 {
		index += len(ds.data[0]) + 1
	}

	if len(ds.data) > 0 && len(column) != len(ds.data) {
		return ErrColLen
	}

	// no data so far
	if len(ds.data) == 0 {
		ds.data = [][]string{column}
		return nil
	}

	for idxRow := range ds.data {
		ds.data[idxRow] = append(ds.data[idxRow], "")
		copy(ds.data[idxRow][index+1:], ds.data[idxRow][index:])
		ds.data[idxRow][index] = column[idxRow]
	}

	return nil
}

// ExtractCol returns the column with the given name.
func (ds *Dataset) ExtractCol(name string) ([]string, error) {
	index, err := ds.indexOfCol(name)
	if err != nil {
		return []string{}, err
	}

	var resultCol []string

	for _, row := range ds.data {
		resultCol = append(resultCol, row[index])
	}

	return resultCol, nil
}

// ModifyCol changes the values of column 'name' according to 'f'.
// 'val' contains the column value and 'row' is the current row number.
func (ds *Dataset) ModifyCol(name string, f func(val string, row int) string) error {
	index, err := ds.indexOfCol(name)
	if err != nil {
		return err
	}

	if index > len(ds.data[0]) {
		return errors.New("index out of bounds")
	}

	for idxRow, row := range ds.data {
		// skip header
		if idxRow == 0 {
			continue
		}

		row[index] = f(row[index], idxRow)
	}
	return nil
}

type writeOptions struct {
	delimiter rune
	useCLRF   bool
}

type option func(*writeOptions)

// Delimiter is the separator between each value (default: ',').
func Delimiter(delimiter rune) option {
	return func(o *writeOptions) {
		o.delimiter = delimiter
	}
}

// UseCLRF is set to false by default. If set to true, the Writer ends each line with \r\n instead of \n.
func UseCLRF(useCLRF bool) option {
	return func(o *writeOptions) {
		o.useCLRF = useCLRF
	}
}

// Write the dataset to the given writer.
func (ds *Dataset) Write(writer io.Writer, opts ...option) error {
	csvWriter := csv.NewWriter(writer)
	defer csvWriter.Flush()

	// default options
	o := &writeOptions{
		delimiter: ',',
		useCLRF:   false,
	}

	// apply users options
	for _, opt := range opts {
		opt(o)
	}

	// pass options
	csvWriter.Comma = o.delimiter
	csvWriter.UseCRLF = o.useCLRF

	if err := csvWriter.WriteAll(ds.data); err != nil {
		return err
	}

	return nil
}
