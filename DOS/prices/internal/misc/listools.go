package misc

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"
)

// Function to generate a unique hash for the tuple
func GenerateHash(priceTuple map[string]string) string {
	tupleString := priceTuple["list"] + priceTuple["sku"] + priceTuple["unit"] + priceTuple["material"] + priceTuple["tservicio"]
	hash := md5.Sum([]byte(tupleString))
	return hex.EncodeToString(hash[:])
}

// GenerateNameWithTimestamp generates a name in the format "prefix-current_timestamp".
func GenerateNameWithTimestamp(prefix string) string {

	return fmt.Sprintf("%s-%d", prefix, time.Now().Unix())
}
