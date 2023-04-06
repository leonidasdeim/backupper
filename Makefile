.DEFAULT_GOAL := build

BUILD_DIR=out
TMP_DIR=tmp
FILE_NAME=backupper

.PHONY: build
build: 
	if ! [ -d ${BUILD_DIR} ] ; then mkdir ${BUILD_DIR} ; fi
	go build -o ${BUILD_DIR}/${FILE_NAME} .;

.PHONY: test
test: 
	go test -v -count=1 ./...;

.PHONY: clean
clean:
	if [ -d ${BUILD_DIR} ] ; then rm -rf ${BUILD_DIR} ; fi
	if [ -d ${TMP_DIR} ] ; then rm -rf ${TMP_DIR} ; fi


######################################################
# jobs bellow are for demo purposes only
######################################################

.PHONY: run_demo
run_demo: 
	if ! [ -d ${TMP_DIR} ] ; then mkdir ${TMP_DIR} ; fi
	if ! [ -d ${TMP_DIR}/hot ] ; then mkdir ${TMP_DIR}/hot ; fi
	cd ${TMP_DIR};\
	../${BUILD_DIR}/${FILE_NAME} --in=hot --out=backup;

.PHONY: build_run_demo
build_run_demo: build run_demo

.PHONY: second_run_demo
second_run_demo: 
	cd ${TMP_DIR};\
	../${BUILD_DIR}/${FILE_NAME};

