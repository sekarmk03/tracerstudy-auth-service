package client

import (
	"context"
	"tracerstudy-auth-service/pb"
	"tracerstudy-auth-service/server"
)

type MhsBiodataServiceClient struct {
	Client pb.MhsBiodataServiceClient
}

func BuildMhsBiodataServiceClient(url string) MhsBiodataServiceClient {
	cc := server.InitGRPCConn(url, false, "")

	c := MhsBiodataServiceClient{
		Client: pb.NewMhsBiodataServiceClient(cc),
	}

	return c
}

func (mc *MhsBiodataServiceClient) FetchMhsBiodataByNim(nim string) (*pb.MhsBiodataResponse, error) {
	req := &pb.MhsBiodataRequest{
		Nim: nim,
	}

	return mc.Client.FetchMhsBiodataByNim(context.Background(), req)
}
