package repository

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go-jwt-mysql/config"
	"go-jwt-mysql/model"
	"golang.org/x/crypto/bcrypt"
)

//GetAllUsers Fetch all user data
func GetAllUsers(user *[]model.User) (err error) {
	if err = config.Db.Find(user).Error; err != nil {
		return err
	}
	return nil
}

//CreateUser ... Insert New data
func CreateUser(user *model.User) (err error) {
	err = BeforeSave(user)
	if err != nil {
		return err
	}
	if err = config.Db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

//GetUserByID ... Fetch only one user by Id
func GetUserByID(user *model.User, id string) (err error) {
	if err = config.Db.Where("id = ?", id).First(user).Error; err != nil {
		return err
	}
	return nil
}

func GetUserByUsername(user *model.User, username string) (err error) {
	if err = config.Db.Where("username = ?", username).First(user).Error; err != nil {
		return err
	}
	return nil
}

//UpdateUser ... Update user
func UpdateUser(user *model.User, id string) (err error) {
	err = BeforeSave(user)
	if err != nil {
		return err
	}
	fmt.Println(user)
	config.Db.Model(user).Updates(user)
	return nil
}

//DeleteUser ... Delete user
func DeleteUser(user *model.User, id string) (err error) {
	config.Db.Where("id = ?", id).Delete(user)
	return nil
}

func BeforeSave(user *model.User) (err error) {
	hashedPassword, err := Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return nil
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
