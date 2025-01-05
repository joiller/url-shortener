package url_shortener

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

type ShortUrl struct {
	gorm.Model
	Short       string `json:"shortCode" gorm:"column:short;type:varchar(255);not null"`
	Long        string `json:"url" gorm:"column:long;type:varchar(255);not null"`
	AccessCount int    `json:"accessCount" gorm:"column:access_count;type:int;default:0"`
}

func init() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/shortener?charset=utf8&parseTime=True&loc=Local",
		Env.User, Env.Password, Env.Host, Env.Port)
	d, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // data source name
		DefaultStringSize:         256,   // default size for string fields
		DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db = d
	err = db.AutoMigrate(&ShortUrl{})
	if err != nil {
		panic("failed to migrate database")
	}
}

func LongToShort(req LongToShortRequest) (shortUrl ShortUrl, err error) {
	shortUrl.Long = req.LongUrl
	if err := db.Where("`long` = ?", req.LongUrl).First(&shortUrl).Error; err != nil {
		shortUrl.Short = GenerateShortUrl()
		if err := db.Create(&shortUrl).Error; err != nil {
			return shortUrl, err
		}
		return shortUrl, nil
	}
	return shortUrl, nil
}

func ShortToLong(short string) (shortUrl ShortUrl, err error) {
	if err := db.Where("`short` = ?", short).First(&shortUrl).Error; err != nil {
		return shortUrl, err
	}
	shortUrl.AccessCount++
	db.Save(&shortUrl)
	return shortUrl, nil
}

func GenerateShortUrl() string {
	return uuid.NewString()[:8]
}

func UpdateShortUrl(short string, request UpdateShortUrlRequest) (shortUrl ShortUrl, err error) {
	if err := db.Model(&ShortUrl{}).Where("`short` = ?", short).First(&shortUrl).Error; err != nil {
		return shortUrl, err
	}
	shortUrl.Long = request.Long
	db.Save(&shortUrl)
	return shortUrl, nil
}

func DeleteShortUrl(short string) (err error) {
	var shortUrl ShortUrl
	if err := db.Where("`short` = ?", short).First(&shortUrl).Error; err != nil {
		return err
	}
	db.Delete(&shortUrl)
	return nil
}

func GetShortUrlStatus(short string) (shortUrl ShortUrl, err error) {
	if err := db.Where("`short` = ?", short).First(&shortUrl).Error; err != nil {
		return shortUrl, err
	}
	return shortUrl, nil
}
