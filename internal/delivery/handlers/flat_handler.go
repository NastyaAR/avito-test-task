package handlers

import (
	"avito-test-task/internal/domain"
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

type FlatHandler struct {
	uc        domain.FlatUsecase
	lg        *zap.Logger
	dbTimeout time.Duration
}

func NewFlatHandler(uc domain.FlatUsecase, timeout time.Duration, lg *zap.Logger) *FlatHandler {
	return &FlatHandler{
		uc:        uc,
		lg:        lg,
		dbTimeout: timeout,
	}
}

func (h *FlatHandler) Create(w http.ResponseWriter, r *http.Request) {
	var (
		respBody     []byte
		flatRequest  domain.CreateFlatRequest
		flatResponse domain.CreateFlatResponse
	)

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.lg.Warn("flat handler: create error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), ReadHTTPBodyError, ReadHTTPBodyMsg)
		w.Write(respBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &flatRequest)
	if err != nil {
		h.lg.Warn("flat handler: create error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), UnmarshalHTTPBodyError, UnmarshalHTTPBodyMsg)
		w.Write(respBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	flatResponse, err = h.uc.Create(ctx, &flatRequest, h.lg)
	if err != nil {
		h.lg.Warn("flat handler: create error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), CreateFlatError, CreateFlatErrorMsg)
		w.Write(respBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respBody, err = json.Marshal(flatResponse)
	if err != nil {
		h.lg.Warn("flat handler: create error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), MarshalHTTPBodyError, MarshalHTTPBodyErrorMsg)
		w.Write(respBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(respBody)
	w.WriteHeader(http.StatusOK)
}

func (h *FlatHandler) Update(w http.ResponseWriter, r *http.Request) {
	var (
		respBody     []byte
		flatRequest  domain.UpdateFlatRequest
		flatResponse domain.CreateFlatResponse
	)

	defer r.Body.Close()

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.lg.Warn("flat handler: update error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), ReadHTTPBodyError, ReadHTTPBodyMsg)
		w.Write(respBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.Unmarshal(body, &flatRequest)
	if err != nil {
		h.lg.Warn("flat handler: update error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), UnmarshalHTTPBodyError, UnmarshalHTTPBodyMsg)
		w.Write(respBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	flatResponse, err = h.uc.Update(ctx, &flatRequest, h.lg)
	if err != nil {
		h.lg.Warn("flat handler: update error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), UpdateFlatError, UpdateFlatErrorMsg)
		w.Write(respBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	respBody, err = json.Marshal(flatResponse)
	if err != nil {
		h.lg.Warn("flat handler: update error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), MarshalHTTPBodyError, MarshalHTTPBodyErrorMsg)
		w.Write(respBody)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(respBody)
	w.WriteHeader(http.StatusOK)
}
