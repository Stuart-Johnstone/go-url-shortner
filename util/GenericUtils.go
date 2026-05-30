package util

import (
	"hash/fnv"
	"strconv"
	"urlshortener/model"
)

func GetHash(s string) KvPair {
	hash := fnv.New32a()
	hash.Write([]byte(s))
	sum := hash.Sum32()

	return KvPair{shortened: strconv.FormatUint(uint64(sum), 36), original: s}
}
