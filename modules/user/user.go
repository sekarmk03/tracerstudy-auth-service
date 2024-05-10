package user

import (
	"tracerstudy-auth-service/common/config"
	"tracerstudy-auth-service/modules/user/builder"
	"tracerstudy-auth-service/pb"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func InitGrpc(server *grpc.Server, cfg config.Config, db *gorm.DB, grpcConn *grpc.ClientConn) {
	user := builder.BuildUserHandler(cfg, db, grpcConn)
	pb.RegisterUserServiceServer(server, user)
}
