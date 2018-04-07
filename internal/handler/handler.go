package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/YuheiNakasaka/radiorec/internal/db"
	"github.com/labstack/echo"
)

// Index is root path handler for ping
func Index(c echo.Context) error {
	return c.String(http.StatusOK, "Welcome to radiorec!")
}

// Programs return program lists
func Programs(c echo.Context) error {
	mydb := &db.MyDB{}
	err := mydb.New()
	if err != nil {
		return fmt.Errorf("Failed to connect database: %v", err)
	}

	// offset and limit params are required
	var offset, limit int
	if offset, err = strconv.Atoi(c.QueryParam("offset")); err != nil {
		return fmt.Errorf("Failed to convert offset to int: %v", err)
	}
	if limit, err = strconv.Atoi(c.QueryParam("limit")); err != nil {
		return fmt.Errorf("Failed to convert limit to int: %v", err)
	}

	conn := mydb.Connection
	results := []db.ProgramJoinsProgramContent{}
	conn.Table("programs").
		Select("*").
		Joins("inner join program_contents on programs.id = program_contents.program_id").
		Order("program_contents.id desc").Offset(offset).Limit(limit).Scan(&results)

	c.Response().Header().Set("Access-Control-Allow-Origin", "*")
	c.Response().Header().Set("Access-Control-Allow-Methods", "GET,HEAD")
	return c.JSON(http.StatusOK, results)
}
