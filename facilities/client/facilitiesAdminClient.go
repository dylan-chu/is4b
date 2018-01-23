package main

import (
	"io"
	"log"
	"os"
	"strconv"

	pb "is4b/facilities/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	svcAddr := "localhost:53001"

	svrPort := os.Getenv("FACILITIES")
	if len(svrPort) > 1 {
			svcAddr = svrPort
	}

	conn, err := grpc.Dial(svcAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewFacilitiesAdminClient(conn)

	cmd := os.Args[1]
	itemType := os.Args[2]

	if cmd == "add" {
		if itemType == "building" {
			addBuilding(client, os.Args[3])
		} else if itemType == "meetingroom" {
			blgdId, _ := strconv.ParseInt(os.Args[3], 10, 32)
			flr, _ := strconv.ParseInt(os.Args[4], 10, 32)
			cpcty, _ := strconv.ParseInt(os.Args[6], 10, 32)

			meetingRoom := &pb.MeetingRoom{
				Id: int32(0),
				BuildingId: int32(blgdId),
				Floor: int32(flr),
				Name: os.Args[5],
				Capacity: int32(cpcty),
				HasProjector: os.Args[7] == "Y",
				HasWhiteboard: os.Args[8] == "Y",
				HasConferenceLine: os.Args[9] == "Y",
			}
			addMeetingRoom(client, meetingRoom)
		}
	} else if cmd == "list" {
		if itemType == "building" {
			listBuildings(client)
		} else if itemType == "meetingroom" {
			listMeetingRooms(client)
		}
	} else {
		log.Fatalf("unknown command: " + cmd )
	}

}

func addBuilding(client pb.FacilitiesAdminClient, buildingName string) {
	confirm, err := client.AddBuilding(context.Background(), &pb.AddBuidingEvent{
		UserId: 1,
		BuildingName: buildingName,
	})

	if err != nil {
		log.Fatalf("could not call service: %v", err)
	}

	log.Printf("Confirmation: %s", confirm.StatusMessage)
}

func listBuildings(client pb.FacilitiesAdminClient) {
	stream, err := client.ListBuildings(context.Background(), &pb.ListBuildingsEvent{
		UserId: 1,
	})

	if err != nil {
		log.Fatalf("%v.ListBuildings(_) = _, %v", client, err)
	}

	for {
		building, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("%v.ListBuildings(_) = _, %v", client, err)
		}

		log.Println(building)
	}
}

func addMeetingRoom(client pb.FacilitiesAdminClient, meetingRoom *pb.MeetingRoom) {
	confirm, err := client.AddMeetingRoom(context.Background(), &pb.AddMeetingRoomEvent{
		UserId: 1,
		MeetingRoom: meetingRoom,
	})

	if err != nil {
		log.Fatalf("could not call service: %v", err)
	}

	log.Printf("Confirmation: %s", confirm.StatusMessage)
}

func listMeetingRooms(client pb.FacilitiesAdminClient) {
	stream, err := client.ListMeetingRooms(context.Background(), &pb.ListMeetingRoomsEvent{
		UserId: 1,
	})

	if err != nil {
		log.Fatalf("%v.ListMeetingRooms(_) = _, %v", client, err)
	}

	for {
		item, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("%v.ListMeetingRooms(_) = _, %v", client, err)
		}

		log.Println(item)
	}
}
