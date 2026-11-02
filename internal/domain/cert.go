package domain

import "time"

// Limites do aviso de expiração do certificado SSL, em dias.
const (
	CertAlertFirstDays = 30
	CertAlertFinalDays = 7
)

// CertAlert é um aviso de certificado perto de expirar, já decidido pela state
// machine e entregue pelos notificadores depois do commit.
type CertAlert struct {
	ExpiresAt     time.Time // NotAfter do certificado, UTC
	DaysLeft      int       // dias inteiros restantes (para baixo); negativo se já venceu
	ThresholdDays int       // CertAlertFirstDays ou CertAlertFinalDays
}
