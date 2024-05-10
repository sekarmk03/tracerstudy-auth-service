package builder

import (
	"tracerstudy-auth-service/common/config"
	"tracerstudy-auth-service/modules/user/handler"
	"tracerstudy-auth-service/modules/user/repository"
	"tracerstudy-auth-service/modules/user/service"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func BuildUserHandler(cfg config.Config, db *gorm.DB, grpcConn *grpc.ClientConn) *handler.UserHandler {
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(cfg, userRepo)

	return handler.NewUserHandler(cfg, userSvc)
}
