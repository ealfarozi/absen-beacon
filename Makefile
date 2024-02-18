local-build:
	rm go.mod;rm go.sum;rm -rf vendor/;go mod init github.com/ealfarozi/absen-beacon;go mod tidy;go mod download;go mod vendor;go mod verify
pull:
	rm go.mod;rm go.sum;rm -rf vendor/;git pull
pull-run:
	rm go.mod;rm go.sum;rm -rf vendor/;git pull;go mod init github.com/ealfarozi/absen-beacon;go mod tidy;go mod download;go mod vendor;go mod verify;go run main.go