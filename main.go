package main

import (
	"fmt"
	"log"
	"os"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/net/context"

	"study_gin_golang/docs"
	"study_gin_golang/handlers"

	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/go-redis/redis/v8"
)

var authHandler *handlers.AuthHandler
var recipesHandler *handlers.RecipesHandler

func init() {

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

	ctx := context.Background()
	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err = client.Ping(context.TODO(),
		readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Conntected to MongoDB")

	collection := client.Database(os.Getenv(
		"MONGO_DATABASE")).Collection("recipes")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	status := redisClient.Ping(ctx)
	fmt.Println("status: ", status)
	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)

	collectionUsers := client.Database(os.Getenv(
		"MONGO_DATABASE")).Collection("users")
	authHandler = handlers.NewAuthHandler(ctx, collectionUsers)
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

	store, _ := redisStore.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	router.Use(sessions.Sessions("recipes_api", store))
	// api handlers v1
	v1 := router.Group("/api/v1")
	{
		v1.POST("/signin", authHandler.SignInHandler)
		v1.POST("/signout", authHandler.SignOutHandler)
		v1.POST("/refresh", authHandler.RefreshHandler)

		recipes := v1.Group("/recipes")
		{
			recipes.GET("", recipesHandler.ListRecipesHandler)
		}
		authorized_recipes := v1.Group("/recipes")
		authorized_recipes.Use(authHandler.AuthMiddleware())
		{
			authorized_recipes.POST("", recipesHandler.NewRecipeHandler)
			authorized_recipes.PUT(":id", recipesHandler.UpdateRecipeHandler)
			authorized_recipes.DELETE(":id", recipesHandler.DeleteRecipeHandler)
			authorized_recipes.GET("search", recipesHandler.SearchRecipesHandler)
			authorized_recipes.GET(":id", recipesHandler.GetRecipeHandler)
		}
	}

	// swagger handlers
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run()
}
