package handlers

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/mikestefanello/otcscanner/config"
	"github.com/mikestefanello/otcscanner/repository"
)

// Page describes a page that is rendered in templates
type Page struct {
	SiteName string
	Title    string
	Messages Messages
	Content  interface{}
}

// AddMessage adds a status message to a given page
func (p *Page) AddMessage(status, text string) {
	p.Messages = append(p.Messages, Message{
		Status: status,
		Text: text,
	)
}

// Message describes a status message which can be added to pages
type Message struct {
	Status string
	Text   string
}

// Messages is a slice of messages
type Messages []Message

// HTTPHandler handles HTTP routes
type HTTPHandler struct {
	baseTemplate  *template.Template
	pageTemplates map[string]*template.Template
	config        config.Config
	repo          repository.OrderRepository
	validator     *validator.Validate
}

// NewHTTPHandler creates a new HTTP handler
func NewHTTPHandler(cfg config.Config, repo repository.OrderRepository) *HTTPHandler {
	base, err := template.ParseGlob("templates/global/*.html")

	if err != nil {
		panic(err)
	}

	return &HTTPHandler{
		baseTemplate:  base,
		pageTemplates: make(map[string]*template.Template),
		config:        cfg,
		repo:          repo,
		validator:     validator.New(),
	}
}

// Render renders a given page struct within a given template, specified without the .html extension.
// Templates are parsed within all templates in the global template directory
func (h *HTTPHandler) Render(w http.ResponseWriter, tmpl string, page Page) {
	// Format the template path
	path := fmt.Sprintf("templates/%s.html", tmpl)

	// Check if the template has been parsed
	if _, ok := h.pageTemplates[path]; !ok {
		// Parse the template
		parsed, err := template.Must(h.baseTemplate.Clone()).ParseFiles(path)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Store it in the handler so it only happens once
		h.pageTemplates[path] = parsed
	}

	// Set the page site name, if needed
	if page.SiteName == "" {
		page.SiteName = h.config.App.Name
	}

	// Execute the templates
	err := h.pageTemplates[path].Execute(w, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
