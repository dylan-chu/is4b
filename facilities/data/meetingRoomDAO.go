package data

import (
  "log"

  pb "is4b/facilities/proto"
  _ "github.com/go-sql-driver/mysql"
)

type MeetingRoomDAO struct {
    ConnStr string
}

// func New(connStr string) MeetingRoomDAO {
//   dao := MeetingRoomDAO{ connStr }
//   return dao
// }

func (dao MeetingRoomDAO) Create(item *pb.MeetingRoom) error {
  db := getDbConn(dao.ConnStr)
  defer db.Close()

  query := "INSERT into MEETING_ROOM (buildingId, floor, name, capacity, hasProjector, hasWhiteboard, hasConferenceLine) VALUES (?, ?, ?, ?, ?, ?, ?)"
  stmt, err := db.Prepare(query)
  defer stmt.Close()
  if err != nil {
    return err
  }

  result, err := stmt.Exec(item.BuildingId, item.Floor, item.Name, item.Capacity, item.HasProjector, item.HasWhiteboard, item.HasConferenceLine)
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

func (dao MeetingRoomDAO) GetAll() []pb.MeetingRoom {
  db := getDbConn(dao.ConnStr)
  defer db.Close()

  items := make([]pb.MeetingRoom, 0)

  rows, err := db.Query("SELECT id, buildingId, floor, name, capacity, hasProjector, hasWhiteboard, hasConferenceLine from MEETING_ROOM")
  if err != nil {
    log.Fatal(err)
  }

  defer rows.Close()

  for rows.Next() {
    var item pb.MeetingRoom
    err := rows.Scan(&item.Id, &item.BuildingId, &item.Floor, &item.Name,
      &item.Capacity, &item.HasProjector, &item.HasWhiteboard, &item.HasConferenceLine)
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
