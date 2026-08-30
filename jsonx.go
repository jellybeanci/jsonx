package jsonx

import (
	jsonv2 "encoding/json/v2"

	stringutil "github.com/jellybeanci/jsonx/internal/string-util"
	jsoniter "github.com/json-iterator/go"
)

var jiter = jsoniter.ConfigCompatibleWithStandardLibrary

func Unmarshal(buf []byte, val any) error {
	return jsonv2.Unmarshal(buf, val)
}

func UnmarshalString(buf string, val any) error {
	return jsonv2.Unmarshal(stringutil.Byte(buf), val)
}

func Marshal(val any) ([]byte, error) {
	return jiter.Marshal(val)
}

func MarshalString(val any) (string, error) {
	res, err := jiter.Marshal(val)
	if err != nil {
		return "", err
	}
	return stringutil.String(res), nil
}

func Cast[Model any](buf []byte, err error) (*Model, error) {
	if err != nil {
		return nil, err
	}
	var result Model
	if err = Unmarshal(buf, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func CastSlice[Model any](buff []byte, err error) ([]Model, error) {
	if err != nil {
		return nil, err
	}
	var result []Model
	if err = Unmarshal(buff, &result); err != nil {
		return nil, err
	}
	return result, nil
}
