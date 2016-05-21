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

func (s *Store) Clear() error {
	return s.DB.Update(func(tx *bolt.Tx) error {
		tx.DeleteBucket(s.BucketName())
		_, err := tx.CreateBucketIfNotExists(s.BucketName())
		return err
	})
}

func (s *Store) Save(is issue.Issue) error {
	return s.DB.Update(func(tx *bolt.Tx) error {
		pb, _ := tx.CreateBucketIfNotExists(s.BucketName())
		b, _ := pb.CreateBucketIfNotExists([]byte(*is.State))
		data, err := json.Marshal(is)
		if err != nil {
			return err
		}
		err = b.Put([]byte(strconv.Itoa(*is.Number)), data)
		if err != nil {
			return err
		}
		inb, _ := pb.CreateBucketIfNotExists([]byte("_map"))
		return inb.Put([]byte(strconv.Itoa(*is.Number)), []byte(*is.State))
	})
}

func (s *Store) Get(number string) (issue.Issue, error) {
	id := []byte(number)
	i := issue.Issue{}
	err := s.DB.View(func(tx *bolt.Tx) error {
		pb := tx.Bucket(s.BucketName())
		inb := pb.Bucket([]byte("_map"))
		bn := inb.Get(id)
		b := pb.Bucket(bn)
		v := b.Get(id)
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

	err := s.DB.View(func(tx *bolt.Tx) error {
		for _, state := range []string{"open", "closed"} {
			si, err := s.AllByState(state)
			if err != nil {
				return err
			}
			issues = append(issues, si...)
		}
		return nil
	})
	return issues, err
}

func (s *Store) AllByState(state string) ([]issue.Issue, error) {
	issues := []issue.Issue{}

	s.DB.View(func(tx *bolt.Tx) error {
		pb := tx.Bucket(s.BucketName())
		b := pb.Bucket([]byte(state))
		if b == nil {
			// the bucket doesn't exist, possibly because there are no
			// issues with this state. exit and move on.
			return nil
		}
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
