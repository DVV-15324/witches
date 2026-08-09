package utils

import (
	"crypto/md5"
	"fmt"
	"strings"

	u_ctx "github.com/DVV-15324/witches/pkg/core/utils"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/text/language"
)

// Hỗ trợ các locale
var supportedLocales = []language.Tag{
	language.English,            // en-US
	language.MustParse("vi-VN"), // vi-VN
}

var matcher = language.NewMatcher(supportedLocales)

// GetDeviceInfo lấy toàn bộ thông tin client: deviceID, ip, userAgent, locale, timezone
func GetDeviceInfo(c *gin.Context) (deviceID, ipAddress, userAgent, locale, timezone string) {
	// Lấy device info
	reqCtx := u_ctx.GetRequestContext(c.Request.Context(), ReqKey)
	if reqCtx.DeviceID != "" {
		deviceID = reqCtx.DeviceID
		ipAddress = reqCtx.IPAddress
		userAgent = reqCtx.UserAgent
	} else {
		ipAddress = c.ClientIP()
		userAgent = c.GetHeader("User-Agent")
		deviceID = generateDeviceID(ipAddress, userAgent)
	}

	// Lấy locale từ Accept-Language
	acceptLang := c.GetHeader("Accept-Language")
	if acceptLang != "" {
		tag, _ := language.MatchStrings(matcher, acceptLang)
		base, _ := tag.Base()
		region, _ := tag.Region()

		// Nếu không có region, gán mặc định
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

	// Lấy timezone từ header
	timezone = c.GetHeader("X-Timezone")
	if timezone == "" {
		timezone = "UTC"
	}

	return deviceID, ipAddress, userAgent, locale, timezone
}

// generateDeviceID - Tạo deviceID từ IP + User-Agent
func generateDeviceID(ip, userAgent string) string {
	data := fmt.Sprintf("%s|%s", ip, userAgent)
	return fmt.Sprintf("%x", md5.Sum([]byte(data)))
}

// extractTokenFromHeader - Trích xuất token từ header Authorization
func ExtractTokenFromHeader(accessToken string) (string, error) {
	if accessToken == "" {
		return "", errors.New("authorization header is required")
	}

	args := strings.Split(accessToken, " ")
	// Thiếu bearer, thiếu token, bị nhiều " "
	if len(args) != 2 || args[0] != "Bearer" {
		return "", errors.New("invalid authorization header format")
	}
	return args[1], nil
}
