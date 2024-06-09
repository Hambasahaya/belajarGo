package userController

import (
	"net/http"
	"strconv"
	"time"

	"restapisch/models"

	"github.com/gin-gonic/gin"
)

func GetData(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID tidak valid"})
		return
	}
	type UserData struct {
		Id_pendaftaran int
		Id_user        int
		Email          string
		Nama_lengkap   string
		Asal_sekolah   string
		Tp_lahir       string
		Tgl_lahir      string
		JK             string
		Foto           string
	}

	var userData []UserData
	if err := models.DB.Table("users").
		Select("users.user_id, users.email, pendaftar_sekolahs.nama_lengkap, pendaftar_sekolahs.asal_sekolah, pendaftar_sekolahs.tp_lahir, pendaftar_sekolahs.tgl_lahir, pendaftar_sekolahs.jk, pendaftar_sekolahs.foto").
		Joins("INNER JOIN pendaftar_sekolahs ON users.user_id = pendaftar_sekolahs.user_id ").
		Where("users.user_id = ?", id).
		Scan(&userData).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengambil data pengguna dan pendaftar sekolah"})
		return
	}
	for i := range userData {
		if userData[i].JK == "LAKI-LAKI" {
			userData[i].JK = "1"
		} else if userData[i].JK == "PEREMPUAN" {
			userData[i].JK = "2"
		}
	}
	c.JSON(http.StatusOK, gin.H{"data": userData})
}
func Detail_pendaftars(c *gin.Context) {
	var D_user *models.Detail_pendaftar
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "ID tidak valid"})
		return
	}
	if err := models.DB.Where("user_id= ?", id).First(&D_user).Error; err != nil {
		// Menangani kasus ketika data tidak ditemukan
		c.JSON(http.StatusNotFound, gin.H{"message": "Detail pendaftar tidak ditemukan"})
		return
	}
	c.JSON(http.StatusOK, D_user)
}
func UpdateData(c *gin.Context) {
	var datauser *models.Detail_pendaftar
	var requestData struct {
		User_id         uint   `json:"User_id"`
		Email           string `json:"email"`
		Nama_lengkap    string `json:"nama_lengkap"`
		Asal_sekolah    string `json:"asal_sekolah"`
		Tp_lahir        string `json:"tp_lahir"`
		Tgl_lahir       string `json:"tgl_lahir"`
		JK              string `json:"jk"`
		Foto            string `json:"foto"`
		Nama_ibu        string `json:"nama_ibu"`
		Nama_ayah       string `json:"nama_ayah"`
		Alamat_orangtua string `json:"alamat_orangtua"`
		Pekerjaan_ibu   string `json:"pekerjaan_ibu"`
		Pekerjaan_ayah  string `json:"pekerjaan_ayah"`
		Alamat_siswa    string `json:"alamat_siswa"`
		Tlp_orangtua    string `json:"tlp_orangtua"`
	}

	if err := c.ShouldBindJSON(&requestData); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Update data pengguna (user)
	if err := models.DB.Model(&models.User{}).Where("user_id = ?", requestData.User_id).
		Updates(models.User{Email: requestData.Email}).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengupdate data pengguna"})
		return
	}
	tgl_lahir, err := time.Parse("2006-01-02", requestData.Tgl_lahir)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengonversi tanggal lahir"})
		return
	}
	// Update data pendaftaran siswa
	if err := models.DB.Model(&models.Pendaftar_sekolah{}).Where("user_id = ?", requestData.User_id).
		Updates(models.Pendaftar_sekolah{
			Nama_lengkap: requestData.Nama_lengkap,
			Asal_sekolah: requestData.Asal_sekolah,
			Tp_lahir:     requestData.Tp_lahir,
			Tgl_lahir:    tgl_lahir,
			JK:           requestData.JK,
			Foto:         requestData.Foto,
		}).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengupdate data pendaftaran siswa"})
		return
	}

	// Insert data detail siswa
	detailPendaftar := models.Detail_pendaftar{
		User_id:         int(requestData.User_id),
		Nama_ibu:        requestData.Nama_ibu,
		Nama_ayah:       requestData.Nama_ayah,
		Alamat_orangtua: requestData.Alamat_orangtua,
		Pekerjaan_ibu:   requestData.Pekerjaan_ibu,
		Pekerjaan_ayah:  requestData.Pekerjaan_ayah,
		Alamat_siswa:    requestData.Alamat_siswa,
		Tlp_orangtua:    requestData.Tlp_orangtua,
	}
	if err := models.DB.Where("user_id = ?", requestData.User_id).First(&datauser).Error; err != nil {
		if err := models.DB.Create(&detailPendaftar).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal menyimpan data detail siswa"})
			return
		}
	} else {
		if err := models.DB.Model(&models.Detail_pendaftar{}).Where("user_id = ?", requestData.User_id).
			Updates(models.Detail_pendaftar{
				User_id:         int(requestData.User_id),
				Nama_ibu:        requestData.Nama_ibu,
				Nama_ayah:       requestData.Nama_ayah,
				Alamat_orangtua: requestData.Alamat_orangtua,
				Pekerjaan_ibu:   requestData.Pekerjaan_ibu,
				Pekerjaan_ayah:  requestData.Pekerjaan_ayah,
				Alamat_siswa:    requestData.Alamat_siswa,
				Tlp_orangtua:    requestData.Tlp_orangtua,
			}).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal mengupdate data pendaftaran siswa"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data berhasil diupdate dan disimpan"})
}
