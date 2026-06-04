package repository

import (
	"fmt"

	"github.com/PZOBJECT/house/backend/model"
	"gorm.io/gorm"
)

type BillRepo struct {
	db *gorm.DB
}

func NewBillRepo(db *gorm.DB) *BillRepo {
	return &BillRepo{db: db}
}

func (r *BillRepo) GetLastBill(roomNo string, year, month int) (*model.Bill, error) {
	if month == 1 {
		year--
		month = 12
	} else {
		month--
	}
	var bill model.Bill
	err := r.db.Where("room_no = ? AND year = ? AND month = ?", roomNo, year, month).First(&bill).Error
	if err != nil {
		return nil, err
	}
	return &bill, nil
}

func (r *BillRepo) List(year, month int, roomNo string, isPaid *int, floor *int) ([]model.Bill, error) {
	var bills []model.Bill
	q := r.db.Model(&model.Bill{})
	if year > 0 {
		q = q.Where("year = ?", year)
	}
	if month > 0 {
		q = q.Where("month = ?", month)
	}
	if roomNo != "" {
		q = q.Where("room_no = ?", roomNo)
	}
	if isPaid != nil {
		q = q.Where("is_paid = ?", *isPaid)
	}
	if floor != nil {
		q = q.Where("room_no LIKE ?", fmt.Sprintf("%d%%", *floor))
	}
	if err := q.Order("room_no").Find(&bills).Error; err != nil {
		return nil, err
	}
	return bills, nil
}

func (r *BillRepo) Create(bill *model.Bill) error {
	return r.db.Create(bill).Error
}

func (r *BillRepo) Update(bill *model.Bill) error {
	return r.db.Save(bill).Error
}

func (r *BillRepo) GetByID(id uint) (*model.Bill, error) {
	var bill model.Bill
	if err := r.db.First(&bill, id).Error; err != nil {
		return nil, err
	}
	return &bill, nil
}

func (r *BillRepo) GetByRoomYearMonth(roomNo string, year, month int) (*model.Bill, error) {
	var bill model.Bill
	err := r.db.Where("room_no = ? AND year = ? AND month = ?", roomNo, year, month).First(&bill).Error
	if err != nil {
		return nil, err
	}
	return &bill, nil
}
