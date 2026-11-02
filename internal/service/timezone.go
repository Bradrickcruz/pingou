package service

import (
	"context"
	"log/slog"
	"strings"
	"time"
)

// DefaultTimezone: sem fuso salvo, tudo continua em UTC como antes.
const DefaultTimezone = "UTC"

const maxTimezoneLen = 64

// normalizeTimezone: vazio vira UTC; o resto precisa ser um nome IANA que o Go
// carregue (o binário embute a base com time/tzdata). "Local" é recusado: no
// servidor ele depende da máquina, não do que o usuário escolheu.
func normalizeTimezone(raw string) (string, error) {
	tz := strings.TrimSpace(raw)
	if tz == "" {
		return DefaultTimezone, nil
	}
	if len(tz) > maxTimezoneLen || tz == "Local" {
		return "", errTimezoneInvalid
	}
	if _, err := time.LoadLocation(tz); err != nil {
		return "", errTimezoneInvalid
	}
	return tz, nil
}

// LocationLoader devolve o fuso salvo nas settings. É lido a cada envio, então
// mudar o fuso vale na próxima notificação, sem reiniciar. Erro de leitura ou
// nome que não carrega cai em UTC com aviso no log: a notificação sai mesmo assim.
func LocationLoader(get func(ctx context.Context, key string) (string, error)) func() *time.Location {
	return func() *time.Location {
		name, err := get(context.Background(), settingTimezone)
		if err != nil {
			slog.Warn("timezone: could not read settings, using UTC", "err", err)
			return time.UTC
		}
		if name == "" {
			return time.UTC
		}
		loc, err := time.LoadLocation(name)
		if err != nil {
			slog.Warn("timezone: invalid saved timezone, using UTC", "timezone", name, "err", err)
			return time.UTC
		}
		return loc
	}
}

// messageTimeLayout: em UTC mostra "UTC"; em America/Sao_Paulo mostra "-03".
const messageTimeLayout = "2006-01-02 15:04:05 MST"

func formatMessageTime(t time.Time, loc *time.Location) string {
	if loc == nil {
		loc = time.UTC
	}
	return t.In(loc).Format(messageTimeLayout)
}
