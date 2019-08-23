#!/bin/sh

`rtmpdump -q -r rtmp://fms-base2.mitene.ad.jp/agqr/aandg1 --live --stop 1 -o /tmp/dead_or_live.flv`

CONTENT=`file /tmp/dead_or_live.flv | grep empty`

if [ -n "${CONTENT}" ]; then
  echo "URL dead"
  curl -s -S -X POST --data-urlencode "payload={\"text\": \"A&G endpoint has been dead\"}" 'YOUR_SLACK_WEBHOOK_URL'
fi