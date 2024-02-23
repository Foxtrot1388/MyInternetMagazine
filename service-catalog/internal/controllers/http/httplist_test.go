package httpapi

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"v1/internal/entity"
	"v1/internal/service"
	"v1/internal/service/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestCatalogListErrorDBHandler(t *testing.T) {

	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}),
	)
	tracer := noop.NewTracerProvider().Tracer("")

	req, err := http.NewRequest("GET", "/catalog/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	db := mocks.NewDBRepository(t)
	db.
		On("List", mock.Anything).
		Return(nil, errors.New("Test error")).
		Once()
	cashedb := mocks.NewCasheRepository(t)

	usercases := service.NewService(log, db, cashedb, tracer)
	srvhttp, _, _ := NewServer(usercases, tracer)
	srvhttp.R.ServeHTTP(rr, req)

	require.Equal(t, rr.Code, http.StatusInternalServerError)

	body := rr.Body.String()
	var resp response
	require.NoError(t, json.Unmarshal([]byte(body), &resp))

	require.Equal(t, notGetList, resp)

}

func TestCatalogListHandler(t *testing.T) {

	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}),
	)
	tracer := noop.NewTracerProvider().Tracer("")

	t.Run("TestCatalogListHandler", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/catalog/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		el1 := entity.ElementOfList{Id: 1, Name: "Test1"}
		el2 := entity.ElementOfList{Id: 2, Name: "Test2"}
		pardb := make([]entity.ElementOfList, 0)
		pardb = append(pardb, el1)
		pardb = append(pardb, el2)
		wantlist := make(listResponse, 0)
		wantlist = append(wantlist, responseElement{Id: el1.Id, Name: el1.Name})
		wantlist = append(wantlist, responseElement{Id: el2.Id, Name: el2.Name})

		db := mocks.NewDBRepository(t)
		db.
			On("List", mock.Anything).
			Return(&pardb, nil).
			Once()

		cashedb := mocks.NewCasheRepository(t)

		usercases := service.NewService(log, db, cashedb, tracer)
		srvhttp, _, _ := NewServer(usercases, tracer)
		srvhttp.R.ServeHTTP(rr, req)

		require.Equal(t, rr.Code, http.StatusOK)

		body := rr.Body.String()
		var resp listResponse
		require.NoError(t, json.Unmarshal([]byte(body), &resp))

		require.Equal(t, wantlist, resp)
	})

	t.Run("TestCatalogEmptyListHandler", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/catalog/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		pardb := make([]entity.ElementOfList, 0)
		wantlist := make(listResponse, 0)

		db := mocks.NewDBRepository(t)
		db.
			On("List", mock.Anything).
			Return(&pardb, nil).
			Once()

		cashedb := mocks.NewCasheRepository(t)

		usercases := service.NewService(log, db, cashedb, tracer)
		srvhttp, _, _ := NewServer(usercases, tracer)
		srvhttp.R.ServeHTTP(rr, req)

		require.Equal(t, rr.Code, http.StatusOK)

		body := rr.Body.String()
		var resp listResponse
		require.NoError(t, json.Unmarshal([]byte(body), &resp))

		require.Equal(t, wantlist, resp)
	})

}
