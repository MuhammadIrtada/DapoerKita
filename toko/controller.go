package toko

import (
	"fmt"
	"net/http"

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

		queryResult := []Toko{}
		menuResult := []Menu{}
		trx := db

		// Tanpa filter
		if !isNamaExists && !isKotaExists && !isRatingExists {
			if result := db.Find(&queryResult); result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Error when querying the database.",
					"error":   result.Error.Error(),
				})
				return
			}
		}

		// Filter Nama
		if isNamaExists {
			trx = trx.Where("nama LIKE ?", "%"+nama+"%")
		} else {
			trx = trx.Find(&queryResult)
		}

		// Filter Kota
		if isKotaExists {
			trx = trx.Where("kota LIKE ?", "%"+kota+"%")
		} else {
			trx = trx.Find(&queryResult)
		}

		// Filter Rating
		if isRatingExists {
			trx = trx.Where("rating = ?", rating)
		} else {
			trx = trx.Find(&queryResult)
		}

		// Filter Menu
		if isMenuExists {
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

			fmt.Println(arrId)

			trx = trx.Model(&Toko{}).Preload("Menu").Where("id IN (SELECT toko_id FROM menus WHERE id IN ?)", arrId).Find(&queryResult)

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

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Toko berhasil ditampilkan",
			"data":    queryResult,
		})
	})

	// Category
	r.GET("/toko/category", func(c *gin.Context) {
		menuMasukan, _ := c.GetQuery("menu")
		// num, _ := strconv.ParseInt(menuMasukan, 10, 64)

		toko := Category{
			Nama: menuMasukan,
		}

		// categoryy := Category{
		// 	ID: 3,
		// }

		// db.Preload("Category").Take(&toko)
		db.Preload("Toko").Take(&toko)

		fmt.Println(toko)

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Toko berhasil ditampilkan",
			"data":    toko,
		})

	})
}
