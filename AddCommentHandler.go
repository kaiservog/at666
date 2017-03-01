package main

type AddCommentHandler struct {
  Dao *Dao
}

func NewAddCommentHandler(dao *Dao) *AddCommentHandler{
  ach := &AddCommentHandler{dao}
  return ach
}

func (ach *AddCommentHandler) AddComment(coordinate *Coordinate, nick, text string) error {
  return ach.Dao.AddComment(nick, text, coordinate.Lat, coordinate.Lon)
}
