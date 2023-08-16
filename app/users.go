package app

import (
	"codefolio/server/database"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// @Summary Retrieve all users
// @Description Retrieve all users with optional pagination
// @Tags users
// @Produce json
// @Param limit query int false "Limit the number of returned users (default is 25)"
// @Param offset query int false "Starting point for user retrieval (default is 0)"
// @Success 200 {array} database.User "Successfully retrieved users"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /users [get]
func (a *App) GetAllUsers(c echo.Context) error {
	limit := c.QueryParam("limit")
	offset := c.QueryParam("offset")

	parsed_limit, err := strconv.ParseInt(limit, 10, 32)
	if err != nil {
		parsed_limit = 25
	}

	parsed_offset, err := strconv.ParseInt(offset, 10, 32)
	if err != nil {
		parsed_offset = 0
	}

	users, err := a.queries.SelectAllUsersWithoutPassword(a.ctx, database.SelectAllUsersWithoutPasswordParams{
		Limit:  int32(parsed_limit),
		Offset: int32(parsed_offset),
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, users)
}
