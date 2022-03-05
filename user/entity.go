package user

import "time"

type User struct {
	ID           uint      `gorm:"primarykey"`
	Nama         string    `json:"nama"`
	Email        string    `json:"email"`
	Password     string    `json:"password"`
	Cek_Password string    `json:"cek_password"`
	No_Telp      string    `json:"no_telp"`
	Created_At   time.Time `json:"created_at"`
}

type postRegisterBody struct {
	Nama       string    `json:"nama"`
	Email      string    `json:"email"`
	Password   string    `json:"password"`
	No_Telp    string    `json:"no_telp"`
	Created_At time.Time `json:"created_at"`
}

type postLoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type patchUserBody struct {
	Nama     string `json:"nama"`
	Email    string `json:"email"`
	Password string `json:"password"`
	No_Telp  string `json:"no_telp"`
}
