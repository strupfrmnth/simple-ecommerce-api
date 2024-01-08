package repository

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// model
type User struct {
	gorm.Model
	Name string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserRepository interface {
	GetAllUser() ([]User, error)
	GetById(id int) (User, error)
	GetByUsername(username string) (User, error)
	AddUser(user *User) error
	UpdateUser(updateuser User) (User, error)
	Delete(id int) error
}

type concreteUserRepository struct {
	db *gorm.DB
}

// interface can store the pointer of struct or struct
// so we cannot return the pointer of interface
func NewUserRepository() UserRepository {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to Database")
		return nil
	}

	// update user model
	db.AutoMigrate(&User{}, &Order{}, &OrderedItem{})

	return &concreteUserRepository{
		db: db,
	}
}

func (repo *concreteUserRepository) GetByUsername(username string) (User, error) {
	var user User
	if err := repo.db.First(&user, "name = ?", username).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (repo *concreteUserRepository) AddUser(user *User) error {
	return repo.db.Create(&user).Error
}

func (repo *concreteUserRepository) GetAllUser() ([]User, error) {
	var users []User
	if err := repo.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *concreteUserRepository) GetById(id int) (User, error) {
	var user User
	if err := repo.db.First(&user, id).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (repo *concreteUserRepository) UpdateUser(updateuser User) (User, error) {
	var user User
	if err := repo.db.First(&user, updateuser.ID).Error; err != nil {
		return user, err
	}

	if err := repo.db.Model(&user).Updates(updateuser).Error; err != nil {
		return user, err
	}
	return user, nil
}

func (repo *concreteUserRepository) Delete(id int) error {
	return repo.db.Delete(&User{}, id).Error
}
