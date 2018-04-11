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

# TODO
- [x] recorderの各種スクリプト
- [x] recorderの設計の見直し
  - internal/recorder/配下の各放送局は`RecorderBase interface`を実装する
  - `RecorderBase`では`exec.Command`で実行できる録音コマンドだけ返すメソッドを定義
  - `RecorderBase`を引数に取る`recorder.Generate(rb RecorderBase)`みたいな共通関数を作る
  - 各放送局のスクリプトでは`RecorderBase`のコマンド実装とそれを`recorder.Generate(rb RecorderBase)`に渡すだけみたいな実装だけに絞る
  - これによって放送局が増えた時は`RecorderBase`のコマンド実装だけすれば基本的にはやることは終わりになるのでメンテしやすくなりそう
- [x] アップロードするスクリプト
- [ ] cron設定するプロセス
- [x] typeのカラムを追加するか考慮する
  - on airしてるかどうかのフラグとラジオ局を表すカラムを追加した
- [ ] 各種ミドルウェアが入った開発用docker-compose
- [ ] デプロイ
