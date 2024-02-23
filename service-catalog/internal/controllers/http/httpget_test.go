package httpapi

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"v1/internal/cashe"
	"v1/internal/entity"
	liberrors "v1/internal/lib/errors"
	"v1/internal/service"
	"v1/internal/service/mocks"
	storage "v1/internal/storage/gorm"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace/noop"
)

func TestCatalogGETFailHandler(t *testing.T) {

	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}),
	)
	notFoundId := 1
	tracer := noop.NewTracerProvider().Tracer("")

	req, err := http.NewRequest("GET", "/catalog/"+strconv.Itoa(notFoundId), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	db := mocks.NewDBRepository(t)
	db.
		On("Get", mock.Anything, notFoundId).
		Return(nil, liberrors.WrapErr("storage.get", storage.ErrRecordNotFound)).
		Once()

	cashedb := mocks.NewCasheRepository(t)
	cashedb.
		On("Get", mock.Anything, strconv.Itoa(notFoundId)).
		Return(nil, liberrors.WrapErr("cashe.get", cashe.ErrRecordNotFound)).
		Once()

	usercases := service.NewService(log, db, cashedb, tracer)
	srvhttp, _, _ := NewServer(usercases, tracer)
	srvhttp.R.ServeHTTP(rr, req)

	require.Equal(t, rr.Code, http.StatusInternalServerError)

	body := rr.Body.String()
	var resp Response
	require.NoError(t, json.Unmarshal([]byte(body), &resp))

}

func TestCatalogGETHandler(t *testing.T) {

	log := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}),
	)
	wantin := &entity.Product{Id: 1, Name: "Name", Description: "Description"}
	wantout := GetResponse{Name: wantin.Name, Description: wantin.Description}
	tracer := noop.NewTracerProvider().Tracer("")

	req, err := http.NewRequest("GET", "/catalog/"+strconv.Itoa(wantin.Id), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	db := mocks.NewDBRepository(t)
	db.
		On("Get", mock.Anything, wantin.Id).
		Return(wantin, nil).
		Once()

	cashedb := mocks.NewCasheRepository(t)
	cashedb.
		On("Get", mock.Anything, strconv.Itoa(wantin.Id)).
		Return(nil, liberrors.WrapErr("cashe.get", cashe.ErrRecordNotFound)).
		Once()
	cashedb.
		On("Set", mock.Anything, strconv.Itoa(wantin.Id), mock.AnythingOfType("*entity.Product")).
		Return(nil).
		Once()

	usercases := service.NewService(log, db, cashedb, tracer)
	srvhttp, _, _ := NewServer(usercases, tracer)
	srvhttp.R.ServeHTTP(rr, req)

	require.Equal(t, rr.Code, http.StatusOK)

	body := rr.Body.String()
	var resp GetResponse
	require.NoError(t, json.Unmarshal([]byte(body), &resp))

	require.Equal(t, wantout, resp)

}
