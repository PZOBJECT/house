package router

import (
	"github.com/PZOBJECT/house/backend/handler"
	"github.com/PZOBJECT/house/backend/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, roomH *handler.RoomHandler, billH *handler.BillHandler) {
	r.Use(middleware.CORS())

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
		api.GET("/rooms", roomH.GetRooms)
		api.PUT("/rooms/:room_no", roomH.UpdateRoom)

		api.POST("/bills/generate", billH.GenerateBills)
		api.GET("/bills/list", billH.GetBillList)
		api.PUT("/bills/:id", billH.UpdateBill)
		api.PUT("/bills/:id/pay", billH.PayBill)
		api.PUT("/bills/:id/unpay", billH.UnpayBill)
		api.GET("/bills/:id/receipt", billH.GetReceipt)
	}
}
