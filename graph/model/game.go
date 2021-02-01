package model

//type Game struct {
//	ID        int64     `json:"id"`
//	Name      string    `json:"name"`
//	Price     int       `json:"price"`
//	Banner    []byte    `json:"banner"`
//	CreatedAt time.Time `json:"createdAt"`
//}
//
//type GameSlideshow struct {
//	Gameid int    `json:"gameid"`
//	Links  []byte `json:"links"`
//	contentType string
//}

type Files struct{
	Id		int
	ContentType	string
	File		[]byte
}
