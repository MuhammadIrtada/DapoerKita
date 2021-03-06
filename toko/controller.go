package toko

import (
	"dapoer-kita/authMiddle"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func InitController(r *gin.Engine, db *gorm.DB) {
	// REGISTER TOKO
	r.POST("/toko", func(c *gin.Context) {
		var body Toko
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid",
				"error":   err.Error(),
			})
			return
		}

		result := db.Create(&body)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when inserting into the database.",
				"error":   result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Toko berhasil dibuat.",
			"data": gin.H{
				"id":       body.ID,
				"nama":     body.Nama,
				"menu":     body.Menu,
				"funfact":  body.Funfact,
				"rating":   body.Rating,
				"komentar": body.Komentar,
				"kota":     body.Kota,
				"category": body.Category,
			},
		})
	})

	// GET DISPLAY TOKO
	r.GET("/toko", func(c *gin.Context) {
		toko := []Toko{}
		if result := db.Preload("Menu").Preload("Komentar").Preload("Category").Find(&toko); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Toko berhasil ditampilkan",
			"data":    toko,
		})
	})

	// Detail satu toko
	r.GET("/toko/:id", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}
		toko := Toko{}
		if result := db.Preload("Menu").Preload("Komentar").Preload("Category").Where("id = ?", id).Take(&toko); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Toko Berhasil Ditampilkan",
			"data":    toko,
		})
	})

	// Filter Toko
	r.GET("/toko/search", func(c *gin.Context) {
		nama, isNamaExists := c.GetQuery("nama")
		kota, isKotaExists := c.GetQuery("kota")
		menu, isMenuExists := c.GetQuery("menu")
		rating, isRatingExists := c.GetQuery("rating")
		category, isCategoryExists := c.GetQuery("category")

		queryResult := []Toko{}
		menuResult := []Menu{}
		trx := db

		isKetemu := true

		// Filter Nama
		if isNamaExists {
			cekNama := Toko{}
			db.Where("nama LIKE ?", "%"+nama+"%").Take(&cekNama)
			if cekNama.Nama == "" {
				isKetemu = false
			} else {
				trx = trx.Where("nama LIKE ?", "%"+nama+"%")
			}
		} else {
			trx = trx.Find(&queryResult)
		}

		// Filter Kota
		if isKotaExists {
			cekKota := Toko{}
			db.Where("kota LIKE ?", "%"+kota+"%").Take(&cekKota)
			if cekKota.Kota == "" {
				isKetemu = false
			} else {
				trx = trx.Where("kota LIKE ?", "%"+kota+"%")
			}
		} else {
			trx = trx.Find(&queryResult)
		}

		// Filter Rating
		if isRatingExists {
			cekRating := Toko{}
			db.Where("rating = ?", rating).Take(&cekRating)
			if cekRating.Rating == 0 {
				isKetemu = false
			} else {
				trx = trx.Where("rating = ?", rating)
			}
		} else {
			trx = trx.Find(&queryResult)
		}

		// Filter Menu
		if isMenuExists {
			cekMenu := Menu{}
			db.Where("nama LIKE ?", "%"+menu+"%").Take(&cekMenu)
			if cekMenu.Nama == "" {
				isKetemu = false
			} else {
				if result := db.Model(&Menu{}).Where("nama LIKE ?", "%"+menu+"%").Find(&menuResult); result.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"success": false,
						"message": "Error when querying the database.",
						"error":   result.Error.Error(),
					})
					return
				}

				var arrId = []uint{}
				for i := 0; i < len(menuResult); i++ {
					var add = append(arrId, menuResult[i].ID)
					arrId = add
				}
				trx = trx.Model(&Toko{}).Preload("Menu").Where("id IN (SELECT toko_id FROM menus WHERE id IN ?)", arrId).Find(&queryResult)
			}
		} else {
			trx = trx.Find(&queryResult)
		}

		// Filter Category
		if isCategoryExists {
			cekCategory := Category{}
			db.Where("nama LIKE ?", "%"+category+"%").Take(&cekCategory)
			if cekCategory.Nama == "" {
				isKetemu = false
			} else {
				getCategory := Category{}
				queryResult := []Toko{}

				db.Preload("Toko").Where("nama LIKE ?", "%"+category+"%").Take(&getCategory)

				var arrId = []uint{}
				for i := 0; i < len(getCategory.Toko); i++ {
					var add = append(arrId, getCategory.Toko[i].ID)
					arrId = add
				}

				trx = trx.Model(&Toko{}).Preload(clause.Associations).Where("id IN ?", arrId).Find(&queryResult)
			}
		} else {
			trx = trx.Find(&queryResult)
		}

		if result := trx.Model(&Toko{}).Preload(clause.Associations).Find(&queryResult); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}

		if isKetemu {
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "Toko berhasil ditampilkan",
				"data":    queryResult,
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Obyek yang anda cari tidak ditemukan",
				"data":    queryResult,
			})
		}
	})

	// Add Komentar
	r.POST("/toko/:id/komentar", authMiddle.AuthMiddleware(), func(c *gin.Context) {
		user_id, _ := c.Get("id")
		tokoId, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}

		toko_id, _ := strconv.ParseUint(tokoId, 10, 64)

		body := Komentar{}

		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid.",
				"error":   err.Error(),
			})
			return
		}
		komentar := Komentar{
			Toko_id: uint(toko_id),
			User_id: uint(user_id.(float64)),
			Teks:    body.Teks,
		}

		result := db.Create(&komentar)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Komentar gagal ditambahkan",
				"error":   result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Komentar berhasil ditambahkan",
			"data": gin.H{
				"toko_id": toko_id,
				"user_id": uint(user_id.(float64)),
				"teks":    komentar.Teks,
			},
		})
	})

	// Add Rating
	r.POST("/toko/:id/rating", authMiddle.AuthMiddleware(), func(c *gin.Context) {
		user_id, _ := c.Get("id")
		tokoId, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}

		// Input dari user
		body := RatingInfo{}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid.",
				"error":   err.Error(),
			})
			return
		}
		if body.Rating > 0 && body.Rating <= 5 {
			toko_id, _ := strconv.ParseUint(tokoId, 10, 64)

			// Cek antisipasi 2x rating
			cekRating := RatingInfo{}

			db.Where("toko_id = ?", uint(toko_id)).Where("user_id = ?", uint(user_id.(float64))).Take(&cekRating)

			if cekRating.Toko_id == 0 && cekRating.User_id == 0 {
				// Tambah rating pengguna
				ratingInfo := RatingInfo{
					Toko_id: uint(toko_id),
					User_id: uint(user_id.(float64)),
					Rating:  body.Rating,
				}
				result := db.Create(&ratingInfo)
				if result.Error != nil {
					c.JSON(http.StatusInternalServerError, gin.H{
						"success": false,
						"message": "Rating gagal ditambahkan",
						"error":   result.Error.Error(),
					})
					return
				}

				c.JSON(http.StatusCreated, gin.H{
					"success": true,
					"message": "Rating berhasil ditambahkan",
					"data": gin.H{
						"toko_id": toko_id,
						"user_id": uint(user_id.(float64)),
						"rating":  ratingInfo.Rating,
					},
				})

				// Menghitung rata-rata rating
				arrRating := []RatingInfo{}
				db.Where("toko_id = ?", uint(toko_id)).Find(&arrRating)
				sum := 0
				for _, rating := range arrRating {
					sum += int(rating.Rating)
				}
				avg := sum / len(arrRating)

				toko := Toko{
					ID:     uint(toko_id),
					Rating: uint(avg),
				}
				db.Updates(toko)

			} else {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Anda sudah menambahkan rating",
				})
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Masukkan rating antara 1 - 5",
			})
			return
		}
	})

	// UPLOAD ARTIKEL
	r.POST("/artikel", func(c *gin.Context) {
		body := Artikel{}
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid",
				"error":   err.Error(),
			})
			return
		}

		artikel := Artikel{
			Title:      body.Title,
			Author:     body.Author,
			Teks:       body.Teks,
			Created_At: time.Time{},
		}

		result := db.Create(&artikel)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when inserting into the database.",
				"error":   result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Artikel berhasil dibuat.",
			"data":    artikel,
		})
	})

	// GET ARTIKEL
	r.GET("/artikel", func(c *gin.Context) {
		artikel := []Artikel{}
		if result := db.Find(&artikel); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Artikel berhasil ditampilkan",
			"data":    artikel,
		})
	})

	// link environtment
	link := "https://f191-125-166-13-9.ngrok.io/"
	r.Static("/material", "./material")

	// UPLOAD GAMBAR ARTIKEL
	r.POST("/artikel/:id/upload", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}
		path := "material/gambar_artikel/" + id + "_" + RandomString(10) + "_" + file.Filename
		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
		parsedId, _ := strconv.ParseUint(id, 10, 64)

		pathDB := link + path

		artikel := Artikel{
			Gambar: pathDB,
		}

		resultUpdate := db.Model(&Artikel{}).Where("id = ?", parsedId).Updates(&artikel)
		if resultUpdate.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when updating the database.",
				"error":   resultUpdate.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, fmt.Sprintf("File %s uploaded successfully", file.Filename))
	})

	// UPLOAD GAMBAR MENU
	r.POST("/menu/:id/upload", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}
		path := "material/gambar_menu/" + id + "_" + RandomString(10) + "_" + file.Filename
		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
		parsedId, _ := strconv.ParseUint(id, 10, 64)

		pathDB := link + path

		menu := Menu{
			ID:     uint(parsedId),
			Gambar: pathDB,
		}

		resultUpdate := db.Model(&Menu{}).Where("id = ?", id).Updates(&menu)
		if resultUpdate.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when updating the database.",
				"error":   resultUpdate.Error.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, fmt.Sprintf("File %s uploaded successfully", file.Filename))
	})

	// UPLOAD VIDEO
	// link youtube
	r.POST("video/upload/link", func(c *gin.Context) {
		tokoId, isTokoIdExists := c.GetQuery("toko_id")
		linkY, islinkYExists := c.GetQuery("link")
		title, isTitleExists := c.GetQuery("title")
		author, isAuthorExists := c.GetQuery("author")

		if !isTokoIdExists && !isTitleExists && !islinkYExists && !isAuthorExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Query is not supplied.",
			})
			return
		}
		banner, err := c.FormFile("banner")
		if err != nil {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}
		pathBanner := "material/video/img_banner/" + tokoId + "_" + RandomString(10) + "_" + banner.Filename
		if err := c.SaveUploadedFile(banner, pathBanner); err != nil {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		parsedId, _ := strconv.ParseUint(tokoId, 10, 64)
		pathBannerDB := link + pathBanner
		resultVideo := Video{
			Toko_id: uint(parsedId),
			Video:   linkY,
			Banner:  pathBannerDB,
			Title:   title,
			Author:  author,
		}
		result := db.Create(&resultVideo)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when inserting into the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Gambar dan video berhasil diupload",
			"data":    resultVideo,
		})
	})
	// local
	r.POST("video/upload", func(c *gin.Context) {
		tokoId, isTokoIdExists := c.GetQuery("toko_id")
		title, isTitleExists := c.GetQuery("title")
		author, isMenuExists := c.GetQuery("author")

		if !isTokoIdExists && !isTitleExists && !isMenuExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Query is not supplied.",
			})
			return
		}

		video, err := c.FormFile("video")
		if err != nil {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}
		banner, err := c.FormFile("banner")
		if err != nil {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}
		pathVideo := "material/video/" + tokoId + "_" + RandomString(10) + "_" + video.Filename
		if err := c.SaveUploadedFile(video, pathVideo); err != nil {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
		pathBanner := "material/video/img_banner/" + tokoId + "_" + RandomString(10) + "_" + banner.Filename
		if err := c.SaveUploadedFile(banner, pathBanner); err != nil {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}

		parsedId, _ := strconv.ParseUint(tokoId, 10, 64)

		pathVideoDB := link + pathVideo
		pathBannerDB := link + pathBanner
		resultVideo := Video{
			Toko_id: uint(parsedId),
			Video:   pathVideoDB,
			Banner:  pathBannerDB,
			Title:   title,
			Author:  author,
		}

		result := db.Create(&resultVideo)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when inserting into the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "Gambar dan video berhasil diupload",
			"data":    resultVideo,
		})

		c.JSON(http.StatusOK, fmt.Sprintf("File %s and %s uploaded successfully", video.Filename, banner.Filename))

	})

}

// RANDOM STRING
func RandomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
