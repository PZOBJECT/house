package service

import (
	"errors"

	"github.com/PZOBJECT/house/backend/model"
	"github.com/PZOBJECT/house/backend/repository"
)

type RoomService struct {
	repo *repository.RoomRepo
}

func NewRoomService(repo *repository.RoomRepo) *RoomService {
	return &RoomService{repo: repo}
}

type UpdateRoomReq struct {
	RentPrice  *float64 `json:"rent_price"`
	ElecPrice  *float64 `json:"elec_price"`
	WaterPrice *float64 `json:"water_price"`
	IsRented   *int     `json:"is_rented"`
}

func (s *RoomService) GetRoomList(rented *int, floor *int) ([]model.Room, error) {
	return s.repo.GetAll(rented, floor)
}

func (s *RoomService) UpdateRoom(roomNo string, req UpdateRoomReq) error {
	room, err := s.repo.GetByRoomNo(roomNo)
	if err != nil {
		return errors.New("房间不存在")
	}
	if req.RentPrice != nil {
		if *req.RentPrice < 0 {
			return errors.New("租金不能为负数")
		}
		room.RentPrice = *req.RentPrice
	}
	if req.ElecPrice != nil {
		if *req.ElecPrice < 0 {
			return errors.New("电费单价不能为负数")
		}
		room.ElecPrice = *req.ElecPrice
	}
	if req.WaterPrice != nil {
		if *req.WaterPrice < 0 {
			return errors.New("水费单价不能为负数")
		}
		room.WaterPrice = *req.WaterPrice
	}
	if req.IsRented != nil {
		if *req.IsRented != 0 && *req.IsRented != 1 {
			return errors.New("出租状态值无效")
		}
		room.IsRented = *req.IsRented
	}
	return s.repo.Update(room)
}
