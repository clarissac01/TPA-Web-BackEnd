package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/clarissac01/Staem/database"
	"github.com/clarissac01/Staem/graph/function"
	"github.com/clarissac01/Staem/graph/generated"
	"github.com/clarissac01/Staem/graph/model"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-gomail/gomail"
)

func (r *filesResolver) File(ctx context.Context, obj *model.Files) (*graphql.Upload, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) Login(ctx context.Context, username string, password string) (string, error) {
	var user model.User
	if err := database.GetDB().First(&user, "username=? and password=?", username, password).Error; err != nil {
		return "", err
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userid": user.ID,
	})

	jwt, _ := token.SignedString([]byte("mimi"))

	//cacheKey := fmt.Sprintf(user.Username)
	//if err := database.UseCache().Set(ctx, cacheKey, jwt, 10*time.Second).Err(); err != nil {
	//	log.Fatal(err)
	//}

	return jwt, nil
}

func (r *mutationResolver) Register(ctx context.Context, fullname string, username string, password string, email string, country string) (string, error) {
	user := model.User{Fullname: fullname, Username: username, Password: password, Email: email, Country: country, IsSuspended: false, Balance: 0}
	if err := database.GetDB().Create(&user).Error; err != nil {
		return "", err
	}

	profile := model.UserProfile{Userid: int(user.ID), Level: 0, Status: "Offline", Summary: "", Image: -1}
	database.GetDB().Create(&profile)

	return "success", nil
}

func (r *mutationResolver) Addgame(ctx context.Context, title string, desc string, price int, banner graphql.Upload, slideshow []*graphql.Upload, tag []string, developer string, publisher string, systemreq string, mature bool) (string, error) {
	var bannertype = ""

	if strings.Contains(banner.Filename, ".jpg") || strings.Contains(banner.ContentType, ".png") || strings.Contains(banner.ContentType, ".jpeg") {
		bannertype = "Image"
	} else {
		bannertype = "Video"
	}

	banner2, err := ioutil.ReadAll(banner.File)
	if err != nil {
		return "", err
	}

	file := model.Files{File: banner2, ContentType: bannertype}
	database.GetDB().Create(&file)

	game := model.Game{Name: title, Price: price, Banner: file.Id, MatureContent: mature}
	if err := database.GetDB().Create(&game).Error; err != nil {
		return "", err
	}

	gameid := game.ID

	for i := 0; i < len(slideshow); i++ {
		if strings.Contains(slideshow[i].Filename, ".jpg") || strings.Contains(slideshow[i].Filename, ".png") || strings.Contains(slideshow[i].Filename, ".jpeg") {
			bannertype = "Image"
		} else {
			bannertype = "Video"
		}

		slideshow2, err := ioutil.ReadAll(slideshow[i].File)
		if err != nil {
			return "", err
		}

		file := model.Files{File: slideshow2}
		if err := database.GetDB().Create(&file).Error; err != nil {
			return "", err
		}

		gameslideshow := model.GameSlideshow{Gameid: int(gameid), Link: file.Id, ContentType: bannertype}
		if err := database.GetDB().Create(&gameslideshow).Error; err != nil {
			return "", err
		}
	}

	for i := 0; i < len(tag); i++ {
		gametag := model.GameTag{Gameid: gameid, Tagname: tag[i]}
		if err := database.GetDB().Create(&gametag).Error; err != nil {
			return "", err
		}
	}

	gamedetail := model.GameDetail{ID: gameid, Hoursplayed: 0, Systemrequirements: systemreq, Publisher: publisher, Developer: developer, Description: desc}
	if err := database.GetDB().Create(&gamedetail).Error; err != nil {
		return "", err
	}

	return "success", nil
}

func (r *mutationResolver) Updategame(ctx context.Context, id int, title string, desc string, price int, banner *graphql.Upload, slideshow []*graphql.Upload, tag []string, developer string, publisher string, systemreq string) (string, error) {
	//database.GetDB().Delete(model.GameTag{Gameid: int64(id)})
	database.GetDB().Where("Gameid = ?", id).Delete(model.GameTag{})
	gd := model.GameDetail{ID: int64(id)}
	database.GetDB().Delete(&gd)

	game := model.Game{ID: int64(id)}
	database.GetDB().First(&game)

	game.Name = title
	game.Price = price

	var bannertype = ""

	if banner != nil {

		if strings.Contains(banner.Filename, ".jpg") || strings.Contains(banner.ContentType, ".png") || strings.Contains(banner.ContentType, ".jpeg") {
			bannertype = "Image"
		} else {
			bannertype = "Video"
		}

		banner2, err := ioutil.ReadAll(banner.File)
		if err != nil {
			return "", err
		}
		file := model.Files{File: banner2}
		database.GetDB().Create(&file)
		game.Banner = file.Id

	}

	database.GetDB().Save(&game)

	if len(slideshow) != 0 {
		//database.GetDB().Delete(model.GameSlideshow{Gameid: int(id)})
		database.GetDB().Where("Gameid = ?", id).Delete(model.GameSlideshow{})
		for i := 0; i < len(slideshow); i++ {
			if strings.Contains(slideshow[i].Filename, ".jpg") || strings.Contains(slideshow[i].Filename, ".png") || strings.Contains(slideshow[i].Filename, ".jpeg") {
				bannertype = "Image"
			} else {
				bannertype = "Video"
			}

			slideshow2, err := ioutil.ReadAll(slideshow[i].File)
			if err != nil {
				return "", err
			}

			file := model.Files{File: slideshow2, ContentType: bannertype}
			if err := database.GetDB().Create(&file).Error; err != nil {
				return "", err
			}

			gameslideshow := model.GameSlideshow{Gameid: int(id), Link: file.Id, ContentType: bannertype}
			if err := database.GetDB().Create(&gameslideshow).Error; err != nil {
				return "", err
			}
		}

	}

	for i := 0; i < len(tag); i++ {
		gametag := model.GameTag{Gameid: int64(id), Tagname: tag[i]}
		if err := database.GetDB().Create(&gametag).Error; err != nil {
			return "", err
		}
	}

	gamedetail := model.GameDetail{ID: int64(id), Hoursplayed: gd.Hoursplayed, Systemrequirements: systemreq, Publisher: publisher, Developer: developer, Description: desc}
	if err := database.GetDB().Create(&gamedetail).Error; err != nil {
		return "", err
	}

	return "success", nil
}

func (r *mutationResolver) Addreview(ctx context.Context, userid int, gameid int, review string, files *graphql.Upload, positive bool) (string, error) {
	var contentType string
	if files != nil {
		println("hello")
		if strings.Contains(files.Filename, ".jpg") || strings.Contains(files.Filename, ".png") || strings.Contains(files.Filename, ".jpeg") {
			contentType = "Image"
		} else {
			contentType = "Video"
		}
		files2, err := ioutil.ReadAll(files.File)
		if err != nil {
			return "Error", err
		}
		file := model.Files{File: files2, ContentType: contentType}
		if err := database.GetDB().Create(&file).Error; err != nil {
			return "", err
		}

		database.GetDB().Create(model.GameReview{Review: review, Positive: positive, ContentType: contentType, Link: file.Id, Upvote: 0, Downvote: 0, Date: time.Now(), Gameid: gameid, Userid: userid, Helpful: 0})

	} else {
		database.GetDB().Create(model.GameReview{Userid: userid, Gameid: gameid, Date: time.Now(), Downvote: 0, Upvote: 0, Link: 0, ContentType: "", Positive: positive, Review: review, Helpful: 0})
	}

	return "success", nil
}

func (r *mutationResolver) Deletegame(ctx context.Context, id int) (string, error) {
	database.GetDB().Delete(model.Game{ID: int64(id)})
	database.GetDB().Where("Gameid = ?", id).Delete(model.GameSlideshow{})
	database.GetDB().Where("Gameid = ?", id).Delete(model.GameTag{})
	database.GetDB().Where("Gameid = ?", id).Delete(model.GamePromo{})
	database.GetDB().Where("Gameid = ?", id).Delete(model.GameReview{})
	database.GetDB().Delete(model.GameDetail{ID: int64(id)})
	return "", nil
}

