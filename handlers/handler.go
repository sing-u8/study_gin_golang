package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"study_gin_golang/models"

	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
)

type RecipesHandler struct {
	collection  *mongo.Collection
	ctx         context.Context
	redisClient *redis.Client
}

func NewRecipesHandler(ctx context.Context, collection *mongo.Collection, redisClient *redis.Client) *RecipesHandler {
	return &RecipesHandler{
		collection:  collection,
		ctx:         ctx,
		redisClient: redisClient,
	}
}

// ListRecipe godoc
// @Summary      recipes listRecipes
// @Description  Returns list of recipes
// @Tags         recipe
// @Accept       json
// @Produce      application/json
// @Success      200  {array}  Recipe  "Successful operation"
// @Router       /recipes [get]
func (handler *RecipesHandler) ListRecipesHandler(c *gin.Context) {
	val, err := handler.redisClient.Get(handler.ctx, "recipes").Result()
	if err == redis.Nil {
		log.Printf("Request to MongoDB")
		cur, err := handler.collection.Find(handler.ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer cur.Close(handler.ctx)

		recipes := make([]models.Recipe, 0)
		for cur.Next(handler.ctx) {
			var recipe models.Recipe
			cur.Decode(&recipe)
			recipes = append(recipes, recipe)
		}

		data, _ := json.Marshal(recipes)
		handler.redisClient.Set(handler.ctx, "recipes", string(data), 0)
		c.JSON(http.StatusOK, recipes)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		log.Printf("Request to Redis")
		recipes := make([]models.Recipe, 0)
		json.Unmarshal([]byte(val), &recipes)
		c.JSON(http.StatusOK, recipes)
	}

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
func (handler *RecipesHandler) NewRecipeHandler(c *gin.Context) {
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()
	_, err := handler.collection.InsertOne(handler.ctx, recipe)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting a new recipe"})
		return
	}

	log.Println("Remove data from Redis")
	handler.redisClient.Del(handler.ctx, "recipes")

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
func (handler *RecipesHandler) UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := handler.collection.UpdateOne(handler.ctx, bson.M{
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

	handler.redisClient.Del(handler.ctx, "recipes")

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
func (handler *RecipesHandler) DeleteRecipeHandler(c *gin.Context) {
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
func (handler *RecipesHandler) SearchRecipesHandler(c *gin.Context) {
	// tag := c.Query("tag")
	// listOfRecipes := make([]models.Recipe, 0)
	// for i := 0; i < len(recipes); i++ {
	// 	found := false
	// 	for _, t := range recipes[i].Tags {
	// 		if strings.EqualFold(t, tag) {
	// 			found = true
	// 		}
	// 	}
	// 	if found {
	// 		listOfRecipes = append(listOfRecipes,
	// 			recipes[i])
	// 	}
	// }
	// c.JSON(http.StatusOK, listOfRecipes)
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
func (handler *RecipesHandler) GetRecipeHandler(c *gin.Context) {
	// id := c.Param("id")
	// for i := 0; i < len(recipes); i++ {
	// 	if recipes[i].ID == id {
	// 		c.JSON(http.StatusOK, recipes[i])
	// 		return
	// 	}
	// }

	// c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
}
