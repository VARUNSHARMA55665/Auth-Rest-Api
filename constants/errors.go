package constants

const (
	InternalServerError = "P11001"
	InvalidRequest      = "P11002"
	InvalidToken        = "P11003"
)

var ErrorCodeMap = map[string]string{
	"P11001": "INTERNAL SERVER ERROR",
	"P11002": "Invalid Request",
	"P11003": "Invalid Token",
}
