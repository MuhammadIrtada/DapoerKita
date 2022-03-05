package toko

import "dapoer-kita/user"

type Menu struct {
	ID      uint `gorm:"primarykey"`
	Toko    Toko
	Toko_id uint
	Nama    string
	Harga   float64
	Gambar  string
}

type Komentar struct {
	ID      uint `gorm:"primarykey"`
	Toko    Toko
	Toko_id uint
	User    user.User
	User_id uint
	Teks    string
}

type Video struct {
	ID      uint `gorm:"primarykey"`
	Toko    Toko
	Toko_id uint
	Menu    Menu
	Menu_id uint
	Video   string
}

type Funfact struct {
	ID      uint `gorm:"primarykey"`
	Toko    Toko
	Toko_id uint
	Menu    Menu
	Menu_id uint
	Teks    string
}

type Category struct {
	ID   uint `gorm:"primarykey"`
	Nama string
	Toko []Toko `gorm:"many2many:toko_category;"`
}

type Toko struct {
	ID       uint `gorm:"primarykey"`
	Nama     string
	Menu     []Menu
	Funfact  []Funfact
	Rating   uint
	Komentar []Komentar
	Video    []Video
	Kota     string
	Category []Category `gorm:"many2many:toko_category;"`
}

// type Orang struct {
// 	ID    uint `gorm:"primarykey"`
// 	Nama  string
// 	Mobil []Mobil
// }

// type Mobil struct {
// 	ID       uint `gorm:"primarykey"`
// 	Nama     string
// 	Orang    Orang
// 	Orang_id uint
// }
