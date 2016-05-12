package users

import (
	"database/sql"
	"os"
	"testing"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func setup() SQLiteUserDB {
	os.Remove("./foo.db")
	db, err := sql.Open("sqlite3", "./foo.db")
	check(err)
	return SQLiteUserDB{db}
}

func cleanup(db *SQLiteUserDB) {
	db.Close()
	os.Remove("./foo.db")
}

func TestDataBaseCRUD(t *testing.T) {
	var udb UserDB
	var sqldb SQLiteUserDB
	sqldb = setup()
	udb = &sqldb
	defer cleanup(&sqldb)

	creds := UserCredentials{username: "foo",
		password: "This is such!_asecUR3 Passw0rd?"}
	creds2 := UserCredentials{username: "alice or 1=1'*@;; FROM sqlimebro",
		password: "baz"}

	udb.CreateUserTable()

	// Create
	err := udb.CreateUser(creds)
	check(err)
	err = udb.CreateUser(creds2)
	check(err)

	exists, err := udb.UserExists(creds.username)
	check(err)
	if exists == false {
		t.Fail()
	}

	no_exists, err := udb.UserExists("nonexistent user")
	check(err)
	if no_exists == true {
		t.Fail()
	}

	// Read(Verify)
	res := udb.VerifyCredentials(creds)
	if res != true {
		t.Fail()
	}

	res2 := udb.VerifyCredentials(creds2)
	if res2 != true {
		t.Fail()
	}
	res3 := udb.VerifyCredentials(UserCredentials{username: "yo",
		password: "dog"})
	if res3 == true {
		t.Fail()
	}

	// Update
	new_creds := UserCredentials{username: "Noam",
		password: "Chomsky"}
	err2 := udb.UpdateUser(creds, new_creds)
	check(err2)
	res4 := udb.VerifyCredentials(new_creds)
	if res4 != true {
		t.Fail()
	}

	err3 := udb.UpdateUser(creds, creds2) // `creds` should be invalid
	if err3 == nil {
		t.Fail()
	}

	// Delete
	err4 := udb.DeleteUser(new_creds)
	if err4 != nil {
		t.Fail()
	}
	res5 := udb.VerifyCredentials(new_creds)
	if res5 == true {
		t.Fail()
	}
}
