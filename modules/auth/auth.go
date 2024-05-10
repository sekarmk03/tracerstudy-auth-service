package auth

import (
	"tracerstudy-auth-service/common/config"
	commonJwt "tracerstudy-auth-service/common/jwt"
	"tracerstudy-auth-service/modules/auth/builder"
	"tracerstudy-auth-service/pb"

	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func InitGrpc(server *grpc.Server, cfg config.Config, db *gorm.DB, jwtManager *commonJwt.JWT, grpcConn *grpc.ClientConn) {
	auth := builder.BuildAuthHandler(cfg, db, jwtManager, grpcConn)
	pb.RegisterAuthServiceServer(server, auth)
}
