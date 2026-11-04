package redirect

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	if code == "" {
		h.serveErrorPage(w, "not_found", http.StatusNotFound)
		return
	}

	source := r.URL.Query().Get("src")
	userAgent := r.UserAgent()
	referrer := r.Referer()

	result, err := h.service.Resolve(r.Context(), code, userAgent, referrer, source)
	if err != nil {
		slog.Error("redirect resolve error", "error", err, "code", code)
		h.serveErrorPage(w, "not_found", http.StatusInternalServerError)
		return
	}

	switch result.Action {
	case "redirect":
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		http.Redirect(w, r, result.URL, http.StatusFound)

	case "not_found":
		h.serveErrorPage(w, "not_found", http.StatusNotFound)

	case "inactive":
		h.serveErrorPage(w, "inactive", http.StatusGone)

	case "no_url":
		h.serveErrorPage(w, "no_url", http.StatusServiceUnavailable)

	default:
		h.serveErrorPage(w, "not_found", http.StatusInternalServerError)
	}
}

func (h *Handler) serveErrorPage(w http.ResponseWriter, page string, status int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	var html string
	switch page {
	case "inactive":
		html = inactiveHTML
	case "no_url":
		html = noURLHTML
	default:
		html = notFoundHTML
	}

	w.Write([]byte(html))
}

const notFoundHTML = `<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Tidak Ditemukan - REPUTAP</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; min-height: 100vh; display: flex; align-items: center; justify-content: center; background: #f8fafc; color: #1e293b; }
        .container { text-align: center; padding: 2rem; max-width: 400px; }
        .icon { font-size: 4rem; margin-bottom: 1rem; }
        h1 { font-size: 1.5rem; margin-bottom: 0.5rem; }
        p { color: #64748b; line-height: 1.6; }
        .brand { margin-top: 2rem; font-size: 0.875rem; color: #94a3b8; }
    </style>
</head>
<body>
    <div class="container">
        <div class="icon">🔍</div>
        <h1>QR Code Tidak Ditemukan</h1>
        <p>QR code yang kamu scan tidak terdaftar di sistem kami. Pastikan QR code masih valid.</p>
        <p class="brand">Powered by REPUTAP</p>
    </div>
</body>
</html>`

const inactiveHTML = `<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Tidak Aktif - REPUTAP</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; min-height: 100vh; display: flex; align-items: center; justify-content: center; background: #f8fafc; color: #1e293b; }
        .container { text-align: center; padding: 2rem; max-width: 400px; }
        .icon { font-size: 4rem; margin-bottom: 1rem; }
        h1 { font-size: 1.5rem; margin-bottom: 0.5rem; }
        p { color: #64748b; line-height: 1.6; }
        .brand { margin-top: 2rem; font-size: 0.875rem; color: #94a3b8; }
    </style>
</head>
<body>
    <div class="container">
        <div class="icon">⏸️</div>
        <h1>QR Code Tidak Aktif</h1>
        <p>QR code ini sedang tidak aktif. Hubungi pemilik toko untuk informasi lebih lanjut.</p>
        <p class="brand">Powered by REPUTAP</p>
    </div>
</body>
</html>`

const noURLHTML = `<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Belum Dikonfigurasi - REPUTAP</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; min-height: 100vh; display: flex; align-items: center; justify-content: center; background: #f8fafc; color: #1e293b; }
        .container { text-align: center; padding: 2rem; max-width: 400px; }
        .icon { font-size: 4rem; margin-bottom: 1rem; }
        h1 { font-size: 1.5rem; margin-bottom: 0.5rem; }
        p { color: #64748b; line-height: 1.6; }
        .brand { margin-top: 2rem; font-size: 0.875rem; color: #94a3b8; }
    </style>
</head>
<body>
    <div class="container">
        <div class="icon">⚠️</div>
        <h1>Belum Dikonfigurasi</h1>
        <p>Toko ini belum mengatur link review. Hubungi pemilik toko untuk mengaktifkan fitur ini.</p>
        <p class="brand">Powered by REPUTAP</p>
    </div>
</body>
</html>`
