package handler

import (
	"context"
	"log"
	"net/http"
	"tracerstudy-auth-service/common/config"
	"tracerstudy-auth-service/common/errors"
	commonJwt "tracerstudy-auth-service/common/jwt"
	"tracerstudy-auth-service/common/utils"

	"tracerstudy-auth-service/modules/auth/client"
	userSvc "tracerstudy-auth-service/modules/user/service"
	"tracerstudy-auth-service/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	config     config.Config
	userSvc    userSvc.UserServiceUseCase
	jwtManager *commonJwt.JWT
	pktsSvc    client.PktsServiceClient
	mhsSvc     client.MhsBiodataServiceClient
}

func NewAuthHandler(
	config config.Config,
	userService userSvc.UserServiceUseCase,
	jwtManager *commonJwt.JWT,
	pktsService client.PktsServiceClient,
	mhsService client.MhsBiodataServiceClient,
) *AuthHandler {
	return &AuthHandler{
		config:     config,
		userSvc:    userService,
		jwtManager: jwtManager,
		pktsSvc:    pktsService,
		mhsSvc:     mhsService,
	}
}

func (ah *AuthHandler) LoginAlumni(ctx context.Context, req *pb.LoginAlumniRequest) (*pb.LoginResponse, error) {
	mhs, err := ah.mhsSvc.FetchMhsBiodataByNim(req.GetNim())
	if err != nil {
		if mhs == nil {
			log.Println("WARNING: [AuthHandler - Login] Mhs resource not found")
			return nil, status.Errorf(codes.NotFound, "mhs resource not found")
		}
		parseError := errors.ParseError(err)
		log.Println("ERROR: [AuthHandler - Login] Error while fetching mhs biodata:", parseError.Message)
		return nil, status.Errorf(parseError.Code, parseError.Message)
	}

	if mhs == nil {
		log.Println("WARNING: [AuthHandler - Login] Mhs resource not found")
		return nil, status.Errorf(codes.NotFound, "mhs resource not found")
	}

	// generate token with role 6 = alumni
	token, err := ah.jwtManager.GenerateToken(mhs.GetData().NIM, 6)
	if err != nil {
		log.Println("ERROR: [AuthHandler - Login] Error while generating token:", err)
		return nil, status.Errorf(codes.Internal, "token failed to generate: %v", err)
	}

	return &pb.LoginResponse{
		Code:    uint32(http.StatusOK),
		Message: "login success",
		Token:   token,
	}, nil
}

func (ah *AuthHandler) LoginUserStudy(ctx context.Context, req *pb.LoginUserStudyRequest) (*pb.LoginResponse, error) {
	user, err := ah.pktsSvc.GetNimByDataAtasan(req.GetNamaAtasan(), req.GetEmailAtasan(), req.GetHpAtasan())
	if err != nil {
		if user.GetNims() == nil || len(user.GetNims()) == 0 {
			log.Println("WARNING: [AuthHandler - LoginUserStudy] User resource not found")
			return nil, status.Errorf(codes.NotFound, "user resource not found")
		}
		parseError := errors.ParseError(err)
		log.Println("ERROR: [AuthHandler - LoginUserStudy] Error while fetching user:", parseError.Message)
		return nil, status.Errorf(parseError.Code, parseError.Message)
	}

	if user.GetNims() == nil || len(user.GetNims()) == 0 {
		log.Println("WARNING: [AuthHandler - LoginUserStudy] User resource not found")
		return nil, status.Errorf(codes.NotFound, "user resource not found")
	}

	// generate token with role 7 = user study
	token, err := ah.jwtManager.GenerateToken(req.GetEmailAtasan(), 7)
	if err != nil {
		log.Println("ERROR: [AuthHandler - LoginUserStudy] Error while generating token:", err)
		return nil, status.Errorf(codes.Internal, "token failed to generate: %v", err)
	}

	return &pb.LoginResponse{
		Code:    uint32(http.StatusOK),
		Message: "login user study success",
		Token:   token,
	}, nil
}

func (ah *AuthHandler) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginResponse, error) {
	user, err := ah.userSvc.FindByUsername(ctx, req.GetUsername())
	if err != nil {
		if user == nil {
			log.Println("WARNING: [AuthHandler - LoginUser] User resource not found")
			return nil, status.Errorf(codes.NotFound, "user resource not found")
		}
		parseError := errors.ParseError(err)
		log.Println("ERROR: [AuthHandler - LoginUser] Error while fetching user:", parseError.Message)
		return nil, status.Errorf(parseError.Code, parseError.Message)
	}

	if user == nil {
		log.Println("WARNING: [AuthHandler - LoginUser] User resource not found")
		return nil, status.Errorf(codes.NotFound, "user resource not found")
	}

	match := utils.CheckPasswordHash(req.GetPassword(), user.Password)
	if !match {
		log.Println("WARNING: [AuthHandler - LoginUser] Invalid credentials")
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}

	token, err := ah.jwtManager.GenerateToken(user.Username, user.RoleId)
	if err != nil {
		log.Println("ERROR: [AuthHandler - LoginUser] Error while generating token:", err)
		return nil, status.Errorf(codes.Internal, "token failed to generate: %v", err)
	}

	return &pb.LoginResponse{
		Code:    uint32(http.StatusOK),
		Message: "login user success",
		Token:   token,
	}, nil
}

func (ah *AuthHandler) RegisterUser(ctx context.Context, req *pb.User) (*pb.RegisterUserResponse, error) {
	user, err := ah.userSvc.FindByUsername(ctx, req.GetUsername())
	if err != nil {
		if user != nil {
			log.Println("WARNING: [AuthHandler - RegisterUser] User already exists")
			return nil, status.Errorf(codes.AlreadyExists, "user already exist")
		}
		parseError := errors.ParseError(err)
		if parseError.Code != codes.NotFound {
			log.Println("ERROR: [AuthHandler - RegisterUser] Error while fetching user:", parseError.Message)
			return nil, status.Errorf(parseError.Code, parseError.Message)
		}
	}

	if user != nil {
		log.Println("WARNING: [AuthHandler - RegisterUser] User already exists")
		return nil, status.Errorf(codes.AlreadyExists, "user already exists")
	}

	user, err = ah.userSvc.Create(ctx, req.GetName(), req.GetUsername(), req.GetEmail(), req.GetPassword(), req.GetRoleId())
	if err != nil {
		parseError := errors.ParseError(err)
		log.Println("ERROR: [AuthHandler - RegisterUser] Error while creating user:", parseError.Message)
		return nil, status.Errorf(parseError.Code, parseError.Message)
	}

	userProto := &pb.User{
		Id:       user.Id,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		RoleId:   user.RoleId,
	}

	return &pb.RegisterUserResponse{
		Code:    uint32(http.StatusCreated),
		Message: "register user success",
		Data:    userProto,
	}, nil
}