package main

import (
	"log"
	"os"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/net/context"

	"time"

	"study_gin_golang/docs"
	"study_gin_golang/handlers"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Recipe struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"` //`swaggerignore:"true"`
	Name         string             `json:"name" bson:"name"`
	Tags         []string           `json:"tags" bson:"tags"`
	Ingredients  []string           `json:"ingredients" bson:"ingredients"`
	Instructions []string           `json:"instructions" bson:"instructions"`
	PublishedAt  time.Time          `json:"publishedAt" bson:"publishedAt"`
}

var recipes []Recipe
var ctx context.Context
var err error
var client *mongo.Client
var collection *mongo.Collection
var recipesHandler *handlers.RecipesHandler

func init() {

	ctx = context.Background()
	client, err = mongo.Connect(
		ctx,
		options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(),
		readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Conntected to MongoDB")
	collection = client.Database(os.Getenv(
		"MONGO_DATABASE")).Collection("recipes")
	recipesHandler = handlers.NewRecipesHandler(ctx, collection)

	/*
		var listOfRecipes []interface{}
		for _, recipe := range recipes {
			listOfRecipes = append(listOfRecipes, recipe)
		}
		collection := client.Database(os.Getenv(
			"MONGO_DATABASE")).Collection("recipes")
		insertManyResult, err := collection.InsertMany(
			ctx, listOfRecipes)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Inserted recipes: ",
			len(insertManyResult.InsertedIDs))
	*/

}

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth
func main() {

	// programmatically set swagger info
	docs.SwaggerInfo.Title = "Swagger For Study Gin"
	docs.SwaggerInfo.Description = "Building Distributed Applications in Gin_ A hands-on guide for Go developers to build and deploy distributed web apps with the Gin framework"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "singyu.swagger.io"
	docs.SwaggerInfo.BasePath = "/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	router := gin.Default()

	// api handlers v1
	v1 := router.Group("/api/v1")
	{
		recipes := v1.Group("/recipes")
		{
			recipes.POST("", recipesHandler.NewRecipeHandler)
			recipes.GET("", recipesHandler.ListRecipesHandler)
			recipes.PUT(":id", recipesHandler.UpdateRecipeHandler)
			recipes.DELETE(":id", recipesHandler.DeleteRecipeHandler)
			recipes.GET("search", recipesHandler.SearchRecipesHandler)
			recipes.GET(":id", recipesHandler.GetRecipeHandler)
		}
	}

	// swagger handlers
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run()
}
