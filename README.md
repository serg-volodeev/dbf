# dbf
Package dbf reads and writes [DBF](http://en.wikipedia.org/wiki/DBase#File_formats) files.
The API of the __dbf__ package is similar to the __csv__ package from the standard library.

## Installation
You can incorporate the library into your local workspace with the following 'go get' command:

    go get github.com/serg-volodeev/dbf

## Using
Code needing to call into the library needs to include the following import statement:

```go
import (
    "github.com/serg-volodeev/dbf"
)
```

## Limitations
The following field types are supported:

- Character
- Numeric
- Logical
- Date

Memo fields are not supported. Index files are not supported.

## Examples
Сreate a file and write one record.
    
```go
f, err := os.Create("products.dbf")
if err != nil {
    log.Fatal(err)
}
defer f.Close()

fields := dbf.NewFields()
fields.AddCharacterField("NAME", 30)
fields.AddNumericField("COUNT", 8, 0)
fields.AddNumericField("PRICE", 12, 2)
fields.AddDateField("DATE")

w, err := dbf.NewWriter(f, fields, 1251)
if err != nil {
    log.Fatal(err)
}
record := make([]interface{}, fields.Count())
record[0] = "Apple"
record[1] = 1200
record[2] = 18.20
record[3] = time.Date(2021, 2, 12, 0, 0, 0, 0, time.UTC)

if err := w.Write(record); err != nil {
    log.Fatal(err)
}
if err := w.Flush(); err != nil {
    log.Fatal(err)
}
```

Read records.

```go
f, err := os.Open("products.dbf")
if err != nil {
    log.Fatal(err)
}
defer f.Close()

r, err := dbf.NewReader(f)
if err != nil {
    log.Fatal(err)
}

for i := uint32(0); i < r.RecordCount(); i++ {
    record, err := r.Read()
    if err != nil {
        log.Fatal(err)
    }
    name := record[0].(string)
    count := record[1].(int64)
    price := record[2].(float64)
    date := record[3].(time.Time)
    fmt.Println(name, count, price, date)
}
```

## License
Copyright (C) Sergey Volodeev. Released under MIT license.