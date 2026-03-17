package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ericolvr/sec-back-v2/internal/core/domain"
	"github.com/ericolvr/sec-back-v2/internal/infrastructure/storage"
	"github.com/gin-gonic/gin"
)

type ActivityMediaHandler struct {
	mediaRepo     domain.ActivityMediaRepository
	storageClient *storage.StorageClient
}

func NewActivityMediaHandler(mediaRepo domain.ActivityMediaRepository, storageClient *storage.StorageClient) *ActivityMediaHandler {
	return &ActivityMediaHandler{
		mediaRepo:     mediaRepo,
		storageClient: storageClient,
	}
}

// Upload faz upload de um arquivo e cria registro no banco
func (h *ActivityMediaHandler) Upload(c *gin.Context) {
	activityID, err := strconv.ParseInt(c.Param("activityId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	// Verificar se storage client está disponível
	if h.storageClient == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Storage service not available. Please configure GCS credentials."})
		return
	}

	// Receber arquivo do form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required", "details": err.Error()})
		return
	}
	defer file.Close()

	// Validar tamanho (max 10MB)
	if header.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size exceeds 10MB limit"})
		return
	}

	// Determinar tipo de mídia baseado na extensão
	ext := strings.ToLower(filepath.Ext(header.Filename))
	mediaType := getMediaType(ext)
	if mediaType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file type"})
		return
	}

	// Ler conteúdo do arquivo
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	// Gerar nome único para o arquivo
	timestamp := time.Now().Unix()
	objectName := fmt.Sprintf("activities/%d/%d_%s", activityID, timestamp, header.Filename)

	// Upload para Google Cloud Storage
	ctx := context.Background()
	fmt.Printf("🔄 [UPLOAD] Uploading file: %s (size: %d bytes)\n", objectName, len(fileBytes))
	mediaURL, err := h.storageClient.UploadFile(ctx, objectName, fileBytes, header.Header.Get("Content-Type"))
	if err != nil {
		fmt.Printf("❌ [UPLOAD ERROR] Failed to upload: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file", "details": err.Error()})
		return
	}
	fmt.Printf("✅ [UPLOAD] File uploaded successfully: %s\n", mediaURL)

	// Criar registro no banco
	media := &domain.ActivityMedia{
		ActivityID: activityID,
		MediaURL:   mediaURL,
		MediaType:  mediaType,
	}

	if err := h.mediaRepo.Create(c.Request.Context(), media); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save media record", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, media)
}

// List retorna todas as mídias de uma atividade
func (h *ActivityMediaHandler) List(c *gin.Context) {
	activityID, err := strconv.ParseInt(c.Param("activityId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}

	mediaList, err := h.mediaRepo.GetByActivityID(c.Request.Context(), activityID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list media", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mediaList)
}

// Delete remove uma mídia
func (h *ActivityMediaHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid media ID"})
		return
	}

	// Buscar mídia para verificar se existe
	_, err = h.mediaRepo.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Media not found"})
		return
	}

	// Deletar do banco
	if err := h.mediaRepo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete media", "details": err.Error()})
		return
	}

	// TODO: Deletar do GCS também (extrair objectName da URL)
	// h.storageClient.DeleteFile(ctx, objectName)

	c.JSON(http.StatusOK, gin.H{"message": "Media deleted successfully"})
}

// getMediaType determina o tipo de mídia baseado na extensão
func getMediaType(ext string) string {
	imageExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true,
	}
	videoExts := map[string]bool{
		".mp4": true, ".mov": true, ".avi": true, ".webm": true,
	}
	docExts := map[string]bool{
		".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true, ".txt": true,
	}

	if imageExts[ext] {
		return "photo"
	}
	if videoExts[ext] {
		return "video"
	}
	if docExts[ext] {
		return "document"
	}
	return ""
}
