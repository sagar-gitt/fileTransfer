package handlers

import (
	"fileTransfer/internal/dto"
	"fileTransfer/internal/models"
	"fileTransfer/internal/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"net/http"
	"time"
)

func (h *Handlers) UploadFileAndSaveInfo(c *gin.Context) {
	//Logic to upload file on AWS
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	res, err := h.AwsS3.UploadFile(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	//Logic for wrapping file info in the file struct
	f, _ := file.Open()
	defer f.Close()
	size, _ := f.Seek(0, io.SeekEnd)
	f.Seek(0, 0) // reset

	expiry := time.Now().UTC().Add(2 * time.Minute)
	expiryDuration := time.Until(expiry)

	signedURL, err := h.AwsS3.GenerateSignedURL(*res.Key, expiryDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	userId := " ee6d4c16-eaf3-482c-9271-b9236175b57c"
	newFile := models.NewFile(uuid.New().String(), *res.Key, file.Filename, size, expiry, userId, res.Location, time.Now().UTC(), 0)

	//Saving file in the db
	err = h.FileDbRepo.AddFile(newFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Uploaded successfully", "key": res.Key, "URL": signedURL, "valid For": "2 minutes"})
}

func (h *Handlers) DownloadFile(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing file key"})
		return
	}

	// Download file from AWS
	resp, err := h.AwsS3.DownloadFile(key)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Download failed", "details": err.Error()})
		return
	}
	defer resp.Body.Close()

	// Update download count in the db
	err = h.FileDbRepo.IncreaseDownloadCount(key)
	if err != nil {
		// Log the error but continue with the download
		fmt.Printf("Failed to update download count: %v\n", err)
	}

	// Set appropriate headers for file download
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.Name))
	c.Header("Content-Type", *resp.ContentType)
	c.Header("Content-Length", fmt.Sprintf("%d", *resp.ContentLength))

	// Stream the file to the client
	c.DataFromReader(http.StatusOK, *resp.ContentLength, *resp.ContentType, resp.Body, nil)
}

func (h *Handlers) ListFile(c *gin.Context) {
	//Logic to List files from AWS
	resp, err := h.AwsS3.ListFiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not list objects", "details": err.Error()})
		return
	}

	var keys []string
	for _, item := range resp.Contents {
		keys = append(keys, *item.Key)
	}

	c.JSON(http.StatusOK, gin.H{"files": keys})
}

func (h *Handlers) SendFileDownloadLink(c *gin.Context) {
	var body dto.EmailRequestBody
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if body.To == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'to' email"})
		return
	}

	if body.DownloadLink == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing download link"})
		return
	}

	emailBody := "Here is your download link: \n\n" + body.DownloadLink + "\nvalid for: " + body.LinkValidity
	htmlBody, err := utils.RenderEmailHTML(body.DownloadLink, body.LinkValidity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := utils.SendEmailWithSendGrid(body.To, body.Bcc, body.Cc, emailBody, htmlBody); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}
