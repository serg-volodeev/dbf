package xbase

import (
	"fmt"
	"time"
)

func interfaceToString(v interface{}) (string, error) {
	result, ok := v.(string)
	if !ok {
		return "", fmt.Errorf("error convert interface{} to string")
	}
	return result, nil
}

func interfaceToBool(v interface{}) (bool, error) {
	result, ok := v.(bool)
	if !ok {
		return false, fmt.Errorf("error convert interface{} to bool")
	}
	return result, nil
}

func interfaceToInt(value interface{}) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		return int64(v), nil
	}
	return 0, fmt.Errorf("error convert interface{} to int")
}

func interfaceToFloat(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float32:
		return float64(v), nil
	case float64:
		return float64(v), nil
	}
	return 0, fmt.Errorf("error convert interface{} to float")
}

func interfaceToDate(v interface{}) (time.Time, error) {
	result, ok := v.(time.Time)
	if !ok {
		return result, fmt.Errorf("error convert interface{} to date")
	}
	return result, nil
}
