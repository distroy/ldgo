#
# Copyright (C) distroy
#

# variables
PROJECT_ROOT= $(patsubst %/,%,$(abspath $(dir $$PWD)))
PROTOS=$(dir $(shell test -d "$(PROJECT_ROOT)/proto" && find "$(PROJECT_ROOT)/proto" -name '*.proto'))
$(info PROJECT_ROOT: $(PROJECT_ROOT))
$(info PROTOS: $(PROTOS))

# go
# GO=env GODEBUG=madvdontneed=1 go
GO=go
GO_FLAGS=${flags}
GO_VERSION=$(shell go version | cut -d" " -f 3)
GO_MAJOR_VERSION=$(shell echo $(GO_VERSION) | cut -d"." -f 1)
GO_SUB_VERSION=$(shell echo $(GO_VERSION) | cut -d"." -f 2)
ifeq ($(shell expr ${GO_SUB_VERSION} '>' 22), 1)
	GO_FLAGS+=-ldflags=-checklinkname=0
endif
$(info GO_VERSION: $(GO_MAJOR_VERSION).$(GO_SUB_VERSION))
$(info GO_FLAGS: $(GO_FLAGS))

# go test
GO_TEST_DIRS+=$(shell find . -name '*_test.go' | grep -v -E 'vendor|bak' | xargs dirname | sort | uniq)
GO_TEST_DIRS_NAME=$(notdir $(GO_TEST_DIRS))
# $(info GO_TEST_DIRS: $(GO_TEST_DIRS_NAME))

ifeq (${test_report},)
	export test_report=$(PROJECT_ROOT)/log
endif
GO_TEST_FLAGS+=-v
GO_TEST_FLAGS+=-gcflags="all=-l"
GO_TEST_REPORT_DIR=${test_report}

# git
GIT_REVISION=$(shell git rev-parse HEAD 2> /dev/null)
GIT_BRANCH=$(shell git symbolic-ref HEAD 2> /dev/null | sed -e 's/refs\/heads\///')
GIT_TAG=$(shell git describe --exact-match --tags 2> /dev/null)
$(info GIT_REVISION: $(GIT_REVISION))
$(info GIT_BRANCH: $(GIT_BRANCH))
$(info GIT_TAG: $(GIT_TAG))

_mk_protobuf = ( \
	echo "=== building protobuf: $(1)"; \
	cd $(1); \
	rm -rf *.pb.go *_pb2.py; \
	echo protoc --go_out . --python_out . *.proto; \
	protoc --go_out . --python_out . ./*.proto || exit $$?; \
	cd $(PROJECT_ROOT); \
	);

_go_install =  ( \
	_work_dir="$$PWD"; \
	_go_install_dir="$$(mktemp -d)"; \
	cd "$$_go_install_dir"; \
	go mod init main 1>/dev/null 2>/dev/null; \
	echo "go install $(1)@latest || go install $(1)"; \
	go install $(1)@latest || go install $(1); \
	cd "$$_work_dir"; \
	rm -rf "$$_go_install_dir"; \
	);

.PHONY: all
all: go-test

.PHONY: $(GO_TEST_DIRS_NAME)
$(GO_TEST_DIRS_NAME):
	@echo GO_TEST_DIRS: $(notdir $@)
	$(GO) test $(GO_FLAGS) $(GO_TEST_FLAGS) ./$(notdir $@) \
		-coverprofile="$(GO_TEST_REPORT_DIR)/go-coverage.out"

.PHONY: pb
pb:
	@$(foreach i, $(PROTOS), $(call _mk_protobuf,$(i)))

.PHONY: dep
dep:
	$(GO) mod tidy
	# $(GO) mod vendor

.PHONY: go-test-report-dir
go-test-report-dir:
	mkdir -pv $(GO_TEST_REPORT_DIR)

.PHONY: go-test
go-test: go-test-report-dir
	$(GO) test $(GO_FLAGS) $(GO_TEST_FLAGS) ./... \
		-coverprofile="$(GO_TEST_REPORT_DIR)/go-coverage.out"
	$(GO) tool cover -html $(GO_TEST_REPORT_DIR)/go-coverage.out \
		-o $(GO_TEST_REPORT_DIR)/go-coverage.html

go-test-report: go-test-report-dir
	$(GO) test $(GO_FLAGS) $(GO_TEST_FLAGS) ./... \
		-coverprofile="$(GO_TEST_REPORT_DIR)/go-coverage.out" \
		-json > "$(GO_TEST_REPORT_DIR)/go-test.json"
	$(GO) tool cover -html $(GO_TEST_REPORT_DIR)/go-coverage.out \
		-o $(GO_TEST_REPORT_DIR)/go-coverage.html

.PHONY: setup
setup:
	git submodule init
	git submodule update
	git config core.hooksPath "git-go-tool/git-hook"
	@$(call _go_install,github.com/distroy/git-go-tool/cmd/git-diff-go-cognitive)
	@$(call _go_install,github.com/distroy/git-go-tool/cmd/git-diff-go-coverage)
	@$(call _go_install,github.com/distroy/git-go-tool/cmd/git-diff-go-format)
	@$(call _go_install,github.com/distroy/git-go-tool/cmd/go-cognitive)
	@$(call _go_install,github.com/distroy/git-go-tool/cmd/go-format)
	@echo $$'\E[32;1m'"setup succ"$$'\E[0m'

.PHONY: cognitive
cognitive: setup
	go-cognitive

.PHONY: format
format: setup
	go-format --func-input-num 4 -func-context-error-match=1
