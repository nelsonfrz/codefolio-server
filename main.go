package main

import (
	"codefolio/server/app"
	"fmt"
	"log"
	"os"

	_ "codefolio/server/docs"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @SecurityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @title Echo Swagger Example API
// @version 1.0
// @description Codefolio backend server written in Go.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
// @schemes http
func main() {
	fmt.Println(os.Getenv("DB_URL"))
	fmt.Println(os.Getenv("JWT_SECRET"))

	err := godotenv.Load()

	e := echo.New()
	e.Validator = &app.CustomValidator{Validator: validator.New()}

	a, err := app.NewApp(os.Getenv("DB_URL"))
	if err != nil {
		log.Fatal(err)
	}

	jwtConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(app.JwtCustomClaims)
		},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}

	if os.Getenv("ENVIRONMENT") == "dev" {
		e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
			LogStatus: true,
			LogURI:    true,
			LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
				fmt.Printf("REQUEST: uri: %v, status: %v\n", v.URI, v.Status)
				return nil
			},
		}))
	}

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/", a.HealthRoute)
	e.POST("/login", a.Login)
	e.POST("/register", a.Register)

	e.GET("/users", a.GetAllUsers)
	e.GET("/projects", a.GetAllProjects)

	auth := e.Group("/authorized")
	auth.Use(echojwt.WithConfig(jwtConfig))
	auth.GET("/projects", a.GetAllProjectsAuthorized)
	auth.POST("/projects", a.CreateProject)
	auth.DELETE("/projects", a.DeleteProject)

	e.Logger.Fatal(e.Start(":3000"))
}
