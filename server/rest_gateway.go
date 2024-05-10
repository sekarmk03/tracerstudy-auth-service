package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

type Rest struct {
	*runtime.ServeMux
	port string
}

type ErrorData struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Error struct {
	Error ErrorData   `json:"error"`
	Meta  interface{} `json:"meta"`
}

func NewRest(port string) *Rest {
	return &Rest{
		ServeMux: runtime.NewServeMux(
			runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
					MarshalOptions: protojson.MarshalOptions{
						UseProtoNames: true,
						EmitUnpopulated: true,
					},
					UnmarshalOptions: protojson.UnmarshalOptions{
						DiscardUnknown: true,
					},
				}),
				runtime.WithErrorHandler(customErrorHandler),
			),
			port: port,
	}
}

func (r *Rest) Run() error {
	go func ()  {
		if err := http.ListenAndServe(fmt.Sprintf(":%s", r.port), allowCORS(r.ServeMux)); err != nil {
			panic(err)
		}
	}()
	return nil
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				headers := []string{"Content-Type", "Accept", "Authorization"}
				w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
				methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func customErrorHandler(ctx context.Context, mux *runtime.ServeMux, mrs runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	const fallback = `{"error": {"code":13,"message":"failed to marshal error message"}, "meta":null}`

	s := status.Convert(err)

	w.Header().Set("Content-type", mrs.ContentType("application/json"))
	w.WriteHeader(runtime.HTTPStatusFromCode(s.Code()))
	jsonErr := json.NewEncoder(w).Encode(Error{
		Error: ErrorData{
			Code:    int(s.Code()),
			Message: s.Message(),
		},
		Meta: nil,
	})

	if jsonErr != nil {
		_, _ = w.Write([]byte(fallback))
	}
}
