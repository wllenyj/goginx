#!/bin/sh
PORT=23523
./http_svr -p $PORT
sleep 2
for i in `seq 0 10`
do
    curl "localhost:${PORT}/hello" &
done

./http_svr -s restart

#sleep 1

for i in `seq 0 10`
do
    curl "localhost:${PORT}/hello1" &
done

sleep 1

./http_svr -s stop
