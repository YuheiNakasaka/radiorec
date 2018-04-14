package cron

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/YuheiNakasaka/radiorec/internal/db"
)

// convertStartTime : 時間と分を分解してmapに格納し返す
func convertStartTime(s string) map[string]string {
	words := strings.Split(s, ":")
	hour := words[0]
	minute := words[1]

	rep := regexp.MustCompile(`^0`)
	formatedHour := rep.ReplaceAllString(hour, "")
	formatedMinute := rep.ReplaceAllString(minute, "")

	return map[string]string{
		"hour":   formatedHour,
		"minute": formatedMinute,
	}
}

// generateCronLine : 設定するcron行を生成する
func generateCronLine(programID int, dayOfWeek int, times map[string]string) map[string]string {
	configDir := os.Getenv("CONFIG_DIR")
	outputParentDir := os.Getenv("OUTPUT_DIR")

	chdirCmd := "cd /var/www/radiorec/;"
	stdoutCmd := ">> /var/log/cron.log  2>&1"
	mainCmd := "CONFIG_DIR=" + configDir + " OUTPUT_DIR=" + outputParentDir + " /var/www/radiorec record -i " + strconv.Itoa(programID) + " -s s3"
	cronTime := times["minute"] + " " + times["hour"] + " * * " + strconv.Itoa(dayOfWeek)

	cmds := map[string]string{}
	cmds["schedule"] = cronTime
	cmds["exec"] = chdirCmd + mainCmd + " " + stdoutCmd

	return cmds
}

// Generate creates cron settings
func Generate() error {
	// get all records being on air.
	mydb := &db.MyDB{}
	err := mydb.New()
	if err != nil {
		return fmt.Errorf("Failed to connect database: %v", err)
	}
	conn := mydb.Connection
	results := []db.Program{}
	conn.Table("programs").
		Select("*").
		Where("programs.on_air_status = 1").
		Order("programs.id").Scan(&results)

	// create tmpfile to add crontab
	tmpFilePath := filepath.Join(filepath.FromSlash("/tmp/new_crontab_lines"))
	tmpfile, err := os.Create(tmpFilePath)
	if err != nil {
		return fmt.Errorf("Failed to create file: %v", err)
	}
	defer tmpfile.Close()
	defer os.Remove(tmpFilePath)

	// create cron lines
	for _, result := range results {
		time := convertStartTime(string(result.StartTime))
		line := generateCronLine(result.ID, result.DayOfWeek, time)
		if _, terr := tmpfile.WriteString(fmt.Sprintf("#%s\n%s\n", result.Name, line["schedule"]+" "+line["exec"])); terr != nil {
			return fmt.Errorf("Failed to write line in crontab: %v", terr)
		}
	}

	// redirect the lines to crontab
	cmd := "crontab < " + tmpFilePath
	exec.Command("sh", "-c", cmd).Run()

	return err
}
