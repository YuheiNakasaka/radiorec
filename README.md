# radiorec

自分専用のシンプルな声優ラジオ録音環境です。

# Requirement
- Go
- MySQL
- swftools
- ffmpeg
- rtmpdump

# Supported Radios
- 超！A&G([番組表](http://www.agqr.jp/timetable/streaming.html))
- radiko([アーカイブ](http://www.joqr.co.jp/programs/daily-programsheet.php?date=20171001))

# Endpoint

#### ※Setup
サーバーの立ち上げ

```
$ go run cmd/server/server.go
```

#### /programs
番組一覧を返す。

##### Parameters
- offset
  - required
  - int
- limit
  - required
  - int

##### Example
```
$ curl -XGET "http://localhost:1323/programs?offset=0&limit=1" | jq
[
  {
    "Name": "佐倉としたい大西",
    "Cast": "佐倉綾音,大西沙織",
    "DayOfWeek": 2,
    "StartTime": "MjM6MzA6MDA=",
    "Airtime": 1800,
    "ProgramID": 1,
    "VideoPath": "/hogehoge/example.mp4",
    "CreatedAt": "2017-01-01T00:00:00+09:00",
    "UpdatedAt": "2017-01-01T00:00:00+09:00"
  }
]
```
