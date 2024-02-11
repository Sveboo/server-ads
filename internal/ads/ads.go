package ads

import (
	"time"
)

type Ad struct {
	ID        int64
	Title     string
	Text      string
	AuthorID  int64
	CDate     time.Time
	UDate     time.Time
	Published bool
}

func New(aID int64, title string, text string) *Ad {
	creationTime := time.Now().UTC()
	return &Ad{
		ID:        0,
		Title:     title,
		Text:      text,
		AuthorID:  aID,
		CDate:     creationTime,
		UDate:     creationTime,
		Published: false,
	}
}
