package model

import (
	"fmt"
	"github.com/memochou1993/gh-rankings/app/resource"
	"github.com/memochou1993/gh-rankings/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type User struct {
	AvatarURL    string       `json:"avatarUrl,omitempty" bson:"avatar_url,omitempty"`
	CreatedAt    *time.Time   `json:"createdAt,omitempty" bson:"created_at,omitempty"`
	Followers    *Items       `json:"followers,omitempty" bson:"followers,omitempty"`
	Location     string       `json:"location,omitempty" bson:"location,omitempty"`
	Login        string       `json:"login" bson:"_id"`
	Name         string       `json:"name,omitempty" bson:"name,omitempty"`
	Gists        []Gist       `json:"gists,omitempty" bson:"gists,omitempty"`
	Repositories []Repository `json:"repositories,omitempty" bson:"repositories,omitempty"`
	Tags         []string     `json:"tags,omitempty" bson:"tags,omitempty"`
}

func (u *User) ID() string {
	return u.Login
}

func (u *User) AddTypeTag() {
	u.Tags = append(u.Tags, fmt.Sprintf("type:%s", TypeUser))
}

func (u *User) AddLocationTag() {
	for _, location := range resource.Locate(u.Location) {
		u.Tags = append(u.Tags, fmt.Sprintf("location:%s", location))
	}
}

type UserModel struct {
	*Model
}

func (u *UserModel) Store(users []User) *mongo.BulkWriteResult {
	if len(users) == 0 {
		return nil
	}
	var models []mongo.WriteModel
	for _, user := range users {
		user.AddTypeTag()
		user.AddLocationTag()
		filter := bson.D{{"_id", user.ID()}}
		update := bson.D{{"$set", user}}
		models = append(models, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true))
	}
	return database.BulkWrite(u.Name(), models)
}

func (u *UserModel) UpdateGists(user User, gists []Gist) {
	filter := bson.D{{"_id", user.ID()}}
	update := bson.D{{"$set", bson.D{{"gists", gists}}}}
	database.UpdateOne(u.Name(), filter, update)
}

func (u *UserModel) UpdateRepositories(user User, repositories []Repository) {
	filter := bson.D{{"_id", user.ID()}}
	update := bson.D{{"$set", bson.D{{"repositories", repositories}}}}
	database.UpdateOne(u.Name(), filter, update)
}

func NewUserModel() *UserModel {
	return &UserModel{
		&Model{
			name: "users",
		},
	}
}