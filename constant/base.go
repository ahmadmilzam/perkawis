package constant

const (
	CodeSuccess = iota
	CodeErrRequestNotValid
	CodeErrQueryDB
	CodeErrDataNotFound
	CodeInternalServerError
	CodeCreated
	CodeBadRequest
)

var CodeMapping = map[int]string{
	CodeSuccess:             "Success",
	CodeErrRequestNotValid:  "Request not valid",
	CodeErrQueryDB:          "There is error while query DB",
	CodeErrDataNotFound:     "No Found",
	CodeInternalServerError: "Internal server error",
	CodeCreated:             "Success created",
	CodeBadRequest:          "Bad Request",
}
