.PHONY: default
ifeq ($(OS),Windows_NT)
EXE=.exe
endif

ifeq ($(shell uname),Linux)
EXE=
endif



default:
	if [ ! -d node_modules ]; then npm install ; fi
	npm run build
	go build -o frontend$(EXE) frontend.go
