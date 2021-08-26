# dbf
Package dbf reads and writes [DBF](http://en.wikipedia.org/wiki/DBase#File_formats) files.

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

w.SetStringFieldValue(0, "Apple")
w.SetIntFieldValue(1, 1200)
w.SetFloatFieldValue(2, 18.20)
w.SetDateFieldValue(3, time.Date(2021, 2, 12, 0, 0, 0, 0, time.UTC))

w.Write()

w.Flash()

if w.Err() != nil {
    log.Fatal(w.Err())
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

for r.Read() {
    name := r.StringFieldValue(0)
    count := r.IntFieldValue(1)
    price := r.FloatFieldValue(2)
    date := r.DateFieldValue(3)
    
    fmt.Println(name, count, price, date)
}

if r.Err() != nil {
    log.Fatal(r.Err())
}
```

## License
Copyright (C) Sergey Volodeev. Released under MIT license.