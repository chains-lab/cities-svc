package interceptors

import (
	"context"
	"fmt"

	"github.com/chains-lab/gatekit/auth"
	"github.com/chains-lab/gatekit/roles"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type MetaData struct {
	Issuer         string     `json:"iss,omitempty"`
	Subject        string     `json:"sub,omitempty"`
	Audience       []string   `json:"aud,omitempty"`
	InitiatorID    uuid.UUID  `json:"initiator_id,omitempty"`
	SessionID      uuid.UUID  `json:"session_id,omitempty"`
	SubscriptionID uuid.UUID  `json:"subscription_id,omitempty"`
	Verified       bool       `json:"verified,omitempty"`
	Role           roles.Role `json:"role,omitempty"`
	RequestID      uuid.UUID  `json:"request_id,omitempty"`
}

func NewAuth(skService, skUser string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.New(codes.Unauthenticated, fmt.Sprintf("no metadata found in incoming context")).Err()
		}
		toksServ := md["authorization"]
		if len(toksServ) == 0 {
			return nil, status.New(codes.Unauthenticated, fmt.Sprintf("authorization token not supplied")).Err()
		}

		data, err := auth.VerifyServiceJWT(ctx, toksServ[0], skService)
		if err != nil {
			return nil, status.New(codes.Unauthenticated, fmt.Sprintf("failed to verify token: %s", err)).Err()
		}

		toksUser := md["x-user-token"]
		if len(toksUser) == 0 {
			return nil, status.New(codes.Unauthenticated, fmt.Sprintf("user token not supplied")).Err()
		}

		requestIDArr := md["x-request-id"]
		if len(requestIDArr) == 0 {
			return nil, status.New(codes.Unauthenticated, fmt.Sprintf("request ID not supplied")).Err()
		}

		userData, err := auth.VerifyUserJWT(ctx, toksUser[0], skUser)
		if err != nil {
			return nil, status.New(codes.Unauthenticated, fmt.Sprintf("invalid user token: %v", err)).Err()
		}

		userID, err := uuid.Parse(userData.Subject)
		if err != nil {
			return nil, status.New(codes.Unauthenticated, fmt.Sprintf("invalid user ID: %v", err)).Err()
		}

		requestID, err := uuid.Parse(requestIDArr[0])
		if err != nil {
			return nil, status.New(codes.Unauthenticated, fmt.Sprintf("invalid request ID: %v", err)).Err()
		}

		ctx = context.WithValue(ctx, MetaCtxKey, MetaData{
			Issuer:      data.Issuer,
			Subject:     data.Subject,
			Audience:    data.Audience,
			InitiatorID: userID,
			SessionID:   userData.Session,
			Verified:    userData.Verified,
			Role:        userData.Role,
			RequestID:   requestID,
		})

		return handler(ctx, req)
	}
}
