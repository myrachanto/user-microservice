package controllers

import (
	"encoding/json"
	"testing"

	"github.com/myrachanto/microservice/user/src/model"
	"github.com/myrachanto/microservice/user/src/service"
	"github.com/stretchr/testify/assert"
)

// end to end testing
var (
	jsondata = `{"FName":"Jane","LName":"Doe","UName":"doe","Usercode": "Doe345","Phone":"1234567","Email":   "email@example.com","Password": "1234567","Address":"psd 456 king view"
	}`
	jsondata1 = `{"Email":   "email@example.com","Password": "1234567"}`
)

//make end to end testing
func TestCreateUser(t *testing.T) {

	user := &model.User{}
	if err := json.Unmarshal([]byte(jsondata), &user); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	res, err := service.UserService.Create(user)
	expected := "user created successifully"
	if res.UName != "doe" {
		assert.EqualValues(t, expected, res, "Something went wrong testing the Create method")
		assert.NotNil(t, err, "Something went wrong testing the Create method")
	}
	// afterparty cleaner
	//use uname to clean after testing
	service.UserService.Cleaner(user.UName)
}
func TestLoginUser(t *testing.T) {

	user1 := &model.User{}
	if err := json.Unmarshal([]byte(jsondata), &user1); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	_, _ = service.UserService.Create(user1)
	user := &model.LoginUser{}
	if err := json.Unmarshal([]byte(jsondata1), &user); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	res, err := service.UserService.Login(user)
	// expected := "user created successifully"
	if res.Token == "" {
		assert.EqualValues(t, "", res, "Something went wrong testing with the Logging method")
		assert.NotNil(t, err, "Something went wrong testing with the Logging method")
	}
	// afterparty cleaner
	//use uname to clean after testing
	service.UserService.Cleaner(user1.UName)
}
func TestGetAllUser(t *testing.T) {

	user1 := &model.User{}
	if err := json.Unmarshal([]byte(jsondata), &user1); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	_, _ = service.UserService.Create(user1)
	user := &model.LoginUser{}
	if err := json.Unmarshal([]byte(jsondata1), &user); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	res, err := service.UserService.GetAll("", 1, 12)
	// expected := "user created successifully"
	if res[0].UName == "" {
		assert.EqualValues(t, "", res, "Something went wrong testing with the Getting all method")
		assert.NotNil(t, err, "Something went wrong testing with the Getting all method")
	}
	// afterparty cleaner
	//use uname to clean after testing
	service.UserService.Cleaner(user1.UName)
}
func TestGetOneUser(t *testing.T) {

	user1 := &model.User{}
	if err := json.Unmarshal([]byte(jsondata), &user1); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	person, _ := service.UserService.Create(user1)
	user := &model.LoginUser{}
	if err := json.Unmarshal([]byte(jsondata1), &user); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	res, err := service.UserService.GetOne(int(person.ID))
	// expected := "user created successifully"
	if res.UName == "" {
		assert.EqualValues(t, "", res, "Something went wrong testing with the Getting one method")
		assert.NotNil(t, err, "Something went wrong testing with the Getting one method")
	}
	// afterparty cleaner
	//use uname to clean after testing
	service.UserService.Cleaner(user1.UName)
}

func TestUpdateUser(t *testing.T) {

	user1 := &model.User{}
	if err := json.Unmarshal([]byte(jsondata), &user1); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	person, _ := service.UserService.Create(user1)
	user := &model.LoginUser{}
	if err := json.Unmarshal([]byte(jsondata1), &user); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	person.FName = "John"
	res, err := service.UserService.Update(int(person.ID), person)
	// expected := "user created successifully"
	if res.FName == "jane" {
		assert.EqualValues(t, "", res.FName, "Something went wrong testing with the Updating one method")
		assert.NotNil(t, err, "Something went wrong testing with the Updating one method")
	}
	// afterparty cleaner
	//use uname to clean after testing
	service.UserService.Cleaner(user1.UName)
}
func TestDeleteUser(t *testing.T) {

	user1 := &model.User{}
	if err := json.Unmarshal([]byte(jsondata), &user1); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	person, _ := service.UserService.Create(user1)
	user := &model.LoginUser{}
	if err := json.Unmarshal([]byte(jsondata1), &user); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	person.FName = "John"
	res, err := service.UserService.Delete(int(person.ID))
	expected := "deleted successfully"
	if res != expected {
		assert.EqualValues(t, "", res, "Something went wrong testing with the Deleting method")
		assert.NotNil(t, err, "Something went wrong testing with the Deleting method")
	}
}
