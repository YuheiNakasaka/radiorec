create database radiorec character set utf8mb4;

create table radiorec.programs (
  id int auto_increment not null primary key,
  name varchar(512),
  cast varchar(255),
  day_of_week int not null,
  start_time time not null,
  airtime int not null,
  station int not null,
  on_air_status int not null default 0
);

create table radiorec.program_contents (
  id int auto_increment not null primary key,
  program_id int not null,
  video_path varchar(512),
  created_at datetime not null,
  updated_at timestamp DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

-- seed
INSERT INTO `programs` VALUES
(1,'佐倉としたい大西','佐倉綾音,大西沙織',2,'23:30:00',1800),
(2,'井口裕香のむ〜〜〜ん','井口裕香',1,'22:00:00',1800),
(3,'西明日香のデリケートゾーン','西明日香',2,'01:00:00',1800),
(4,'花澤香菜のひとりでできるかな？','花澤香菜',4,'23:00:00',1800),
(5,'早見沙織のふり～すたいる','早見沙織',6,'19:30:00',1800),
(6,'上坂すみれの♡をつければ可愛かろう','上坂すみれ',0,'00:00:00',1800),
(7,'水瀬いのりのmelody flag','水瀬いのり',0,'22:00:00',1800),
(8,'碧と彩奈のラ・プチミレディオ flag','悠木碧,竹達彩奈',0,'22:30:00',1800),
(9,'内田真礼のおはなししません？','内田真礼',6,'20:30:00',1800),
(10,'村川梨衣の a りえしょんぷり〜ず','村川梨衣',2,'23:00:00',1800),
(11,'内田雄馬 君の話を焼かせて','内田雄馬',3,'21:00:00',1800);
(12,'洲崎西','洲崎綾,西明日香',3,'01:00:00',1800),
(13,'大橋彩香のAny Beats!','大橋彩香',3,'19:00:00',1800),
(14,'花澤香菜・内山夕実のクロ香菜さんとシロ夕実さん','花澤香菜,内山夕実',1,19:30:00,1800),
(15,'鷲崎健のヨルナイトxヨルナイト','鷲崎健',2,'00:00:00',3600),
(16,'鷲崎健のヨルナイトxヨルナイト','鷲崎健',3,'00:00:00',3600),
(17,'鷲崎健のヨルナイトxヨルナイト','鷲崎健',4,'00:00:00',3600),
(18,'鷲崎健のヨルナイトxヨルナイト','鷲崎健',5,'00:00:00',3600),
(19,'小倉唯のyui※room','小倉唯',1,'00:00:30',1800),
(20,'三澤紗千香のラジオを聴くじゃんね！','三澤紗千香',5,'02:00:00',1800),
(21,'徳井青空のまぁるくなぁれ！','徳井青空',5,'01:00:00',1800);
