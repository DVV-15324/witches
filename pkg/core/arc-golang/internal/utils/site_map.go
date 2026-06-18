package utils

import (
	"context"
	"encoding/xml"
	"time"
)

// Url struct định nghĩa 1 URL trong sitemap
// loc: URL đầy đủ ví dụ https://example.com/all_blog
// LastMod: Thời gian cập nhật lần cuối
// ChangeFreq: Tuần suất thay đổi ví dụ:
// always
// hourly
// daily
// weekly
// monthly
// yearly
// never
// Priority: Độ ưu tiên khi Search (thường dựa trên bài blog hoặc bài hay nhất trong chủ đề đang tìm)
type Url struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod,omitempty"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

// Urlset struct theo chuẩn sitemap
// Khi marshal bằng package Go encoding/xml, nó sẽ tạo: <urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
// XMLName xác định root tag là: <urlset>
// Xmlns field này là XML attribute => sinh ra xmlns="..."
// Urls []Url mỗi phần tử sẽ thành: <url>
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

func (s *SiteMap) BusinessSiteMap(cxt context.Context, loc string, changeFreq string, priviority string) *Urlset {
	var urls []Url
	if len(changeFreq) < 1 {
		changeFreq = "always"
	}
	if len(priviority) < 1 {
		priviority = "0.8"
	}
	for _, slug := range s.Slugs {
		url := Url{
			Loc:        loc + slug.Slug,
			LastMod:    slug.UpdatedAt.Format("2006-01-02"),
			ChangeFreq: changeFreq,
			Priority:   priviority,
		}
		urls = append(urls, url)
	}

	urlset := Urlset{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		Urls:  urls,
	}
	return &urlset
}

func NewBusinessSiteMap(slugs []SiteMapData) *SiteMap {
	return &SiteMap{
		Slugs: slugs,
	}
}
