package handlers

import (
	"avito-test-task/internal/domain"
	"avito-test-task/pkg"
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}
	err = json.Unmarshal(body, &houseRequest)
	if err != nil {
		h.lg.Warn("unmarshal request body error")
		respBody = CreateErrorResponse(r.Context(), UnmarshalHTTPBodyError, UnmarshalHTTPBodyMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.dbTimeout*time.Second)
	defer cancel()

	houseResponse, err = h.uc.Create(ctx, &houseRequest, h.lg)
	if err != nil {
		h.lg.Warn("house handler: create error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), CreateHouseError, CreateHouseErrorMsg)
		w.WriteHeader(GetReturnHTTPCode(err))
		w.Write(respBody)
		return
	}

	respBody, err = json.Marshal(houseResponse)
	if err != nil {
		h.lg.Warn("house handler: create error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), MarshalHTTPBodyError, MarshalHTTPBodyErrorMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}

	w.Write(respBody)
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
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.dbTimeout*time.Second)
	defer cancel()

	role, err := pkg.ExtractPayloadFromToken(r.Header.Get("authorization"), "role")
	if err != nil {
		h.lg.Warn("house handler: get flats by id error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), ExtractRoleFromTokenError, ExtractRoleFromTokenErrorMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}

	var status string
	if role == domain.Moderator {
		status = domain.AnyStatus
	} else {
		status = domain.ApprovedStatus
	}

	flats, err := h.uc.GetFlatsByHouseID(ctx, id, status, h.lg)
	if err != nil {
		h.lg.Warn("house handler: get flats by id error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), GetFlatsByHouseIDError, GetFlatsByHouseIDErrorMsg)
		w.WriteHeader(GetReturnHTTPCode(err))
		w.Write(respBody)
		return
	}

	respBody, err = json.Marshal(flats)
	if err != nil {
		h.lg.Warn("house handler: get flats by id error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), MarshalHTTPBodyError, MarshalHTTPBodyErrorMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}

	w.Write(respBody)
}

func (h *HouseHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	var (
		respBody []byte
	)
	defer r.Body.Close()

	fmt.Println(r.URL.Path)

	userID, err := pkg.ExtractPayloadFromToken(r.Header.Get("authorization"), "userID")
	if err != nil {
		h.lg.Warn("house handler: subscribe error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), SubscribeOnHouseError, SubscribeOnHouseErrorMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}
	userUuid, err := uuid.Parse(userID)
	if err != nil {
		h.lg.Warn("house handler: create error: extract id", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), CreateFlatError, CreateFlatErrorMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	idString := pathParts[len(pathParts)-2]
	id, err := strconv.Atoi(idString)
	if err != nil {
		h.lg.Warn("house handler: subscribe error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), ParseURLError, ParseURLErrorMsg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(respBody)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), h.dbTimeout*time.Second)
	defer cancel()

	err = h.uc.SubscribeByID(ctx, id, userUuid, h.lg)
	if err != nil {
		h.lg.Warn("house handler: subscribe error", zap.Error(err))
		respBody = CreateErrorResponse(r.Context(), SubscribeOnHouseError, SubscribeOnHouseErrorMsg)
		w.WriteHeader(GetReturnHTTPCode(err))
		w.Write(respBody)
		return
	}

	w.WriteHeader(http.StatusOK)
}
