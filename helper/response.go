package helper

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
	"perkawis/constant"
)

type (
	successJson struct {
		Status  string      `json:"status"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}

	errorJson struct {
		Status  string   `json:"status"`
		Message string   `json:"message"`
		Data    struct{} `json:"data"`
	}

	successDeleteJson struct {
		Status  string   `json:"status"`
		Message string   `json:"message"`
		Data    struct{} `json:"data"`
	}

	ErrorWithCode struct {
		CodeID     int         `json:"-"`
		Msg        string      `json:"message"`
		Status     string      `json:"status"`
		StatusCode int         `json:"-"`
		Data       interface{} `json:"data,omitempty"`
	}
)

func NewErrorMsg(code int, err error) error {
	msg := constant.CodeMapping[code]
	return ErrorWithCode{
		CodeID:     code,
		Msg:        msg,
		Status:     msg,
		StatusCode: http.StatusUnprocessableEntity,
	}
}

func NewErrorRecordNotFound(code int, err error) error {
	msg := constant.CodeMapping[code]
	return ErrorWithCode{
		CodeID:     code,
		Msg:        msg,
		Status:     msg,
		StatusCode: http.StatusNotFound,
		Data:       err.Error(),
	}
}

func (c ErrorWithCode) Error() string {
	b, _ := json.Marshal(&c)

	return string(b)
}

func JsonSUCCESS(c echo.Context, data interface{}) error {
	res := successJson{
		Message: "Success",
		Status:  "Success",
		Data:    data,
	}

	return c.JSON(http.StatusOK, res)
}

func JsonSuccessDelete(c echo.Context) error {
	res := successDeleteJson{
		Message: "Success",
		Status:  "Success",
	}

	return c.JSON(http.StatusOK, res)
}
func JsonCreated(c echo.Context, data interface{}) error {
	res := successJson{
		Message: "Success",
		Status:  "Success",
		Data:    data,
	}
	return c.JSON(http.StatusCreated, res)
}

func JsonValidationError(c echo.Context, message string) error {
	res := errorJson{
		Message: message,
		Status:  "Bad Request",
	}
	return c.JSON(http.StatusBadRequest, res)
}

func JsonNotFound(c echo.Context, message string) error {
	res := errorJson{
		Message: message,
		Status:  "Not Found",
	}

	return c.JSON(http.StatusNotFound, res)
}

func JsonERROR(c echo.Context, err error) error {

	switch err.(type) {
	case ErrorWithCode:
		errMsg := err.(ErrorWithCode)
		res := ErrorWithCode{
			CodeID: http.StatusUnprocessableEntity,
			Msg:    errMsg.Msg,
		}

		c.JSON(errMsg.StatusCode, res)
		return err
	default:
		res := ErrorWithCode{
			CodeID: http.StatusInternalServerError,
			Msg:    "Error message type not defined.",
		}
		c.JSON(res.CodeID, res)
		return err
	}
}
