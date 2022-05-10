package controllers

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/myrachanto/microservice/user/src/model"
	"github.com/myrachanto/microservice/user/src/repository"
	"github.com/myrachanto/microservice/user/src/service"

	"github.com/stretchr/testify/assert"
)

// end to end testing
var (
	jsondata = `{"f_name":"Jane","l_name":"Doe","u_name":"doe","usercode": "Doe345","phone":"1234567","email":"email@example1.com","password":"1234567","address":"psd 456 king view"
	}`
	jsondata1 = `{"email":   "email@example1.com","password": "1234567"}`
)

//make end to end testing
func TestCreateUser(t *testing.T) {
	repository.IndexRepo.InitDB()
	user := &model.User{}
	if err := json.Unmarshal([]byte(jsondata), &user); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	// fmt.Println(">>>>>>>>>", user)
	u, err := service.UserService.Create(user)
	fmt.Println(">>>>>>>>>", u)
	assert.EqualValues(t, "Jane", u.FName, "failed to validate create method")
	assert.Nil(t, err)
	user2 := &model.User{}
	if err := json.Unmarshal([]byte(jsondata), &user2); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	u2, err2 := service.UserService.Create(user2)
	assert.EqualValues(t, "Your email already exists!", err2.Message(), "failed to validate create method")
	assert.Nil(t, u2)
	// afterparty cleaner
	//use uname to clean after testing
	service.UserService.Cleaner(u.Usercode)
}

func TestGetAllUser(t *testing.T) {

	user1 := &model.User{}
	if err := json.Unmarshal([]byte(jsondata), &user1); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	_, _ = service.UserService.Create(user1)
	_, err := service.UserService.GetAll("", 1, 12)
	assert.Nil(t, err, "Something went wrong testing with the Getting all method")
	// afterparty cleaner
	//use uname to clean after testing
	service.UserService.Cleaner(user1.Usercode)
}
func TestGetOneUser(t *testing.T) {

	user1 := &model.User{}
	if err := json.Unmarshal([]byte(jsondata), &user1); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	person, _ := service.UserService.Create(user1)

	_, e := service.UserService.GetOne(person.Usercode)
	// assert.EqualValues(t, "doe", u.UName, "Something went wrong testing with the Getting one method")
	assert.Nil(t, e, "Something went wrong testing with the Getting one method")
	// afterparty cleaner
	//use uname to clean after testing
	service.UserService.Cleaner(user1.Usercode)
}

func TestUpdateUser(t *testing.T) {
	user1 := &model.User{}
	if err := json.Unmarshal([]byte(jsondata), &user1); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	person, _ := service.UserService.Create(user1)
	// fmt.Println(">>>>>>>>>>>>>>sd", person)
	person.FName = "John"
	// fmt.Println(">>>>>>>>>>>>>>", person)
	u, err := service.UserService.Update(person.Usercode, person)
	// expected := "user created successifully"
	assert.EqualValues(t, "John", u.FName, "Something went wrong testing with the Getting one method")
	assert.Nil(t, err)
	// afterparty cleaner
	//use uname to clean after testing
	service.UserService.Cleaner(user1.Usercode)
}
func TestDeleteUser(t *testing.T) {

	user1 := &model.User{}
	if err := json.Unmarshal([]byte(jsondata), &user1); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	person, _ := service.UserService.Create(user1)
	res, err := service.UserService.Delete(int(person.ID))
	expected := "deleted successfully"
	assert.EqualValues(t, expected, res, "Something went wrong testing with the Deleting method")
	assert.Nil(t, err, "Something went wrong testing with the Deleting method")
}
