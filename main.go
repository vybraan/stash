package main

import (
	"database/sql"
	"fmt"
	"github/com/vybraan/snare"
	"github/com/vybraan/snare/database"
	"github/com/vybraan/snare/middlewares"
	"github/com/vybraan/stash/components/toast"
	"github/com/vybraan/stash/helpers"
	"github/com/vybraan/stash/ui"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./app.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := database.Migrate(db); err != nil {
		log.Fatal(err)
	}
	if err := database.Seed(db); err != nil {
		log.Fatal(err)
	}

	auth := snare.New(db)

	r := gin.Default()
	r.Use(gin.Logger())

	auth.RegisterRoutes(r)
	RegisterRoutes(r)

	log.Println("stash running on port 8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}

func RegisterRoutes(r *gin.Engine) *gin.Engine {

	r.Static("/assets", "./assets")

	r.DELETE("/delete/:filename", deleteFile)

	r.GET("/", middlewares.AuthRequired, func(c *gin.Context) {
		template := ui.FileManager()

		helpers.Render(c, http.StatusOK, template)
	})

	r.POST("/upload", uploadFile)
	r.GET("/download/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		c.File("/home/public/" + filename)
	})
	r.GET("/files", middlewares.AuthRequired, renderFileList)

	r.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "404: Not Found")
	})

	return r
}

func uploadFile(c *gin.Context) {
	r := c.Request
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {

		c.String(http.StatusBadRequest, `<div class="alert alert-error">Failed to upload file</div>`)
		return
	}
	defer file.Close()

	dir := "/home/public"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			c.String(http.StatusInternalServerError, `<div class="alert alert-error">Failed to create upload directory</div>`)
			return
		}
	}

	dst, err := os.Create(fmt.Sprintf("%s/%s", dir, handler.Filename))
	defer dst.Close()
	if err != nil {
		c.String(http.StatusInternalServerError, `<div class="alert alert-error">Failed to save file</div>`)
		return
	}

	if _, err := io.Copy(dst, file); err != nil {
		c.String(http.StatusInternalServerError, `<div class="alert alert-error">Error saving file</div>`)
		return
	}

	template := toast.Toast(toast.Props{
		Title:       "success",
		Description: "File upload successfully",
	})

	helpers.Render(c, http.StatusOK, template)

	// Return success toast and trigger refresh
	c.String(http.StatusOK, `<script>
			document.getElementById('progressBar').style.width = '0%';
			htmx.ajax('GET', '/files', {target: '#files'});
			document.getElementById('fileInput').value = '';
		</script>`)
}

func renderFileList(c *gin.Context) {
	dir := "/home/public"

	currentFolder := ""
	q := c.Request.URL.Query().Get("folder")

	currentFolder += q

	files, err := os.ReadDir(dir + currentFolder)
	if err != nil || len(files) == 0 {
		c.String(http.StatusOK, "<p>No files uploaded.</p>")
		return
	}
	template := ui.FileList(currentFolder, files)

	helpers.Render(c, http.StatusOK, template)

}

func deleteFile(c *gin.Context) {
	filename := c.Param("filename")
	filePath := "/home/public/" + filename

	if err := os.RemoveAll(filePath); err != nil {
		c.String(http.StatusInternalServerError, "Error deleting file")
		return
	}

	// Return updated file list after delete
	renderFileList(c)
}
