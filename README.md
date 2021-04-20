# dbf
Package dbf reads and writes [DBF](http://en.wikipedia.org/wiki/DBase#File_formats) files.
The API of the __dbf__ package is similar to the __csv__ package from the standard library.

## Limitations
The following field types are supported: __C__, __N__, __L__, __D__.
Memo fields are not supported. Index files are not supported.

## Examples
Сreate a file and write one record.

    f, err := os.Create("products.dbf")
    if err != nil {
        return err
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
        return err
    }
    record := make([]interface{}, len(fields))
    record[0] = "Apple"
    record[1] = 1200
    record[2] = 18.20
    record[3] = time.Date(2021, 2, 12, 0, 0, 0, 0, time.UTC)

    if err := w.Write(record); err != nil {
        return err
    }

    if err := w.Flush(); err != nil {
        return err
    }

Read records.

    f, err := os.Open("products.dbf")
    if err != nil {
        return err
    }
    defer f.Close()

    r, err := dbf.NewReader(f, 0)
    if err != nil {
        return err
    }

    for {
        record, err := r.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }
        fmt.Println(record)
    }

## License
Copyright (C) Sergey Volodeev. Released under MIT license.