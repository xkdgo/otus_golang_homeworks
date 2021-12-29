package helper

import (
	"github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/xkdgo/otus_golang_homeworks/hw12_13_14_15_calendar/internal/storage/memory"
)

func InitStorage(storagetype string, dsn string) (storage.Storage, error) {
	if storagetype == "in-memory" {
		return memorystorage.New(), nil
	}
	return nil, storage.ErrUnkownTypeOfStorage
}
