package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

// エラーコード定数
const (
	ErrCodeBadRequest        = "BAD_REQUEST"
	ErrCodeNotFound          = "NOT_FOUND"
	ErrCodeUnauthorized      = "UNAUTHORIZED"
	ErrCodeForbidden         = "FORBIDDEN"
	ErrCodeConflict          = "CONFLICT"
	ErrCodeInternalError     = "INTERNAL_ERROR"
	ErrCodeDatabaseError     = "DATABASE_ERROR"
	ErrCodeValidationError   = "VALIDATION_ERROR"
	ErrCodeGroupAccessDenied = "GROUP_ACCESS_DENIED"
)

// AppError はアプリケーション全体で使用するカスタムエラー型
type AppError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Status  int                    `json:"-"`
	Details map[string]interface{} `json:"details,omitempty"`
	Err     error                  `json:"-"` // 元のエラー(エラーラップ対応)
}

// Error は error インターフェースの実装
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap は errors.Unwrap 対応
func (e *AppError) Unwrap() error {
	return e.Err
}

// WithDetails はエラーに詳細情報を追加する
func (e *AppError) WithDetails(key string, value interface{}) *AppError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}
	e.Details[key] = value
	return e
}

// WithErr は元のエラーを追加する
func (e *AppError) WithErr(err error) *AppError {
	e.Err = err
	return e
}

// BadRequestError は不正なリクエストエラーを返す
func BadRequestError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeBadRequest,
		Message: message,
		Status:  http.StatusBadRequest,
	}
}

// NotFoundError はリソースが見つからないエラーを返す
func NotFoundError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeNotFound,
		Message: message,
		Status:  http.StatusNotFound,
	}
}

// UnauthorizedError は認証エラーを返す
func UnauthorizedError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeUnauthorized,
		Message: message,
		Status:  http.StatusUnauthorized,
	}
}

// ForbiddenError は権限エラーを返す
func ForbiddenError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeForbidden,
		Message: message,
		Status:  http.StatusForbidden,
	}
}

// ConflictError は競合エラーを返す
func ConflictError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeConflict,
		Message: message,
		Status:  http.StatusConflict,
	}
}

// InternalServerError は内部サーバーエラーを返す
func InternalServerError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeInternalError,
		Message: message,
		Status:  http.StatusInternalServerError,
	}
}

// DatabaseError はデータベースエラーを返す
// err パラメータで元のエラーをラップする
func DatabaseError(message string, err error) *AppError {
	return &AppError{
		Code:    ErrCodeDatabaseError,
		Message: message,
		Status:  http.StatusInternalServerError,
		Err:     err,
	}
}

// GroupAccessDeniedError はグループアクセス拒否エラーを返す
func GroupAccessDeniedError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeGroupAccessDenied,
		Message: message,
		Status:  http.StatusForbidden,
	}
}

// ValidationError はフィールド別のバリデーションエラーを扱う
type ValidationError struct {
	*AppError
	Fields map[string]string `json:"fields,omitempty"` // フィールド名 -> エラーメッセージ
}

// Error は error インターフェースの実装
func (e *ValidationError) Error() string {
	return e.AppError.Error()
}

// NewValidationError はバリデーションエラーを生成する
func NewValidationError(message string) *ValidationError {
	return &ValidationError{
		AppError: &AppError{
			Code:    ErrCodeValidationError,
			Message: message,
			Status:  http.StatusBadRequest,
		},
		Fields: make(map[string]string),
	}
}

// AddField はフィールドエラーを追加する
func (e *ValidationError) AddField(field, message string) *ValidationError {
	e.Fields[field] = message
	return e
}

// HasErrors はフィールドエラーが存在するかチェックする
func (e *ValidationError) HasErrors() bool {
	return len(e.Fields) > 0
}

// As はerrors.As()のラッパー関数
func As(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return nil
}
