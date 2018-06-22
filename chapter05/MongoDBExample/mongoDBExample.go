package main

import (
	"context"
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var ctx = context.Background()

// Movie holds a movie data
type Movie struct {
	Name      string   `bson:"name"`
	Year      string   `bson:"year"`
	Directors []string `bson:"directors"`
	Writers   []string `bson:"writers"`
	BoxOffice `bson:"boxOffice"`
}

// BoxOffice is nested im Movie
type BoxOffice struct {
	Budget uint64 `bson:"budget"`
	Gross  uint64 `bson:"gross"`
}

func main() {
	// client, _ := mongo.NewClient("mongodb://127.0.0.1:27017")
	// err := client.Connect(ctx)
	// if err != nil {
	// 	log.Fatal("Connection failed with error: ", err)
	// }
	// defer client.Disconnect(ctx)
	// collection := client.Database("appdb").Collection("movies")

	// cur, err := collection.Find(ctx, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// for cur.Next(ctx) {
	// 	elem := bson.NewDocument()
	// 	err := cur.Decode(elem)
	// 	if err != nil {
	// 		log.Println(elem)
	// 	}
	// 	// log.Println(elem)
	// 	bsonElem, _ := elem.MarshalBSON()
	// 	jsonString, _ := extjson.BsonToExtJSON(false, bsonElem)
	// 	log.Println(jsonString)
	// }
	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		log.Fatal("Error connecting to DB server: ", err.Error())
	}
	defer session.Close()

	c := session.DB("appdb").C("movies")

	result := Movie{}
	err = c.Find(bson.M{"boxOffice.budget": bson.M{"$gt": 150000000}}).One(&result)
	if err != nil {
		log.Fatal("Error retrieving data from DB: ", err.Error())
	}

	fmt.Println("Movie:", result.Name)
}
