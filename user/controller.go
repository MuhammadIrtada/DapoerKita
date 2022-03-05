package user

import (
	"dapoer-kita/authMiddle"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hash), err
}

func InitController(r *gin.Engine, db *gorm.DB) {
	// POST REGISTER
	r.POST("/user/register", func(c *gin.Context) {
		var body postRegisterBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid",
				"error":   err.Error(),
			})
			return
		}

		hash, _ := HashPassword(body.Password)

		user := User{
			Nama:         body.Nama,
			Email:        body.Email,
			Password:     hash,
			Cek_Password: body.Password,
			No_Telp:      body.No_Telp,
			Created_At:   time.Now(),
		}

		var cek = []byte(body.Password)
		var angka = false

		for i := 0; i < len(body.Password); i++ {
			if cek[i] >= 48 && cek[i] <= 57 {
				angka = true
			}
		}

		if (len(body.Password) >= 8) && (cek[0] >= 65 && cek[0] <= 90) && angka {
			result := db.Create(&user)
			if result.Error != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Error when inserting into the database.",
					"error":   result.Error.Error(),
				})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Password harus lebih dari 8 karakter, Huruf pertama harus besar, Password harus terdapat angka",
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"success": true,
			"message": "User berhasil dibuat.",
			"data": gin.H{
				"nama":       user.Nama,
				"email":      user.Email,
				"no_telp":    user.No_Telp,
				"created_at": user.Created_At,
			},
		})
	})

	// POST LOGIN
	r.POST("/user/login", func(c *gin.Context) {
		var body postLoginBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid",
				"error":   err.Error(),
			})
			return
		}
		user := User{}
		if result := db.Where("email = ?", body.Email).Take(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying into the database.",
				"error":   result.Error.Error(),
			})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err == nil {
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"id":  user.ID,
				"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
			})
			tokenString, err := token.SignedString([]byte("passwordBuatSigning"))
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Error when generating the token.",
					"error":   err.Error(),
				})
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"message": "Login Berhasil.",
				"data": gin.H{
					"id":    user.ID,
					"nama":  user.Nama,
					"token": tokenString,
				},
			})
			return
		} else {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "Password salah.",
			})
			return
		}
	})

	// GET SHOWPROFILE
	r.GET("/user/showprofile", authMiddle.AuthMiddleware(), func(c *gin.Context) {
		id, _ := c.Get("id")
		user := User{}
		if result := db.Where("id = ?", id).Take(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Query successful",
			"data": gin.H{
				"ID":         user.ID,
				"nama":       user.Nama,
				"email":      user.Email,
				"password":   user.Cek_Password,
				"no_telp":    user.No_Telp,
				"created_at": user.Created_At,
			},
		})
	})

	r.PATCH("/user/update", authMiddle.AuthMiddleware(), func(c *gin.Context) {
		id, _ := c.Get("id")
		var body patchUserBody
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Body is invalid.",
				"error":   err.Error(),
			})
			return
		}

		hash, _ := HashPassword(body.Password)

		user := User{
			ID:           uint(id.(float64)),
			Nama:         body.Nama,
			Email:        body.Email,
			Password:     hash,
			Cek_Password: body.Password,
			No_Telp:      body.No_Telp,
		}
		result := db.Model(&user).Updates(user)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when updating the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		if result = db.Where("id = ?", id).Take(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		if result.RowsAffected < 1 {
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": "User not found.",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Update Success.",
			"data":    user,
		})
	})

	r.DELETE("/user/delete", authMiddle.AuthMiddleware(), func(c *gin.Context) {
		id, _ := c.Get("id")

		user := User{
			ID: uint(id.(float64)),
		}
		if result := db.Delete(&user); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when deleting from the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Delete success.",
		})
	})
}
