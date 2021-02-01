package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/clarissac01/Staem/database"
	"github.com/clarissac01/Staem/graph/function"
	"github.com/clarissac01/Staem/graph/generated"
	"github.com/clarissac01/Staem/graph/model"
	jwt "github.com/dgrijalva/jwt-go"
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

	return token.SignedString([]byte("mimi"))
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

	return "", nil
}

func (r *mutationResolver) Updatepromo(ctx context.Context, id int, discount int, validTo time.Time) (string, error) {
	database.GetDB().Model(model.GamePromo{}).Where("gameid = ?", id).Update("discount", discount)
	database.GetDB().Model(model.GamePromo{}).Where("gameid = ?", id).Update("validTo", validTo)

	return "", nil
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

func (r *queryResolver) Auth(ctx context.Context, token string) (*model.User, error) {
	userid, err := function.ParseToken(token)
	var user model.User
	database.GetDB().First(&user, int(userid))

	return &user, err
}

func (r *queryResolver) GetUserProfile(ctx context.Context, id int) (*model.UserProfile, error) {
	userprofile := model.UserProfile{Userid: id}
	database.GetDB().Find(&userprofile)

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
	var byid []*model.GamePromo
	database.GetDB().Select("gameid").Order("discount").Find(&byid)
	var gameids []int
	for i := 0; i < len(byid); i++ {
		gameids[i] = int(byid[i].Gameid)
	}

	var games []*model.Game
	if len(gameids) == 0 {
		return nil, nil

	}
	database.GetDB().Find(&games).Where("id = ?", gameids)

	return games, nil
}

func (r *queryResolver) SearchGame(ctx context.Context, keyword string) ([]*model.Game, error) {
	var games []*model.Game
	database.GetDB().Where("name LIKE ? OR name LIKE ?", "%"+strings.ToLower(keyword)+"%", "%"+strings.ToUpper(keyword)+"%").Find(&games).Limit(5)

	return games, nil
}

// Files returns generated.FilesResolver implementation.
func (r *Resolver) Files() generated.FilesResolver { return &filesResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type filesResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
