#! /usr/bin/bash

#vars
CERT="CERT LOCATION"
KEY="KEY LOCATION"

#build native version
go build -C ./src/ -o ../bin/svcbundle

#illmos binary signing
if [ `uname` == "SunOS" ]; then
        elfsign \
                sign \
                -v \
                -e ./bin/svcbundle \
                -c $CERT \
                -k $KEY
        #copy contents of bin folder to pkg5 folder for packaging
	if [ -d "./pkg5" ];
        then
            mkdir "./pkg5"
        fi
        cp -r ./bin/ ./pkg5/
fi