func (r *mutationResolver) Addpromo(ctx context.Context, id int, discount int, validTo time.Time) (string, error) {
	gamepromo := model.GamePromo{Gameid: int64(id), Discount: int(discount), ValidTo: validTo}
	database.GetDB().Create(&gamepromo)

	var promo model.User
	database.GetDB().Raw("select * from users where id in (select userid from user_wishlists where gameid = ?)", id).Find(&promo)

	if promo.ID > 0 {
		m := gomail.NewMessage()
		m.SetHeader("From", "store.staempowered@gmail.com")
		m.SetHeader("To", "clarissachuardi01@gmail.com")
		m.SetHeader("Subject", "Sale")
		m.SetBody("text", "This game is on sale. Get it now for amazing deal!")

		d := gomail.NewDialer("smtp.gmail.com", 587, "store.staempowered@gmail.com", "zwhrruqnoyhebgea")

		if err := d.DialAndSend(m); err != nil {
			panic(err)
		}
	}

	return "", nil
}

func (r *mutationResolver) Updatepromo(ctx context.Context, id int, discount int, validTo time.Time) (string, error) {
	var promo model.GamePromo
	database.GetDB().Where("gameid = ?", id).Find(&promo)
	promo.Discount = discount
	promo.ValidTo = validTo

	database.GetDB().Where("gameid = ?", id).Save(&promo)

	return "success", nil
}

func (r *mutationResolver) Deletepromo(ctx context.Context, id int) (string, error) {
	promo := model.GamePromo{Gameid: int64(id)}

	database.GetDB().Where("gameid = ?", id).Delete(&promo)

	return "", nil
}

func (r *mutationResolver) SetUserStatus(ctx context.Context, id int) (string, error) {
	userprofile := model.UserProfile{Userid: id}
	database.GetDB().Find(&userprofile)

	userprofile.Status = "Online"
	database.GetDB().Save(&userprofile)

	return "success", nil
}

func (r *mutationResolver) AddCart(ctx context.Context, gameid int) (string, error) {
	var findgame model.Cart

	database.GetDB().Find(&findgame, gameid)

	cart := model.Cart{Gameid: gameid}
	database.GetDB().Create(&cart)

	return "success", nil
}

func (r *mutationResolver) RemoveCart(ctx context.Context, gameid int) (string, error) {
	cart := model.Cart{Gameid: gameid}
	database.GetDB().Where("gameid = ?", gameid).Delete(&cart)
	return "success", nil
}

func (r *mutationResolver) BuyGameNotWallet(ctx context.Context, gameid int, userid int) (string, error) {
	database.GetDB().Create(model.UserGame{Gameid: gameid, Userid: userid})
	database.GetDB().Where("gameid = ?", gameid).Delete(model.Cart{})
	m := gomail.NewMessage()
	m.SetHeader("From", "store.staempowered@gmail.com")
	m.SetHeader("To", "clarissachuardi01@gmail.com")
	m.SetHeader("Subject", "Purchase Transaction")
	m.SetBody("text", "You have succeeded to make a transaction! Thankyou!")

	d := gomail.NewDialer("smtp.gmail.com", 587, "store.staempowered@gmail.com", "zwhrruqnoyhebgea")

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	return "success", nil
}

func (r *mutationResolver) BuyGameWallet(ctx context.Context, gameid int, userID int, price int) (string, error) {
	database.GetDB().Create(model.UserGame{Gameid: gameid, Userid: userID})
	database.GetDB().Where("gameid = ?", gameid).Delete(model.Cart{})
	user := model.User{ID: (int64)(userID)}
	database.GetDB().Find(user)
	user.Balance = user.Balance - price

	database.GetDB().Save(&user)

	m := gomail.NewMessage()
	m.SetHeader("From", "store.staempowered@gmail.com")
	m.SetHeader("To", "clarissachuardi01@gmail.com")
	m.SetHeader("Subject", "Purchase Transaction")
	m.SetBody("text", "You have succeeded to make a transaction! Thankyou!")

	d := gomail.NewDialer("smtp.gmail.com", 587, "store.staempowered@gmail.com", "zwhrruqnoyhebgea")
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
	return "success", nil
}

func (r *mutationResolver) GiftfriendNotWallet(ctx context.Context, gameid int, senderid int, receiverid int, message *string) (string, error) {
	database.GetDB().Create(model.UserGame{Gameid: gameid, Userid: receiverid})
	database.GetDB().Where("gameid = ?", gameid).Delete(model.Cart{})
	database.GetDB().Create(model.UserNotif{Userid: receiverid, ContentType: "gift", News: "You got a gift!", Friendid: senderid})

	return "success", nil
}

func (r *mutationResolver) GiftfriendWallet(ctx context.Context, gameid int, senderid int, receiverid int, price int, message *string) (string, error) {
	database.GetDB().Create(model.UserGame{Gameid: gameid, Userid: receiverid})
	database.GetDB().Where("gameid = ?", gameid).Delete(model.Cart{})
	database.GetDB().Create(model.UserNotif{Userid: receiverid, ContentType: "gift", News: "You got a gift!", Friendid: senderid})
	user := model.User{ID: (int64)(senderid)}
	database.GetDB().Find(user)
	user.Balance = user.Balance - price

	database.GetDB().Save(&user)

	return "success", nil
}

func (r *mutationResolver) Downvotereview(ctx context.Context, gameid int, userid int, review string) (string, error) {
	var reviews model.GameReview
	database.GetDB().Raw("UPDATE game_reviews SET downvote = downvote+1 where gameid = ? and userid = ? and review = ?", gameid, userid, review).Find(&reviews)

	return "success", nil
}

func (r *mutationResolver) Upvotereview(ctx context.Context, gameid int, userid int, review string) (string, error) {
	var reviews model.GameReview
	database.GetDB().Raw("UPDATE game_reviews SET upvote = upvote+1 where gameid = ? and userid = ? and review = ?", gameid, userid, review).Find(&reviews)

	return "success", nil
}

func (r *mutationResolver) SendOtp(ctx context.Context, input int) (string, error) {
	code := generateOTP()

	m := gomail.NewMessage()
	m.SetHeader("From", "store.staempowered@gmail.com")
	m.SetHeader("To", "clarissachuardi01@gmail.com")
	m.SetHeader("Subject", "OTP Code From Staem")
	m.SetBody("text", code)

	d := gomail.NewDialer("smtp.gmail.com", 587, "store.staempowered@gmail.com", "zwhrruqnoyhebgea")

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	return code, nil
}

func (r *mutationResolver) AddReviewComment(ctx context.Context, commenterid int, posterid int, review string, comment string, gameid int) (string, error) {
	database.GetDB().Create(model.UserComment{Gameid: gameid, Review: review, Comment: comment, Commenterid: commenterid, Posterid: posterid})

	return "success", nil
}

func (r *mutationResolver) Unhelpful(ctx context.Context, gameid int, review string, userid int) (string, error) {
	var reviews model.GameReview
	database.GetDB().Raw("UPDATE game_reviews SET helpful = helpful-1 where gameid = ? and userid = ? and review = ?", gameid, userid, review).Find(&reviews)

	return "success", nil
}

func (r *mutationResolver) Helpful(ctx context.Context, gameid int, review string, userid int) (string, error) {
	var reviews model.GameReview
	database.GetDB().Raw("UPDATE game_reviews SET helpful = helpful+1 where gameid = ? and userid = ? and review = ?", gameid, userid, review).Find(&reviews)

	return "success", nil
}

