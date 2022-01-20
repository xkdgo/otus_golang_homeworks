package utilstorage

import "github.com/google/uuid"

func GenerateUUID() string {
	id := uuid.New()
	return id.String()
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
