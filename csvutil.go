package csvutil

import (
	"encoding/csv"
	"errors"
	"io"
)

type Dataset struct {
	hasHeader bool
	data      [][]string
}

var (
	ErrNoHeader       = errors.New("dataset is configured without header")
	ErrColumnNotFound = errors.New("column not found")
)

// New creates a new CSV dataset.
func New(dataset [][]string, header bool) Dataset {
	// dataset optional, can create an empty dataset
	// TODO: option if it contains a header
	// TODO: check for duplacate column names
	var ds Dataset
	ds.data = dataset
	ds.hasHeader = header
	return ds
}

// Raw returns the raw data usable with Go's stdlib csv package.
func (ds *Dataset) Raw() [][]string {
	return ds.data
}

// DeleteColumnID deletes the column with the given index.
// -1 deletes last column
// -2 deletes second last column, and so on...
func (ds *Dataset) DeleteColumnID(index int) error {
	if index < 0 {
		index += len(ds.data[0])
	}

	for idxRow, _ := range ds.data {
		ds.data[idxRow] = append(ds.data[idxRow][:index], ds.data[idxRow][index+1:]...)
	}
	return nil
}

func (ds *Dataset) indexOfCol(name string) (int, error) {
	index := -1
	for idxCol, col := range ds.data[0] {
		if col == name {
			if index != -1 {
				return -1, errors.New("dataset contains several column with this name. aborting.")
			}
			index = idxCol
		}
	}
	if index == -1 {
		return -1, ErrColumnNotFound // TODO: wrap error with 'name' as help
	}

	return index, nil
}

// DeleteColumn deletes the column with the given name.
// Requires a dataset with headers
func (ds *Dataset) DeleteColumn(name string) error {
	if !ds.hasHeader {
		return ErrNoHeader
	}

	idxToDelete, err := ds.indexOfCol(name)
	if err != nil {
		return err
	}
	return ds.DeleteColumnID(idxToDelete)
}

// RenameColumn renames the header of column 'old' to 'new'.
func (ds *Dataset) RenameColumn(old, new string) error {
	if !ds.hasHeader {
		return ErrNoHeader
	}

	index, err := ds.indexOfCol(old)
	if err != nil {
		return err
	}
	ds.data[0][index] = new

	return nil
}

// HasHeader returns if the dataset is configured with a header.
func (ds *Dataset) HasHeader() bool {
	return ds.hasHeader
}

// DeleteHeader deletes the header from the dataset.
func (ds *Dataset) DeleteHeader() error {
	if !ds.hasHeader {
		return ErrNoHeader
	}
	ds.data = ds.data[1:]
	ds.hasHeader = false
	return nil
}

// AddHeader adds a header to the dataset.
func (ds *Dataset) AddHeader(header []string) error {
	if len(header) == 0 {
		return nil
	}

	if ds.HasHeader() {
		return errors.New("dataset already contains a header")
	}

	// TODO: check for duplicate column names

	if len(ds.data) > 0 && len(ds.data) != len(header) {
		return errors.New("number of column names doesn't match with existing data")
	}
	// TODO: deny overwriting existing header, add option to allow overwriting
	ds.data[0] = append(ds.data[0], header...)
	ds.hasHeader = true
	return nil
}

// Header returns the header of the dataset.
// Returns an empty slice of strings when the dataset is
// configured without a header or the dataset is empty.
func (ds *Dataset) Header() []string {
	if len(ds.data) == 0 || ds.hasHeader {
		return []string{}
	}
	return ds.data[0]
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
	if len(ds.data) > 0 && len(ds.data) != len(rows) {
		return errors.New("number of rows doesn't match existing column length")
	}

	ds.data = append(ds.data, rows...)
	return nil
}

// AddColumn inserts the given column at the position of the index.
// -1 adds the column at the last column
// -2 adds the column as the second last column, and so on...
func (ds *Dataset) AddColumn(column []string, index int) error {
	// TODO: name optional as we can have a csv without header
	// TODO: index as option
	if index < 0 {
		index += len(ds.data[0])
	}

	if len(ds.data) > 0 && len(column) != len(ds.data) {
		return errors.New("column needs to have same length as existing data")
	}

	// no data so far
	if len(ds.data) == 0 {
		ds.data = [][]string{column}
		return nil
	}

	for idxRow, _ := range ds.data {
		ds.data[idxRow] = append(ds.data[idxRow], "")
		copy(ds.data[idxRow][index+1:], ds.data[idxRow][index:])
		ds.data[idxRow][index] = column[idxRow]
	}

	return nil
}

func (ds *Dataset) ExtractColumn(name string) ([]string, error) {
	index, err := ds.indexOfCol(name)
	if err != nil {
		return []string{}, err
	}

	return ds.ExtractColumnID(index)
}

func (ds *Dataset) ExtractColumnID(index int) ([]string, error) {
	if index < 0 {
		index += len(ds.data[0])
	}

	var resultCol []string

	for _, row := range ds.data {
		resultCol = append(resultCol, row[index])
	}

	return resultCol, nil
}

// ModifyColumn changes the values of the given column according to 'f'.
// 'val' contains the column value and 'row' is the current row number.
func (ds *Dataset) ModifyColumnID(index int, f func(val string, row int) string) error {
	if index < 0 {
		index += len(ds.data[0])
	}

	if index > len(ds.data[0]) {
		return errors.New("index out of bounds")
	}

	for idxRow, row := range ds.data {
		if ds.hasHeader && idxRow == 0 {
			continue
		}

		row[index] = f(row[index], idxRow)
	}
	return nil
}

// ModifyColumn changes the values of column 'name' according to 'f'.
// 'val' contains the column value and 'row' is the current row number.
func (ds *Dataset) ModifyColumn(name string, f func(val string, row int) string) error {
	index, err := ds.indexOfCol(name)
	if err != nil {
		return err
	}

	return ds.ModifyColumnID(index, f)
}

// Write the dataset to the given writer.
func (ds *Dataset) Write(writer io.Writer) error {
	csvWriter := csv.NewWriter(writer) // output as option
	defer csvWriter.Flush()

	// w.Comma = // TODO
	// w.UseCRLF

	if err := csvWriter.WriteAll(ds.data); err != nil {
		return err
	}

	// TODO: flush or lock internal dataset?

	return nil
}