func (r *mutationResolver) UpdateUserProfile(ctx context.Context, userid int, url string, fullname string, username string, summary string) (string, error) {
	var user model.User
	database.GetDB().Where("id = ?", userid).Find(&user)
	var userprofile model.UserProfile
	database.GetDB().Where("userid = ?", userid).Find(&userprofile)
	var userurl model.UserURL
	database.GetDB().Where("userid = ?", userid).Find(&userurl)

	user.Username = username
	user.Fullname = fullname

	userprofile.Summary = summary

	database.GetDB().Save(&user)
	database.GetDB().Where("userid = ?", userid).Save(&userprofile)

	if userurl.Userid == userid {
		userurl.URL = url
		database.GetDB().Where("userid = ?", userid).Save(&userurl)
	} else {
		database.GetDB().Create(model.UserURL{URL: url, Userid: userid})
	}

	return "success", nil
}

func (r *mutationResolver) UpdateUserAvatar(ctx context.Context, userid int, avatarframeid int) (string, error) {
	var useravatar model.UserProfile
	var avatarmodel model.UserAvatar

	database.GetDB().Where("userid = ?", userid).Find(&useravatar)
	database.GetDB().Where("userid = ? and active = true", userid).Find(&avatarmodel)
	avatarmodel.Active = false
	database.GetDB().Where("userid = ? and active = true", userid).Save(&avatarmodel)

	var mdl model.UserAvatar
	database.GetDB().Where("userid = ? and avatarid = ?", userid, avatarframeid).Find(&mdl)
	mdl.Active = true
	database.GetDB().Where("userid = ? and avatarid = ?", userid, avatarframeid).Save(&mdl)

	return "success", nil
}

func (r *mutationResolver) UpdateUserBackground(ctx context.Context, userid int, backgroundid int) (string, error) {
	var useravatar model.UserBackground
	database.GetDB().Where("userid = ? and active = true", userid).Find(&useravatar)
	useravatar.Active = false
	database.GetDB().Where("userid = ? and active = true", userid).Save(&useravatar)

	var background2 model.UserBackground
	database.GetDB().Where("userid = ? and backgroundid = ?", userid, backgroundid).Find(&background2)
	background2.Active = true
	database.GetDB().Where("userid = ? and backgroundid = ?", userid, backgroundid).Save(&background2)

	return "success", nil
}

func (r *mutationResolver) UpdateUserMiniBackground(ctx context.Context, userid int, minibackgroundid int) (string, error) {
	var useravatar model.UserMiniBackground
	database.GetDB().Where("userid = ? and active = true", userid).Find(&useravatar)
	useravatar.Active = false
	database.GetDB().Where("userid = ? and active = true", userid).Save(&useravatar)

	var background2 model.UserMiniBackground
	database.GetDB().Where("userid = ? and backgroundid = ?", userid, minibackgroundid).Find(&background2)
	background2.Active = true
	database.GetDB().Where("userid = ? and backgroundid = ?", userid, minibackgroundid).Save(&background2)

	return "success", nil
}

func (r *mutationResolver) UpdateUserAnimated(ctx context.Context, userid int, avatarid int) (string, error) {
	var a model.AnimatedAvatar
	a.Active = false
	database.GetDB().Where("active = true").Save(&a)

	var animated model.AnimatedAvatar
	database.GetDB().Where("userid = ? and avatarid = ?", userid, avatarid).Find(&animated)

	animated.Active = true
	database.GetDB().Where("userid = ? and avatarid = ?", userid, avatarid).Save(&animated)

	return "success", nil
}

func (r *mutationResolver) UpdateUserHex(ctx context.Context, userid *int, hex string) (string, error) {
	var userhex model.UserTheme
	database.GetDB().Where("userid = ?", userid).Find(&userhex)
	userhex.Hex = hex

	return "success", nil
}

func (r *mutationResolver) UpdateWallet(ctx context.Context, userid int, code int) (string, error) {
	var wallet model.User
	database.GetDB().Where("id = ?", userid).Find(&wallet)

	var balance model.WalletCode
	if err := database.GetDB().Where("code = ?", code).First(&balance).Error; err != nil {
		print(err)
		return "error", err
	}

	wallet.Balance = wallet.Balance + balance.Amount

	return "success", database.GetDB().Where("id = ?", userid).Save(&wallet).Error
}

func (r *mutationResolver) BuyChatStickers(ctx context.Context, userid int, stickerid int) (string, error) {
	var cs model.UserChatSticker
	database.GetDB().Where("userid = 0 and stickerid = ?", stickerid).Find(&cs)

	database.GetDB().Create(model.UserChatSticker{Active: false, Chatsticker: cs.Chatsticker, Stickerid: stickerid, Userid: userid, Price: cs.Price})

	var user model.User
	database.GetDB().Where("id = ?", userid).Find(&user)
	user.Point = user.Point - cs.Price
	database.GetDB().Where("id = ?", user).Save(&user)

	return "", nil
}

func (r *mutationResolver) BuyProfileBackground(ctx context.Context, userid int, backgroundid int) (string, error) {
	var cs model.UserBackground
	database.GetDB().Where("userid = 0 and backgroundid = ?", backgroundid).Find(&cs)

	database.GetDB().Create(model.UserBackground{Price: cs.Price, Userid: userid, Active: false, Backgroundid: backgroundid})

	var user model.User
	database.GetDB().Where("id = ?", userid).Find(&user)
	user.Point = user.Point - cs.Price
	database.GetDB().Where("id = ?", user).Save(&user)

	return "", nil
}

func (r *mutationResolver) BuyMiniBackground(ctx context.Context, userid int, backgroundid int) (string, error) {
	var cs model.UserMiniBackground
	database.GetDB().Where("userid = 0 and backgroundid = ?", backgroundid).Find(&cs)

	database.GetDB().Create(model.UserMiniBackground{Price: cs.Price, Userid: userid, Active: false, Backgroundid: backgroundid})

	var user model.User
	database.GetDB().Where("id = ?", userid).Find(&user)
	user.Point = user.Point - cs.Price
	database.GetDB().Where("id = ?", user).Save(&user)

	return "", nil
}

func (r *mutationResolver) BuyAvatarFrame(ctx context.Context, userid int, avatarid int) (string, error) {
	var cs model.UserAvatar
	database.GetDB().Where("userid = 0 and avatarid = ?", avatarid).Find(&cs)

	database.GetDB().Create(model.UserAvatar{Price: cs.Price, Userid: userid, Active: false, Avatarid: avatarid})

	var user model.User
	database.GetDB().Where("id = ?", userid).Find(&user)
	user.Point = user.Point - cs.Price
	database.GetDB().Where("id = ?", user).Save(&user)

	return "", nil
}

func (r *mutationResolver) BuyAnimated(ctx context.Context, userid int, animatedid int) (string, error) {
	var cs model.AnimatedAvatar
	database.GetDB().Where("userid = 0 and avatarid = ?", animatedid).Find(&cs)

	database.GetDB().Create(model.AnimatedAvatar{Price: cs.Price, Userid: userid, Active: false, Avatarid: animatedid, Avatar: cs.Avatar})

	var user model.User
	database.GetDB().Where("id = ?", userid).Find(&user)
	user.Point = user.Point - cs.Price
	database.GetDB().Where("id = ?", user).Save(&user)

	return "", nil
}

func (r *mutationResolver) AddWishlist(ctx context.Context, userid int, gameid int) (string, error) {
	database.GetDB().Create(model.UserWishlist{Userid: userid, Gameid: gameid})

	var promo model.GamePromo
	database.GetDB().Where("gameid = ?", gameid).First(&promo)
	fmt.Print(promo.Discount)

	if int(promo.Gameid) == gameid && promo.Discount > 0 {
		m := gomail.NewMessage()
		m.SetHeader("From", "store.staempowered@gmail.com")
		m.SetHeader("To", "clarissachuardi01@gmail.com")
		m.SetHeader("Subject", "Sale")
		m.SetBody("text", "This game is on sale. Get it now for amazing deal!")

		d := gomail.NewDialer("smtp.gmail.com", 587, "store.staempowered@gmail.com", "zwhrruqnoyhebgea")

		if err := d.DialAndSend(m); err != nil {
			panic(err)
		}
	}

	return "success", nil
}

