package handler

import (
	"strconv"

	"github.com/PZOBJECT/house/backend/middleware"
	"github.com/PZOBJECT/house/backend/service"
	"github.com/gin-gonic/gin"
)

type RoomHandler struct {
	svc *service.RoomService
}

func NewRoomHandler(svc *service.RoomService) *RoomHandler {
	return &RoomHandler{svc: svc}
}

func (h *RoomHandler) GetRooms(c *gin.Context) {
	var rented *int
	if v := c.Query("rented"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || (n != 0 && n != 1) {
			middleware.Error(c, "rented参数无效，只能为0或1")
			return
		}
		rented = &n
	}
	rooms, err := h.svc.GetRoomList(rented)
	if err != nil {
		middleware.Error(c, err.Error())
		return
	}
	middleware.Success(c, rooms)
}

func (h *RoomHandler) UpdateRoom(c *gin.Context) {
	roomNo := c.Param("room_no")
	var req service.UpdateRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.Error(c, "请求参数格式错误")
		return
	}
	if err := h.svc.UpdateRoom(roomNo, req); err != nil {
		middleware.Error(c, err.Error())
		return
	}
	middleware.Success(c, nil)
}
