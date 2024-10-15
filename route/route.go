package route

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"

	reflectv1beta1connect "buf.build/gen/go/bufbuild/reflect/connectrpc/go/buf/reflect/v1beta1/reflectv1beta1connect"
	reflectv1beta1 "buf.build/gen/go/bufbuild/reflect/protocolbuffers/go/buf/reflect/v1beta1"
	"connectrpc.com/connect"
	"github.com/vic3lord/bufile/proto/bufile/v1"
	"google.golang.org/protobuf/types/descriptorpb"

	_ "embed"
)

type ServiceProfile struct {
	Service   string
	Namespace string
	Method    string
	Routes    []Rule
}

type Rule struct {
	Name       string
	PathRegex  string
	Retryable  bool
	Deprecated bool
	Timeout    string
}

//go:embed template.yaml
var routesTemplate string

func Generate(ctx context.Context, w io.Writer) error {
	tmpl, err := template.New("routes-template").Parse(routesTemplate)
	if err != nil {
		return err
	}

	res, err := authorizedRequest(ctx, "buf.build/tannin/backend")
	if err != nil {
		return err
	}

	svc := ServiceProfile{Namespace: "default"}
	for _, file := range res.Msg.GetFileDescriptorSet().GetFile() {
		fileOpts := file.GetOptions()
		if fileOpts == nil {
			continue
		}

		service := fileOpts.ProtoReflect().Get(bufilev1.E_LinkerdService.TypeDescriptor())
		if service.IsValid() {
			svc.Service = service.String()
		}

		ns := fileOpts.ProtoReflect().Get(bufilev1.E_LinkerdNamespace.TypeDescriptor())
		if ns.IsValid() {
			svc.Namespace = ns.String()
		}

		for _, service := range file.GetService() {
			packageAndService := fmt.Sprintf("%s.%s", file.GetPackage(), service.GetName())
			for _, method := range service.GetMethod() {
				rt := Rule{
					Name:      method.GetName(),
					PathRegex: fmt.Sprintf("%s/%s", packageAndService, method.GetName()),
				}

				methodOpts := method.GetOptions()
				if methodOpts == nil {
					continue
				}

				rt.Deprecated = methodOpts.GetDeprecated()
				if methodOpts.GetIdempotencyLevel() != descriptorpb.MethodOptions_IDEMPOTENCY_UNKNOWN {
					rt.Retryable = true
				}

				timeout := methodOpts.ProtoReflect().Get(bufilev1.E_LinkerdTimeout.TypeDescriptor())
				if timeout.IsValid() {
					rt.Timeout = timeout.String()
				}
				svc.Routes = append(svc.Routes, rt)
			}
		}
	}
	return tmpl.Execute(w, svc)
}

func authorizedRequest(ctx context.Context, mod string) (*connect.Response[reflectv1beta1.GetFileDescriptorSetResponse], error) {
	client := reflectv1beta1connect.NewFileDescriptorSetServiceClient(
		http.DefaultClient,
		"https://buf.build:443",
	)

	msg := &reflectv1beta1.GetFileDescriptorSetRequest{
		Module:  mod,
		Version: "dem-6413-write-a-tool-that-syncs-buf-into-linkerd-service-profiles",
	}

	// Add Authorization header to the request.
	req := connect.NewRequest(msg)
	req.Header().Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("BUF_TOKEN")))

	return client.GetFileDescriptorSet(ctx, req)
}
