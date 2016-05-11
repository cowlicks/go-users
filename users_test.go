package users

import (
    "testing"
    "os"
    "database/sql"
)

func cleanDB(db * sql.DB) {
    db.Close()
    os.Remove("./foo.db")
}

func check(err error) {
    if err != nil {
        panic(err)
    }
}


func TestDataBaseCRUD(t * testing.T) {
    os.Remove("./foo.db")
    db, err := sql.Open("sqlite3", "./foo.db")
    check(err)
    defer cleanDB(db)

    creds := UserCredentials{username: "foo",
                             password: "This is such!_asecUR3 Passw0rd?"}
    creds2 := UserCredentials{username: "alice or 1=1'*@;; FROM sqlimebro",
                              password: "baz"}

    CreateUserTable(db)

    // Create
    err = CreateUser(db, creds)
    check(err)
    err = CreateUser(db, creds2)
    check(err)

    exists, err := UserExists(db, creds.username)
    check(err)
    if exists == false {
       t.Fail()
   }

    no_exists, err := UserExists(db, "nonexistent user")
    check(err)
    if no_exists == true {
        t.Fail()
    }

    // Read(Verify)
    res := VerifyCredentials(db, creds)
    if res != true {
        t.Fail()
    }

    res2 := VerifyCredentials(db, creds2)
    if res2 != true {
        t.Fail()
    }
    res3 := VerifyCredentials(db, UserCredentials{username: "yo",
                                                  password: "dog"})
    if res3 == true {
        t.Fail()
    }

    // Update
    new_creds := UserCredentials{username: "Noam",
                                 password: "Chomsky"}
    err2 := UpdateUser(db, creds, new_creds)
    check(err2)
    res4 := VerifyCredentials(db, new_creds)
    if res4 != true {
        t.Fail()
    }

    err3 := UpdateUser(db, creds, creds2)  // `creds` should be invalid
    if err3 == nil {
        t.Fail()
    }

    // Delete
    err4 := DeleteUser(db, new_creds)
    if err4 != nil {
        t.Fail()
    }
    res5 := VerifyCredentials(db, new_creds)
    if res5 == true {
        t.Fail()
    }
}
