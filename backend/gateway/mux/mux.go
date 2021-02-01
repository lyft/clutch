package mux

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/pprof"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	gatewayv1 "github.com/lyft/clutch/backend/api/config/gateway/v1"
	"github.com/lyft/clutch/backend/service"
	awsservice "github.com/lyft/clutch/backend/service/aws"
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
		// If not a known static asset and an asset provider is configured, try streaming from the configured provider.
		if a.assetCfg != nil && a.assetCfg.Provider != nil && strings.HasPrefix(r.URL.Path, "/static/") {
			// We attach this header simply for observability purposes.
			// Otherwise its difficult to know if the assets are being served from the configured provider.
			w.Header().Set("x-clutch-asset-passthrough", "true")

			asset, err := a.assetProviderHandler(r.Context(), r.URL.Path)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(fmt.Sprintf("Error getting assets from the configured asset provider: %v", err)))
				return
			}
			defer asset.Close()

			_, err = io.Copy(w, asset)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				_, _ = w.Write([]byte(fmt.Sprintf("Error getting assets from the configured asset provider: %v", err)))
				return
			}
			return
		}

		// If not a known static asset serve the SPA.
		r.URL.Path = "/"
	} else {
		_ = f.Close()
	}

	a.fileServer.ServeHTTP(w, r)
}

func (a *assetHandler) assetProviderHandler(ctx context.Context, urlPath string) (io.ReadCloser, error) {
	switch a.assetCfg.Provider.(type) {
	case *gatewayv1.Assets_S3:
		aws, err := getAssetProviderService(a.assetCfg)
		if err != nil {
			return nil, err
		}

		awsClient, ok := aws.(awsservice.Client)
		if !ok {
			return nil, fmt.Errorf("Unable to aquire the aws client")
		}

		return awsClient.S3StreamingGet(
			ctx,
			a.assetCfg.GetS3().Region,
			a.assetCfg.GetS3().Bucket,
			path.Join(a.assetCfg.GetS3().Key, strings.TrimPrefix(urlPath, "/static")),
		)
	default:
		return nil, fmt.Errorf("configured asset provider has not been implemented")
	}
}

// getAssetProviderService is used in two different contexts
// Its invoked in the mux constructor which checks if the necessary service has been configured,
// if there is an asset provider which requires ones.
//
// Otherwise its used to get the service for an asset provider in assetProviderHandler() if necessary.
func getAssetProviderService(assetCfg *gatewayv1.Assets) (service.Service, error) {
	switch assetCfg.Provider.(type) {
	case *gatewayv1.Assets_S3:
		aws, ok := service.Registry[awsservice.Name]
		if !ok {
			return nil, fmt.Errorf("The AWS service must be configured to use the asset s3 provider.")
		}
		return aws, nil

	default:
		// An asset provider does not necessarily require a service to function properly
		// if there is nothing configured for a provider type we cant necessarily throw an error here.
		return nil, nil
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
				runtime.DefaultHTTPErrorHandler(ctx, mux, m, w, req, err)
				return
			}
			if redirectPath != referer.Path {
				redirectPath = fmt.Sprintf("%s?redirect_url=%s", redirectPath, referer.Path)
			}
		}

		http.Redirect(w, req, redirectPath, http.StatusFound)
		return
	}
	runtime.DefaultHTTPErrorHandler(ctx, mux, m, w, req, err)
}

func New(unaryInterceptors []grpc.UnaryServerInterceptor, assets http.FileSystem, gatewayCfg *gatewayv1.GatewayOptions) (*Mux, error) {
	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(unaryInterceptors...))
	jsonGateway := runtime.NewServeMux(
		runtime.WithForwardResponseOption(customResponseForwarder),
		runtime.WithErrorHandler(customErrorHandler),
		runtime.WithMarshalerOption(
			runtime.MIMEWildcard,
			&runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					// Use camelCase for the JSON version.
					UseProtoNames: false,
					// Transmit zero-values over the wire.
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{},
			},
		),
	)

	// If there is a configured asset provider, we check to see if the service is configured before proceeding.
	// Bailing out early during the startup process instead of hitting this error at runtime when serving assets.
	if gatewayCfg.Assets != nil && gatewayCfg.Assets.Provider != nil {
		_, err := getAssetProviderService(gatewayCfg.Assets)
		if err != nil {
			return nil, err
		}
	}

	httpMux := http.NewServeMux()
	httpMux.Handle("/", &assetHandler{
		assetCfg:   gatewayCfg.Assets,
		next:       jsonGateway,
		fileSystem: assets,
		fileServer: http.FileServer(assets),
	})

	if gatewayCfg.EnablePprof {
		httpMux.HandleFunc("/debug/pprof/", pprof.Index)
	}

	mux := &Mux{
		GRPCServer:  grpcServer,
		JSONGateway: jsonGateway,
		HTTPMux:     httpMux,
	}
	return mux, nil
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
