package misc

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

// ISO4217Symbols is a lookup table for ISO 4217 currency symbols.
var symbolsISO4217 = map[string]string{
	"AED": "د.إ",
	"AFN": "؋",
	"ALL": "L",
	"AMD": "֏",
	"ANG": "ƒ",
	"AOA": "Kz",
	"ARS": "$",
	"AUD": "$",
	"AWG": "ƒ",
	"AZN": "₼",
	"BAM": "KM",
	"BBD": "$",
	"BDT": "৳",
	"BGN": "лв",
	"BHD": ".د.ب",
	"BIF": "FBu",
	"BMD": "$",
	"BND": "$",
	"BOB": "Bs.",
	"BRL": "R$",
	"BSD": "$",
	"BTN": "Nu.",
	"BWP": "P",
	"BYN": "Br",
	"BZD": "$",
	"CAD": "$",
	"CDF": "FC",
	"CHF": "Fr.",
	"CLP": "$",
	"CNY": "¥",
	"COP": "$",
	"CRC": "₡",
	"CUP": "$",
	"CVE": "$",
	"CZK": "Kč",
	"DJF": "Fdj",
	"DKK": "kr",
	"DOP": "$",
	"DZD": "د.ج",
	"EGP": "£",
	"ERN": "Nfk",
	"ETB": "Br",
	"EUR": "€",
	"FJD": "$",
	"FKP": "£",
	"FOK": "kr",
	"GBP": "£",
	"GEL": "₾",
	"GGP": "£",
	"GHS": "₵",
	"GIP": "£",
	"GMD": "D",
	"GNF": "FG",
	"GTQ": "Q",
	"GYD": "$",
	"HKD": "$",
	"HNL": "L",
	"HRK": "kn",
	"HTG": "G",
	"HUF": "Ft",
	"IDR": "Rp",
	"ILS": "₪",
	"IMP": "£",
	"INR": "₹",
	"IQD": "ع.د",
	"IRR": "﷼",
	"ISK": "kr",
	"JEP": "£",
	"JMD": "$",
	"JOD": "د.ا",
	"JPY": "¥",
	"KES": "Sh",
	"KGS": "сом",
	"KHR": "៛",
	"KID": "$",
	"KMF": "CF",
	"KRW": "₩",
	"KWD": "د.ك",
	"KYD": "$",
	"KZT": "₸",
	"LAK": "₭",
	"LBP": "ل.ل",
	"LKR": "Rs",
	"LRD": "$",
	"LSL": "L",
	"LYD": "ل.د",
	"MAD": "د.م.",
	"MDL": "L",
	"MGA": "Ar",
	"MKD": "ден",
	"MMK": "K",
	"MNT": "₮",
	"MOP": "P",
	"MRU": "UM",
	"MUR": "₨",
	"MVR": "Rf",
	"MWK": "MK",
	"MXN": "$",
	"MYR": "RM",
	"MZN": "MT",
	"NAD": "$",
	"NGN": "₦",
	"NIO": "C$",
	"NOK": "kr",
	"NPR": "₨",
	"NZD": "$",
	"OMR": "ر.ع.",
	"PAB": "B/.",
	"PEN": "S/",
	"PGK": "K",
	"PHP": "₱",
	"PKR": "₨",
	"PLN": "zł",
	"PYG": "₲",
	"QAR": "ر.ق",
	"RON": "lei",
	"RSD": "дин.",
	"RUB": "₽",
	"RWF": "FRw",
	"SAR": "ر.س",
	"SBD": "$",
	"SCR": "₨",
	"SDG": "ج.س.",
	"SEK": "kr",
	"SGD": "$",
	"SHP": "£",
	"SLL": "Le",
	"SOS": "Sh",
	"SRD": "$",
	"SSP": "£",
	"STN": "Db",
	"SYP": "£",
	"SZL": "L",
	"THB": "฿",
	"TJS": "ЅМ",
	"TMT": "T",
	"TND": "د.ت",
	"TOP": "T$",
	"TRY": "₺",
	"TTD": "$",
	"TVD": "$",
	"TWD": "NT$",
	"TZS": "Sh",
	"UAH": "₴",
	"UGX": "Sh",
	"USD": "$",
	"UYU": "$",
	"UZS": "сўм",
	"VES": "Bs.",
	"VND": "₫",
	"VUV": "Vt",
	"WST": "T",
	"XAF": "FCFA",
	"XCD": "$",
	"XOF": "CFA",
	"XPF": "₣",
	"YER": "﷼",
	"ZAR": "R",
	"ZMW": "ZK",
	"ZWL": "$",
}

// Function to generate a unique hash for the tuple
func GenerateHash(priceTuple map[string]string) string {
	tupleString := priceTuple["list"] + priceTuple["sku"] + priceTuple["unit"] + priceTuple["material"] + priceTuple["tservicio"]
	hash := md5.Sum([]byte(tupleString))
	return hex.EncodeToString(hash[:])
}

// GenerateNameWithTimestamp generates a name in the format "prefix-current_timestamp".
func GenerateNameWithTimestamp(prefix string) string {

	return fmt.Sprintf("%d-%s", time.Now().Unix(), prefix)
}

func GenerateNameWithCurrency(currency string, prefix string) string {

	return fmt.Sprintf("%s-%s", currency, prefix)
}

// ValidateISO4217Code checks if a given currency code exists in the ISO4217SymbolsSet.
// Returns nil if valid, otherwise returns an error.
func ValidateISO4217Code(currencyCode string) (string, error) {
	val, exists := symbolsISO4217[currencyCode]
	if !exists {
		return "", errors.New(fmt.Sprintf("invalid ISO 4217 currency code: %s", currencyCode))
	}
	return val, nil
}
