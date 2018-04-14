# radiorec
自分専用の声優ラジオ録音環境です。

データベースに登録した番組を放送時間になったら録音して既定のディレクトリ&(任意で外部ストレージ)へ保存するだけです。

保存した情報を返す簡易のAPIサーバーも一応あります。

# How it works

流れは以下。

- 録音したい番組をテーブルに登録
- 各放送時間に起動するようにcrontabを書き換える
- 放送時間になったら録音コマンドが動く
- rtmpdumpでストリームを録音
- ffmpegでmp4へ変換する
- ローカルに保存
- 任意でs3へ保存

また、**radikoには地域制限があるのでTokyo regionのEC2を利用することを推奨します(GCPサービスやさくらのVPSではIPが海外だったり大阪だったりすることがあるのでダメ)**。

# Requirements
- Go(>=1.10.x)
- MySQL(>=5.6)
- swftools
- ffmpeg
- rtmpdump
- **EC2 Tokyo region instance**
  - to avoid region restrictions

# Supported Radios
- 超！A&G([番組表](http://www.agqr.jp/timetable/streaming.html))
- radiko([アーカイブ](http://www.joqr.co.jp/programs/daily-programsheet.php?date=20171001))

# How to use

** ※Requirementsはすべて導入してる前提です **

### 1. 設定ファイルの準備

`config/config.exmpale.yml`を書き換えて、`config/config.yml`として保存する。

### 2. テーブルのセットアップ

`radiorec/migrate/schemra.sql`内のSQLをデータベースに実行。

(僕が録音している声優のおすすめラジオが20番組ほど登録されます。)

### 3. 依存関係

```
dep ensure
```

### 4. cron設定

下記コマンドを実行するとcrontabが書き換わる。

`CONFIG_PATH`には先ほど作った`config/config.yml`の場所を指定する。

ちなみに** 既存のcrontabの内容を上書きするので他に何か大事な設定をしている場合はヤバイ... **

##### Example

```
CONFIG_PATH=/var/www/app/radiodic/config go run cmd/cli/cli.go cron
```

### 4. 試し

試しに下記の録音コマンドを実行してみて、エラーが出なければ大丈夫です。

```
$ CONFIG_PATH=/var/www/app/radiodic/config go run cmd/cli record -id 1

/public/1/09f706c5-b7f3-43ec-9beb-7f1fe44fb0ab /Users/razokulover/src/github.com/YuheiNakasaka/radiorec/public/1/09f706c5-b7f3-43ec-9beb-7f1fe44fb0ab
rtmpdump -q -r rtmp://fms-base2.mitene.ad.jp/agqr/aandg2 --live --stop 3 -o /Users/razokulover/src/github.com/YuheiNakasaka/radiorec/public/1/09f706c5-b7f3-43ec-9beb-7f1fe44fb0ab.flv
ffmpeg -y -i /Users/razokulover/src/github.com/YuheiNakasaka/radiorec/public/1/09f706c5-b7f3-43ec-9beb-7f1fe44fb0ab.flv -acodec aac -vcodec h264 /Users/razokulover/src/github.com/YuheiNakasaka/radiorec/public/1/09f706c5-b7f3-43ec-9beb-7f1fe44fb0ab.mp4
Recording...
[] <nil>
Converting...
[] <nil>
Registering...
```

# Cli
録音やcron設定のコマンド

### Record
指定されたラジオの録音を行う。

##### Parameters

- id, i
  - 必須
  - Programsテーブルに登録されているラジオのID
- storage, s
  - 任意
  - 外部のストレージに録音ファイルを保存する場合。デフォルトはローカルのみ保存。
- CONFIG_PATH
  - 任意
  - 環境変数
  - 設定ファイル(config.yml)の場所を指定
  - デフォルトは`$GOPATH/src/github.com/YuheiNakasaka/radiorec/config`

##### Example
```
CONFIG_PATH=/etc/app/config go run cmd/cli/cli.go record -id 123 -storage s3
```

### Cron
Progrmsテーブルに設定されている放送時間に録音コマンドが動くようにcrontabを書き換える

##### Parameters
- CONFIG_PATH
  - 任意
  - 環境変数
  - 設定ファイル(config.yml)の場所を指定
  - デフォルトは`$GOPATH/src/github.com/YuheiNakasaka/radiorec/config`

##### Example
```
CONFIG_PATH=/etc/app/config go run cmd/cli/cli.go cron
```

# Endpoint
録音した番組情報を返すAPI。録音した音源を自作したアプリなどで利用する時に使える。

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

# Disclaimer
これは個人用途に利用する目的で作っています。

録音したラジオを複数人へ配信することを推奨しているわけではありません。

自分の責任の範囲で利用してください。

# TODO
- [x] recorderの各種スクリプト
- [x] recorderの設計の見直し
  - internal/recorder/配下の各放送局は`RecorderBase interface`を実装する
  - `RecorderBase`では`exec.Command`で実行できる録音コマンドだけ返すメソッドを定義
  - `RecorderBase`を引数に取る`recorder.Generate(rb RecorderBase)`みたいな共通関数を作る
  - 各放送局のスクリプトでは`RecorderBase`のコマンド実装とそれを`recorder.Generate(rb RecorderBase)`に渡すだけみたいな実装だけに絞る
  - これによって放送局が増えた時は`RecorderBase`のコマンド実装だけすれば基本的にはやることは終わりになるのでメンテしやすくなりそう
- [x] アップロードするスクリプト
- [x] cron設定するプロセス
  - crontabをclearして再生成する雑な作り。既存のcrontabを消すから危険。
- [x] typeのカラムを追加するか考慮する
  - on airしてるかどうかのフラグとラジオ局を表すカラムを追加した
- [ ] テスト
- [ ] 各種ミドルウェアが入った開発用docker-compose
- [ ] デプロイ

# License

The library is available as open source under the terms of the MIT License.
