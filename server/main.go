package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"os"
	"path/filepath"
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

	//fileDB := client.Database("file_sharing")

	//var result bson.M
	//err = coll.FindOne(context.TODO(), bson.D{{"name", name}}).Decode(&result)
	//if err != nil {
	//	if err == mongo.ErrNoDocuments {
	//		return
	//	}
	//
	//	panic(err)
	//}

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "fuck ya chickin strips",
		})
	})
	r.POST("/upload", func(c *gin.Context) {
		uploadFile(c)
	})
	r.POST("/save/:file", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "you little turd",
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

func uploadFile(c *gin.Context) {
	//Retrieve file
	file, err := c.FormFile("file")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file was recieved",
		})
		return
	}
	//Generate random name to avoid conflicts
	ext := filepath.Ext(file.Filename)
	randomName := uuid.New().String() + ext
	if err := c.SaveUploadedFile(file, "/save/"+randomName); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "Something went wrong uploading the file",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
	})
}

func getEnvVar(key string) string {
	err := godotenv.Load()
	if err != nil {
		panic("Unable to load environment variable %d")
	}

	return os.Getenv(key)
}
