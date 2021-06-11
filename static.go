package go_utils

import "net/http"

// BaseCSSTemplate contains styles that should be applied to all pages
const baseCSSTemplate = `
/* Styles to be applied to all pages */

`

// CSS serves the server side page base CSS file
func CSS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/css")
		_, _ = w.Write([]byte(baseCSSTemplate))
	}
}

const visitCSSTemplate = `
/* Styles to be applied to visit pages */

`

// VisitCSS serves the visit summary page CSS
func VisitCSS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/css")
		_, _ = w.Write([]byte(visitCSSTemplate))
	}
}

const profileCSSTemplate = `
/* Styles to be applied to profile pages */

`

// ProfileCSS serves the profile page CSS
func ProfileCSS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/css")
		_, _ = w.Write([]byte(profileCSSTemplate))
	}
}

const historyCSSTemplate = `
/* Styles to be applied to history pages */

`

// HistoryCSS serves the history page CSS
func HistoryCSS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/css")
		_, _ = w.Write([]byte(historyCSSTemplate))
	}
}

const invalidCSSTemplate = `
/* Styles to be applied to invalid link pages */

`

// InvalidCSS serves the error page CSS
func InvalidCSS() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/css")
		_, _ = w.Write([]byte(invalidCSSTemplate))
	}
}
