package {{.MethodsPackageName}}

import (
	
	"github.com/rs/zerolog"
	pb "{{.PBImport}}"
)

type serverData struct{
	zerolog.Logger
}

// New{{.GoServiceName}}Server returns an object that implements the pb.{{.GoServiceName}}Server interface
func New{{.GoServiceName}}Server(logger zerolog.Logger) (pb.{{.GoServiceName}}Server, error) {
	return &serverData{
		logger,
	}, nil
}
