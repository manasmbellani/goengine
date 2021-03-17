#!/bin/bash
cd /opt/goengine/src
go build -o goengine * 
mv goengine ~/go/bin/goengine