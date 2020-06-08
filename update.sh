#!/bin/bash 

killall useless_bot
git pull -f https://github.com/odavzo/Useless-Bot.git
go build -o ./useless_bot ./main.go
./useless_bot -t NzE3MzI4ODMxMjYxNzA0MjUz.XtjreQ.mLUQDcrXUmtnx6wBBoNG4gWvKx4 2> ./log/error.log > /log/out.log &

