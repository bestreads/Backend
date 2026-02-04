package database

// https://limbd.org/isbn-10-and-isbn-13-check-digit-and-missing-digit-calculation/

type Book struct {
	ID          uint   `gorm:"primaryKey"`
	ISBN        string `gorm:"uniqueIndex"` // das zu validieren ist nicht so ez, ich mache es noch nicht
	Title       string
	Author      string
	CoverURL    string
	Rating      Rating `gorm:embedded`
	Description string
	ReleaseDate uint64 // am besten unix epoch oder so, erfordert keine speziellen datentypen
	Genre       string
}

type Rating struct {
	Avg   float32
	Count uint
}
