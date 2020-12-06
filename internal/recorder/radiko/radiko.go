package radiko

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/YuheiNakasaka/radiorec/config"
	"github.com/YuheiNakasaka/radiorec/internal/filemanager"

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
	storage   string
	channel   string
}

// URLTags is xml type
type URLTags struct {
	Lists []URLTag `xml:"url"`
}

// URLTag is xml type
type URLTag struct {
	AreaType string `xml:"areafree,attr"`
	Text     string `xml:"playlist_create_url"`
}

// ProgramID is method to fill recorder.Recorder interface.
func (r *Radiko) ProgramID() int {
	return r.programID
}

// Airtime is method to fill recorder.Recorder interface.
func (r *Radiko) Airtime() int {
	return r.airtime
}

// Storage is method to fill recorder.Recorder interface.
func (r *Radiko) Storage() string {
	return r.storage
}

// RecordCommand is method to fill recorder.Recorder interface.
// It returns rtmpdump command to record during airtime.
func (r *Radiko) RecordCommand(outputPath string) string {
	url := fetchStreamURL(r.channel)
	authToken := authorize()
	return "ffmpeg -loglevel error -y -fflags +discardcorrupt -headers 'X-Radiko-Authtoken: " + authToken + "' -allowed_extensions ALL -protocol_whitelist file,crypto,http,https,tcp,tls -i " + url + " -t " + strconv.Itoa(r.airtime) + " -vcodec copy -acodec copy -bsf:a aac_adtstoasc " + outputPath + ".mp4"
}

func fetchStreamURL(channel string) string {
	client := &http.Client{}

	values := url.Values{}
	req1, err := http.NewRequest("GET", "http://radiko.jp/v2/station/stream_smh_multi/"+channel+".xml", strings.NewReader(values.Encode()))
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req1)
	if err != nil {
		panic(err)
	}
	buf, _ := ioutil.ReadAll(resp.Body)

	data := new(URLTags)
	if err := xml.Unmarshal(buf, data); err != nil {
		panic(err)
	}

	defer resp.Body.Close()
	for i := 0; i < len(data.Lists); i++ {
		if data.Lists[i].AreaType == "0" {
			return data.Lists[i].Text
		}
	}
	return ""
}

// Start : record ag program
func (r *Radiko) Start(programID int, airtime int, storage string, channel int) error {
	radiko := &Radiko{}
	radiko.programID = programID
	radiko.airtime = airtime
	radiko.storage = storage
	if channel == 2 {
		radiko.channel = "QRR"
	} else if channel == 3 {
		radiko.channel = "HELLOFIVE"
	}
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
	keyPath := filepath.Join(parentDir, "radikokey2.txt")

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
	req1, err := http.NewRequest("GET", "https://radiko.jp/v2/api/auth1", strings.NewReader(values.Encode()))
	if err != nil {
		panic(err)
	}
	req1.Header.Add("pragma", "no-cache")
	req1.Header.Add("X-Radiko-App", "pc_html5")
	req1.Header.Add("X-Radiko-App-Version", "0.0.1")
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
	req2, err := http.NewRequest("GET", "https://radiko.jp/v2/api/auth2", strings.NewReader(values2.Encode()))
	if err != nil {
		panic(err)
	}
	req2.Header.Add("pragma", "no-cache")
	req2.Header.Add("X-Radiko-User", "test-stream")
	req2.Header.Add("X-Radiko-Device", "pc")
	req2.Header.Add("X-Radiko-Authtoken", autoToken)
	req2.Header.Add("X-Radiko-Partialkey", partialkey)

	resp2, err := client.Do(req2)
	if err != nil {
		panic(err)
	}

	buf2, _ := ioutil.ReadAll(resp2.Body)
	fmt.Println(string(buf2))

	defer resp.Body.Close()

	return autoToken
}
