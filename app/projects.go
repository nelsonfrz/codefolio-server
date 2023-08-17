package app

import (
	"codefolio/server/database"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// @Summary Retrieve all projects
// @Description Retrieve all projects with optional pagination
// @Tags projects
// @Produce json
// @Param limit query int false "Limit the number of returned projects (default is 25)"
// @Param offset query int false "Starting point for project retrieval (default is 0)"
// @Success 200 {array} database.Project "Successfully retrieved projects"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /projects [get]
func (a *App) GetAllProjects(c echo.Context) error {
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
	users, err := a.queries.SelectAllProjects(a.ctx, database.SelectAllProjectsParams{
		Limit:  int32(parsed_limit),
		Offset: int32(parsed_offset),
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, users)
}

// @Summary Retrieve all projects (authorized)
// @Description Retrieve all authorized projects with optional pagination
// @Tags projects
// @Produce json
// @Param limit query int false "Limit the number of returned projects (default is 25)"
// @Param offset query int false "Starting point for project retrieval (default is 0)"
// @Security ApiKeyAuth
// @Header 200 {string} Token "Bearer"
// @Success 200 {array} database.Project "Successfully retrieved authorized projects"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /authorized/projects [get]
func (a *App) GetAllProjectsAuthorized(c echo.Context) error {
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

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)

	users, err := a.queries.SelectAllProjects(a.ctx, database.SelectAllProjectsParams{
		Limit:  int32(parsed_limit),
		Offset: int32(parsed_offset),
		UserID: sql.NullInt32{
			Valid: true,
			Int32: claims.UserId,
		},
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, users)
}

type Project struct {
	Name          string `json:"name" validate:"required"`
	ThumbnailUrl  string `json:"thumbnail_url"`
	Description   string `json:"description"`
	Content       string `json:"content"`
	Visibility    string `json:"visibility"`
	SourceCodeUrl string `json:"source_code_url"`
	DeploymentUrl string `json:"deployment_url"`
}

// @Summary Create a new project
// @Description Create a new project associated with the authenticated user
// @Tags projects
// @Accept json
// @Produce json
// @Param project body Project true "Project body to be added"
// @Security ApiKeyAuth
// @Header 200 {string} Token "Bearer"
// @Success 201 {object} database.Project "Successfully created project"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /authorized/projects [post]
func (a *App) CreateProject(c echo.Context) error {
	inputProject := new(Project)
	if err := c.Bind(inputProject); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(inputProject); err != nil {
		return err
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)

	visibility := database.VisibilityPublic
	if inputProject.Visibility != "public" {
		visibility = database.VisibilityPrivate
	}

	project, err := a.queries.InsertProject(a.ctx, database.InsertProjectParams{
		Name: inputProject.Name,
		ThumbnailUrl: sql.NullString{
			Valid:  inputProject.ThumbnailUrl != "",
			String: inputProject.ThumbnailUrl,
		},
		Description: sql.NullString{
			Valid:  inputProject.Description != "",
			String: inputProject.Description,
		},
		Visibility: visibility,
		Content: sql.NullString{
			Valid:  inputProject.Content != "",
			String: inputProject.Content,
		},
		SourceCodeUrl: sql.NullString{
			Valid:  inputProject.SourceCodeUrl != "",
			String: inputProject.SourceCodeUrl,
		},
		DeploymentUrl: sql.NullString{
			Valid:  inputProject.DeploymentUrl != "",
			String: inputProject.DeploymentUrl,
		},
		UserID: sql.NullInt32{
			Valid: true,
			Int32: claims.UserId,
		},
	})
	if err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, project)
}

type DeleteProject struct {
	ID int64 `json:"id"`
}

// @Summary Delete a project
// @Description Delete a project by ID associated with the authenticated user
// @Tags projects
// @Accept json
// @Produce json
// @Param project body DeleteProject true "Project ID to be deleted"
// @Security ApiKeyAuth
// @Header 200 {string} Token "Bearer"
// @Success 200 {string} string "Successfully deleted project with message"
// @Failure 400 {object} map[string]string "Bad Request"
// @Failure 401 {object} map[string]string "Unauthorized"
// @Failure 404 {object} map[string]string "Not Found"
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /authorized/projects [delete]
func (a *App) DeleteProject(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)

	u := new(DeleteProject)
	if err := c.Bind(u); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(u); err != nil {
		return err
	}

	project, err := a.queries.SelectProject(a.ctx, int32(u.ID))
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	if project.UserID.Int32 != claims.UserId {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	err = a.queries.DeleteProject(a.ctx, database.DeleteProjectParams{
		ID: int32(u.ID),
		UserID: sql.NullInt32{
			Valid: true,
			Int32: claims.UserId,
		},
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.String(http.StatusOK, "deleted project with id "+strconv.FormatInt(u.ID, 10))
}

type UpdatedProject struct {
	ID            int32  `json:"id" validate:"required"`
	Name          string `json:"name"`
	ThumbnailUrl  string `json:"thumbnail_url"`
	Description   string `json:"description"`
	Content       string `json:"content"`
	Visibility    string `json:"visibility"`
	SourceCodeUrl string `json:"source_code_url"`
	DeploymentUrl string `json:"deployment_url"`
}

func (a *App) UpdateProject(c echo.Context) error {
	input := new(UpdatedProject)
	if err := c.Bind(input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(input); err != nil {
		return err
	}

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JwtCustomClaims)

	visibility := database.VisibilityPublic
	if input.Visibility != "public" {
		visibility = database.VisibilityPrivate
	}

	project, err := a.queries.SelectProject(a.ctx, input.ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	if project.UserID.Int32 != claims.UserId {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}

	project.Name = input.Name
	project.Content = sql.NullString{String: input.Content, Valid: input.Content != ""}
	project.Description = sql.NullString{String: input.Description, Valid: input.Description != ""}
	project.DeploymentUrl = sql.NullString{String: input.DeploymentUrl, Valid: input.DeploymentUrl != ""}
	project.SourceCodeUrl = sql.NullString{String: input.SourceCodeUrl, Valid: input.SourceCodeUrl != ""}
	project.ThumbnailUrl = sql.NullString{String: input.ThumbnailUrl, Valid: input.ThumbnailUrl != ""}
	project.Visibility = visibility

	project, err = a.queries.UpdateProject(a.ctx, database.UpdateProjectParams{
		ID:            input.ID,
		Name:          input.Name,
		UserID:        sql.NullInt32{Int32: claims.UserId, Valid: true},
		ThumbnailUrl:  project.ThumbnailUrl,
		Description:   project.Description,
		Content:       project.Content,
		Visibility:    project.Visibility,
		SourceCodeUrl: project.SourceCodeUrl,
		DeploymentUrl: project.DeploymentUrl,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, project)
}
