package httpapi

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"v1/internal/service"
	"v1/internal/service/mocks"
	storage "v1/internal/storage/gorm"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestCatalogCreateFailHandler(t *testing.T) {

	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}),
	)
	tracer := noop.NewTracerProvider().Tracer("")

	t.Run("DescriptionValidation", func(t *testing.T) {

		newcatalog := newCatalog{Name: "Name", Description: ""}
		example, _ := json.Marshal(newcatalog)

		req, err := http.NewRequest("POST", "/catalog/new", bytes.NewReader(example))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		db := mocks.NewDBRepository(t)
		cashedb := mocks.NewCasheRepository(t)

		usercases := service.NewService(log, db, cashedb, tracer)
		srvhttp, _, _ := NewServer(usercases, tracer)
		srvhttp.R.ServeHTTP(rr, req)

		require.Equal(t, rr.Code, http.StatusInternalServerError)

		body := rr.Body.String()
		var resp response
		require.NoError(t, json.Unmarshal([]byte(body), &resp))

	})

	t.Run("NameValidation", func(t *testing.T) {

		newcatalog := newCatalog{Name: "", Description: "Description"}
		example, _ := json.Marshal(newcatalog)

		req, err := http.NewRequest("POST", "/catalog/new", bytes.NewReader(example))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		db := mocks.NewDBRepository(t)
		cashedb := mocks.NewCasheRepository(t)

		usercases := service.NewService(log, db, cashedb, tracer)
		srvhttp, _, _ := NewServer(usercases, tracer)
		srvhttp.R.ServeHTTP(rr, req)

		require.Equal(t, rr.Code, http.StatusInternalServerError)

		body := rr.Body.String()
		var resp response
		require.NoError(t, json.Unmarshal([]byte(body), &resp))

	})

	t.Run("ErrorDBCreate", func(t *testing.T) {

		newcatalog := newCatalog{Name: "Name", Description: "Description"}
		example, _ := json.Marshal(newcatalog)

		req, err := http.NewRequest("POST", "/catalog/new", bytes.NewReader(example))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		db := mocks.NewDBRepository(t)
		db.
			On("Create", mock.Anything, mock.AnythingOfType("*entity.Product")).
			Return(0, storage.ErrRecordNotFound).
			Once()

		cashedb := mocks.NewCasheRepository(t)
		usercases := service.NewService(log, db, cashedb, tracer)
		srvhttp, _, _ := NewServer(usercases, tracer)
		srvhttp.R.ServeHTTP(rr, req)

		require.Equal(t, rr.Code, http.StatusInternalServerError)

		body := rr.Body.String()
		var resp response
		require.NoError(t, json.Unmarshal([]byte(body), &resp))

		require.Equal(t, notCreate, resp)

	})

}

func TestCatalogCreateHandler(t *testing.T) {

	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}),
	)
	newcatalog := newCatalog{Name: "Name", Description: "Description"}
	wantresponse := createResponse{
		Id: 1,
	}
	example, _ := json.Marshal(newcatalog)
	tracer := noop.NewTracerProvider().Tracer("")

	req, err := http.NewRequest("POST", "/catalog/new", bytes.NewReader(example))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	db := mocks.NewDBRepository(t)
	db.
		On("Create", mock.Anything, mock.AnythingOfType("*entity.Product")).
		Return(wantresponse.Id, nil).
		Once()

	cashedb := mocks.NewCasheRepository(t)
	usercases := service.NewService(log, db, cashedb, tracer)
	srvhttp, _, _ := NewServer(usercases, tracer)
	srvhttp.R.ServeHTTP(rr, req)

	require.Equal(t, rr.Code, http.StatusOK)

	body := rr.Body.String()
	var resp createResponse
	require.NoError(t, json.Unmarshal([]byte(body), &resp))

	require.Equal(t, wantresponse, resp)

}
