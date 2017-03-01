package main

type LastIdHandler struct {
  Dao *Dao
}

func NewLastIdHandler(dao *Dao) *LastIdHandler{
  lih := &LastIdHandler{dao}
  return lih
}

func (lih *LastIdHandler) GetLastId(coordinate *Coordinate) int {
  up, down, left, right := coordinate.GetSquare()
  lastId := lih.Dao.GetLastId(up, down, left, right)

  return lastId
}
