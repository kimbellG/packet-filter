package codes

type Code int

const (
	Unknown Code = iota

	JSONEncodeError

	InvalidProtocolInIP

	BCCGetFromTableError
	BCCSetToTableError
	BCCDeleteFromTableError
	BCCNilValueFromTableError
)

var statusMessages = map[Code]string{
	Unknown: "Unknown",

	JSONEncodeError: "JSON Encoding Error",

	InvalidProtocolInIP: "InvalidProtocolInIP",

	BCCGetFromTableError:      "BCCGetFromTableError",
	BCCSetToTableError:        "BCCSetToTableError",
	BCCDeleteFromTableError:   "BCCDeleteFromTableError",
	BCCNilValueFromTableError: "BCCNilValueFromTableError",
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
