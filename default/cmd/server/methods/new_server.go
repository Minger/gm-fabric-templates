package {{.MethodsPackageName}}

import (
	pb "{{.PBImport}}"
)

type serverData struct{
}

// New{{.GoServiceName}}Server returns an object that implements the pb.{{.GoServiceName}}Server interface
func New{{.GoServiceName}}Server() (pb.{{.GoServiceName}}Server, error) {
	return &serverData{}, nil
}
