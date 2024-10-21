package app_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"go-postgres-boilerplate/app"
	"go-postgres-boilerplate/dao"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type MockPersonDAO struct {
	mock.Mock
}

func (m *MockPersonDAO) Save(p *dao.Person) error {
	args := m.Called(p)
	return args.Error(0)
}

func (m *MockPersonDAO) GetPersonByName(name string) (*dao.Person, error) {
	args := m.Called(name)
	return args.Get(0).(*dao.Person), args.Error(1)
}

func (m *MockPersonDAO) GetPersonById(id string) (*dao.Person, error) {
	args := m.Called(id)
	return args.Get(0).(*dao.Person), args.Error(1)
}

func (m *MockPersonDAO) GetAllPersons() ([]*dao.Person, error) {
	args := m.Called()
	return args.Get(0).([]*dao.Person), args.Error(1)
}

func (m *MockPersonDAO) DeletePerson(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockPersonDAO) GetPersonByUID(uid string) (*dao.Person, error) {
	args := m.Called(uid)
	return args.Get(0).(*dao.Person), args.Error(1)
}

type WebTestSuite struct {
	suite.Suite
	app     *app.App
	mockDAO *MockPersonDAO
}

func (suite *WebTestSuite) SetupTest() {
	suite.mockDAO = new(MockPersonDAO)
	suite.app = app.NewApp(suite.mockDAO)
}

type HappyPathTestSuite struct {
	WebTestSuite
}

func (suite *HappyPathTestSuite) TestGetAllPersons() {
	persons := []*dao.Person{
		{ID: 1, Name: "John Doe"},
		{ID: 2, Name: "Jane Doe"},
	}

	suite.mockDAO.On("GetAllPersons").Return(persons, nil)

	req, _ := http.NewRequest("GET", "/person", nil)
	rr := httptest.NewRecorder()

	handler := suite.app.SetupRouter()
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusOK, rr.Code)

	var response []*dao.Person
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal(persons, response)

	suite.mockDAO.AssertExpectations(suite.T())
}

func (suite *HappyPathTestSuite) TestCreatePerson() {
	newPerson := &dao.Person{Name: "New Person"}

	suite.mockDAO.On("GetPersonByName", "New Person").Return((*dao.Person)(nil), dao.ErrPersonNotFound)
	suite.mockDAO.On("Save", mock.AnythingOfType("*dao.Person")).Return(nil)

	personJSON, _ := json.Marshal(newPerson)
	req, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(personJSON))
	rr := httptest.NewRecorder()

	handler := suite.app.SetupRouter()
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusCreated, rr.Code)

	var response dao.Person
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	suite.NoError(err)
	suite.Equal(newPerson.Name, response.Name)

	suite.mockDAO.AssertExpectations(suite.T())
}

type SadPathTestSuite struct {
	WebTestSuite
}

func (suite *SadPathTestSuite) TestGetPersonNotFound() {
	suite.mockDAO.On("GetPersonById", "999").Return((*dao.Person)(nil), sql.ErrNoRows)

	req, _ := http.NewRequest("GET", "/person/999", nil)
	rr := httptest.NewRecorder()

	handler := suite.app.SetupRouter()
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusNotFound, rr.Code)
	suite.mockDAO.AssertExpectations(suite.T())
}

func (suite *SadPathTestSuite) TestCreatePersonAlreadyExists() {
	existingPerson := &dao.Person{Name: "Existing Person"}

	suite.mockDAO.On("GetPersonByName", "Existing Person").Return(existingPerson, nil)

	personJSON, _ := json.Marshal(existingPerson)
	req, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(personJSON))
	rr := httptest.NewRecorder()

	handler := suite.app.SetupRouter()
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusConflict, rr.Code)
	suite.mockDAO.AssertExpectations(suite.T())
}

type NetworkIssueTestSuite struct {
	WebTestSuite
}

func (suite *NetworkIssueTestSuite) TestGetAllPersonsNetworkError() {
	suite.mockDAO.On("GetAllPersons").Return([]*dao.Person(nil), errors.New("network error"))

	req, _ := http.NewRequest("GET", "/person", nil)
	rr := httptest.NewRecorder()

	handler := suite.app.SetupRouter()
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusInternalServerError, rr.Code)
	suite.mockDAO.AssertExpectations(suite.T())
}

func (suite *NetworkIssueTestSuite) TestCreatePersonNetworkError() {
	newPerson := &dao.Person{Name: "New Person"}

	suite.mockDAO.On("GetPersonByName", "New Person").Return((*dao.Person)(nil), dao.ErrPersonNotFound)
	suite.mockDAO.On("Save", mock.AnythingOfType("*dao.Person")).Return(errors.New("network error"))

	personJSON, _ := json.Marshal(newPerson)
	req, _ := http.NewRequest("POST", "/person", bytes.NewBuffer(personJSON))
	rr := httptest.NewRecorder()

	handler := suite.app.SetupRouter()
	handler.ServeHTTP(rr, req)

	suite.Equal(http.StatusInternalServerError, rr.Code)
	suite.mockDAO.AssertExpectations(suite.T())
}

func TestWebSuites(t *testing.T) {
	suite.Run(t, new(HappyPathTestSuite))
	suite.Run(t, new(SadPathTestSuite))
	suite.Run(t, new(NetworkIssueTestSuite))
}
