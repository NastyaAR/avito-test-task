package handlers

import (
	"avito-test-task/internal/domain"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type HouseHandler struct {
	uc domain.HouseUsecase
	lg *zap.Logger
}

func NewHouseHandler(uc domain.HouseUsecase, lg *zap.Logger) *HouseHandler {
	return &HouseHandler{uc, lg}
}

func (h *HouseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var houseRequest domain.CreateHouseRequest
	var houseResponse domain.CreateHouseResponse
	var respBody []byte
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

	houseResponse, err := usecase.
		w.Write(respBody)
	w.WriteHeader(http.StatusOK)
}

type GetRoomsRequest struct {
	ID int `json:"id"`
}

type GetRoomsResponse struct {
	Flats []domain.Flat
}

func (h *HouseHandler) GetRoomsByID(w http.ResponseWriter, r *http.Request) {
	requestID := middleware.GetReqID(r.Context())
	lg := h.lg.With(zap.String("RequestID", requestID))
	lg.Info("HouseHandler: GetRoomsByID - start")

	defer r.Body.Close()

	var req GetRoomsRequest
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		return
	}

	flats, err := h.uc.GetFlatsByID(context.TODO(), req.ID, lg)
	if err != nil {

	}

	resp, err := json.Marshal(flats)
	if err != nil {
		lg.Warn("HouseHandler: GetRoomsByID - marshal error", zap.Error(err), zap.Any(flats))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
	}

	lg.Info("HouseHandler: GetRoomsByID - success")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
