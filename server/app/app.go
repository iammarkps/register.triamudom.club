package app

import (
	"encoding/gob"
	"time"

	//Redis
	_ "github.com/gomodule/redigo/redis"
	"github.com/gorilla/sessions"
	"github.com/iammarkps/clubreg/server/handler"

	"github.com/iammarkps/clubreg/server/middlewares"
	"github.com/iammarkps/clubreg/server/models"
	"github.com/jinzhu/gorm"

	//MariaDB dialects
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/boj/redistore.v1"

	"log"
	"net/http"
)

//New return new echo and gorm object
func New() (*echo.Echo, *gorm.DB) {
	gob.Register(time.Time{})
	e := echo.New()

	dbURI := "root:password@tcp(127.0.0.1:3306)/crn?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open("mysql", dbURI)
	log.Printf("Connecting to database...")
	if err != nil {
		log.Fatal("Could not connect database")
	} else {
		log.Printf("Successfully connected to database")
	}
	// defer db.Close()

	db.BlockGlobalUpdate(true)
	db.DB().SetMaxIdleConns(100)
	db.AutoMigrate(&models.Club{}, &models.User{}, &models.Audition{}, &models.Settings{})

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowOrigins:     []string{"http://localhost:3000"},
	}))
	e.Use(middleware.Gzip())
	e.Use(middleware.Secure())
	e.Use(session.Middleware(newRedisStore()))

	e.Debug = true

	h := &handler.Handler{DB: db}
	m := &middlewares.Middleware{DB: db}

	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "😼Triam Udom Suksa School's club registration API is running!")
	})
	e.POST("/login", h.Login)
	e.POST("/register", h.Register)
	e.GET("/logout", h.Logout)
	e.GET("/health", h.HealthCheck)
	e.GET("/user", h.User, m.Auth(1))

	return e, db
}

func newRedisStore() sessions.Store {
	store, err := redistore.NewRediStore(10, "tcp", ":6379", "", []byte("secret"))
	log.Printf("Connecting to redis...")
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Successfully connected to redis")
	}
	return store
}
