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
	FName    string `json:"fname"`
	LName    string `json:"lname"`
	UName    string `json:"uname"`
	Usercode string `json:"usercode"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
	// Dob           *time.Time       `json:"dob"`
	Picture     string `json:"picture"`
	Business    string `json:"business"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Admin       string `json:"admin"`
	Employee    string `json:"employee"`
	Supervisor  string `json:"supervisor"`
	Accesslevel string `json:"level"`
	Shopalias   string `json:"shopalias"`
	// Message       []*Message       `gorm:"many2many:user_messages;" json:"message"`
	// Nortification []*Nortification `gorm:"many2many:user_nortifications;" json:"nortification"`
	gorm.Model
}

//Auth ..
type Auth struct {
	//User User `gorm:"foreignKey:UserID; not null"`
	UserID     uint   `json:"userid"`
	UName      string `json:"uname"`
	Usercode   string `json:"usercode"`
	Picture    string `json:"picture"`
	Token      string `gorm:"size:500;not null" json:"token"`
	Admin      string `json:"admin"`
	Business   string `json:"business"`
	Employee   string `json:"employee"`
	Supervisor string `json:"supervisor"`
	Level      string `json:"level"`
	gorm.Model
}

//LoginUser ..
type LoginUser struct {
	Email    string `gorm:"not null"`
	Password string `gorm:"not null"`
}

//LoginUser ..
type LoginEmployee struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Alias    string `json:"alias"`
}

//Token struct declaration
type Token struct {
	UserID     uint   `json:"user_id"`
	UName      string `json:"uname"`
	Email      string `json:"email"`
	Usercode   string `json:"usercode"`
	Admin      string `json:"admin"`
	Employee   string `json:"employee"`
	Supervisor string `json:"supervisor"`
	Role       string `json:"role"`
	SystemAuth string `json:"system_auth"`
	Bizname    string `json:"bizname"`
	Bizstatus  string `json:"bizstatus"`
	Shopalias  string `json:"shopalias"`
	Central    string `json:"central"`
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
func (user User) ValidateEmail(email string) (matchedString bool) {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&amp;'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	matchedString = re.MatchString(email)
	return matchedString
}

//ValidatePassword ...
func (user User) ValidatePassword(password string) (bool, httperors.HttpErr) {
	if len(password) < 5 {
		return false, httperors.NewBadRequestError("your password need more characters!")
	} else if len(password) > 32 {
		return false, httperors.NewBadRequestError("your password is way too long!")
	}
	return true, nil
}

//HashPassword ..
func (user User) HashPassword(password string) (string, httperors.HttpErr) {
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return "", httperors.NewNotFoundError("type a stronger password!")
	}
	return string(pass), nil

}

//Compare ..
func (user User) Compare(p1, p2 string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p2), []byte(p1))
	if err != nil {
		return false
	}
	return true
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
func (e LoginEmployee) Validate() httperors.HttpErr {
	if e.Email == "" {
		return httperors.NewNotFoundError("Invalid Email")
	}
	if e.Password == "" {
		return httperors.NewNotFoundError("Invalid password")
	}
	if e.Alias == "" {
		return httperors.NewNotFoundError("Invalid password")
	}
	return nil
}

//Validate ..
func (user User) Validate() httperors.HttpErr {
	if user.FName == "" {
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
	// if user.Picture == "" {
	// 	return httperors.NewNotFoundError("Invalid picture")
	// }
	if user.Email == "" {
		return httperors.NewNotFoundError("Invalid picture")
	}
	return nil
}
