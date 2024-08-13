package seeder

import (
	"go-clean/domain/user"
	"go-clean/utils/encrypt"

	"gorm.io/gorm"
)

func DatabaseSeeder(db *gorm.DB) {
	hashPassword, _ := encrypt.Brypt("123")
	user := []*user.User{
		{Name: "super admin", Email: "super@mail.com", Password: hashPassword},
		{Name: "admin", Email: "admin@mail.com", Password: hashPassword},
		{Name: "user", Email: "user@mail.com", Password: hashPassword},
	}
	db.Create(&user)
}
