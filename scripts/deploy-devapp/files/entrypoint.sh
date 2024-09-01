apk add watchexec

sleep 1s

watchexec -w /git/cokane-authz -r --stop-signal SIGKILL "cd /git/cokane-authz; go run ./main.go manage"