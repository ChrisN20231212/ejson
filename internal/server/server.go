package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf16"

	assetweb "ejson/web"
)

type apiResponse struct {
	OK     bool   `json:"ok"`
	Result string `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

type apiRequest struct {
	Input string `json:"input"`
}

func New() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(handleIndex))
	mux.Handle("/app.css", http.HandlerFunc(handleCSS))
	mux.Handle("/app.js", http.HandlerFunc(handleJS))
	mux.Handle("/api/validate", apiHandler(validateJSON))
	mux.Handle("/api/minify", apiHandler(minifyJSON))
	mux.Handle("/api/escape", apiHandler(escapeText))
	mux.Handle("/api/unescape", apiHandler(unescapeText))
	mux.Handle("/api/unicode/encode", apiHandler(unicodeEncode))
	mux.Handle("/api/unicode/decode", apiHandler(unicodeDecode))
	return mux
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeContent(w, r, "index.html", assetweb.ModTime(), bytes.NewReader(assetweb.IndexHTML()))
}

func handleCSS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	http.ServeContent(w, r, "app.css", assetweb.ModTime(), bytes.NewReader(assetweb.AppCSS()))
}

func handleJS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	http.ServeContent(w, r, "app.js", assetweb.ModTime(), bytes.NewReader(assetweb.AppJS()))
}

func apiHandler(fn func(string) (string, error)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req apiRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, apiResponse{Error: "invalid request body"})
			return
		}

		result, err := fn(req.Input)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, apiResponse{Error: err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, apiResponse{OK: true, Result: result})
	})
}

func writeJSON(w http.ResponseWriter, status int, payload apiResponse) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func validateJSON(input string) (string, error) {
	var out bytes.Buffer
	if err := json.Indent(&out, []byte(input), "", "  "); err != nil {
		return "", err
	}
	return out.String(), nil
}

func minifyJSON(input string) (string, error) {
	var out bytes.Buffer
	if err := json.Compact(&out, []byte(input)); err != nil {
		return "", err
	}
	return out.String(), nil
}

func escapeText(input string) (string, error) {
	quoted := strconv.Quote(input)
	return strings.TrimSuffix(strings.TrimPrefix(quoted, `"`), `"`), nil
}

func unescapeText(input string) (string, error) {
	return decodeEscaped(input)
}

func unicodeEncode(input string) (string, error) {
	var b strings.Builder
	for _, r := range input {
		if r <= 0x7f {
			b.WriteRune(r)
			continue
		}
		if r <= 0xffff {
			b.WriteString(`\u`)
			b.WriteString(padHex(uint16(r)))
			continue
		}
		for _, unit := range utf16.Encode([]rune{r}) {
			b.WriteString(`\u`)
			b.WriteString(padHex(unit))
		}
	}
	return b.String(), nil
}

func unicodeDecode(input string) (string, error) {
	return decodeEscaped(input)
}

func decodeEscaped(input string) (string, error) {
	literal := `"` + input + `"`
	value, err := strconv.Unquote(literal)
	if err != nil {
		return "", err
	}
	return value, nil
}

func padHex(v uint16) string {
	hex := strings.ToLower(strconv.FormatUint(uint64(v), 16))
	for len(hex) < 4 {
		hex = "0" + hex
	}
	return hex
}
