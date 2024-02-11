package httpgin

import (
	"ads-server/internal/users"
	"github.com/gin-gonic/gin"
	"time"

	"ads-server/internal/ads"
)

type createAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

type adResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	AuthorID  int64     `json:"author_id"`
	Published bool      `json:"published"`
	CDate     time.Time `json:"create"`
	UDate     time.Time `json:"update"`
}

type userResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type userRequest struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type changeAdStatusRequest struct {
	Published bool  `json:"published"`
	UserID    int64 `json:"user_id"`
}

type updateAdRequest struct {
	Title  string `json:"title"`
	Text   string `json:"text"`
	UserID int64  `json:"user_id"`
}

func AdSuccessResponse(ad *ads.Ad) *gin.H {
	return &gin.H{
		"data": adResponse{
			ID:        ad.ID,
			Title:     ad.Title,
			Text:      ad.Text,
			AuthorID:  ad.AuthorID,
			Published: ad.Published,
		},
		"error": nil,
	}
}

func AdsSuccessResponse(allAds []*ads.Ad) *gin.H {
	var res []adResponse
	for _, val := range allAds {
		res = append(res,
			adResponse{
				ID:        val.ID,
				Title:     val.Title,
				Text:      val.Text,
				AuthorID:  val.AuthorID,
				Published: val.Published,
				CDate:     val.CDate,
				UDate:     val.UDate,
			})
	}
	return &gin.H{
		"data":  res,
		"error": nil,
	}
}

func UserSuccessResponse(user *users.User) *gin.H {
	return &gin.H{
		"data": userResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
		"error": nil,
	}
}

func AdErrorResponse(err error) *gin.H {
	return &gin.H{
		"data":  nil,
		"error": err.Error(),
	}
}
