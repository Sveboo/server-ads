package users

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
	"unicode/utf8"
)

func FuzzNew(f *testing.F) {
	f.Fuzz(func(t *testing.T, name string, email string) {
		u := New(name, email)
		if !utf8.ValidString(u.Name) {
			t.Fatalf("Wrong name format %s", u.Name)
		}
		if !utf8.ValidString(u.Email) {
			t.Fatalf("Wrong email format %s", u.Name)
		}
	})
}

type NewTestSuite struct {
	suite.Suite
	baseEmail string
}

func (n *NewTestSuite) SetupTest() {
	n.baseEmail = "example@gmail.com"

}
func (n *NewTestSuite) TestExample() {
	tests := []struct {
		name     string
		expected string
	}{
		{
			New("John", n.baseEmail).Name,
			"John",
		},

		{New("Привет", n.baseEmail).Name,
			"auto-generated-string",
		},
	}
	n.Equal("example@gmail.com", n.baseEmail)

	for _, tt := range tests {
		fmt.Println(tt)
		n.Equal(tt.name, tt.expected)
	}
}

func TestNewTestSuite(t *testing.T) {
	suite.Run(t, new(NewTestSuite))
}
