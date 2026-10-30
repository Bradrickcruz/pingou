package handler

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

type infoResponse struct {
	Version           string `json:"version"`
	Commit            string `json:"commit"`
	BuildDate         string `json:"build_date"`
	DatabaseSizeBytes int64  `json:"database_size_bytes"`
}

// handleInfo não passa por service/repository: não há SQL, só os.Stat.
func (s *Server) handleInfo(w http.ResponseWriter, r *http.Request) {
	size, err := databaseSize(s.cfg.DatabaseURL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not read database size", "INTERNAL_ERROR")
		return
	}
	writeJSON(w, http.StatusOK, infoResponse{
		Version:           s.cfg.Version,
		Commit:            s.cfg.Commit,
		BuildDate:         s.cfg.BuildDate,
		DatabaseSizeBytes: size,
	})
}

// databaseSize soma o arquivo do banco com o -wal e o -shm, quando existem.
// Aceita DSN com prefixo "file:" e com "?query".
func databaseSize(dsn string) (int64, error) {
	path, _, _ := strings.Cut(dsn, "?")
	path = strings.TrimPrefix(path, "file:")

	info, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("stat database: %w", err)
	}
	total := info.Size()

	for _, suffix := range []string{"-wal", "-shm"} {
		info, err := os.Stat(path + suffix)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return 0, fmt.Errorf("stat database%s: %w", suffix, err)
		}
		total += info.Size()
	}
	return total, nil
}
