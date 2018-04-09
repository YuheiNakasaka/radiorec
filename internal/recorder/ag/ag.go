package ag

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"sync"

	"github.com/YuheiNakasaka/radiorec/internal/db"
	"github.com/YuheiNakasaka/radiorec/internal/filemanager"
	"github.com/mattn/go-shellwords"
)

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
	fileManager := &filemanager.FileManager{}
	err = fileManager.PreparePaths(programID)
	if err != nil {
		return fmt.Errorf("Failed to make directory: %v", err)
	}
	fmt.Println(fileManager.FilePath, fileManager.OutputPath)

	// record as live streaming
	recExt := ".flv"
	recCmd := "rtmpdump -q -r rtmp://fms-base2.mitene.ad.jp/agqr/aandg2 --live --stop " + strconv.Itoa(airtime) + " -o " + fileManager.OutputPath + recExt
	fmt.Println(recCmd)

	// convert flv to mp4
	mp4Ext := ".mp4"
	mp4Cmd := "ffmpeg -y -i " + fileManager.OutputPath + ".flv -acodec aac -vcodec h264 " + fileManager.OutputPath + mp4Ext
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
		if rmErr := os.Remove(fileManager.OutputPath + ".flv"); rmErr != nil {
			return
		}

		// register data to table
		fmt.Println("Registering...")
		mydb.InsertProgramContent(programID, fileManager.FilePath+".mp4")

		// S3にアップロード
		// uploader.Upload(outputPath+".mp4", filePath+".mp4")

		wg.Done()
	}()
	wg.Wait()

	return err
}
