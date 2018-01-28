package main

import (
	"io"
	"log"
	"os"
	"strconv"

	pb "is4b/facilitiesReservation/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	svcAddr := "localhost:53002"

	svrPort := os.Getenv("FACILITIES_RESERVATION")
	if len(svrPort) > 1 {
		svcAddr = svrPort
	}

	conn, err := grpc.Dial(svcAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewFacilitiesReservationClient(conn)

	cmd := os.Args[1]

	if cmd == "reserve" {
		meetingRoomId, _ := strconv.ParseInt(os.Args[2], 10, 32)

		req := &pb.ReserveMeetingRoomEvent{
			UserId: int32(0),
			MeetingRoomId: int32(meetingRoomId),
			StartTime: os.Args[3],
			EndTime: os.Args[4],
		}

		reserveMeetingRoom(client, req)

	} else if cmd == "list" {
		meetingRoomId, _ := strconv.ParseInt(os.Args[2], 10, 32)

		req := &pb.ListMeetingRoomReservationsEvent{
			UserId: int32(0),
			MeetingRoomId: int32(meetingRoomId),
		}

		listMeetingRoomReservations(client, req)

	} else if cmd == "find" {
		req := &pb.FindAvailableMeetingRoomsEvent {
			UserId: int32(0),
			StartTime: os.Args[2],
			EndTime: os.Args[3],
		}
		findAvailableMeetingRooms(client, req)
	} else {
		log.Fatalf("unknown command: " + cmd )
	}

}


func reserveMeetingRoom(client pb.FacilitiesReservationClient, evt *pb.ReserveMeetingRoomEvent) {
	confirm, err := client.ReserveMeetingRoom(context.Background(), evt)

	if err != nil {
		log.Fatalf("could not call service: %v", err)
	}

	log.Printf("Confirmation: %s", confirm.StatusMessage)
}


func listMeetingRoomReservations(client pb.FacilitiesReservationClient, evt *pb.ListMeetingRoomReservationsEvent) {
	stream, err := client.ListMeetingRoomReservations(context.Background(), evt)

	if err != nil {
		log.Fatalf("%v.ListMeetingRoomReservations(_) = _, %v", client, err)
	}

	for {
		building, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("%v.ListMeetingRoomReservations(_) = _, %v", client, err)
		}

		log.Println(building)
	}
}


func findAvailableMeetingRooms(client pb.FacilitiesReservationClient, evt *pb.FindAvailableMeetingRoomsEvent) {
	stream, err := client.FindAvailableMeetingRooms(context.Background(), evt)

	if err != nil {
		log.Fatalf("%v.FindAvailableMeetingRooms(_) = _, %v", client, err)
	}

	for {
		item, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("%v.FindAvailableMeetingRooms(_) = _, %v", client, err)
		}

		log.Println(item)
	}
}
