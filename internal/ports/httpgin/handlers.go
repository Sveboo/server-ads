package httpgin

import (
	"ads-server/internal/errs"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	"ads-server/internal/app"
)

// createUser handles route to create new user
func createUser(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody userRequest
		if err := c.Bind(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		user, err := a.CreateUser(c, reqBody.Name, reqBody.Email)
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		c.JSON(http.StatusOK, UserSuccessResponse(user))
	}
}

// getAdsByName handles route to find all ads by title given
func getAdsByName(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		title := c.Param("title")
		ads := a.GetAdByName(c, title)
		if len(ads) == 0 {
			c.Status(http.StatusNotFound)
			c.JSON(http.StatusNotFound, AdErrorResponse(fmt.Errorf("no ads with such name")))
			return
		}
		c.JSON(http.StatusOK, AdsSuccessResponse(ads))
	}
}

// getAdByID handles route to get ad by ID given
func getAdByID(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {

		id, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}
		ad, err := a.GetAdByID(c, int64(id))

		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, err)
			return
		}

		if errors.Is(err, errs.AdNotFoundError) {
			c.Status(http.StatusNotFound)
			c.JSON(http.StatusNotFound, err)
			return
		}
		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для создания объявления (ad)
func createAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody createAdRequest
		err := c.Bind(&reqBody)
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}
		if _, err = a.FindUser(c, reqBody.UserID); errors.Is(err, errs.UserNotFoundError) {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}
		ad, err := a.CreateAd(c, reqBody.UserID, reqBody.Title, reqBody.Text)
		if errors.Is(err, errs.ValidationError) {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		if errors.Is(err, errs.AccessError) {
			c.Status(http.StatusForbidden)
			c.JSON(http.StatusForbidden, AdErrorResponse(err))
			return
		}
		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// changeAdStatus handles route to change ad status
func changeAdStatus(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody changeAdStatusRequest
		if err := c.Bind(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		adID, err := strconv.Atoi(c.Param("ad_id"))
		if err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}
		//if adID == 0 {
		//	c.Status(http.StatusBadRequest)
		//	c.JSON(http.StatusBadRequest, AdErrorResponse(fmt.Errorf("empty field ad_id")))
		//}

		if _, err := a.FindUser(c, reqBody.UserID); errors.Is(err, errs.UserNotFoundError) {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		ad, err := a.PublishAd(c, int64(adID), reqBody.UserID, reqBody.Published)

		if errors.Is(err, errs.AccessError) {
			c.Status(http.StatusForbidden)
			c.JSON(http.StatusForbidden, AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}
}

// Метод для обновления текста(Text) или заголовка(Title) объявления
func updateAd(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody updateAdRequest
		if err := c.Bind(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		adID := c.GetInt64("ad_id")
		if _, err := a.FindUser(c, reqBody.UserID); errors.Is(err, errs.UserNotFoundError) {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		ad, err := a.UpdateAd(c, adID, reqBody.UserID, reqBody.Title, reqBody.Text)

		if errors.Is(err, errs.ValidationError) {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		}

		if errors.Is(err, errs.AccessError) {
			c.Status(http.StatusForbidden)
			c.JSON(http.StatusForbidden, AdErrorResponse(err))
			return
		}

		c.JSON(http.StatusOK, AdSuccessResponse(ad))
	}

}

// filterAds handles route to return filtered ads
func filterAds(a app.App) gin.HandlerFunc {
	return func(c *gin.Context) {
		params := c.Request.URL.Query()
		if allAds, err := a.Filter(c, params); err != nil {
			c.Status(http.StatusBadRequest)
			c.JSON(http.StatusBadRequest, AdErrorResponse(err))
			return
		} else {
			c.JSON(http.StatusOK, AdsSuccessResponse(allAds))
		}
	}
}
