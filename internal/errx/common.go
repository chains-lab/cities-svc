package errx

import (
	"github.com/chains-lab/ape"
)

var ErrorInternal = ape.DeclareError("INTERNAL_ERROR")

var ErrorPermissionDenied = ape.DeclareError("PERMISSION_DENIED")

var ErrorUnauthenticated = ape.DeclareError("UNAUTHENTICATED")

var ErrorInvalidArgument = ape.DeclareError("INVALID_ARGUMENT")

var ErrorNotFound = ape.DeclareError("NOT_FOUND")
