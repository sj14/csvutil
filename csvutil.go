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

func (ds *Dataset) Raw() [][]string {
	return ds.data
}

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

// DeleteColumn deletes the column with the given name.
// Requires a dataset with headers
func (ds *Dataset) DeleteColumn(name string) error {
	if !ds.hasHeader {
		return ErrNoHeader
	}
	idxToDelete := -1
	for idxCol, col := range ds.data[0] {
		if col == name {
			if idxToDelete != -1 {
				return errors.New("dataset contains several column with this name. aborting.")
			}
			idxToDelete = idxCol
		}
	}

	if idxToDelete == -1 {
		return ErrColumnNotFound // TODO: wrap error with 'name' as help
	}

	return ds.DeleteColumnID(idxToDelete)
}

var (
	ErrNoHeader       = errors.New("dataset is configured without header")
	ErrColumnNotFound = errors.New("column not found")
)

// RenameColumn renames the header of column 'old' to 'new'.
func (ds *Dataset) RenameColumn(old, new string) error {
	if !ds.hasHeader {
		return ErrNoHeader
	}

	renamed := false
	for idxCol, column := range ds.data[0] {
		if column == old {
			if renamed {
				return errors.New("old column name exists multiple times")
			}
			ds.data[0][idxCol] = new
			renamed = true
		}
	}
	if !renamed {
		return ErrColumnNotFound
	}

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

// func  insert(dataset [][]string, toAdd []string, index int) {
// 	ds.data = append(ds.data, []string{})
// 	copy(ds.data[index+1:], ds.data[index:])
// 	ds.data[index] = toAdd
// }

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
