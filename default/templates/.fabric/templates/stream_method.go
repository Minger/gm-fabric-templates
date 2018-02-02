package {{"{{.MethodsPackageName}}"}}

import (
    "github.com/pkg/errors"

	pb "{{"{{.PBImport}}"}}"
)

{{"{{.Comments}}"}}
func (s *serverData) {{"{{.MethodDeclaration}}"}} {
	return errors.New("not implemented")
}

