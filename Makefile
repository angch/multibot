build:
	cd lib/uwu && cargo build --release
	go build .

prerequisites:
	sudo apt install -y cargo

# Add the env vars, including setting up LD_LIBRARY_PATH
include dev.sh

run:
	LD_LIBRARY_PATH=${LD_LIBRARY_PATH} ./discordbot run

