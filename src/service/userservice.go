package service

import (
	"fmt"
	"log"

	httperors "github.com/myrachanto/custom-http-error"
	"github.com/myrachanto/microservice/user/src/model"
	r "github.com/myrachanto/microservice/user/src/repository"
)

//UserService ...
var (
	UserService UserServiceInterface = &userService{}
)

type UserServiceInterface interface {
	Create(user *model.User) (*model.User, httperors.HttpErr)
	Login(auser *model.LoginUser) (*model.Auth, httperors.HttpErr)
	Logout(token string) httperors.HttpErr
	GetOne(code int) (*model.User, httperors.HttpErr)
	GetAll(search string, page, pagesize int) ([]model.User, httperors.HttpErr)
	Update(id int, user *model.User) (*model.User, httperors.HttpErr)
	Delete(id int) (string, httperors.HttpErr)
	Cleaner(id string) (string, httperors.HttpErr)
}

type userService struct {
	repository r.UserRepoInterface
}

func NewUserService(repo r.UserRepoInterface) UserServiceInterface {
	return &userService{
		repo,
	}
}

func (service userService) Create(user *model.User) (*model.User, httperors.HttpErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	log.Println("service level")
	s, err1 := r.Userrepo.Create(user)
	if err1 != nil {
		return nil, err1
	}
	return s, nil

}
func (service userService) Login(auser *model.LoginUser) (*model.Auth, httperors.HttpErr) {
	user, err1 := r.Userrepo.Login(auser)
	if err1 != nil {
		return nil, err1
	}
	return user, nil
}

// func (service userService) Employee_login(auser *model.LoginEmployee) (*model.Auth, *httperors.HttpError) {
// 	user, err1 := r.Userrepo.Employee_login(auser)
// 	if err1 != nil {
// 		return nil, err1
// 	}
// 	return user, nil
// }
func (service userService) Logout(token string) httperors.HttpErr {
	err1 := r.Userrepo.Logout(token)
	if err1 != nil {
		return err1
	}
	return nil
}

// func (service userService) Forgot(email string) *httperors.HttpError {
// 	email, pass, err1 := r.Userrepo.Forgot(email)
// 	// email.EmailingPassword(pass, email)
// 	// email.
// 	e.EmailingPassword(pass, email)
// 	return err1
// }

// func (service userService) PasswordUpdate(oldpassword, email, newpassword string) *httperors.HttpError {
// 	email, pass, err1 := r.Userrepo.PasswordUpdate(oldpassword, email, newpassword)
// 	e.EmailingPassword(pass, email)
// 	return err1
// }

// func (service userService) AdminUpdate(code, status string) *httperors.HttpError {
// 	err1 := r.Userrepository.AdminUpdate(code, status)
// 	return err1
// }
func (service userService) GetOne(code int) (*model.User, httperors.HttpErr) {
	user, err1 := r.Userrepo.GetOne(code)
	if err1 != nil {
		return nil, err1
	}
	return user, nil
}

func (service userService) GetAll(search string, page, pagesize int) ([]model.User, httperors.HttpErr) {
	results, err := r.Userrepo.GetAll(search, page, pagesize)
	return results, err
}

// func (service userService) UpdateRole(code, admin, supervisor, employee, level, usercode string) (string, *httperors.HttpError) {
// 	user, err1 := r.Userrepo.UpdateRole(code, admin, supervisor, employee, level, usercode)
// 	return user, err1
// }

func (service userService) Update(id int, user *model.User) (*model.User, httperors.HttpErr) {
	fmt.Println("update1-controller")
	fmt.Println(id)
	user, err1 := r.Userrepo.Update(id, user)
	if err1 != nil {
		return nil, err1
	}

	return user, nil
}
func (service userService) Delete(id int) (string, httperors.HttpErr) {
	success, failure := r.Userrepo.Delete(id)
	return success, failure
}
func (service userService) Cleaner(id string) (string, httperors.HttpErr) {
	success, failure := r.Userrepo.Cleaner(id)
	return success, failure
}
