all:
	@make build
build:
	@echo "anna: building site"
	@go build
	@./anna
serve:
	@echo "anna: serving site"
	@go build
	@./anna -s "site/"
tests:
	@echo "anna: running all tests"
	@go test ./...
bench:
	@echo "anna: running benchmark"
	@go test -bench . -benchmem -cpuprofile pprof.cpu
	@#to profile anna, run "go tool pprof app.test pprof.cpu"
clean:
	@echo "bash: purging site/rendered and test output"
	@rm -rf site/rendered
	@#cleaning test output
	@cd test/
	@rm -rf `find . -type d -name rendered`
	@cd ../
