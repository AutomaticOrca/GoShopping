package server

import "github.com/gin-gonic/gin"

func (server *Server) setupRouter() {
	router := gin.Default()
	router.GET("/_healthz", Healthz)

	apiRoutes := router.Group("/api")

	userRoutes := apiRoutes.Group("/users")
	{
		userRoutes.POST("/register", server.createUser)
		userRoutes.POST("/login", server.loginUser)
	}

	categoryRoutes := apiRoutes.Group("/categories")
	{
		categoryRoutes.POST("/", server.createCategory)
		categoryRoutes.GET("/:id", server.getCategory)
		categoryRoutes.GET("/", server.listCategories)
		categoryRoutes.GET("/parent/:parent_id", server.getSubcategories)
		categoryRoutes.PUT("/", server.updateCategory)
		categoryRoutes.DELETE("/:id", server.deleteCategory)
	}

	productRoutes := apiRoutes.Group("/products")
	{
		productRoutes.POST("/", server.createProduct)
		productRoutes.PUT("/:id", server.updateProduct)
		productRoutes.DELETE("/:id", server.deleteProduct)
		productRoutes.GET("/:id", server.getProduct)
		productRoutes.GET("/", server.listProducts)
	}

	server.router = router
}
