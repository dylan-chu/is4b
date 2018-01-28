package main

import (
	"log"
	"net"
	"os"

	pb "is4b/facilitiesReservation/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	data "is4b/facilitiesReservation/data"
)

type FacilitiesReservationServer struct {
	mtgRmRsvDAO data.MeetingRoomReservationDAO
}

func main() {
	svcAddr := ":53002"
	if len(os.Args) > 1 {
		svcAddr = os.Args[1]
	}

	lis, err := net.Listen("tcp", svcAddr)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer() //NewServer(opts...)
	pb.RegisterFacilitiesReservationServer(grpcServer, newServer())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func newServer() *FacilitiesReservationServer {
	connStr := "root:mysql1234@tcp(127.0.0.1:3306)/IS4B"
	mrrDAO := data.MeetingRoomReservationDAO{ connStr }
	server := &FacilitiesReservationServer{
		mtgRmRsvDAO: mrrDAO,
	}
	return server
}

func (s *FacilitiesReservationServer) ReserveMeetingRoom(ctx context.Context, evt *pb.ReserveMeetingRoomEvent) (*pb.ConfirmationEvent, error) {
	item := &pb.MeetingRoomReservation{
		Id: int32(0),
		UserId: evt.UserId,
		MeetingRoomId: evt.MeetingRoomId,
		StartTime: evt.StartTime,
		EndTime: evt.EndTime,
	}

	err := s.mtgRmRsvDAO.Create(item)
	var returnCode = 0
	var statusMsg = "Reservation Successful"
	if err != nil {
		returnCode = -1
		statusMsg = err.Error()
	}
	return &pb.ConfirmationEvent{ReturnCode: int32(returnCode), StatusMessage: statusMsg}, err
}

func (s *FacilitiesReservationServer) ListMeetingRoomReservations(evt *pb.ListMeetingRoomReservationsEvent,
	stream pb.FacilitiesReservation_ListMeetingRoomReservationsServer) error {
	items := s.mtgRmRsvDAO.GetByMeetingRoom(evt.MeetingRoomId)

	for _, item := range items {
		if err := stream.Send(&item); err != nil {
			return err
		}
	}

	return nil
}

func (s *FacilitiesReservationServer) FindAvailableMeetingRooms(evt *pb.FindAvailableMeetingRoomsEvent,
	stream pb.FacilitiesReservation_FindAvailableMeetingRoomsServer) error {
	items := s.mtgRmRsvDAO.GetByReservationTime(evt.StartTime, evt.EndTime)

	for _, item := range items {
		if err := stream.Send(&item); err != nil {
			return err
		}
	}

	return nil
}
