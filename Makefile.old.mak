#

tname = .+
# timestamp = $(shell date -u '+%Y-%m-%dT%H:%M:%S')
timestamp = $(shell date -u '+%Y-%m-%dT%H-%M-%S')
LOGS_DIR = $(shell test -d logs || mkdir -pv logs)
COUNT_TEST ?= 0
COUNT_COV ?= 0
COUNT_TEST_FILE = ./logs/count-test
COUNT_COV_FILE = ./logs/count-cov


# all: test
all: help




## test: make testing
test: test-slog

## cov/cover/caverage: make coverage testing
cov: cover
coverage: cover
cover: cover-slog

## bench: make benchmark testing
bench: bench-top-level

bench-top-level:
	@-$(MAKE) -s bench-meta package=./bench btitle=top-level bname=.+

test-slog:
	@-$(MAKE) -s test-meta tname=$(tname) ttitle=slog package=./slog/...

cover-slog:
	@-$(MAKE) -s cover-meta tname=$(tname) ttitle=slog package=./slog/...

cover-diff:
	@-$(MAKE) -s cover-meta tname=$(tname) seq=-1 ttitle=slog package=./slog/...
	@-$(MAKE) -s cover-meta tname=$(tname) seq=-2 ttitle=slog package=./slog/...
	@which benchstat >/dev/null || go install golang.org/x/perf/cmd/benchstat@latest
	@benchstat ./logs/cover-slog-$(timestamp).log ./logs/cover-slog-$(timestamp).log

test-meta: inc-counter-test | ./logs
	@-$(MAKE) -s test-meta-1 COUNT_TEST=$(shell cat $(COUNT_TEST_FILE))

test-meta-1: | ./logs
	echo go test -v -test.v -test.run '^$(tname)$$' $(package) 2>&1 '>>>' tee ./logs/test-$(ttitle)-$(COUNT_TEST).log
	go test -v -test.v -test.run '^$(tname)$$' $(package) 2>&1 | tee ./logs/test-$(ttitle)-$(COUNT_TEST).log
	@echo "test-$(ttitle)-$(COUNT_TEST)" end.

cover-meta: inc-counter-cov | ./logs
	@$(MAKE) -s cover-meta-1 COUNT_COV=$(shell cat $(COUNT_COV_FILE))

cover-meta-1:
	echo go test -v -cover -test.v -test.run '^$(tname)$$' $(package) -race -coverprofile=./logs/coverage$(seq).txt -covermode=atomic -timeout=20m -test.short -vet=off 2>&1 | tee ./logs/cover-$(ttitle)-$(COUNT_COV)$(seq).log
	go test -v -cover -test.v -test.run '^$(tname)$$' $(package) -race -coverprofile=./logs/coverage$(seq).txt -covermode=atomic -timeout=20m -test.short -vet=off 2>&1 | tee ./logs/cover-$(ttitle)-$(COUNT_COV)$(seq).log
	echo Generating cover.html...
	go tool cover -html=./logs/coverage$(seq).txt -o ./logs/cover.html 2>&1
	# open ./logs/cover.html
	@echo coverage "cover-$(ttitle)-$(COUNT_COV)$(seq)" end.

bench-meta: inc-counter-test | ./logs
	@-$(MAKE) -s bench-meta-1 COUNT_TEST=$(shell cat $(COUNT_TEST_FILE))

bench-meta-1:
	echo go test -v -bench=$(package) -benchmem -run='^$(bname)$$'
	go test -v $(package) -bench=. -benchmem -run='^$(bname)$$' | tee ./logs/bench-$(btitle)-$(COUNT_TEST).log
	@echo "bench-$(ttitle)-$(COUNT_TEST)" end.




./logs:
	mkdir -pv $@
	@# touch $@/count-{test,cov,bench}

inc-counter-test:
	@if ! test -f $(COUNT_TEST_FILE); then echo 0 > $(COUNT_TEST_FILE); fi
	@echo $$(( $(shell cat $(COUNT_TEST_FILE))+1 )) > $(COUNT_TEST_FILE)

inc-counter-cov:
	@if ! test -f $(COUNT_COV_FILE); then echo 0 > $(COUNT_COV_FILE); fi
	@echo $$(( $(shell cat $(COUNT_COV_FILE))+1 )) > $(COUNT_COV_FILE)





.PHONY: printvars info help all
printvars:
	$(foreach V, $(sort $(filter-out .VARIABLES,$(.VARIABLES))), $(info $(v) = $($(v))) )
	# Simple:
	#   (foreach v, $(filter-out .VARIABLES,$(.VARIABLES)), $(info $(v) = $($(v))) )

print-%:
	@echo $* = $($*)

info:
	@echo "     GO_VERSION: $(GOVERSION)"
	@echo "        GOPROXY: $(GOPROXY)"
	@echo "         GOROOT: $(shell go env GOROOT) | GOPATH: $(shell go env GOPATH)"
	@echo "    GO111MODULE: $(GO111MODULE)"
	@echo
	@echo "         GOBASE: $(GOBASE)"
	@echo "          GOBIN: $(GOBIN)"
	@echo "    PROJECTNAME: $(PROJECTNAME)"
	@echo "        APPNAME: $(APPNAME)"
	@echo "        VERSION: $(VERSION)"
	@echo "      BUILDTIME: $(TIMESTAMP)"
	@echo "    GIT_VERSION: $(GIT_VERSION)"
	@echo "   GIT_REVISION: $(GIT_REVISION)"
	@echo "        GIT_HASH: $(GIT_HASH)"
	@echo "    GIT_SUMMARY: $(GIT_SUMMARY)"
	@echo "       GIT_DESC: $(GIT_DESC)"
	@echo
	@echo "             OS: $(OS)"
	@echo
	@echo " MAIN_BUILD_PKG: $(MAIN_BUILD_PKG)"
	@echo "      MAIN_APPS: $(MAIN_APPS)"
	@echo "       SUB_APPS: $(SUB_APPS)"
	@echo "MAIN_ENTRY_FILE: $(MAIN_ENTRY_FILE)"
	@echo
	#@echo "export GO111MODULE=on"
	@echo "export GOPROXY=$(shell go env GOPROXY)"
	#@echo "export GOPATH=$(shell go env GOPATH)"
	@echo

.PHONY: help
# all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo