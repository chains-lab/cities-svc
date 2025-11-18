package errx

import "github.com/chains-lab/ape"

var ErrorInviteNotFound = ape.DeclareError("INVITE_NOT_FOUND")

var ErrorInviteExpired = ape.DeclareError("INVITE_EXPIRED")

var ErrorInviteAlreadyReplied = ape.DeclareError("INVITE_ALREADY_REPLIED")

var ErrorInvalidInviteReply = ape.DeclareError("INVALID_INVITE_ANSWER")
