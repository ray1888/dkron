package dkron

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/tidwall/buntdb"
)

const (
	// MaxExecutions to maintain in the storage

	userPrefix     = "user"
	jwtTokenPrefix = "session:user:"
)

// var secretKey = []byte("your-secret-key") // 替换为你的密钥
var (
	DefaultTTL = 7 * 24 * time.Hour
)

type UserEncryptPass struct {
	Username    string
	EncrtpyPass string
}

func (s *Store) setUserTxFunc(pbj *UserEncryptPass) func(tx *buntdb.Tx) error {
	return func(tx *buntdb.Tx) error {
		userKey := fmt.Sprintf("%s:%s", userPrefix, pbj.Username)

		s.logger.WithField("user", pbj.Username).Debug("store: Setting user")
		tString, err := json.Marshal(&pbj)
		if err != nil {
			return err
		}
		if _, _, err := tx.Set(userKey, string(tString), nil); err != nil {
			return err
		}
		return nil
	}
}

func (s *Store) getUserTxFunc(username string, pbj *UserEncryptPass) func(tx *buntdb.Tx) error {
	return func(tx *buntdb.Tx) error {
		item, err := tx.Get(fmt.Sprintf("%s:%s", userPrefix, username))
		if err != nil {
			return err
		}

		if err := json.Unmarshal([]byte(item), pbj); err != nil {
			return err
		}

		s.logger.WithFields(logrus.Fields{
			"user": username,
		}).Debug("store: Retrieved session from datastore")

		return nil
	}
}

func (s *Store) deleteUserTxFunc(username string) func(tx *buntdb.Tx) error {
	return func(tx *buntdb.Tx) error {
		var delkeys []string
		prefix := fmt.Sprintf("%s:%s", userPrefix, username)
		if err := tx.Ascend("", func(key, value string) bool {
			if strings.HasPrefix(key, prefix) {
				delkeys = append(delkeys, key)
			}
			return true
		}); err != nil {
			return err
		}

		for _, k := range delkeys {
			_, err := tx.Delete(k)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func (s *Store) getUser(name string) (*UserEncryptPass, error) {
	uep := new(UserEncryptPass)
	err := s.db.View(s.getUserTxFunc(name, uep))
	if err != nil {
		return uep, err
	}
	return uep, nil
}

func (s *Store) CheckUserExist(name string) (bool, error) {
	data, err := s.getUser(name)
	if err != nil {
		return false, err
	}
	if data != nil {
		return true, nil
	}
	return false, nil
}

func (s *Store) GetUserWithPassword(name, password string) (bool, error) {
	data, err := s.getUser(name)
	if err != nil || data == nil {
		return false, err
	}
	if data.EncrtpyPass != password {
		return false, nil
	}
	return true, nil
}

func (s *Store) AddUser(name string, password string) error {
	data, err := s.getUser(name)
	if err != nil {
		return err
	}
	if data != nil {
		s.logger.Infof("user already added, user is %s", name)
		return nil
	}
	err = s.db.Update(func(tx *buntdb.Tx) error {
		if err := s.setUserTxFunc(&UserEncryptPass{
			Username:    name,
			EncrtpyPass: password,
		})(tx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) DeleteUser(name string) error {
	data, err := s.getUser(name)
	if err != nil {
		return err
	}
	if data == nil {
		s.logger.Infof("user not exist, user is %s", name)
		return nil
	}
	err = s.db.Update(func(tx *buntdb.Tx) error {
		if err := s.deleteUserTxFunc(name)(tx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

type SessionToken struct {
	Token string
}

func (s *Store) setSessionTxFunc(username string, token SessionToken) func(tx *buntdb.Tx) error {
	return func(tx *buntdb.Tx) error {
		sessionKey := fmt.Sprintf("%s:%s", jwtTokenPrefix, username)

		s.logger.WithField("user", username).Debug("store: Setting user session")
		tString, err := json.Marshal(&token)
		if err != nil {
			return err
		}
		if _, _, err := tx.Set(sessionKey, string(tString), &buntdb.SetOptions{
			TTL: DefaultTTL,
		}); err != nil {
			return err
		}
		return nil
	}
}

func (s *Store) deleteSessionTxFunc(username string) func(tx *buntdb.Tx) error {
	return func(tx *buntdb.Tx) error {
		delUserKey := fmt.Sprintf("%s:%s", jwtTokenPrefix, username)

		s.logger.WithField("user", username).Debug("delete: session user")

		_, err := tx.Delete(delUserKey)

		return err
	}
}

func (s *Store) getSessionTxFunc(name string, pbj *SessionToken) func(tx *buntdb.Tx) error {
	return func(tx *buntdb.Tx) error {
		item, err := tx.Get(fmt.Sprintf("%s:%s", jwtTokenPrefix, name))
		if err != nil {
			return err
		}

		if err := json.Unmarshal([]byte(item), pbj); err != nil {
			return err
		}

		s.logger.WithFields(logrus.Fields{
			"user": name,
		}).Debug("store: Retrieved session from datastore")

		return nil
	}
}

func (s *Store) AddSession(name string, jwtToken SessionToken) error {
	err := s.db.Update(func(tx *buntdb.Tx) error {
		if err := s.setSessionTxFunc(name, jwtToken)(tx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) DeleteSession(name string, jwtToken SessionToken) error {
	err := s.db.Update(func(tx *buntdb.Tx) error {
		if err := s.deleteSessionTxFunc(name)(tx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetSession(name string) (SessionToken, error) {
	var token SessionToken
	err := s.db.View(s.getSessionTxFunc(name, &token))
	if err != nil {
		return token, err
	}
	if err != nil {
		return token, err
	}

	return token, nil
}
