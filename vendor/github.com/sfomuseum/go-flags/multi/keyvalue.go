package multi

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const SEP string = "="

type KeyValueStringFlag struct {
	Key   string
	Value string
}

type KeyValueString []*KeyValueStringFlag

func (e *KeyValueString) String() string {
	return fmt.Sprintf("%v", *e)
}

func (e *KeyValueString) Set(value string) error {

	value = strings.Trim(value, " ")
	kv := strings.Split(value, SEP)

	if len(kv) != 2 {
		return errors.New("Invalid key=value argument")
	}

	a := KeyValueStringFlag{
		Key:   kv[0],
		Value: kv[1],
	}

	*e = append(*e, &a)
	return nil
}

func (e *KeyValueString) Get() interface{} {
	return *e
}

type KeyValueInt64Flag struct {
	Key   string
	Value int64
}

type KeyValueInt64 []*KeyValueInt64Flag

func (e *KeyValueInt64) String() string {
	return fmt.Sprintf("%v", *e)
}

func (e *KeyValueInt64) Set(value string) error {

	value = strings.Trim(value, " ")
	kv := strings.Split(value, SEP)

	if len(kv) != 2 {
		return errors.New("Invalid key=value argument")
	}

	v, err := strconv.ParseInt(kv[1], 10, 64)

	if err != nil {
		return err
	}

	a := KeyValueInt64Flag{
		Key:   kv[0],
		Value: v,
	}

	*e = append(*e, &a)
	return nil
}

func (e *KeyValueInt64) Get() interface{} {
	return *e
}

type KeyValueFloat64Flag struct {
	Key   string
	Value float64
}

type KeyValueFloat64 []*KeyValueFloat64Flag

func (e *KeyValueFloat64) String() string {
	return fmt.Sprintf("%v", *e)
}

func (e *KeyValueFloat64) Set(value string) error {

	value = strings.Trim(value, " ")
	kv := strings.Split(value, SEP)

	if len(kv) != 2 {
		return errors.New("Invalid key=value argument")
	}

	v, err := strconv.ParseFloat(kv[1], 64)

	if err != nil {
		return err
	}

	a := KeyValueFloat64Flag{
		Key:   kv[0],
		Value: v,
	}

	*e = append(*e, &a)
	return nil
}

func (e *KeyValueFloat64) Get() interface{} {
	return *e
}
