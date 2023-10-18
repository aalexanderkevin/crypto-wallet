package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/mr-tron/base58"
)

type NumberType interface {
	~string | int | int32 | int64 | float32 | float64 | bool
}

func Pointer[T any](data T) *T {
	return &data
}

func Val[T any](pointer *T) T {
	if pointer == nil {
		return *new(T)
	}

	return *pointer
}

func ToString[T NumberType](data *T) *string {
	if data == nil {
		return Pointer("")
	}
	stringData := fmt.Sprintf("%v", *data)
	return &stringData
}

func ToInt64[T NumberType](data *T) (result *int64) {
	if data == nil {
		return
	}
	switch v := reflect.ValueOf(data); v.Elem().Kind() {
	case reflect.String:
		intVal, err := strconv.ParseInt(v.Elem().String(), 10, 64)
		if err != nil {
			result = nil
			return
		}
		result = Pointer(intVal)
	case reflect.Int, reflect.Int32, reflect.Int64:
		result = Pointer(v.Elem().Int())
	case reflect.Float32, reflect.Float64:
		result = Pointer(int64(v.Elem().Float()))
	case reflect.Bool:
		if v.Elem().Bool() {
			result = Pointer[int64](1)
		} else {
			result = Pointer[int64](0)
		}
	}
	return
}

func ToInt[T NumberType](data *T) (result *int) {
	if data == nil {
		return
	}
	switch v := reflect.ValueOf(data); v.Elem().Kind() {
	case reflect.String:
		intVal, err := strconv.Atoi(v.Elem().String())
		if err != nil {
			result = nil
		}
		result = Pointer(intVal)
	case reflect.Int, reflect.Int32, reflect.Int64:
		result = Pointer(int(v.Elem().Int()))
	case reflect.Float32, reflect.Float64:
		result = Pointer(int(v.Elem().Float()))
	case reflect.Bool:
		if v.Elem().Bool() {
			result = Pointer(1)
		} else {
			result = Pointer(0)
		}
	}
	return
}

func ValTimeUnix(val *time.Time) int64 {
	if val == nil {
		return 0
	}

	return val.Unix()
}

func ValOrDefault[T NumberType](value *T, defaultVal T) T {
	if value == nil || *value == *new(T) {
		return defaultVal
	}

	return *value
}

// ValYearMonthUnix trim time to year and mont and convert it to unix format
func ValYearMonthUnix(val *time.Time) int64 {
	if val == nil {
		return 0
	}

	return time.Date(val.Year(), val.Month(), 1, 0, 0, 0, 0, time.UTC).Unix()
}

func EqualPointerValue[T comparable](a *T, b *T) bool {
	return Val(a) == Val(b)
}

func TimeToMilisecond(t *time.Time) int64 {
	return int64(t.UnixNano()) / int64(time.Millisecond)
}

func WeiToEther(wei *big.Int) *float64 {
	fwei := new(big.Float)
	fwei.SetString(wei.String())
	ethValue := new(big.Float).Quo(fwei, big.NewFloat(math.Pow10(18)))
	res, _ := ethValue.Float64()
	return &res
}

func HexToAddress(hex string) common.Address {
	return common.HexToAddress(hex)
}

func EthToWei(amount float64) (bn *big.Int) {
	bf := new(big.Float).Mul(big.NewFloat(amount), big.NewFloat(1000000000000000000))
	bn, _ = bf.Int(bn)
	return bn
}

func ToTrxAddress(publicAddress string) *string {
	addb, _ := hex.DecodeString(publicAddress)
	hash11 := s256(s256(addb))
	secrets := hash11[:4]
	addb = append(addb, secrets...)

	return Pointer(base58.Encode(addb))
}

func s256(s []byte) []byte {
	h := sha256.New()
	h.Write(s)
	bs := h.Sum(nil)
	return bs
}
