package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var r *gin.Engine

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
	Foto       string    `json:"foto"`
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

func InitDB() error {
	_db, err := gorm.Open(mysql.Open("root:@tcp(127.0.0.1:3306)/dapoer_kita?parseTime=true"), &gorm.Config{})
	if err != nil {
		return err
	}
	db = _db
	if err = db.AutoMigrate(&User{}); err != nil {
		return err
	}
	return nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hash), err
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header.Get("Authorization")
		header = header[len("Bearer "):]
		token, err := jwt.Parse(header, func(t *jwt.Token) (interface{}, error) {
			return []byte("passwordBuatSigning"), nil
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "JWT validation error.",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Set("id", claims["id"])
			c.Next()
			return
		} else {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "JWT invalid.",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}
	}
}

func InitRouter() {

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
				"id": user.ID,
			},
		})
	})

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
					"name":  user.Nama,
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

	r.GET("/user/:id", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}
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
			"message": "Query successful.",
			"data":    user,
		})
	})

	r.GET("/user/token", AuthMiddleware(), func(c *gin.Context) {
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
			"data":    user,
		})
	})

	r.GET("/user", func(c *gin.Context) {
		var allUsersFromDB []User

		if result := db.Find(&allUsersFromDB); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Error when querying the database.",
				"error":   result.Error.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Query successful.",
			"data":    allUsersFromDB,
		})
	})

	r.PATCH("/user/update", AuthMiddleware(), func(c *gin.Context) {
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

	r.DELETE("/user/delete/token", AuthMiddleware(), func(c *gin.Context) {
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

	r.DELETE("/user/delete/:id", func(c *gin.Context) {
		id, isIdExists := c.Params.Get("id")
		if !isIdExists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is not supplied.",
			})
			return
		}
		parsedId, err := strconv.ParseUint(id, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "ID is invalid.",
			})
			return
		}
		user := User{
			ID: uint(parsedId),
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

func InitGin() {
	r = gin.Default()
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
	InitRouter()

	if err := StartServer(); err != nil {
		fmt.Println("Server error!")
		fmt.Println(err.Error())
		return
	}
}
