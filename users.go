/*
An interface to a database for managing user data.

Usernames must be unique.
*/
package users

import (
    "errors"
    "database/sql"
    _"github.com/mattn/go-sqlite3"
    "golang.org/x/crypto/bcrypt"
)

func Check(err error) {
    if err != nil {
        panic(err)
    }
}

type UserCredentials struct {
    username string
    password string
}

func CreateUserTable(db * sql.DB) {
    sqlStmt := `
    CREATE TABLE IF NOT EXISTS userinfo(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE,
		password TEXT
	);
    `
    _, err := db.Exec(sqlStmt)
    Check(err)
}

func UserExists(db * sql.DB, username string) (bool, error) {
    var un string
    sqlStmt := `SELECT username FROM userinfo WHERE username = ?`
    err := db.QueryRow(sqlStmt, username).Scan(&un)

    switch {
    case err == sql.ErrNoRows: // no rows matched username, doesn't exist
        return false, nil
    case err != nil:           // a real error
        return false, err
    default:
        return true, nil       // user exists
    }
}

func CreateUser(db * sql.DB, uc UserCredentials) error {
    // Check username isn't taken
    exists, err := UserExists(db, uc.username)
    if exists {
        return errors.New("Username taken")
    }

    // hash password
    hashed_pw, err := bcrypt.GenerateFromPassword([]byte(uc.password),
                                                  bcrypt.DefaultCost)
    Check(err)

    sqlStmt := `INSERT INTO userinfo(username, password) values(?,?)`
    stmt, err := db.Prepare(sqlStmt)
    Check(err)

    // User
    _, err = stmt.Exec(uc.username, hashed_pw)
    Check(err)
    return nil

}

func VerifyCredentials(db * sql.DB, uc UserCredentials) bool {
    var un string
    var pw string
    sqlStmt :=`SELECT username, password FROM userinfo WHERE username = ?`

    err := db.QueryRow(sqlStmt, uc.username).Scan(&un, &pw)
    if err != nil {
        return false
    }

    err2 := bcrypt.CompareHashAndPassword([]byte(pw), []byte(uc.password))
    if err2 != nil {
        return false
    } else {
        return true
    }
}

func UpdateUser(db * sql.DB, old_creds, new_creds UserCredentials) error {
    updateStmt := `UPDATE userinfo SET username = ?, password = ? WHERE username = ?`

    verified := VerifyCredentials(db, old_creds)
    if !verified {
        return errors.New("The supplied credentials match no existing user.")
    }

    hashed_pw, err2 := bcrypt.GenerateFromPassword([]byte(new_creds.password),
                                                   bcrypt.DefaultCost)
    if err2 != nil {
        return err2
    }

    _, err3:= db.Exec(updateStmt, new_creds.username, hashed_pw, old_creds.username)
    if err3 != nil {
        return err3
    }
    return nil
}

func DeleteUser(db * sql.DB, uc UserCredentials) error {
    deleteStmt := `DELETE FROM userinfo WHERE username = ?`
    verified := VerifyCredentials(db, uc)
    if !verified {
        return errors.New("The supplied credentials match no existing user.")
    }
    _, err := db.Exec(deleteStmt, uc.username)
    return err
}
