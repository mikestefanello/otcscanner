package handlers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mikestefanello/otcscanner/models"
	"github.com/mikestefanello/otcscanner/repository"
	"github.com/rs/zerolog/log"
)

const cookieNamePreviousScan = "previous_scan"

// ScanForm handles both get and post requests on the scan form route
func (h *HTTPHandler) ScanForm(w http.ResponseWriter, r *http.Request) {
	page := Page{
		Title:   "Scan",
		Content: models.Scan{},
	}

	if r.Method == http.MethodPost {
		// Process the scan
		scan, err := h.processScan(r)
		if err != nil {
			if _, ok := err.(validator.ValidationErrors); ok {
				for _, valErr := range err.(validator.ValidationErrors) {
					page.AddMessage("danger", fmt.Sprintf("%s failed validation: %s", valErr.Field(), valErr.Tag()))
				}
			} else {
				page.AddMessage("danger", err.Error())
			}

		} else {
			page.AddMessage("success", "Scan processed successfully.")
		}

		// Set the scan in a cookie so the values default the form
		h.setPreviousScanCookie(w, scan)

		// Pass the scan to the page
		page.Content = scan
	} else {
		scan, err := h.getPreviousScanFromCookie(r)
		if err == nil {
			page.Content = scan
		}
	}

	h.Render(w, "scan", page)
}

// setPreviousScanCookie encodes a scan in to a cookie and sets it in the response
func (h *HTTPHandler) setPreviousScanCookie(w http.ResponseWriter, scan models.Scan) error {
	json, err := json.Marshal(scan)

	if err != nil {
		log.Error().Err(err).Msg("Unable to encode scan as JSON for cookie.")
		return err
	}

	encoded := base64.StdEncoding.EncodeToString(json)
	c := http.Cookie{
		Name:  cookieNamePreviousScan,
		Value: encoded,
	}
	http.SetCookie(w, &c)

	return nil
}

// getPreviousScanFromCookie loads the cookie that stores a scan and decodes it
func (h *HTTPHandler) getPreviousScanFromCookie(r *http.Request) (models.Scan, error) {
	scan := models.Scan{}

	cookie, err := r.Cookie(cookieNamePreviousScan)
	if err != nil {
		return scan, err
	}

	decoded, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		log.Error().Err(err).Msg("Unable to decode base64 scan cookie.")
		return scan, err
	}

	err = json.Unmarshal(decoded, &scan)
	if err != nil {
		log.Error().Err(err).Msg("Unable to decode json scan cookie.")
		return scan, err
	}

	return scan, nil
}

// processScan processes scan input and attempts to update a matching order in the database
func (h *HTTPHandler) processScan(r *http.Request) (models.Scan, error) {
	// Build a scan model from the form values
	var s = models.Scan{
		Barcode: strings.ToUpper(r.FormValue("barcode")),
		Country: r.FormValue("country"),
		Weight:  r.FormValue("weight"),
		Length:  r.FormValue("length"),
		Width:   r.FormValue("width"),
		Height:  r.FormValue("height"),
		Service: r.FormValue("service"),
		Account: r.FormValue("account"),
	}

	// Validate the input
	err := h.validator.Struct(s)
	if err != nil {
		return s, err
	}

	// Load an order with the given barcode
	order, err := h.repo.LoadByID(s.Barcode)
	if err != nil {
		if err == repository.ErrNotFound {
			return s, errors.New("Unable to match barcode to order")
		} else {
			log.Error().Err(err).Msg("Unable to load order from database.")
			return s, errors.New("Unable to communicate with database")
		}
	}

	// Update the order with the scan
	order.Country = s.Country
	order.Weight = s.Weight
	order.Length = s.Length
	order.Width = s.Width
	order.Height = s.Height
	order.Service = s.Service
	order.Account = s.Account
	order.CalculateDim()

	// Save the order
	err = h.repo.UpdateOne(order)
	if err != nil {
		log.Error().Err(err).Msg("Unable to update order in database from scan.")
		return s, errors.New("Unable to save order in the database")
	}

	return s, nil
}
