package middlewares

// Middleware 中间件
type Middleware struct {
	RequestIdKey string
	JwtSecretKey string
}

func NewMiddleware(requestIdKey, jwtSecretKey string) *Middleware {
	return &Middleware{
		RequestIdKey: requestIdKey,
		JwtSecretKey: jwtSecretKey,
	}
}
