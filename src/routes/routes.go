package routes

import (
	"fmt"

	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/myrachanto/microservice/user/src/controllers"
	m "github.com/myrachanto/microservice/user/src/middlewares"
	"github.com/myrachanto/microservice/user/src/repository"
	"github.com/myrachanto/microservice/user/src/service"

	"github.com/spf13/viper"
)

//StoreAPI =>entry point to routes
type Open struct {
	Port     string `mapstructure:"PORT"`
	Key      string `mapstructure:"EncryptionKey"`
	DURATION string `mapstructure:"DURATION"`
}

// func LoadConfig(path string) (open Open, err error) {
func LoadConfig() (open Open, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&open)
	return
}
func StoreApi() {
	open, err := LoadConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	repository.IndexRepo.InitDB()
	//check db connection//////////////////////
	fmt.Println("initialization----------------")
	controllers.NewUserController(service.NewUserService(repository.NewUserRepo()))
	e := echo.New()

	e.Static("/", "public")
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	JWTgroup := e.Group("/api/")
	JWTgroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS256",
		SigningKey:    []byte(open.Key),
	}))
	front := e.Group("/front")
	front.POST("/register", controllers.UserController.Create)
	front.POST("/login", controllers.UserController.Login, m.Tokenizing)
	front.GET("/users", controllers.UserController.GetAll, m.Tokenizing)
	JWTgroup.GET("logout/:token", controllers.UserController.Logout)
	JWTgroup.PUT("users/:id", controllers.UserController.Update, m.Tokenizing)
	JWTgroup.DELETE("users/:id", controllers.UserController.Delete, m.Tokenizing)
	//e.DELETE("loggoutall/:id", controllers.UserController.DeleteALL) logout all accounts

	// Start server
	e.Logger.Fatal(e.Start(open.Port))
}
