package platform

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PlatformLogs struct {
	MongoDb    string `json:"mongoDb"`    //
	Collection string `json:"collection"` //
}

type HistoryLogs struct {
	ID        primitive.ObjectID `bson:"_id" json:"ID"`                //
	Userid    int                `bson:"userid" json:"Userid"`         //
	Country   string             `bson:"country" json:"Country"`       //
	Message   string             `bson:"message" json:"Message"`       //
	IPAddress string             `bson:"IPAddress" json:"IPAddress"`   //
	CreatedAt time.Time          `bson:"created_at" json:"created_at"` //
}

type ServiceLogs struct {
}

var collection *mongo.Collection
var ctx = context.TODO()

func (a PlatformLogs) connection() {
	clientOptions := options.Client().ApplyURI(a.MongoDb)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	collection = client.Database("platform").Collection(a.Collection)
}

func (a PlatformLogs) InsertRecord(history *HistoryLogs) error {
	a.connection()
	history.ID = primitive.NewObjectID()
	history.CreatedAt = time.Now()
	_, err := collection.InsertOne(ctx, history)
	return err
}

func (a PlatformLogs) ListRecords(filters interface{}) ([]*HistoryLogs, error) {
	a.connection()
	var history []*HistoryLogs
	filter := bson.D{{}}

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		return history, err
	}

	for cur.Next(ctx) {
		var t HistoryLogs
		err := cur.Decode(&t)
		if err != nil {
			return history, err
		}
		history = append(history, &t)
	}

	if err := cur.Err(); err != nil {
		return history, err
	}

	cur.Close(ctx)

	if len(history) == 0 {
		return history, mongo.ErrNoDocuments
	}

	return history, nil

}
