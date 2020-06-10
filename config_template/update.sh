#!/bin/bash 

killall useless_bot
git pull -f <git url>
go build -o ./useless_bot ./main.go
./useless_bot -t <bot_token> 2> ./log/error.log > /log/out.log &

