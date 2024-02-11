package tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
	"time"
)

func TestCreateUser(t *testing.T) {
	client := getTestClient()
	user, err := client.createUser(0, "James", "Ostin")

	assert.NoError(t, err)
	assert.Equal(t, "James", user.Data.Name)
	assert.Equal(t, "Ostin", user.Data.Email)
}

func TestCreateAd(t *testing.T) {
	client := getTestClient()
	_, err := client.createUser(0, "James", "Ostin")
	if err != nil {
		fmt.Println(err)
	}
	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	assert.Zero(t, response.Data.ID)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")
	assert.Equal(t, response.Data.AuthorID, int64(0))
	assert.False(t, response.Data.Published)
}

func TestChangeAdStatus(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser(0, "James", "Ostin")
	if err != nil {
		fmt.Println(err)
	}
	_, err = client.createUser(1, "James2", "Ostin2")
	if err != nil {
		fmt.Println(err)
	}

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	response, err = client.changeAdStatus(0, response.Data.ID, true)
	assert.NoError(t, err)
	assert.True(t, response.Data.Published)

	response, err = client.changeAdStatus(0, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)

	response, err = client.changeAdStatus(0, response.Data.ID, false)
	assert.NoError(t, err)
	assert.False(t, response.Data.Published)

	response, err = client.changeAdStatus(1, response.Data.ID, true)
	assert.Error(t, err)
	assert.False(t, response.Data.Published)

	response, err = client.changeAdStatus(1, 0, true)
	assert.Error(t, err)
}

func TestUpdateAd(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser(0, "James", "Ostin")
	if err != nil {
		fmt.Println(err)
	}

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	response, err = client.updateAd(0, response.Data.ID, "привет", "мир")
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Title, "привет")
	assert.Equal(t, response.Data.Text, "мир")

	response, err = client.updateAd(1, response.Data.ID, "привет", "мир")
	assert.Error(t, err)
}

func TestUpdateAdNoUser(t *testing.T) {
	client := getTestClient()

	//response, err := client.createAd(0, "hello", "world")
	//assert.NoError(t, err)

	response, err := client.updateAd(0, 0, "привет", "мир")
	fmt.Println(response)
	fmt.Println(err)
	assert.ErrorIs(t, err, ErrBadRequest)
	assert.Empty(t, response)
}

func TestListAds(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser(0, "James", "Ostin")
	if err != nil {
		fmt.Println(err)
	}

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	publishedAd, err := client.changeAdStatus(0, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(0, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.listAds()
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, publishedAd.Data.ID)
	assert.Equal(t, ads.Data[0].Title, publishedAd.Data.Title)
	assert.Equal(t, ads.Data[0].Text, publishedAd.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, publishedAd.Data.AuthorID)
	assert.True(t, ads.Data[0].Published)
}

// Тест получения объявления по ID
func TestAdBYID(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser(0, "James", "Ostin")
	if err != nil {
		fmt.Println(err)
	}

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	response, err = client.getAdByID(response.Data.ID)
	assert.NoError(t, err)
	assert.Equal(t, response.Data.Title, "hello")
	assert.Equal(t, response.Data.Text, "world")

}

// Тест получения объявления по названию
func TestAdsByTitle(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser(0, "James", "Ostin")
	if err != nil {
		fmt.Println(err)
	}

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(0, response.Data.ID, true)
	assert.NoError(t, err)

	response, err = client.createAd(0, "hello", "not for sale")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(0, response.Data.ID, true)
	assert.NoError(t, err)

	response, err = client.createAd(0, "not hello", "not for sale")
	assert.NoError(t, err)

	ads, err := client.adsWithFilters("?title=hello")
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 2)
}

// Тест фильтрации: по автору
func TestAuthorFilter(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser(0, "James", "Ostin")
	if err != nil {
		fmt.Println(err)
	}

	_, err = client.createUser(1, "Jeremy", "gmail")
	if err != nil {
		fmt.Println(err)
	}

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	publishedAd, err := client.changeAdStatus(0, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(1, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.adsWithFilters("?author=0")
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, publishedAd.Data.ID)
	assert.Equal(t, ads.Data[0].Title, publishedAd.Data.Title)
	assert.Equal(t, ads.Data[0].Text, publishedAd.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, publishedAd.Data.AuthorID)
	assert.True(t, ads.Data[0].Published)
}

// Тест фильтрации: по дате
func TestDateFilter(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser(0, "James", "Ostin")
	if err != nil {
		log.Println(err)
	}

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(0, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(0, "best cat", "not for sale")
	assert.NoError(t, err)

	year, month, day := time.Now().UTC().Date()
	date := fmt.Sprintf("%d-%d-%d", year, int(month), day)
	// published не установлен, поэтому выведутся оба объявления
	ads, err := client.adsWithFilters(fmt.Sprintf("?date=%s", date))
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 2)
}

// Тест фильтрации: нет фильтров
func TestNoFilters(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser(0, "James", "Ostin")
	if err != nil {
		log.Println(err)
	}

	response, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(0, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(0, "best cat", "not for sale")
	assert.NoError(t, err)

	ads, err := client.adsWithFilters("")
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 2)
}

// Тест фильтрации: все фильтры
func TestAllFilters(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser(0, "James", "Ostin")
	if err != nil {
		log.Println(err)
	}

	_, err = client.createUser(1, "Jeremy", "mail")
	if err != nil {
		log.Println(err)
	}

	response, err := client.createAd(0, "hello", "world1")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(0, response.Data.ID, true)
	assert.NoError(t, err)

	response, err = client.createAd(0, "hello", "world2")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(0, response.Data.ID, true)
	assert.NoError(t, err)

	_, err = client.createAd(1, "best cat", "not for sale")
	assert.NoError(t, err)

	year, month, day := time.Now().UTC().Date()
	date := fmt.Sprintf("%d-%d-%d", year, int(month), day)
	ads, err := client.adsWithFilters("?published=true&author=0&title=hello;date=" + date)
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 2)
}
