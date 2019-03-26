package main

import (
	"sto_order_manager/config"
	"sto_order_manager/controllers"
	"sto_order_manager/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	defer models.DB_mysql.Close()

	router := gin.Default()

	cf := cors.DefaultConfig()
	cf.AllowOrigins = []string{"*"}
	cf.AddAllowHeaders("Origin")
	cf.AllowCredentials = true
	// cf.AllowAllOrigins = true
	router.Use(cors.New(cf))

	router.GET("/orderbook", controllers.GetOrderBook)
	router.POST("/fillorder", controllers.FillOrder)
	router.POST("/addorder", controllers.AddOrder)

	/*
		router.GET("/", func(c *gin.Context) {
			c.HTML(http.StatusOK, "./front/index.html", nil)
		})

		router.Static("/index", "./front/index.html")
	*/
	router.Static("/passportImages/", "./front/css")
	router.Static("/css/", "./front/css")
	router.Static("/js/", "./front/js")
	router.Static("/front/", "./front")

	router.StaticFile("/", "./front/index.html")

	port := ":" + config.V.Port
	router.Run(port)
}
