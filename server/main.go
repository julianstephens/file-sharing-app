package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"os"
)

func main() {
	client, err := ConnectDB()
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("sample_airbnb").Collection("listingsAndReviews")
	name := "Ribeira Charming Duplex"

	var result bson.M
	err = coll.FindOne(context.TODO(), bson.D{{"name", name}}).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return
		}

		panic(err)
	}
	//fileDB := client.Database("file-sharing")
	//filesColl := fileDB.Collection("files")

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"data": result,
		})
	})

	r.Run()
}

func ConnectDB() (*mongo.Client, error) {
	uri := getEnvVar("MONGODB_URI")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}

	return client, nil
}

func getEnvVar(key string) string {
	err := godotenv.Load()
	if err != nil {
		panic("Unable to load environment variable %d")
	}

	return os.Getenv(key)
}
