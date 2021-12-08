package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return newCountDomains(r, domain)
}

type Emails struct {
	Email string
}

func newCountDomains(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	b := strings.Builder{}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var user Emails
		line := scanner.Bytes()
		if err := json.Unmarshal(line, &user); err != nil {
			return nil, err
		}
		b.WriteRune('.')
		b.WriteString(domain)
		afterAt := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
		matched := strings.HasSuffix(afterAt, b.String())
		b.Reset()

		if matched {
			result[afterAt]++
		}
	}
	return result, nil
}
