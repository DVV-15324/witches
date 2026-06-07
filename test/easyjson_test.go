package test

import (
	easyjson "core-v/pkg/core/easyjson"
	"encoding/json"
	easyjsonpkg "github.com/mailru/easyjson"
	"testing"
)

var user = easyjson.User{
	ID:    1,
	Name:  "Vu",
	Email: "vu@example.com",
	Age:   22,
}

func BenchmarkReflectionMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(user)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEasyJSONMarshal(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := easyjsonpkg.Marshal(user)
		if err != nil {
			b.Fatal(err)
		}
	}
}
