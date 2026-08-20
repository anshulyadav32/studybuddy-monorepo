package main

import (
    "encoding/json"
    "errors"
    "fmt"
    "log"
    "net/http"
    "os"
    "strings"
    "time"

    jwt "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

type userRecord struct {
    ID           string      `json:"id"`
    Name         string      `json:"name"`
    Email        string      `json:"email"`
    PasswordHash string      `json:"-"`
    Profile      profileData `json:"profile,omitempty"`
}

type profileData struct {
    Phone       string         `json:"phone,omitempty"`
    AvatarURL   string         `json:"avatarUrl,omitempty"`
    ClassLevel  string         `json:"classLevel,omitempty"`
    Board       string         `json:"board,omitempty"`
    Preferences map[string]any `json:"preferences,omitempty"`
}

type courseRecord struct {
    ID         string `json:"id"`
    Title      string `json:"title"`
    Slug       string `json:"slug"`
    SubjectID  string `json:"subjectId"`
    TeacherID  string `json:"teacherId"`
    CreatedAt  string `json:"createdAt"`
}

type courseRequest struct {
    Title      string `json:"title"`
    Slug       string `json:"slug,omitempty"`
    SubjectID  string `json:"subjectId"`
    TeacherID  string `json:"teacherId,omitempty"`
}

type chapterRecord struct {
    ID       string `json:"id"`
    CourseID string `json:"courseId"`
    Title    string `json:"title"`
    Order    int    `json:"order"`
}

type chapterRequest struct {
    Title    string `json:"title"`
    Order    int    `json:"order,omitempty"`
}

type registerRequest struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
}

type loginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type profileRequest struct {
    Phone      string         `json:"phone,omitempty"`
    AvatarURL  string         `json:"avatarUrl,omitempty"`
    ClassLevel string         `json:"classLevel,omitempty"`
    Board      string         `json:"board,omitempty"`
    Preferences map[string]any `json:"preferences,omitempty"`
}

type authResponse struct {
    Token string       `json:"token"`
    User  userResponse `json:"user"`
}

type userResponse struct {
    ID      string      `json:"id"`
    Name    string      `json:"name"`
    Email   string      `json:"email"`
    Profile profileData `json:"profile,omitempty"`
}

var users = make(map[string]userRecord)
var courses = make(map[string]courseRecord)
var chapters = make(map[string]map[string]chapterRecord)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8081"
    }

    mux := http.NewServeMux()
    mux.HandleFunc("/health", handleHealth)
    mux.HandleFunc("/auth/register", handleRegister)
    mux.HandleFunc("/auth/login", handleLogin)
    mux.HandleFunc("/auth/validate", handleValidate)
    mux.HandleFunc("/auth/me", handleMe)
    mux.HandleFunc("/users/me", handleUserMe)
    mux.HandleFunc("/users/profile", handleUserProfile)
    mux.HandleFunc("/courses", handleCourses)
    mux.HandleFunc("/courses/", handleCourseByID)
    mux.HandleFunc("/courses/chapters", handleCourseChapters)
    mux.HandleFunc("/courses/chapters/", handleCourseChapterByID)

    fmt.Printf("core listening on :%s\n", port)
    if err := http.ListenAndServe(":"+port, mux); err != nil {
        log.Fatal(err)
    }
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
    writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "core"})
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
        return
    }

    var req registerRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
        return
    }

    req.Email = strings.TrimSpace(strings.ToLower(req.Email))
    if req.Name == "" || req.Email == "" || len(req.Password) < 8 {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name, valid email, and password (8+ chars) required"})
        return
    }

    if _, exists := findUserByEmail(req.Email); exists {
        writeJSON(w, http.StatusConflict, map[string]string{"error": "user already exists"})
        return
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to hash password"})
        return
    }

    user := userRecord{
        ID:           generateID(),
        Name:         req.Name,
        Email:        req.Email,
        PasswordHash: string(hash),
    }
    users[user.ID] = user

    token, err := issueToken(user)
    if err != nil {
        writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to issue token"})
        return
    }

    writeJSON(w, http.StatusCreated, authResponse{Token: token, User: userResponse{ID: user.ID, Name: user.Name, Email: user.Email}})
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
        return
    }

    var req loginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
        return
    }

    req.Email = strings.TrimSpace(strings.ToLower(req.Email))
    if req.Email == "" || req.Password == "" {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "email and password required"})
        return
    }

    user, exists := findUserByEmail(req.Email)
    if !exists {
        writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
        return
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
        return
    }

    token, err := issueToken(user)
    if err != nil {
        writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "unable to issue token"})
        return
    }

    writeJSON(w, http.StatusOK, authResponse{Token: token, User: userResponse{ID: user.ID, Name: user.Name, Email: user.Email}})
}

