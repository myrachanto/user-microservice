package model

import ( 
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateUserInputRequiredFields(t *testing.T) {
	jsondata := `{"FName":"jane","LName":"Doe","UName":"doe","Usercode": "Doe345","Phone":"1234567","Email":   "email@example.com","Password": "1234567","Address":"psd 456 king view"
	}`
	user := &User{}
	if err := json.Unmarshal([]byte(jsondata), &user); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	// fmt.Println("------------------", user)
	expected := ""
	if err := user.Validate(); err != nil {
		fmt.Println("------------------", err.Message())
		expected = "Invalid first Name"
		if err.Message() == expected {
			assert.EqualValues(t, "", err.Message(), "Error validating last name")
		}
		expected = "Invalid last name"
		if err.Message() == expected {
			assert.EqualValues(t, "", err.Message(), "Error validating last name")
		}
		expected = "Invalid username"
		if err.Message() == expected {
			assert.EqualValues(t, "", err.Message(), "Error validating username")
		}
		expected = "Invalid phone number"
		if err.Message() == expected {
			assert.EqualValues(t, "", err.Message(), "Error validating Phone number")
		}
		expected = "Invalid Email"
		if err.Message() == expected {
			assert.EqualValues(t, "", err.Message(), "Error validating email")
		}
		expected = "Invalid password"
		if err.Message() == expected {
			assert.EqualValues(t, "", err.Message(), "Error validating password")
		}
		expected = "Invalid phone number"
		if err.Message() == expected {
			assert.EqualValues(t, "", err.Message(), "Error validating Phone number")
		}
		expected = "Invalid Address"
		if err.Message() == expected {
			assert.EqualValues(t, "", err.Message(), "Error validating  address")
		}

	}

}

func TestValidateLoginUserInputRequiredFields(t *testing.T) {
	jsondata := `{"Email":"email@example.com","Password":"1234567"}`
	user := &LoginUser{}
	if err := json.Unmarshal([]byte(jsondata), &user); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	// fmt.Println("------------------", user)
	expected := ""
	if err := user.Validate(); err != nil {
		expected = "Invalid Email"
		if err.Message() == expected {
			assert.EqualValues(t, "", err.Message(), "Error validating email")
		}
		expected = "Invalid password"
		if err.Message() == expected {
			assert.EqualValues(t, "", err.Message(), "Error validating password")
		}

	}

}
func TestComparingPasswords(t *testing.T) {
	// fmt.Println("------------------", user)
	password := "anton345"
	user := User{}
	pas1, _ := user.HashPassword()
	ok := user.Compare(password, pas1)
	if !ok {
		assert.EqualValues(t, false, ok, "Error comparing passwords")
	}
}

func TestValidateEmailInputRequiredFields(t *testing.T) {
	jsondata := `{"FName":"jane","LName":"Doe","UName":"doe","Usercode": "Doe345","Phone":"1234567","Email":   "email@example.com","Password": "1234567","Address":"psd 456 king view"
	}`
	user := &User{}
	if err := json.Unmarshal([]byte(jsondata), &user); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	// fmt.Println("------------------", user)
	ok := user.ValidateEmail()
	if !ok {
		assert.EqualValues(t, false, ok, "Error Validating emails")
	}

}

func TestValidatePasswordInputRequiredFields(t *testing.T) {
	jsondata := `{"FName":"jane","LName":"Doe","UName":"doe","Usercode": "Doe345","Phone":"1234567","Email":   "email@example.com","Password": "1234567","Address":"psd 456 king view"
	}`
	user := &User{}
	if err := json.Unmarshal([]byte(jsondata), &user); err != nil {
		t.Errorf("failed to unmarshal user data %v", err.Error())
	}
	// fmt.Println("------------------", user)
	ok, _ := user.ValidatePassword()
	if !ok {
		assert.EqualValues(t, true, ok, "Error Validating passwords")
	}

}
