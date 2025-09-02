package errx

import "github.com/chains-lab/ape"

var ErrorInviteNotFound = ape.DeclareError("INVITE_NOT_FOUND")

var ErrorInvalidInviteToken = ape.DeclareError("INVALID_INVITE_TOKEN")

var ErrorInviteAlreadyAnswered = ape.DeclareError("INVITE_ALREADY_ANSWERED")

var ErrorInviteExpired = ape.DeclareError("INVITE_EXPIRED")

var ErrorInviteStatusInvalid = ape.DeclareError("INVITE_STATUS_INVALID")
