package util

import (
	"hash/fnv"
	"strconv"

	"urlshortener/model"
)

func GetHash(s string) model.KvPair {
	hash := fnv.New32a()
	hash.Write([]byte(s))
	sum := hash.Sum32()

	return model.KvPair{Shortened: strconv.FormatUint(uint64(sum), 36), Original: s}
}
