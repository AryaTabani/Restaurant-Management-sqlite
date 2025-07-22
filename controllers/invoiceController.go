package controllers

import (
	"errors"
	"net/http"

	"example.com/m/v2/services"
	"github.com/gin-gonic/gin"
)

func GetInvoicesHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		invoices, err := services.GetAllInvoices()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not list invoices"})
			return
		}
		c.JSON(http.StatusOK, invoices)
	}
}

func GetInvoiceHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		invoiceId := c.Param("invoice_id")

		invoice, err := services.GetInvoice(invoiceId)
		if err != nil {
			if errors.Is(err, services.ErrInvoiceNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch invoice details"})
			}
			return
		}

		c.JSON(http.StatusOK, invoice)
	}
}
