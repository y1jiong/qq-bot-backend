#!/bin/bash

BIN=qq-bot-backend
OS_ARCH=linux-amd64v3

echo "> tar -xJvf ${BIN}.${OS_ARCH}.tar.xz"
tar -xJvf ${BIN}.${OS_ARCH}.tar.xz
ret=$?
if [[ $ret -ne 0 ]]; then
    echo "failed"
    exit 1
fi

echo "> chmod 700 ${BIN}"
chmod 700 ${BIN}

echo "> sudo systemctl restart ${BIN}.service"
sudo systemctl restart ${BIN}.service

echo "> rm ._*"
rm ._*

echo "> rm ${BIN}.${OS_ARCH}.tar.xz"
rm ${BIN}.${OS_ARCH}.tar.xz
