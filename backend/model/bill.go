package model

import "time"

type Bill struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	RoomNo       string     `gorm:"column:room_no;type:varchar(10);not null;uniqueIndex:uk_room_year_month" json:"room_no"`
	Year         int        `gorm:"column:year;not null;uniqueIndex:uk_room_year_month" json:"year"`
	Month        int        `gorm:"column:month;not null;uniqueIndex:uk_room_year_month" json:"month"`
	ElecLast     *float64   `gorm:"column:elec_last;type:decimal(10,2)" json:"elec_last"`
	ElecCurrent  *float64   `gorm:"column:elec_current;type:decimal(10,2)" json:"elec_current"`
	ElecUsage    *float64   `gorm:"column:elec_usage;type:decimal(10,2)" json:"elec_usage"`
	ElecCost     *float64   `gorm:"column:elec_cost;type:decimal(10,2)" json:"elec_cost"`
	WaterLast    *float64   `gorm:"column:water_last;type:decimal(10,2)" json:"water_last"`
	WaterCurrent *float64   `gorm:"column:water_current;type:decimal(10,2)" json:"water_current"`
	WaterUsage   *float64   `gorm:"column:water_usage;type:decimal(10,2)" json:"water_usage"`
	WaterCost    *float64   `gorm:"column:water_cost;type:decimal(10,2)" json:"water_cost"`
	RentCost     float64    `gorm:"column:rent_cost;type:decimal(10,2)" json:"rent_cost"`
	TotalCost    float64    `gorm:"column:total_cost;type:decimal(10,2)" json:"total_cost"`
	IsPaid       int        `gorm:"column:is_paid;type:tinyint;default:0" json:"is_paid"`
	PaidAt       *time.Time `gorm:"column:paid_at" json:"paid_at"`
	CreatedAt    time.Time  `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time  `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (Bill) TableName() string {
	return "bills"
}