func (r *mutationResolver) RemoveWishlist(ctx context.Context, userid int, gameid int) (string, error) {
	cart := model.UserWishlist{Gameid: gameid}
	database.GetDB().Where("gameid = ? and userid = ?", gameid, userid).Delete(&cart)
	return "success", nil
}

func (r *mutationResolver) ReportUser(ctx context.Context, userid int, reporterid int) (string, error) {
	database.GetDB().Create(model.ReportUser{Userwhoreportid: reporterid, Reporteduserid: userid, CreatedAt: time.Now()})

	var reports []*model.ReportUser
	database.GetDB().Where("reporteduserid = ? and created_at BETWEEN (now() - '1 week'::interval) AND now()", userid).Find(&reports)

	if len(reports) >= 5 {
		var user model.User
		database.GetDB().Where("id = ?", userid).Find(&user)
		user.IsSuspended = true
		database.GetDB().Where("id = ?", userid).Save(&user)
		database.GetDB().Where("reporteduserid = ?", userid).Delete(&reports)
	}

	return "success", nil
}

func (r *mutationResolver) SendRequest(ctx context.Context, userid int, receiverid int) (string, error) {
	database.GetDB().Create(model.FriendRequest{Receiverid: receiverid, Senderid: userid})
	database.GetDB().Create(model.UserNotif{Userid: receiverid, Friendid: userid, ContentType: "Invites", News: "You have a new friend request!"})

	return "success", nil
}

func (r *mutationResolver) AcceptRequest(ctx context.Context, userid int, senderid int) (string, error) {
	database.GetDB().Where("senderid = ? and receiverid = ?", senderid, userid).Delete(model.FriendRequest{Senderid: senderid, Receiverid: userid})

	database.GetDB().Create(model.UserFriends{Userid: userid, Friendid: senderid})
	database.GetDB().Create(model.UserFriends{Userid: senderid, Friendid: userid})

	return "success", nil
}

func (r *mutationResolver) DeclineRequest(ctx context.Context, userid int, senderid int) (string, error) {
	database.GetDB().Where("senderid = ? and receiverid = ?", senderid, userid).Delete(model.FriendRequest{Senderid: senderid, Receiverid: userid})

	return "success", nil
}

func (r *mutationResolver) RequestUnsuspension(ctx context.Context, userid int) (string, error) {
	database.GetDB().Create(model.UnsuspensionRequest{Userid: userid, Status: "Pending"})

	return "sucess", nil
}

func (r *mutationResolver) AcceptUnsuspension(ctx context.Context, userid int) (string, error) {
	var user model.User
	database.GetDB().Where("id = ?", userid).Find(&user)
	user.IsSuspended = false
	database.GetDB().Where("id = ?", userid).Save(&user)

	var req model.UnsuspensionRequest
	database.GetDB().Where("userid = ? and status = 'Pending'", userid).Find(&req)
	req.Status = "Accepted"
	database.GetDB().Where("userid = ? and status = 'Pending'", userid).Save(&req)

	return "success", nil
}

func (r *mutationResolver) DeclineUnsuspension(ctx context.Context, userid int) (string, error) {
	var req model.UnsuspensionRequest
	database.GetDB().Where("userid = ? and status = 'Pending'", userid).Find(&req)
	req.Status = "Declined"
	database.GetDB().Where("userid = ? and status = 'Pending'", userid).Save(&req)

	return "success", nil
}

func (r *mutationResolver) SellItem(ctx context.Context, userid int, itemid int, price int) (string, error) {
	var item model.GameItem
	database.GetDB().Where("itemid = ?", itemid).Find(&item)
	database.GetDB().Create(model.Market{Itemid: itemid, Gameid: item.Gameid, Price: price, Sellerid: userid, Type: "Sell"})
	return "success", nil
}

func (r *mutationResolver) SendMessage(ctx context.Context, userid int, message string) (string, error) {
	r.ChatSocket[userid] <- message

	return message, nil
}

func (r *mutationResolver) BuyfromMarketDetail(ctx context.Context, sellerid int, userid int, price int, itemid int) (string, error) {
	var market model.Market
	database.GetDB().Where("itemid = ? and price = ? and type = 'Sell'", itemid, price).Find(&market)

	if market.Sellerid > 0 {
		database.GetDB().Where("sellerid = ? and itemid = ? and price = ? and type = 'Sell'", market.Sellerid, itemid, price).Delete(model.Market{})
		database.GetDB().Create(model.Transaction{Sellerid: market.Sellerid, Price: price, Itemid: itemid, CreatedAt: time.Now(), Buyerid: userid})

		database.GetDB().Where("userid = ? and itemid = ?", market.Sellerid, itemid).Delete(model.UserGameItem{})

		m := gomail.NewMessage()
		m.SetHeader("From", "store.staempowered@gmail.com")
		m.SetHeader("To", "clarissachuardi01@gmail.com")
		m.SetHeader("Subject", "Sell Success")
		m.SetBody("text", "The item you sell have been sold by other user!")

		d := gomail.NewDialer("smtp.gmail.com", 587, "store.staempowered@gmail.com", "zwhrruqnoyhebgea")

		if err := d.DialAndSend(m); err != nil {
			panic(err)
		}

		var seller model.User
		database.GetDB().Where("id = ?", market.Sellerid).Find(&seller)
		seller.Balance = seller.Balance + price
		database.GetDB().Where("id = ?", market.Sellerid).Save(&seller)

		ma := gomail.NewMessage()
		ma.SetHeader("From", "clarissachuardi01@gmail.com")
		ma.SetHeader("To", "gaby@gmail.com")
		ma.SetHeader("Subject", "Buy Success")
		ma.SetBody("text", "You have succeeded to buy an item!")

		e := gomail.NewDialer("smtp.gmail.com", 587, "store.staempowered@gmail.com", "zwhrruqnoyhebgea")

		if err := e.DialAndSend(ma); err != nil {
			panic(err)
		}

		var buyer model.User
		database.GetDB().Where("id = ?", userid).Find(&buyer)
		buyer.Balance = buyer.Balance - price
		database.GetDB().Where("id = ?", userid).Save(&buyer)

		var game model.GameItem
		database.GetDB().Where("itemid = ?", itemid).Find(&game)
		database.GetDB().Create(model.UserGameItem{Itemid: itemid, Gameid: game.Gameid, Userid: userid})
		return "success", nil

	}

	database.GetDB().Create(model.Market{Sellerid: 0, Itemid: itemid, Price: price, Type: "Bid", Gameid: market.Gameid, Buyerid: userid})

	return "success", nil
}

func (r *mutationResolver) SellfromMarketDetail(ctx context.Context, sellerid int, userid int, price int, itemid int) (string, error) {
	var game model.GameItem
	database.GetDB().Where("itemid = ?", itemid).Find(&game)
	database.GetDB().Create(model.Market{Type: "Sell", Gameid: game.Gameid, Price: price, Itemid: itemid, Sellerid: userid})

	return "success", nil
}

func (r *queryResolver) Auth(ctx context.Context, token string) (*model.User, error) {
	userid, err := function.ParseToken(token)
	var user model.User
	database.GetDB().First(&user, int(userid))

	return &user, err
}

func (r *queryResolver) GetUserProfile(ctx context.Context, id int) (*model.UserProfile, error) {
	userprofile := model.UserProfile{Userid: id}
	database.GetDB().Where("userid=?", id).Find(&userprofile)

	return &userprofile, nil
}

func (r *queryResolver) Games(ctx context.Context) ([]*model.Game, error) {
	var games []*model.Game
	if err := database.GetDB().Find(&games).Error; err != nil {
		return nil, err
	}

	return games, nil
}

func (r *queryResolver) GetGame(ctx context.Context, id int) (*model.Game, error) {
	game := model.Game{ID: int64(id)}
	if err := database.GetDB().Find(&game).Error; err != nil {
		return nil, err
	}

	return &game, nil
}

