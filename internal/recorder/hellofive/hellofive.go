package hellofive

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/YuheiNakasaka/radiorec/config"
	"github.com/YuheiNakasaka/radiorec/internal/filemanager"
	"github.com/mattn/go-shellwords"

	"github.com/YuheiNakasaka/radiorec/internal/recorder"
)

var (
	swfURL  = "http://radiko.jp/apps/js/flash/myplayer-release.swf"
	rtmpURL = "rtmpe://f-radiko.smartstream.ne.jp"
)

// Hellofive is radiko struct
type Hellofive struct {
	programID int
	airtime   int
	storage   string
}

// ProgramID is method to fill recorder.Recorder interface.
func (r *Hellofive) ProgramID() int {
	return r.programID
}

// Airtime is method to fill recorder.Recorder interface.
func (r *Hellofive) Airtime() int {
	return r.airtime
}

// Storage is method to fill recorder.Recorder interface.
func (r *Hellofive) Storage() string {
	return r.storage
}

// RecordCommand is method to fill recorder.Recorder interface.
// It returns rtmpdump command to record during airtime.
func (r *Hellofive) RecordCommand(outputPath string) string {
	authToken := authorize()
	// HELLOFIVE
	return "rtmpdump -q -r " + rtmpURL + " --playpath 'simul-stream.stream' --app 'HELLOFIVE/_definst_' -W " + swfURL + " -C S:'' -C S:'' -C S:'' -C S:" + authToken + " --live --stop " + strconv.Itoa(r.airtime) + " -o " + outputPath + ".flv"
}

// Start : record ag program
func (r *Hellofive) Start(programID int, airtime int, storage string) error {
	radiko := &Hellofive{}
	radiko.programID = programID
	radiko.airtime = airtime
	radiko.storage = storage
	return recorder.Record(radiko)
}

// authorize : 認証用のswf取得
func authorize() string {
	// Read config file
	myconf := config.Config{}
	err := myconf.Init()
	if err != nil {
		panic(err)
	}

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

	// ログイン
	// 取得したcookieはcookiejarにセットされる。後のリクエストではclientが同じであれば使い回される
	logingURL := "https://radiko.jp/ap/member/login/login"
	values0 := url.Values{}
	values0.Add("mail", fmt.Sprintf("%v", myconf.List.Get("radiko.login.mail")))
	values0.Add("pass", fmt.Sprintf("%v", myconf.List.Get("radiko.login.password")))
	req0, err := http.NewRequest("POST", logingURL, strings.NewReader(values0.Encode()))
	if err != nil {
		panic(err)
	}
	req0.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	jar, _ := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	resp0, err := client.Do(req0)
	if err != nil {
		panic(err)
	}
	defer resp0.Body.Close()

	// 認証1
	values := url.Values{}
	req1, err := http.NewRequest("POST", "https://radiko.jp/v2/api/auth1_fms", strings.NewReader(values.Encode()))
	if err != nil {
		panic(err)
	}
	req1.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req1.Header.Add("pragma", "no-cache")
	req1.Header.Add("X-Radiko-App", "pc_ts")
	req1.Header.Add("X-Radiko-App-Version", "4.0.0")
	req1.Header.Add("X-Radiko-User", "test-stream")
	req1.Header.Add("X-Radiko-Device", "pc")

	resp, err := client.Do(req1)
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
	req1.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req2.Header.Add("pragma", "no-cache")
	req2.Header.Add("X-Radiko-App", "pc_ts")
	req2.Header.Add("X-Radiko-App-Version", "4.0.0")
	req2.Header.Add("X-Radiko-User", "test-stream")
	req2.Header.Add("X-Radiko-Device", "pc")
	req2.Header.Add("X-Radiko-Authtoken", autoToken)
	req2.Header.Add("X-Radiko-Partialkey", partialkey)

	// client2 := &http.Client{}
	resp2, err := client.Do(req2)
	if err != nil {
		panic(err)
	}

	buf2, _ := ioutil.ReadAll(resp2.Body)
	fmt.Println(string(buf2))

	defer resp.Body.Close()

	return autoToken
}
