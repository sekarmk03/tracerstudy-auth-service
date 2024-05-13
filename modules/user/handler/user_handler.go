package handler

import (
	"context"
	"log"
	"net/http"
	"tracerstudy-auth-service/common/config"
	"tracerstudy-auth-service/common/errors"
	"tracerstudy-auth-service/modules/user/entity"
	"tracerstudy-auth-service/modules/user/service"
	"tracerstudy-auth-service/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserHandler struct {
	pb.UnimplementedUserServiceServer
	config  config.Config
	userSvc service.UserServiceUseCase
}

func NewUserHandler(config config.Config, userService service.UserServiceUseCase) *UserHandler {
	return &UserHandler{
		config:  config,
		userSvc: userService,
	}
}

func (uh *UserHandler) GetAllUsers(ctx context.Context, req *emptypb.Empty) (*pb.GetAllUsersResponse, error) {
	user, err := uh.userSvc.FindAll(ctx, req)
	if err != nil {
		parseError := errors.ParseError(err)
		log.Println("ERROR: [UserHandler - GetAllUser] Error while get all user: ", parseError.Message)
		// return nil, status.Errorf(parseError.Code, parseError.Message)
		return &pb.GetAllUsersResponse{
			Code:    uint32(http.StatusInternalServerError),
			Message: parseError.Message,
		}, status.Errorf(parseError.Code, parseError.Message)
	}

	var userArr []*pb.User
	for _, u := range user {
		userProto := entity.ConvertEntityToProto(u)
		userArr = append(userArr, userProto)
	}

	return &pb.GetAllUsersResponse{
		Code:    uint32(http.StatusOK),
		Message: "get all user success",
		Data:    userArr,
	}, nil
}

func (uh *UserHandler) GetUserById(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.GetUserResponse, error) {
	user, err := uh.userSvc.FindById(ctx, req.GetId())
	if err != nil {
		if user == nil {
			log.Println("WARNING: [UserHandler - GetUserById] Resource user not found for ID:", req.GetId())
			// return nil, status.Errorf(codes.NotFound, "user not found")
			return &pb.GetUserResponse{
				Code:    uint32(http.StatusNotFound),
				Message: "user not found",
			}, status.Errorf(codes.NotFound, "user not found")
		}
		parseError := errors.ParseError(err)
		log.Println("ERROR: [UserHandler - GetUserById] Internal server error:", parseError.Message)
		// return nil, status.Errorf(parseError.Code, parseError.Message)
		return &pb.GetUserResponse{
			Code:    uint32(http.StatusInternalServerError),
			Message: parseError.Message,
		}, status.Errorf(parseError.Code, parseError.Message)
	}

	userProto := entity.ConvertEntityToProto(user)

	return &pb.GetUserResponse{
		Code:    uint32(http.StatusOK),
		Message: "get user success",
		Data:    userProto,
	}, nil
}

func (uh *UserHandler) CreateUser(ctx context.Context, req *pb.User) (*pb.GetUserResponse, error) {
	user, err := uh.userSvc.Create(ctx, req.GetName(), req.GetUsername(), req.GetEmail(), req.GetPassword(), req.GetRoleId())

	if err != nil {
		parseError := errors.ParseError(err)
		log.Println("ERROR: [UserHandler - CreateUser] Error while create user: ", parseError.Message)
		// return nil, status.Errorf(parseError.Code, parseError.Message)
		return &pb.GetUserResponse{
			Code:    uint32(http.StatusInternalServerError),
			Message: parseError.Message,
		}, status.Errorf(parseError.Code, parseError.Message)
	}

	userProto := entity.ConvertEntityToProto(user)

	return &pb.GetUserResponse{
		Code:    uint32(http.StatusOK),
		Message: "create user success",
		Data:    userProto,
	}, nil
}

func (uh *UserHandler) UpdateUser(ctx context.Context, req *pb.User) (*pb.GetUserResponse, error) {
	userDataUpdate := &entity.User{
		Name:     req.GetName(),
		Username: req.GetUsername(),
		Email:    req.GetEmail(),
		RoleId:   req.GetRoleId(),
	}

	user, err := uh.userSvc.Update(ctx, req.GetId(), userDataUpdate)

	if err != nil {
		parseError := errors.ParseError(err)
		log.Println("ERROR: [UserHandler - UpdateUser] Error while update user: ", parseError.Message)
		// return nil, status.Errorf(parseError.Code, parseError.Message)
		return &pb.GetUserResponse{
			Code:    uint32(http.StatusInternalServerError),
			Message: parseError.Message,
		}, status.Errorf(parseError.Code, parseError.Message)
	}

	userProto := entity.ConvertEntityToProto(user)

	return &pb.GetUserResponse{
		Code:    uint32(http.StatusOK),
		Message: "update user success",
		Data:    userProto,
	}, nil
}

func (uh *UserHandler) DeleteUser(ctx context.Context, req *pb.GetUserByIdRequest) (*pb.DeleteUserResponse, error) {
	err := uh.userSvc.Delete(ctx, req.GetId())
	if err != nil {
		parseError := errors.ParseError(err)
		log.Println("ERROR: [UserHandler - DeleteUser] Internal server error:", parseError.Message)
		// return nil, status.Errorf(parseError.Code, parseError.Message)
		return &pb.DeleteUserResponse{
			Code:    uint32(http.StatusInternalServerError),
			Message: parseError.Message,
		}, status.Errorf(parseError.Code, parseError.Message)
	}

	return &pb.DeleteUserResponse{
		Code:    uint32(http.StatusOK),
		Message: "delete user success",
	}, nil
}
