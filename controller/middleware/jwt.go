package middleware

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aalexanderkevin/crypto-wallet/helper"

	"github.com/golang-jwt/jwt/v5"
	"github.com/segmentio/ksuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type JWTData struct {
	jwt.RegisteredClaims
	Email string `json:"email"`
}

func JWTMiddleware(secretKey string, excludedMethods []string) func(context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Check if the method is in the excluded list.
		for _, method := range excludedMethods {
			if method == info.FullMethod {
				// Skip token validation for excluded methods.
				return handler(ctx, req)
			}
		}

		token, err := getTokenAuth(ctx)
		if err != nil {
			return nil, err
		}
		claim, err := decodeJwtData(secretKey, token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, err.Error())
		}
		if claim.Email == "" {
			return nil, status.Errorf(codes.Unauthenticated, "missing email")
		}

		// Set the email in the gRPC context.
		ctx = context.WithValue(ctx, helper.ContextKeyJwtData, claim.Email)

		return handler(ctx, req)
	}
}

func getTokenAuth(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("token not found in metadata")
	}

	token := md.Get("authorization")
	if len(token) == 0 {
		return "", errors.New("token not found in authorization header")
	}

	s := strings.SplitN(token[0], " ", 2)
	if len(s) != 2 || strings.ToLower(s[0]) != "bearer" {
		return "", errors.New("invalid token format")
	}
	//Use authorization header token only if token type is bearer else query string access token would be returned
	if len(s) > 0 && strings.ToLower(s[0]) == "bearer" {
		return s[1], nil
	}

	// authHeader := strings.ToLower(token[0])
	// splitToken := strings.Split(authHeader, "bearer ")
	// if len(splitToken) != 2 {
	return "", errors.New("malformed token")
	// }

	// return splitToken[1], nil
}

func decodeJwtData(secret string, tokenStr string) (*JWTData, error) {
	var claim JWTData

	secretFn := func(token *jwt.Token) (interface{}, error) {
		if _, validSignMethod := token.Method.(*jwt.SigningMethodHMAC); !validSignMethod {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	}

	token, err := jwt.ParseWithClaims(tokenStr, &claim, secretFn)
	if err != nil {
		return nil, err
	}

	if claim, ok := token.Claims.(*JWTData); ok && token.Valid {
		return claim, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func GenerateJwt(email string, secretKey string) (*string, error) {
	claims := jwt.RegisteredClaims{
		ID:        ksuid.New().String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		Subject:   email,
	}

	// generate token
	accessClaims := JWTData{
		RegisteredClaims: claims,
		Email:            email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	var byteSecret = []byte(secretKey)
	accessToken, err := token.SignedString(byteSecret)
	if err != nil {
		return nil, err
	}

	return &accessToken, nil
}

func GetJWTData(ctx context.Context) string {
	claim := ctx.Value(helper.ContextKeyJwtData)
	return claim.(string)
}
