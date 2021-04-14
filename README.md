# xbase
Go library for working with DBF files

Example create file

    db := xbase.New()

    if err := db.AddField("NAME", "C", 30); err != nil {
        return
    }
    if err := db.AddField("SALARY", "N", 9, 2); err != nil {
        return
    }
    if err := db.AddField("BDATE", "D"); err != nil {
        return
    }
    db.SetCodePage(1251)
    if err := db.CreateFile("persons.dbf"); err != nil {
        return
    }
    defer db.CloseFile()
    
    // Add record
    db.Add()
    if err := db.SetFieldValue(1, "John Smith"); err != nil {
        return
    }
    if err := db.SetFieldValue(2, 1234.56); err != nil {
        return
    }
    if err := db.SetFieldValue(3, time.Date(1998, 2, 20, 0, 0, 0, 0, time.UTC)); err != nil {
        return
    }
    if err := db.Save(); err != nil {
        return
    }

Example read file

    db := xbase.New()
    if err := db.OpenFile("persons.dbf", true); err != nil {
        return
    }
    defer db.CloseFile()
    if err := db.First(); err != nil {
        return
    }
    for !db.EOF() {
        name, err := db.FieldValueAsString(1)
        if err != nil {
            return
        }
        salary, err := db.FieldValueAsFloat(2)
        if err != nil {
            return
        }
        bDate, err := db.FieldValueAsDate(3)
        if err != nil {
            return
        }
        fmt.Println(name, salary, bDate)
        if err := db.Next(); err != nil {
            return
        }
    }
