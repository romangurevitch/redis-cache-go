package server

import (
	"bytes"
	config "github.com/romangurevitch/redis-cache-go"
	"github.com/romangurevitch/redis-cache-go/crypto"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

// Get contact handler, get contact from the cache if available or redirect the request to Autopilot servers.
func (s *server) getContactHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		if match, err := regexp.MatchString("/contact/[^/]+$", r.URL.Path); err != nil || !match {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		contactId := strings.TrimPrefix(r.URL.Path, "/contact/")
		if contactId == "" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if contact, ok := s.load(contactId, r.Header.Get(config.ApiKeyHeader)); ok {
			s.logger.Printf("loading %s from cache\n", contactId)
			w.WriteHeader(http.StatusOK)
			_, err := w.Write(contact)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		s.logger.Printf("redirecting to %v\n", s.apiUrl)
		s.redirect(s.apiUrl, w, r, func(resp *http.Response) error {
			if resp.StatusCode != http.StatusOK {
				return nil
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			s.logger.Printf("caching %s", contactId)
			s.store(contactId, r.Header.Get(config.ApiKeyHeader), body)

			resp.Body = ioutil.NopCloser(bytes.NewReader(body))
			return nil
		})
	}
}

// Create new contact, invalidate the cache.
func (s *server) postContactHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		s.logger.Printf("redirecting to %v\n", s.apiUrl)
		s.redirect(s.apiUrl, w, r, func(resp *http.Response) error {
			if resp.StatusCode == http.StatusOK {
				s.logger.Printf("invalidating cache")
				return s.cache.Invalidate()
			}
			return nil
		})
	}
}

// Cache the response body
func (s *server) store(contactCacheKey, apiKey string, contact []byte) {
	err := s.cache.Store(crypto.Hash(contactCacheKey, apiKey), contact)
	if err != nil {
		// failed cache should not fail the request
		s.logger.Println(err.Error())
	}
}

// Load contact from the cache if available
func (s *server) load(contactCacheKey, apiKey string) ([]byte, bool) {
	contact, err := s.cache.Load(crypto.Hash(contactCacheKey, apiKey))
	if err != nil {
		// failed cache should not fail the request
		s.logger.Println(err.Error())
		return nil, false
	}

	return contact, contact != nil
}

// Logging handler
func (s *server) log(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Println(r.Method, r.URL.Path, r.RemoteAddr)
		h(w, r)
	}
}
