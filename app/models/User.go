package models

import (
	"errors"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"

)

type User struct {
	gorm.Model
	Name string `gorm:"size:100;not null" json:"name"`
	Email string `gorm:"size:100;not null" json:"email"`
	Password string `gorm:"size:100;not null" json:"password"`
	ProfileImage string `gorm:"size:255" json:"profileImage"`
}

// HashPassword hashes password from the user input
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash checks hashed password and password from user input if they match
func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return errors.New("Password incorrect")
	}
	return nil
}

// BeforeSave hash the user password
func (user *User) BeforeSave() error {
	password := strings.TrimSpace(user.Password)
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

// Prepare strips user input of any white spaces
func (user *User) Prepare() {
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.TrimSpace(user.Email)
	user.ProfileImage = strings.TrimSpace(user.ProfileImage)
}

// Validate the user input
func (user *User) Validate(action string) error {
	switch strings.ToLower(action) {
	case "login":
		if user.Email == "" {
			return errors.New("Email is required")
		}
		if user.Password == "" {
			return errors.New("Password is required")
		}
		return nil
	default:
		if user.Name == "" {
			return errors.New("Name is required")
		}
		if user.Password == "" {
			return errors.New("Password is required")
		}
		if err := checkmail.ValidateFormat(user.Email); err != nil {
			return errors.New("Invalid email")
		}
		return nil
	}
}

// SaveUser add a user to the database
func (user *User) SaveUser(db *gorm.DB) (*User, error) {
	var err error

	// Debug a single operation, show detailed log for this operation
	err = db.Debug().Create(&user).Error
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

// GetUser return a user based on email
func (user *User) GetUser(db *gorm.DB) (*User, error) {
	account := &User{}
	if err := db.Debug().Table("users").Where("email = ?", user.Email).First(account).Error; err != nil {
		return nil, err
	}
	return account, nil
}

// GetUsers return all users
func (user *User) GetUsers(db *gorm.DB) (*[]User, error) {
	users := []User{}

	if err := db.Debug().Table("users").Find(&users).Error; err != nil {
		return &[]User{}, err
	}

	return &users, nil
}