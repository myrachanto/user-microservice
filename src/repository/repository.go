package repository

import (
	"fmt"
	"log"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	httperors "github.com/myrachanto/custom-http-error"
	"github.com/myrachanto/microservice/user/src/model"
	"github.com/spf13/viper"
)

//Userrepo ...
var (
	Userrepo UserRepoInterface = &userrepo{}
)

type Key struct {
	EncryptionKey string `mapstructure:"EncryptionKey"`
}

func LoadKey() (key Key, err error) {
	viper.AddConfigPath("../../")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&key)
	return
}

type userrepo struct{}

type UserRepoInterface interface {
	Create(user *model.User) (*model.User, httperors.HttpErr)
	Login(user *model.LoginUser) (*model.Auth, httperors.HttpErr)
	Logout(token string) httperors.HttpErr
	all() (t []model.User, r httperors.HttpErr)
	GetOne(id int) (*model.User, httperors.HttpErr)
	userExistbycode(code string) bool
	userbycode(code string) *model.User
	GetAll(search string, page, pagesize int) ([]model.User, httperors.HttpErr)
	Update(id int, user *model.User) (*model.User, httperors.HttpErr)
	Delete(id int) (string, httperors.HttpErr)
	Cleaner(id string) (string, httperors.HttpErr)
	geneCode() (string, httperors.HttpErr)
	userExist(email string) bool
	userExistByid(id int) bool
}

