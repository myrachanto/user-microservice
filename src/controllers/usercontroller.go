package controllers

import (
	"fmt"
	"strconv"

	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	httperors "github.com/myrachanto/custom-http-error"
	"github.com/myrachanto/microservice/user/src/model"
	"github.com/myrachanto/microservice/user/src/service"
)

//UserController ..
var (
	UserController UsercontrollerInterface = &userController{}
)

type userController struct {
	service service.UserServiceInterface
}
type UsercontrollerInterface interface {
	Create(c echo.Context) error
	Login(c echo.Context) error
	Logout(c echo.Context) error
	GetAll(c echo.Context) error
	Update(c echo.Context) error
	Delete(c echo.Context) error
}

func NewUserController(ser service.UserServiceInterface) UsercontrollerInterface {
	return &userController{
		ser,
	}
}

/////////controllers/////////////////
func (controller userController) Create(c echo.Context) error {
	user := &model.User{}
	user.FName = c.FormValue("fname")
	user.LName = c.FormValue("lname")
	user.UName = c.FormValue("uname")
	user.Phone = c.FormValue("phone")
	user.Address = c.FormValue("address")
	user.Email = c.FormValue("email")
	user.Password = c.FormValue("password")
	user.Business = c.FormValue("business")

	pic, err2 := c.FormFile("picture")
	//    fmt.Println(pic.Filename)
	if err2 != nil {
		httperror := httperors.NewBadRequestError("Invalid picture")
		return c.JSON(httperror.Code(), err2)
	}
	src, err := pic.Open()
	if err != nil {
		httperror := httperors.NewBadRequestError("the picture is corrupted")
		return c.JSON(httperror.Code(), err)
	}
	defer src.Close()
	// filePath := "./public/imgs/users/"
	filePath := "./public/imgs/users/" + pic.Filename
	filePath1 := "/imgs/users/" + pic.Filename
	// Destination
	dst, err4 := os.Create(filePath)
	if err4 != nil {
		httperror := httperors.NewBadRequestError("the Directory mess")
		return c.JSON(httperror.Code(), err4)
	}
	defer dst.Close()
	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		if err2 != nil {
			httperror := httperors.NewBadRequestError("error filling")
			return c.JSON(httperror.Code(), httperror)
		}
	}

	user.Picture = filePath1
	s, err1 := service.UserService.Create(user)
	if err1 != nil {
		return c.JSON(err1.Code(), err1)
	}
	if _, err = io.Copy(dst, src); err != nil {
		if err2 != nil {
			httperror := httperors.NewBadRequestError("error filling")
			return c.JSON(httperror.Code(), httperror)
		}
	}
	return c.JSON(http.StatusCreated, s)
}
func (controller userController) Login(c echo.Context) error {
	// service.UserService.Bizname = c.Get("bizname").(string)
	user := &model.LoginUser{}
	// auth := &model.Auth{}

	// alias := c.Param("alias")
	user.Email = c.FormValue("email")
	user.Password = c.FormValue("password")
	// alias := c.FormValue("alias")
	auth, problem := service.UserService.Login(user)
	if problem != nil {
		fmt.Println(problem)
		return c.JSON(problem.Code(), problem)
	}
	return c.JSON(http.StatusOK, auth)
}

// func (controller userController) Employee_login(c echo.Context) error {
// 	// service.UserService.Bizname = c.Get("bizname").(string)
// 	user := &model.LoginEmployee{}
// 	// auth := &model.Auth{}

// 	user.Email = c.FormValue("email")
// 	user.Password = c.FormValue("password")
// 	user.Alias = c.FormValue("alias")
// 	auth, problem := service.UserService.Employee_login(user)
// 	if problem != nil {
// 		fmt.Println(problem)
// 		return c.JSON(problem.Code, problem)
// 	}
// 	return c.JSON(http.StatusOK, auth)
// }
func (controller userController) Logout(c echo.Context) error {
	token := c.Param("token")
	problem := service.UserService.Logout(token)
	if problem != nil {
		return c.JSON(problem.Code(), problem)
	}
	return c.JSON(http.StatusOK, "succeessifully logged out")
}

