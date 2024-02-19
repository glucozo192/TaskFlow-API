package http_server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/glu-project/idl/pb"
	"github.com/glu-project/internal/user/models"
	"github.com/glu-project/utils/authenticate"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/rs/zerolog/log"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
)

type middlewareFunc func(http.Handler) http.Handler

var (
	Authorization = "Authorization"
	Bearer        = "Bearer"
	InfoKey       = "info_key"
	CookieKey     = "h5token"

	ignoredAPIs []string
)

func init() {
	ignoredAPIs = []string{
		pb.AllPathToMethodMap[pb.UserService_Login_API],
		pb.AllPathToMethodMap[pb.UserService_Register_API],
	}
}

type payloadKeys struct{}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, authorization")
		if r.Method != "OPTIONS" {
			h.ServeHTTP(w, r)
		}
	})
}

type Response struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details"`
	Data    any      `json:"data"`
}

func ErrorResponse(w http.ResponseWriter, code int, err error) {
	resp := &Response{
		Code:    code,
		Message: err.Error(),
		Details: []string{},
	}

	jData, _ := json.Marshal(resp)

	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func DataResponse(w http.ResponseWriter, data any) {
	resp := &Response{
		Data: data,
	}

	jData, _ := json.Marshal(resp)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func forwardErrorResponse(ctx context.Context, s *runtime.ServeMux, m runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	sta := status.Convert(err)
	errStr := sta.Message()
	log.Print(err)
	firstColonPos := strings.Index(errStr, ":")

	if firstColonPos > 0 {
		errStr = errStr[:firstColonPos]
	}

	runtime.DefaultHTTPErrorHandler(ctx, s, m, w, r, errors.New(errStr))
}

func authorized(
	authenticator authenticate.Authenticator,
	db models.DBTX,
	rolePermissionsCache *expirable.LRU[string, []string],
) middlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			incomingPath := fmt.Sprintf("%s %s", r.Method, r.URL.Path)

			if isCheckingPassed(ignoredAPIs, incomingPath) {
				h.ServeHTTP(w, r)
				return
			}

			authorization := r.Header.Get(Authorization)
			scheme, token, found := strings.Cut(authorization, " ")
			if !found {
				ErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("invalid authorization token"))
				return
			}
			if !strings.EqualFold(scheme, Bearer) {
				ErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("invalid authorization token"))
				return
			}
			payload, err := authenticator.Verify(token)
			if err != nil {
				ErrorResponse(w, http.StatusUnauthorized, err)
				return
			}

			payload.Token = token
			if payload.RoleID == "" {
				ErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("you don't have roles to access this URL"))
				return
			}

			_, ok := rolePermissionsCache.Get(payload.RoleID)
			// TODO: refactor in future
			var rolePermissionRepo rolePermissionRepo = new(postgres.RolePermissionRepository)
			if !ok {
				resp, err := rolePermissionRepo.GetListRolePermissions(ctx, db)
				if err != nil {
					ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("got unexpected error when listing enabled paths: %v", err))
					return
				}
				for i := 0; i < len(resp); i++ {
					rolePermissionsCache.Add(resp[i].RoleID, resp[i].Permissions)
				}
			}

			// get current token in DB and check. If token is not match return http.StatusUnauthorized
			var userLoginRepo userLoginRepo = new(postgres.LoginUserRepository)
			user, err := userLoginRepo.GetByID(ctx, db, payload.UserID)
			if err != nil {
				ErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("got unexpected: %v", err))
				return
			}

			if user.Token.String != payload.Token {
				ErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("your token is not valid"))
				return
			}

			permissions, ok := rolePermissionsCache.Get(payload.RoleID)
			if !ok {
				ErrorResponse(w, http.StatusUnauthorized, fmt.Errorf("your token is not valid"))
				return
			}

			if !isCheckingPassed(permissions, incomingPath) {
				ErrorResponse(w, http.StatusForbidden, fmt.Errorf("you don't have permission to access this"))
				return
			}

			h.ServeHTTP(w, r.WithContext(context.WithValue(ctx, payloadKeys{}, payload)))
			if isCheckingPassed(invalidateCacheAPIs, incomingPath) {
				for _, key := range rolePermissionsCache.Keys() {
					rolePermissionsCache.Remove(key)
				}
			}
		})
	}
}

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func (rec *ResponseRecorder) Write(body []byte) (int, error) {
	rec.Body = body
	return rec.ResponseWriter.Write(body)
}

func HttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		rec := &ResponseRecorder{
			ResponseWriter: res,
			StatusCode:     http.StatusOK,
		}
		handler.ServeHTTP(rec, req)
		duration := time.Since(startTime)

		logger := log.Info()
		if rec.StatusCode != http.StatusOK {
			logger = log.Error().Bytes("body", rec.Body)
		}

		logger.Str("protocol", "http").
			Str("method", req.Method).
			Str("path", req.RequestURI).
			Int("status_code", rec.StatusCode).
			Str("status_text", http.StatusText(rec.StatusCode)).
			Dur("duration", duration).
			Msg("received a HTTP request")
	})
}