func NewUserRepo() *userrepo {
	return &userrepo{}
}
func (userRepo userrepo) Create(user *model.User) (*model.User, httperors.HttpErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}
	ok, err1 := user.ValidatePassword()
	if !ok {
		return nil, err1
	}
	log.Println("step 1")
	ok = user.ValidateEmail()
	if !ok {
		return nil, httperors.NewNotFoundError("Your email format is wrong!")
	}
	ok = userRepo.userExist(user.Email)
	if ok {
		return nil, httperors.NewNotFoundError("Your email already exists!")
	}
	log.Println("step 2")

	hashpassword, err2 := user.HashPassword()
	if err2 != nil {
		return nil, err2
	}
	log.Println("step 3")
	user.Password = hashpassword
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	log.Println("step 4")
	defer IndexRepo.DbClose(GormDB)
	code, x := userRepo.geneCode()
	if x != nil {
		return nil, x
	}
	user.Usercode = code
	// fmt.Println(user)
	user.Admin = "true"
	// _, err := Businessrepo.Create(&model.Business{Name: user.Business})
	// if err != nil {
	// 	return "", err
	// }
	GormDB.Create(&user)
	return user, nil
}
func (userRepo userrepo) Login(auser *model.LoginUser) (*model.Auth, httperors.HttpErr) {
	if err := auser.Validate(); err != nil {
		return nil, err
	}
	ok := userRepo.userExist(auser.Email)
	if !ok {
		return nil, httperors.NewNotFoundError("Your email does not exists!")
	}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	defer IndexRepo.DbClose(GormDB)
	user := model.User{}

	GormDB.Model(&user).Where("email = ?", auser.Email).First(&user)
	ok = user.Compare(auser.Password, user.Password)
	if !ok {
		return nil, httperors.NewNotFoundError("wrong email password combo!")
	}
	tk := &model.Token{
		UserID:   user.ID,
		UName:    user.UName,
		Admin:    user.Admin,
		Usercode: user.Usercode,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: model.ExpiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatal("Error loading key")
	// }
	encyKey, err := LoadKey()
	if err != nil {
		log.Fatal("Error loading key")
	}
	tokenString, error := token.SignedString([]byte(encyKey.EncryptionKey))
	if error != nil {
		fmt.Println(error)
	}
	// messages ,e := userRepo.UnreadMessages(user.ID)
	// if e != nil {
	// 	return nil, e
	// }
	// norti ,e := userRepo.UnreadNortifications(user.ID)
	// if e != nil {
	// 	return nil, e
	// }
	auth := &model.Auth{UserID: user.ID, UName: user.UName, Usercode: user.Usercode, Admin: user.Admin, Picture: user.Picture, Token: tokenString}
	GormDB.Create(&auth)

	return auth, nil
}

func (userRepo userrepo) all() (t []model.User, r httperors.HttpErr) {

	user := model.User{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	defer IndexRepo.DbClose(GormDB)
	GormDB.Model(&user).Find(&t)
	return t, nil

}
func (userRepo userrepo) Logout(token string) httperors.HttpErr {
	auth := model.Auth{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return err1
	}
	defer IndexRepo.DbClose(GormDB)
	res := GormDB.First(&auth, "token =?", token)
	if res.Error != nil {
		return httperors.NewNotFoundError("Something went wrong logging out!")
	}

	GormDB.Model(&auth).Where("token =?", token).Delete(&auth)

	return nil
}
func (userRepo userrepo) geneCode() (string, httperors.HttpErr) {
	user := model.User{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	defer IndexRepo.DbClose(GormDB)
	err := GormDB.Last(&user)
	if err.Error != nil {
		var c1 uint = 1
		code := "UserCode" + strconv.FormatUint(uint64(c1), 10)
		return code, nil
	}
	c1 := user.ID + 1
	code := "UserCode" + strconv.FormatUint(uint64(c1), 10)
	return code, nil

}
func (userRepo userrepo) GetOne(id int) (*model.User, httperors.HttpErr) {
	ok := userRepo.userExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("User with that code does not exists!")
	}
	user := model.User{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	defer IndexRepo.DbClose(GormDB)

	GormDB.Model(&user).Where("id = ?", id).First(&user)
	return &user, nil
}
func (userRepo userrepo) userExistbycode(code string) bool {
	u := model.User{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	defer IndexRepo.DbClose(GormDB)
	GormDB.Where("usercode = ?", code).First(&u)
	if u.ID == 0 {
		return false
	}
	return true

}
func (userRepo userrepo) userbycode(code string) *model.User {
	u := model.User{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil
	}
	defer IndexRepo.DbClose(GormDB)
	GormDB.Where("usercode = ?", code).First(&u)
	if u.ID == 0 {
		return nil
	}
	return &u

}
func (userRepo userrepo) GetAll(search string, page, pagesize int) ([]model.User, httperors.HttpErr) {
	results := []model.User{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	defer IndexRepo.DbClose(GormDB)
	if search == "" {
		GormDB.Find(&results)
	}
	// db.Scopes(Paginate(r)).Find(&users)
	GormDB.Scopes(Paginate(page, pagesize)).Where("name LIKE ?", "%"+search+"%").Or("email LIKE ?", "%"+search+"%").Or("company LIKE ?", "%"+search+"%").Find(&results)

	return results, nil
}

// func (userRepo userrepo) UpdateRole(id int, role, usercode string) (string, httperors.HttpErr) {
// 	user := model.User{}
// 	ok := Userrepo.userExistByid(id)
// 	if !ok {
// 		return "", httperors.NewNotFoundError("customer with that id does not exists!")
// 	}

// 	GormDB, err1 := IndexRepo.Getconnected()
// 	if err1 != nil {
// 		return "", err1
// 	}

// 	// "employee",
// 	// "supervisor",
// 	// "admin",
// 	if usercode == "employee" {
// 		GormDB.Model(&user).Where("id = ?", id).Update("employee", true)
// 	}
// 	if usercode == "supervisor" {
// 		GormDB.Model(&user).Where("id = ?", id).Update("supervisor", true)
// 	}
// 	GormDB.Model(&user).Where("id = ?", id).Update("admin", true)
// 	IndexRepo.DbClose(GormDB)

// 	return "user updated succesifully", nil
// }
func (userRepo userrepo) Update(id int, user *model.User) (*model.User, httperors.HttpErr) {
	ok := userRepo.userExistByid(id)
	if !ok {
		return nil, httperors.NewNotFoundError("User with that id does not exists!")
	}

	hashpassword, err2 := user.HashPassword()
	if err2 != nil {
		return nil, err2
	}
	user.Password = hashpassword
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return nil, err1
	}
	defer IndexRepo.DbClose(GormDB)
	User := model.User{}
	uuser := model.User{}

	GormDB.Model(&User).Where("id = ?", id).First(&uuser)
	if user.FName == "" {
		user.FName = uuser.FName
	}
	if user.LName == "" {
		user.LName = uuser.LName
	}
	if user.UName == "" {
		user.UName = uuser.UName
	}
	if user.Phone == "" {
		user.Phone = uuser.Phone
	}
	if user.Address == "" {
		user.Address = uuser.Address
	}
	if user.Picture == "" {
		user.Picture = uuser.Picture
	}
	if user.Email == "" {
		user.Email = uuser.Email
	}
	// if user.Admin == false {
	// 	user.Admin = true
	// }
	GormDB.Save(&user)

	return user, nil
}
func (userRepo userrepo) Delete(id int) (string, httperors.HttpErr) {
	ok := userRepo.userExistByid(id)
	if !ok {
		return "", httperors.NewNotFoundError("User with that id does not exists!")
	}
	user := model.User{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	defer IndexRepo.DbClose(GormDB)
	GormDB.Model(&user).Where("id = ?", id).Delete(&user)
	return "deleted successfully", nil
}
func (userRepo userrepo) Cleaner(id string) (string, httperors.HttpErr) {
	user := model.User{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return "", err1
	}
	defer IndexRepo.DbClose(GormDB)
	GormDB.Unscoped().Where("u_name = ?", id).Delete(&user)

	return "deleted successfully", nil
}
func (userRepo userrepo) userExist(email string) bool {
	user := model.User{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	res := GormDB.First(&user, "email =?", email)
	IndexRepo.DbClose(GormDB)
	return res.Error == nil

}
func (userRepo userrepo) userExistByid(id int) bool {
	user := model.User{}
	GormDB, err1 := IndexRepo.Getconnected()
	if err1 != nil {
		return false
	}
	defer IndexRepo.DbClose(GormDB)
	res := GormDB.First(&user, "id =?", id)
	return res.Error == nil

}

// func (userRepo userrepo)UnreadMessages(id uint)  (int, *httperors.HttpError)  {
// 	messages := []model.Message{}
// 	GormDB, err1 := IndexRepo.Getconnected()
// 	if err1 != nil {
// 		return 0, err1
// 	}
// 	GormDB.Where("id = ? AND read = ? ", id, false).Find(&messages)
// 	 c := 0
// 	 for i, _:= range messages{
// 		 c += i
// 	 }
// 	IndexRepo.DbClose(GormDB)
// 	return c, nil

// }
// func (userRepo userrepo)UnreadNortifications(id uint)  (int, *httperors.HttpError)  {
// 	ns := []model.Nortification{}
// 	GormDB, err1 := IndexRepo.Getconnected()
// 	if err1 != nil {
// 		return 0, err1
// 	}
// 	GormDB.Where("id = ? AND read = ? ", id, false).Find(&ns)
// 	 c := 0
// 	 for i, _:= range ns{
// 		 c += i
// 	 }
// 	IndexRepo.DbClose(GormDB)
// 	return c, nil

// }
