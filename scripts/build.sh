#!/bin/bash
CURRDIR=$(cd "$(dirname "$0")"/.. && pwd -P)
Version=`date +%Y%m%d%H%M%S`-`git rev-parse HEAD | cut -c1-12`
echo $Version > $CURRDIR/cmd/leo/command/version/version.txt
cd $CURRDIR/cmd/leo && packr2 && cd -

# 目标模块
target_path="cmd"
default_target=(
    "leo"
)

tool_path="tools"
default_tool=(
    "mario"
)

target=()
default_flag=1
aarm=0
windows=0
while getopts "xwa:" opt
do
    case $opt in
        x)
        aarm=1
        ;;
        w)
        windows=1
        ;;
        a)
        default_flag=0
        target+=($OPTARG)
        ;;
    esac
done

if [ $default_flag -eq 1 ]; then
    target=(${default_target[*]})
fi

for v in ${target[@]}
do
    echo "target: $v"
done

build_args=""
if [ $aarm -eq 1 ]; then
    export CGO_ENABLED=1
    export GOOS=linux
    export GOARCH=arm64
    export CC=/mnt/eng-nfs/external/hisi-linux/x86-arm/Hi3559A_V100R001C02SPC020/aarch64-himix100-linux/bin/aarch64-himix100-linux-gcc
    export CXX=/mnt/eng-nfs/external/hisi-linux/x86-arm/Hi3559A_V100R001C02SPC020/aarch64-himix100-linux/bin/aarch64-himix100-linux-g++
    build_args="${build_args} -tags aarch64"
fi
if [ $windows -eq 1 ]; then
    export CGO_ENABLED=1
    export GOOS=windows
    export GOARCH=386
fi

INSTALLDIR=${CURRDIR}/build/package
if [ $aarm -eq 1 ]; then
    INSTALLDIR=${CURRDIR}/build/package/linux_arm64
fi
if [ $windows -eq 1 ]; then
    INSTALLDIR=${CURRDIR}/build/package/windows
fi
echo "$INSTALLDIR"

export GO111MODULE=on

echo 'Build box golang service'
for v in ${target[@]}
do
    pushd "$CURRDIR"/"$target_path"/"$v"
    go build ${build_args} -o $INSTALLDIR/$v
    if [ "$?" -eq 0 ]; then
        echo -e "\033[32m go install $v success. \033[0m"
    else
        echo -e "\033[31m go install $v failed!!! \033[0m"
        exit -1
    fi
    popd
done

echo 'finishied'
