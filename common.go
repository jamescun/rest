package rest

import (
	"net/url"
	"strconv"
)

type Id string

func (id Id) Value(v url.Values) url.Values {
	v.Add("id", string(id))
	return v
}

type Length int

func (l Length) Value(v url.Values) url.Values {
	v.Add("length", strconv.Itoa(int(l)))
	return v
}

type Offset int

func (o Offset) Value(v url.Values) url.Values {
	v.Add("offset", strconv.Itoa(int(o)))
	return v
}

type param struct {
	key, value string
}

func (p param) Value(v url.Values) url.Values {
	v.Add(p.key, p.value)
	return v
}

func Param(key, value string) QueryParam {
	return param{key, value}
}
