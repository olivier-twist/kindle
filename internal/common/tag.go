package common

type Tag struct {
	ID  int64  `db:"id" json:"id"`
	Tag string `db:"tag" json:"tag"`
}
