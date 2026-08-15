package utils

import (
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/DVV-15324/witches/cmd/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

type DeviceHelper struct {
	config  *utils.Config
	matcher language.Matcher
}

// NewDeviceHelper khởi tạo helper với config
func NewDeviceHelper(cfg *utils.Config) *DeviceHelper {

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
	return &DeviceHelper{
		config:  cfg,
		matcher: language.NewMatcher(tags),
	}
}

func (h *DeviceHelper) GetDeviceInfo(c *gin.Context) (deviceID, ipAddress, userAgent, locale, timezone string) {
	reqCtx := GetRequestContext(c.Request.Context(), h.config)
	if reqCtx.DeviceID != "" {
		deviceID = reqCtx.DeviceID
		ipAddress = reqCtx.IPAddress
		userAgent = reqCtx.UserAgent
	} else {
		ipAddress = c.ClientIP()
		userAgent = c.GetHeader("User-Agent")
		deviceID = generateDeviceID(ipAddress, userAgent)
	}

	acceptLang := c.GetHeader("Accept-Language")
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

func ExtractTokenFromHeader(accessToken string) (string, error) {
	if accessToken == "" {
		return "", errors.New("authorization header is required")
	}
	args := strings.Split(accessToken, " ")
	if len(args) != 2 || args[0] != "Bearer" {
		return "", errors.New("invalid authorization header format")
	}
	return args[1], nil
}

func generateDeviceID(ip, userAgent string) string {
	data := fmt.Sprintf("%s|%s", ip, userAgent)
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}
