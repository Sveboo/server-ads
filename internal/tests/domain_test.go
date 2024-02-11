package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangeStatusAdOfAnotherUser(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser(0, "James", "Ostin")
	if err != nil {
		fmt.Println(err)
	}

	_, err = client.createUser(1, "Name", "Mas")
	if err != nil {
		fmt.Println(err)
	}

	resp, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(1, resp.Data.ID, true)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestUpdateAdOfAnotherUser(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser(0, "James", "Ostin")
	if err != nil {
		fmt.Println(err)
	}

	_, err = client.createUser(1, "Name", "Mas")
	if err != nil {
		fmt.Println(err)
	}

	resp, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)

	_, err = client.updateAd(1, resp.Data.ID, "title", "text")
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestCreateAd_ID(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser(0, "James", "Ostin")
	if err != nil {
		fmt.Println(err)
	}

	resp, err := client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(0))

	resp, err = client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(1))

	resp, err = client.createAd(0, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(2))
}
