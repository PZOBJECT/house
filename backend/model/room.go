package model

import "time"

type Room struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	RoomNo      string    `gorm:"column:room_no;type:varchar(10);uniqueIndex;not null" json:"room_no"`
	RentPrice   float64   `gorm:"column:rent_price;type:decimal(10,2);not null" json:"rent_price"`
	ElecPrice   float64   `gorm:"column:elec_price;type:decimal(10,2);not null" json:"elec_price"`
	WaterPrice  float64   `gorm:"column:water_price;type:decimal(10,2);not null" json:"water_price"`
	IsRented    int       `gorm:"column:is_rented;type:tinyint;not null;default:1" json:"is_rented"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Room) TableName() string {
	return "rooms"
}
