.PHONY: default
ifeq ($(OS),Windows_NT)
EXE=.exe
endif

ifeq ($(shell uname),Linux)
EXE=
endif

.PHONY: frontend

default:
	make -C mnmsctl

lint:
	go vet 
	revive
	golangci-lint run
	staticcheck

update:
	go env -w GOPRIVATE=github.com/Atop-NMS-team/*
	git config --global url."git@github.com:".insteadOf "https://github.com/"
	go get -u github.com/Atop-NMS-team/simulator

frontend:
	make -C frontend

.PHONY: release
release: 
	mkdir -p release
	GOOS=linux go build -ldflags "-X 'main.Version=v1.0.1' -s -w" -o release/mnmsctl_linux_amd64 mnmsctl/*.go
	GOOS=linux go build -ldflags "-X 'main.Version=v1.0.1' -s -w" -o release/frontend_linux_amd64 frontend/frontend.go
	GOOS=windows go build -ldflags "-X 'main.Version=v1.0.1' -s -w" -o release/mnmsctl_windows_amd64.exe mnmsctl/*.go
	GOOS=windows go build -ldflags "-X 'main.Version=v1.0.1' -s -w" -o release/frontend_windows_amd64.exe frontend/frontend.go
	cp doc/README.md release
	#go doc -all > release/godoc.txt
	cp doc/authentication.md release
	cd release; rm mnms.zip; zip mnms.zip *