func handleValidate(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
        return
    }

    tokenString, err := bearerTokenFromRequest(r)
    if err != nil {
        writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
        return
    }

    user, err := validateToken(tokenString)
    if err != nil {
        writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
        return
    }

    writeJSON(w, http.StatusOK, map[string]any{"valid": true, "user": userResponse{ID: user.ID, Name: user.Name, Email: user.Email}})
}

func handleMe(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
        return
    }

    tokenString, err := bearerTokenFromRequest(r)
    if err != nil {
        writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
        return
    }

    user, err := validateToken(tokenString)
    if err != nil {
        writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
        return
    }

    writeJSON(w, http.StatusOK, userResponse{ID: user.ID, Name: user.Name, Email: user.Email, Profile: user.Profile})
}

func handleUserMe(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
        return
    }

    tokenString, err := bearerTokenFromRequest(r)
    if err != nil {
        writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
        return
    }

    user, err := validateToken(tokenString)
    if err != nil {
        writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
        return
    }

    writeJSON(w, http.StatusOK, userResponse{ID: user.ID, Name: user.Name, Email: user.Email, Profile: user.Profile})
}

func handleUserProfile(w http.ResponseWriter, r *http.Request) {
    tokenString, err := bearerTokenFromRequest(r)
    if err != nil {
        writeJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
        return
    }

    user, err := validateToken(tokenString)
    if err != nil {
        writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
        return
    }

    switch r.Method {
    case http.MethodGet:
        writeJSON(w, http.StatusOK, user.Profile)
    case http.MethodPut, http.MethodPatch:
        var req profileRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
            return
        }

        if req.Phone != "" {
            user.Profile.Phone = req.Phone
        }
        if req.AvatarURL != "" {
            user.Profile.AvatarURL = req.AvatarURL
        }
        if req.ClassLevel != "" {
            user.Profile.ClassLevel = req.ClassLevel
        }
        if req.Board != "" {
            user.Profile.Board = req.Board
        }
        if req.Preferences != nil {
            user.Profile.Preferences = req.Preferences
        }

        users[user.ID] = user
        writeJSON(w, http.StatusOK, user.Profile)
    default:
        writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
    }
}

func handleCourses(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        courseList := make([]courseRecord, 0, len(courses))
        for _, course := range courses {
            courseList = append(courseList, course)
        }
        writeJSON(w, http.StatusOK, courseList)
        return
    }

    if r.Method == http.MethodPost {
        var req courseRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
            return
        }

        if strings.TrimSpace(req.Title) == "" || strings.TrimSpace(req.SubjectID) == "" {
            writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title and subjectId required"})
            return
        }

        slug := req.Slug
        if slug == "" {
            slug = strings.ToLower(strings.ReplaceAll(req.Title, " ", "-"))
        }

        record := courseRecord{
            ID:        generateID(),
            Title:     req.Title,
            Slug:      slug,
            SubjectID: req.SubjectID,
            TeacherID: req.TeacherID,
            CreatedAt: time.Now().UTC().Format(time.RFC3339),
        }
        courses[record.ID] = record
        writeJSON(w, http.StatusCreated, record)
        return
    }

    writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
}

func handleCourseByID(w http.ResponseWriter, r *http.Request) {
    path := strings.TrimPrefix(r.URL.Path, "/courses")
    if path == "" || path == "/" {
        handleCourses(w, r)
        return
    }

    id := strings.TrimPrefix(path, "/")
    if r.Method == http.MethodGet {
        course, ok := courses[id]
        if !ok {
            writeJSON(w, http.StatusNotFound, map[string]string{"error": "course not found"})
            return
        }
        writeJSON(w, http.StatusOK, course)
        return
    }

    writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
}

