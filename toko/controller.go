package toko

import (
	"dapoer-kita/authMiddle"
	"fmt"
	"net/http"
	"strconv"

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
				"nama":     body.Nama,
				"menu":     body.Menu,
				"funfact":  body.Funfact,
				"rating":   body.Rating,
				"komentar": body.Komentar,
				"video":    body.Video,
				"kota":     body.Kota,
				"category": body.Category,
			},
		})
	})

	// GET DISPLAY TOKO
	r.GET("/toko", func(c *gin.Context) {
		toko := []Toko{}
		if result := db.Preload("Menu").Preload("Funfact").Preload("Komentar").Preload("Video").Preload("Category").Find(&toko); result.Error != nil {
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
		if result := db.Preload("Menu").Preload("Funfact").Preload("Komentar").Preload("Video").Preload("Category").Where("id = ?", id).Take(&toko); result.Error != nil {
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
			db.Where("nama = ?", nama).Take(&cekNama)
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
			db.Where("kota = ?", kota).Take(&cekKota)
			if cekKota.Kota == "" {
				isKetemu = false
			} else {
				trx = trx.Where("nama LIKE ?", "%"+kota+"%")
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
			db.Where("nama = ?", menu).Take(&cekMenu)
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
			db.Where("rating = ?", rating).Take(&cekCategory)
			if cekCategory.Nama == "" {
				isKetemu = false
			} else {
				getCategory := Category{}
				queryResult := []Toko{}

				db.Preload("Toko").Where("nama = ?", category).Take(&getCategory)

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
	})

	r.Static("/material", "./material")
	r.POST("/toko/upload", func(c *gin.Context) {
		//Upload file
		file, err := c.FormFile("file")
		// text, err := c.FormFile("text")
		if err != nil {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}
		path := "material/" + file.Filename
		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
		c.JSON(http.StatusOK, fmt.Sprintf("File %s uploaded successfully", file.Filename))
	})
}

/*
r.Static("/assets", "./assets")
	fmt.Println("SUKSES")
	r.POST("/upload", func(c *gin.Context) {
		//Upload file
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}
		if err := c.SaveUploadedFile(file, file.Filename); err != nil {
			c.JSON(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
		c.JSON(http.StatusOK, fmt.Sprintf("File %s uploaded successfully", file.Filename))
	})
*/
