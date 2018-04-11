package db

import (
	"fmt"
	"time"

	"github.com/YuheiNakasaka/radiorec/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// MyDB : db struct
type MyDB struct {
	Connection *gorm.DB
}

// Program : model
type Program struct {
	ID          int    `gorm:"primary_key"`
	Name        string `gorm:"size:512"`
	Cast        string `gorm:"size:255"`
	DayOfWeek   int
	StartTime   []uint8 `gorm:"type:time"`
	Airtime     int
	Station     int
	OnAirStatus int
}

// ProgramContent : model
type ProgramContent struct {
	ID        int `gorm:"primary_key"`
	ProgramID int
	VideoPath string `gorm:"size:255"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ProgramJoinsProgramContent : 番組+番組コンテンツ情報のレスポンス
type ProgramJoinsProgramContent struct {
	Program
	ProgramContent
}

// New : create db and keep connection
func (mydb *MyDB) New() error {
	// Read config file
	myconf := config.Config{}
	err := myconf.Init()
	if err != nil {
		return fmt.Errorf("Failed to load config %v", err)
	}

	// Get db connection
	dbms := "mysql"
	user := myconf.List.Get("database.user")
	password := myconf.List.Get("database.password")
	protocol := fmt.Sprintf("tcp(%v:%v)", myconf.List.Get("database.host"), myconf.List.Get("database.port"))
	dbname := myconf.List.Get("database.name")
	dialect := fmt.Sprintf("%v:%v@%v/%v?parseTime=true&loc=Japan", user, password, protocol, dbname)
	db, err := gorm.Open(dbms, dialect)

	if err != nil {
		return fmt.Errorf("Failed to connect db: %v", err)
	}

	// Set db connection
	mydb.Connection = db

	return err
}

// FindProgram : get program matched to id
func (mydb *MyDB) FindProgram(id int) Program {
	var program Program
	mydb.Connection.Find(&program, "id=?", id)
	return program
}

// ValidProgramID : check existense of Program ID
func (mydb *MyDB) ValidProgramID(id int) bool {
	var program Program
	mydb.Connection.Find(&program, "id=?", id)
	if len(program.StartTime) == 0 {
		return false
	}
	return true
}

// InsertProgramContent : register program content to table
func (mydb *MyDB) InsertProgramContent(programID int, videoPath string) {
	programContent := ProgramContent{
		ProgramID: programID,
		VideoPath: videoPath,
	}

	mydb.Connection.Create(&programContent)
}
