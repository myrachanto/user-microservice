package model

import (
	"regexp"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	httperors "github.com/myrachanto/custom-http-error"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// var (
// 	Userkeyconfig User = User{}
// )
//ExpiresAt ..
var ExpiresAt = time.Now().Add(time.Minute * 100000).Unix()

//User ..
type User struct {
	FName    string `json:"f_name"`
	LName    string `json:"l_name"`
	UName    string `json:"u_name"`
	Usercode string `json:"usercode"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	// Dob           *time.Time       `json:"dob"`
	Picture  string `json:"picture"`
	Business string `json:"business"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Admin    string `json:"admin"`
	gorm.Model
}

//Auth ..
type Auth struct {
	//User User `gorm:"foreignKey:UserID; not null"`
	UserID   uint   `json:"userid"`
	UName    string `json:"uname"`
	Usercode string `json:"usercode"`
	Picture  string `json:"picture"`
	Token    string `gorm:"size:500;not null" json:"token"`
	Admin    string `json:"admin"`
	gorm.Model
}

//LoginUser ..
type LoginUser struct {
	Email    string `gorm:"not null"`
	Password string `gorm:"not null"`
}

//Token struct declaration
type Token struct {
	UserID   uint   `json:"user_id"`
	UName    string `json:"uname"`
	Email    string `json:"email"`
	Usercode string `json:"usercode"`
	Admin    string `json:"admin"`
	*jwt.StandardClaims
}

//UserProfile user profile and messages
// type UserProfile struct {
// 	User          User            `json:"user"`
// 	Inbox         []Message       `json:"inbox"`
// 	Sent          []Message       `json:"sent"`
// 	Users         []User          `json:"users"`
// 	Nortification []Nortification `json:"nortifications"`
// }

// type Userkey struct{
// 	EncryptionKey string `mapstructure:"EncryptionKey"`
// }
//ValidateEmail ..
func (user User) ValidateEmail() (matchedString bool) {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&amp;'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	matchedString = re.MatchString(user.Email)
	return matchedString
}

//ValidatePassword ...
func (user User) ValidatePassword() (bool, httperors.HttpErr) {
	if len(user.Password) < 5 {
		return false, httperors.NewBadRequestError("your password need more characters!")
	} else if len(user.Password) > 32 {
		return false, httperors.NewBadRequestError("your password is way too long!")
	}
	return true, nil
}

//HashPassword ..
func (user User) HashPassword() (string, httperors.HttpErr) {
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return "", httperors.NewNotFoundError("type a stronger password!")
	}
	return string(pass), nil

}

//Compare ..
func (user User) Compare(p1, p2 string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p2), []byte(p1))
	return err == nil
}

//Validate ..
func (loginuser LoginUser) Validate() httperors.HttpErr {
	if loginuser.Email == "" {
		return httperors.NewNotFoundError("Invalid Email")
	}
	if loginuser.Password == "" {
		return httperors.NewNotFoundError("Invalid password")
	}
	return nil
}

//Validate ..
func (user *User) Validate() httperors.HttpErr {
	if user.FName == "" {
		// log.Println("Invalid username", user.FName)
		return httperors.NewNotFoundError("Invalid first Name")
	}
	if user.LName == "" {
		return httperors.NewNotFoundError("Invalid last name")
	}
	if user.UName == "" {
		return httperors.NewNotFoundError("Invalid username")
	}
	if user.Phone == "" {
		return httperors.NewNotFoundError("Invalid phone number")
	}
	if user.Email == "" {
		return httperors.NewNotFoundError("Invalid Email")
	}
	if user.Address == "" {
		return httperors.NewNotFoundError("Invalid Address")
	}
	if user.Password == "" {
		return httperors.NewNotFoundError("Invalid password")
	}
	return nil
}
