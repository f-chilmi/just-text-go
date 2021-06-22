package models

import (
	"fmt"
	"time"

	"github.com/f-chilmi/just-text-go/db"
	"github.com/f-chilmi/just-text-go/helpers"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Phone     string    `json:"phone"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) InsertUser(user User) (int64, error) {
	// create the db connection
	db := db.CreateConnection()

	// close the db connection
	defer db.Close()

	// create the insert query
	// returning userid will return the id of the inserted user
	sqlStatement := `INSERT INTO users (username, phone, password) VALUES ($1, $2, $3) RETURNING id;`

	// inserted id will store in this id
	var id int64

	// execute the sql statement
	// scan function will save the inserted id in the id
	err := db.QueryRow(sqlStatement, user.Username, user.Phone, user.Password).Scan(&id)

	helpers.CheckError("Unable to execute the query.", err)

	fmt.Printf("Inserted a single record %v", id)

	// return the inserted id
	return id, err
}

func (u *User) GetUsers() ([]User, error) {
	// create the postgres db connection
	db := db.CreateConnection()

	// close the db connection
	defer db.Close()

	var users []User

	// create the select sql query
	sqlStatement := `SELECT * FROM users`

	// execute the sql statement
	rows, err := db.Query(sqlStatement)

	helpers.CheckError("Unable to execute the query.", err)

	// close the statement
	defer rows.Close()

	// iterate over the rows
	for rows.Next() {
		var user User

		// unmarshal the row object to user
		err := rows.Scan(&user.ID, &user.Username, &user.Phone, &user.Password, &user.CreatedAt, &user.UpdatedAt)

		helpers.CheckError("Unable to scan the row.", err)

		// append the user in the users slice
		users = append(users, user)

	}

	// return empty user on error
	return users, err
}

func (u *User) GetUser(id int64) (User, error) {
	// create the postgres db connection
	db := db.CreateConnection()

	// close the db connection
	defer db.Close()

	// create a user of models.User type
	var user User

	// create the select sql query
	sqlStatement := `SELECT * FROM users WHERE id=$1`

	// execute the sql statement
	row := db.QueryRow(sqlStatement, id)

	// unmarshal the row object to user
	err := row.Scan(&user.ID, &user.Username, &user.Phone, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	// return empty user on error
	return user, err
}

func (u *User) UpdateUser(id int64, user User) int64 {

	// create the postgres db connection
	db := db.CreateConnection()

	// close the db connection
	defer db.Close()

	// create the update sql query
	sqlStatement := `UPDATE users SET username=$2, phone=$3, updated_at=CURRENT_TIMESTAMP WHERE id=$1`

	// execute the sql statement
	res, err := db.Exec(sqlStatement, id, user.Username, user.Phone)

	helpers.CheckError("Unable to execute the query", err)

	// check how many rows affected
	rowsAffected, err := res.RowsAffected()

	helpers.CheckError("Error while cheking the affected rows", err)

	fmt.Printf("Total rows/record affected %v", rowsAffected)

	return rowsAffected
}
