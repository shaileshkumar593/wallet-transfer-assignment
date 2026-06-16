package handler

import (
	"net/http"

	"wallet-transfer-assignment/wallet-transfer/internal/dto"
	"wallet-transfer-assignment/wallet-transfer/internal/service"

	"github.com/gin-gonic/gin"
)

type TransferHandler struct {
	transferService *service.TransferService
}

func NewTransferHandler(
	transferService *service.TransferService,
) *TransferHandler {

	return &TransferHandler{
		transferService: transferService,
	}
}

func (h *TransferHandler) CreateTransfer(
	c *gin.Context,
) {

	var req dto.CreateTransferRequest

	if err := c.ShouldBindJSON(
		&req,
	); err != nil {

		c.JSON(
			http.StatusBadRequest,
			dto.ErrorResponse{
				Error: err.Error(),
			},
		)

		return
	}

	resp, err :=
		h.transferService.Transfer(
			c.Request.Context(),
			req,
		)

	if err != nil {

		c.JSON(
			http.StatusBadRequest,
			dto.ErrorResponse{
				Error: err.Error(),
			},
		)

		return
	}

	c.JSON(
		http.StatusOK,
		resp,
	)
}
