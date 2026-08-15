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

func (u *UID) ToBase58(bits int64) string {
	uid := uint64(u.LocalID)<<bits | uint64(u.ObjectID)
	buf := make([]byte, 0, 20)
	buf = fmt.Appendf(buf, "%d", uid)

	return base58.Encode(buf)
}

func DecodeFromBase58(fakeUID string, bits int64) *UID {
	mask := uint64((1 << bits) - 1)
	uidByte := base58.Decode(fakeUID)
	uid, _ := strconv.ParseUint(string(uidByte), 10, 64)
	return &UID{
		LocalID:  uint32(uid >> bits),
		ObjectID: uint(uid & mask),
	}
}
