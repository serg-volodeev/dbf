# dbf
Package dbf reads and writes [DBF](http://en.wikipedia.org/wiki/DBase#File_formats) files.
The API of the __dbf__ package is similar to the __csv__ package from the standard library.

## Installation
You can incorporate the library into your local workspace with the following 'go get' command:

    go get github.com/serg-volodeev/dbf

## Using
Code needing to call into the library needs to include the following import statement:

    import (
        "github.com/serg-volodeev/dbf"
    )

## Limitations
The following field types are supported: __C__, __N__, __L__, __D__.
Memo fields are not supported. Index files are not supported.

## Examples
Сreate a file and write one record.

    f, err := os.Create("products.dbf")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    fields := []dbf.FieldInfo{
        {"NAME", "C", 30, 0},
        {"COUNT", "N", 8, 0},
        {"PRICE", "N", 12, 2},
        {"DATE", "D", 8, 0},
    }
    w, err := dbf.NewWriter(f, fields, 1251)
    if err != nil {
        log.Fatal(err)
    }
    record := make([]interface{}, len(fields))
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

Read records.

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
        fmt.Println(record)
    }

## License
Copyright (C) Sergey Volodeev. Released under MIT license.