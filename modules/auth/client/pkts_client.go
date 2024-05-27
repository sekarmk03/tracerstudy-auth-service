package client

import (
	"context"
	"tracerstudy-auth-service/pb"
	"tracerstudy-auth-service/server"
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

func (c *PktsServiceClient) GetNimByDataAtasan(nama, email, hp string) (*pb.GetNimByDataAtasanResponse, error) {
	req := &pb.GetNimByDataAtasanRequest{
		NamaAtasan:  nama,
		EmailAtasan: email,
		HpAtasan:    hp,
	}

	return c.Client.GetNimByDataAtasan(context.Background(), req)
}
