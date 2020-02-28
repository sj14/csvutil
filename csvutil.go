package csvutil

import (
	"encoding/csv"
	"errors"
	"os"
)

type Dataset struct {
	hasHeader bool
	data      [][]string
}

func (ds *Dataset) Raw() [][]string {
	return ds.data
}

func New(dataset [][]string, header bool) Dataset {
	// dataset optional, can create an empty dataset
	// TODO: option if it contains a header
	var ds Dataset
	ds.data = dataset
	ds.hasHeader = header
	return ds
}

func (ds *Dataset) DeleteColumn(name string) error {
	idxToDelete := -1
	for idxCol, col := range ds.data[0] {
		if col == name {
			idxToDelete = idxCol
		}
	}

	if idxToDelete == -1 {
		return errors.New("column not found")
	}

	for idxRow, _ := range ds.data {
		ds.data[idxRow] = append(ds.data[idxRow][:idxToDelete], ds.data[idxRow][idxToDelete+1:]...)
	}
	return nil
}

func (ds *Dataset) RenameColumn(old, new string) error {
	if !ds.hasHeader {
		return errors.New("dataset is configured without header")
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
		return errors.New("column not found")
	}

	return nil
}

func (ds *Dataset) HasHeader() bool {
	return ds.hasHeader
}

func (ds *Dataset) DeleteHeader() error {
	if !ds.hasHeader {
		return errors.New("dataset is configured without header")
	}
	ds.data = ds.data[1:]
	ds.hasHeader = false
	return nil
}

func (ds *Dataset) AddHeader(header []string) error {
	if len(header) == 0 {
		return nil
	}
	if len(ds.data) > 0 && len(ds.data) != len(header) {
		return errors.New("number of column names doesn't match with existing data")
	}
	// TODO: deny overwriting existing header, add option to allow overwriting
	ds.data[0] = append(ds.data[0], header...)
	ds.hasHeader = true
	return nil
}

func (ds *Dataset) Header() []string {
	if len(ds.data) == 0 {
		return []string{}
	}
	return ds.data[0]
}

func (ds *Dataset) AddRow(row []string) error {
	return ds.AddRows([][]string{row})
}

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

func (ds *Dataset) AddColumn(column []string, index int) error {
	// TODO: name optional as we can have a csv without header
	// TODO: index as option

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

func (ds *Dataset) Write() error {
	w := csv.NewWriter(os.Stdout) // output as option
	defer w.Flush()

	// w.Comma = // TODO
	// w.UseCRLF

	if err := w.WriteAll(ds.data); err != nil {
		return err
	}

	return nil
}
