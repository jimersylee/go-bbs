#!/bin/sh

rm -rf ./build/go-bbs.zip

cd ./web/admin/
npm run build

zip -r ../../build/go-bbs.zip ./dist