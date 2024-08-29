export GOPATH="$HOME/go"
export PATH="$GOPATH/bin:$PATH"

go install github.com/mitranim/gow@latest

gow run /git/main.go -e=go,mod,html,yaml,json manage