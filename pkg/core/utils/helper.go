package utils

import (
	"github.com/DVV-15324/witches/cmd/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
	"strings"
)

type Helper struct {
	config  *utils.Config
	matcher language.Matcher
}

func NewHelper(cfg *utils.Config) *Helper {
	tags := make([]language.Tag, 0, len(cfg.SupportedLanguages))
	for _, locale := range cfg.SupportedLanguages {
		tag, err := language.Parse(locale)
		if err == nil {
			tags = append(tags, tag)
		}
	}
	if len(tags) == 0 {
		tags = []language.Tag{language.English, language.MustParse("vi-VN")}
	}
	return &Helper{
		config:  cfg,
		matcher: language.NewMatcher(tags),
	}
}

func (h *Helper) GetInfo(c *gin.Context) (dpopProof string, ipAddress, userAgent, locale, timezone string) {
	ipAddress = c.ClientIP()
	userAgent = c.GetHeader("User-Agent")
	acceptLang := c.GetHeader("Accept-Language")
	dpopProof = c.GetHeader("DPoP")
	if acceptLang != "" {
		tag, _ := language.MatchStrings(h.matcher, acceptLang)
		base, _ := tag.Base()
		region, _ := tag.Region()
		if region == (language.Region{}) {
			switch base.String() {
			case "vi":
				region = language.MustParseRegion("VN")
			default:
				region = language.MustParseRegion("US")
			}
		}
		locale = base.String() + "-" + region.String()
	} else {
		locale = "en-US"
	}
	timezone = c.GetHeader("X-Timezone")
	if timezone == "" {
		timezone = "UTC"
	}
	return
}

func (h *Helper) ExtractDPoPTokenFromHeader(header string) (string, error) {
	if header == "" {
		return "", errors.New("authorization header is required")
	}

	parts := strings.Fields(header)

	if len(parts) != 2 || parts[0] != "DPoP" {
		return "", errors.New("invalid authorization header format")
	}

	return parts[1], nil
}
