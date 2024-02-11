package repo

import (
	"ads-server/internal/ads"
	"ads-server/internal/app"
	"ads-server/internal/errs"
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type AdRepo struct {
	storage map[int64]*ads.Ad
	mx      *sync.Mutex
	lastID  int64
}

// Create is a function to create a new ad
func (ar *AdRepo) Create(_ context.Context, ad *ads.Ad) (id int64, err error) {
	ar.mx.Lock()
	defer ar.mx.Unlock()
	if _, ok := ar.storage[ar.lastID]; ok {
		return -1, errs.AdNotFoundError
	}
	ad.ID = ar.lastID
	ad.CDate = time.Now().UTC()
	ad.UDate = ad.CDate
	ar.storage[ar.lastID] = ad
	ar.lastID++

	return ar.lastID - 1, nil
}

// Update is a function to update an existing ad
func (ar *AdRepo) Update(_ context.Context, id int64, aID int64, title string, text string) (*ads.Ad, error) {
	ar.mx.Lock()
	defer ar.mx.Unlock()
	if _, ok := ar.storage[id]; !ok {
		return nil, errs.AdNotFoundError
	}
	if ar.storage[id].AuthorID != aID {
		return nil, errs.AccessError
	}
	ar.storage[id].Text = text
	ar.storage[id].Title = title
	ar.storage[id].UDate = time.Now().UTC()
	return ar.storage[id], nil
}

// Delete deletes ad from storage
func (ar *AdRepo) Delete(_ context.Context, id, uID int64) error {
	ar.mx.Lock()
	defer ar.mx.Unlock()
	ad, ok := ar.storage[id]
	if !ok {
		return errs.AdNotFoundError
	}

	if ad.AuthorID != uID {
		return errs.AccessError
	}

	delete(ar.storage, id)
	return nil
}

// Publish is a function to change ad status
func (ar *AdRepo) Publish(_ context.Context, adID, aID int64, action bool) (*ads.Ad, error) {
	ar.mx.Lock()
	defer ar.mx.Unlock()
	_, ok := ar.storage[adID]
	if !ok || ar.storage[adID].AuthorID != aID {
		return nil, fmt.Errorf("access forbidden")
	}
	ar.storage[adID].Published = action
	ar.storage[adID].UDate = time.Now().UTC()
	return ar.storage[adID], nil

	//return nil, errors.NewAd(fmt.Sprintf("no such ad %d", adID))
}

// GetByID is a function to find ad in storage using ID
func (ar *AdRepo) GetByID(_ context.Context, id int64) (*ads.Ad, error) {
	ar.mx.Lock()
	defer ar.mx.Unlock()
	for adID, ad := range ar.storage {
		if adID == id {
			return ad, nil
		}
	}
	return nil, errs.AdNotFoundError
}

// GetByName is a function to find ad in storage using name
func (ar *AdRepo) GetByName(_ context.Context, title string) []*ads.Ad {
	ar.mx.Lock()
	defer ar.mx.Unlock()
	var resAds []*ads.Ad
	for _, val := range ar.storage {
		if (title == "" || strings.Contains(val.Title, title)) && val.Published {
			resAds = append(resAds, val)
		}
	}
	return resAds
}

// Date represents date
type Date struct {
	Y int
	M int
	D int
}

// parseDate parses ISO format (yyyy-mm-dd) to Date
func parseDate(dateISO string) (Date, error) {
	splitted := strings.Split(dateISO, "-")
	Y, err := strconv.Atoi(splitted[0])
	if err != nil {
		return Date{}, err
	}
	M, err := strconv.Atoi(splitted[1])
	if err != nil {
		return Date{}, err
	}
	D, err := strconv.Atoi(splitted[2])
	if err != nil {
		return Date{}, err
	}
	return Date{
		Y: Y,
		M: M,
		D: D,
	}, nil
}

// Filter returns ads satisfying filters
func (ar *AdRepo) Filter(_ context.Context, params url.Values) ([]*ads.Ad, error) {
	ar.mx.Lock()
	defer ar.mx.Unlock()

	author, mustAuthor := params["author"]
	authorID := -1
	var err error
	if mustAuthor {
		authorID, err = strconv.Atoi(author[0])
		if err != nil {
			return nil, err
		}
	}

	date, mustDate := params["date"]
	var parsedDate Date
	if mustDate {
		parsedDate, err = parseDate(date[0])
		if err != nil {
			return nil, err
		}
	}

	title, mustTitle := params["title"]

	_, mustPublished := params["published"]

	var allAds []*ads.Ad

	for _, ad := range ar.storage {
		if mustPublished && !ad.Published {
			continue
		}
		if mustAuthor && ad.AuthorID != int64(authorID) {
			continue
		}
		year, month, day := ad.CDate.Date()
		if mustDate && (parsedDate.Y != year || parsedDate.M != int(month) || parsedDate.D != day) {
			continue
		}
		if mustTitle && ad.Title != title[0] {
			continue
		}
		allAds = append(allAds, ad)
	}

	return allAds, nil
}

// NewAd is a constructor
func NewAd() app.AdRepository {
	return &AdRepo{
		mx:      &sync.Mutex{},
		storage: make(map[int64]*ads.Ad, 1),
		lastID:  0,
	}
}
