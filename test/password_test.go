package test

import (
	"fmt"
	"testing"
	"todo-app/pkg/utils"
)

func TestHashPassword(t *testing.T) {
	// Test với mật khẩu hợp lệ
	password := "secretPassword"
	hashedPassword, err := utils. HashPassword(password)
	fmt.Println(hashedPassword)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Kiểm tra nếu mật khẩu đã băm không rỗng
	if hashedPassword == "" {
		t.Fatalf("Expected hashed password to be non-empty")
	}
}

func TestComparePassword(t *testing.T) {
	// Test với mật khẩu hợp lệ
	password := "secretPassword"
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Kiểm tra nếu mật khẩu khớp
	match, err := utils.ComparePassword(hashedPassword, password)
	fmt.Println(match)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !match {
		t.Fatalf("Expected password to match")
	}

	// Kiểm tra nếu mật khẩu không khớp
	wrongPassword := "wrongPassword"
	match, err = utils.ComparePassword(hashedPassword, wrongPassword)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
	if match {
		t.Fatalf("Expected password to not match")
	}
}
