package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/zhashkevych/courses-backend/internal/service"
	"github.com/zhashkevych/courses-backend/pkg/payment/fondy"
	"net/http"
)

func (h *Handler) initCallbackRoutes(api *gin.RouterGroup) {
	callback := api.Group("/callback")
	{
		callback.POST("/fondy", h.handleFondyCallback)
	}
}

func (h *Handler) handleFondyCallback(c *gin.Context) {
	if c.Request.UserAgent() != fondy.FondyUserAgent {
		newResponse(c, http.StatusForbidden, "forbidden")
		return
	}

	var inp fondy.Callback
	if err := c.BindJSON(&inp); err != nil {
		newResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.Payments.ProcessTransaction(c.Request.Context(), inp); err != nil {
		if err == service.ErrTransactionInvalid {
			newResponse(c, http.StatusBadRequest, err.Error())
			return
		}

		newResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
