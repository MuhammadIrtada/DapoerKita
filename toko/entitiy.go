package toko

// POST
type Menu struct {
	ID      uint    `gorm:"primarykey" json:"id"`
	Toko_id uint    `json:"toko_id"`
	Nama    string  `json:"nama"`
	Harga   float64 `json:"harga"`
	Gambar  string  `json:"gambar"`
}

type Komentar struct {
	ID      uint   `gorm:"primarykey" json:"id"`
	Toko_id uint   `json:"toko_id"`
	User_id uint   `json:"user_id"`
	Teks    string `json:"teks"`
}

type Video struct {
	ID      uint   `gorm:"primarykey" json:"id"`
	Toko_id uint   `json:"toko_id"`
	Menu_id uint   `json:"menu_id"`
	Video   string `json:"video"`
}

type Funfact struct {
	ID      uint   `gorm:"primarykey" json:"id"`
	Toko_id uint   `json:"toko_id"`
	Menu_id uint   `json:"menu_id"`
	Teks    string `json:"teks"`
}

type Category struct {
	ID   uint   `gorm:"primarykey" json:"id"`
	Nama string `json:"nama"`
	Toko []Toko `gorm:"many2many:toko_category;"`
}

type RatingInfo struct {
	ID      uint `gorm:"primarykey" json:"id"`
	Toko_id uint `json:"toko_id"`
	User_id uint `json:"user_id"`
	Rating  uint `json:"rating"`
}

type Toko struct {
	ID       uint       `gorm:"primarykey" json:"id"`
	Nama     string     `json:"nama"`
	Menu     []Menu     `json:"menu"`
	Funfact  []Funfact  `json:"funfact"`
	Rating   uint       `json:"rating"`
	Komentar []Komentar `json:"komentar"`
	Video    []Video    `json:"video"`
	Kota     string     `json:"kota"`
	Category []Category `gorm:"many2many:toko_category;" json:"category"`
}
