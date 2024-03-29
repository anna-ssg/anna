all:
	@make build
build:
	@echo "anna: building site"
	@go build
	@./anna
serve:
	@echo "anna: serving site"
	@go build
	@./anna -s
clean:
	@echo "bash: purging site/rendered"
	@rm -rf site/rendered
