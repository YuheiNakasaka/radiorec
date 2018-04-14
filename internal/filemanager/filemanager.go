package filemanager

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	uuid "github.com/satori/go.uuid"
)

// FileManager : manage files and directories
type FileManager struct {
	FilePath   string
	OutputPath string
}

// StaticFilePath : return directory path to place static files (e.g mp4, flv, swf..)
func StaticFilePath() (string, error) {
	outputParentDir := os.Getenv("OUTPUT_DIR")
	if outputParentDir == "" {
		b, err := exec.Command("go", "env", "GOPATH").CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("Failed to find GOPATH: %v", err)
		}
		for _, p := range filepath.SplitList(strings.TrimSpace(string(b))) {
			p = filepath.Join(p, filepath.FromSlash("/src/github.com/YuheiNakasaka/radiorec/public"))
			outputParentDir = p
		}
	}
	if outputParentDir == "" {
		return "", fmt.Errorf("Failed to create directory path")
	}
	return outputParentDir, nil
}

// PreparePaths makes directories if not exists
func (f *FileManager) PreparePaths(programID int) error {
	outputParentDir, err := StaticFilePath()
	if err != nil {
		return fmt.Errorf("Failed to get file path: %v", err)
	}

	// なんかイケてない...
	// filePath: /public/123/abcd-efgh-ijkl-mnop => need to register record
	// outputDir: /src/github.com/YuheiNakasaka/radiorec/public/123/
	// outputPath: /src/github.com/YuheiNakasaka/radiorec/public/123/abcd-efgh-ijkl-mnop
	filename := uuid.NewV4().String()
	strProgramID := strconv.Itoa(programID)
	f.FilePath = filepath.Join(string(os.PathSeparator), "public", strProgramID, filename)
	f.OutputPath = filepath.Join(outputParentDir, strProgramID, filename)
	outputDir := filepath.Join(outputParentDir, strProgramID)

	_, err = os.Stat(outputDir)
	if err != nil {
		return os.MkdirAll(outputDir, 0777)
	}
	return err
}
