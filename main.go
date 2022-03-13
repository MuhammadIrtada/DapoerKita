package main

import (
	"dapoer-kita/toko"
	"dapoer-kita/user"
	"fmt"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var r *gin.Engine

func InitDB() error {

	//Pakai Env Database
	_db, err := gorm.Open(mysql.Open("admin:HnVXVx8rF4G3YjS3nKuQrKVS7apg4Vzt@tcp(13.212.140.154:3306)/intern_bcc_10?parseTime=true"), &gorm.Config{})
	if err != nil {
		return err
	}
	db = _db
	if err = db.AutoMigrate(&user.User{}, &toko.Toko{}, &toko.Menu{}, &toko.Komentar{}, &toko.Video{},
		&toko.Category{}, &toko.RatingInfo{}, &toko.Artikel{}); err != nil {

		return err
	}
	return nil
}

func InitGin() {
	r = gin.Default()
	// config := cors.DefaultConfig()
	// config.AllowAllOrigins = true
	// congig.AllowHeaders = true

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}))
}

func StartServer() error {
	return r.Run()
}

func main() {
	if err := InitDB(); err != nil {
		fmt.Println("Database error on init!")
		fmt.Println(err.Error())
		return
	}

	InitGin()
	user.InitController(r, db)
	toko.InitController(r, db)

	if err := StartServer(); err != nil {
		fmt.Println("Server error!")
		fmt.Println(err.Error())
		return
	}
}
