package model

import "gorm.io/gorm"

//User -> model for users table
type User struct {
	gorm.Model
	Name     string `json:"name" `
	Email    string `json:"email"  gorm:"unique"`
	Role     string `json:"role" gorm:"-"`
	Password string `json:"password" `
}

//TableName --> Table for Product Model
func (User) TableName() string {
	return "users"
}
