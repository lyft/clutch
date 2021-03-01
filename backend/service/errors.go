package service

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

func CodeFromHTTPStatus(status int) codes.Code {
	switch status {
	case http.StatusOK: // 200
		return codes.OK
	case http.StatusBadRequest: // 400
		return codes.FailedPrecondition
	case http.StatusUnauthorized: // 401
		return codes.Unauthenticated
	case http.StatusForbidden: // 403
		return codes.PermissionDenied
	case http.StatusNotFound: // 404
		return codes.NotFound
	case http.StatusConflict: // 409
		return codes.AlreadyExists
	case http.StatusGone: // 410
		return codes.NotFound
	case http.StatusInternalServerError: // 500
		return codes.Internal
	case http.StatusServiceUnavailable: // 503
		return codes.Unavailable
	case http.StatusGatewayTimeout: // 504
		return codes.DeadlineExceeded
	default:
		return codes.Unknown
	}
}
