package handler

import (
	"database/sql"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ExportHandler struct {
	db *sql.DB
}

func (s *Server) handleExportDB(w http.ResponseWriter, r *http.Request) {
	tmpDir, err := os.MkdirTemp("", "pingou-export-*")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create export temp dir", "INTERNAL_ERROR")
		return
	}
	defer os.RemoveAll(tmpDir)

	exportPath := filepath.Join(tmpDir, "backup.db")
	exportSQL := "VACUUM INTO '" + escapeSQLiteLiteral(exportPath) + "'"
	if _, err := s.db.ExecContext(r.Context(), exportSQL); err != nil {
		writeError(w, http.StatusInternalServerError, "could not export database", "INTERNAL_ERROR")
		return
	}

	f, err := os.Open(exportPath)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not open exported database", "INTERNAL_ERROR")
		return
	}
	defer f.Close()

	filename := "pingou-backup-" + time.Now().UTC().Format("20060102-150405") + ".db"
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)

	if _, err := io.Copy(w, f); err != nil {
		writeError(w, http.StatusInternalServerError, "could not stream exported database", "INTERNAL_ERROR")
		return
	}
}

func escapeSQLiteLiteral(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}
