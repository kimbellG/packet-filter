package codes

type Code int

const (
	Unknown Code = iota
)

var statusMessages = map[Code]string{
	Unknown: "Unknown error",
}

func (c Code) Int() int {
	return int(c)
}

func (c Code) Message() string {
	mes, ok := statusMessages[c]
	if !ok {
		return statusMessages[Unknown]
	}

	return mes
}
