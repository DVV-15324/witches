package response

import (
	"testing"
	"time"
)

func TestReponse(t *testing.T) {
	data := map[string]interface{}{
		"name": "vu",
	}
	testReponse := NewEntityResponse(200, "Thanh Cong", time.Now(), data)
	testReponse.PrintReponse()

}
