package httpgin

import (
	"ads-server/internal/app"
	"github.com/gin-gonic/gin"
)

func AppRouter(r gin.IRouter, a app.App) {
	r.GET("/ads/:ad_id/info", getAdByID(a))        // Метод для получения объявления по ID
	r.POST("/user", createUser(a))                 // Метод для создания пользователя (user)
	r.POST("/ads", createAd(a))                    // Метод для создания объявления (ad)
	r.PUT("/ads/:ad_id/status", changeAdStatus(a)) // Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
	r.PUT("/ads/:ad_id", updateAd(a))              // Метод для обновления текста(Text) или заголовка(Title) объявления
	r.GET("ads/find/:title", getAdsByName(a))      // Метод для получения списка объявлений по имени
	r.GET("ads/filter", filterAds(a))              // Метод для фильтрации объявлений по query-параметрам
}
