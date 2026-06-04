package service

import (
	"errors"
	"time"

	"github.com/PZOBJECT/house/backend/model"
	"github.com/PZOBJECT/house/backend/repository"
	"gorm.io/gorm"
)

type BillService struct {
	billRepo *repository.BillRepo
	roomRepo *repository.RoomRepo
	db       *gorm.DB
}

func NewBillService(billRepo *repository.BillRepo, roomRepo *repository.RoomRepo, db *gorm.DB) *BillService {
	return &BillService{billRepo: billRepo, roomRepo: roomRepo, db: db}
}

type GenerateBillsReq struct {
	Year  int `json:"year"`
	Month int `json:"month"`
}

type UpdateBillReq struct {
	ElecCurrent  *float64 `json:"elec_current"`
	WaterCurrent *float64 `json:"water_current"`
}

func (s *BillService) GenerateBills(year, month int) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		rooms, err := s.roomRepo.GetRentedRooms()
		if err != nil {
			return err
		}
		for _, room := range rooms {
			// 检查是否已存在
			existing, err := s.billRepo.GetByRoomYearMonth(room.RoomNo, year, month)
			if err == nil && existing != nil {
				continue // 幂等：已存在则跳过
			}

			bill := &model.Bill{
				RoomNo:    room.RoomNo,
				Year:      year,
				Month:     month,
				RentCost:  room.RentPrice,
				TotalCost: room.RentPrice,
			}

			// 获取上月账单读数作为本月last
			lastBill, err := s.billRepo.GetLastBill(room.RoomNo, year, month)
			if err == nil && lastBill != nil {
				if lastBill.ElecCurrent != nil {
					v := *lastBill.ElecCurrent
					bill.ElecLast = &v
				}
				if lastBill.WaterCurrent != nil {
					v := *lastBill.WaterCurrent
					bill.WaterLast = &v
				}
			}

			if err := tx.Create(bill).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *BillService) GetBillList(year, month int, roomNo string, isPaid *int, floor *int) ([]model.Bill, error) {
	return s.billRepo.List(year, month, roomNo, isPaid, floor)
}

func (s *BillService) UpdateBill(id uint, req UpdateBillReq) error {
	bill, err := s.billRepo.GetByID(id)
	if err != nil {
		return errors.New("账单不存在")
	}
	if bill.IsPaid == 1 {
		return errors.New("已收费账单不可修改")
	}

	room := s.roomRepo.GetByRoomNoOrNil(bill.RoomNo)
	if room == nil {
		return errors.New("房间不存在")
	}

	if req.ElecCurrent != nil {
		if *req.ElecCurrent < 0 {
			return errors.New("电表读数不能为负数")
		}
		bill.ElecCurrent = req.ElecCurrent
		if bill.ElecLast != nil {
			usage := *req.ElecCurrent - *bill.ElecLast
			if usage < 0 {
				usage = 0
			}
			bill.ElecUsage = &usage
			cost := usage * room.ElecPrice
			bill.ElecCost = &cost
		}
	}

	if req.WaterCurrent != nil {
		if *req.WaterCurrent < 0 {
			return errors.New("水表读数不能为负数")
		}
		bill.WaterCurrent = req.WaterCurrent
		if bill.WaterLast != nil {
			usage := *req.WaterCurrent - *bill.WaterLast
			if usage < 0 {
				usage = 0
			}
			bill.WaterUsage = &usage
			cost := usage * room.WaterPrice
			bill.WaterCost = &cost
		}
	}

	// 重新计算总费用
	var total float64
	total += bill.RentCost
	if bill.ElecCost != nil {
		total += *bill.ElecCost
	}
	if bill.WaterCost != nil {
		total += *bill.WaterCost
	}
	bill.TotalCost = total

	return s.billRepo.Update(bill)
}

func (s *BillService) PayBill(id uint) error {
	bill, err := s.billRepo.GetByID(id)
	if err != nil {
		return errors.New("账单不存在")
	}
	if bill.IsPaid == 1 {
		return errors.New("账单已收费")
	}
	now := time.Now()
	bill.IsPaid = 1
	bill.PaidAt = &now
	return s.billRepo.Update(bill)
}

func (s *BillService) UnpayBill(id uint) error {
	bill, err := s.billRepo.GetByID(id)
	if err != nil {
		return errors.New("账单不存在")
	}
	if bill.IsPaid == 0 {
		return errors.New("账单未收费")
	}
	bill.IsPaid = 0
	bill.PaidAt = nil
	return s.billRepo.Update(bill)
}

func (s *BillService) GetReceiptData(id uint) (map[string]interface{}, error) {
	bill, err := s.billRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("账单不存在")
	}
	room := s.roomRepo.GetByRoomNoOrNil(bill.RoomNo)

	data := map[string]interface{}{
		"id":             bill.ID,
		"room_no":        bill.RoomNo,
		"year":           bill.Year,
		"month":          bill.Month,
		"elec_last":      bill.ElecLast,
		"elec_current":   bill.ElecCurrent,
		"elec_usage":     bill.ElecUsage,
		"elec_cost":      bill.ElecCost,
		"water_last":     bill.WaterLast,
		"water_current":  bill.WaterCurrent,
		"water_usage":    bill.WaterUsage,
		"water_cost":     bill.WaterCost,
		"rent_cost":      bill.RentCost,
		"total_cost":     bill.TotalCost,
		"is_paid":        bill.IsPaid,
		"paid_at":        bill.PaidAt,
		"elec_price":     nil,
		"water_price":    nil,
	}
	if room != nil {
		data["elec_price"] = room.ElecPrice
		data["water_price"] = room.WaterPrice
	}
	return data, nil
}
