package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// UserRest is REST controller used to handler request about the user entities.
type UserRest struct {
	db *sql.DB
}

// NewUserRest constructs new UserRest controller used to handle requests about the user entities.
func NewUserRest(db *sql.DB) *UserRest {
	return &UserRest{db: db}
}

// Users handler responds with list of all users in system.
func (r *UserRest) Users(w http.ResponseWriter, req *http.Request) {
	rows, err := r.db.Query("SELECT * FROM users")
	if err != nil {
		WriteErr(w, err, http.StatusInternalServerError)
		return
	}

	users := make([]User, 0, 0)
	for rows.Next() {
		var uid int
		var firstName string
		var secondName string
		var birthDate time.Time
		err = rows.Scan(&uid, &firstName, &secondName, &birthDate)
		if err != nil {
			WriteErr(w, err, http.StatusInternalServerError)
			return
		}

		users = append(users, User{ID: uid, FirstName: firstName, SecondName: secondName, BirthDate: birthDate})
	}

	if rows.Err() != nil {
		WriteErr(w, rows.Err(), http.StatusInternalServerError)
		return
	}

	WriteJSON(w, users)
}

// GetUser handler responds with particular user based on the id of the user.
func (r *UserRest) GetUser(w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	if id == "" {
		WriteErr(w, errors.New("please provide User id"), http.StatusBadRequest)
		return
	}

	row := r.db.QueryRow("SELECT * FROM users WHERE id = ?", id)
	var uid int
	var firstName string
	var secondName string
	var birthDate time.Time
	err := row.Scan(&uid, &firstName, &secondName, &birthDate)
	if err != nil {
		if err == sql.ErrNoRows {
			WriteErr(w, errors.New("can't find user with id: "+id), http.StatusNotFound)
			return
		}

		WriteErr(w, err, http.StatusInternalServerError)
		return
	}

	user := User{ID: uid, FirstName: firstName, SecondName: secondName, BirthDate: birthDate}

	WriteJSON(w, user)
}

// AddUser handler is responsible for adding new user to system.
func (r *UserRest) AddUser(w http.ResponseWriter, req *http.Request) {
	var user User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		WriteErr(w, errors.New("can't deserialize json with user"), http.StatusBadRequest)
		return
	}

	_, err = r.db.Exec("INSERT INTO users(firstName, secondName, birthDate) values(?,?,?)", user.FirstName, user.SecondName, user.BirthDate)
	if err != nil {
		WriteErr(w, err, http.StatusInternalServerError)
		return
	}
}

// User holds basic info about the User entity.
type User struct {
	ID         int       `json:"id"`
	FirstName  string    `json:"firstName"`
	SecondName string    `json:"secondName"`
	BirthDate  time.Time `json:"birthDate"`
}
