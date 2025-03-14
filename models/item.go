package models

import "time"

type Item struct {
	ID          uint      `json:"id" gorm:"column:id;primaryKey;autoIncrement"`
	Name        string    `json:"name" gorm:"column:name"`
	Price       float64   `json:"price" gorm:"column:price"`
	Description string    `json:"description,omitempty" gorm:"column:description"`
	CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updated_at"`
}

func (Item) TableName() string {
	return "items"
}
