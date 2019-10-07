package server

import (
	"bytes"
	"encoding/json"
	"github.com/alicebob/miniredis"
	config "github.com/romangurevitch/redis-cache-go"
	"github.com/romangurevitch/redis-cache-go/cache"
	"github.com/romangurevitch/redis-cache-go/test"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type testResources struct {
	testServer    *httptest.Server
	miniRedis     *miniredis.Miniredis
	redisCache    cache.Cache
	contactServer *server
}

func TestGetSimpleContact(t *testing.T) {
	testServerRespContent := map[string]string{"contact_id": "0", "Email": "some@email.com"}
	testServer := mockSimpleTestServer(t, testServerRespContent)

	contactServer := mockContactServer(t, testServer)
	defer contactServer.close(t)

	req := newGetRequest(t, "contactId", "apiKey")
	resp := exec(contactServer, req)

	checkStatusCode(t, resp, http.StatusOK)
	checkContent(t, resp, testServerRespContent)
}

func TestVerifyCacheHit(t *testing.T) {
	testServerRespContent := map[string]string{"contact_id": "0", "Email": "some@email.com"}
	testServer := mockSimpleTestServer(t, testServerRespContent)

	contactServer := mockContactServer(t, testServer)
	defer contactServer.close(t)

	const contactId = "contactId"
	const apiKey = "apiKey"

	firstReq := newGetRequest(t, contactId, apiKey)
	firstResp := exec(contactServer, firstReq)

	checkStatusCode(t, firstResp, http.StatusOK)
	checkContent(t, firstResp, testServerRespContent)

	testServer.Close()
	secondReq := newGetRequest(t, contactId, apiKey)
	secondResp := exec(contactServer, secondReq)

	checkStatusCode(t, secondResp, http.StatusOK)
	checkContent(t, secondResp, testServerRespContent)
}

func TestVerifyCacheMissForDifferentApiKeys(t *testing.T) {
	testServerRespContent := map[string]string{"contact_id": "0", "Email": "some@email.com"}
	testServer := mockSimpleTestServer(t, testServerRespContent)

	contactServer := mockContactServer(t, testServer)
	defer contactServer.close(t)

	const contactId = "contactId"

	firstReq := newGetRequest(t, contactId, "apiKey")
	firstResp := exec(contactServer, firstReq)

	checkStatusCode(t, firstResp, http.StatusOK)
	checkContent(t, firstResp, testServerRespContent)

	testServer.Close()
	secondReq := newGetRequest(t, contactId, "differentApiKey")
	secondResp := exec(contactServer, secondReq)

	checkStatusCode(t, secondResp, http.StatusBadGateway)
}

func TestInvalidateCacheAfterPost(t *testing.T) {
	testServerRespContent := map[string]string{"contact_id": "0"}
	testServer := mockSimpleTestServer(t, testServerRespContent)

	contactServer := mockContactServer(t, testServer)
	defer contactServer.close(t)

	const contactId = "contactId"
	const apiKey = "apiKey"

	firstReq := newGetRequest(t, contactId, apiKey)
	firstResp := exec(contactServer, firstReq)

	checkStatusCode(t, firstResp, http.StatusOK)
	checkContent(t, firstResp, testServerRespContent)

	content, err := json.Marshal(testServerRespContent)
	test.CheckError(t, err)

	postReq := newPostRequest(t, content, apiKey)
	postResp := exec(contactServer, postReq)

	checkStatusCode(t, postResp, http.StatusOK)
	checkContent(t, postResp, testServerRespContent)

	testServer.Close()
	secondReq := newGetRequest(t, contactId, apiKey)
	secondResp := exec(contactServer, secondReq)

	checkStatusCode(t, secondResp, http.StatusBadGateway)
}

func checkStatusCode(t *testing.T, response *httptest.ResponseRecorder, expected int) {
	if response.Code != expected {
		t.Fatalf("expected: %v, got: %v", http.StatusText(expected), http.StatusText(response.Code))
	}
}

func checkContent(t *testing.T, response *httptest.ResponseRecorder, expected map[string]string) {
	var actual map[string]string
	err := json.Unmarshal(response.Body.Bytes(), &actual)
	test.CheckError(t, err)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("expected: %v, got: %v", expected, actual)
	}
}

func exec(resource *testResources, req *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	resource.contactServer.router.ServeHTTP(w, req)
	return w
}

func newGetRequest(t *testing.T, contactId, apiKey string) *http.Request {
	req, err := http.NewRequest("GET", "http://localhost:"+config.HttpPort+"/contact/"+contactId, nil)
	test.CheckError(t, err)
	req.Header.Set(config.ApiKeyHeader, apiKey)
	return req
}

func newPostRequest(t *testing.T, content []byte, apiKey string) *http.Request {
	req, err := http.NewRequest("POST", "http://localhost:"+config.HttpPort+"/contact", nil)
	test.CheckError(t, err)

	req.Body = ioutil.NopCloser(bytes.NewReader(content))
	req.Header.Set(config.ApiKeyHeader, apiKey)
	return req
}

func mockSimpleTestServer(t *testing.T, response map[string]string) *httptest.Server {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		w.Header().Set("Content-Type", "application/json")

		err := encoder.Encode(response)
		test.CheckError(t, err)
	}))
	return testServer
}

func mockContactServer(t *testing.T, testServer *httptest.Server) *testResources {
	miniRedis, err := miniredis.Run()
	test.CheckError(t, err)
	redisCache, err := cache.NewRedis("tcp", miniRedis.Addr(), config.RedisPoolSize)
	test.CheckError(t, err)
	server, err := NewContactServer(testServer.URL, redisCache)
	test.CheckError(t, err)

	// set up the routes
	server.routes()
	return &testResources{
		testServer:    testServer,
		miniRedis:     miniRedis,
		redisCache:    redisCache,
		contactServer: server,
	}
}

func (r *testResources) close(t *testing.T) {
	r.testServer.Close()
	r.miniRedis.Close()
	test.CheckError(t, r.redisCache.Close())
}
