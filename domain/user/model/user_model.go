package userModel

type Id struct {
	ID int `gorm:"primaryKey;autoIncrement; column:id; type:bigint" json:"id" example:"1" validate:"required"`
}

type TimeStamp struct {
	CreatedAt string `gorm:"column:created_at; type:varchar(100)" json:"created_at" autoCreateTime:"true" example:"2024-01-01T00:00:00Z"`
	UpdatedAt string `gorm:"column:updated_at; type:varchar(100)" json:"updated_at" autoCreateTime:"true" example:"2024-01-01T00:00:00Z"`
	DeletedAt string `gorm:"column:deleted_at; type:varchar(100)" json:"deleted_at,omitempty" autoCreateTime:"true" example:"2024-01-01T00:00:00Z"`
}

type UserEntity struct {
	Name     string `gorm:"column:name; type:varchar(255)" json:"name" example:"John Doe" validate:"required"`
	Email    string `gorm:"column:email; type:varchar(255)" json:"email" example:"lQwLd@example.com" validate:"required,email"`
	Password string `gorm:"column:password; type:varchar(255)" json:"password" example:"password"`
}

type User struct {
	Id
	UserEntity
	TimeStamp
}

type UserCreate struct {
	UserEntity
}

type UserUpdate struct {
	Id
	UserEntity
}

type UserSearchRequest struct {
	Id    int    `json:"id" validate:"omitempty,min=1"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Page  int    `json:"page" validate:"min=1"`
	Size  int    `json:"size" validate:"min=1"`
}

type UserResponse struct {
	Id
	Name  string `json:"name"`
	Email string `json:"email"`
	TimeStamp
}

func UserToResponse(user *UserResponse) *UserResponse {
	return &UserResponse{
		Id:        user.Id,
		Name:      user.Name,
		Email:     user.Email,
		TimeStamp: user.TimeStamp,
	}
}
