package handler

import (
	"fmt"
	"net/http"

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

	// Fetch programs
	offset := 0
	limit := 10
	conn := mydb.Connection
	results := []db.ProgramJoinsProgramContent{}
	conn.Table("programs").
		Select("*").
		Joins("inner join program_contents on programs.id = program_contents.program_id").
		Order("program_contents.id desc").Offset(offset).Limit(limit).Scan(&results)

	return c.String(http.StatusOK, "Programs")
}
