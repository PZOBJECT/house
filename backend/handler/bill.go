package handler

import (
	"strconv"

	"github.com/PZOBJECT/house/backend/middleware"
	"github.com/PZOBJECT/house/backend/service"
	"github.com/gin-gonic/gin"
)

type BillHandler struct {
	svc *service.BillService
}

func NewBillHandler(svc *service.BillService) *BillHandler {
	return &BillHandler{svc: svc}
}

func (h *BillHandler) GenerateBills(c *gin.Context) {
	var req struct {
		Year  int `json:"year"`
		Month int `json:"month"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Error(c, "请求参数格式错误")
		return
	}
	if req.Year < 2000 || req.Year > 2100 {
		middleware.Error(c, "年份无效")
		return
	}
	if req.Month < 1 || req.Month > 12 {
		middleware.Error(c, "月份无效")
		return
	}
	if err := h.svc.GenerateBills(req.Year, req.Month); err != nil {
		middleware.Error(c, err.Error())
		return
	}
	middleware.Success(c, nil)
}

func (h *BillHandler) GetBillList(c *gin.Context) {
	year, _ := strconv.Atoi(c.DefaultQuery("year", "0"))
	month, _ := strconv.Atoi(c.DefaultQuery("month", "0"))
	roomNo := c.Query("room_no")

	var isPaid *int
	if v := c.Query("is_paid"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || (n != 0 && n != 1) {
			middleware.Error(c, "is_paid参数无效，只能为0或1")
			return
		}
		isPaid = &n
	}

	bills, err := h.svc.GetBillList(year, month, roomNo, isPaid)
	if err != nil {
		middleware.Error(c, err.Error())
		return
	}
	middleware.Success(c, bills)
}

func (h *BillHandler) UpdateBill(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		middleware.Error(c, "账单ID无效")
		return
	}
	var req service.UpdateBillReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Error(c, "请求参数格式错误")
		return
	}
	if err := h.svc.UpdateBill(uint(id), req); err != nil {
		middleware.Error(c, err.Error())
		return
	}
	middleware.Success(c, nil)
}

func (h *BillHandler) PayBill(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		middleware.Error(c, "账单ID无效")
		return
	}
	if err := h.svc.PayBill(uint(id)); err != nil {
		middleware.Error(c, err.Error())
		return
	}
	middleware.Success(c, nil)
}

func (h *BillHandler) UnpayBill(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		middleware.Error(c, "账单ID无效")
		return
	}
	if err := h.svc.UnpayBill(uint(id)); err != nil {
		middleware.Error(c, err.Error())
		return
	}
	middleware.Success(c, nil)
}

func (h *BillHandler) GetReceipt(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		middleware.Error(c, "账单ID无效")
		return
	}
	data, err := h.svc.GetReceiptData(uint(id))
	if err != nil {
		middleware.Error(c, err.Error())
		return
	}
	middleware.Success(c, data)
}
