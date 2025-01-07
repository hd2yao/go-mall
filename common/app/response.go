package app

import "github.com/gin-gonic/gin"

type response struct {
    ctx        *gin.Context
    Code       int         `json:"code"`
    Msg        string      `json:"msg"`
    RequestId  string      `json:"request_id"`
    Data       interface{} `json:"data,omitempty"`
    Pagination *Pagination `json:"pagination,omitempty"`
}
