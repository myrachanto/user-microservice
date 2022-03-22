package service

import (
	"testing"

	httperors "github.com/myrachanto/custom-http-error"
	"github.com/myrachanto/microservice/user/src/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

var (
	User = &model.User{
		FName:    "jane",
		LName:    "Doe",
		UName:    "doe123",
		Usercode: "Doe345",
		Phone:    "1234567",
		Email:    "email@example.com",
		Password: "1234567",
		Address:  "psd 456 king view",
	}
)

type UserMockInterface interface {
	Create(User *model.User) (*model.User, httperors.HttpErr)
	GetOne(id string) (*model.User, httperors.HttpErr)
	GetAll() ([]*model.User, httperors.HttpErr)
	Update(code string, User *model.User) httperors.HttpErr
}

func (mock MockRepository) Create(User *model.User) (*model.User, httperors.HttpErr) {
	args := mock.Called()
	result := args.Get(0)
	User, err := result.(*model.User), args.Error(1)
	if err != nil {
		return nil, httperors.NewNotFoundError("Something went wrong creating the resourse")
	}
	return User, nil
}
func (mock MockRepository) GetOne() (*model.User, httperors.HttpErr) {
	args := mock.Called()
	result := args.Get(0)
	User, err := result.(*model.User), args.Error(1)
	if err != nil {
		return nil, httperors.NewNotFoundError("Something went wrong getting the resourse")
	}
	return User, nil
}
func (mock MockRepository) GetAll() ([]*model.User, httperors.HttpErr) {
	args := mock.Called()
	result := args.Get(0)
	Users, err := result.([]*model.User), args.Error(1)
	if err != nil {
		return nil, httperors.NewNotFoundError("Something went wrong getting the resourses")
	}
	return Users, nil
}
func (mock MockRepository) Update(code string, user *model.User) httperors.HttpErr {
	args := mock.Called()
	result := args.Get(0)
	_, err := result.(*model.User), args.Error(1)
	if err != nil {
		return httperors.NewNotFoundError("Something went wrong updating the resourse")
	}
	return nil
}
func TestGetAll(t *testing.T) {
	mockRepo := new(MockRepository)
	mockRepo.On("Create").Return(User, nil)
	_, _ = mockRepo.Create(User)
	mockRepo.On("GetAll").Return([]*model.User{User}, nil)
	results, _ := mockRepo.GetAll()
	//mock assertion: behavioral
	mockRepo.AssertExpectations(t)
	//data assertion
	assert.Equal(t, User.FName, results[0].FName)
	assert.Equal(t, User.LName, results[0].LName)
	assert.Equal(t, User.UName, results[0].UName)
	assert.Equal(t, User.Phone, results[0].Phone)
	assert.Equal(t, User.Password, results[0].Password)
	assert.Equal(t, User.Email, results[0].Email)

}
func TestGetOne(t *testing.T) {
	mockRepo := new(MockRepository)
	mockRepo.On("Create").Return(User, nil)
	_, _ = mockRepo.Create(User)
	mockRepo.On("GetOne").Return(User, nil)
	results, _ := mockRepo.GetOne()
	//mock assertion: behavioral
	mockRepo.AssertExpectations(t)
	//data assertion
	assert.Equal(t, User.FName, results.FName)
	assert.Equal(t, User.LName, results.LName)
	assert.Equal(t, User.UName, results.UName)
	assert.Equal(t, User.Phone, results.Phone)
	assert.Equal(t, User.Password, results.Password)
	assert.Equal(t, User.Email, results.Email)

}
func TestCreate(t *testing.T) {
	mockRepo := new(MockRepository)
	mockRepo.On("Create").Return(User, nil)
	result, err := mockRepo.Create(User)
	mockRepo.AssertExpectations(t)
	assert.Equal(t, User.FName, result.FName)
	assert.Equal(t, User.LName, result.LName)
	assert.Equal(t, User.UName, result.UName)
	assert.Equal(t, User.Phone, result.Phone)
	assert.Equal(t, User.Password, result.Password)
	assert.Equal(t, User.Email, result.Email)
	assert.Nil(t, err)

}
