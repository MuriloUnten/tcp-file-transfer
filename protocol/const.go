package protocol

type Method int
type StatusCode int

const (
	Chat Method = iota
	Fetch
	Quit
)

func TranslateMethod(m Method) string {
 	switch m {
	case Chat:
		return "Chat"
	case Fetch:
		return "Fetch"
	case Quit:
		return "Quit"
	default:
		return "Unknown Method"
	}
}

const (
	Ok StatusCode = iota
	BadRequest
	NotFound
	InternalError
)

func TranslateStatusCode(s StatusCode) string {
	switch s {
	case Ok:
		return "Ok"
	case BadRequest:
		return "Bad Request"
	case NotFound:
		return "Not Found"
	case InternalError:
		return "Internal Error"
	default:
		return "Unknown Status Code"
	}
}
