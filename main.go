package main

import (
	"fmt"
	"net/http"

	"github.com/62teknologi/62whale/62golib/utils"
	"github.com/62teknologi/62whale/app/http/controllers"
	"github.com/62teknologi/62whale/app/http/middlewares"
	"github.com/62teknologi/62whale/app/interfaces"
	"github.com/62teknologi/62whale/config"

	"github.com/gin-gonic/gin"
)

func main() {

	configs, err := config.LoadConfig(".", &config.Data)
	if err != nil {
		fmt.Printf("cannot load config: %w", err)
		return
	}

	// todo : replace last variable with spread notation "..."
	utils.ConnectDatabase(configs.DBDriver, configs.DBSource1, configs.DBSource2)

	utils.InitPluralize()

	r := gin.Default()

	apiV1 := r.Group("/api/v1").Use(middlewares.DbSelectorMiddleware())
	{
		RegisterRoute(apiV1, "comment", controllers.CommentController{})
		RegisterRoute(apiV1, "category", controllers.CategoryController{})
		RegisterRoute(apiV1, "catalog", controllers.CatalogController{})
		RegisterRoute(apiV1, "group", controllers.GroupController{})
		RegisterRoute(apiV1, "item", controllers.ItemController{})
		RegisterRoute(apiV1, "review", controllers.ReviewController{})
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, utils.ResponseData("success", "Server running well", nil))
	})

	err = r.Run(configs.HTTPServerAddress)

	if err != nil {
		fmt.Printf("cannot run server: %w", err)
		return
	}
}

func RegisterRoute(r gin.IRoutes, t string, c interfaces.Crud) {
	r.GET("/"+t+"/:table/:id", c.Find)
	r.GET("/"+t+"/:table/slug/:slug", c.Find)
	r.GET("/"+t+"/:table", c.FindAll)
	r.POST("/"+t+"/:table", c.Create)
	r.PUT("/"+t+"/:table/:id", c.Update)
	r.DELETE("/"+t+"/:table/:id", c.Delete)
}
