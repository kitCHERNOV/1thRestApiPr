package save

import (
	resp "RestApiFP/1thRestApiPr/internal/lib/api/response"
	"RestApiFP/1thRestApiPr/internal/lib/logger/sl"
	"RestApiFP/1thRestApiPr/internal/lib/random"
	"RestApiFP/1thRestApiPr/internal/storage"
	"errors"
	"net/http"

	"log/slog"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

type Request struct {
	URL string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	Alias string `json:"alias,omitempty"`
}

// TODO: move to config if needed
const aliasLength = 6

// go generate
//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLSaver --output=../1thRestApiPr/internal/http-server/handlers/url/save/mocksss --log-level=debug
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to decode request"))

			return 
		}

		log.Info("request body decoded", slog.Any("request", req))
	
		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			
			log.Error("invalid request", sl.Err(err))

			render.JSON(w, r, resp.ValidattonError(validateErr))

			return 
		}

		alias := req.Alias
		if alias == "" {
			// TODO: testing colision for aliases
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if errors.Is(err, storage.ErrURLExists) {
			log.Info("url already exists", slog.String("url",req.URL))

			render.JSON(w, r, resp.Error("url already exists"))

			return
		}
		if err != nil {
			log.Error("failed to add url", sl.Err(err))

			render.JSON(w, r, resp.Error("failed to add url"))

			return
		}

		log.Info("url added", slog.Int64("id", id))

		// TODO: wrapping this strings to own function
		render.JSON(w, r, Response {
			Response: resp.OK(),
			Alias: alias,
		})
	}
}