

CGO_CPPFLAGS ?= ${CPPFLAGS}
export CGO_CPPFLAGS
CGO_CFLAGS ?= ${CFLAGS}
export CGO_CFLAGS
CGO_LDFLAGS ?= $(filter -g -L% -l% -O%,${LDFLAGS})
export CGO_LDFLAGS

app = simulator
udp_port=55954


PROJECT=$(shell basename  $(CURDIR))


.PHONY: clean   bin/$(PROJECT)

bin/$(PROJECT): 
ifeq ($(OS),Windows_NT)
	go build -o ./bin/$(PROJECT).exe ./cmd/$(PROJECT)/main.go
else
	go build -o ./bin/$(PROJECT) ./cmd/$(PROJECT)/main.go
endif
run: build
	${app} run

clean:
	@go clean -x
	rm -rf ./bin/* $(PROJECT)
	rm -rf  ${app}


install: bin/$(PROJECT)
	mv ./bin/* ../bin



help:
	@echo Usage:
	@echo "make  <target>"
	@echo Targets:	
	@echo bin/$(PROJECT): 	build bin/$(PROJECT)
	@echo clean:	 		clean test data
