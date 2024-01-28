package listen

import (
	"context"
	"github.com/CloudStriver/cloudmind-mq/internal/config"
	"github.com/CloudStriver/cloudmind-mq/internal/svc"

	"github.com/zeromicro/go-zero/core/service"
)

func Mqs(c config.Config) []service.Service {
	svcContext := svc.NewServiceContext(c)
	ctx := context.Background()
	var services []service.Service
	services = append(services, KqMqs(c, ctx, svcContext)...)
	return services
}