func (r *queryResolver) GetGameDetail(ctx context.Context, id int) (*model.GameDetail, error) {
	gamedetail := model.GameDetail{ID: int64(id)}
	if err := database.GetDB().Find(&gamedetail).Error; err != nil {
		return nil, err
	}

	return &gamedetail, nil
}

func (r *queryResolver) GetGameTags(ctx context.Context, id int) ([]*model.GameTag, error) {
	var gametags []*model.GameTag
	if err := database.GetDB().Find(&gametags, int64(id)).Error; err != nil {
		return nil, err
	}

	return gametags, nil
}

func (r *queryResolver) GetGameSlideshows(ctx context.Context, id int) ([]*model.GameSlideshow, error) {
	var gameslideshows []*model.GameSlideshow
	if err := database.GetDB().Find(&gameslideshows, int64(id)).Error; err != nil {
		return nil, err
	}

	return gameslideshows, nil
}

func (r *queryResolver) GetPromo(ctx context.Context) ([]*model.GamePromo, error) {
	var gamepromos []*model.GamePromo
	database.GetDB().Debug().Find(&gamepromos)

	return gamepromos, nil
}

func (r *queryResolver) GetPromobyID(ctx context.Context, id int) (*model.GamePromo, error) {
	var gamepromo model.GamePromo
	if err := database.GetDB().Where("gameid = ?", id).Find(&gamepromo).Error; err != nil {
		return nil, err
	}

	return &gamepromo, nil
}

func (r *queryResolver) GameNotPromo(ctx context.Context) ([]*model.Game, error) {
	var games []*model.Game
	database.GetDB().Where("id not in (?)", database.GetDB().Table("game_promos").Select("gameid")).Find(&games)

	return games, nil
}

func (r *queryResolver) GetUserNotif(ctx context.Context, id int) ([]*model.UserNotif, error) {
	usernotif := []*model.UserNotif{}
	database.GetDB().Where("userid = ?", id).Find(&usernotif)

	return usernotif, nil
}

func (r *queryResolver) Fandr(ctx context.Context) ([]*model.Game, error) {
	var byid []*model.GameDetail
	database.GetDB().Select("id").Order("hoursplayed").Find(&byid).Limit(12)

	var games []*model.Game
	database.GetDB().Raw("SELECT * FROM games WHERE games.id in (SELECT id FROM game_details)").Find(&games)

	return games, nil
}

func (r *queryResolver) Faandrimg(ctx context.Context) ([]*model.GameSlideshow, error) {
	var byid []*model.GameDetail
	database.GetDB().Select("id").Order("hoursplayed").Find(&byid).Limit(12)

	var gamedetails []*model.GameSlideshow
	database.GetDB().Raw("SELECT * FROM game_slideshows WHERE game_slideshows.gameid in (SELECT id FROM game_details) and game_slideshows.link in (select id from files where content_type= 'Image')").Find(&gamedetails)

	return gamedetails, nil
}

func (r *queryResolver) SpecialOffer(ctx context.Context) ([]*model.Game, error) {
	var games []*model.Game

	database.GetDB().Raw("select * from games where id in (select gameid from game_promos order by discount asc limit 12)").Find(&games)

	return games, nil
}

func (r *queryResolver) SearchGame(ctx context.Context, keyword string) ([]*model.Game, error) {
	var games []*model.Game
	database.GetDB().Where("name LIKE ? OR name LIKE ?", "%"+strings.ToLower(keyword)+"%", "%"+strings.ToUpper(keyword)+"%").Find(&games).Limit(5)

	return games, nil
}

func (r *queryResolver) SearchGamePage(ctx context.Context, keyword string, countGame int) ([]*model.Game, error) {
	var games []*model.Game
	database.GetDB().Where("name LIKE ? OR name LIKE ?", "%"+strings.ToLower(keyword)+"%", "%"+strings.ToUpper(keyword)+"%").Find(&games).Limit(countGame)

	return games, nil
}

func (r *queryResolver) Communityrecommended(ctx context.Context) ([]*model.Game, error) {
	var games []*model.Game
	database.GetDB().Raw("SELECT * FROM games where id in (select gameid from game_reviews where positive = 'true' group by gameid order by count(gameid) desc limit 12)").Find(&games)

	return games, nil
}

func (r *queryResolver) FiltergameByPrice(ctx context.Context, price int, keyword string, countGame int) ([]*model.Game, error) {
	var games []*model.Game
	database.GetDB().Where("name LIKE ? OR name LIKE ?", "%"+strings.ToLower(keyword)+"%", "%"+strings.ToUpper(keyword)+"%").Find(&games, "price <= ? ", price).Limit(countGame)

	return games, nil
}

func (r *queryResolver) GetGameGenre(ctx context.Context, keyword string) ([]*model.GameTag, error) {
	var games []*model.GameTag
	database.GetDB().Raw("SELECT distinct tagname FROM game_tags WHERE gameid IN (select id from games where name like '%" + keyword + "%')").Find(&games)

	return games, nil
}

func (r *queryResolver) FiltergamebyGenre(ctx context.Context, genre string, keyword string, countGame int) ([]*model.Game, error) {
	var games []*model.Game
	database.GetDB().Raw("SELECT * FROM games WHERE id IN (SELECT gameid FROM game_tags WHERE tagname like '" + genre + "' ) and name like '%" + keyword + "%' ").Find(&games).Limit(countGame)
	println(len(games))
	return games, nil
}

func (r *queryResolver) Filtergamebyfandr(ctx context.Context, keyword string, countGame int) ([]*model.Game, error) {
	var byid []*model.GameDetail
	database.GetDB().Select("id").Order("hoursplayed").Find(&byid).Limit(12)

	var games []*model.Game
	database.GetDB().Raw("SELECT * FROM games WHERE games.id in (SELECT id FROM game_details) and name like '%" + keyword + "%'").Find(&games).Limit(countGame)

	return games, nil
}

func (r *queryResolver) FiltergamebySpecialOffer(ctx context.Context, keyword string, countGame int) ([]*model.Game, error) {
	var games []*model.Game

	database.GetDB().Raw("select * from games where id in (select gameid from game_promos order by discount asc) and name like '%" + keyword + "%'").Find(&games).Limit(countGame)

	return games, nil
}

func (r *queryResolver) FiltergamebyCommunityRec(ctx context.Context, keyword string, countGame int) ([]*model.Game, error) {
	var games []*model.Game
	database.GetDB().Raw("SELECT * FROM games where id in (select gameid from game_reviews where positive = 'true' group by gameid order by count(gameid) desc) and name like '%" + keyword + "%'").Find(&games).Limit(countGame)

	return games, nil
}

func (r *queryResolver) FiltergenrebyPrice(ctx context.Context, genre string, price int, countGame int) ([]*model.Game, error) {
	var games []*model.Game
	database.GetDB().Raw("SELECT * FROM games WHERE id IN (SELECT gameid FROM game_tags WHERE tagname like '"+genre+"' ) and price <= ?", price).Find(&games).Limit(countGame)

	return games, nil
}

func (r *queryResolver) GetCart(ctx context.Context) ([]*model.Cart, error) {
	var cart []*model.Cart

	database.GetDB().Raw("SELECT DISTINCT * from carts").Find(&cart)

	return cart, nil
}

func (r *queryResolver) Getallgamefromcart(ctx context.Context) ([]*model.Game, error) {
	var games []*model.Game

	database.GetDB().Raw("select * from games where id in (select distinct gameid from carts)").Find(&games)

	return games, nil
}

func (r *queryResolver) GetUserFriends(ctx context.Context, userid int) ([]*model.User, error) {
	var users []*model.User
	database.GetDB().Raw("select * from users where id in (select friendid from user_friends where userid = ?)", userid).Find(&users)

	return users, nil
}

func (r *queryResolver) GetUserFriendProfiles(ctx context.Context, userID int) ([]*model.UserProfile, error) {
	var users []*model.UserProfile
	database.GetDB().Raw("select * from user_profiles where userid in (select friendid from user_friends where userid = ?)", userID).Find(&users)

	return users, nil
}

