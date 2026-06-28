package auth

import (
	"errors"
)

// Lỗi đầu vào, thường là định dạng văn bản khi Request
var (
	ErrorPasswordNotValid = errors.New("password toi thieu 2 ki tu va nho hon 20 ki tu")
	ErrorEmailNotValid    = errors.New("email khong hop le")
	ErrorNameNotValid     = errors.New("name toi thieu 2 ki tu va nho hon 20 ki tu")

	ErrorEmailIsNotEmpty    = errors.New("email khong duoc trong")
	ErrorPasswordIsNotEmpty = errors.New("password khong duoc trong")
	ErrorNameIsNotEmpty     = errors.New("name khong duoc trong")
)
