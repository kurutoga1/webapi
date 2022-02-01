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
  serverDir=${buildDir}/pro${i}
  mkdir -p ${serverDir}
  rm -fr ${serverDir}/pro
  rm -fr ${serverDir}/tmplates
  go build -o ${serverDir}/pro pro/main.go
  cp -r pro/config/templates ${serverDir}/
  cp -r pro/static ${serverDir}/
done

echo "cli"
localDir=${buildDir}/cli
mkdir -p ${localDir}
go build -o ${localDir}/cli cli/main.go
