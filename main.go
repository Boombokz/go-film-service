package main

import (
	"context"
	"filmservice/config"
	"filmservice/docs"
	"filmservice/handlers"
	"filmservice/logger"
	"filmservice/middlewares"
	"filmservice/repositories"
	"time"

	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
)

// @title 			FilmService API
// @version 		1.0
// @description 	This is a sample server
// @termsOfService 	http://swagger.io/terms/

// @contact.name 	API Support
// @contact.url 	http://www.swagger.io/support
// @contact.email 	support@swagger.io

// @license.name 	Apache 2.0
// @license.url 	http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host 		localhost:8081
// @BasePath 	/

// @securityDefinitions.apikey 	Bearer
// @in 							header
// @name 						Authorization
// @description 				Type "Bearer" followed by a space and JWT token.

// @externalDocs.description 	OpenAPI
// @externalDocs.url 			https://swagger.io/resources/open-api/
func main() {
	r := gin.New()
	gin.SetMode(gin.ReleaseMode)

	logger := logger2.GetLogger()
	r.Use(
		ginzap.Ginzap(logger, time.RFC3339, true),
		ginzap.RecoveryWithZap(logger, true),
	)

	corsConfig := cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"*"},
		AllowHeaders:    []string{"*"},
	}

	r.Use(cors.New(corsConfig))

	err := loadConfig()
	if err != nil {
		panic(err)
	}

	conn, err := connectToDb()
	if err != nil {
		panic(err)
	}

	moviesRepository := repositories.NewMoviesRepository(conn)
	genresRepository := repositories.NewGenresRepository(conn)
	watchListRepository := repositories.NewWatchListRepository(conn)
	usersRepository := repositories.NewUsersRepository(conn)

	moviesHandler := handlers.NewMoviesHandler(moviesRepository, genresRepository)
	genresHandler := handlers.NewGenreHandler(genresRepository)
	imageHandler := handlers.NewImageHandler()
	watchListHandler := handlers.NewWatchListHandlers(watchListRepository)
	usersHandler := handlers.NewUsersHandlers(usersRepository)
	authHandler := handlers.NewAuthHandlers(usersRepository)

	authorized := r.Group("")
	authorized.Use(middlewares.AuthMiddleware)

	authorized.GET("/movies", moviesHandler.FindAll)
	authorized.GET("/movies/:id", moviesHandler.FindById)
	authorized.POST("/movies", moviesHandler.Create)
	authorized.PUT("/movies/:id", moviesHandler.Update)
	authorized.DELETE("/movies/:id", moviesHandler.Delete)
	authorized.PATCH("/movies/:id/rate", moviesHandler.HandleSetRating)
	authorized.PATCH("/movies/:id/setWatched", moviesHandler.HandleSetWatched)

	authorized.GET("/genres", genresHandler.FindAll)
	authorized.GET("/genres/:id", genresHandler.FindById)
	authorized.POST("/genres", genresHandler.Create)
	authorized.PUT("/genres/:id", genresHandler.Update)
	authorized.DELETE("/genres/:id", genresHandler.Delete)

	authorized.GET("/watchlist", watchListHandler.GetAll)
	authorized.POST("/watchlist/:movieId", watchListHandler.Toggle)
	authorized.DELETE("/watchlist/:movieId", watchListHandler.Delete)

	authorized.GET("/users", usersHandler.FindAll)
	authorized.GET("/users/:id", usersHandler.FindById)
	authorized.POST("/users", usersHandler.Create)
	authorized.PUT("/users/:id", usersHandler.Update)
	authorized.PATCH("/users/:id/changePassword", usersHandler.ChangePassword)
	authorized.DELETE("/users/:id", usersHandler.Delete)
	authorized.GET("/users/userInfo", usersHandler.GetUserInfo)

	authorized.POST("/auth/signOut", authHandler.SignOut)

	unauthorized := r.Group("")

	unauthorized.POST("/auth/signIn", authHandler.SignIn)
	unauthorized.GET("/images/:imageId", imageHandler.HandleGetImageById)

	docs.SwaggerInfo.BasePath = ""
	unauthorized.GET("/swagger/*any", swagger.WrapHandler(swaggerfiles.Handler))

	logger.Info("Application starting...")

	err = r.Run(config.Config.AppHost)
	if err != nil {
		return
	}
}

func loadConfig() error {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")

	if err := viper.BindEnv("APP_HOST"); err != nil {
		viper.SetDefault("APP_HOST", ":8080")
	}
	if err := viper.BindEnv("DB_CONNECTION_STRING"); err != nil {
		viper.SetDefault("DB_CONNECTION_STRING", "postgres://postgres:postgres@localhost:5432/postgres")
	}
	if err := viper.BindEnv("JST_SECRET_KEY"); err != nil {
		viper.SetDefault("JST_SECRET_KEY", "supersecretkey")
	}
	if err := viper.BindEnv("JWT_EXPIRE_DURATION"); err != nil {
		viper.SetDefault("JWT_EXPIRE_DURATION", "24h")
	}

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	var mapConfig config.MapConfig
	err = viper.Unmarshal(&mapConfig)
	if err != nil {
		return err
	}

	config.Config = &mapConfig

	return nil
}

func connectToDb() (*pgxpool.Pool, error) {
	conn, err := pgxpool.New(context.Background(), config.Config.DbConnectionString)
	if err != nil {
		return nil, err
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return conn, nil
}
