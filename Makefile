BINARY_NAME = qq-bot-backend
CODE_FILE = ./main.go
OUTPUT_PATH = ./manifest/build
BUILD_TIME = $(shell date "+%Y-%m-%d %H:%M:%S")
COMMIT_HASH = $(shell git rev-parse HEAD)
LDFLAGS = -w -s
LDFLAGS += -X "$(BINARY_NAME)/internal/consts.BuildTime=$(BUILD_TIME)"
LDFLAGS += -X "$(BINARY_NAME)/internal/consts.CommitHash=$(COMMIT_HASH)"


.PHONY: default
default: check linux-amd64

.PHONY: all
all: check linux windows darwin

.PHONY: linux
linux: linux-amd64 linux-arm64 linux-armv7

.PHONY: windows
windows: windows-amd64 windows-arm64

.PHONY: darwin
darwin: darwin-amd64 darwin-arm64

.PHONY: check
check:
	@mkdir -p $(OUTPUT_PATH)

.PHONY: third-party-upgrade
third-party-upgrade:
	@go get -u all
	@echo $@ completed.

.PHONY: clean
clean:
	@rm -f $(OUTPUT_PATH)/*
	@echo $@ completed.

# linux
.PHONY: linux-amd64
linux-amd64:
	@rm -f $(OUTPUT_PATH)/$(BINARY_NAME).$@.tar.gz
	@CGO_ENABLE=0 GOOS=linux GOARCH=amd64 go build -ldflags '$(LDFLAGS)' -o $(OUTPUT_PATH)/$(BINARY_NAME) $(CODE_FILE)
	@tar -czvf $(OUTPUT_PATH)/$(BINARY_NAME).$@.tar.gz -C $(OUTPUT_PATH) $(BINARY_NAME) >/dev/null
	@rm -f $(OUTPUT_PATH)/$(BINARY_NAME)
	@echo $@ build completed.

.PHONY: linux-arm64
linux-arm64:
	@rm -f $(OUTPUT_PATH)/$(BINARY_NAME).$@.tar.gz
	@CGO_ENABLE=0 GOOS=linux GOARCH=arm64 go build -ldflags '$(LDFLAGS)' -o $(OUTPUT_PATH)/$(BINARY_NAME) $(CODE_FILE)
	@tar -czvf $(OUTPUT_PATH)/$(BINARY_NAME).$@.tar.gz -C $(OUTPUT_PATH) $(BINARY_NAME) >/dev/null
	@rm -f $(OUTPUT_PATH)/$(BINARY_NAME)
	@echo $@ build completed.

.PHONY: linux-armv7
linux-armv7:
	@rm -f $(OUTPUT_PATH)/$(BINARY_NAME).$@.tar.gz
	@CGO_ENABLE=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags '$(LDFLAGS)' -o $(OUTPUT_PATH)/$(BINARY_NAME) $(CODE_FILE)
	@tar -czvf $(OUTPUT_PATH)/$(BINARY_NAME).$@.tar.gz -C $(OUTPUT_PATH) $(BINARY_NAME) >/dev/null
	@rm -f $(OUTPUT_PATH)/$(BINARY_NAME)
	@echo $@ build completed.

# windows
.PHONY: windows-amd64
windows-amd64:
	@rm -f $(OUTPUT_PATH)/$(BINARY_NAME).$@.tar.gz
	@CGO_ENABLE=0 GOOS=windows GOARCH=amd64 go build -ldflags '$(LDFLAGS)' -o $(OUTPUT_PATH)/$(BINARY_NAME).exe $(CODE_FILE)
	@tar -czvf $(OUTPUT_PATH)/$(BINARY_NAME).$@.tar.gz -C $(OUTPUT_PATH) $(BINARY_NAME).exe >/dev/null
	@rm -f $(OUTPUT_PATH)/$(BINARY_NAME).exe
	@echo $@ build completed.

.PHONY: windows-arm64
windows-arm64:
	@rm -f $(OUTPUT_PATH)/$(BINARY_NAME).$@.tar.gz
	@CGO_ENABLE=0 GOOS=windows GOARCH=arm64 go build -ldflags '$(LDFLAGS)' -o $(OUTPUT_PATH)/$(BINARY_NAME).exe $(CODE_FILE)
	@tar -czvf $(OUTPUT_PATH)/$(BINARY_NAME).$@.tar.gz -C $(OUTPUT_PATH) $(BINARY_NAME).exe >/dev/null
	@rm -f $(OUTPUT_PATH)/$(BINARY_NAME).exe
	@echo $@ build completed.

# darwin
.PHONY: darwin-amd64
darwin-amd64:
	@rm -f $(OUTPUT_PATH)/$(BINARY_NAME).$@.tar.gz
	@CGO_ENABLE=0 GOOS=darwin GOARCH=amd64 go build -ldflags '$(LDFLAGS)' -o $(OUTPUT_PATH)/$(BINARY_NAME) $(CODE_FILE)
	@tar -czvf $(OUTPUT_PATH)/$(BINARY_NAME).$@.tar.gz -C $(OUTPUT_PATH) $(BINARY_NAME) >/dev/null
	@rm -f $(OUTPUT_PATH)/$(BINARY_NAME)
	@echo $@ build completed.

.PHONY: darwin-arm64
darwin-arm64:
	@rm -f $(OUTPUT_PATH)/$(BINARY_NAME).$@.tar.gz
	@CGO_ENABLE=0 GOOS=darwin GOARCH=arm64 go build -ldflags '$(LDFLAGS)' -o $(OUTPUT_PATH)/$(BINARY_NAME) $(CODE_FILE)
	@tar -czvf $(OUTPUT_PATH)/$(BINARY_NAME).$@.tar.gz -C $(OUTPUT_PATH) $(BINARY_NAME) >/dev/null
	@rm -f $(OUTPUT_PATH)/$(BINARY_NAME)
	@echo $@ build completed.
