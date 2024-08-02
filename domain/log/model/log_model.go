package logModel

type Id struct {
	ID int `gorm:"primaryKey;autoIncrement; column:id; type:bigint" json:"id" example:"1" validate:"required"`
}

type TimeStamp struct {
	CreatedAt string `gorm:"column:created_at; type:varchar(100)" json:"created_at" autoCreateTime:"true" example:"2024-01-01T00:00:00Z"`
	UpdatedAt string `gorm:"column:updated_at; type:varchar(100)" json:"updated_at" autoCreateTime:"true" example:"2024-01-01T00:00:00Z"`
}

type LogEntity struct {
	Level   string `gorm:"column:level; type:varchar(255)" json:"level"`
	Message string `gorm:"column:message; type:text" json:"message"`
}

type Log struct {
	Id
	LogEntity
	TimeStamp
}
