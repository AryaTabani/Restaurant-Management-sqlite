package main

import (
	"example.com/m/v2/controllers"
	"example.com/m/v2/DB"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()

	router := gin.Default()

	router.POST("/signup", controllers.SignUpHandler())
	router.POST("/login", controllers.LoginHandler())
	router.GET("/users", controllers.GetAllUsersHandler())
	router.GET("/users/:userid", controllers.GetUserByIDHandler())
	
	router.POST("/menus", controllers.CreateMenuHandler())
	router.GET("/menus", controllers.GetAllMenusHandler())
	router.GET("/menus/:menuid", controllers.GetMenuByIDHandler())
	router.PATCH("/menus/:menuid", controllers.UpdateMenuHandler())

	router.POST("/foods", controllers.CreateFoodHandler())
	router.GET("/foods", controllers.GetAllFoodsHandler())
	router.GET("/foods/:foodid", controllers.GetFoodByIDHandler())
	router.PATCH("/foods/:foodid", controllers.UpdateFoodHandler())
	
	router.GET("/tables", controllers.GetAllTablesHandler())
	router.GET("/tables/:tableid", controllers.GetTableByIDHandler())

	router.GET("/orders", controllers.GetAllOrdersHandler())
	router.GET("/orders/:orderid", controllers.GetOrderByIDHandler())
	router.PATCH("/orders/:orderid", controllers.UpdateOrderHandler())

	router.POST("/order-items", controllers.CreateOrderWithItemsHandler()) 
	router.GET("/order-items", controllers.GetAllOrderItemsHandler())
	router.GET("/order-items/:orderitemid", controllers.GetOrderItemHandler())
	router.PATCH("/order-items/:orderitemid", controllers.UpdateOrderItemHandler())
	
	router.POST("/invoices", controllers.CreateInvoiceHandler())
	router.GET("/invoices", controllers.GetAllInvoicesHandler())
	router.GET("/invoices/:invoiceid", controllers.GetInvoiceHandler())
	router.PATCH("/invoices/:invoiceid", controllers.UpdateInvoiceHandler())

	router.Run(":8080")
}