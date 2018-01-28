package data

import (
  "log"

  pb "is4b/facilitiesReservation/proto"
  _ "github.com/go-sql-driver/mysql"
  "time"
  "errors"
)

const standardTimeLayout = "2006-01-02 15:04:05"

type MeetingRoomReservationDAO struct {
    ConnStr string
}

func (dao MeetingRoomReservationDAO) Create(item *pb.MeetingRoomReservation) error {
  startTime, err := time.Parse(standardTimeLayout, item.StartTime)
  endTime, err := time.Parse(standardTimeLayout, item.EndTime)

  if endTime.Before(startTime) {
    return errors.New("end Time is before start Time")
  }

  db := getDbConn(dao.ConnStr)
  defer db.Close()

  // TODO add timeReserved to query
  rows, err := db.Query("SELECT startTime, endTime from MEETING_ROOM_RESERVATION where meetingRoomId=?", item.MeetingRoomId)
  if err != nil {
    log.Fatal(err)
  }

  defer rows.Close()
  hasConflict := false

  for rows.Next() {
    var sTime, eTime time.Time
    err := rows.Scan(&sTime, &eTime)
    if err != nil {
      log.Fatal(err)
    }

    if startTime.Before(sTime) {
      if endTime.After(sTime) {
        hasConflict = true
      }
    } else if startTime.Before(eTime) {
      hasConflict = true
    }
  }

  err = rows.Err()
  if err != nil {
    log.Fatal(err)
  }

  if hasConflict {
    return errors.New("conflict with reservation time")
  }

  query := "INSERT into MEETING_ROOM_RESERVATION (userId, meetingRoomId, startTime, endTime) VALUES (?, ?, ?, ?)"
  stmt, err := db.Prepare(query)
  defer stmt.Close()
  if err != nil {
    return err
  }

  result, err := stmt.Exec(item.UserId, item.MeetingRoomId, startTime, endTime)
  if err != nil {
    return err
  }

  rowId, err := result.LastInsertId()
  if err != nil {
      return err
  }

  item.Id = int32(rowId)
  return nil
}


func (dao MeetingRoomReservationDAO) GetByMeetingRoom(meetingRoomId int32) []pb.MeetingRoomReservation {
  db := getDbConn(dao.ConnStr)
  defer db.Close()

  items := make([]pb.MeetingRoomReservation, 0)

  rows, err := db.Query("SELECT id, userId, meetingRoomId, startTime, endTime from MEETING_ROOM_RESERVATION where meetingRoomId=?", meetingRoomId)
  if err != nil {
    log.Fatal(err)
  }

  defer rows.Close()

  for rows.Next() {
    var item pb.MeetingRoomReservation
    err := rows.Scan(&item.Id, &item.UserId, &item.MeetingRoomId, &item.StartTime, &item.EndTime)
    if err != nil {
      log.Fatal(err)
    }

    items = append(items, item)
  }

  err = rows.Err()
  if err != nil {
    log.Fatal(err)
  }

  return items
}


func (dao MeetingRoomReservationDAO) GetByReservationTime(startTimeStr string, endTimeStr string) []pb.MeetingRoom {
  startTime, err := time.Parse(standardTimeLayout, startTimeStr)
  if err != nil {
    log.Fatalf("invalid start time:%s", startTimeStr)
  }

  endTime, err := time.Parse(standardTimeLayout, endTimeStr)
  if err != nil {
    log.Fatalf("invalid end time:%s", endTimeStr)
  }

  if endTime.Before(startTime) {
    log.Fatal("end time is before start time")
  }

  db := getDbConn(dao.ConnStr)
  defer db.Close()

  items := make(map[int32]pb.MeetingRoom)

  rows1, err := db.Query("SELECT id, buildingId, floor, name from MEETING_ROOM")
  if err != nil {
    log.Fatal(err)
  }

  defer rows1.Close()

  for rows1.Next() {
    var mr pb.MeetingRoom
    err := rows1.Scan(&mr.Id, &mr.BuildingId, &mr.Floor, &mr.Name)
    if err != nil {
      log.Fatal(err)
    }

    items[mr.Id] = mr
  }

  err = rows1.Err()
  if err != nil {
    log.Fatal(err)
  }

  rows, err := db.Query("SELECT meetingRoomId from MEETING_ROOM_RESERVATION where (startTime>=? and startTime<?) or (endTime>? and endTime<=?)",
    startTime, endTime, startTime, endTime)
  if err != nil {
    log.Fatal(err)
  }

  defer rows.Close()

  for rows.Next() {
    var mrId int32
    err := rows.Scan(&mrId)
    if err != nil {
      log.Fatal(err)
    }

    delete (items, mrId)
  }

  err = rows.Err()
  if err != nil {
    log.Fatal(err)
  }

  values := []pb.MeetingRoom{}
  for _, value := range items {
    values = append(values, value)
  }

  return values
}
