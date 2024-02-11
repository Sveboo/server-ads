package app

import (
	"ads-server/internal/errs"
	"context"
	"errors"
	"github.com/AntonShadrinNN/validatelength"
	"net/url"

	"ads-server/internal/ads"
	"ads-server/internal/users"
)

const (
	titleConst = iota
	textConst
)

type App struct {
	adRepo   AdRepository
	userRepo UserRepository
}

func validate(s string, m int) error {
	constraint := 0
	if m == titleConst {
		constraint = 100
	} else if m == textConst {
		constraint = 500
	}
	ok, err := validatelength.ValidateLen(s, 1, 0)
	if err != nil || ok {

		return errs.ValidationError
	}
	ok, err = validatelength.ValidateLen(s, 4, constraint)
	if err != nil || ok {
		return errs.ValidationError
	}
	return nil
}

// CreateAd creates new ad using repository
func (a App) CreateAd(ctx context.Context, uID int64, title string, text string) (*ads.Ad, error) {
	err := validate(title, titleConst)
	if err != nil {
		return nil, errs.ValidationError
	}
	err = validate(text, textConst)
	if err != nil {
		return nil, errs.ValidationError
	}

	ad := ads.New(uID, title, text)
	_, err = a.adRepo.Create(ctx, ad)
	if err != nil {
		return nil, errs.AccessError
	}
	return ad, nil
}

// UpdateAd updates ad using repository
func (a App) UpdateAd(ctx context.Context, adID int64, uID int64, title string, text string) (*ads.Ad, error) {
	err := validate(title, titleConst)
	if err != nil {
		return nil, err
	}
	err = validate(text, textConst)
	if err != nil {
		return nil, err
	}

	ad, err := a.adRepo.Update(ctx, adID, uID, title, text)
	if err != nil {
		return nil, errs.AccessError
	}
	return ad, nil
}

func (a App) DeleteAd(ctx context.Context, adID, uID int64) error {
	err := a.adRepo.Delete(ctx, adID, uID)
	if err != nil {
		return err
	}
	return nil
}

// PublishAd changes ad status using repository
func (a App) PublishAd(ctx context.Context, adID int64, uID int64, action bool) (*ads.Ad, error) {
	ad, err := a.adRepo.Publish(ctx, adID, uID, action)
	if err != nil {
		return nil, errs.AccessError
	}
	return ad, nil
}

// GetAdByID returns ad by ID given using repository
func (a App) GetAdByID(ctx context.Context, id int64) (*ads.Ad, error) {
	ad, err := a.adRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errs.AdNotFoundError

	}
	return ad, nil
}

// GetAdByName returns ad by name given using repository
func (a App) GetAdByName(ctx context.Context, title string) []*ads.Ad {
	return a.adRepo.GetByName(ctx, title)
}

// FindUser returns user by ID given using repository
func (a App) FindUser(ctx context.Context, id int64) (*users.User, error) {
	if user, err := a.userRepo.Get(ctx, id); err == nil {
		return user, nil
	}
	return nil, errs.UserNotFoundError
}

func (a App) DeleteUser(ctx context.Context, id int64) error {
	if err := a.userRepo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func (a App) UpdateUser(ctx context.Context, id int64, name, email string) (*users.User, error) {
	if user, err := a.userRepo.Update(ctx, id, name, email); err == nil {
		return user, nil
	} else {
		return nil, err
	}

}

// CreateUser creates a new user using repository
func (a App) CreateUser(ctx context.Context, name string, email string) (*users.User, error) {
	user := users.New(name, email)
	_, err := a.userRepo.Create(ctx, user)
	if errors.Is(err, errs.UserNotFoundError) {
		return nil, errs.UserNotFoundError
	}
	if err != nil {
		return nil, errs.AccessError
	}
	return user, nil

}

// Filter filters all ads by query params given
func (a App) Filter(ctx context.Context, params url.Values) ([]*ads.Ad, error) {
	allAds, err := a.adRepo.Filter(ctx, params)
	if err != nil {
		return nil, err
	}
	return allAds, nil
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name UserRepository
type UserRepository interface {
	Create(ctx context.Context, u *users.User) (id int64, err error)
	Update(ctx context.Context, id int64, name string, email string) (*users.User, error)
	Get(ctx context.Context, id int64) (*users.User, error)
	Delete(ctx context.Context, id int64) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name AdRepository
type AdRepository interface {
	Create(context.Context, *ads.Ad) (int64, error)
	Publish(context.Context, int64, int64, bool) (*ads.Ad, error)
	Update(context.Context, int64, int64, string, string) (*ads.Ad, error)
	Delete(context.Context, int64, int64) error
	GetByID(context.Context, int64) (*ads.Ad, error)
	GetByName(context.Context, string) []*ads.Ad
	Filter(ctx context.Context, params url.Values) ([]*ads.Ad, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.20.2 --name IApp
type IApp interface {
	CreateAd(ctx context.Context, uID int64, title string, text string) (*ads.Ad, error)
	UpdateAd(ctx context.Context, adID int64, uID int64, title string, text string) (*ads.Ad, error)
	DeleteAd(ctx context.Context, adID, uID int64) error
	PublishAd(ctx context.Context, adID int64, uID int64, action bool) (*ads.Ad, error)
	GetAdByID(ctx context.Context, id int64) (*ads.Ad, error)
	GetAdByName(ctx context.Context, title string) []*ads.Ad
	FindUser(ctx context.Context, id int64) (*users.User, error)
	DeleteUser(ctx context.Context, id int64) error
	UpdateUser(ctx context.Context, id int64, name, email string) (*users.User, error)
	CreateUser(ctx context.Context, name string, email string) (*users.User, error)
	Filter(ctx context.Context, params url.Values) ([]*ads.Ad, error)
}

func NewApp(repo AdRepository, userRepo UserRepository) App {
	return App{adRepo: repo, userRepo: userRepo}
}
