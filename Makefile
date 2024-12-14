
.PHONY: blinker-app
blinker-app:
	tinygo flash -size short -target=tkey ./examples/blinker/app

.PHONY: blinker-cmd
blinker-cmd:
	go run ./examples/blinker/cmd

.PHONY: test
test:
	go test ./...