func (r *queryResolver) GetUserGames(ctx context.Context, userid int) ([]*model.Game, error) {
	var games []*model.Game
	database.GetDB().Raw("Select * from games where id in (select gameid from user_games where userid = ?)", userid).Find(&games)

	return games, nil
}

func (r *queryResolver) GetGameReview1(ctx context.Context, gameid int) ([]*model.GameReview, error) {
	var gamereview []*model.GameReview
	database.GetDB().Raw("select * from game_reviews where gameid = ? and date BETWEEN (now() - '1 month'::interval) AND now() order by upvote desc", gameid).Find(&gamereview)

	return gamereview, nil
}

func (r *queryResolver) GetUserReview1(ctx context.Context, gameid int) ([]*model.User, error) {
	var user []*model.User
	database.GetDB().Raw("select * from users where id in (select userid from game_reviews where gameid = ? and date BETWEEN (now() - '1 month'::interval) AND now() order by upvote desc)", gameid).Find(&user)

	return user, nil
}

func (r *queryResolver) GetUserReview2(ctx context.Context, gameid int) ([]*model.User, error) {
	var user []*model.User
	database.GetDB().Raw("select * from users where id in (select userid from game_reviews where gameid = ? order by date asc)", gameid).Find(&user)

	return user, nil
}

func (r *queryResolver) GetGameReview2(ctx context.Context, gameid int) ([]*model.GameReview, error) {
	var gamereview []*model.GameReview
	database.GetDB().Raw("select * from game_reviews where gameid = ? order by date asc", gameid).Find(&gamereview)

	return gamereview, nil
}

func (r *queryResolver) GetMedia(ctx context.Context) ([]*model.GameReview, error) {
	var media []*model.GameReview
	database.GetDB().Raw("SELECT * FROM game_reviews where link > 0").Find(&media)

	return media, nil
}

func (r *queryResolver) GetMediaGame(ctx context.Context) ([]*model.Game, error) {
	var games []*model.Game

	database.GetDB().Raw("select * from games where id in (SELECT gameid FROM game_reviews where link > 0)").Find(&games)

	return games, nil
}

func (r *queryResolver) GetReviews(ctx context.Context) ([]*model.GameReview, error) {
	var reviews []*model.GameReview
	database.GetDB().Raw("select * from game_reviews where link = 0 order by gameid").Find(&reviews)

	return reviews, nil
}

func (r *queryResolver) GetMediaCommenter(ctx context.Context, gameid int, userid int, review string) ([]*model.User, error) {
	var user []*model.User
	database.GetDB().Raw("select * from users where id in (select commenterid from user_comments where gameid = ? and posterid = ? and review = ?)", gameid, userid, review).Find(&user)

	return user, nil
}

func (r *queryResolver) GetMediaCommenterDetail(ctx context.Context, gameid int, userid int, review string) ([]*model.UserProfile, error) {
	var user []*model.UserProfile
	database.GetDB().Raw("select * from user_profiles where userid in (select commenterid from user_comments where gameid = ? and posterid = ? and review = ?)", gameid, userid, review).Find(&user)

	return user, nil
}

func (r *queryResolver) GetMediaComments(ctx context.Context) ([]*model.UserComment, error) {
	var mediacomments []*model.UserComment
	database.GetDB().Raw("select * from user_comments where gameid in (SELECT gameid FROM game_reviews where link > 0)").Find(&mediacomments)

	return mediacomments, nil
}

func (r *queryResolver) GetReviewsComments(ctx context.Context) ([]*model.UserComment, error) {
	var mediacomments []*model.UserComment
	database.GetDB().Raw("select * from user_comments where gameid in (SELECT gameid FROM game_reviews where link = 0)").Find(&mediacomments)

	return mediacomments, nil
}

func (r *queryResolver) GetReviewsCommenter(ctx context.Context, gameid int, userid int, review string) ([]*model.User, error) {
	var user []*model.User
	database.GetDB().Raw("select u.* from user_comments uc join users u on uc.commenterid = u.id where gameid = ? and posterid = ? and uc.review = ?", gameid, userid, review).Find(&user)

	return user, nil
}

func (r *queryResolver) GetReviewsGame(ctx context.Context) ([]*model.Game, error) {
	var games []*model.Game
	database.GetDB().Raw("select * from games where id in (select gameid from game_reviews where link = 0) order by id").Find(&games)

	return games, nil
}

func (r *queryResolver) GetReviewsCommenterDetail(ctx context.Context, gameid int, userid int, review string) ([]*model.UserProfile, error) {
	var user []*model.UserProfile
	database.GetDB().Raw("select u.* from user_comments uc join user_profiles u on uc.commenterid = u.userid where gameid = ? and posterid = ? and uc.review = ?", gameid, userid, review).Find(&user)

	return user, nil
}

func (r *queryResolver) GetPoster(ctx context.Context, userid int) (*model.User, error) {
	var user model.User
	database.GetDB().Where("id = ?", userid).Find(&user)

	return &user, nil
}

func (r *queryResolver) GetUserComment(ctx context.Context, userid int) ([]*model.UserProfileComment, error) {
	var res []*model.UserProfileComment
	database.GetDB().Where("userid = ?", userid).Find(&res)

	return res, nil
}

func (r *queryResolver) GetCommentUserProfile(ctx context.Context, userid int) ([]*model.UserProfile, error) {
	var res []*model.UserProfile
	database.GetDB().Raw("select * from user_profiles where userid in (select userid from user_profile_comments where userid = ?)", userid).Find(&res)

	return res, nil
}

func (r *queryResolver) GetUserCommenter(ctx context.Context, userid int) ([]*model.User, error) {
	var res []*model.User
	database.GetDB().Raw("select * from users where id in (select userid from user_profile_comments where userid = ?)", userid).Find(&res)

	return res, nil
}

func (r *queryResolver) GetUserURL(ctx context.Context, userid int) (*model.UserURL, error) {
	var userurl model.UserURL
	database.GetDB().Where("userid = ?", userid).Find(&userurl)

	return &userurl, nil
}

func (r *queryResolver) GetUserBackground(ctx context.Context, userid int) (*model.UserBackground, error) {
	var useravatar model.UserBackground
	database.GetDB().Where("userid = ? and active = true", userid).Find(&useravatar)

	return &useravatar, nil
}

func (r *queryResolver) GetAllUserBackground(ctx context.Context, userid int) ([]*model.UserBackground, error) {
	var useravatar []*model.UserBackground
	database.GetDB().Debug().Where("userid = ?", userid).Find(&useravatar)

	return useravatar, nil
}

func (r *queryResolver) GetAllUserMiniBackground(ctx context.Context, userid int) ([]*model.UserMiniBackground, error) {
	var useravatar []*model.UserMiniBackground
	database.GetDB().Where("userid = ?", userid).Find(&useravatar)

	return useravatar, nil
}

func (r *queryResolver) GetUserMiniBackground(ctx context.Context, userid int) (*model.UserMiniBackground, error) {
	var useravatar model.UserMiniBackground
	database.GetDB().Where("userid = ? and active = true", userid).Find(&useravatar)

	return &useravatar, nil
}

func (r *queryResolver) GetUserAvatar(ctx context.Context, userid int) (*model.UserAvatar, error) {
	var useravatar model.UserAvatar
	database.GetDB().Where("userid = ? and active = true", userid).Find(&useravatar)

	return &useravatar, nil
}

func (r *queryResolver) GetAllUserAvatar(ctx context.Context, userid int) ([]*model.UserAvatar, error) {
	var useravatar []*model.UserAvatar
	database.GetDB().Where("userid = ?", userid).Find(&useravatar)

	return useravatar, nil
}

func (r *queryResolver) GetUserHex(ctx context.Context, userid int) (*model.UserTheme, error) {
	var userhex model.UserTheme
	database.GetDB().Where("userid = ?", userid).Find(&userhex)

	return &userhex, nil
}

