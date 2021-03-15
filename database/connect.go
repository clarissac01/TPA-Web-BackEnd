package database

import (
	"github.com/clarissac01/Staem/graph/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	dsn := "host=localhost user=postgres password=clarissa dbname=Staem port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, _ = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	db.Migrator().DropTable(&model.User{})
	db.Migrator().DropTable(&model.UserProfile{})
	db.Migrator().DropTable(&model.UserFriends{})
	db.Migrator().DropTable(&model.UserURL{})
	db.Migrator().DropTable(&model.UserAvatar{})
	db.Migrator().DropTable(&model.UserBackground{})
	db.Migrator().DropTable(&model.UserMiniBackground{})
	db.Migrator().DropTable(&model.WalletCode{})
	db.Migrator().DropTable(&model.UserChatSticker{})
	db.Migrator().DropTable(&model.AnimatedAvatar{})
	db.Migrator().DropTable(&model.UserCode{})
	db.Migrator().DropTable(&model.GameBadge{})
	db.Migrator().DropTable(&model.UserBadge{})
	db.Migrator().DropTable(&model.UnsuspensionRequest{})
	db.Migrator().DropTable(&model.GameItem{})
	db.Migrator().DropTable(&model.UserGameItem{})
	//db.Migrator().DropTable(&model.GameReview{})
	//db.Migrator().DropTable(&model.GamePromo{})
	//db.Migrator().DropTable(&model.Game{})
	//db.Migrator().DropTable(&model.GameDetail{})
	//db.Migrator().DropTable(&model.GameTag{})
	//db.Migrator().DropTable(&model.Files{})
	//db.Migrator().DropTable(&model.GameReview{})
	db.AutoMigrate(&model.UserComment{})
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
	db.AutoMigrate(&model.Cart{})
	db.AutoMigrate(&model.UserProfileComment{})
	db.AutoMigrate(&model.UserURL{})
	db.AutoMigrate(&model.UserAvatar{})
	db.AutoMigrate(&model.UserBackground{})
	db.AutoMigrate(&model.UserMiniBackground{})
	db.AutoMigrate(&model.WalletCode{})
	db.AutoMigrate(&model.UserChatSticker{})
	db.AutoMigrate(&model.AnimatedAvatar{})
	db.AutoMigrate(&model.UserCode{})
	db.AutoMigrate(&model.FriendRequest{})
	db.AutoMigrate(&model.GameBadge{})
	db.AutoMigrate(&model.UserBadge{})
	db.AutoMigrate(&model.UnsuspensionRequest{})
	db.AutoMigrate(&model.GameItem{})
	db.AutoMigrate(&model.UserGameItem{})
	db.AutoMigrate(&model.Market{})
	db.AutoMigrate(&model.Transaction{})
	db.AutoMigrate(&model.Country{})
	db.AutoMigrate(&model.UserActivities{})

	db.Create([]*model.User{
		{
			Fullname:    "Staem Admin",
			Username:    "admin",
			Password:    "admin",
			Email:       "admin@admin.com",
			Country:     "Indonesia",
			IsSuspended: false,
			Balance:     0,
		},
		{
			Fullname:    "Clarissa Chuardi",
			Username:    "clarissa",
			Password:    "clarissa",
			Email:       "clarissachuardi01@gmail.com",
			Country:     "Indonesia",
			IsSuspended: false,
		},
		{
			Fullname:    "Gabriella",
			Username:    "gaby",
			Password:    "gaby",
			Email:       "gaby@gmail.com",
			Country:     "Indonesia",
			IsSuspended: false,
		},
		{
			Fullname:    "Stanley",
			Username:    "stanley",
			Password:    "stanley",
			Email:       "stanley@gmail.com",
			Country:     "Indonesia",
			IsSuspended: false,
		},
		{
			Fullname:    "Cleo",
			Username:    "cleo",
			Password:    "cleo",
			Email:       "cleo@gmail.com",
			Country:     "Indonesia",
			IsSuspended: false,
		},
	})
	db.Create([]*model.UserProfile{
		{
			Userid:  1,
			Status:  "Offline",
			Image:   -1,
			Summary: "",
			Level:   0,
		},
		{
			Userid:  2,
			Status:  "Offline",
			Image:   -1,
			Summary: "",
			Level:   0,
		},
		{
			Userid:  3,
			Status:  "Offline",
			Image:   -1,
			Summary: "",
			Level:   0,
		},
		{
			Userid:  5,
			Status:  "Offline",
			Image:   -1,
			Summary: "",
			Level:   0,
		},
		{
			Userid:  4,
			Status:  "Offline",
			Image:   -1,
			Summary: "",
			Level:   0,
		},
	})
	db.Create([]*model.UserFriends{
		{
			Userid:   2,
			Friendid: 3,
		},
		{
			Userid:   3,
			Friendid: 2,
		},
	})
	db.Create([]*model.UserURL{
		{
			Userid: 1,
			URL:    "",
		},
		{
			Userid: 2,
			URL:    "",
		},
		{
			Userid: 3,
			URL:    "",
		},
		{
			Userid: 4,
			URL:    "",
		},
		{
			Userid: 5,
			URL:    "",
		},
	})
	db.Create([]*model.UserAvatar{
		{
			Userid:   0,
			Avatarid: 1,
			Active:   false,
		},
		{
			Userid:   0,
			Avatarid: 2,
			Active:   false,
		},
		{
			Userid:   0,
			Avatarid: 3,
			Active:   false,
		},
		{
			Userid:   1,
			Avatarid: 1,
			Active:   true,
		},
		{
			Userid:   1,
			Avatarid: 2,
			Active:   false,
		},
		{
			Userid:   1,
			Avatarid: 3,
			Active:   false,
		},
		{
			Userid:   2,
			Avatarid: 1,
			Active:   true,
		},
		{
			Userid:   2,
			Avatarid: 2,
			Active:   false,
		},
		{
			Userid:   2,
			Avatarid: 3,
			Active:   false,
		},
		{
			Userid:   3,
			Avatarid: 1,
			Active:   true,
		},
		{
			Userid:   3,
			Avatarid: 2,
			Active:   false,
		},
		{
			Userid:   3,
			Avatarid: 3,
			Active:   false,
		},
	})
	db.Create([]*model.UserBackground{
		{
			Userid:       0,
			Backgroundid: 1,
			Active:       true,
		},
		{
			Userid:       0,
			Backgroundid: 2,
			Active:       false,
		},
		{
			Userid:       0,
			Backgroundid: 3,
			Active:       true,
		},
		{
			Userid:       2,
			Backgroundid: 1,
			Active:       true,
		},
		{
			Userid:       2,
			Backgroundid: 2,
			Active:       false,
		},
		{
			Userid:       3,
			Backgroundid: 1,
			Active:       true,
		},
		{
			Userid:       3,
			Backgroundid: 2,
			Active:       false,
		},
	})
	db.Create([]*model.UserMiniBackground{
		{
			Userid:       0,
			Backgroundid: 1,
			Active:       true,
		},
		{
			Userid:       0,
			Backgroundid: 2,
			Active:       false,
		},
		{
			Userid:       0,
			Backgroundid: 3,
			Active:       true,
		},
		{
			Userid:       2,
			Backgroundid: 1,
			Active:       true,
		},
		{
			Userid:       2,
			Backgroundid: 2,
			Active:       false,
		},
		{
			Userid:       3,
			Backgroundid: 1,
			Active:       true,
		},
		{
			Userid:       3,
			Backgroundid: 2,
			Active:       false,
		},
	})
	db.Create([]*model.WalletCode{
		{
			Code:   100,
			Amount: 100000,
		},
		{
			Code:   200,
			Amount: 200000,
		},
		{
			Code:   300,
			Amount: 300000,
		},
		{
			Code:   50,
			Amount: 50000,
		},
	})
	db.Create([]*model.UserChatSticker{
		{
			Userid:      0,
			Stickerid:   1,
			Chatsticker: "https://cdn.cloudflare.steamstatic.com/steamcommunity/public/images/items/1091500/bc003c956445fe939986abf0ccfd903eec0de22b.png",
			Active:      false,
		},
		{
			Userid:      0,
			Stickerid:   2,
			Chatsticker: "https://cdn.cloudflare.steamstatic.com/steamcommunity/public/images/items/440/3ffc77baebb36beb3e054d30a3bae2e20a9a2d51.png",
			Active:      false,
		},
		{
			Userid:      0,
			Stickerid:   3,
			Chatsticker: "https://cdn.cloudflare.steamstatic.com/steamcommunity/public/images/items/730/c81eec46e75c35bcd996ae2621d124bdcfa5589d.png",
			Active:      false,
		},
		{
			Userid:      2,
			Stickerid:   1,
			Chatsticker: "https://cdn.cloudflare.steamstatic.com/steamcommunity/public/images/items/730/c81eec46e75c35bcd996ae2621d124bdcfa5589d.png",
			Active:      true,
		},
	})
	db.Create([]*model.AnimatedAvatar{
		{
			Userid:   0,
			Avatarid: 1,
			Avatar:   "https://cdn.cloudflare.steamstatic.com/steamcommunity/public/images/items/1526200/1b9d00f090479d24604d15b3b8a385ad7de6433d.gif",
			Active:   false,
		},
		{
			Userid:   0,
			Avatarid: 2,
			Avatar:   "https://cdn.cloudflare.steamstatic.com/steamcommunity/public/images/items/1263950/4e0f9df2a984e2208844614afdfb59c8f903b7a6.gif",
			Active:   false,
		},
		{
			Userid:   0,
			Avatarid: 3,
			Avatar:   "https://cdn.cloudflare.steamstatic.com/steamcommunity/public/images/items/619150/9663aeacebed0d39ebbb07cda0352cc4a26249a3.gif",
			Active:   false,
		},
		{
			Userid:   2,
			Avatarid: 1,
			Avatar:   "https://cdn.cloudflare.steamstatic.com/steamcommunity/public/images/items/1526200/1b9d00f090479d24604d15b3b8a385ad7de6433d.gif",
			Active:   true,
		},
	})
	db.Create([]*model.UserCode{
		{
			Userid: 2,
			Code:   12345,
		},
		{
			Userid: 3,
			Code:   12346,
		},
		{
			Userid: 4,
			Code:   12345,
		},
		{
			Userid: 5,
			Code:   12346,
		},
	})
	db.Create([]*model.GameBadge{
		{
			Gameid:  7,
			Badgeid: 1,
			Badge:   "https://www.pngitem.com/pimgs/m/394-3941811_venusaur-gx-pokemon-card-hd-png-download.png",
		},
		{
			Gameid:  7,
			Badgeid: 2,
			Badge:   "https://i.pinimg.com/originals/c0/00/a8/c000a83578aaad8beb8e5da41f06f285.png",
		},
		{
			Gameid:  7,
			Badgeid: 3,
			Badge:   "https://www.pngitem.com/pimgs/m/134-1349562_pokemon-cards-unbroken-bonds-hd-png-download.png",
		},
		{
			Gameid:  11,
			Badgeid: 4,
			Badge:   "https://www.kindpng.com/picc/m/220-2204354_pikachu-zekrom-gx-tag-team-pokmon-card-pokemon.png",
		},
		{
			Gameid:  11,
			Badgeid: 5,
			Badge:   "https://www.kindpng.com/picc/m/79-792301_pokemon-cards-unbroken-bonds-hd-png-download.png",
		},
		{
			Gameid:  11,
			Badgeid: 6,
			Badge:   "https://www.pngitem.com/pimgs/m/156-1563584_pokemon-cards-unified-minds-hd-png-download.png",
		},
	})
	db.Create([]*model.UserBadge{
		{
			Userid: 2,
			Gameid: 7,
			Badge:  1,
		},
		{
			Userid: 2,
			Gameid: 7,
			Badge:  2,
		},
		{
			Userid: 2,
			Gameid: 7,
			Badge:  3,
		},
		{
			Userid: 2,
			Gameid: 11,
			Badge:  4,
		},
	})
	db.Create([]*model.GameItem{
		{
			Gameid:  7,
			Itemid:  1,
			Itemn:   "https://i.pinimg.com/originals/94/9b/80/949b80956f246b74dc1f4f1f476eb9c1.png",
			Summary: "Coin",
			Name:    "Gain 1 Coin",
		},
		{
			Gameid:  7,
			Itemid:  2,
			Itemn:   "https://i.pinimg.com/originals/5c/a0/27/5ca0275c8af6763c734aecae4235ef05.png",
			Summary: "Bob-omb",
			Name:    "Explodes after a brief pause",
		},
		{
			Gameid:  7,
			Itemid:  3,
			Itemn:   "https://static.wikia.nocookie.net/nintendo/images/5/55/MK8_Fire_Flower.png/revision/latest/scale-to-width-down/110?cb=20170206125257&path-prefix=en",
			Summary: "Fire Flower",
			Name:    "Shoot Fireballs",
		},
		{
			Gameid:  11,
			Itemid:  4,
			Itemn:   "https://mario.wiki.gallery/images/thumb/2/2a/GoldFlower.png/1200px-GoldFlower.png",
			Summary: "Gold Flower",
			Name:    "Shoot Gold Fireballs",
		},
		{
			Gameid:  11,
			Itemid:  5,
			Itemn:   "https://static.wikia.nocookie.net/nintendo/images/8/84/Penguin_Suit_NSMBW.png/revision/latest/scale-to-width-down/110?cb=20120110143846&path-prefix=en",
			Summary: "Penguin Suit",
			Name:    "Knock Out Enemies",
		},
		{
			Gameid:  11,
			Itemid:  6,
			Itemn:   "https://static.wikia.nocookie.net/nintendo/images/e/ec/Invincibility_Leaf.png/revision/latest/scale-to-width-down/100?cb=20150419131352&path-prefix=en",
			Summary: "Invincibility Leaf",
			Name:    "Break Blocks with Tail",
		},
	})
	db.Create([]*model.UserGameItem{
		{
			Gameid: 7,
			Itemid: 1,
			Userid: 2,
		},
		{
			Gameid: 7,
			Itemid: 2,
			Userid: 2,
		},
		{
			Gameid: 11,
			Itemid: 4,
			Userid: 2,
		},
	})
}

func GetDB() *gorm.DB {
	return db
}
