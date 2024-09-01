apk add watchexec

sleep 1s

watchexec -w /git/cokane-authz -r --stop-signal SIGKILL "go run /git/cokane-authz/main.go manage"