package database

import (
	"github.com/clarissac01/Staem/graph/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func init()  {
	dsn := "host=localhost user=postgres password=clarissa dbname=Staem port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, _ = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.Migrator().DropTable(&model.User{})
	//db.Migrator().DropTable(&model.GamePromo{})
	//db.Migrator().DropTable(&model.Game{})
	//db.Migrator().DropTable(&model.GameDetail{})
	//db.Migrator().DropTable(&model.GameTag{})
	//db.Migrator().DropTable(&model.Files{})
	//db.Migrator().DropTable(&model.GameReview{})
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.UserProfile{})
	db.AutoMigrate(&model.UserNotif{})
	db.AutoMigrate(&model.UserBadge{})
	db.AutoMigrate(&model.UserFriends{})
	db.AutoMigrate(&model.UserGame{})
	db.AutoMigrate(&model.UserInventory{})
	db.AutoMigrate(&model.UserWishlist{})
	db.AutoMigrate(&model.ReportUser{})
	db.AutoMigrate(&model.Game{})
	db.AutoMigrate(&model.Files{})
	db.AutoMigrate(&model.GameDetail{})
	db.AutoMigrate(&model.GameTag{})
	db.AutoMigrate(&model.GamePromo{})
	db.AutoMigrate(&model.GameReview{})
	db.AutoMigrate(&model.GameSales{})
	db.AutoMigrate(&model.GameSlideshow{})


	db.Create([]*model.User{
		{
			Fullname: "Staem Admin",
			Username: "admin",
			Password: "admin",
			Email:    "admin@admin.com",
			Country:  "Indonesia",
			IsSuspended: false,
			Balance: 0,
		},
		{
			Fullname: "Clarissa Chuardi",
			Username: "clarissa",
			Password: "clarissa",
			Email:    "clarissachuardi01@gmail.com",
			Country:  "Indonesia",
			IsSuspended: false,
		},
	})
}


func GetDB() *gorm.DB{
	return db
}
