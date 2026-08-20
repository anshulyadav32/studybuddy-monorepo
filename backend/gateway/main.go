package main

import (
    "bytes"
    "crypto/rand"
    "encoding/base64"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "os"
    "strings"
    "time"
)

const defaultCoreURL = "http://localhost:8081"

type apiResponse struct {
    Message string `json:"message,omitempty"`
    Error   string `json:"error,omitempty"`
}

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    mux := http.NewServeMux()
    mux.HandleFunc("/", handleRoot)
    mux.HandleFunc("/health", handleHealth)
    mux.HandleFunc("/api/v2/auth/register", proxyAuth("/auth/register"))
    mux.HandleFunc("/api/v2/auth/login", proxyAuth("/auth/login"))
    mux.HandleFunc("/api/v2/auth/validate", proxyAuth("/auth/validate"))
    mux.HandleFunc("/api/v2/auth/me", proxyAuth("/auth/me"))
    mux.HandleFunc("/api/v2/users/me", proxyAuth("/users/me"))
    mux.HandleFunc("/api/v2/users/profile", proxyAuth("/users/profile"))

    // Google OAuth endpoints
    mux.HandleFunc("/api/v2/auth/google/login", googleLoginHandler)
    mux.HandleFunc("/api/v2/auth/google/callback", googleCallbackHandler)

    mux.HandleFunc("/api/v2/courses", proxyPath("/courses"))
    mux.HandleFunc("/api/v2/courses/", proxyPath("/courses/"))

    fmt.Printf("gateway listening on :%s\n", port)
    if err := http.ListenAndServe(":"+port, mux); err != nil {
        panic(err)
    }
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
    if r.URL.Path != "/" {
        writeJSONError(w, http.StatusNotFound, "endpoint not found")
        return
    }

    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte(`
        <html>
          <head><title>StudyBuddy API</title></head>
          <body style="font-family: sans-serif; padding: 32px; color: #0f172a; background: #f8fafc;">
            <h1>StudyBuddy API Gateway</h1>
            <p>The StudyBuddy platform is running.</p>
            <ul>
              <li><a href="/health">/health</a></li>
              <li><a href="/api/v2/auth/register">/api/v2/auth/register</a></li>
              <li><a href="/api/v2/courses">/api/v2/courses</a></li>
            </ul>
          </body>
        </html>
    `))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    _, _ = w.Write([]byte(`{"status":"ok","service":"gateway"}`))
}

func proxyAuth(path string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        proxyTo(w, r, path)
    }
}

func proxyPath(path string) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        suffix := strings.TrimPrefix(r.URL.Path, "/api/v2")
        if suffix == "" {
            suffix = path
        }
        if path == "/courses/" && suffix == "/courses" {
            suffix = "/courses"
        }
        proxyTo(w, r, suffix)
    }
}

func proxyTo(w http.ResponseWriter, r *http.Request, targetPath string) {
    coreURL := strings.TrimRight(os.Getenv("CORE_URL"), "/")
    if coreURL == "" {
        coreURL = defaultCoreURL
    }

    target := coreURL + targetPath
    if r.URL.RawQuery != "" {
        target += "?" + r.URL.RawQuery
    }
    body, err := io.ReadAll(r.Body)
    if err != nil {
        writeJSONError(w, http.StatusBadRequest, "unable to read request body")
        return
    }
    defer r.Body.Close()

    req, err := http.NewRequest(r.Method, target, bytes.NewReader(body))
    if err != nil {
        writeJSONError(w, http.StatusInternalServerError, "unable to create upstream request")
        return
    }

    for key, values := range r.Header {
        for _, value := range values {
            req.Header.Add(key, value)
        }
    }
    if req.Header.Get("Authorization") == "" {
        if token := tokenFromCookieHeader(r.Header.Get("Cookie")); token != "" {
            req.Header.Set("Authorization", "Bearer "+token)
        }
    }

    res, err := http.DefaultClient.Do(req)
    if err != nil {
        writeJSONError(w, http.StatusBadGateway, "upstream auth service unavailable")
        return
    }
    defer res.Body.Close()

    payload, err := io.ReadAll(res.Body)
    if err != nil {
        writeJSONError(w, http.StatusBadGateway, "unable to read upstream response")
        return
    }

    copyHeaders(w, res.Header)
    w.WriteHeader(res.StatusCode)
    _, _ = w.Write(payload)
}

func tokenFromCookieHeader(cookieHeader string) string {
    for _, part := range strings.Split(cookieHeader, ";") {
        segment := strings.TrimSpace(part)
        if segment == "" {
            continue
        }
        key, value, ok := strings.Cut(segment, "=")
        if ok && strings.EqualFold(key, "token") {
            return strings.TrimSpace(value)
        }
    }
    return ""
}

func copyHeaders(dst http.ResponseWriter, src http.Header) {
    for key, values := range src {
        for _, value := range values {
            dst.Header().Add(key, value)
        }
    }
}

func writeJSONError(w http.ResponseWriter, status int, msg string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(apiResponse{Error: msg})
}

// -------------------- Google OAuth helpers --------------------

