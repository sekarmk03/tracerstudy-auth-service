package client

import (
	"context"
	"tracerstudy-auth-service/pb"
	"tracerstudy-auth-service/server"

	"google.golang.org/protobuf/types/known/emptypb"
)

type PktsServiceClient struct {
	Client pb.PKTSServiceClient
}

func BuildPktsServiceClient(url string) PktsServiceClient {
	cc := server.InitGRPCConn(url, false, "")

	c := PktsServiceClient{
		Client: pb.NewPKTSServiceClient(cc),
	}

	return c
}

func (c *PktsServiceClient) GetAllPkts() (*pb.GetAllPKTSResponse, error) {
	return c.Client.GetAllPKTS(context.Background(), &emptypb.Empty{})
}

func (c *PktsServiceClient) GetNimByDataAtasan(nama, email, hp string) (*pb.GetNimByDataAtasanResponse, error) {
	req := &pb.GetNimByDataAtasanRequest{
		NamaAtasan:  nama,
		EmailAtasan: email,
		HpAtasan:    hp,
	}

	return c.Client.GetNimByDataAtasan(context.Background(), req)
}
