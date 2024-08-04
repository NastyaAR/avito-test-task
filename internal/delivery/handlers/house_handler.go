package handlers

import (
	"avito-test-task/internal/domain"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type HouseHandler struct {
	uc        domain.HouseUsecase
	lg        *zap.Logger
	dbTimeout time.Duration
}

func NewHouseHandler(uc domain.HouseUsecase, timeout time.Duration, lg *zap.Logger) *HouseHandler {
	return &HouseHandler{uc, lg, timeout}
}

func (h *HouseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var (
		houseRequest  domain.CreateHouseRequest
		houseResponse domain.CreateHouseResponse
		respBody      []byte
	)

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.lg.Warn("read http body error")
		respBody = CreateErrorResponse(r.Context(), ReadHTTPBodyError, ReadHTTPBodyMsg)
		w.Write(respBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &houseRequest)
	if err != nil {
		h.lg.Warn("unmarshal request body error")
		respBody = CreateErrorResponse(r.Context(), UnmarshalHTTPBodyError, UnmarshalHTTPBodyMsg)
		w.Write(respBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.dbTimeout)
	defer cancel()

	houseResponse, err = h.uc.Create(ctx, &houseRequest, h.lg)
	if err != nil {
		h.lg.Warn("house handler: create error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), CreateHouseError, CreateHouseErrorMsg)
		w.Write(respBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respBody, err = json.Marshal(houseResponse)
	if err != nil {
		h.lg.Warn("house handler: create error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), MarshalHTTPBodyError, MarshalHTTPBodyErrorMsg)
		w.Write(respBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(respBody)
	w.WriteHeader(http.StatusOK)
}

func (h *HouseHandler) GetFlatsByID(w http.ResponseWriter, r *http.Request) {
	var (
		respBody []byte
	)
	defer r.Body.Close()

	pathParts := strings.Split(r.URL.Path, "/")
	idString := pathParts[len(pathParts)-1]
	id, err := strconv.Atoi(idString)
	if err != nil {
		h.lg.Warn("house handler: get flats by id error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), ParseURLError, ParseURLErrorMsg)
		w.Write(respBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.dbTimeout)
	defer cancel()

	flats, err := h.uc.GetFlatsByHouseID(ctx, id, r.Header.Get("status"), h.lg)
	if err != nil {
		h.lg.Warn("house handler: get flats by id error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), GetFlatsByHouseIDError, GetFlatsByHouseIDErrorMsg)
		w.Write(respBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respBody, err = json.Marshal(flats)
	if err != nil {
		h.lg.Warn("house handler: get flats by id error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), MarshalHTTPBodyError, MarshalHTTPBodyErrorMsg)
		w.Write(respBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(respBody)
	w.WriteHeader(http.StatusOK)
}
