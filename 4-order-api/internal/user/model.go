package user

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Phone     string `gorm:"index"`
	Name      string
	SessionId string
	Code      int
}
