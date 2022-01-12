package helper

import (
	"context"

	"github.com/pkg/errors"
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage/sql"
)

func InitStorage(storagetype string, dsn string) (storage.Storage, error) {
	if storagetype == "in-memory" {
		return memorystorage.New(), nil
	}
	if storagetype == "sql" {
		storage := sqlstorage.New()
		err := storage.Connect(context.Background(), dsn)
		if err != nil {
			return nil, errors.Wrapf(err, ": init storage error with dsn= %s", dsn)
		}
		return storage, nil
	}
	return nil, storage.ErrUnkownTypeOfStorage
}
