package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mafi020/social/internal/dto"
	"github.com/mafi020/social/internal/errs"
	"github.com/mafi020/social/internal/utils"
)

type userKey string

const targetUserCtx userKey = "user"

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	paramUserId := chi.URLParam(r, "userID")
	userId, err := strconv.ParseInt(paramUserId, 10, 64)
	if err != nil {
		app.badRequestError(w, r, errors.New("invalid user ID"))
		return
	}
	ctx := r.Context()
	user, err := app.store.Users.GetById(ctx, userId)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}
	if err := utils.JSONResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	paramUserID := chi.URLParam(r, "userID")

	userID, err := strconv.ParseInt(paramUserID, 10, 64)
	if err != nil {
		app.badRequestError(w, r, errors.New("invalid user ID"))
		return
	}

	ctx := r.Context()
	if err := app.store.Users.Delete(ctx, userID); err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			app.notFoundError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := utils.JSONResponse(w, http.StatusOK, map[string]string{"message": "User deleted successfully"}); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	targetUser := getTargetUserFromContext(r) // user from route param
	loggedInUserID := int64(1)                // TODO: get from auth

	log.Printf("[FOLLOW] targetUser.ID=%d, loggedInUserID=%d\n", targetUser.ID, loggedInUserID)

	if targetUser.ID == loggedInUserID {
		app.badRequestError(w, r, errors.New("you cannot follow yourself"))
		return
	}

	if err := app.store.Followers.Follow(r.Context(), targetUser.ID, loggedInUserID); err != nil {
		switch {
		case errors.Is(err, errs.ErrDuplicateEntry):
			app.badRequestError(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := utils.JSONResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	targetUser := getTargetUserFromContext(r) // user from URL param
	loggedInUserID := int64(4)                // TODO: Replace with actual logged-in user ID from auth/session

	if targetUser.ID == loggedInUserID {
		app.badRequestError(w, r, errors.New("you cannot unfollow yourself"))
		return
	}

	ctx := r.Context()
	if err := app.store.Followers.UnFollow(ctx, targetUser.ID, loggedInUserID); err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := utils.JSONResponse(w, http.StatusOK, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

type feedResponse struct {
	Feed       []dto.Feed     `json:"feed"`
	Pagination dto.Pagination `json:"pagination"`
}

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {
	userID := int64(1) // TODO: replace with auth user ID

	params := utils.ParseQueryParams(r)

	page := utils.ParseIntWithDefaultAndMax(params["page"], 1, 0)
	limit := utils.ParseIntWithDefaultAndMax(params["limit"], 25, 100)
	tags := utils.ParseCSV(params["tags"])
	search := params["search"]

	queryParams := dto.FeedQueryParams{
		Page:   page,
		Limit:  limit,
		Tags:   tags,
		Search: search,
	}

	feed, totalCount, err := app.store.Posts.Feed(r.Context(), userID, queryParams)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	log.Printf("Feeds %v", feed)

	if len(feed) == 0 {
		feed = []dto.Feed{}
	}

	resp := feedResponse{
		Feed: feed,
		Pagination: dto.Pagination{
			Page:       queryParams.Page,
			Limit:      queryParams.Limit,
			TotalCount: totalCount,
			TotalPages: (totalCount + queryParams.Limit - 1) / queryParams.Limit,
		},
	}

	if err := utils.JSONResponse(w, http.StatusOK, resp); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

/* ---------------User Context Middleware----------- */
func (app *application) userFromRouteMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		paramUserID := chi.URLParam(r, "userID")
		userID, err := strconv.ParseInt(paramUserID, 10, 64)
		if err != nil {
			app.badRequestError(w, r, errors.New("invalid user ID"))
			return
		}

		ctx := r.Context()
		user, err := app.store.Users.GetById(ctx, userID)
		if err != nil {
			switch {
			case errors.Is(err, errs.ErrNotFound):
				app.notFoundError(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		// Store the *target user* in context
		ctx = context.WithValue(ctx, targetUserCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getTargetUserFromContext(r *http.Request) *dto.User {
	user, _ := r.Context().Value(targetUserCtx).(*dto.User)
	return user
}
