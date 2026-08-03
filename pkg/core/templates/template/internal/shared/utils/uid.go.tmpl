package utils

import (
	"fmt"
	"strconv"

	"github.com/btcsuite/btcutil/base58"
)

type UID struct {
	LocalID  uint32
	ObjectID uint
}

func NewUID(localID uint32, objectID uint) *UID {
	return &UID{
		LocalID:  localID,
		ObjectID: objectID,
	}
}

func (u *UID) ToBase58() string {
	// Dùng fmt.Appendf thay vì fmt.Sprintf + []byte
	uid := uint64(u.LocalID)<<26 | uint64(u.ObjectID)

	// Tạo buffer với capacity đủ lớn để tránh reallocate
	buf := make([]byte, 0, 20) // uint64 max ~20 digits
	buf = fmt.Appendf(buf, "%d", uid)

	return base58.Encode(buf)
}

func DecodeFromBase58(fakeUID string) *UID {
	uidByte := base58.Decode(fakeUID)
	uid, _ := strconv.ParseUint(string(uidByte), 10, 64)
	return &UID{
		LocalID:  uint32(uid >> 26),
		ObjectID: uint(uid & 0x2ffffff),
	}
}
