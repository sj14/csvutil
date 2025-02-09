# csvutil

[![GoDoc](https://godoc.org/github.com/sj14/csvutil?status.png)](https://godoc.org/github.com/sj14/csvutil)
[![Go Report Card](https://goreportcard.com/badge/github.com/sj14/csvutil)](https://goreportcard.com/report/github.com/sj14/csvutil)
![Action](https://github.com/sj14/csvutil/workflows/Go/badge.svg)

## Installation

```bash
go get -u github.com/sj14/csvutil
```

## Examples

For all options please check the godoc. The examples ignore all error handling!

---

### Create a new dataset

All datasets should contain a header for further processing.

```go
records := [][]string{
    {"first_name", "last_name", "username"},
    {"Rob", "Pike", "rob"},
    {"Ken", "Thompson", "ken"},
    {"Robert", "Griesemer", "gri"},
}

ds := csvutil.New(records)
```

### Add a new row

```go
ds.AddRow([]string{"my first name", "my last name", "my nick"})
```

```text
first_name,last_name,username
Rob,Pike,rob
Ken,Thompson,ken
Robert,Griesemer,gri
my first name,my last name,my nick
```

### Add multiple rows

```go
ds.AddRows([][]string{
    {"my first name 0", "my last name 0", "my nick 0"},
    {"my first name 1", "my last name 1", "my nick 1"},
    {"my first name 2", "my last name 2", "my nick 2"},
})
```

```text
first_name,last_name,username
Rob,Pike,rob
Ken,Thompson,ken
Robert,Griesemer,gri
my first name 0,my last name 0,my nick 0
my first name 1,my last name 1,my nick 1
my first name 2,my last name 2,my nick 2
```

### Add a new column at index 1

```go
ds.AddCol([]string{"column_headline", "my row 1", "my row 2", "my row 3"}, 1)
ds.Write(os.Stdout)
```

```text
first_name,column_headline,last_name,username
Rob,my row 1,Pike,rob
Ken,my row 2,Thompson,ken
Robert,my row 3,Griesemer,gri
```

### Get Column

```go
lastNames, _ := ds.GetCol("last_name")
fmt.Println(lastNames)
````

```text
[last_name Pike Thompson Griesemer]
```

### Rename Column

```go
ds.RenameCol("username", "nick")
```

```text
first_name,last_name,nick
Rob,Pike,rob
Ken,Thompson,ken
Robert,Griesemer,gri
```

### Move Column

```
ds.MoveCol("first_name", 1)
```

```
last_name,first_name,nick
Pike,Rob,rob
Thompson,Ken,ken
Griesemer,Robert,gri
```

### Delete Column

```go
ds.DeleteCol("first_name")
```

```text
last_name,nick
Pike,rob
Thompson,ken
Griesemer,gri
```

### Modify Column

```go
addRowNumber := func(val string, i int) string { return fmt.Sprintf("%v (%v)", val, i) }
ds.ModifyCol("first_name", addRowNumber)
ds.Write(os.Stdout)
```

```text
first_name,last_name,username
Rob (1),Pike,rob
Ken (2),Thompson,ken
Robert (3),Griesemer,gri
```

### Write the dataset

```go
ds.Write(os.Stdout)
```

```text
first_name,last_name,username
Rob,Pike,rob
Ken,Thompson,ken
Robert,Griesemer,gri
```

#### Write options

```go
ds.Write(os.Stdout, csvutil.Delimiter('|'), csvutil.UseCLRF(true))
```

```text
first_name|last_name|username
Rob|Pike|rob
Ken|Thompson|ken
Robert|Griesemer|gri
```

### Get dataset as [][]string

```go
fmt.Println(ds.Raw())
```

```text
[[first_name last_name username] [Rob Pike rob] [Ken Thompson ken] [Robert Griesemer gri]]
```

### Compare two (raw) datasets

```go
a := [][]string{
    {"a", "b", "c"},
    {"d", "e", "f"},
}

b := [][]string{
    {"g", "h", "i"},
    {"j", "k", "l"},
}

result := csvutil.Equals(a, b)
fmt.Println(result)
```

```text
false
```