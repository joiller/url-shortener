package url_shortener

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LongToShortRequest struct {
	LongUrl string `json:"url" binding:"required"`
}

type LongToShortResponse struct {
	ID        uint      `json:"id"`
	Short     string    `json:"shortCode" `
	Long      string    `json:"url"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func LongToShortHandler(c *gin.Context) {
	var req LongToShortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := LongToShort(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, LongToShortResponseFromModel(resp))
}

func LongToShortResponseFromModel(shortUrl ShortUrl) LongToShortResponse {
	response := LongToShortResponse{
		ID:        shortUrl.ID,
		Short:     shortUrl.Short,
		Long:      shortUrl.Long,
		CreatedAt: shortUrl.CreatedAt,
		UpdatedAt: shortUrl.UpdatedAt,
	}
	return response
}

type ShortToLongResponse struct {
	ID        uint      `json:"id"`
	Short     string    `json:"shortCode" `
	Long      string    `json:"url"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func ShortToLongHandler(c *gin.Context) {
	short := c.Param("short")

	resp, err := ShortToLong(short)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ShortToLongResponseFromModel(resp))
}

func ShortToLongResponseFromModel(shortUrl ShortUrl) ShortToLongResponse {
	response := ShortToLongResponse{
		ID:        shortUrl.ID,
		Short:     shortUrl.Short,
		Long:      shortUrl.Long,
		CreatedAt: shortUrl.CreatedAt,
		UpdatedAt: shortUrl.UpdatedAt,
	}
	return response
}

type UpdateShortUrlRequest struct {
	Long string `json:"url" binding:"required"`
}

func UpdateShortUrlHandler(c *gin.Context) {
	short := c.Param("short")
	req := UpdateShortUrlRequest{}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if short == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "short code is required"})
		return
	}

	resp, err := UpdateShortUrl(short, req)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ShortToLongResponseFromModel(resp))
}

func DeleteShortUrlHandler(c *gin.Context) {
	short := c.Param("short")
	if short == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "short code is required"})
		return
	}

	err := DeleteShortUrl(short)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

type ShortUrlStatusResponse struct {
	ID        uint      `json:"id"`
	Short     string    `json:"shortCode" `
	Long      string    `json:"url"`
	Access    int       `json:"accessCount"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func ShortUrlStatusResponseFromModel(shortUrl ShortUrl) ShortUrlStatusResponse {
	response := ShortUrlStatusResponse{
		ID:        shortUrl.ID,
		Short:     shortUrl.Short,
		Long:      shortUrl.Long,
		Access:    shortUrl.AccessCount,
		CreatedAt: shortUrl.CreatedAt,
		UpdatedAt: shortUrl.UpdatedAt,
	}
	return response
}

func GetShortUrlStatusHandler(c *gin.Context) {
	short := c.Param("short")
	if short == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "short code is required"})
	}

	resp, err := GetShortUrlStatus(short)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ShortUrlStatusResponseFromModel(resp))
}
