package physical

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Netcracker/dbaas-adapter-core/pkg/dao"
	"github.com/Netcracker/dbaas-opensearch-adapter/basic"
	cl "github.com/Netcracker/dbaas-opensearch-adapter/client"
	"github.com/Netcracker/dbaas-opensearch-adapter/common"
	"github.com/stretchr/testify/assert"
)

func TestRegistration(t *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		if _, err := res.Write([]byte("")); err != nil {
			assert.Fail(t, "Can't write to ResponseWriter", err)
		}
	}))
	defer func() { testServer.Close() }()

	registrationService := NewRegistrationProvider(
		testServer.URL,
		dao.BasicAuth{
			Username: "cluster-dba",
			Password: "test",
		},
		"",
		nil,
		150000,
		60000,
		5000,
		"tmp-test",
		"http://dbaas-opensearch-adapter.elasticsearch-cluster:8080",
		dao.BasicAuth{
			Username: "dbaas-aggregator",
			Password: "dbaas-aggregator",
		},
		basic.NewBaseProvider(nil),
	)

	registrationService.doRegistrationRequest()
	assert.Equal(t, registrationService.Health, common.ComponentHealth{Status: "OK"})
}

func TestFailedRegistration(t *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(500)
		if _, err := res.Write([]byte("")); err != nil {
			assert.Fail(t, "Can't write to ResponseWriter", err)
		}
	}))
	defer func() { testServer.Close() }()

	registrationService := NewRegistrationProvider(
		testServer.URL,
		dao.BasicAuth{
			Username: "cluster-dba",
			Password: "test",
		},
		"",
		nil,
		150000,
		60000,
		5000,
		"tmp-test",
		"http://dbaas-opensearch-adapter.elasticsearch-cluster:8080",
		dao.BasicAuth{
			Username: "dbaas-aggregator",
			Password: "dbaas-aggregator",
		},
		basic.NewBaseProvider(nil),
	)

	registrationService.doRegistrationRequest()
	assert.Equal(t, registrationService.Health, common.ComponentHealth{Status: "PROBLEM"})
}

func TestApiVersion(t *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		if _, err := res.Write([]byte("{\"major\":3,\"minor\":4,\"supportedMajors\":[1,2,3]}")); err != nil {
			assert.Fail(t, "Can't write to ResponseWriter", err)
		}
	}))
	defer func() { testServer.Close() }()

	apiVersion := getApiVersion(testServer.URL, cl.ConfigureClient())
	assert.Equal(t, apiVersion, common.ApiV2)
}

func TestApiVersionWhenV3NotSupported(t *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		if _, err := res.Write([]byte("{\"major\":2,\"minor\":4,\"supportedMajors\":[1,2]}")); err != nil {
			assert.Fail(t, "Can't write to ResponseWriter", err)
		}
	}))
	defer func() { testServer.Close() }()

	apiVersion := getApiVersion(testServer.URL, cl.ConfigureClient())
	assert.Equal(t, apiVersion, common.ApiV1)
}

func TestApiVersionWithIncorrectData(t *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		if _, err := res.Write([]byte("{\"major\":2,\"minor\":4,\"supportedMajors\":\"\"}")); err != nil {
			assert.Fail(t, "Can't write to ResponseWriter", err)
		}
	}))
	defer func() { testServer.Close() }()

	apiVersion := getApiVersion(testServer.URL, cl.ConfigureClient())
	assert.Equal(t, apiVersion, common.ApiV2)
}

func TestApiVersionWhenApiVersionNotFound(t *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(401)
		if _, err := res.Write([]byte("")); err != nil {
			assert.Fail(t, "Can't write to ResponseWriter", err)
		}
	}))
	defer func() { testServer.Close() }()

	apiVersion := getApiVersion(testServer.URL, cl.ConfigureClient())
	assert.Equal(t, apiVersion, common.ApiV1)
}
