package listen

import (
	"context"
	"github.com/CloudStriver/cloudmind-mq/app/config"
	"github.com/CloudStriver/cloudmind-mq/app/svc"
	"github.com/cloudwego/kitex/pkg/remote/trans/nphttp2/metadata"

	"github.com/zeromicro/go-zero/core/service"
)

const (
	EnvHeader = "X_XH_ENV" // 使用正确的header名称
)

func Mqs(c config.Config) []service.Service {
	svcContext := svc.NewServiceContext(c)
	md := metadata.MD{
		EnvHeader: []string{c.EnvHeader},
	}
	ctx := metadata.NewIncomingContext(context.Background(), md)

	var services []service.Service
	services = append(services, KqMqs(c, ctx, svcContext)...)
	return services
}
