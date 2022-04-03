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
	ID        primitive.ObjectID `bson:"_id" json:"ID"`                        //
	Userid    int                `bson:"userid" json:"userid"`                 //
	Serviceid int                `bson:"serviceid" json:"serviceid,omitempty"` //
	Country   string             `bson:"country" json:"country"`               //
	Message   string             `bson:"message" json:"type"`                  //
	IPAddress string             `bson:"ipaddress" json:"ipaddress"`           //
	CreatedAt time.Time          `bson:"created_at" json:"date"`               //
	Tag       string             `bson:"tag" json:"tag"`                       //
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
	//filter := bson.D{{}}

	opts := options.Find().SetSort(bson.D{{"created_at", -1}})

	cur, err := collection.Find(ctx, filters, opts)
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