// func (controller userController) Forgot(c echo.Context) error {
// 	email := c.FormValue("email")
// 	problem := service.UserService.Forgot(email)
// 	if problem != nil {
// 		return c.JSON(problem.Code, problem)
// 	}
// 	return c.JSON(http.StatusOK, "updated succesifully")
// }
// func (controller userController) PasswordUpdate(c echo.Context) error {
// 	oldpassword := c.FormValue("oldpassword")
// 	email := c.FormValue("email")
// 	newpassword := c.FormValue("newpassword")
// 	problem := service.UserService.PasswordUpdate(oldpassword, email, newpassword)
// 	if problem != nil {
// 		return c.JSON(problem.Code, problem)
// 	}
// 	return c.JSON(http.StatusOK, "updated succesifully")
// }

// func (controller userController) AdminUpdate(c echo.Context) error {
// 	code := c.Param("code")
// 	status := c.FormValue("admin")
// 	problem := service.UserService.AdminUpdate(code, status)
// 	if problem != nil {
// 		return c.JSON(problem.Code, problem)
// 	}
// 	return c.JSON(http.StatusOK, "updated succesifully")
// }
func (controller userController) GetAll(c echo.Context) error {
	search := string(c.QueryParam("q"))
	page, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid page number")
		return c.JSON(httperror.Code(), httperror)
	}
	pagesize, err := strconv.Atoi(c.QueryParam("pagesize"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid pagesize")
		return c.JSON(httperror.Code(), httperror)
	}

	results, err3 := service.UserService.GetAll(search, page, pagesize)
	if err3 != nil {
		return c.JSON(err3.Code(), err3)
	}
	return c.JSON(http.StatusOK, results)
}
func (controller userController) GetOne(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code(), httperror)
	}
	user, problem := service.UserService.GetOne(id)
	if problem != nil {
		return c.JSON(problem.Code(), problem)
	}
	return c.JSON(http.StatusOK, user)
}

// func (controller userController) UpdateRole(c echo.Context) error {
// 	service.UserService.Bizname = c.Get("bizname").(string)
// 	service.UserService.Central = c.Get("central").(string)
// 	code := c.Param("code")
// 	admin := c.FormValue("admin")
// 	supervisor := c.FormValue("supervisor")
// 	employee := c.FormValue("employee")
// 	level := c.FormValue("level")
// 	usercode := c.FormValue("usercode")
// 	updateduser, problem := service.UserService.UpdateRole(code, admin, supervisor, employee, level, usercode)
// 	if problem != nil {
// 		return c.JSON(problem.Code, problem)
// 	}
// 	return c.JSON(http.StatusOK, updateduser)
// }
func (controller userController) Update(c echo.Context) error {
	user := &model.User{}
	user.FName = c.FormValue("fname")
	user.LName = c.FormValue("lname")
	user.UName = c.FormValue("uname")
	user.Phone = c.FormValue("phone")
	user.Address = c.FormValue("address")
	user.Email = c.FormValue("email")
	user.Password = "njenga456"
	// user.Password = c.FormValue("password")
	// user.Admin = true

	pic, err2 := c.FormFile("picture")
	//    fmt.Println(pic.Filename)
	if err2 != nil {
		httperror := httperors.NewBadRequestError("Invalid picture")
		return c.JSON(httperror.Code(), err2)
	}
	src, err := pic.Open()
	if err != nil {
		httperror := httperors.NewBadRequestError("the picture is corrupted")
		return c.JSON(httperror.Code(), err)
	}
	defer src.Close()
	filePath := "./public/imgs/users/" + pic.Filename
	filePath1 := "/imgs/users/" + pic.Filename
	// Destination
	dst, err4 := os.Create(filePath)
	if err4 != nil {
		httperror := httperors.NewBadRequestError("the Directory mess")
		return c.JSON(httperror.Code(), err4)
	}
	defer dst.Close()
	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		if err2 != nil {
			httperror := httperors.NewBadRequestError("error filling")
			return c.JSON(httperror.Code(), httperror)
		}
	}

	user.Picture = filePath1
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code(), httperror)
	}
	updateduser, problem := service.UserService.Update(id, user)
	if problem != nil {
		return c.JSON(problem.Code(), problem)
	}
	if _, err = io.Copy(dst, src); err != nil {
		if err2 != nil {
			httperror := httperors.NewBadRequestError("error filling")
			return c.JSON(httperror.Code(), httperror)
		}
	}

	return c.JSON(http.StatusOK, updateduser)
}

func (controller userController) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		httperror := httperors.NewBadRequestError("Invalid ID")
		return c.JSON(httperror.Code(), httperror)
	}
	success, failure := service.UserService.Delete(id)
	if failure != nil {
		return c.JSON(failure.Code(), failure)
	}
	return c.JSON(http.StatusOK, success)

}
