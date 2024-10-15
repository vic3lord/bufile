package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"text/template"

	reflectv1beta1connect "buf.build/gen/go/bufbuild/reflect/connectrpc/go/buf/reflect/v1beta1/reflectv1beta1connect"
	reflectv1beta1 "buf.build/gen/go/bufbuild/reflect/protocolbuffers/go/buf/reflect/v1beta1"
	"connectrpc.com/connect"
	"github.com/vic3lord/probufile/proto/probufile/v1"
	"google.golang.org/protobuf/types/descriptorpb"

	_ "embed"
)

type ServiceProfile struct {
	Service   string
	Namespace string
	Method    string
	Routes    []Route
}

type Route struct {
	Name       string
	PathRegex  string
	Retryable  bool
	Deprecated bool
}

//go:embed template/service-profile.yaml
var serviceProfileTemplate string

var print bool

func main() {
	tmpl, err := template.New("service-profile").Parse(serviceProfileTemplate)
	if err != nil {
		panic(err)
	}

	client := reflectv1beta1connect.NewFileDescriptorSetServiceClient(http.DefaultClient, "https://buf.build:443")
	req := connect.NewRequest(&reflectv1beta1.GetFileDescriptorSetRequest{
		Module:  "buf.build/tannin/backend",
		Version: "dem-6413-write-a-tool-that-syncs-buf-into-linkerd-service-profiles",
	})

	// Add Authorization header to the request.
	req.Header().Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("BUF_TOKEN")))

	res, err := client.GetFileDescriptorSet(context.Background(), req)
	if err != nil {
		panic(err)
	}

	svc := ServiceProfile{
		Service:   "backend",
		Namespace: "default",
	}
	for _, file := range res.Msg.GetFileDescriptorSet().GetFile() {
		if opts := file.GetOptions(); opts != nil {
			svcURL := opts.ProtoReflect().Get(probufilev1.E_LinkerdService.TypeDescriptor())
			if svcURL.IsValid() {
				svc.Service = svcURL.String()
			}

			ns := opts.ProtoReflect().Get(probufilev1.E_LinkerdNamespace.TypeDescriptor())
			if svcURL.IsValid() {
				svc.Namespace = ns.String()
			}
		}

		for _, service := range file.GetService() {
			packageAndService := fmt.Sprintf("%s.%s", file.GetPackage(), service.GetName())
			for _, method := range service.GetMethod() {
				rt := Route{
					Name:      method.GetName(),
					PathRegex: fmt.Sprintf("%s/%s", packageAndService, method.GetName()),
				}

				if opt := method.GetOptions(); opt != nil {
					rt.Deprecated = opt.GetDeprecated()
					if opt.GetIdempotencyLevel() != descriptorpb.MethodOptions_IDEMPOTENCY_UNKNOWN {
						rt.Retryable = true
					}
					// fmt.Printf("method opts: %+v\n", opt)
				}
				svc.Routes = append(svc.Routes, rt)
			}
		}
	}
	if !print {
		tmpl.Execute(os.Stdout, svc)
	}
}
