package handler

import (
	"database/sql"
	"io"
	"net/http"
	"os"
	"time"
)

type ExportHandler struct {
	db *sql.DB
}

func (s *Server) handleExportDB(w http.ResponseWriter, r *http.Request) {
	dbPath := os.Getenv("PINGOU_DATABASE_URL")
	if dbPath == "" {
		dbPath = "pingou.db"
	}

	f, err := os.Open(dbPath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not open database file", "INTERNAL_ERROR")
		return
	}
	defer f.Close()

	filename := "pingou-backup-" + time.Now().UTC().Format("20060102-150405") + ".db"
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)

	io.Copy(w, f)
}
