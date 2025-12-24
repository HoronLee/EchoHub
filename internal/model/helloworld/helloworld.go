package helloworld

import "time"

// HelloWorld 定义HelloWorld实体
type HelloWorld struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
