package handler

import (
	"context"
	"net/http"

	"github.com/Fog3211/learn-golang-mongodb-api/config"
	"github.com/Fog3211/learn-golang-mongodb-api/controllers"
	"github.com/Fog3211/learn-golang-mongodb-api/routes"
	"github.com/Fog3211/learn-golang-mongodb-api/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoclient *mongo.Client
	redisclient *redis.Client

	userService         services.UserService
	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	authCollection      *mongo.Collection
	authService         services.AuthService
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController
)

func init() {
	config, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	ctx = context.TODO()

	// Connect to MongoDB
	mongoconn := options.Client().ApplyURI(config.DBUri)
	mongoclient, err = mongo.Connect(ctx, mongoconn)
	if err != nil {
		panic(err)
	}

	if err := mongoclient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	// Connect to Redis
	redisclient = redis.NewClient(&redis.Options{
		Addr: config.RedisUri,
	})

	// Collections
	authCollection = mongoclient.Database("golang_mongodb").Collection("users")
	userService = services.NewUserServiceImpl(authCollection, ctx)
	authService = services.NewAuthService(authCollection, ctx)
	AuthController = controllers.NewAuthController(authService, userService)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(userService)
	UserRouteController = routes.NewRouteUserController(UserController)

	server = gin.Default()
}

// Handler - Vercel serverless function handler
func Handler(w http.ResponseWriter, r *http.Request) {
	setupServer()
	server.ServeHTTP(w, r)
}

func setupServer() {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Server is healthy"})
	})

	AuthRouteController.AuthRoute(router, userService)
	UserRouteController.UserRoute(router, userService)
}
