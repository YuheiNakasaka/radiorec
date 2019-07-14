# 参考: https://qiita.com/taiko19xx/items/215b9943c8aa0d8edcf6
#!/bin/bash

BACKUP_PATH="/your/backup/path/"
S3_PATH="s3://your-s3-bucket/"
FILE_NAME="mysql_dump_`date +%Y%m%d`.sql.gz"

AWS_CONFIG_FILE="/home/user/.aws/config"

cd $BACKUP_PATH
mysqldump -u mysql_user -p mysql_password --all-databases | gzip > $FILE_NAME

find $BACKUP_PATH -type f -name "mysql_dump_*.sql.gz"  -mtime +9 -daystart | xargs rm -rf

aws s3 sync $BACKUP_PATH $S3_PATH  --delete