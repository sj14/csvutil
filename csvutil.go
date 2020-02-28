package csvutil

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
)

type Dataset struct {
	header []string
	data   [][]string
}

func (ds *Dataset) Raw() [][]string {
	result := append(ds.data, []string{})
	copy(result[1:], result[0:])
	result[0] = ds.header
	return result
}

func New(dataset [][]string) Dataset {
	// dataset optional, can create an empty dataset
	// TODO: option if it contains a header
	var ds Dataset
	ds.header = dataset[0]
	ds.data = dataset[1:]
	return ds
}

func (ds *Dataset) DeleteColumn(name string) error {
	if len(ds.header) != len(ds.data[0]) {
		return fmt.Errorf("dataset not consistent: %+v\n", ds.Raw())
	}

	idxToDelete := -1
	for idxCol, col := range ds.header {
		if col == name {
			idxToDelete = idxCol
			ds.header = append(ds.header[:idxToDelete], ds.header[idxToDelete+1:]...)
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
	renamed := false
	for i, column := range ds.header {
		if column == old {
			if renamed {
				return errors.New("old column name exists multiple times")
			}
			ds.header[i] = new
			renamed = true
		}
	}
	if !renamed {
		return errors.New("column not found")
	}

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
	ds.header = append(ds.header, header...)
	return nil
}

func (ds *Dataset) AddRow(row []string) error {
	return ds.AddRows([][]string{row})
}

func (ds *Dataset) AddRows(rows [][]string) error {
	if len(rows) == 0 {
		return nil
	}
	if len(ds.header) > 0 && len(ds.header) != len(rows) {
		return errors.New("number of rows doesn't match with header")
	}
	// TODO: check if length of header match data (when header already exists)

	ds.data = append(ds.data, rows...)
	return nil
}

func (ds *Dataset) AddColumn(name string, column []string, index int) error {
	// TODO: name optional as we can have a csv without header
	// TODO: index as option

	if len(column) != len(ds.data) && len(ds.data) > 0 {
		return errors.New("column needs to have same length as existing data")
	}

	// no data so far
	if len(ds.data) == 0 {
		ds.header = []string{name}
		ds.data = [][]string{column}
		return nil
	}

	// header
	ds.header = append(ds.header, "")
	copy(ds.header[index+1:], ds.header[index:])
	ds.header[index] = name

	for idxRow, _ := range ds.data {
		// for idxCol, col := range row {
		ds.data[idxRow] = append(ds.data[idxRow], "")
		copy(ds.data[idxRow][index+1:], ds.data[idxRow][index:])
		ds.data[idxRow][index] = column[idxRow]
		// }
	}

	return nil
}

// func  insert(dataset [][]string, toAdd []string, index int) {
// 	ds.data = append(ds.data, []string{})
// 	copy(ds.data[index+1:], ds.data[index:])
// 	ds.data[index] = toAdd
// }

func (ds *Dataset) WriteAll() error {
	w := csv.NewWriter(os.Stdout) // output as option
	defer w.Flush()

	// w.Comma = // TODO
	// w.UseCRLF

	// TODO: what heapens when header is empty?
	if err := w.Write(ds.header); err != nil {
		return err
	}

	if err := w.WriteAll(ds.data); err != nil {
		return err
	}

	return nil
}
