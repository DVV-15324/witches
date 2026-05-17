package sitemap

import (
	"context"
	"encoding/xml"
	"time"
)

// Url struct định nghĩa 1 URL trong sitemap
type Url struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod,omitempty"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

// Urlset struct theo chuẩn sitemap
type Urlset struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	Urls    []Url    `xml:"url"`
}

type SiteMapData struct {
	Slug      string
	UpdatedAt time.Time
}

type SiteMap struct {
	Slugs []SiteMapData
}

func (s *SiteMap) BusinessSitemap(cxt context.Context) *Urlset {
	var urls []Url
	for _, slug := range s.Slugs {
		url := Url{
			Loc:        "https://bloghomnay.com/post/" + slug.Slug, // hoặc id hoặc gì phù hợp
			LastMod:    slug.UpdatedAt.Format("2006-01-02"),
			ChangeFreq: "always",
			Priority:   "0.8",
		}
		urls = append(urls, url)
	}

	urlset := Urlset{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Urls:  urls,
	}
	return &urlset
}

func NewBusinessSitePost(slugs []SiteMapData) *SiteMap {
	return &SiteMap{
		Slugs: slugs,
	}
}
