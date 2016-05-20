package issue

import (
	"fmt"
	"time"

	"github.com/google/go-github/github"
)

type Issue struct {
	github.Issue
}

func (i Issue) FmtTitle() string {
	return fmt.Sprintf("%d\t%s\n", *i.Number, *i.Title)
}

func (i Issue) FmtByLine() string {
	return fmt.Sprintf("\tCreated %s by %s\n", i.CreatedAt.In(time.Local), *i.User.Login)
}
