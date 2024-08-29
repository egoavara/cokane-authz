export GOPATH="$HOME/go"
export PATH="$GOPATH/bin:$PATH"

go install github.com/mitranim/gow@latest

cd /git/cokane-authz

gow -e=go,mod,html,yaml,json run ./main.go manage