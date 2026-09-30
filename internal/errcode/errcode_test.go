package errcode

import (
	"errors"
	"testing"
)

func TestNewCustomError(t *testing.T) {
	err := NewCustomError(TokenEmpty)

	var ce *CustomError
	if !errors.As(err, &ce) {
		t.Fatalf("expect *CustomError, got %T", err)
	}
	if ce.Code != TokenEmpty {
		t.Errorf("code = %d, want %d", ce.Code, TokenEmpty)
	}
	if ce.Message != TokenEmpty.String() {
		t.Errorf("message = %q, want %q", ce.Message, TokenEmpty.String())
	}
}

func TestNewCustomErrorWithMessage(t *testing.T) {
	err := NewCustomErrorWithMessage(ErrDb, "connect refused")

	var ce *CustomError
	if !errors.As(err, &ce) {
		t.Fatalf("expect *CustomError, got %T", err)
	}
	if ce.Message != ErrDb.String()+":connect refused" {
		t.Errorf("unexpected message: %q", ce.Message)
	}
	// errors.Wrap格式为 "wrapMsg: innerErr"
	if got := err.Error(); got != "connect refused: 数据库出错" {
		t.Errorf("Error() = %q, want %q", got, "connect refused: 数据库出错")
	}
}

// 普通错误不应被误判为CustomError（FailErr依赖此行为决定是否透出错误信息）
func TestPlainErrorNotCustom(t *testing.T) {
	err := errors.New("plain")
	var ce *CustomError
	if errors.As(err, &ce) {
		t.Fatal("plain error should not match *CustomError")
	}
}
