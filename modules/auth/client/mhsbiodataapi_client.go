package client

import (
	"context"
	"tracerstudy-auth-service/pb"
	"tracerstudy-auth-service/server"
)

type MhsBiodataApiServiceClient struct {
	Client pb.MhsBiodataApiServiceClient
}

func BuildMhsBiodataServiceClient(url string) MhsBiodataApiServiceClient {
	cc := server.InitGRPCConn(url, false, "")

	c := MhsBiodataApiServiceClient{
		Client: pb.NewMhsBiodataApiServiceClient(cc),
	}

	return c
}

func (mc *MhsBiodataApiServiceClient) FetchMhsBiodataByNim(nim string) (*pb.MhsBiodataApiResponse, error) {
	req := &pb.MhsBiodataApiRequest{
		Nim: nim,
	}

	return mc.Client.FetchMhsBiodataByNim(context.Background(), req)
}

func (mc *MhsBiodataApiServiceClient) CheckMhsAlumni(nim string, tglSidang string) (*pb.CheckMhsAlumniResponse, error) {
	req := &pb.CheckMhsAlumniRequest{
		Nim: nim,
		TglSidang: tglSidang,
	}

	return mc.Client.CheckMhsAlumni(context.Background(), req)
}
