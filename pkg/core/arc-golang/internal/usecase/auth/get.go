package auth

// import (
// 	entityAuth "arc-golang/internal/entity/auth"
// 	c_response "arc-golang/internal/utils"
// 	"context"
// 	"time"
// )

// func (a *UsecaseAuth) BzGetAuthByEmail(ctx context.Context, email string) (*entityAuth.Auth, *c_response.Response) {
// 	auth, err := a.ReponsitoryAuth.GetAuthByEmail(ctx, email)
// 	if err != nil {
// 		resp := c_response.NewResponse(404, err, time.Now(), nil)
// 		return nil, resp
// 	}
// 	return auth, nil
// }
