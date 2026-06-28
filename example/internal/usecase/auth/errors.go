package auth

import (
	"errors"
)

var (
	ErrorEmailIsExisted   = errors.New("email da ton tai")
	ErrorEmailAndPassword = errors.New("tai Khoan va mat khau khong dung")
)