func (r *queryResolver) GetAllUserChatStickers(ctx context.Context, userid int) ([]*model.UserChatSticker, error) {
	var chatsticker []*model.UserChatSticker
	database.GetDB().Where("userid = ?", userid).Find(&chatsticker)

	return chatsticker, nil
}

func (r *queryResolver) GetAllChatStickers(ctx context.Context) ([]*model.UserChatSticker, error) {
	var chatsticker []*model.UserChatSticker
	database.GetDB().Where("userid = 0").Find(&chatsticker)

	return chatsticker, nil
}

func (r *queryResolver) GetAllBackground(ctx context.Context) ([]*model.UserBackground, error) {
	var chatsticker []*model.UserBackground
	database.GetDB().Where("userid = 0").Find(&chatsticker)

	return chatsticker, nil
}

func (r *queryResolver) GetAllMiniBackground(ctx context.Context) ([]*model.UserMiniBackground, error) {
	var chatsticker []*model.UserMiniBackground
	database.GetDB().Where("userid = 0").Find(&chatsticker)

	return chatsticker, nil
}

func (r *queryResolver) GetAllAvatar(ctx context.Context) ([]*model.UserAvatar, error) {
	var chatsticker []*model.UserAvatar
	database.GetDB().Where("userid = 0").Find(&chatsticker)

	return chatsticker, nil
}

func (r *queryResolver) GetUserAnimated(ctx context.Context, userid int) (*model.AnimatedAvatar, error) {
	var chatsticker model.AnimatedAvatar
	database.GetDB().Where("userid = ? and active = true", userid).Find(&chatsticker)

	return &chatsticker, nil
}

func (r *queryResolver) GetAllUserAnimated(ctx context.Context, userid int) ([]*model.AnimatedAvatar, error) {
	var chatsticker []*model.AnimatedAvatar
	database.GetDB().Where("userid = ?", userid).Find(&chatsticker)

	return chatsticker, nil
}

func (r *queryResolver) GetAllAnimated(ctx context.Context) ([]*model.AnimatedAvatar, error) {
	var chatsticker []*model.AnimatedAvatar
	database.GetDB().Where("userid = 0").Find(&chatsticker)

	return chatsticker, nil
}

func (r *queryResolver) GetAllUserWishlist(ctx context.Context, userid int) ([]*model.Game, error) {
	var wishlist []*model.Game
	database.GetDB().Raw("select * from games where id in (select distinct gameid from user_wishlists where userid = ?)", userid).Find(&wishlist)

	return wishlist, nil
}

func (r *queryResolver) GetUserCode(ctx context.Context, userid int) (*model.UserCode, error) {
	var codes *model.UserCode
	database.GetDB().Where("userid = ?", userid).Find(&codes)

	return codes, nil
}

func (r *queryResolver) GetUserFriendCode(ctx context.Context, userid int) ([]*model.UserCode, error) {
	var users []*model.UserCode
	database.GetDB().Raw("select * from user_codes where userid in (select friendid from user_friends where userid = ?)", userid).Find(&users)

	return users, nil
}

func (r *queryResolver) GetPendingInvite(ctx context.Context, userid int) ([]*model.User, error) {
	var users []*model.User
	database.GetDB().Raw("select * from users where id in (select senderid from friend_requests where receiverid = ?) and id != ? and id != 1", userid, userid).Find(&users)

	return users, nil
}

func (r *queryResolver) GetSentInvite(ctx context.Context, userid int) ([]*model.User, error) {
	var users []*model.User
	database.GetDB().Raw("select * from users where id in (select receiverid from friend_requests where senderid = ?) and id != ? and id != 1", userid, userid).Find(&users)

	return users, nil
}

func (r *queryResolver) GetUserNotFriend(ctx context.Context, userid int) ([]*model.User, error) {
	var users []*model.User
	database.GetDB().Raw("select * from users where id not in (select friendid from user_friends where userid = ?) and id not in (select receiverid from friend_requests where senderid = ?) and id not in (select senderid from friend_requests where receiverid = ?) and id != ? and id != 1", userid, userid, userid, userid).Find(&users)

	return users, nil
}

func (r *queryResolver) GetUserBadges(ctx context.Context, userid int) ([]*model.GameBadge, error) {
	var badges []*model.GameBadge
	database.GetDB().Raw("select * from game_badges where badgeid in (select badge from user_badges where userid = ?) group by game_badges.badgeid, game_badges.gameid, game_badges.badge order by gameid", userid).Find(&badges)

	return badges, nil
}

func (r *queryResolver) GetUserAllBadges(ctx context.Context, userid int) ([]*model.GameBadge, error) {
	var badges []*model.GameBadge
	database.GetDB().Raw("select * from game_badges where badgeid not in (select badge from user_badges where userid = ?) group by game_badges.badgeid, game_badges.gameid, game_badges.badge order by gameid", userid).Find(&badges)

	return badges, nil
}

func (r *queryResolver) GetUserGamesHaveBadge(ctx context.Context, userid int) ([]*model.Game, error) {
	var game []*model.Game
	database.GetDB().Raw("select * from games where id in (select gameid from game_badges where gameid in (select gameid from user_games where userid = ?))", userid).Find(&game)

	return game, nil
}

func (r *queryResolver) GetAllUser(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	database.GetDB().Where("id != 1").Find(&users).Order("id")

	return users, nil
}

func (r *queryResolver) GetAllUserProfile(ctx context.Context) ([]*model.UserProfile, error) {
	var users []*model.UserProfile
	database.GetDB().Where("userid != 1").Find(&users).Order("userid")

	return users, nil
}

func (r *queryResolver) GetReportedUser(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	database.GetDB().Raw("select * from users where id in (select reporteduserid from report_users order by reporteduserid)").Find(&users)

	return users, nil
}

func (r *queryResolver) GetReportsUser(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	database.GetDB().Raw("select * from users where id in (select userwhoreportid from report_users order by reporteduserid)").Find(&users)

	return users, nil
}

func (r *queryResolver) GetUnsuspensionRequest(ctx context.Context) ([]*model.UnsuspensionRequest, error) {
	var reports []*model.UnsuspensionRequest
	database.GetDB().Find(&reports)

	return reports, nil
}

func (r *queryResolver) GetUnsuspensionUser(ctx context.Context) ([]*model.User, error) {
	var reports []*model.User
	database.GetDB().Raw("select u.* from users u join unsuspension_requests ur on u.id = ur.userid").Find(&reports).Order("id")

	return reports, nil
}

func (r *queryResolver) GetUserItem(ctx context.Context, userid int) ([]*model.GameItem, error) {
	var gameitems []*model.GameItem
	database.GetDB().Raw("select * from game_items where itemid in (select itemid from user_game_items where userid = ?)", userid).Find(&gameitems)

	return gameitems, nil
}

func (r *queryResolver) GetItemGameName(ctx context.Context, userid int) ([]*model.Game, error) {
	var gameitems []*model.Game
	database.GetDB().Raw("select * from games where id in (select distinct gameid from user_game_items where userid = ?)", userid).Find(&gameitems)

	return gameitems, nil
}

func (r *queryResolver) GetItemSalesTransaction(ctx context.Context, itemid int) ([]*model.Transaction, error) {
	var transaction []*model.Transaction
	database.GetDB().Where("itemid = ?", itemid).Find(&transaction)

	return transaction, nil
}

func (r *queryResolver) GetMarketItem(ctx context.Context) ([]*model.GameItem, error) {
	var gameitems []*model.GameItem
	database.GetDB().Raw("select gi.* from game_items gi left join transactions t on gi.itemid = t.itemid group by gi.gameid, gi.itemid, gi.itemn, gi.summary, gi.name order by count(t.itemid) desc").Find(&gameitems)

	return gameitems, nil
}

func (r *queryResolver) GetMarketDetailItem(ctx context.Context, itemid int) (*model.GameItem, error) {
	var detail model.GameItem
	database.GetDB().Where("itemid = ?", itemid).Find(&detail)

	return &detail, nil
}

