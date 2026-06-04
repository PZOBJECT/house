package main

import (
	"fmt"
	"log"

	"github.com/PZOBJECT/house/backend/config"
	"github.com/PZOBJECT/house/backend/handler"
	"github.com/PZOBJECT/house/backend/model"
	"github.com/PZOBJECT/house/backend/repository"
	"github.com/PZOBJECT/house/backend/router"
	"github.com/PZOBJECT/house/backend/service"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 连接MySQL
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Mysql.User,
		cfg.Mysql.Password,
		cfg.Mysql.Address,
		cfg.Mysql.Port,
		cfg.Mysql.Dbname,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	log.Println("数据库连接成功")

	// 自动迁移表结构
	if err := db.AutoMigrate(&model.Room{}, &model.Bill{}); err != nil {
		log.Fatalf("自动迁移失败: %v", err)
	}
	log.Println("表结构迁移完成")

	// 初始化房间数据
	initRooms(db)

	// 初始化各层依赖
	roomRepo := repository.NewRoomRepo(db)
	billRepo := repository.NewBillRepo(db)
	roomSvc := service.NewRoomService(roomRepo)
	billSvc := service.NewBillService(billRepo, roomRepo, db)
	roomH := handler.NewRoomHandler(roomSvc)
	billH := handler.NewBillHandler(billSvc)

	// 启动Gin
	r := gin.Default()
	router.SetupRoutes(r, roomH, billH)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("服务启动于端口 %d", cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

func initRooms(db *gorm.DB) {
	var count int64
	db.Model(&model.Room{}).Count(&count)
	if count > 0 {
		return
	}

	var rooms []model.Room
	// 改为：2-5层，每层 9 个房间（201-209、301-309...501-509）
	for floor := 2; floor <= 5; floor++ {
		for no := 1; no <= 9; no++ {
			roomNo := fmt.Sprintf("%d%02d", floor, no)
			rooms = append(rooms, model.Room{
				RoomNo:     roomNo,
				RentPrice:  1200,
				ElecPrice:  1.0,
				WaterPrice: 3.0,
				IsRented:   0,
			})
		}
	}
	if err := db.Create(&rooms).Error; err != nil {
		log.Fatalf("初始化房间数据失败: %v", err)
	}
	log.Println("已初始化36间房间数据")
}
