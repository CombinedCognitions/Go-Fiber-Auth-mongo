package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                 primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Email              string             `json:"email" bson:"email,omitempty"`
	Password           string             `json:"password" bson:"password,omitempty"`
	CreatedAt          time.Time          `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt          time.Time          `json:"updatedAt" bson:"updatedAt,omitempty"`
	Username           string             `json:"username" bson:"username,omitempty"`
	Bio                string             `json:"bio" bson:"bio,omitempty"`
	Age                int64              `json:"age" bson:"age,omitempty"`
	Phone              string             `json:"phone" bson:"phone,omitempty"`
	Gender             string             `json:"gender" bson:"gender,omitempty"`
	Profilepicturelink string             `json:"profilepicturelink" bson:"profilepicturelink,omitempty"`
	Uploadsno          int64              `json:"uploadsno" bson:"uploadsno,omitempty"`
}
