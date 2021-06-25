package models

import (
	"database/sql"
	"errors"
	"fmt"
	"html"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/f-chilmi/just-text-go/db"
	"github.com/f-chilmi/just-text-go/helpers"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Phone     string    `json:"phone"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Token struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
}

var secretkey = os.Getenv("SECRET_KEY")

func (u *User) Prepare() {
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Phone = html.EscapeString(strings.TrimSpace(u.Phone))
	u.Password = html.EscapeString(strings.TrimSpace(u.Password))
}

func (u *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "update":
		switch "" {
		case u.Username:
			return errors.New("required username")
		case u.Password:
			return errors.New("required password")
		case u.Phone:
			return errors.New("required phone")
		default:
			return nil
		}

	case "login":
		switch "" {
		case u.Password:
			return errors.New("required password")
		case u.Phone:
			return errors.New("required phone")
		default:
			return nil
		}

	case "register":
		switch "" {
		case u.Username:
			return errors.New("required username")
		case u.Password:
			return errors.New("required password")
		case u.Phone:
			return errors.New("required phone")
		default:
			return nil
		}

	default:
		switch "" {
		case u.Username:
			return errors.New("required username")
		case u.Password:
			return errors.New("required password")
		case u.Phone:
			return errors.New("required phone")
		default:
			return nil
		}
	}
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func HashAndSalt(pwd []byte) (string, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), err
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}

func GeneratehashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func GenerateJWT(id int64, username string, phone string) (string, error) {
	var mySigningKey = []byte(secretkey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["id"] = id
	claims["username"] = username
	claims["phone"] = phone
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (u *User) Login(phone string, password string) (string, error) {
	var err error

	user := User{}

	db := db.CreateConnection()

	err = db.QueryRow(`
		SELECT id, 
		username, 
		phone,
		password,
		created_at,
		updated_at
		FROM users WHERE phone=$1
		`, phone).
		Scan(&user.ID, &user.Username, &user.Phone, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return "", err
	}

	err = CheckPasswordHash(password, user.Password)
	if err != nil {
		return "", err
	}

	validToken, err := GenerateJWT(user.ID, user.Username, user.Phone)
	if err != nil {
		return "", err
	}

	return validToken, err

}

func (u *User) Register(username string, phone string, password string) (string, error) {
	var err error

	// user := User{}
	var id int64

	db := db.CreateConnection()
	err = db.QueryRow(`SELECT id FROM users WHERE phone=$1`, phone).Scan(&id)

	switch err {
	// if no user found, so user can create new (register)
	case sql.ErrNoRows:
		newP, err := GeneratehashPassword(password)

		if err != nil {
			return "", err
		}

		newU := User{
			Username: username,
			Phone:    phone,
			Password: newP,
		}
		_, err = u.InsertUser(newU)
		if err != nil {
			return "", err
		}
		return "", err

	// if user found
	case nil:
		return "user already exist", errors.New("user already exist")

	default:
		break
	}

	return "success", err
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

func (u *User) GetUserByPhone(phone string) (User, error) {
	// create the postgres db connection
	db := db.CreateConnection()

	// close the db connection
	defer db.Close()

	var user User

	// create the select sql query
	sqlStatement := `SELECT * FROM users WHERE phone=$1`

	// execute the sql statement
	row := db.QueryRow(sqlStatement, phone)

	err := row.Scan(&user.ID, &user.Username, &user.Phone, &user.Password, &user.CreatedAt, &user.UpdatedAt)

	// return empty user on error
	return user, err
}
