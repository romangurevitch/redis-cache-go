package contact

import (
	"bytes"
	config "github.com/romangurevitch/redis-cache-go"
	"github.com/romangurevitch/redis-cache-go/crypto"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

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

		redirectUrl, err := s.apiUrl.Parse(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s.logger.Printf("redirecting to %v\n", redirectUrl)
		s.redirect(redirectUrl, w, r, func(resp *http.Response) error {
			if resp.StatusCode != http.StatusOK {
				return nil
			}

			s.logger.Printf("caching %s", contactId)
			body, err := s.store(contactId, r.Header.Get(config.ApiKeyHeader), resp.Body)
			if err != nil {
				return err
			}

			resp.Body = ioutil.NopCloser(bytes.NewReader(body))
			return nil
		})
	}
}

func (s *server) postContactHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		redirectUrl, err := s.apiUrl.Parse(r.URL.Path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		s.logger.Printf("redirecting to %v\n", redirectUrl)
		s.redirect(redirectUrl, w, r, func(resp *http.Response) error {
			if resp.StatusCode == http.StatusOK {
				s.logger.Printf("invalidating cache")
				return s.cache.Invalidate()
			}
			return nil
		})
	}
}

func (s *server) store(contactCacheKey, apiKey string, body io.ReadCloser) ([]byte, error) {
	content, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}

	err = s.cache.Store(crypto.Hash(contactCacheKey, apiKey), content)
	if err != nil {
		// failed cache should not fail the request
		s.logger.Println(err.Error())
	}

	return content, nil
}

func (s *server) load(contactCacheKey, apiKey string) ([]byte, bool) {
	contact, err := s.cache.Load(crypto.Hash(contactCacheKey, apiKey))
	if err != nil {
		// failed cache should not fail the request
		s.logger.Println(err.Error())
		return nil, false
	}

	return contact, contact != nil
}

func (s *server) log(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.logger.Println(r.Method, r.URL.Path, r.RemoteAddr)
		h(w, r)
	}
}
