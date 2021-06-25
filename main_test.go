package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllUsers(t *testing.T) {

	res, err := http.Get("http://localhost:8080/users")
	testError(err, t)

	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode, "http.StatusOK response expected")
	expected := `[{"id":"1","name":"User 1","gender":"Male","email":"o.goyal@","address":"Alwar"},{"id":"101","name":"vaibhav","gender":"Male","email":"o.goyal@","address":"Alwar"},{"id":"1601","name":"chirag -16","gender":"Male","email":"c.goyal@","address":"rajasthan"},{"id":"1602","name":"chirag -1602","gender":"Male","email":"c.goyal@","address":"rajasthan"},{"id":"67","name":"rinku-67","gender":"Male","email":"o.goyal@","address":"Alwar"},{"id":"67","name":"rinku-67","gender":"Male","email":"o.goyal@","address":"Alwar"},{"id":"61","name":"rinku-61","gender":"Male","email":"o.goyal@","address":"Alwar"},{"id":"62","name":"rinku-62","gender":"Male","email":"o.goyal@","address":"Alwar"},{"id":"64","name":"rinku-64","gender":"Male","email":"o.goyal@","address":"Alwar"},{"id":"64","name":"rinku-64","gender":"Male","email":"o.goyal@","address":"Alwar"}]`
	assert.Equal(t, string(data), expected, "result not matched")
}

func TestGetUserByID(t *testing.T) {

	res, err := http.Get("http://localhost:8080/users/1")
	testError(err, t)
	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	assert.Equal(t, http.StatusOK, res.StatusCode, "http.StatusOK response expected")
	expected := `[{"id":"1","name":"User 1","gender":"Male","email":"o.goyal@","address":"Alwar"}]`
	assert.Equal(t, string(data), expected, "result not matched")
}

func TestCreateUser(t *testing.T) {
	var jsonStr = []byte(`{"id":"11","name:"xyz","gender":"Male","email":"xyz@pqr.com","address":"Jaipur"}`)

	req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(jsonStr))
	testError(err, t)

	rr := httptest.NewRecorder()

	createUser(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusCreated)
	}
}

func TestDeleteUser(t *testing.T) {

	req, err := http.NewRequest("DELETE", "/user/1", nil)
	testError(err, t)

	rr := httptest.NewRecorder()

	deleteUser(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestUpdateUser(t *testing.T) {

	var jsonStr = []byte(`{"id":101,"name":"xyz change","gender":"female","email":"xyz@pqr.com","address":"delhi"}`)

	req, err := http.NewRequest("PATCH", "/user/1", bytes.NewBuffer(jsonStr))
	testError(err, t)

	rr := httptest.NewRecorder()
	updateUser(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func testError(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}
