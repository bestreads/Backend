package database

// https://limbd.org/isbn-10-and-isbn-13-check-digit-and-missing-digit-calculation/

type Book struct {
	ID          uint   `gorm:"primaryKey"`
	ISBN        string // das zu validieren ist nicht so ez, ich mache es noch nicht
	Title       string
	Author      string
	CoverURL    string
	RatingAvg   float32
	Description string
}
