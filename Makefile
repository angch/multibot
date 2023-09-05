.PHONY: run test

build:
	cd lib/uwu && cargo build --release
	go build -trimpath .

test:
	cd lib/uwu && cargo build --release
	LD_LIBRARY_PATH=${LD_LIBRARY_PATH} go test ./...

prerequisites:
	sudo apt install -y cargo

# Add the env vars, including setting up LD_LIBRARY_PATH
include dev.sh

run:
	while true; do LD_LIBRARY_PATH=${LD_LIBRARY_PATH} ./multibot run; done

