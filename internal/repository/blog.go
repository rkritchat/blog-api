package repository

import (
	"blog-api/internal/config"
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Content struct {
	Message  string
	Filename string
}

type BlogEntity struct {
	Title    string
	Contents []Content
}

type Blog interface {
	Create(entity *BlogEntity) error
	FindByTitle(title string) (*BlogEntity, error)
}

type blog struct {
	dbName         string
	collectionName string
	db             *mongo.Client
}

func NewBlog(db *mongo.Client, env config.Env) Blog {
	return &blog{
		dbName:         env.MongoDBName,
		collectionName: "blog",
		db:             db,
	}
}

func (repo blog) Create(entity *BlogEntity) error {
	r, err := repo.db.Database(repo.dbName).Collection(repo.collectionName).InsertOne(context.TODO(), entity)
	if err != nil {
		return err
	}
	if r.InsertedID == nil {
		return errors.New("insert failed")
	}
	return nil
}

func (repo blog) FindByTitle(title string) (*BlogEntity, error) {
	on := bson.M{"title": title}
	var r BlogEntity
	err := repo.db.Database(repo.dbName).Collection(repo.collectionName).FindOne(context.TODO(), on).Decode(&r)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &r, nil
}
