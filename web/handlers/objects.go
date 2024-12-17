// handlers/objects.go

package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

// ObjectAction represents an available action for an object
type ObjectAction struct {
	ID          string   `json:"id"`
	Text        string   `json:"text"`
	Icon        string   `json:"icon,omitempty"`
	Divider     bool     `json:"divider,omitempty"`
	Dangerous   bool     `json:"dangerous,omitempty"`
	Conditions  []string `json:"conditions,omitempty"`
	HxGet       string   `json:"hx-get,omitempty"`
	HxPost      string   `json:"hx-post,omitempty"`
	HxTarget    string   `json:"hx-target,omitempty"`
	HxSwap      string   `json:"hx-swap,omitempty"`
	HxConfirm   string   `json:"hx-confirm,omitempty"`
}

// ObjectContext holds the context for object actions
type ObjectContext struct {
	ObjectType  string                 `json:"type"`
	ObjectID    string                 `json:"id"`
	Properties  map[string]interface{} `json:"properties"`
	Actions     []ObjectAction         `json:"actions"`
	Selected    []string              `json:"selected,omitempty"`
}

// HandleObjectMenu handles rendering the context menu for objects
func HandleObjectMenu(w http.ResponseWriter, r *http.Request) {
	// Parse object context from request
	objType := r.URL.Query().Get("type")
	objID := r.URL.Query().Get("id")
	selected := r.URL.Query().Get("selected")

	// Get object properties
	properties, err := getObjectProperties(objType, objID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build available actions
	actions := getAvailableActions(objType, properties)

	// If multiple objects selected, filter to bulk actions
	selectedIDs := strings.Split(selected, ",")
	if len(selectedIDs) > 1 {
		actions = filterBulkActions(actions)
	}

	context := ObjectContext{
		ObjectType: objType,
		ObjectID:   objID,
		Properties: properties,
		Actions:    actions,
		Selected:   selectedIDs,
	}

	// Render menu template
	if err := templates.ExecuteTemplate(w, "objects/menu.html", context); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// HandleObjectAction handles executing an action on objects
func HandleObjectAction(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")
	objType := r.URL.Query().Get("type")
	objIDs := strings.Split(r.URL.Query().Get("ids"), ",")

	// Validate action is allowed
	if !isActionAllowed(action, objType) {
		http.Error(w, "Action not allowed", http.StatusForbidden)
		return
	}

	// Execute action
	result, err := executeAction(action, objType, objIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return updated content
	if err := templates.ExecuteTemplate(w, "objects/result.html", result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// HandleObjectModal handles rendering modals for object actions
func HandleObjectModal(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")
	objType := r.URL.Query().Get("type")
	objID := r.URL.Query().Get("id")

	data := struct {
		Action     string
		ObjectType string
		ObjectID   string
		Properties map[string]interface{}
	}{
		Action:     action,
		ObjectType: objType,
		ObjectID:   objID,
	}

	// Get object properties if needed
	if objID != "" {
		props, err := getObjectProperties(objType, objID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data.Properties = props
	}

	// Render appropriate modal template
	template := fmt.Sprintf("objects/modals/%s_%s.html", objType, action)
	if err := templates.ExecuteTemplate(w, template, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Helper functions below would integrate with your actual object/data layer
func getObjectProperties(objType, objID string) (map[string]interface{}, error) {
	// Implementation depends on your data layer
	return nil, nil
}

func getAvailableActions(objType string, properties map[string]interface{}) []ObjectAction {
	// Implementation depends on your business logic
	return nil
}

func filterBulkActions(actions []ObjectAction) []ObjectAction {
	var bulkActions []ObjectAction
	for _, action := range actions {
		// Add logic to determine if action supports bulk operations
		bulkActions = append(bulkActions, action)
	}
	return bulkActions
}

func isActionAllowed(action, objType string) bool {
	// Implementation depends on your permissions system
	return true
}

func executeAction(action, objType string, objIDs []string) (interface{}, error) {
	// Implementation depends on your business logic
	return nil, nil
}