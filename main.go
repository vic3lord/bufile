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

	_ "embed"
)

type ServiceProfile struct {
	ServiceName string
	Method      string
	Routes      []Route
}

type Route struct {
	Name      string
	PathRegex string
}

//go:embed template/service-profile.yaml
var serviceProfileTemplate string

func main() {
	tmpl, err := template.New("service-profile").Parse(serviceProfileTemplate)
	if err != nil {
		panic(err)
	}

	client := reflectv1beta1connect.NewFileDescriptorSetServiceClient(http.DefaultClient, "https://buf.build:443")
	req := connect.NewRequest(&reflectv1beta1.GetFileDescriptorSetRequest{
		Module: "buf.build/tannin/template",
	})

	// Add Authorization header to the request.
	req.Header().Set("Authorization", fmt.Sprintf("Bearer %s", os.Getenv("BUF_TOKEN")))

	res, err := client.GetFileDescriptorSet(context.Background(), req)
	if err != nil {
		panic(err)
	}

	svc := ServiceProfile{ServiceName: "template"}
	for _, file := range res.Msg.GetFileDescriptorSet().GetFile() {
		for _, service := range file.GetService() {
			opt := service.GetOptions().GetFeatures()
			fmt.Println(opt)
			packageAndService := fmt.Sprintf("%s.%s", file.GetPackage(), service.GetName())
			for _, method := range service.GetMethod() {
				svc.Routes = append(svc.Routes, Route{
					Name:      method.GetName(),
					PathRegex: fmt.Sprintf("%s/%s", packageAndService, method.GetName()),
				})
			}
		}
	}
	tmpl.Execute(os.Stdout, svc)
}
