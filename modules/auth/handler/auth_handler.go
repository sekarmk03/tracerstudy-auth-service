package handler

import (
	"context"
	"log"
	"net/http"
	"strings"
	"tracerstudy-auth-service/common/config"
	"tracerstudy-auth-service/common/errors"
	commonJwt "tracerstudy-auth-service/common/jwt"
	"tracerstudy-auth-service/common/utils"

	"tracerstudy-auth-service/modules/auth/client"
	userSvc "tracerstudy-auth-service/modules/user/service"
	"tracerstudy-auth-service/pb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	config     config.Config
	userSvc    userSvc.UserServiceUseCase
	jwtManager *commonJwt.JWT
	pktsSvc    client.PktsServiceClient
	mhsApiSvc  client.MhsBiodataApiServiceClient
}

func NewAuthHandler(
	config config.Config,
	userService userSvc.UserServiceUseCase,
	jwtManager *commonJwt.JWT,
	pktsService client.PktsServiceClient,
	mhsApiService client.MhsBiodataApiServiceClient,
) *AuthHandler {
	return &AuthHandler{
		config:     config,
		userSvc:    userService,
		jwtManager: jwtManager,
		pktsSvc:    pktsService,
		mhsApiSvc:  mhsApiService,
	}
}

func (ah *AuthHandler) LoginAlumni(ctx context.Context, req *pb.LoginAlumniRequest) (*pb.LoginResponse, error) {
	res, err := ah.mhsApiSvc.CheckMhsAlumni(req.GetNim())
	if err != nil {
		parseError := errors.ParseError(err)
		log.Println("ERROR: [AuthHandler - LoginAlumni] Error checking mhs biodata:", parseError.Message)
		// return nil, status.Errorf(parseError.Code, parseError.Message)
		return &pb.LoginResponse{
			Code:    uint32(http.StatusInternalServerError),
			Message: parseError.Message,
		}, status.Errorf(parseError.Code, parseError.Message)
	}

	if !res.GetIsAlumni() {
		log.Println("WARNING: [AuthHandler - LoginAlumni] Mahasiswa is not an alumni")
		// return nil, status.Errorf(codes.PermissionDenied, "mahasiswa is not an alumni")
		return &pb.LoginResponse{
			Code:    uint32(http.StatusForbidden),
			Message: "mahasiswa is not an alumni",
		}, status.Errorf(codes.PermissionDenied, "mahasiswa is not an alumni")
	}

	// generate token with cred = nim, role = 6 (alumni)
	token, err := ah.jwtManager.GenerateToken(req.GetNim(), 6)

	if err != nil {
		parseError := errors.ParseError(err)
		log.Println("ERROR: [AuthHandler - LoginAlumni] Error while generating token:", parseError.Message)
		// return nil, status.Errorf(codes.Internal, "token failed to generate: %v", parseError.Message)
		return &pb.LoginResponse{
			Code:    uint32(http.StatusInternalServerError),
			Message: "token failed to generate: " + parseError.Message,
		}, status.Errorf(codes.Internal, "token failed to generate: %v", parseError.Message)
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
			// return nil, status.Errorf(codes.NotFound, "user resource not found")
			return &pb.LoginResponse{
				Code:    uint32(http.StatusNotFound),
				Message: "user resource not found",
			}, status.Errorf(codes.NotFound, "user resource not found")
		}
		parseError := errors.ParseError(err)
		log.Println("ERROR: [AuthHandler - LoginUserStudy] Error while fetching user:", parseError.Message)
		// return nil, status.Errorf(parseError.Code, parseError.Message)
		return &pb.LoginResponse{
			Code:    uint32(http.StatusInternalServerError),
			Message: parseError.Message,
		}, status.Errorf(parseError.Code, parseError.Message)
	}

	if user.GetNims() == nil || len(user.GetNims()) == 0 {
		log.Println("WARNING: [AuthHandler - LoginUserStudy] User resource not found")
		// return nil, status.Errorf(codes.NotFound, "user resource not found")
		return &pb.LoginResponse{
			Code:    uint32(http.StatusNotFound),
			Message: "user resource not found",
		}, status.Errorf(codes.NotFound, "user resource not found")
	}

	// generate token with cred = email, role = 7 (pengguna alumni)
	token, err := ah.jwtManager.GenerateToken(req.GetEmailAtasan(), 7)

	if err != nil {
		parseError := errors.ParseError(err)
		log.Println("ERROR: [AuthHandler - LoginUserStudy] Error while generating token:", parseError.Message)
		// return nil, status.Errorf(codes.Internal, "token failed to generate: %v", parseError.Message)
		return &pb.LoginResponse{
			Code:    uint32(http.StatusInternalServerError),
			Message: "token failed to generate: " + parseError.Message,
		}, status.Errorf(codes.Internal, "token failed to generate: %v", parseError.Message)
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
			log.Println("WARNING: [AuthHandler - LoginUser] User not found")
			// return nil, status.Errorf(codes.NotFound, "user not found")
			return &pb.LoginResponse{
				Code:    uint32(http.StatusNotFound),
				Message: "user not found",
			}, status.Errorf(codes.NotFound, "user not found")
		}
		parseError := errors.ParseError(err)
		log.Println("ERROR: [AuthHandler - LoginUser] Error while fetching user:", parseError.Message)
		// return nil, status.Errorf(parseError.Code, parseError.Message)
		return &pb.LoginResponse{
			Code:    uint32(http.StatusInternalServerError),
			Message: parseError.Message,
		}, status.Errorf(parseError.Code, parseError.Message)
	}

	if user == nil {
		log.Println("WARNING: [AuthHandler - LoginUser] User resource not found")
		// return nil, status.Errorf(codes.NotFound, "user resource not found")
		return &pb.LoginResponse{
			Code:    uint32(http.StatusNotFound),
			Message: "user not found",
		}, status.Errorf(codes.NotFound, "user not found")
	}

	match := utils.CheckPasswordHash(req.GetPassword(), user.Password)
	if !match {
		log.Println("WARNING: [AuthHandler - LoginUser] Invalid credentials")
		// return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
		return &pb.LoginResponse{
			Code:    uint32(http.StatusBadRequest),
			Message: "invalid credentials",
		}, status.Errorf(codes.InvalidArgument, "invalid credentials")
	}

	// generate token with cred = username, role = roleId
	token, err := ah.jwtManager.GenerateToken(user.Username, user.RoleId)

	if err != nil {
		parseError := errors.ParseError(err)
		log.Println("ERROR: [AuthHandler - LoginUser] Error while generating token:", parseError.Message)
		// return nil, status.Errorf(codes.Internal, "token failed to generate: %v", parseError.Message)
		return &pb.LoginResponse{
			Code:    uint32(http.StatusInternalServerError),
			Message: "token failed to generate: " + parseError.Message,
		}, status.Errorf(codes.Internal, "token failed to generate: %v", parseError.Message)
	}

	return &pb.LoginResponse{
		Code:    uint32(http.StatusOK),
		Message: "login user success",
		Token:   token,
	}, nil
}

