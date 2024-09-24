package dkron

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/distribworks/dkron/v3/dkron/model"
)

func (h *HTTPTransport) loginHandler(c *gin.Context) {

	var loginRequest model.User
	if err := c.BindJSON(&loginRequest); err != nil {
		_, _ = c.Writer.WriteString(fmt.Sprintf("Unable to parse payload: %s.", err))
		h.logger.Error(err)
		return
	}
	var jobName string
	job, err := h.agent.Store.GetJob(jobName, nil)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	// Toggle job status
	job.Disabled = !job.Disabled

	// Call gRPC SetJob
	if err := h.agent.GRPCClient.SetJob(job); err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	// c.SetCookie("dkron-user-token", "")

	renderJSON(c, http.StatusOK, job)
}

func (h *HTTPTransport) logoutHandler(c *gin.Context) {
	jobName := c.Param("job")

	job, err := h.agent.Store.GetJob(jobName, nil)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	// Toggle job status
	job.Disabled = !job.Disabled

	// Call gRPC SetJob
	if err := h.agent.GRPCClient.SetJob(job); err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	c.Header("Location", c.Request.RequestURI)
	renderJSON(c, http.StatusOK, job)
}

func (h *HTTPTransport) initHandler(c *gin.Context) {
	jobName := c.Param("job")

	job, err := h.agent.Store.GetJob(jobName, nil)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	// Toggle job status
	job.Disabled = !job.Disabled

	// Call gRPC SetJob
	if err := h.agent.GRPCClient.SetJob(job); err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	c.Header("Location", c.Request.RequestURI)
	renderJSON(c, http.StatusOK, job)
}
