package i18n

import (
	"net/http"
)

const (
	// LangCookieName est le nom du cookie de langue
	LangCookieName = "wedding_lang"
	// LangContextKey est la clé de contexte pour la langue
	LangContextKey = "lang"
)

// GetLangFromRequest extrait la langue depuis la requête (cookie ou query param)
func GetLangFromRequest(r *http.Request) Lang {
	// 1. Vérifier le query param ?lang=de
	if langParam := r.URL.Query().Get("lang"); langParam != "" {
		if langParam == string(DE) {
			return DE
		}
		return FR
	}

	// 2. Vérifier le cookie
	if cookie, err := r.Cookie(LangCookieName); err == nil {
		if cookie.Value == string(DE) {
			return DE
		}
		return FR
	}

	// 3. Par défaut : français (on ignore Accept-Language pour toujours démarrer en français)
	return FR
}

// SetLangCookie définit le cookie de langue
func SetLangCookie(w http.ResponseWriter, lang Lang) {
	cookie := &http.Cookie{
		Name:     LangCookieName,
		Value:    string(lang),
		Path:     "/",
		MaxAge:   365 * 24 * 60 * 60, // 1 an
		HttpOnly: false,              // Accessible en JS si besoin
		Secure:   false,              // En prod, mettre à true avec HTTPS
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}
