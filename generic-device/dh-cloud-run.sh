#!/bin/bash

script_dir=$(pwd)

if [ ! -f $GOPATH/bin/devicehive-cloud ]; then
	echo "devicehive-cloud package not found. Installing it now..."
	go get github.com/devicehive/IoT-framework/devicehive-cloud
	cd $GOPATH/src/github.com/devicehive/IoT-framework/devicehive-cloud
	go install
	echo "Installation of devicehive-cloud complete. Running application..."
fi

cd $script_dir
$GOPATH/bin/devicehive-cloud -conf=deviceconf.yml
