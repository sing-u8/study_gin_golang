package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/net/context"

	"time"

	"study_gin_golang/docs"

	"github.com/gin-gonic/gin"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Recipe struct {
	ID           primitive.ObjectID `json:"id" bson:"_id"` //`swaggerignore:"true"`
	Name         string             `json:"name" bson: "name"`
	Tags         []string           `json:"tags" "tags"`
	Ingredients  []string           `json:"ingredients" bson:"ingredients"`
	Instructions []string           `json:"instructions" bson: "instructions"`
	PublishedAt  time.Time          `json:"publishedAt" bson:"publishedAt"`
}

// ListRecipe godoc
// @Summary      recipes listRecipes
// @Description  Returns list of recipes
// @Tags         recipe
// @Accept       json
// @Produce      application/json
// @Success      200  {array}  Recipe  "Successful operation"
// @Router       /recipes [get]
func ListRecipesHandler(c *gin.Context) {
	cur, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(ctx)

	recipes := make([]Recipe, 0)
	for cur.Next(ctx) {
		var recipe Recipe
		cur.Decode(&recipe)
		recipes = append(recipes, recipe)
	}
	c.JSON(http.StatusOK, recipes)
}

// NewRecipe godoc
// @Summary      recipes newRecipe
// @Description  create an rew recipe
// @Tags         recipe
// @Accept       json
// @Produce      application/json
// @Param		 recipe body Recipe true "Recipe Schema"
// @Success      200  {object}  Recipe  "Successful operation"
// @Failure      400  {string}  string	"Invalid input"
// @Router       /recipes [post]
func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	_, err = collection.InsertOne(ctx, recipe)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new recipe"})
		return
	}
	c.JSON(http.StatusOK, recipe)
}

// UpdateRecipe godoc
// @Summary      recipes updateRecipe
// @Description  Update an existing recipe
// @Tags         recipe
// @Accept       json
// @Produce      application/json
// @Param		 id path string true "ID of the recipe"
// @Param		 recipe body Recipe true "Recipe Schema"
// @Success      200  {object}  Recipe  "Successful operation"
// @Failure      400  {string}  string	"Invalid input"
// @Failure      404  {string}  string	"Invalid recipe ID"
// @Router       /recipes/{id} [put]
func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err = collection.UpdateOne(ctx, bson.M{
		"_id": objectId,
	}, bson.D{{Key: "$set", Value: bson.D{
		{Key: "name", Value: recipe.Name},
		{Key: "instructions", Value: recipe.Instructions},
		{Key: "ingredients", Value: recipe.Ingredients},
		{Key: "tags", Value: recipe.Tags},
	}}})

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been updated"})
}

// DeleteRecipe godoc
// @Summary      recipes deleteRecipe
// @Description  Delete an existing recipe
// @Tags         recipe
// @Accept       json
// @Produce      application/json
// @Param		 id path string true "ID of the recipe"
// @Success      200  {object}  Recipe  "Successful operation"
// @Failure      404  {string}  string	"Invalid recipe ID"
// @Router       /recipes/{id} [delete]
func DeleteRecipeHandler(c *gin.Context) {
	// id := c.Param("id")
	// index := -1
	// for i := 0; i < len(recipes); i++ {
	// 	if recipes[i].ID == id {
	// 		index = i
	// 		break
	// 	}
	// }
	// if index == -1 {
	// 	c.JSON(http.StatusNotFound, gin.H{
	// 		"error": "Recipe not found"})
	// 	return
	// }
	// recipes = append(recipes[:index], recipes[index+1:]...)
	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "Recipe has been deleted"})
}

// SearchRecipes godoc
// @Summary      recipes findRecipe
// @Description  Search recipes based on tags
// @Tags         recipe
// @Accept       json
// @Produce      application/json
// @Param		 tag query string true "recipe tag"
// @Success      200 {object} Recipe "Successful operation"
// @Router       /recipes/search [get]
func SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")
	listOfRecipes := make([]Recipe, 0)
	for i := 0; i < len(recipes); i++ {
		found := false
		for _, t := range recipes[i].Tags {
			if strings.EqualFold(t, tag) {
				found = true
			}
		}
		if found {
			listOfRecipes = append(listOfRecipes,
				recipes[i])
		}
	}
	c.JSON(http.StatusOK, listOfRecipes)
}

// GetRecipe godoc
// @Summary      recipes findRecipe
// @Description  Search recipes based on tags
// @Tags         recipe
// @Accept       json
// @Produce      application/json
// @Param		 id path string true "ID of recipe"
// @Success      200 {object} Recipe "Successful operation"
// @Failure      404  {string}  string	"Invalid recipe ID"
// @Router       /recipes/search [get]
func GetRecipeHandler(c *gin.Context) {
	// id := c.Param("id")
	// for i := 0; i < len(recipes); i++ {
	// 	if recipes[i].ID == id {
	// 		c.JSON(http.StatusOK, recipes[i])
	// 		return
	// 	}
	// }

	// c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
}

var recipes []Recipe

var ctx context.Context
var err error
var client *mongo.Client

var collection *mongo.Collection

func init() {

	/*
		recipes = make([]Recipe, 0)
		file, _ := ioutil.ReadFile("recipes.json")
		_ = json.Unmarshal([]byte(file), &recipes)
	*/

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
			recipes.POST("", NewRecipeHandler)
			recipes.GET("", ListRecipesHandler)
			recipes.PUT(":id", UpdateRecipeHandler)
			recipes.DELETE(":id", DeleteRecipeHandler)
			recipes.GET("search", SearchRecipesHandler)
			recipes.GET(":id", GetRecipeHandler)
		}
	}

	// swagger handlers
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run()
}