func googleLoginHandler(w http.ResponseWriter, r *http.Request) {
    clientID := os.Getenv("GOOGLE_CLIENT_ID")
    if clientID == "" {
        writeJSONError(w, http.StatusInternalServerError, "GOOGLE_CLIENT_ID not configured")
        return
    }

    // build redirect URI
    redirectURI := os.Getenv("GOOGLE_OAUTH_REDIRECT_URI")
    if redirectURI == "" {
        // derive from request
        scheme := "http"
        if r.TLS != nil {
            scheme = "https"
        }
        host := r.Host
        redirectURI = fmt.Sprintf("%s://%s/api/v2/auth/google/callback", scheme, host)
    }

    // generate state
    state, err := generateState(16)
    if err != nil {
        writeJSONError(w, http.StatusInternalServerError, "unable to generate state")
        return
    }

    // set state cookie (short lived)
    cookie := &http.Cookie{
        Name:     "oauth_state",
        Value:    state,
        Path:     "/",
        HttpOnly: true,
        Secure:   false,
        Expires:  time.Now().Add(5 * time.Minute),
    }
    http.SetCookie(w, cookie)

    authURL := url.URL{
        Scheme: "https",
        Host:   "accounts.google.com",
        Path:   "/o/oauth2/v2/auth",
    }

    q := authURL.Query()
    q.Set("client_id", clientID)
    q.Set("response_type", "code")
    q.Set("scope", "openid email profile")
    q.Set("redirect_uri", redirectURI)
    q.Set("state", state)
    q.Set("access_type", "offline")
    q.Set("prompt", "consent")
    authURL.RawQuery = q.Encode()

    http.Redirect(w, r, authURL.String(), http.StatusFound)
}

func googleCallbackHandler(w http.ResponseWriter, r *http.Request) {
    // verify state
    stateCookie, err := r.Cookie("oauth_state")
    if err != nil {
        writeJSONError(w, http.StatusBadRequest, "missing oauth state cookie")
        return
    }
    state := r.URL.Query().Get("state")
    if state == "" || stateCookie.Value != state {
        writeJSONError(w, http.StatusBadRequest, "invalid oauth state")
        return
    }

    code := r.URL.Query().Get("code")
    if code == "" {
        writeJSONError(w, http.StatusBadRequest, "missing code in callback")
        return
    }

    clientID := os.Getenv("GOOGLE_CLIENT_ID")
    clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
    if clientID == "" || clientSecret == "" {
        writeJSONError(w, http.StatusInternalServerError, "Google OAuth client not configured")
        return
    }

    redirectURI := os.Getenv("GOOGLE_OAUTH_REDIRECT_URI")
    if redirectURI == "" {
        scheme := "http"
        if r.TLS != nil {
            scheme = "https"
        }
        host := r.Host
        redirectURI = fmt.Sprintf("%s://%s/api/v2/auth/google/callback", scheme, host)
    }

    tokenResp, err := exchangeCodeForToken(code, clientID, clientSecret, redirectURI)
    if err != nil {
        writeJSONError(w, http.StatusBadGateway, "token exchange failed: "+err.Error())
        return
    }

    userInfo, err := fetchUserInfo(tokenResp.AccessToken)
    if err != nil {
        // still return tokens even if userinfo failed
        userInfo = map[string]interface{}{"error": "failed to fetch userinfo"}
    }

    // return JSON containing tokens and user info
    resp := map[string]interface{}{
        "token": tokenResp,
        "userinfo": userInfo,
    }
    writeJSON(w, http.StatusOK, resp)
}

func generateState(n int) (string, error) {
    b := make([]byte, n)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return base64.RawURLEncoding.EncodeToString(b), nil
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(v)
}

type tokenResponse struct {
    AccessToken  string `json:"access_token"`
    ExpiresIn    int    `json:"expires_in"`
    RefreshToken string `json:"refresh_token,omitempty"`
    Scope        string `json:"scope,omitempty"`
    TokenType    string `json:"token_type,omitempty"`
    IdToken      string `json:"id_token,omitempty"`
}

func exchangeCodeForToken(code, clientID, clientSecret, redirectURI string) (*tokenResponse, error) {
    data := url.Values{}
    data.Set("code", code)
    data.Set("client_id", clientID)
    data.Set("client_secret", clientSecret)
    data.Set("redirect_uri", redirectURI)
    data.Set("grant_type", "authorization_code")

    req, err := http.NewRequest("POST", "https://oauth2.googleapis.com/token", strings.NewReader(data.Encode()))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

    res, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    if res.StatusCode >= 400 {
        b, _ := io.ReadAll(res.Body)
        return nil, fmt.Errorf("token endpoint returned %d: %s", res.StatusCode, string(b))
    }

    var tr tokenResponse
    if err := json.NewDecoder(res.Body).Decode(&tr); err != nil {
        return nil, err
    }
    return &tr, nil
}

func fetchUserInfo(accessToken string) (map[string]interface{}, error) {
    req, err := http.NewRequest("GET", "https://openidconnect.googleapis.com/v1/userinfo", nil)
    if err != nil {
        return nil, err
    }
    req.Header.Set("Authorization", "Bearer "+accessToken)
    res, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()
    if res.StatusCode >= 400 {
        b, _ := io.ReadAll(res.Body)
        return nil, fmt.Errorf("userinfo returned %d: %s", res.StatusCode, string(b))
    }
    var info map[string]interface{}
    if err := json.NewDecoder(res.Body).Decode(&info); err != nil {
        return nil, err
    }
    return info, nil
}
