package builder

import (
	"tracerstudy-auth-service/common/config"
	commonJwt "tracerstudy-auth-service/common/jwt"
	"tracerstudy-auth-service/modules/auth/client"
	"tracerstudy-auth-service/modules/auth/handler"

	// mhsSvc "tracerstudy-auth-service/modules/mhsbiodata/service"
	// pktsRepo "tracerstudy-auth-service/modules/pkts/repository"
	// pktsSvc "tracerstudy-auth-service/modules/pkts/service"
	userRepo "tracerstudy-auth-service/modules/user/repository"
	userSvc "tracerstudy-auth-service/modules/user/service"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func BuildAuthHandler(cfg config.Config, db *gorm.DB, jwtManager *commonJwt.JWT, grpcConn *grpc.ClientConn) *handler.AuthHandler {
	// mhsBiodataSvc := mhsSvc.NewMhsBiodataService(cfg)

	// pktsRepository := pktsRepo.NewPKTSRepository(db)
	// pktsSvc := pktsSvc.NewPKTSService(cfg, pktsRepository)

	userRepository := userRepo.NewUserRepository(db)
	userSvc := userSvc.NewUserService(cfg, userRepository)

	pktsSvc := client.BuildPktsServiceClient(cfg.ClientURL.Pkts)
	mhsSvc := client.BuildMhsBiodataServiceClient(cfg.ClientURL.MhsBiodata)

	return handler.NewAuthHandler(cfg /*mhsBiodataSvc, pktsSvc,*/, userSvc, jwtManager, pktsSvc, mhsSvc)
}
