package repository

import (
	"fmt"

	"github.com/PZOBJECT/house/backend/model"
	"gorm.io/gorm"
)

type RoomRepo struct {
	db *gorm.DB
}

func NewRoomRepo(db *gorm.DB) *RoomRepo {
	return &RoomRepo{db: db}
}

func (r *RoomRepo) GetAll(rented *int, floor *int) ([]model.Room, error) {
	var rooms []model.Room
	q := r.db.Model(&model.Room{})
	if rented != nil {
		q = q.Where("is_rented = ?", *rented)
	}
	if floor != nil {
		q = q.Where("room_no LIKE ?", fmt.Sprintf("%d%%", *floor))
	}
	if err := q.Order("room_no").Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *RoomRepo) GetByRoomNo(roomNo string) (*model.Room, error) {
	var room model.Room
	if err := r.db.Where("room_no = ?", roomNo).First(&room).Error; err != nil {
		return nil, err
	}
	return &room, nil
}

func (r *RoomRepo) Update(room *model.Room) error {
	return r.db.Save(room).Error
}

func (r *RoomRepo) BatchCreate(rooms []model.Room) error {
	return r.db.Create(&rooms).Error
}

func (r *RoomRepo) Count() (int64, error) {
	var count int64
	err := r.db.Model(&model.Room{}).Count(&count).Error
	return count, err
}

func (r *RoomRepo) GetByRoomNoOrNil(roomNo string) *model.Room {
	room, err := r.GetByRoomNo(roomNo)
	if err != nil {
		return nil
	}
	return room
}

func (r *RoomRepo) GetRentedRooms() ([]model.Room, error) {
	var rooms []model.Room
	err := r.db.Where("is_rented = 1").Order("room_no").Find(&rooms).Error
	return rooms, err
}
