package radiko

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/YuheiNakasaka/radiorec/internal/filemanager"
	"github.com/mattn/go-shellwords"

	"github.com/YuheiNakasaka/radiorec/internal/recorder"
)

var (
	swfURL  = "http://radiko.jp/apps/js/flash/myplayer-release.swf"
	rtmpURL = "rtmpe://f-radiko.smartstream.ne.jp"
)

// Radiko is radiko struct
type Radiko struct {
	programID int
	airtime   int
}

// ProgramID is method to fill recorder.Recorder interface.
func (r *Radiko) ProgramID() int {
	return r.programID
}

// Airtime is method to fill recorder.Recorder interface.
func (r *Radiko) Airtime() int {
	return r.airtime
}

// RecordCommand is method to fill recorder.Recorder interface.
// It returns rtmpdump command to record during airtime.
func (r *Radiko) RecordCommand(outputPath string) string {
	authToken := authorize()
	return "rtmpdump -q -r " + rtmpURL + " --playpath 'simul-stream.stream' --app 'QRR/_definst_' -W " + swfURL + " -C S:'' -C S:'' -C S:'' -C S:" + authToken + " --live --stop " + strconv.Itoa(r.airtime) + " -o " + outputPath + ".flv"
}

// Start : record ag program
func (r *Radiko) Start(programID int, airtime int) error {
	radiko := &Radiko{}
	radiko.programID = programID
	radiko.airtime = airtime
	return recorder.Record(radiko)
}

// authorize : 認証用のswf取得
func authorize() string {
	// keyファイル作成
	parentDir, _ := filemanager.StaticFilePath()
	swfPath := filepath.Join(parentDir, "player2.swf")
	keyPath := filepath.Join(parentDir, "radikokey2.png")

	// 認証用swfダウンロード
	response, err := http.Get(swfURL)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile(swfPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	file.Write(body)

	// 認証用フレーム抽出
	swfCmd := "swfextract -b 12 " + swfPath + " -o " + keyPath
	swfC, err := shellwords.Parse(swfCmd)
	if err != nil {
		panic(err)
	}
	exec.Command(swfC[0], swfC[1:]...).Run()

	// 認証1
	values := url.Values{}
	req, err := http.NewRequest("POST", "https://radiko.jp/v2/api/auth1_fms", strings.NewReader(values.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Add("pragma", "no-cache")
	req.Header.Add("X-Radiko-App", "pc_ts")
	req.Header.Add("X-Radiko-App-Version", "4.0.0")
	req.Header.Add("X-Radiko-User", "test-stream")
	req.Header.Add("X-Radiko-Device", "pc")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	autoToken := resp.Header.Get("X-Radiko-Authtoken")
	keyOffset, _ := strconv.Atoi(resp.Header.Get("X-Radiko-KeyOffset"))
	keyLength, _ := strconv.Atoi(resp.Header.Get("X-Radiko-Keylength"))
	defer resp.Body.Close()

	// 認証2
	// partial key作成
	fd, err := os.Open(keyPath)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, keyLength)
	_, err = fd.ReadAt(buf, int64(keyOffset))
	if err != nil {
		panic(err)
	}
	partialkey := base64.StdEncoding.EncodeToString(buf)

	values2 := url.Values{}
	req2, err := http.NewRequest("POST", "https://radiko.jp/v2/api/auth2_fms", strings.NewReader(values2.Encode()))
	if err != nil {
		panic(err)
	}
	req2.Header.Add("pragma", "no-cache")
	req2.Header.Add("X-Radiko-App", "pc_ts")
	req2.Header.Add("X-Radiko-App-Version", "4.0.0")
	req2.Header.Add("X-Radiko-User", "test-stream")
	req2.Header.Add("X-Radiko-Device", "pc")
	req2.Header.Add("X-Radiko-Authtoken", autoToken)
	req2.Header.Add("X-Radiko-Partialkey", partialkey)

	client2 := &http.Client{}
	resp2, err := client2.Do(req2)
	if err != nil {
		panic(err)
	}

	buf2, _ := ioutil.ReadAll(resp2.Body)
	fmt.Println(string(buf2))

	defer resp.Body.Close()

	return autoToken
}
