generate-test: embed
	@go generate ./test
	@echo "[OK] Files for test generated!"

embed:
	@go generate ./internal/box
	@echo "[OK] Files added to embed box!"
