package builder

import (
	"tracerstudy-auth-service/common/config"
	commonJwt "tracerstudy-auth-service/common/jwt"
	"tracerstudy-auth-service/modules/auth/client"
	"tracerstudy-auth-service/modules/auth/handler"
	userRepo "tracerstudy-auth-service/modules/user/repository"
	userSvc "tracerstudy-auth-service/modules/user/service"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func BuildAuthHandler(cfg config.Config, db *gorm.DB, jwtManager *commonJwt.JWT, grpcConn *grpc.ClientConn) *handler.AuthHandler {
	userRepository := userRepo.NewUserRepository(db)
	userSvc := userSvc.NewUserService(cfg, userRepository)

	pktsSvc := client.BuildPktsServiceClient(cfg.ClientURL.Pkts)
	mhsSvc := client.BuildMhsBiodataServiceClient(cfg.ClientURL.MhsBiodata)

	return handler.NewAuthHandler(cfg, userSvc, jwtManager, pktsSvc, mhsSvc)
}
