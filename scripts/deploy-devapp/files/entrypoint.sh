apk add watchexec

sleep 1s

watchexec -w /git/cokane-authz -r -- "cd /git/cokane-authz; go run ./main.go manage"