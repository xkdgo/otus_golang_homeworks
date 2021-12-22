package hw10programoptimization

import (
	"archive/zip"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkGetDomainStatFile(b *testing.B) {
	r, err := zip.OpenReader("testdata/users.dat.zip")
	require.NoError(b, err)
	defer r.Close()
	for i := 0; i < b.N; i++ {
		data, err := r.File[0].Open()
		require.NoError(b, err)
		b.StartTimer()
		GetDomainStat(data, "biz")
		b.StopTimer()
	}
}
