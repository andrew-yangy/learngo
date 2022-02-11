package postgres

import (
	"github.com/sirupsen/logrus"
)

var store *PgStore

func NewPostgresStore() (*PgStore, error) {
	if store == nil || store.DB == nil {
		pgConnectionString := GetUrlByStage().String()
		logrus.Println(pgConnectionString)
		s, err := NewStore(pgConnectionString)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		if err := s.DB.Ping(); err != nil {
			logrus.Error(err)
			return nil, err
		}

		return s, nil
	} else {
		if err := store.DB.Ping(); err != nil {
			logrus.Errorf("pinging cached DB client failed: %+v", err)

			// try to establish a new connection if a previously cached store errored out.
			pgConnectionString := GetUrlByStage().String()
			s, err := NewStore(pgConnectionString)
			if err == nil {
				return s, nil
			}
			return nil, err
		}
		return store, nil
	}
}
