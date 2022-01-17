#!/bin/zsh

rm -fr test

for i in `seq 100`
do
mkdir -p test/

timeDir=${i}-$(date +%F-%T)
mkdir -p test/${timeDir}

go run main.go -url http://127.0.0.1:8001 -name convertToJson -i test.txt -o test/${timeDir} &

#if [ -e test/${timeDir}/test.txt.json.pcap ]; then
#  mv test/${timeDir} test/${timeDir}-ok
#else
#  mv test/${timeDir} test/${timeDir}-no
#fi

done

echo end
