package controllers

import (
	"errors"
	"net/http"
	"strings"

	"example.com/m/v2/models"
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
		invoiceId := c.Param("invoiceid")

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
func CreateInvoiceHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var invoice models.Invoice
		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		createdInvoice, err := services.CreateInvoice(&invoice)
		if err != nil {
			switch {
			case errors.Is(err, services.ErrOrderNotFound):
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			case strings.Contains(err.Error(), "validation failed"):
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create invoice"})
			}
			return
		}
		c.JSON(http.StatusCreated, createdInvoice)
	}
}
func UpdateInvoiceHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		invoiceId := c.Param("invoiceid")
		var invoiceUpdates models.Invoice

		if err := c.BindJSON(&invoiceUpdates); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		updatedInvoice, err := services.UpdateInvoice(invoiceId, &invoiceUpdates)
		if err != nil {
			if errors.Is(err, services.ErrInvoiceNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update invoice"})
			}
			return
		}
		c.JSON(http.StatusOK, updatedInvoice)
	}
}
