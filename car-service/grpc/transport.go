package grpc

import (
	"context"

	"dmitrii.com/car-service/car"
	"dmitrii.com/car-service/pb"
	"dmitrii.com/car-service/transport"

	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	pb.UnimplementedCarServiceServer
	create   grpctransport.Handler
	getByID  grpctransport.Handler
	markSold grpctransport.Handler
}

func NewGRPCServer(eps transport.Endpoints) pb.CarServiceServer {
	return &grpcServer{
		create: grpctransport.NewServer(
			eps.Create,
			decodeGRPCCreateRequest,
			encodeGRPCCreateResponse,
		),
		getByID: grpctransport.NewServer(
			eps.GetByID,
			decodeGRPCGetByIDRequest,
			encodeGRPCGetByIDResponse,
		),
		markSold: grpctransport.NewServer(
			eps.MarkSold,
			decodeGRPCMarkSoldRequest,
			encodeGRPCMarkSoldResponse,
		),
	}
}

func (s *grpcServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	_, rep, err := s.create.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.CreateResponse), nil
}

func (s *grpcServer) GetByID(ctx context.Context, req *pb.GetByIDRequest) (*pb.GetByIDResponse, error) {
	_, rep, err := s.getByID.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetByIDResponse), nil
}

func (s *grpcServer) MarkSold(ctx context.Context, req *pb.MarkSoldRequest) (*pb.MarkSoldResponse, error) {
	_, rep, err := s.markSold.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.MarkSoldResponse), nil
}

// --- Декодеры (Protobuf -> Endpoint Request) ---

func decodeGRPCCreateRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CreateRequest)
	var car car.Car
	if req.Car != nil {
		car.ID = req.Car.Id
		car.Brand = req.Car.Brand
		car.Model = req.Car.Model
		car.Year = int64(req.Car.Year)
		car.Price = req.Car.Price
		car.Sold = req.Car.Sold
	}
	return transport.CreateRequest{Car: car}, nil
}

func decodeGRPCGetByIDRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetByIDRequest)
	return transport.GetByIDRequest{ID: req.Id}, nil
}

func decodeGRPCMarkSoldRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.MarkSoldRequest)
	return transport.MarkSoldRequest{ID: req.Id}, nil
}

// --- Энкодеры (Endpoint Response -> Protobuf) ---

func encodeGRPCCreateResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(transport.CreateResponse)
	return &pb.CreateResponse{
		Id:  resp.ID,
		Err: err2str(resp.Err),
	}, nil
}

func encodeGRPCGetByIDResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(transport.GetByIDResponse)
	return &pb.GetByIDResponse{
		Car: &pb.Car{
			Id:    resp.Car.ID,
			Brand: resp.Car.Brand,
			Model: resp.Car.Model,
			Year:  int32(resp.Car.Year),
			Price: resp.Car.Price,
			Sold:  resp.Car.Sold,
		},
		Err: err2str(resp.Err),
	}, nil
}

func encodeGRPCMarkSoldResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(transport.MarkSoldResponse)
	return &pb.MarkSoldResponse{
		Err: err2str(resp.Err),
	}, nil
}

// Вспомогательная функция для безопасной конвертации ошибки в строку
func err2str(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
