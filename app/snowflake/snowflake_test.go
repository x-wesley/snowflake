package snowflake

import (
	"testing"
)

//goos: darwin
//goarch: amd64
//pkg: snowflake/app/snowflake
//BenchmarkSnowFlake_NextId-8   	 4924948	       244 ns/op
//PASS
func BenchmarkSnowFlake_NextId(b *testing.B) {
	snowFlake := NewSnowFlake(1, 1)
	for i := 0; i < b.N; i++ {
		snowFlake.NextId()
	}
}
