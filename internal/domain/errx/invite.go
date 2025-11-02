package errx

import "github.com/chains-lab/ape"

var ErrorInvalidInviteAnswer = ape.DeclareError("INVITE_INVALID_ANSWER")

var ErrorInviteNotFound = ape.DeclareError("INVITE_NOT_FOUND")

var ErrorInviteAlreadyAnswered = ape.DeclareError("INVITE_ALREADY_ANSWERED")

var ErrorInviteExpired = ape.DeclareError("INVITE_EXPIRED")

var ErrorUserAlreadyCityAdmin = ape.DeclareError("USER_ALREADY_CITY_ADMIN")

var ErrorCityAdminNotAllowed = ape.DeclareError("CITY_ADMIN_NOT_ALLOWED")
