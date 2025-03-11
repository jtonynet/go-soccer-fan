package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ginRoutesSuite struct {
	suite.Suite
	r *ginRouter
}

func TestGinRoutesSuite(t *testing.T) {
	suite.Run(t, new(ginRoutesSuite))
}

func (suite *ginRoutesSuite) SetupSuite() {
	suite.r = newGinRouter()
}

func (suite *ginRoutesSuite) TestHappyPath() {
	req, err := http.NewRequest("GET", "/campeonatos", nil)
	assert.NoError(suite.T(), err)

	resp := httptest.NewRecorder()
	suite.r.Router.ServeHTTP(resp, req)
	assert.Equal(suite.T(), http.StatusOK, resp.Code)
}
