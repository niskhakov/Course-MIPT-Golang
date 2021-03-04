package urlshortener

import (
	"github.com/go-chi/chi"
	"math/rand"
	"net/http"
	"net/url"
	"time"
)

type URLShortener struct {
	addr     string
	size     int
	database map[string]string
	capacity int
}

func NewShortener(addr string) *URLShortener {
	rand.Seed(time.Now().UnixNano())
	return &URLShortener{
		addr:     addr,
		size:     5,
		database: make(map[string]string),
		capacity: 916132832 / 2, // 916132832 = len(letters)^5 - theoretical max for 5-symbol combinations
	}
}

func (s *URLShortener) HandleSave(rw http.ResponseWriter, req *http.Request) {
	_, err := url.Parse(req.URL.Query().Get("u"))

	if err != nil {
		http.Error(rw, "Missing key", http.StatusBadRequest)
		return
	}

	if s.capacity == 0 {
		http.Error(rw, "Out of capacity to create short url for specified URLShortener",
			http.StatusInternalServerError)
		return
	}

	shortUrlCandidate := randSequence(s.size)

	// Checking for collision, if exists try again
	// Capacity should be significantly less than theoretical max to be able to quickly find rand sequence
	for {
		_, prs := s.database[shortUrlCandidate]
		if prs == false {
			// No collision
			break
		}
		// Next try to generate random sequence
		shortUrlCandidate = randSequence(s.size)
	}

	param := req.URL.Query().Get("u")
	s.database[shortUrlCandidate] = param

	s.capacity--
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(s.addr + "/" + shortUrlCandidate))

}

func (s *URLShortener) HandleExpand(rw http.ResponseWriter, req *http.Request) {
	param := chi.URLParam(req, "key")
	addr, prs := s.database[param]
	if prs == false {
		http.Error(rw, "Not Found", http.StatusNotFound)
		return
	}

	http.Redirect(rw, req, addr, http.StatusMovedPermanently)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randSequence(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
