#!/bin/sh

cd /go/src/github.com/open-falcon/falcon-plus/bin;
/go/src/github.com/open-falcon/falcon-plus/bin/open-falcon start
/usr/bin/tail -f
