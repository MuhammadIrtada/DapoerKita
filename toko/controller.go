package toko

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitController(r *gin.Engine, db *gorm.DB) {
	r.POST("/toko/register", func(c *gin.Context) {
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
}
