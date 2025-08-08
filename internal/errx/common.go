package errx

import (
	"github.com/chains-lab/cities-dir-svc/internal/errx/statusx"
	"github.com/chains-lab/gatekit/roles"
	"github.com/chains-lab/svc-errors/ape"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var ErrorInternal = ape.Declare("INTERNAL_ERROR")

func RaiseInternal(cause error) error {
	return ErrorInternal.Raise(
		cause,
		status.New(codes.Internal, "internal server error"),
	)
}

var ErrorRoleIsNotApplicable = ape.Declare("ROLE_IS_NOT_APPLICABLE")

func RaiseRoleIsNotApplicable(cause error, userID uuid.UUID, role roles.Role) error {
	return ErrorRoleIsNotApplicable.Raise(
		cause,
		statusx.PermissionDeniedByRole(userID, role),
	)
}
