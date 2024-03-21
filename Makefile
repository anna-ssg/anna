all: 
	@make install
	@make build
install: 
	@echo "go: installing anna"
	@go install github.com/acmpesuecc/anna@latest 
build:
	@echo "anna: building site"
	@$(GOPATH)/bin/anna
serve: 
	@echo "anna: serving site"
	@$(GOPATH)/bin/anna -s
clean: 
	@echo "bash: purging site/rendered"
	@rm -rf site/rendered
