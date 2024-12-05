
.PHONY: app
app:
	tinygo flash -size short -target=tkey ./app/blinker

.PHONY: cmd
cmd:
	go run ./cmd/tkeyled

.PHONY: test
test:
	go test ./...
