#!/bin/zsh

buildDir=./build
mkdir -p ${buildDir}

# ゲートウェイサーバは３台
for i in 1 2 3
do
  echo "gate way"${i}
  lbDir=${buildDir}/gw${i}
  mkdir -p ${lbDir}
  rm -fr ${lbDir}/gw
  go build -o ${lbDir}/gw gw/main.go
  cp -r gw/templates ${lbDir}/
  cp -r gw/static ${lbDir}/
#  cp gw/config.json ${lbDir}/
done

# プログラムサーバは３台
for i in 1 2 3
do
  echo "server"${i}
  serverDir=${buildDir}/server${i}
  mkdir -p ${serverDir}
  rm -fr ${serverDir}/server
  rm -fr ${serverDir}/tmplates
  go build -o ${serverDir}/server server/main.go
  cp -r server/config/templates ${serverDir}/
  cp -r server/static ${serverDir}/
done

echo "cli"
localDir=${buildDir}/cli
mkdir -p ${localDir}
go build -o ${localDir}/cli cli/main.go
