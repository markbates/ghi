package store

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/boltdb/bolt"
	"github.com/markbates/ghi/cmd/issue"
	"github.com/mitchellh/go-homedir"
)

type Store struct {
	Owner string
	Repo  string
	DB    *bolt.DB
}

func (s Store) BucketName() []byte {
	return []byte(fmt.Sprintf("%s-%s", s.Owner, s.Repo))
}

func New(owner, repo string) (*Store, error) {
	s := &Store{Owner: owner, Repo: repo}
	db, err := bolt.Open(location(), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return s, err
	}
	s.DB = db
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(s.BucketName())
		return err
	})
	return s, err
}

func (s *Store) Persist(issues []issue.Issue) error {
	// clear the bucket:
	s.DB.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket(s.BucketName())
		b, _ := tx.CreateBucketIfNotExists(s.BucketName())
		for _, issue := range issues {
			data, err := json.Marshal(issue)
			if err != nil {
				return err
			}
			b.Put([]byte(strconv.Itoa(*issue.Number)), data)
		}
		return nil
	})

	return nil
}

func (s *Store) Get(number string) (issue.Issue, error) {
	i := issue.Issue{}
	err := s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.BucketName())
		v := b.Get([]byte(number))
		if v == nil {
			return fmt.Errorf("Issue #%s was not found!", number)
		}
		err := json.Unmarshal(v, &i)
		return err
	})
	return i, err
}

func (s *Store) All() ([]issue.Issue, error) {
	issues := []issue.Issue{}

	s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(s.BucketName())
		return b.ForEach(func(k, v []byte) error {
			i := issue.Issue{}
			err := json.Unmarshal(v, &i)
			if err != nil {
				return err
			}
			issues = append(issues, i)
			return nil
		})
	})
	return issues, nil
}

func location() string {
	dir, _ := homedir.Dir()
	dir, _ = homedir.Expand(dir)
	return fmt.Sprintf("%s/.issues.db", dir)
}
