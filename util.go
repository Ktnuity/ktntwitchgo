package ktntwitchgo

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"reflect"
	"strings"
)

func DecodeDataInstanceBytes[T any](data []byte) (*T, error) {
	var result T
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return &result, nil
}

func EncodeDataInstanceBytes[T any](value *T) ([]byte, error) {
	result, err := json.Marshal(*value)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal instance: %w", err)
	}

	return result, nil
}

func DecodeDataInstanceString[T any](data string) (*T, error) {
	return DecodeDataInstanceBytes[T]([]byte(data))
}

func EncodeDataInstanceString[T any](value *T) (string, error) {
	result, err := json.Marshal(*value)
	if err != nil { return "", err }
	return string(result), nil
}

type LocalCache struct {
	AccessToken		string		`json:"access_token"`
	RefreshToken	string		`json:"refresh_token"`
	ClientID		string		`json:"client_id"`
	ClientSecret	string		`json:"client_secret"`
}

var userFile = "./data/apiUser.json"

func LoadLocalCache() (*LocalCache, error) {
	data, err := os.ReadFile(userFile)
	if err != nil {
		return nil, err
	}

	obj, err := DecodeDataInstanceBytes[LocalCache](data)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func GetLocalAccessToken() (string, error) {
	cache, err := LoadLocalCache()
	if err != nil {
		return "", err
	}

	return cache.AccessToken, nil
}

func GetLocalRefreshToken() (string, error) {
	cache, err := LoadLocalCache()
	if err != nil {
		return "", err
	}

	return cache.RefreshToken, nil
}

func GetLocalClientID() (string, error) {
	cache, err := LoadLocalCache()
	if err != nil {
		return "", err
	}

	return cache.ClientID, nil
}

func GetLocalClientSecret() (string, error) {
	cache, err := LoadLocalCache()
	if err != nil {
		return "", err
	}

	return cache.ClientSecret, nil
}

func ParseOptions[T any](options T) string {
	params := url.Values{}

	v := reflect.ValueOf(options)
	t := reflect.TypeOf(options)

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}

		v = v.Elem()
		t = t.Elem()
	}

	if v.Kind() != reflect.Struct {
		return ""
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		if !field.CanInterface() {
			continue
		}

		if (field.Kind() == reflect.Ptr && field.IsNil()) ||
		   (field.Kind() == reflect.Interface && field.IsNil()) {
			continue
		}

		key := fieldType.Name

		if field.Kind() == reflect.Slice && !field.IsNil() {
			for j := 0; j < field.Len(); j++ {
				element := field.Index(j)
				if element.CanInterface() {
					params.Add(key, fmt.Sprintf("%v", element.Interface()))
				}
			}
		} else {
			params.Set(key, fmt.Sprintf("%v", field.Interface()))
		}
	}

	encoded := params.Encode()

	return strings.TrimSuffix(encoded, "&")
}

