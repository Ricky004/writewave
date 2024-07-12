package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBConfig struct {
	User string
	Pwd  string
	Host string
	Port uint16
	Name string
}

type Document[T any] interface {
	EnsureIndexes(Database)
	GetValue() *T
	Validate() error
}

type Database interface {
	GetInstance() *database
	Connect()
	Disconnect() 
}

type database struct {
	*mongo.Database
	context context.Context
	config  DBConfig
	client  *mongo.Client
}

func (db *database) Connect() {
	uri := fmt.Sprintf(
		"mongodb+srv://%s:%s@%s/%s",
		db.config.User, db.config.Pwd, db.config.Host, db.config.Name,
	)

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(db.context, clientOptions)
	if err != nil {
		log.Fatal("Connecting to mongoDB failed!!!", err)
	}

	err = client.Ping(db.context, nil)
	if err != nil {
		log.Panic("pinging to mongo failed!: ", err)
	}
	fmt.Println("connected to mongo!")

	db.Database = client.Database(db.config.Name)
}

func (db *database) Disconnect() {
	if db.client != nil {
		if err := db.client.Disconnect(db.context); err != nil {
			log.Fatal("Failed to disconnect from mongoDB", err)
		}
	}
	fmt.Println("disconnect mongodb")
}

func NewDatabase(ctx context.Context, config DBConfig) Database {
	db := database{
		context: ctx,
		config:  config,
	}
	return &db
}

func (db *database) GetInstance() *database {
	return db
}
