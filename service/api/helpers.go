package api

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Daniel200273/WASA-project/service/api/reqcontext"
)

// === VALIDATION HELPERS ===

// validateUsername checks if the given username is valid according to API specifications.
// Returns an error if the username is invalid.
// Rules: 3-16 characters long, alphanumeric, underscores, and hyphens allowed.
func validateUsername(username string) error {
	// Length check
	if len(username) < 3 || len(username) > 16 {
		return fmt.Errorf("username must be between 3 and 16 characters")
	}

	// Regular expression check: must match the pattern ^[a-zA-Z0-9_-]+$.
	re := regexp.MustCompile("^[a-zA-Z0-9_-]+$")
	if !re.MatchString(username) {
		return fmt.Errorf("username can only contain letters, numbers, underscores, and hyphens")
	}

	return nil
}

// validateID validates a URL parameter ID (UUID format)
func validateID(id string, fieldName string) error {
	if id == "" {
		return fmt.Errorf("%s is required", fieldName)
	}
	if len(id) < 1 || len(id) > 36 {
		return fmt.Errorf("%s must be between 1 and 36 characters", fieldName)
	}
	// Pattern for UUID-like IDs: alphanumeric, hyphens, underscores
	re := regexp.MustCompile("^[a-zA-Z0-9_-]+$")
	if !re.MatchString(id) {
		return fmt.Errorf("%s contains invalid characters", fieldName)
	}
	return nil
}

// validateGroupName validates group name format
func validateGroupName(name string) error {
	if name == "" {
		return fmt.Errorf("group name is required")
	}
	if len(name) < 1 || len(name) > 50 {
		return fmt.Errorf("group name must be between 1 and 50 characters")
	}
	return nil
}

// validateMessageContent validates message text content
func validateMessageContent(content string) error {
	if content == "" {
		return fmt.Errorf("message content is required")
	}
	if len(content) < 1 || len(content) > 1000 {
		return fmt.Errorf("message content must be between 1 and 1000 characters")
	}
	return nil
}

// validateEmoticon validates emoticon format
func validateEmoticon(emoticon string) error {
	if emoticon == "" {
		return fmt.Errorf("emoticon is required")
	}
	if len(emoticon) < 1 || len(emoticon) > 10 {
		return fmt.Errorf("emoticon must be between 1 and 10 characters")
	}
	return nil
}

// === HTTP HELPERS ===

// parseJSONRequest parses JSON request body into the provided struct
// Use this ONLY for endpoints that expect application/json
func parseJSONRequest(r *http.Request, target interface{}) error {
	// Check Content-Type header for JSON endpoints
	contentType := r.Header.Get("Content-Type")
	if contentType != "" && !strings.HasPrefix(contentType, "application/json") {
		return fmt.Errorf("Content-Type must be application/json, got: %s", contentType)
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() // Strict parsing

	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("invalid JSON format: %w", err)
	}

	return nil
}

// sendJSONResponse sends a JSON response with the specified status code
func sendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		encoder := json.NewEncoder(w)
		return encoder.Encode(data)
	}

	return nil
}

// sendErrorResponse sends a standardized error response
func sendErrorResponse(w http.ResponseWriter, statusCode int, message string, ctx reqcontext.RequestContext) {
	ctx.Logger.WithField("error", message).Error("API error response")

	response := ErrorResponse{Message: message}
	if err := sendJSONResponse(w, statusCode, response); err != nil {
		ctx.Logger.WithError(err).Error("failed to send error response")
	}
}

// getQueryParam extracts a query parameter from URL
func getQueryParam(r *http.Request, paramName string) string {
	return r.URL.Query().Get(paramName)
}

// === FILE UPLOAD HELPERS ===

// validateImageFile validates uploaded image file (for profile/group photos)
func validateImageFile(r *http.Request, fieldName string) error {
	// Parse multipart form (max 10MB)
	err := r.ParseMultipartForm(10 * 1024 * 1024) // 10MB
	if err != nil {
		return fmt.Errorf("failed to parse multipart form: %w", err)
	}

	file, header, err := r.FormFile(fieldName)
	if err != nil {
		return fmt.Errorf("missing or invalid %s file", fieldName)
	}
	defer file.Close()

	// Validate file size
	if header.Size > 10*1024*1024 { // 10MB
		return fmt.Errorf("file size exceeds 10MB limit")
	}

	if header.Size == 0 {
		return fmt.Errorf("file is empty")
	}

	// Validate file type by extension
	filename := strings.ToLower(header.Filename)
	validExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}

	isValid := false
	for _, ext := range validExtensions {
		if strings.HasSuffix(filename, ext) {
			isValid = true
			break
		}
	}

	if !isValid {
		return fmt.Errorf("invalid file type. Allowed: JPG, PNG, GIF, WebP")
	}

	return nil
}

// saveUploadedImage saves an uploaded image file to the temporary uploads directory
// Returns the URL path for accessing the saved image
func saveUploadedImage(file multipart.File, category, filename string) (string, error) {
	// Create the uploads directory structure in tmp
	uploadsDir := filepath.Join("tmp", "uploads", category)
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	// Create the destination file path
	destPath := filepath.Join(uploadsDir, filename)

	// Create the destination file
	destFile, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	// Reset file pointer to beginning
	if _, err := file.Seek(0, 0); err != nil {
		return "", fmt.Errorf("failed to reset file pointer: %w", err)
	}

	// Copy the uploaded file to destination
	if _, err := io.Copy(destFile, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Return the URL path (without tmp prefix for serving)
	urlPath := fmt.Sprintf("/uploads/%s/%s", category, filename)
	return urlPath, nil
}

// initializeUploadsDirectory creates the temporary uploads directory structure
// Call this at application startup to ensure directories exist
func initializeUploadsDirectory() error {
	categories := []string{"profiles", "groups", "messages"}

	for _, category := range categories {
		dir := filepath.Join("tmp", "uploads", category)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create upload directory %s: %w", dir, err)
		}
	}

	return nil
}

// cleanupUploadsDirectory removes the temporary uploads directory
// Call this at application shutdown if needed
func cleanupUploadsDirectory() error {
	uploadsDir := filepath.Join("tmp", "uploads")
	if err := os.RemoveAll(uploadsDir); err != nil {
		return fmt.Errorf("failed to cleanup uploads directory: %w", err)
	}
	return nil
}

// getUploadedFile extracts and validates an uploaded file from the request
// Returns the file, header, and any validation error
func getUploadedFile(r *http.Request, fieldName string) (multipart.File, *multipart.FileHeader, error) {
	// Validate the uploaded file first
	if err := validateImageFile(r, fieldName); err != nil {
		return nil, nil, err
	}

	// Get the file from the form
	file, header, err := r.FormFile(fieldName)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get uploaded file: %w", err)
	}

	return file, header, nil
}
