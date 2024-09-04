package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	liberrors "v1/internal/lib/errors"
	"v1/internal/service"
	"v1/internal/service/mocks"
	storage "v1/internal/storage/gorm"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestCatalogDeleteFailHandler(t *testing.T) {

	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}),
	)
	notFoundId := 1
	tracer := noop.NewTracerProvider().Tracer("")

	req, err := http.NewRequest("DELETE", "/catalog/"+strconv.Itoa(notFoundId), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	db := mocks.NewDBRepository(t)
	db.
		On("Delete", mock.Anything, notFoundId).
		Return(false, liberrors.WrapErr("storage.get", storage.ErrRecordNotFound)).
		Once()

	cashedb := mocks.NewCasheRepository(t)

	usercases := service.NewService(log, db, cashedb, tracer)
	srvhttp, _, _ := NewServer(usercases, tracer)
	srvhttp.R.ServeHTTP(rr, req)

	require.Equal(t, rr.Code, http.StatusInternalServerError)

	body := rr.Body.String()
	var resp response
	require.NoError(t, json.Unmarshal([]byte(body), &resp))

	require.Equal(t, notDeleteProduct, resp)

}

func TestCatalogDeleteWrongParametrID(t *testing.T) {

	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}),
	)
	tracer := noop.NewTracerProvider().Tracer("")

	db := mocks.NewDBRepository(t)
	cashedb := mocks.NewCasheRepository(t)
	usercases := service.NewService(log, db, cashedb, tracer)
	srvhttp, _, _ := NewServer(usercases, tracer)

	req, err := http.NewRequest("DELETE", "/catalog/test", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	srvhttp.R.ServeHTTP(rr, req)

	require.Equal(t, rr.Code, http.StatusInternalServerError)

	body := rr.Body.String()
	var resp response
	require.NoError(t, json.Unmarshal([]byte(body), &resp))

	require.Equal(t, notDeleteProduct, resp)

}

func TestCatalogDeleteHandler(t *testing.T) {

	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}),
	)
	wantinId := 1
	tracer := noop.NewTracerProvider().Tracer("")

	req, err := http.NewRequest("DELETE", "/catalog/"+strconv.Itoa(wantinId), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	db := mocks.NewDBRepository(t)
	db.
		On("Delete", mock.Anything, wantinId).
		Return(true, nil).
		Once()

	cashedb := mocks.NewCasheRepository(t)
	cashedb.
		On("Invalidate", mock.Anything, strconv.Itoa(wantinId)).
		Return(nil).
		Once()

	usercases := service.NewService(log, db, cashedb, tracer)
	srvhttp, _, _ := NewServer(usercases, tracer)
	srvhttp.R.ServeHTTP(rr, req)

	require.Equal(t, rr.Code, http.StatusOK)

	body := rr.Body.String()
	var resp deleteResponse
	require.NoError(t, json.Unmarshal([]byte(body), &resp))

	require.Equal(t, deleteResponse{
		OK: true,
	}, resp)

}
