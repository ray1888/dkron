package dkron

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/distribworks/dkron/v3/dkron/model"
	dkronpb "github.com/distribworks/dkron/v3/plugin/types"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (h *HTTPTransport) loginHandler(c *gin.Context) {
	var loginRequest model.User
	var resp Response
	if err := c.BindJSON(&loginRequest); err != nil {
		_, _ = c.Writer.WriteString(fmt.Sprintf("Unable to parse payload: %s.", err))
		h.logger.Error(err)
		resp.Code = 1001
		resp.Msg = "Unable to login, invalid request body"
		renderJSON(c, http.StatusForbidden, resp)
		return
	}

	password, err := hashPassword(loginRequest.Password)
	if err != nil {
		h.logger.Error(err)
		resp.Code = 1002
		resp.Msg = fmt.Sprintf("encrypt password failed, err is %v", err)
		renderJSON(c, http.StatusForbidden, resp)
		return
	}

	// Call gRPC Login
	if err := h.agent.GRPCClient.Login(loginRequest.Username, password); err != nil {
		_ = c.AbortWithError(http.StatusTemporaryRedirect, err)
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": loginRequest.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 设置 token 过期时间
	})

	tokenString, err := token.SignedString(password)
	if err != nil {
		return
	}

	// Call Set Token to storage
	if err := h.agent.GRPCClient.SetSession(loginRequest.Username, tokenString, dkronpb.UserAction_Add); err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	// TODO set cookie maybe need to debug
	c.SetCookie("dkron-user-token", tokenString, 3600, "/", "localhost", false, true)
	renderJSON(c, http.StatusOK, resp)
}

func (h *HTTPTransport) logoutHandler(c *gin.Context) {
	var logoutRequest model.User
	if err := c.BindJSON(&logoutRequest); err != nil {
		_, _ = c.Writer.WriteString(fmt.Sprintf("Unable to parse payload: %s.", err))
		h.logger.Error(err)
		return
	}

	// Call gRPC SetJob
	if err := h.agent.GRPCClient.Logout(logoutRequest.Username); err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	var resp Response
	renderJSON(c, http.StatusOK, resp)
}

// func (h *HTTPTransport) initHandler(c *gin.Context) {
// 	jobName := c.Param("job")

// 	job, err := h.agent.Store.GetJob(jobName, nil)
// 	if err != nil {
// 		_ = c.AbortWithError(http.StatusNotFound, err)
// 		return
// 	}

// 	// Toggle job status
// 	job.Disabled = !job.Disabled

// 	// Call gRPC UserModify
// 	if err := h.agent.GRPCClient.UserModify(job); err != nil {
// 		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
// 		return
// 	}

// 	c.Header("Location", c.Request.RequestURI)
// 	renderJSON(c, http.StatusOK, job)
// }
