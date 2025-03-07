package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	swagFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "my-auth-api/docs" // Gantilah dengan nama modul Anda

	"golang.org/x/crypto/bcrypt"
)

// Koneksi ke PostgreSQL
var db *sql.DB

func init() {
	var err error
	connStr := "user=postgres password=sidoagung dbname=my_auth_db sslmode=disable"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
}

// Model User
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

// Struktur JWT
var jwtKey = []byte("secret_key")

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// @title MyAuth API
// @version 1.0
// @description API untuk otentikasi user dengan JWT di Golang menggunakan Gin Framework.
// @host localhost:8080
// @BasePath /

// @Summary Register User
// @Description Mendaftarkan user baru ke database
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body User true "User Data"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /register [post]
func register(c *gin.Context) {
	var user User
	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error hashing password"})
		return
	}

	_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, string(hashedPassword))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User berhasil didaftarkan"})
}

// @Summary Logout User
// @Description Logout user dengan menghapus token dari frontend
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]string
// @Router /logout [post]
func logout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Logout berhasil, silakan hapus token di frontend"})
}

// @Summary Login User
// @Description Login user dan mendapatkan token JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body User true "User Data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /login [post]
func login(c *gin.Context) {
	var user User
	var dbUser User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	row := db.QueryRow("SELECT id, username, password FROM users WHERE username=$1", user.Username)
	err := row.Scan(&dbUser.ID, &dbUser.Username, &dbUser.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
		return
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

// Middleware Auth JWT
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
			c.Abort()
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			c.Abort()
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}

// @Summary Get User Profile
// @Description Mengambil profil user berdasarkan token JWT
// @Tags User
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/profile [get]
func getProfile(c *gin.Context) {
	username, _ := c.Get("username")
	c.JSON(http.StatusOK, gin.H{"message": "Berhasil masuk!", "username": username})
}

func main() {
	r := gin.Default()

	// Menambahkan endpoint Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swagFiles.Handler))
	// Public Route
	r.POST("/register", register)
	r.POST("/login", login)
	r.POST("/logout", logout)

	// Protected Route
	protected := r.Group("/api")
	protected.Use(authMiddleware())
	protected.GET("/profile", getProfile)

	r.Run(":8080")
}