func handleCourseChapters(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        courseID := r.URL.Query().Get("courseId")
        if courseID == "" {
            writeJSON(w, http.StatusBadRequest, map[string]string{"error": "courseId query param required"})
            return
        }
        result := make([]chapterRecord, 0)
        for _, chapter := range chapters[courseID] {
            result = append(result, chapter)
        }
        writeJSON(w, http.StatusOK, result)
        return
    }

    if r.Method == http.MethodPost {
        var req chapterRequest
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
            return
        }

        courseID := r.URL.Query().Get("courseId")
        if courseID == "" || strings.TrimSpace(req.Title) == "" {
            writeJSON(w, http.StatusBadRequest, map[string]string{"error": "courseId and title required"})
            return
        }

        chapter := chapterRecord{
            ID:       generateID(),
            CourseID: courseID,
            Title:    req.Title,
            Order:    req.Order,
        }
        if chapter.Order == 0 {
            chapter.Order = len(chapters[courseID]) + 1
        }

        if chapters[courseID] == nil {
            chapters[courseID] = make(map[string]chapterRecord)
        }
        chapters[courseID][chapter.ID] = chapter
        writeJSON(w, http.StatusCreated, chapter)
        return
    }

    writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
}

func handleCourseChapterByID(w http.ResponseWriter, r *http.Request) {
    courseID := r.URL.Query().Get("courseId")
    if courseID == "" {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "courseId query param required"})
        return
    }

    path := strings.TrimPrefix(r.URL.Path, "/courses/chapters")
    if path == "" || path == "/" {
        handleCourseChapters(w, r)
        return
    }

    chapterID := strings.TrimPrefix(path, "/")
    if chapterID == "" {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "chapter id required"})
        return
    }

    if r.Method == http.MethodGet {
        chapter, ok := chapters[courseID][chapterID]
        if !ok {
            writeJSON(w, http.StatusNotFound, map[string]string{"error": "chapter not found"})
            return
        }
        writeJSON(w, http.StatusOK, chapter)
        return
    }

    writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
}

func generateID() string {
    return fmt.Sprintf("user_%d", time.Now().UnixNano())
}

func findUserByEmail(email string) (userRecord, bool) {
    for _, user := range users {
        if user.Email == email {
            return user, true
        }
    }
    return userRecord{}, false
}

func issueToken(user userRecord) (string, error) {
    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        jwtSecret = "studybuddy-dev-secret"
    }

    claims := jwt.RegisteredClaims{
        Subject:   user.ID,
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
        IssuedAt:  jwt.NewNumericDate(time.Now()),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(jwtSecret))
}

func validateToken(tokenString string) (userRecord, error) {
    if tokenString == "" {
        return userRecord{}, errors.New("empty token")
    }

    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" {
        jwtSecret = "studybuddy-dev-secret"
    }

    token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return []byte(jwtSecret), nil
    })
    if err != nil || !token.Valid {
        return userRecord{}, err
    }

    claims, ok := token.Claims.(*jwt.RegisteredClaims)
    if !ok || claims.Subject == "" {
        return userRecord{}, errors.New("missing subject")
    }

    for _, user := range users {
        if user.ID == claims.Subject {
            return user, nil
        }
    }

    return userRecord{}, errors.New("user not found")
}

func bearerTokenFromRequest(r *http.Request) (string, error) {
    headerValue := strings.TrimSpace(r.Header.Get("Authorization"))
    if headerValue != "" {
        return strings.TrimPrefix(headerValue, "Bearer "), nil
    }

    for _, part := range strings.Split(r.Header.Get("Cookie"), ";") {
        key, value, ok := strings.Cut(strings.TrimSpace(part), "=")
        if ok && strings.EqualFold(key, "token") {
            return strings.TrimSpace(value), nil
        }
    }

    return "", errors.New("missing bearer token")
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    if err := json.NewEncoder(w).Encode(payload); err != nil {
        log.Printf("writeJSON error: %v", err)
    }
}
