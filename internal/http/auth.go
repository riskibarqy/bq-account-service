package http

import (
	"net/http"
	"strings"

	"github.com/riskibarqy/bq-account-service/internal/http/response"
	"github.com/riskibarqy/bq-account-service/internal/types"
	"github.com/riskibarqy/bq-account-service/internal/user"
)

func (hs *Server) authorizedOnly(userService user.ServiceInterface) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var token string

			ctx := r.Context()
			token = getBearerToken(r)
			if token == "" {
				response.Error(w, "Unauthorized", http.StatusUnauthorized, types.Error{
					Path:    ".Server->authorizeOnly()",
					Message: "",
					Error:   nil,
					Type:    "",
				})
				return
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func getBearerToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	splitToken := strings.Split(token, "Bearer")

	if len(splitToken) < 2 {
		return ""
	}

	token = strings.Trim(splitToken[1], " ")
	return token
}

// func getBasicToken(r *http.Request) string {
// 	token := r.Header.Get("Authorization")
// 	splitToken := strings.Split(token, "Basic")

// 	if len(splitToken) < 2 {
// 		return ""
// 	}

// 	token = strings.Trim(splitToken[1], " ")
// 	return token
// }

// func getXAccessToken(r *http.Request) string {
// 	token := r.Header.Get("X-Access-Token")
// 	return token
// }

// func getVersion(r *http.Request) int {
// 	stringVersion := strings.Split(r.UserAgent(), "(")
// 	versionStr := strings.Replace(stringVersion[0][strings.LastIndex(stringVersion[0], ".")+1:], " ", "", -1)
// 	var version int

// 	version, err := strconv.Atoi(versionStr)
// 	if err != nil {
// 		version = -1
// 	}
// 	return version
// }
