package main

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	UnimplementedLocalVendorServer
	vendor *Vendor
}

func (s *server) CreateVolume(ctx context.Context, req *VolumeRequest) (*VolumeReply, error) {
	err := s.vendor.CreateVolume(req.Name, req.Size)
	if err != nil {
		return nil, err
	}
	return &VolumeReply{Message: "created"}, nil
}

func (s *server) GetVolumeInfo(ctx context.Context, req *VolumeRequest) (*VolumeInfo, error) {
	meta, err := s.vendor.GetVolumeInfo(req.Name)
	if err != nil {
		return nil, err
	}
	return &VolumeInfo{
		Status:   meta["status"].(string),
		Size:     int64(meta["size"].(float64)),
		ServerIp: fmt.Sprintf("%v", meta["server_ip"]),
	}, nil
}

func (s *server) DeleteVolume(ctx context.Context, req *VolumeRequest) (*VolumeReply, error) {
	err := s.vendor.DeleteVolume(req.Name)
	if err != nil {
		return nil, err
	}
	return &VolumeReply{Message: "deleted"}, nil
}

func runGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	vendor := NewVendor("/mnt/localstorage/volumes")
	RegisterLocalVendorServer(grpcServer, &server{vendor: vendor})
	fmt.Println("gRPC server listening on :50051")
	grpcServer.Serve(lis)
}
