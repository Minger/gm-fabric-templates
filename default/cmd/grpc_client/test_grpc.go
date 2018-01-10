package main

import (

    "github.com/pkg/errors"
	"github.com/rs/zerolog"

    pb "{{.PBImport}}"
)

func runTest(logger zerolog.Logger, client pb.{{.GoServiceName}}Client) error {
    return errors.New("not implemented")
}
