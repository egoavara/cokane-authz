package util

import "os"

func MustHostname() string {
	return Must(os.Hostname())
}

func Must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

func MustDo(err error) {
	if err != nil {
		panic(err)
	}
}
