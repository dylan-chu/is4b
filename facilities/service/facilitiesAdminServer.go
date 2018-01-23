package main

import (
	"log"
	"net"
	"os"

	pb "is4b/facilities/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	data "is4b/facilities/data"
)

type FacilitiesAdminServer struct {
	bldgDAO data.BuildingDAO
	mtgRmDAO data.MeetingRoomDAO
}

func main() {
	svcAddr := ":53001"
	if len(os.Args) > 1 {
		svcAddr = os.Args[1]
	}

	lis, err := net.Listen("tcp", svcAddr)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer() //NewServer(opts...)
	pb.RegisterFacilitiesAdminServer(grpcServer, newServer())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func newServer() *FacilitiesAdminServer {
	connStr := "root:mysql1234@tcp(127.0.0.1:3306)/IS4B"
	bgDAO := data.BuildingDAO{ connStr }
	mrDAO := data.MeetingRoomDAO{ connStr }
	server := &FacilitiesAdminServer{
		bldgDAO: bgDAO,
		mtgRmDAO: mrDAO,
	}
	return server
}

func (s *FacilitiesAdminServer) AddBuilding(ctx context.Context, evt *pb.AddBuidingEvent) (*pb.ConfirmationEvent, error) {
	item := &pb.Building{
		Id: int32(0),
		Name: evt.BuildingName,
	}

	err := s.bldgDAO.Create(item)
	var returnCode = 0
	var statusMsg = "Building Added"
	if err != nil {
		returnCode = -1
		statusMsg = "error"
	}
	return &pb.ConfirmationEvent{ReturnCode: int32(returnCode), StatusMessage: statusMsg}, err
}

func (s *FacilitiesAdminServer) ListBuildings(evt *pb.ListBuildingsEvent, stream pb.FacilitiesAdmin_ListBuildingsServer) error {
	items := s.bldgDAO.GetAll()

	for _, item := range items {
		if err := stream.Send(&item); err != nil {
			return err
		}
	}

	return nil
}

func (s *FacilitiesAdminServer) AddMeetingRoom(ctx context.Context, evt *pb.AddMeetingRoomEvent) (*pb.ConfirmationEvent, error) {
	item := evt.MeetingRoom

	err := s.mtgRmDAO.Create(item)
	var returnCode = 0
	var statusMsg = "Meeting Room Added"
	if err != nil {
		returnCode = -1
		statusMsg = "error"
	}
	return &pb.ConfirmationEvent{ReturnCode: int32(returnCode), StatusMessage: statusMsg}, nil
}

func (s *FacilitiesAdminServer) ListMeetingRooms(evt *pb.ListMeetingRoomsEvent, stream pb.FacilitiesAdmin_ListMeetingRoomsServer) error {
	items := s.mtgRmDAO.GetAll()

	for _, item := range items {
		if err := stream.Send(&item); err != nil {
			return err
		}
	}

	return nil
}
