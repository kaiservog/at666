package main


type LastCommentsHandler struct {
  Dao *Dao
}

func NewLastCommentsHandler(dao *Dao) *LastCommentsHandler{
  lch := &LastCommentsHandler{dao}
  return lch
}

func (lch *LastCommentsHandler) GetLastComments(coordinate *Coordinate, quantity int) (*Comments, error) {
  up, down, left, right := coordinate.GetSquare()
  comments, err := lch.Dao.GetLastsComments(quantity, up, down, left, right)
  return comments, err
}
