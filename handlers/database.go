package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/mikestefanello/otcscanner/models"
	"github.com/rs/zerolog/log"
)

type orderStats struct {
	All        int64
	Completed  int64
	Incomplete int64
}

// DatabasePage handles get requests for the database route
func (h *HTTPHandler) DatabasePage(w http.ResponseWriter, r *http.Request) {
	page := Page{
		Title: "Database",
	}

	stats, err := h.getOrderStats()
	if err != nil {
		page.AddMessage("danger", "Unable to communicate with the database.")
	} else {
		page.Content = stats
	}

	h.Render(w, "database", page)
}

// getOrderStats gets order stats from the database
func (h *HTTPHandler) getOrderStats() (orderStats, error) {
	var stats orderStats

	all, err := h.repo.CountAll()
	if err != nil {
		log.Error().Err(err).Msg("Unable to get count of all orders from the database.")
		return stats, err
	}

	completed, err := h.repo.CountCompleted()
	if err != nil {
		log.Error().Err(err).Msg("Unable to get count of all completed orders from the database.")
		return stats, err
	}

	incomplete, err := h.repo.CountIncomplete()
	if err != nil {
		log.Error().Err(err).Msg("Unable to get count of all incomplete orders from the database.")
		return stats, err
	}

	stats.All = all
	stats.Completed = completed
	stats.Incomplete = incomplete

	return stats, nil
}

// DatabaseUpload handles post requests to upload a CSV of orders to the database
func (h *HTTPHandler) DatabaseUpload(w http.ResponseWriter, r *http.Request) {
	page := Page{
		Title: "Database",
	}

	added, err := h.processDatabaseUpload(r)

	if err != nil {
		page.AddMessage("danger", err.Error())
	} else {
		log.Info().Int("count", added).Msg("Uploaded orders to the database.")
		page.AddMessage("success", fmt.Sprintf("Added %d orders to the database.", added))
	}

	h.Render(w, "text", page)
}

// DatabaseDeleteAll handles post requests to delete the entire order database
func (h *HTTPHandler) DatabaseDeleteAll(w http.ResponseWriter, r *http.Request) {
	page := Page{
		Title: "Database",
	}

	err := h.repo.DeleteAll()

	if err != nil {
		log.Error().Err(err).Msg("Unable to delete entire database.")
		page.AddMessage("danger", "Unable to delete entire database.")
	} else {
		log.Info().Msg("Deleted all orders from the database.")
		page.AddMessage("success", "Database deleted.")
	}

	h.Render(w, "text", page)
}

// DatabaseDeleteCompleted handles post requests to delete completed orders from the database
func (h *HTTPHandler) DatabaseDeleteCompleted(w http.ResponseWriter, r *http.Request) {
	page := Page{
		Title: "Database",
	}

	err := h.repo.DeleteCompleted()

	if err != nil {
		log.Error().Err(err).Msg("Unable to delete completed orders from database.")
		page.AddMessage("danger", "Unable to delete completed orders.")
	} else {
		log.Info().Msg("Deleted completed orders from the database.")
		page.AddMessage("success", "Completed orders have been deleted.")
	}

	h.Render(w, "text", page)
}

// DatabaseDownloadAll handles post requests to download the entire database as a CSV file
func (h *HTTPHandler) DatabaseDownloadAll(w http.ResponseWriter, r *http.Request) {
	page := Page{
		Title: "Database",
	}

	// Load all orders
	err := h.serveOrdersCsv(w, r, h.repo.LoadAll)

	if err != nil {
		page.AddMessage("danger", err.Error())
		h.Render(w, "text", page)
		return
	}
}

// DatabaseDownloadCompleted handles post requests to download completed orders from the database as a CSV file
func (h *HTTPHandler) DatabaseDownloadCompleted(w http.ResponseWriter, r *http.Request) {
	page := Page{
		Title: "Database",
	}

	// Load all orders
	err := h.serveOrdersCsv(w, r, h.repo.LoadCompleted)

	if err != nil {
		page.AddMessage("danger", err.Error())
		h.Render(w, "text", page)
		return
	}
}

// DatabaseDownloadIncomplete handles post requests to download incomplete orders from the database as a CSV file
func (h *HTTPHandler) DatabaseDownloadIncomplete(w http.ResponseWriter, r *http.Request) {
	page := Page{
		Title: "Database",
	}

	// Load all orders
	err := h.serveOrdersCsv(w, r, h.repo.LoadIncomplete)

	if err != nil {
		page.AddMessage("danger", err.Error())
		h.Render(w, "text", page)
		return
	}
}

// serveOrdersCsv gets data from a loader function and serves a CSV file with the data returned
func (h *HTTPHandler) serveOrdersCsv(w http.ResponseWriter, r *http.Request, loader func() (*models.Orders, error)) error {
	// Load all orders
	orders, err := loader()

	if err != nil {
		log.Error().Err(err).Msg("Unable to load orders from the database.")
		return errors.New("Unable to load orders")
	}

	csv, err := gocsv.MarshalString(orders)

	if err != nil {
		log.Error().Err(err).Msg("Unable to encode orders as CSV.")
		return errors.New("Unable to process orders for export")
	}

	http.ServeContent(w, r, "db.csv", time.Now(), bytes.NewReader([]byte(csv)))
	return nil
}

// processDatabaseUpload processes CSV uploads and inserts records in to the database
func (h *HTTPHandler) processDatabaseUpload(r *http.Request) (int, error) {
	r.ParseMultipartForm(10 << 20)

	// Get the uploaded file
	file, _, err := r.FormFile("upload")
	if err != nil {
		log.Error().Err(err).Msg("Unable to load database upload file.")
		return 0, errors.New("Error reading the file")
	}
	defer file.Close()

	// Read the entire file
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error().Err(err).Msg("Unable to read database upload file.")
		return 0, errors.New("Error processing the file")
	}

	// Parse in to orders
	orders := &models.Orders{}
	gocsv.UnmarshalString(string(fileBytes), orders)

	// Validate the orders
	for _, order := range *orders {
		err = h.validator.Struct(order)

		if err != nil {
			return 0, err
		}
	}

	// Save the orders
	err = h.repo.InsertMany(orders)

	if err != nil {
		log.Error().Err(err).Msg("Unable to save orders to database.")
		return 0, errors.New("Unable to add items to the database")
	}

	return len(*orders), nil
}
