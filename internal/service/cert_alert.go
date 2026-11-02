package service

import (
	"math"
	"time"

	"github.com/Bradrickcruz/pingou/internal/domain"
)

// CertAlertDecision é o resultado de EvaluateCertAlert.
type CertAlertDecision struct {
	Alert  *domain.CertAlert // nil = nada a enviar
	Mark30 bool              // gravar cert_alerted_30_for = expiresAt
	Mark7  bool              // gravar cert_alerted_7_for = expiresAt
}

// EvaluateCertAlert decide se um aviso de certificado deve sair agora. Função
// pura: sem banco, sem relógio global, sem envio.
//
// As marcas guardam a data de expiração já alertada (não um sim/não). Um
// certificado renovado tem outra data, as marcas antigas não batem e o ciclo
// recomeça sozinho.
//
//   - faltam 7 dias ou menos (inclui já vencido) e o de 7 ainda não saiu para
//     esta data: aviso de 7 dias, e marca os dois. Quem liga a opção na reta
//     final recebe um aviso só, e o de 30 não chega atrasado depois.
//   - senão, faltam 30 dias ou menos e o de 30 ainda não saiu: aviso de 30.
//   - senão: nada.
func EvaluateCertAlert(now, expiresAt time.Time, alerted30For, alerted7For *time.Time) CertAlertDecision {
	remaining := expiresAt.Sub(now)
	alert := func(threshold int) *domain.CertAlert {
		return &domain.CertAlert{
			ExpiresAt:     expiresAt.UTC(),
			DaysLeft:      int(math.Floor(remaining.Hours() / 24)),
			ThresholdDays: threshold,
		}
	}

	if remaining <= domain.CertAlertFinalDays*24*time.Hour && !sameInstant(alerted7For, expiresAt) {
		return CertAlertDecision{Alert: alert(domain.CertAlertFinalDays), Mark30: true, Mark7: true}
	}
	if remaining <= domain.CertAlertFirstDays*24*time.Hour && !sameInstant(alerted30For, expiresAt) {
		return CertAlertDecision{Alert: alert(domain.CertAlertFirstDays), Mark30: true}
	}
	return CertAlertDecision{}
}

func sameInstant(a *time.Time, b time.Time) bool {
	return a != nil && a.Equal(b)
}
