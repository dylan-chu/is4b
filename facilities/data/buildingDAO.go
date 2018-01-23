package data

import (
  "log"

  pb "is4b/facilities/proto"
  _ "github.com/go-sql-driver/mysql"
)

type BuildingDAO struct {
    ConnStr string
}

// func New(connStr string) BuildingDAO {
//   dao := BuildingDAO{ connStr }
//   return dao
// }

func (dao BuildingDAO) Create(building *pb.Building) error {
  db := getDbConn(dao.ConnStr)
  defer db.Close()

  query := "INSERT into BUILDING (name) VALUES (?)"
  stmt, err := db.Prepare(query)
  defer stmt.Close()
  if err != nil {
    return err
  }

  result, err := stmt.Exec(building.Name)
  if err != nil {
    return err
  }

  rowId, err := result.LastInsertId()
  if err != nil {
      return err
  }

  building.Id = int32(rowId)
  return nil
}

func (dao BuildingDAO) GetAll() []pb.Building {
  db := getDbConn(dao.ConnStr)
  defer db.Close()

  buildings := make([]pb.Building, 0)

  rows, err := db.Query("SELECT id, name from BUILDING")
  if err != nil {
    log.Fatal(err)
  }

  defer rows.Close()

  for rows.Next() {
    var building pb.Building
    err := rows.Scan(&building.Id, &building.Name)
    if err != nil {
      log.Fatal(err)
    }

    buildings = append(buildings, building)
  }

  err = rows.Err()
  if err != nil {
    log.Fatal(err)
  }

  return buildings
}
