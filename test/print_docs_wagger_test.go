package test

import (
	docs "core-v/pkg/core/wagger/docs"
	"testing"
)

func TestDocsWaggerStructTag(t *testing.T) {
	docs.PrintDocsSwaggerStructTag()
}
func TestDocsWaggerMetaData(t *testing.T) {
	docs.PrintDocsSwaggerMetaData()
}
