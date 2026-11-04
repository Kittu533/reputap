package redirect

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler struct {
	service     *Service
	frontendURL string
}

func NewHandler(service *Service, frontendURL string) *Handler {
	return &Handler{service: service, frontendURL: frontendURL}
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

	case "unassigned":
		if h.frontendURL != "" {
			w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Expires", "0")
			http.Redirect(w, r, fmt.Sprintf("%s/activate?code=%s", h.frontendURL, code), http.StatusFound)
			return
		}
		h.serveErrorPage(w, "unassigned", http.StatusOK)

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
	case "unassigned":
		html = unassignedHTML
	case "inactive":
		html = inactiveHTML
	case "no_url":
		html = noURLHTML
	default:
		html = notFoundHTML
	}

	w.Write([]byte(html))
}

const unassignedHTML = `<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Kartu Belum Diaktivasi - REPUTAP</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; min-height: 100vh; display: flex; align-items: center; justify-content: center; background: #f8fafc; color: #1e293b; padding: 1rem; }
        .container { text-align: center; padding: 2.5rem; max-width: 440px; background: #ffffff; border-radius: 16px; box-shadow: 0 4px 20px rgba(0,0,0,0.06); }
        .icon { font-size: 3.5rem; margin-bottom: 1rem; }
        h1 { font-size: 1.4rem; font-weight: 700; margin-bottom: 0.5rem; color: #0f172a; }
        p { color: #64748b; line-height: 1.5; font-size: 0.95rem; margin-bottom: 1.5rem; }
        .btn { display: inline-block; background: #2563eb; color: #ffffff; padding: 0.75rem 1.5rem; border-radius: 8px; font-weight: 600; text-decoration: none; transition: background 0.2s; }
        .btn:hover { background: #1d4ed8; }
        .brand { margin-top: 1.75rem; font-size: 0.8rem; color: #94a3b8; }
    </style>
</head>
<body>
    <div class="container">
        <div class="icon">✨</div>
        <h1>Kartu Belum Diaktivasi</h1>
        <p>Kartu Smart Review REPUTAP ini baru dan siap digunakan. Silakan login ke dashboard toko Anda untuk menghubungkan kartu ini ke Google Review.</p>
        <a href="http://localhost:3000/qr-codes" class="btn">Aktivasi Toko Sekarang</a>
        <p class="brand">Powered by REPUTAP</p>
    </div>
</body>
</html>`

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
