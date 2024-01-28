package svc

import (
	"github.com/CloudStriver/cloudmind-mq/app/config"
	"github.com/CloudStriver/go-pkg/utils/kitex/client"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/system/systemservice"
)

type ServiceContext struct {
	Config             config.Config
	CloudMindSystemRPC systemservice.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:             c,
		CloudMindSystemRPC: client.NewClient(c.Name, "cloudmind-system", systemservice.NewClient),
	}
}