func (ah *AuthHandler) RegisterUser(ctx context.Context, req *pb.User) (*pb.SingleUserResponse, error) {
	user, err := ah.userSvc.FindByUsername(ctx, req.GetUsername())
	if err != nil {
		if user != nil {
			log.Println("WARNING: [AuthHandler - RegisterUser] User already exists")
			// return nil, status.Errorf(codes.AlreadyExists, "user already exist")
			return &pb.SingleUserResponse{
				Code:    uint32(http.StatusConflict),
				Message: "user already exists",
			}, status.Errorf(codes.AlreadyExists, "user already exists")
		}
		parseError := errors.ParseError(err)
		if parseError.Code != codes.NotFound {
			log.Println("ERROR: [AuthHandler - RegisterUser] Error while fetching user:", parseError.Message)
			// return nil, status.Errorf(parseError.Code, parseError.Message)
			return &pb.SingleUserResponse{
				Code:    uint32(http.StatusInternalServerError),
				Message: parseError.Message,
			}, status.Errorf(parseError.Code, parseError.Message)
		}
	}

	if user != nil {
		log.Println("WARNING: [AuthHandler - RegisterUser] User already exists")
		// return nil, status.Errorf(codes.AlreadyExists, "user already exists")
		return &pb.SingleUserResponse{
			Code:    uint32(http.StatusConflict),
			Message: "user already exists",
		}, status.Errorf(codes.AlreadyExists, "user already exists")
	}

	user, err = ah.userSvc.Create(ctx, req.GetName(), req.GetUsername(), req.GetEmail(), req.GetPassword(), req.GetRoleId())
	if err != nil {
		parseError := errors.ParseError(err)
		log.Println("ERROR: [AuthHandler - RegisterUser] Error while creating user:", parseError.Message)
		// return nil, status.Errorf(parseError.Code, parseError.Message)
		return &pb.SingleUserResponse{
			Code:    uint32(http.StatusInternalServerError),
			Message: parseError.Message,
		}, status.Errorf(parseError.Code, parseError.Message)
	}

	userProto := &pb.User{
		Id:        user.Id,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		RoleId:    user.RoleId,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}

	return &pb.SingleUserResponse{
		Code:    uint32(http.StatusCreated),
		Message: "register user success",
		Data:    userProto,
	}, nil
}