func (r *queryResolver) SalesListing(ctx context.Context, itemid int, userid int) ([]*model.Market, error) {
	var sales []*model.Market
	database.GetDB().Raw("select * from markets where itemid = ? and type = 'Sell' and sellerid != ? group by markets.price,markets.sellerid, markets.buyerid, markets.itemid, markets.gameid, markets.type order by price desc", itemid, userid).Find(&sales)

	return sales, nil
}

func (r *queryResolver) BidListing(ctx context.Context, itemid int, userid int) ([]*model.Market, error) {
	var sales []*model.Market
	database.GetDB().Raw("select * from markets where itemid = ? and type = 'Bid' and sellerid != ? group by markets.buyerid ,markets.price,markets.sellerid, markets.itemid, markets.gameid, markets.type order by price asc", itemid, userid).Find(&sales)

	return sales, nil
}

func (r *queryResolver) UsersalesListing(ctx context.Context, itemid int, userid int) ([]*model.Market, error) {
	var sales []*model.Market
	database.GetDB().Raw("select * from markets where itemid = ? and type = 'Sell' and sellerid = ? group by markets.buyerid,markets.price,markets.sellerid, markets.itemid, markets.gameid, markets.type order by price desc", itemid, userid).Find(&sales)

	return sales, nil
}

func (r *queryResolver) UserbidListing(ctx context.Context, itemid int, userid int) ([]*model.Market, error) {
	var sales []*model.Market
	database.GetDB().Raw("select * from markets where itemid = ? and buyerid = ? and type = 'Bid' group by markets.buyerid,markets.price,markets.sellerid, markets.itemid, markets.gameid, markets.type order by price asc", itemid, userid).Find(&sales)

	return sales, nil
}

func (r *queryResolver) Itemsalesinamonth(ctx context.Context, itemid int) ([]*model.Transaction, error) {
	var sales []*model.Transaction
	database.GetDB().Where("itemid = ? and created_at BETWEEN (now() - '1 month'::interval) AND now()", itemid).Find(&sales)

	return sales, nil
}

func (r *queryResolver) NewreleasesGames(ctx context.Context) ([]*model.Game, error) {
	var games []*model.Game
	database.GetDB().Raw("select * from games order by created_at desc").Find(&games).Limit(10)

	return games, nil
}

func (r *queryResolver) GetUniqueCountry(ctx context.Context, gameid int) ([]*model.User, error) {
	var users []*model.User
	database.GetDB().Raw("select distinct u.country from users u join user_games ug on u.id = ug.userid where gameid = ? order by country", gameid).Find(&users)
	return users, nil
}

func (r *queryResolver) GetUserbyCountry(ctx context.Context, gameid int) ([]*model.User, error) {
	var users []*model.User
	database.GetDB().Raw("select u.* from users u join user_games ug on u.id = ug.userid where gameid = ? order by country", gameid).Find(&users)
	return users, nil
}

func (r *queryResolver) GetAllCountry(ctx context.Context) ([]*model.Country, error) {
	var countries []*model.Country
	database.GetDB().Find(&countries)

	return countries, nil
}

func (r *queryResolver) GetOtp(ctx context.Context) (string, error) {
	code := generateOTP()

	m := gomail.NewMessage()
	m.SetHeader("From", "store.staempowered@gmail.com")
	m.SetHeader("To", "clarissachuardi01@gmail.com")
	m.SetHeader("Subject", "OTP Code From Staem")
	m.SetBody("text", code)

	d := gomail.NewDialer("smtp.gmail.com", 587, "store.staempowered@gmail.com", "zwhrruqnoyhebgea")

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	return code, nil
}

func (r *queryResolver) PaginateAllGame(ctx context.Context, index int) ([]*model.Game, error) {
	var games []*model.Game

	database.GetDB().Scopes(Paginate10(index)).Find(&games)

	return games, nil
}

func (r *queryResolver) PaginateAllPromo(ctx context.Context, index int) ([]*model.Game, error) {
	var games []*model.Game

	database.GetDB().Scopes(Paginate10(index)).Raw("select * from games where id in (select gameid from game_promos)").Find(&games)

	return games, nil
}

func (r *queryResolver) PaginateAllUser(ctx context.Context, index int) ([]*model.User, error) {
	var users []*model.User
	database.GetDB().Scopes(Paginate10(index)).Find(&users)

	return users, nil
}

func (r *queryResolver) PaginateAllUserDetail(ctx context.Context, index int) ([]*model.UserProfile, error) {
	var users []*model.UserProfile
	database.GetDB().Scopes(Paginate10(index)).Find(&users)

	return users, nil
}

func (r *queryResolver) PaginateAllMarketItems(ctx context.Context, index int) ([]*model.GameItem, error) {
	var items []*model.GameItem
	database.GetDB().Scopes(Paginate10(index)).Raw("select * from game_items where itemid in (select itemid from markets)").Find(&items)

	return items, nil
}

func (r *queryResolver) PaginateUserInventory(ctx context.Context, userid int, index int) ([]*model.GameItem, error) {
	var gameitems []*model.GameItem
	database.GetDB().Scopes(Paginate10(index)).Raw("select * from game_items where itemid in (select itemid from user_game_items where userid = ?)", userid).Find(&gameitems)

	return gameitems, nil
}

func (r *queryResolver) PaginateMediaCommenter(ctx context.Context, gameid int, userid int, review string, index int) ([]*model.User, error) {
	var user []*model.User
	database.GetDB().Scopes(Paginate10(index)).Raw("select * from users where id in (select commenterid from user_comments where gameid = ? and posterid = ? and review = ?)", gameid, userid, review).Find(&user)

	return user, nil
}

func (r *queryResolver) PaginateMediaCommenterDetail(ctx context.Context, gameid int, userid int, review string, index int) ([]*model.UserProfile, error) {
	var user []*model.UserProfile
	database.GetDB().Scopes(Paginate10(index)).Raw("select * from user_profiles where userid in (select commenterid from user_comments where gameid = ? and posterid = ? and review = ?)", gameid, userid, review).Find(&user)

	return user, nil
}

func (r *queryResolver) GetAllUserActivities(ctx context.Context, userid *int) ([]*model.UserActivities, error) {
	var act []*model.UserActivities
	database.GetDB().Where("userid = ?", userid).Find(&act)

	return act, nil
}

func (r *queryResolver) PaginateUserActivities(ctx context.Context, userid int, index int) ([]*model.UserActivities, error) {
	var act []*model.UserActivities
	database.GetDB().Scopes(Paginate10(index)).Raw("select * from user_activities where userid = ? order by created_at desc", userid).Find(&act)

	return act, nil
}

func (r *queryResolver) Recommendedgames(ctx context.Context) ([]*model.Game, error) {
	var byid []*model.GameDetail
	database.GetDB().Select("id").Order("hoursplayed").Find(&byid).Limit(10)

	var games []*model.Game
	database.GetDB().Raw("SELECT * FROM games WHERE games.id in (SELECT id FROM game_details)").Find(&games)

	return games, nil
}

func (r *subscriptionResolver) MessageReceived(ctx context.Context, userid int) (<-chan string, error) {
	var event = make(chan string, 1)
	r.ChatSocket[userid] = event
	return event, nil
}

// Files returns generated.FilesResolver implementation.
func (r *Resolver) Files() generated.FilesResolver { return &filesResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type filesResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
func generateOTP() string {
	// omitted few confusing characters
	charSet := "ABCDEFGHJKLMNPQRSTUVWXYZ123456789"
	pass := randomStringGenerator(charSet, 5)
	return pass
}
func randomStringGenerator(charSet string, codeLength int32) string {
	rand.Seed(time.Now().UnixNano())
	code := ""
	charSetLength := int32(len(charSet))
	for i := int32(0); i < codeLength; i++ {
		index := rand.Intn(int(charSetLength))
		code += string(charSet[index])
	}

	return code
}
func randomNumber(min, max int32) int32 {
	rand.Seed(time.Now().UnixNano())
	return min + int32(rand.Intn(int(max-min)))
}
