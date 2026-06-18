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
	// dich chuyen
	uid := uint64(u.LocalID)<<27 | uint64(u.ObjectID)
	uidNew := base58.Encode([]byte(fmt.Sprintf("%d", uid)))
	return uidNew
}

func DecodeFromBase58(fakeUID string) *UID {
	uidByte := base58.Decode(fakeUID)
	uid, _ := strconv.ParseUint(string(uidByte), 10, 64)
	return &UID{
		LocalID:  uint32(uid >> 27),
		ObjectID: uint(uid & 0x3ffffff),
	}
}