func (ah *AuthHandler) GetCurrentUser(ctx context.Context, req *emptypb.Empty) (*pb.SingleUserResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Println("ERROR: [AuthHandler - GetCurrentUser] No metadata found")
		return &pb.SingleUserResponse{
			Code:    uint32(http.StatusBadRequest),
			Message: "no metadata found",
		}, status.Errorf(codes.InvalidArgument, "no metadata found")
	}

	values, ok := md["authorization"]
	if !ok || len(values) == 0 {
		log.Println("ERROR: [AuthHandler - GetCurrentUser] No authorization header found")
		return &pb.SingleUserResponse{
			Code:    uint32(http.StatusBadRequest),
			Message: "no authorization header found",
		}, status.Errorf(codes.InvalidArgument, "no authorization header found")
	}

	authHeader := values[0]
	parts := strings.Fields(authHeader)
	if len(parts) != 2 || parts[0] != "Bearer" {
		log.Println("ERROR: [AuthHandler - GetCurrentUser] Invalid authorization header")
		return &pb.SingleUserResponse{
			Code:    uint32(http.StatusBadRequest),
			Message: "invalid authorization header",
		}, status.Errorf(codes.InvalidArgument, "invalid authorization header")
	}

	accessToken := parts[1]
	claims, err := ah.jwtManager.Verify(accessToken)
	if err != nil {
		log.Println("ERROR: [AuthHandler - GetCurrentUser] Invalid token:", err)
		return &pb.SingleUserResponse{
			Code:    uint32(http.StatusUnauthorized),
			Message: "invalid token",
		}, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	user, err := ah.userSvc.FindByUsername(ctx, claims.Cred)
	if err != nil {
		if user == nil {
			log.Println("WARNING: [AuthHandler - GetCurrentUser] User not found")
			return &pb.SingleUserResponse{
				Code:    uint32(http.StatusNotFound),
				Message: "user not found",
			}, status.Errorf(codes.NotFound, "user not found")
		}
		parseError := errors.ParseError(err)
		log.Println("ERROR: [AuthHandler - GetCurrentUser] Error while fetching user:", parseError.Message)
		return &pb.SingleUserResponse{
			Code:    uint32(http.StatusInternalServerError),
			Message: parseError.Message,
		}, status.Errorf(parseError.Code, parseError.Message)
	}

	userProto := &pb.User{
		Id:       user.Id,
		Name:     user.Name,
		Username: user.Username,
		Email:    user.Email,
		RoleId:   user.RoleId,
	}

	return &pb.SingleUserResponse{
		Code:    uint32(http.StatusOK),
		Message: "get current user success",
		Data:    userProto,
	}, nil
}
