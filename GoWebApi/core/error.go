package core

// Error 只用于 swagger doc 生成
// Deprecated: 是 fiber.Error 的副本
type Error struct {
	Message string `json:"message" validate:"required"`
	Code    int    `json:"code" validate:"required"`
}

func (e *Error) Error() string {
	return e.Message
}
