package config

import "unique"

type metadata struct {
	Version string
}

var Metadata = unique.Make(metadata{
	Version: "v1.0.0",
})
