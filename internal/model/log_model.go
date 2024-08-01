package model

type Log struct {
	ID        int    `gorm:"primaryKey;autoIncrement; column:id; type:bigint" json:"id"`
	Level     string `gorm:"column:level; type:varchar(255)" json:"level"`
	Message   string `gorm:"column:message; type:text" json:"message"`
	CreatedAt string `gorm:"column:created_at; type:varchar(100)" json:"created_at" autoCreateTime:"true"`
	UpdatedAt string `gorm:"column:updated_at; type:varchar(100)" json:"updated_at" autoCreateTime:"true"`
}
