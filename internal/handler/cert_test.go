package handler

import (
	"net/http"
	"testing"
)

func TestMonitorCertAlertField(t *testing.T) {
	s, h, _ := newTestServer(t)

	// criar com a opção ligada
	rec := do(t, h, "POST", "/api/monitors", `{"name":"API","url":"https://x.test","interval_seconds":60,"timeout_seconds":5,"failure_threshold":3,"cert_alert_enabled":true}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("criar: %d %s", rec.Code, rec.Body)
	}
	m := decode[monitorResponse](t, rec)
	if !m.CertAlertEnabled || m.CertExpiresAt != nil {
		t.Fatalf("criar: %+v, want opção ligada e data null", m)
	}

	// GET devolve os campos; marcas internas não aparecem
	rec = do(t, h, "GET", "/api/monitors/"+m.ID, "")
	body := decode[map[string]any](t, rec)
	if body["cert_alert_enabled"] != true || body["cert_expires_at"] != nil {
		t.Fatalf("get: %v", body)
	}
	if _, ok := body["cert_alerted_30_for"]; ok {
		t.Fatal("marca interna exposta na API")
	}

	// a data gravada pela state machine aparece no GET
	if _, err := s.db.Exec(`UPDATE monitors SET cert_expires_at='2026-12-23T10:00:00Z', cert_alerted_30_for='2026-12-23T10:00:00Z' WHERE id=?`, m.ID); err != nil {
		t.Fatal(err)
	}
	rec = do(t, h, "GET", "/api/monitors/"+m.ID, "")
	if got := decode[monitorResponse](t, rec); got.CertExpiresAt == nil || *got.CertExpiresAt != "2026-12-23T10:00:00Z" {
		t.Fatalf("get com data: %+v", got.CertExpiresAt)
	}

	// editar outro campo não mexe na data nem na opção
	rec = do(t, h, "PATCH", "/api/monitors/"+m.ID, `{"name":"API 2"}`)
	if got := decode[monitorResponse](t, rec); rec.Code != 200 || !got.CertAlertEnabled || got.CertExpiresAt == nil {
		t.Fatalf("patch nome: %d %+v", rec.Code, got)
	}

	// ligar em URL http → 400
	rec = do(t, h, "PATCH", "/api/monitors/"+m.ID, `{"url":"http://x.test","cert_alert_enabled":true}`)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("http com opção: %d, want 400", rec.Code)
	}

	// trocar a URL zera a data e as marcas
	rec = do(t, h, "PATCH", "/api/monitors/"+m.ID, `{"url":"https://y.test"}`)
	if got := decode[monitorResponse](t, rec); rec.Code != 200 || got.CertExpiresAt != nil || !got.CertAlertEnabled {
		t.Fatalf("trocar url: %d %+v", rec.Code, got)
	}
	var a30 *string
	if err := s.db.QueryRow(`SELECT cert_alerted_30_for FROM monitors WHERE id=?`, m.ID).Scan(&a30); err != nil || a30 != nil {
		t.Fatalf("marca após trocar url: %v err=%v, want NULL", a30, err)
	}

	// URL trocada para http:// sem mencionar a opção: desliga sozinha
	rec = do(t, h, "PATCH", "/api/monitors/"+m.ID, `{"url":"http://y.test"}`)
	if got := decode[monitorResponse](t, rec); rec.Code != 200 || got.CertAlertEnabled {
		t.Fatalf("url http: %d %+v, want opção desligada", rec.Code, got)
	}
}

func TestMonitorCertAlertValidation(t *testing.T) {
	_, h, _ := newTestServer(t)
	bad := []string{
		`{"name":"x","url":"http://x.test","interval_seconds":60,"timeout_seconds":5,"failure_threshold":3,"cert_alert_enabled":true}`,
		`{"name":"x","type":"heartbeat","interval_seconds":86400,"cert_alert_enabled":true}`,
		`{"name":"x","type":"tcp","config":{"host":"db","port":443},"interval_seconds":60,"timeout_seconds":5,"failure_threshold":3,"cert_alert_enabled":true}`,
	}
	for _, body := range bad {
		rec := do(t, h, "POST", "/api/monitors", body)
		if rec.Code != http.StatusBadRequest {
			t.Fatalf("%s: %d, want 400", body, rec.Code)
		}
		if e := decode[map[string]any](t, rec); e["code"] != "VALIDATION_ERROR" {
			t.Fatalf("%s: %v, want VALIDATION_ERROR", body, e)
		}
	}

	// heartbeat existente não pode ligar
	rec := do(t, h, "POST", "/api/monitors", `{"name":"x","type":"heartbeat","interval_seconds":86400}`)
	hb := decode[monitorResponse](t, rec)
	if rec = do(t, h, "PATCH", "/api/monitors/"+hb.ID, `{"cert_alert_enabled":true}`); rec.Code != http.StatusBadRequest {
		t.Fatalf("heartbeat patch: %d, want 400", rec.Code)
	}
}
