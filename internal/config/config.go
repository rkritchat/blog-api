package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Conf struct {
	DB         *mongo.Client
	AwsSession *session.Session
	Env        Env
}
type Env struct {
	Port             string `env:"PORT"`
	MongoUser        string `env:"MONGO_USER"`
	MongoPwd         string `env:"MONGO_PWD"`
	MongoURI         string `env:"MONGO_URI"`
	MongoDBName      string `env:"MONGO_DB_NAME"`
	MongoDBTimeout   int    `env:"MONGO_DB_TIMEOUT"`
	MongoMaxConIdle  int    `env:"MONGO_MAX_CON_IDEL"`
	MongoMaxPoolSize uint64 `env:"MONGO_MAX_POOL_SIZE"`
	MongoMinPoolSize uint64 `env:"MONGO_MIN_POOL_SIZE"`
	S3BucketName     string `env:"S3_BUCKET_NAME"`
	AwsRegion        string `env:"AWS_REGION"`
}

func init() {
	//load env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

func InitConfig() *Conf {
	var localEnv Env
	err := env.Parse(&localEnv)
	if err != nil {
		log.Fatal(err)
	}

	return &Conf{
		DB:         initConnection(localEnv),
		AwsSession: initAwsSession(localEnv),
		Env:        localEnv,
	}
}

func initConnection(env Env) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(env.MongoDBTimeout)*time.Second)
	defer cancel()

	URI := fmt.Sprintf("mongodb://%v:%v@%v", env.MongoUser, env.MongoPwd, env.MongoURI)
	conn, err := mongo.Connect(ctx, options.Client().ApplyURI(URI).SetMaxConnIdleTime(time.Duration(env.MongoMaxConIdle)*time.Minute).SetMaxPoolSize(env.MongoMaxPoolSize).SetMinPoolSize(env.MongoMinPoolSize))
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func initAwsSession(env Env) *session.Session {
	s, err := session.NewSession(
		&aws.Config{
			Region: aws.String(env.AwsRegion),
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	return s
}

func (c *Conf) Free() {
	if c.DB != nil {
		err := c.DB.Disconnect(context.Background())
		if err != nil {
			log.Fatal(err)
		}
	}
}
