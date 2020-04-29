BINARY = timesheet
VET_REPORT = vet.report
TEST_REPORT = tests.xml
GOARCH = amd64

VERSION?=?

# Symlink into GOPATH
GITHUB_USERNAME=praveenprem
BUILD_DIR=${GOPATH}/src/github.com/${GITHUB_USERNAME}/${BINARY}
BIN_DIR=${GOPATH}/bin
CURRENT_DIR=\$(shell pwd)
BUILD_DIR_LINK=\$(shell readlink ${BUILD_DIR})

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X main.VERSION=${VERSION}"

# Build the project
# all: link clean test vet linux darwin windows

#link:
#	BUILD_DIR=${BUILD_DIR}; \
#	BUILD_DIR_LINK=${BUILD_DIR_LINK}; \
#	CURRENT_DIR=${CURRENT_DIR}; \
#	if [ "${BUILD_DIR_LINK}" != "${CURRENT_DIR}" ]; then \
#	    echo "Fixing symlinks for build"; \
#	    rm -f ${BUILD_DIR}; \
#	    ln -s ${CURRENT_DIR} ${BUILD_DIR}; \
#	fi

linux:
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BIN_DIR}/${BINARY}-linux-${GOARCH} .

build:
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BIN_DIR}/${BINARY} .

darwin:
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BIN_DIR}/${BINARY} .

windows:
	GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-windows-${GOARCH}.exe .

install:
	install ${BIN_DIR}/${BINARY} /usr/local/bin/${BINARY}

test:
	if ! hash go2xunit 2>/dev/null; then go install github.com/tebeka/go2xunit; fi
	cd ${BUILD_DIR}; \
	godep go test -v ./... 2>&1 | go2xunit -output ${TEST_REPORT} ; \
	cd - >/dev/null

vet:
	-cd ${BUILD_DIR}; \
	godep go vet ./... > ${VET_REPORT} 2>&1 ; \
	cd - >/dev/null

fmt:
	cd ${BUILD_DIR}; \
	go fmt \$$(go list ./... | grep -v /vendor/) ; \
	cd - >/dev/null

# deps:
# 	go get gopkg.in/go-ini/ini.v1

# clean:
# 	-rm -f ${TEST_REPORT}
# 	-rm -f ${VET_REPORT}
# 	-rm -f ${BINARY}-*

.PHONY: build
