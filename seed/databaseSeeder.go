package seeder

import (
	"go-clean/domain/user"
	"go-clean/pkg"

	"gorm.io/gorm"
)

func DatabaseSeeder(db *gorm.DB) {
	hashPassword, _ := pkg.NewBcryptService().Bcrypt("password")
	user := []*user.User{
		{Name: "super admin", Email: "super@mail.com", Password: hashPassword},
		{Name: "admin", Email: "admin@mail.com", Password: hashPassword},
		{Name: "user", Email: "user@mail.com", Password: hashPassword},
	}
	db.Create(&user)
}
