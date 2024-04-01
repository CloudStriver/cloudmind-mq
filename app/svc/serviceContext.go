package svc

import (
	"github.com/CloudStriver/cloudmind-mq/app/config"
	"github.com/CloudStriver/go-pkg/utils/kitex/client"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/content/contentservice"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/cloudmind/system/systemservice"
	"github.com/CloudStriver/service-idl-gen-go/kitex_gen/platform/relation/relationservice"
)

type ServiceContext struct {
	Config              config.Config
	CloudMindSystemRPC  systemservice.Client
	CloudMindContentRPC contentservice.Client
	RelationRPC         relationservice.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:              c,
		CloudMindSystemRPC:  client.NewClient(c.Name, "cloudmind-system", systemservice.NewClient),
		CloudMindContentRPC: client.NewClient(c.Name, "cloudmind-content", contentservice.NewClient),
		RelationRPC:         client.NewClient(c.Name, "platform-relation", relationservice.NewClient),
	}
}
