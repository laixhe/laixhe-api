package middlewares

import (
	contribJwt "github.com/gofiber/contrib/v3/jwt"
	"github.com/laixhe/gonet/xfiber"
)

type Middleware struct {
	UseJwtConfig     contribJwt.Config
	UseJwtConfigNext contribJwt.Config
}

func NewMiddleware(jwtSecretKey string) *Middleware {
	return &Middleware{
		UseJwtConfig: contribJwt.Config{
			SigningKey: contribJwt.SigningKey{Key: []byte(jwtSecretKey)},
			Claims:     &JwtClaims{},
		},
		UseJwtConfigNext: contribJwt.Config{
			ErrorHandler: xfiber.JwtErrorHandlerNext,
			SigningKey:   contribJwt.SigningKey{Key: []byte(jwtSecretKey)},
			Claims:       &JwtClaims{},
		},
	}
}
