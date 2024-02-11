package repo

import (
	"ads-server/internal/app"
	"ads-server/internal/errs"
	"ads-server/internal/users"
	"context"
	"sync"
)

type UsersRepo struct {
	storage map[int64]*users.User
	mx      *sync.Mutex
	lastID  int64
}

// Create creates a new user
func (ur *UsersRepo) Create(_ context.Context, u *users.User) (id int64, err error) {
	ur.mx.Lock()
	defer ur.mx.Unlock()
	if _, ok := ur.storage[ur.lastID]; ok {
		return -1, errs.UserNotFoundError
	}
	u.ID = ur.lastID
	ur.storage[u.ID] = u
	ur.lastID++

	return ur.lastID - 1, nil
}

// Update updates an existing user
func (ur *UsersRepo) Update(_ context.Context, id int64, name string, email string) (*users.User, error) {
	ur.mx.Lock()
	defer ur.mx.Unlock()
	if _, ok := ur.storage[id]; !ok {
		return nil, errs.UserNotFoundError
	}

	ur.storage[id].Name = name
	ur.storage[id].Email = email
	return ur.storage[id], nil
}

// Delete deletes user from storage
func (ur *UsersRepo) Delete(_ context.Context, id int64) error {
	ur.mx.Lock()
	defer ur.mx.Unlock()

	if _, ok := ur.storage[id]; ok {
		delete(ur.storage, id)
		return nil
	}

	return errs.UserNotFoundError

}

// Get returns a user by ID given
func (ur *UsersRepo) Get(_ context.Context, id int64) (*users.User, error) {
	ur.mx.Lock()
	defer ur.mx.Unlock()
	if user, ok := ur.storage[id]; ok {
		return user, nil
	}
	return nil, errs.UserNotFoundError
}

// NewUser is a constructor
func NewUser() app.UserRepository {
	return &UsersRepo{
		mx:      &sync.Mutex{},
		storage: make(map[int64]*users.User, 1),
		lastID:  0,
	}
}
