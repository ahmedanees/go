package models

import "time"

// define User model
type User struct{
	Id int16 `json:"id" bson:"id"`
	Email string `json:"email" bson:"email"`
	Password string `json:"password" bson:"password"`
	Created_On    time.Time          `json:"created_at"`
	Updated_On    time.Time          `json:"updated_at"`
}
type Tweet struct {
	Id int16 `json:"id" bson:"id"`
	UserId string `json:"email" bson:"user_id"`
	Description string     `json:"description"`
	Created_On    time.Time          `json:"created_at"`
	Updated_On    time.Time          `json:"updated_at"`
}


type TweetData struct {
	Description string `json:"description"`
	Hashtag string     `json:"hashtag"`
	
}
type TweetInfo struct {
	Tweet string     `json:"tweet"`
	CreatedDate string     `json:"created_date"`
}