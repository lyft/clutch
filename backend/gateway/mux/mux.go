package mux

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/lyft/clutch/backend/service"

	gatewayv1 "github.com/lyft/clutch/backend/api/config/gateway/v1"
	awsservice "github.com/lyft/clutch/backend/service/aws"

	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var apiPattern = regexp.MustCompile(`^/v\d+/`)

type assetHandler struct {
	assetCfg *gatewayv1.Assets

	next       http.Handler
	fileSystem http.FileSystem
	fileServer http.Handler
}

func copyHTTPResponse(resp *http.Response, w http.ResponseWriter) {
	for key, values := range resp.Header {
		for _, val := range values {
			w.Header().Add(key, val)
		}
	}
	w.WriteHeader(resp.StatusCode)
	_, _ = io.Copy(w, resp.Body)
}

func (a *assetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if apiPattern.MatchString(r.URL.Path) || r.URL.Path == "/healthcheck" {
		// Serve from the embedded API handler.
		a.next.ServeHTTP(w, r)
		return
	}

	// Check if assets are okay to serve by calling the Fetch endpoint and verifying it returns a 200.
	rec := httptest.NewRecorder()
	origPath := r.URL.Path
	r.URL.Path = "/v1/assets/fetch"
	a.next.ServeHTTP(rec, r)

	if rec.Code != http.StatusOK {
		copyHTTPResponse(rec.Result(), w)
		return
	}

	// Set the original path.
	r.URL.Path = origPath

	// Serve!
	if f, err := a.fileSystem.Open(r.URL.Path); err != nil {
		// If not a known static asset and a asset provider is configured, try streaming from the configured provider.
		if a.assetCfg.Provider != nil && strings.HasPrefix(r.URL.Path, "/static") {
			// We attatch this header simply for observibility purposes.
			// Otherwise its difficult to know if the assets are being served from the configured provider.
			w.Header().Set("x-clutch-asset-passthrough", "true")

			// TODO: handle errors
			asset, _ := a.assetProviderHandler(r.URL.Path)
			_, _ = io.Copy(w, asset)
			return
		}

		// If not a known static asset serve the SPA.
		r.URL.Path = "/"
	} else {
		_ = f.Close()
	}

	a.fileServer.ServeHTTP(w, r)
}

func (a *assetHandler) assetProviderHandler(urlPath string) (io.Reader, error) {
	switch a.assetCfg.Provider.(type) {
	case *gatewayv1.Assets_S3:
		aws, ok := service.Registry[awsservice.Name]
		if !ok {
			return nil, fmt.Errorf("The AWS service must be configured to use the asset s3 provider.")
		}

		awsClient, ok := aws.(awsservice.Client)
		if !ok {
			return nil, fmt.Errorf("Unable to aquire the aws client")
		}

		return awsClient.S3StreamingGet(
			context.Background(),
			awsClient.GetCurrentRegion(),
			a.assetCfg.GetS3().Bucket,
			path.Join(a.assetCfg.GetS3().Key, strings.TrimPrefix(urlPath, "/static")),
		)
	default:
		return nil, fmt.Errorf("No asset provider is configured")
	}
}

func customResponseForwarder(ctx context.Context, w http.ResponseWriter, resp proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	if cookies := md.HeaderMD.Get("Set-Cookie-Token"); len(cookies) > 0 {
		cookie := &http.Cookie{
			Name:     "token",
			Value:    cookies[0],
			Path:     "/",
			HttpOnly: false,
		}
		http.SetCookie(w, cookie)
	}

	if redirects := md.HeaderMD.Get("Location"); len(redirects) > 0 {
		w.Header().Set("Location", redirects[0])

		code := http.StatusFound
		if st := md.HeaderMD.Get("Location-Status"); len(st) > 0 {
			if newCode, err := strconv.Atoi(st[0]); err != nil {
				code = newCode
			}
		}
		w.WriteHeader(code)
	}
	return nil
}

func customErrorHandler(ctx context.Context, mux *runtime.ServeMux, m runtime.Marshaler, w http.ResponseWriter, req *http.Request, err error) {
	//  TODO(maybe): once we have non-browser clients we probably want to avoid the redirect and directly return the error.
	if s, ok := status.FromError(err); ok && s.Code() == codes.Unauthenticated {
		referer := req.Referer()
		redirectPath := "/v1/authn/login"
		if len(referer) != 0 {
			referer, err := url.Parse(referer)
			if err != nil {
				runtime.DefaultHTTPProtoErrorHandler(ctx, mux, m, w, req, err)
				return
			}
			if redirectPath != referer.Path {
				redirectPath = fmt.Sprintf("%s?redirect_url=%s", redirectPath, referer.Path)
			}
		}

		http.Redirect(w, req, redirectPath, http.StatusFound)
		return
	}
	runtime.DefaultHTTPProtoErrorHandler(ctx, mux, m, w, req, err)
}

func New(unaryInterceptors []grpc.UnaryServerInterceptor, assets http.FileSystem, assetCfg *gatewayv1.Assets) *Mux {
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(unaryInterceptors...))
	jsonGateway := runtime.NewServeMux(
		runtime.WithForwardResponseOption(customResponseForwarder),
		runtime.WithProtoErrorHandler(customErrorHandler),
		runtime.WithMarshalerOption(
			runtime.MIMEWildcard,
			&runtime.JSONPb{
				// Use camelCase for the JSON version.
				OrigName: false,
				// Transmit zero-values over the wire.
				EmitDefaults: true,
			},
		),
	)

	httpMux := http.NewServeMux()
	httpMux.Handle("/", &assetHandler{
		assetCfg:   assetCfg,
		next:       jsonGateway,
		fileSystem: assets,
		fileServer: http.FileServer(assets),
	})

	mux := &Mux{
		GRPCServer:  grpcServer,
		JSONGateway: jsonGateway,
		HTTPMux:     httpMux,
	}
	return mux
}

// Mux allows sharing one port between gRPC and the corresponding JSON gateway via header-based multiplexing.
type Mux struct {
	// Create empty handlers for gRPC and grpc-gateway (JSON) traffic.
	JSONGateway *runtime.ServeMux
	HTTPMux     http.Handler
	GRPCServer  *grpc.Server
}

// Adapted from https://github.com/grpc/grpc-go/blob/197c621/server.go#L760-L778.
func (m *Mux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.ProtoMajor == 2 && strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
		m.GRPCServer.ServeHTTP(w, r)
	} else {
		m.HTTPMux.ServeHTTP(w, r)
	}
}

func (m *Mux) EnableGRPCReflection() {
	reflection.Register(m.GRPCServer)
}

// "h2c" is the unencrypted form of HTTP/2.
func InsecureHandler(handler http.Handler) http.Handler {
	return h2c.NewHandler(handler, &http2.Server{})
}
