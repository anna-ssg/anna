install: 
	@echo "installing go bin"
	@go install github.com/acmpesuecc/anna@latest 
build:
	@echo "building site"
	@$(GOPATH)/bin/anna
serve: 
	@echo "serving site"
	@$(GOPATH)/bin/anna -s
clean: 
	@echo "cleaning site/rendered"
	@rm -rf site/rendered
all: 
	@install
	@build
