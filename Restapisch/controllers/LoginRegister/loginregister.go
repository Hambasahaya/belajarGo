package LoginRegister

import (
	"net/http"
	"time"

	"restapisch/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var user models.User

	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BindJSON(&loginData); err != nil {
		c.JSON(500, gin.H{"message": "Data tidak valid"})
		return
	} else {
		if err := models.DB.Where("email = ?", loginData.Email).First(&user).Error; err != nil {
			c.JSON(404, gin.H{"message": "Email atau password salah"})
			return
		} else {
			if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)); err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"message": "password salah"})
				return
			} else {
				c.JSON(200, gin.H{"message": "Login berhasil", "userId": user.User_id, "email": user.Email, "level": user.Level})
				return
			}
		}

	}
}

func Register(c *gin.Context) {
	validate := validator.New()
	var user models.User

	var requestData struct {
		Email        string `json:"email" validate:"required,email"`
		Password     string `json:"password" validate:"required,min=6,max=16"`
		Nama_lengkap string `json:"nama"`
		Asal_sklh    string `json:"asl_sklh"`
		Tpl_lahir    string `json:"tpl_lahir"`
		Tgl_lahir    string `json:"tgl_lahir"`
		JK           string `json:"jk"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := validate.Struct(requestData); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		var errorMessages []string
		for _, err := range validationErrors {
			errorMessages = append(errorMessages, err.Field()+" "+err.Tag())
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": "Periksa kembali data Anda", "errors": errorMessages})
		return
	}

	if err := models.DB.Where("email = ?", requestData.Email).First(&user).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"message": "Email sudah terdaftar"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestData.Password), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat hash password"})
		return
	}

	newUser := models.User{
		Email:    requestData.Email,
		Password: string(hashedPassword),
	}
	if err := models.DB.Create(&newUser).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat pengguna baru"})
		return
	}

	id_user := newUser.User_id
	email := newUser.Email

	parsedDate, err := time.Parse("2006-01-02", requestData.Tgl_lahir)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Format tanggal lahir tidak valid"})
		return
	}
	newpendaftar := models.Pendaftar_sekolah{
		User_id:      int(id_user),
		Nama_lengkap: requestData.Nama_lengkap,
		Asal_sekolah: requestData.Asal_sklh,
		Tp_lahir:     requestData.Tpl_lahir,
		Tgl_lahir:    parsedDate,
		Foto:         "users.jpg",
	}

	if err := models.DB.Create(&newpendaftar).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat pengguna baru"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pengguna berhasil didaftarkan", "user_id": id_user, "email": email})
}

func Cekemail(c *gin.Context) {
	var user models.User
	var requestData struct {
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := models.DB.Where("email = ?", requestData.Email).First(&user).Error; err == nil {
		c.JSON(http.StatusOK, gin.H{"message": "Email sudah terdaftar"})
	}

}
func UpdatePassword(c *gin.Context) {
	var requestData struct {
		Email       string `json:"email"`
		NewPassword string `json:"password"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	var user models.User
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestData.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat hash password"})
		return
	}

	if err := models.DB.Model(&user).Where("email=?", requestData.Email).Update("password", string(hashedPassword)).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengupdate password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password berhasil diupdate"})
}
