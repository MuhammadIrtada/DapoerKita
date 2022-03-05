package toko

type Menu struct {
	ID      uint `gorm:"primarykey"`
	Toko_id uint
	Nama    string
	Harga   float64
	Gambar  string
}

type Komentar struct {
	ID      uint `gorm:"primarykey"`
	Toko_id uint
	User_id uint
	Teks    string
}

type Video struct {
	ID      uint `gorm:"primarykey"`
	Toko_id uint
	Menu_id uint
	Video   string
}
