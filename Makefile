PROJ := "$(notdir $(shell pwd))"
BRANCH := "$(shell git rev-parse --abbrev-ref HEAD)"
STATUS := "$(shell git status -s)"

BUILD_OUTDIR = "dist"
BUILD_FILE_PATTERN := "${PROJ}_${BRANCH}_{{.OS}}_{{.Arch}}"

BUILD_ARCH = "amd64"
BUILD_OS = "linux darwin windows"
BUILD_LDFLAGS := "-s -w -X main.version=$(BRANCH)"

TAG_REGEX = "^v[0-9]\.[0-9]\.[0-9]$$"

export GO111MODULE=on

.PHONY: test
test:
	go test ./...

.PHONY: tidy
tidy:
	@go mod tidy
	@sleep 1

.PHONY: credits
credits: tidy
	@gocredits -w
	@sleep 1

.PHONY: prepare
prepare: test tidy credits

.PHONY: build
build: test
	gox -ldflags=${BUILD_LDFLAGS} -os=${BUILD_OS} -arch=${BUILD_ARCH} -output=${BUILD_OUTDIR}/${BRANCH}/${BUILD_FILE_PATTERN}

.PHONY: release
release: checkbranch checkstatus build
	ghr "${BRANCH}" "${BUILD_OUTDIR}/${BRANCH}/"

.PHONY: checkbranch
checkbranch:
ifeq (${BRANCH}, "$(shell echo ${BRANCH} | grep ${TAG_REGEX})")
	@echo "branch name ${BRANCH} successfully checked for release"
else
	@echo "branch name ${BRANCH} does not follow semver naming convention, will not release"
	@exit 1
endif

.PHONY: checkstatus
checkstatus:
ifneq (${STATUS}, "")
	@echo "dirty branch: check git status"
	@exit 1
endif
	@:

