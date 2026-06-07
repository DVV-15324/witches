package test

import (
	response "core-v/pkg/core/reponse"
	"testing"
	"time"
)

func TestReponse(t *testing.T) {
	data := map[string]interface{}{
		"name": "vu",
	}
	testReponse := response.NewEntityResponse(200, "Thanh Cong", time.Now(), data)
	testReponse.PrintReponse()

}
