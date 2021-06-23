package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type user struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Gender  string `json:"gender"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

type allUsers []user

var users = allUsers{
	{
		ID:      "1",
		Name:    "User 1",
		Gender:  "Male",
		Email:   "a.goyal@iiitmanipur.ac.in",
		Address: "Alwar",
	},
	{
		ID:      "2",
		Name:    "User 2",
		Gender:  "Male",
		Email:   "b.goyal@iiitmanipur.ac.in",
		Address: "Blwar",
	},
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to HTTP Sever 8080 for CRUD Users!!")
}

func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newUser user
	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "Kindly enter new user data")
	}

	json.Unmarshal(reqBody, &newUser)
	users = append(users, newUser)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newUser)
}

func getUserById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := mux.Vars(r)["id"]

	for _, user := range users {
		if user.ID == userID {
			json.NewEncoder(w).Encode(user)
		}
	}
}

func getAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	eventID := mux.Vars(r)["id"]

	for i, user := range users {
		if user.ID == eventID {
			users = append(users[:i], users[i+1:]...)
			fmt.Fprintf(w, "The user with ID %v has been deleted successfully, Now Final Users are Below:", eventID)
		}
	}
	json.NewEncoder(w).Encode(users)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	eventID := mux.Vars(r)["id"]
	var updatedUser user
	// var updateduser user = json.NewDecoder(r.body).Decode(&newUser)
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data to update User")
	}
	json.Unmarshal(reqBody, &updatedUser)

	for i, user := range users {
		if user.ID == eventID {
			user.Name = updatedUser.Name
			user.Gender = updatedUser.Gender
			user.Email = updatedUser.Email
			user.Address = updatedUser.Address
			users = append(users[:i], users[i+1:]...)
			users = append(users, user)
			json.NewEncoder(w).Encode(user)
		}
	}
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
