package uuid

import (
	"math/rand"
	"strconv"
	"strings"
)

const UPPER_LIMIT = 16

func V4() string {
	ds := [32]byte{}

	for i := range ds {
		n := rand.Float32()
		ds[i] = byte(n * UPPER_LIMIT)
	}

	ds[16] = 0b1000 | (ds[16] & 0b0011)
	ds[12] = 0b0100

	as := [32]string{}
	for i, v := range ds {
		as[i] = strconv.FormatInt(int64(v), 16)
	}

	return strings.Join(as[0:8], "") + "-" + strings.Join(as[8:12], "") + "-" + strings.Join(as[12:16], "") + "-" +
		strings.Join(as[16:20], "") + "-" + strings.Join(as[20:], "")
}
