package app

import (
	"codefolio/server/database"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type JwtCustomClaims struct {
	UserId int32 `json:"userId"`
	Admin  bool  `json:"admin"`
	jwt.RegisteredClaims
}

type Login struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// @Summary User login
// @Description Authenticate user and provide JWT token
// @Tags authentication
// @Accept json
// @Produce json
// @Param credentials body Login true "User login credentials"
// @Success 200 {object} map[string]string "Successfully authenticated user with JWT token"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /login [post]
func (a *App) Login(c echo.Context) error {
	login := new(Login)
	if err := c.Bind(login); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(login); err != nil {
		return err
	}

	user, err := a.queries.SelectUserByUsername(a.ctx, login.Username)
	if err != nil {
		return echo.ErrNotFound
	}

	if !CheckPasswordHash(login.Password, user.Password) {
		return echo.ErrUnauthorized
	}

	claims := &JwtCustomClaims{
		user.ID,
		false,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

// @Summary Register a new user
// @Description Register a new user with the provided username and password
// @Tags authentication
// @Accept json
// @Produce json
// @Param credentials body Login true "User registration credentials"
// @Success 201 {object} database.User "Successfully registered user details"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 409 {object} map[string]string "Conflict - User with the given username already exists"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /register [post]
func (a *App) Register(c echo.Context) error {
	login := new(Login)
	if err := c.Bind(login); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(login); err != nil {
		return err
	}

	_, err := a.queries.SelectUserByUsername(a.ctx, login.Username)
	if err == nil {
		return echo.ErrConflict
	}

	hashed_password, err := HashPassword(login.Password)
	if err != nil {
		return echo.ErrInternalServerError
	}

	user, err := a.queries.InsertUser(a.ctx, database.InsertUserParams{
		Username: login.Username,
		Password: hashed_password,
	})
	if err != nil {
		fmt.Println(err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, user)
}
