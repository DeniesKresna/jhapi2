package types

import "gorm.io/gorm"

type User struct {
	Email     string `json:"email,omitempty" gorm:"email"`
	Password  string `json:"password,omitempty" gorm:"password"`
	Phone     int    `json:"phone,omitempty" gorm:"phone"`
	FirstName string `json:"first_name,omitempty" gorm:"first_name"`
	LastName  string `json:"last_name,omitempty" gorm:"last_name"`
	RoleID    int    `json:"role_id,omitempty" gorm:"role_id"`
	Status    int    `json:"status" gorm:"status"`
	gorm.Model
}

func (u *User) TableName() string {
	return "users"
}
