package ag

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/YuheiNakasaka/radiorec/internal/db"
	"github.com/mattn/go-shellwords"
	"github.com/satori/go.uuid"
)

var filePath = ""
var outputPath = ""

// preparePaths: make directory if not exists
func preparePaths(programID int) error {
	b, err := exec.Command("go", "env", "GOPATH").CombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed to find GOPATH: %v", err)
	}
	outputParentDir := ""
	for _, p := range filepath.SplitList(strings.TrimSpace(string(b))) {
		p = filepath.Join(p, filepath.FromSlash("/src/github.com/YuheiNakasaka/radiorec/public"))
		outputParentDir = p
	}
	if outputParentDir == "" {
		return fmt.Errorf("Failed to create directory path: %v", err)
	}

	// なんかイケてない...
	// filePath: /public/123/abcd-efgh-ijkl-mnop => need to register record
	// outputDir: /src/github.com/YuheiNakasaka/radiorec/public/123/
	// outputPath: /src/github.com/YuheiNakasaka/radiorec/public/123/abcd-efgh-ijkl-mnop
	filename := uuid.NewV4().String()
	strProgramID := strconv.Itoa(programID)
	filePath = filepath.Join(string(os.PathSeparator), "public", strProgramID, filename)
	outputPath = filepath.Join(outputParentDir, strProgramID, filename)
	outputDir := filepath.Join(outputParentDir, strProgramID)

	_, err = os.Stat(outputDir)
	if err != nil {
		return os.MkdirAll(outputDir, 0777)
	}
	return err
}

// Start : record ag program
func Start(programID int, airtime int) error {
	// check args
	if programID == 0 {
		return fmt.Errorf("Could not set 0 as programID")
	}
	if airtime == 0 {
		return fmt.Errorf("Could not set 0 as airtime")
	}

	// get db connection
	mydb := &db.MyDB{}
	err := mydb.New()
	if err != nil {
		return fmt.Errorf("Failed to connect database: %v", err)
	}

	if mydb.ValidProgramID(programID) == false {
		return fmt.Errorf("Failed to find the program id: %v", programID)
	}

	// create output path and filenames
	err = preparePaths(programID)
	if err != nil {
		return fmt.Errorf("Failed to make directory: %v", err)
	}
	fmt.Println(filePath, outputPath)

	// record as live streaming
	recExt := ".flv"
	recCmd := "rtmpdump -q -r rtmp://fms-base2.mitene.ad.jp/agqr/aandg2 --live --stop " + strconv.Itoa(airtime) + " -o " + outputPath + recExt
	fmt.Println(recCmd)

	// convert flv to mp4
	mp4Ext := ".mp4"
	mp4Cmd := "ffmpeg -y -i " + outputPath + ".flv -acodec aac -vcodec h264 " + outputPath + mp4Ext
	fmt.Println(mp4Cmd)

	// wait for finishing to record
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		// start recording
		fmt.Println("Recording...")
		recC, parseErr := shellwords.Parse(recCmd)
		if parseErr != nil {
			return
		}
		out, cmdErr := exec.Command(recC[0], recC[1:]...).Output()
		fmt.Println(out, cmdErr)

		// start converting
		fmt.Println("Converting...")
		convC, convErr := shellwords.Parse(mp4Cmd)
		if convErr != nil {
			return
		}
		convO, convE := exec.Command(convC[0], convC[1:]...).Output()
		fmt.Println(convO, convE)

		// remove src flv file
		if rmErr := os.Remove(outputPath + ".flv"); rmErr != nil {
			return
		}

		// register data to table
		fmt.Println("Registering...")
		mydb.InsertProgramContent(programID, filePath+".mp4")

		// S3にアップロード
		// uploader.Upload(outputPath+".mp4", filePath+".mp4")

		wg.Done()
	}()
	wg.Wait()

	return err
}
