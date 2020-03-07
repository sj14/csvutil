# csvutil

[![GoDoc](https://godoc.org/github.com/sj14/csvutil?status.png)](https://godoc.org/github.com/sj14/csvutil)
[![Go Report Card](https://goreportcard.com/badge/github.com/sj14/csvutil)](https://goreportcard.com/report/github.com/sj14/csvutil)

## Examples

For all options please check the godoc.
The examples ignore all error handling.

---

Create a new dataset:

```go
records := [][]string{
    {"first_name", "last_name", "username"},
    {"Rob", "Pike", "rob"},
    {"Ken", "Thompson", "ken"},
    {"Robert", "Griesemer", "gri"},
}

ds := csvutil.New(records)
```

Write dataset:

```go
ds.Write(os.Stdout)
```

```text
first_name,last_name,username
Rob,Pike,rob
Ken,Thompson,ken
Robert,Griesemer,gri
```

Add a new column at index 1:

```go
ds.AddCol([]string{"asd", "1", "2", "3"}, 1)
ds.Write(os.Stdout)
```

```text
first_name,asd,last_name,username
Rob,1,Pike,rob
Ken,2,Thompson,ken
Robert,3,Griesemer,gri
```

Modify columns:

```go
addRowNumber := func(val string, i int) string { return fmt.Sprintf("%v (%v)", val, i) }
ds.ModifyCol("first_name", addRowNumber)
ds.Write(os.Stdout)
```

```text
first_name,asd,last_name,username
Rob (1),1,Pike,rob
Ken (2),2,Thompson,ken
Robert (3),3,Griesemer,gri
```

Add write options:

```go
ds.Write(os.Stdout, csvutil.Delimiter('|'), csvutil.UseCLRF(true))
```

```text
first_name|last_name|username
Rob|Pike|rob
Ken|Thompson|ken
Robert|Griesemer|gri
```
