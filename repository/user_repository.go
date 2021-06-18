package repository

import (
	"casbin-golang/model"
	"log"

	"gorm.io/gorm"
)

type userRepository struct {
	DB *gorm.DB
}

// UserRepository : represent the user's repository contract
type UserRepository interface {
	AddUser(model.User) (model.User, error)
	GetUser(int) (model.User, error)
	GetByEmail(string) (model.User, error)
	GetAllUser() ([]model.User, error)
	UpdateUser(model.User) (model.User, error)
	DeleteUser(model.User) (model.User, error)
	Migrate() error
}

// NewUserRepository -> returns new user repository
func NewUserRepository(db *gorm.DB) UserRepository {
	return userRepository{
		DB: db,
	}
}

func (u userRepository) Migrate() error {
	log.Print("[UserRepository]...Migrate")
	return u.DB.AutoMigrate(&model.User{})
}

func (u userRepository) GetUser(id int) (user model.User, err error) {
	return user, u.DB.First(&user, id).Error
}

func (u userRepository) GetByEmail(email string) (user model.User, err error) {
	return user, u.DB.First(&user, "email=?", email).Error
}

func (u userRepository) GetAllUser() (users []model.User, err error) {
	return users, u.DB.Find(&users).Error
}

func (u userRepository) AddUser(user model.User) (model.User, error) {
	return user, u.DB.Create(&user).Error
}

func (u userRepository) UpdateUser(user model.User) (model.User, error) {
	if err := u.DB.First(&user, user.ID).Error; err != nil {
		return user, err
	}
	return user, u.DB.Model(&user).Updates(&user).Error
}

func (u userRepository) DeleteUser(user model.User) (model.User, error) {
	if err := u.DB.First(&user, user.ID).Error; err != nil {
		return user, err
	}
	return user, u.DB.Delete(&user).Error
}
