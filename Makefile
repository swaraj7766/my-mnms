.PHONY: default
ifeq ($(OS),Windows_NT)
EXE=.exe
endif

ifeq ($(shell uname),Linux)
EXE=
endif

default:
	go build -o mnmsctl/mnmsctl$(EXE) mnmsctl/main.go && go vet ./...
	
update:
	go env -w GOPRIVATE=github.com/Atop-NMS-team/*
	git config --global url."git@github.com:".insteadOf "https://github.com/"
	go get -u github.com/Atop-NMS-team/simulator
