package toko

import (
	"dapoer-kita/user"
	"time"
)

// POST
type Menu struct {
	ID      uint    `gorm:"primarykey" json:"id"`
	Toko    Toko    `json:"toko"`
	Toko_id uint    `json:"toko_id"`
	Nama    string  `json:"nama"`
	Harga   float64 `json:"harga"`
	Gambar  string  `json:"gambar"`
	Desk    string  `json:"desk"`
}

type Komentar struct {
	ID      uint      `gorm:"primarykey" json:"id"`
	Toko    Toko      `json:"toko"`
	Toko_id uint      `json:"toko_id"`
	User    user.User `json:"user"`
	User_id uint      `json:"user_id"`
	Teks    string    `json:"teks"`
}

type Video struct {
	ID      uint   `gorm:"primarykey" json:"id"`
	Toko    Toko   `json:"toko"`
	Toko_id uint   `json:"toko_id"`
	Video   string `json:"video"`
	Title   string `json:"title"`
	Author  string `json:"author"`
	Banner  string `json:"banner"`
}

type Category struct {
	ID   uint   `gorm:"primarykey" json:"id"`
	Nama string `json:"nama"`
	Toko []Toko `gorm:"many2many:toko_category;"`
}

type RatingInfo struct {
	ID      uint      `gorm:"primarykey" json:"id"`
	Toko    Toko      `json:"toko"`
	Toko_id uint      `json:"toko_id"`
	User    user.User `json:"user"`
	User_id uint      `json:"user_id"`
	Rating  uint      `json:"rating"`
}

type Artikel struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	Title      string    `json:"title"`
	Author     string    `json:"author"`
	Gambar     string    `json:"gambar"`
	Teks       string    `json:"teks"`
	Created_At time.Time `json:"created_at"`
}

type Toko struct {
	ID       uint       `gorm:"primarykey" json:"id"`
	Nama     string     `json:"nama"`
	Menu     []Menu     `json:"menu"`
	Funfact  string     `json:"funfact"`
	Rating   uint       `json:"rating"`
	Komentar []Komentar `json:"komentar"`
	Kota     string     `json:"kota"`
	Category []Category `gorm:"many2many:toko_category;" json:"category"`
}
