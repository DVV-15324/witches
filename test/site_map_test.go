package test

import (
	"context"
	sitemap "core-v/pkg/core/utils"
	"encoding/xml"
	"fmt"
	"testing"
	"time"
)

func TestSiteMap(t *testing.T) {
	background := context.Background()

	slugs := []sitemap.SiteMapData{
		{
			Slug:      "Hello Nguoi Anh Em 1",
			UpdatedAt: time.Now(),
		},
		{
			Slug:      "Hello Nguoi Anh Em 2",
			UpdatedAt: time.Now(),
		},
	}

	businessSiteMap := sitemap.NewBusinessSiteMap(slugs)

	urlSet := businessSiteMap.BusinessSiteMap(
		background,
		"https://example.com/all_blog/",
		"always",
		"0.8",
	)

	data, err := xml.MarshalIndent(urlSet, "", "  ")
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(xml.Header + string(data))
}
