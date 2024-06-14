#! /usr/bin/bash

#vars
CERT="CERT-FILE LOCATION HERE"
KEY="CERT-FILE LOCATION HERE"
NAME="scrawler"

#build native version
go build -C ./src/ -o ../bin/$NAME

#illmos binary signing
if [ `uname` == "SunOS" ]; then
        elfsign \
                sign \
                -v \
                -e ./bin/$NAME \
                -c $CERT \
                -k $KEY
        #copy contents of bin folder to pkg5 folder for packaging
	if [ -d "./pkg5" ];
        then
            mkdir "./pkg5"
        fi
        cp -r ./bin/ ./pkg5/
fi


