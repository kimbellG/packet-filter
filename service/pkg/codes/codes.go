package codes

type Code int

const (
	Unknown Code = iota

	JSONEncodeError
)

var statusMessages = map[Code]string{
	Unknown:         "Unknown",
	JSONEncodeError: "JSON Encoding Error",
}

func (c Code) Int() int {
	return int(c)
}

func (c Code) String() string {
	mes, ok := statusMessages[c]
	if !ok {
		return statusMessages[Unknown]
	}

	return mes
}
