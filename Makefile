.PHONY: run test

build:
	cd lib/uwu && cargo build --release
	go build -trimpath -v .

test:
	cd lib/uwu && cargo build --release
	LD_LIBRARY_PATH=${LD_LIBRARY_PATH} go test ./...
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...

prerequisites:
	sudo apt install -y cargo

# Add the env vars, including setting up LD_LIBRARY_PATH
include dev.sh

run:
	while true; do LD_LIBRARY_PATH=${LD_LIBRARY_PATH} ./multibot run; done

updatemod:
	go get -u ./...
	go mod tidy
	go mod vendor
	cd lib/uwu && cargo update
	cd lib/uwu && cargo build --release
	LD_LIBRARY_PATH=${LD_LIBRARY_PATH} go test ./...
	git add go.mod go.sum vendor lib/uwu/Cargo.lock

nilaway:
	go install go.uber.org/nilaway/cmd/nilaway@latest
	nilaway -exclude-pkgs \
		github.com/getsentry/sentry-go,github.com/pkg/sftp,google.golang.org/protobuf,google.golang.org/grpc,github.com/form3tech-oss/jwt-go,firebase.google.com/go/v4 \
		./...
