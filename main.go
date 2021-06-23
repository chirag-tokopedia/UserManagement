package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	userpg   = "postgres"
	password = "<password>"
	dbname   = "usermanagment"
)

type user struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Gender  string `json:"gender"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

type allUsers []user

var users allUsers

// API for checking the connection with server.
func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to New Psql HTTP Sever 8080 for CRUD Users!!")
}

func getUsers() allUsers {
	db := openConnection()

	//select data
	rows, err := db.Query(`SELECT id, name ,gender, email, address FROM users`)
	CheckError(err)

	defer rows.Close()
	defer db.Close()

	var tempUsers allUsers

	for rows.Next() {
		var user1 user
		rows.Scan(&user1.ID, &user1.Name, &user1.Gender, &user1.Email, &user1.Address)
		tempUsers = append(tempUsers, user1)
	}
	return tempUsers
}

// API for getting all users information.
func getAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	users = getUsers()
	json.NewEncoder(w).Encode(users)
}

// API for getting user information by user id.
func getUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := mux.Vars(r)["id"]

	db := openConnection()

	//select data
	rows, err := db.Query(`SELECT id, name ,gender, email, address FROM users WHERE id=$1`, userID)
	CheckError(err)

	defer rows.Close()
	defer db.Close()

	var tempUsers allUsers

	for rows.Next() {
		var user1 user
		rows.Scan(&user1.ID, &user1.Name, &user1.Gender, &user1.Email, &user1.Address)
		tempUsers = append(tempUsers, user1)
	}
	json.NewEncoder(w).Encode(tempUsers)
}

// API for deleting a user by user id.
func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	eventID := mux.Vars(r)["id"]

	db := openConnection()

	// Delete User
	deleteUser := `DELETE FROM users WHERE id=$1`
	_, d := db.Exec(deleteUser, eventID)
	CheckError(d)

	users = getUsers()

	defer db.Close()

	//Users after deletion
	json.NewEncoder(w).Encode(users)
}

// API for updating the parital user information
func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	eventID := mux.Vars(r)["id"]
	var updatedUser user

	reqBody, err := ioutil.ReadAll(r.Body)
	fatalError(err)

	json.Unmarshal(reqBody, &updatedUser)

	db := openConnection()

	// update data
	updateUser := `UPDATE users set name=$1, gender=$2, email=$3, address=$4  where id=$5`
	_, e := db.Exec(updateUser, updatedUser.Name, updatedUser.Gender, updatedUser.Email, updatedUser.Address, eventID)
	CheckError(e)

	defer db.Close()

	users = getUsers()

	json.NewEncoder(w).Encode(users)
}

// API for creating new user.
func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newUser user
	reqBody, err := ioutil.ReadAll(r.Body)
	fatalError(err)

	json.Unmarshal(reqBody, &newUser)

	db := openConnection()

	//insert data
	insertDynStmt := `INSERT into users(id, name, gender, email, address) values($1, $2, $3, $4, $5)`
	_, i := db.Exec(insertDynStmt, newUser.ID, newUser.Name, newUser.Gender, newUser.Email, newUser.Address)
	CheckError(i)

	defer db.Close()

	w.WriteHeader(http.StatusCreated)
	users = getUsers()

	json.NewEncoder(w).Encode(newUser)
}

func openConnection() *sql.DB {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, userpg, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// check db
	err = db.Ping()
	CheckError(err)

	fmt.Println("Connected!")
	return db
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/user", createUser).Methods("POST")
	router.HandleFunc("/users/{id}", getUserById).Methods("GET")
	router.HandleFunc("/users", getAllUsers).Methods("GET")
	router.HandleFunc("/users/{id}", updateUser).Methods("PATCH")
	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func fatalError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
