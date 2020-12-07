package ag

import (
	"strconv"

	"github.com/YuheiNakasaka/radiorec/internal/recorder"
)

// Ag is a&g+ struct
type Ag struct {
	programID int
	airtime   int
	storage   string
}

// ProgramID is method to fill recorder.Recorder interface.
func (a *Ag) ProgramID() int {
	return a.programID
}

// Airtime is method to fill recorder.Recorder interface.
func (a *Ag) Airtime() int {
	return a.airtime
}

// Storage is method to fill recorder.Recorder interface.
func (a *Ag) Storage() string {
	return a.storage
}

// RecordCommand is method to fill recorder.Recorder interface.
// It returns rtmpdump command to record during airtime.
func (a *Ag) RecordCommand(outputPath string) string {
	return "ffmpeg -loglevel error -y -fflags +discardcorrupt -allowed_extensions ALL -protocol_whitelist file,crypto,http,https,tcp,tls -i https://fms2.uniqueradio.jp/agqr10/aandg1.m3u8 -t " + strconv.Itoa(a.airtime) + " -vcodec copy -acodec aac -bsf:a aac_adtstoasc " + outputPath + ".mp4"
}

// Start : record ag program
func (a *Ag) Start(programID int, airtime int, storage string) error {
	ag := &Ag{}
	ag.programID = programID
	ag.airtime = airtime
	ag.storage = storage
	return recorder.Record(ag)
}
