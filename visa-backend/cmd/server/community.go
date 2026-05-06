package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
)

type communityAuthorItem struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type communityPostItem struct {
	ID           int64               `json:"id"`
	UserID       int64               `json:"userId"`
	Category     string              `json:"category"`
	Title        string              `json:"title"`
	Content      string              `json:"content"`
	Image        string              `json:"image"`
	Images       []string            `json:"images"`
	Status       string              `json:"status"`
	ReviewNote   string              `json:"reviewNote"`
	LikeCount    int                 `json:"likeCount"`
	ReportCount  int                 `json:"reportCount"`
	CommentCount int                 `json:"commentCount"`
	Liked        bool                `json:"liked"`
	Favorited    bool                `json:"favorited"`
	CreatedAt    string              `json:"createdAt"`
	UpdatedAt    string              `json:"updatedAt"`
	Author       communityAuthorItem `json:"author"`
}

type communityCreatePostRequest struct {
	Category string   `json:"category"`
	Title    string   `json:"title"`
	Content  string   `json:"content"`
	Image    string   `json:"image"`
	Images   []string `json:"images"`
}

type communityCommentItem struct {
	ID              int64                  `json:"id"`
	PostID          int64                  `json:"postId"`
	UserID          int64                  `json:"userId"`
	ParentID        int64                  `json:"parentId"`
	RootID          int64                  `json:"rootId"`
	ReplyToNickname string                 `json:"replyToNickname"`
	Content         string                 `json:"content"`
	Image           string                 `json:"image"`
	Status          string                 `json:"status"`
	ReviewNote      string                 `json:"reviewNote"`
	ReplyCount      int                    `json:"replyCount"`
	CreatedAt       string                 `json:"createdAt"`
	UpdatedAt       string                 `json:"updatedAt"`
	Author          communityAuthorItem    `json:"author"`
	Replies         []communityCommentItem `json:"replies,omitempty"`
}

type communityCreateCommentRequest struct {
	Content  string `json:"content"`
	Image    string `json:"image"`
	ParentID int64  `json:"parentId"`
}

type communityReportRequest struct {
	Reason string `json:"reason"`
	Detail string `json:"detail"`
}

type communityReportItem struct {
	ID         int64               `json:"id"`
	PostID     int64               `json:"postId"`
	PostTitle  string              `json:"postTitle"`
	ReporterID int64               `json:"reporterId"`
	Reason     string              `json:"reason"`
	Detail     string              `json:"detail"`
	Status     string              `json:"status"`
	CreatedAt  string              `json:"createdAt"`
	UpdatedAt  string              `json:"updatedAt"`
	Reporter   communityAuthorItem `json:"reporter"`
}

type communityCommentReportItem struct {
	ID             int64               `json:"id"`
	CommentID      int64               `json:"commentId"`
	PostID         int64               `json:"postId"`
	PostTitle      string              `json:"postTitle"`
	CommentContent string              `json:"commentContent"`
	ReporterID     int64               `json:"reporterId"`
	Reason         string              `json:"reason"`
	Detail         string              `json:"detail"`
	Status         string              `json:"status"`
	CreatedAt      string              `json:"createdAt"`
	UpdatedAt      string              `json:"updatedAt"`
	Reporter       communityAuthorItem `json:"reporter"`
	CommentAuthor  communityAuthorItem `json:"commentAuthor"`
}

type communityKeywordRuleItem struct {
	ID        int64  `json:"id"`
	Keyword   string `json:"keyword"`
	Action    string `json:"action"`
	Enabled   bool   `json:"enabled"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type communityModerationResult struct {
	Status     string
	ReviewNote string
}

func (s *appServer) handleCommunityPosts(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		uid, _ := s.optionalUserIDFromHeader(r)
		items, err := s.listCommunityPosts(r.Context(), uid, strings.TrimSpace(r.URL.Query().Get("q")), strings.TrimSpace(r.URL.Query().Get("category")), false)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: items})
	case http.MethodPost:
		uid, err := s.userIDFromHeader(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err)
			return
		}
		var in communityCreatePostRequest
		if err := readJSON(r, &in); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		item, message, err := s.createCommunityPost(r.Context(), uid, in)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusCreated, apiResponse{Message: message, Data: item})
	default:
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}

func (s *appServer) handleCommunityFavorites(w http.ResponseWriter, r *http.Request, uid int64) {
	switch r.Method {
	case http.MethodGet:
		items, err := s.listCommunityFavorites(r.Context(), uid)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: items})
	case http.MethodPost:
		var in struct {
			PostID int64 `json:"postId"`
		}
		if err := readJSON(r, &in); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		item, err := s.favoriteCommunityPost(r.Context(), in.PostID, uid)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "favorited", Data: item})
	default:
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}

func (s *appServer) handleCommunityFavoriteByID(w http.ResponseWriter, r *http.Request, uid int64) {
	postID, ok := parseID(w, r.URL.Path, "/api/community/favorites/")
	if !ok {
		return
	}
	if r.Method != http.MethodDelete {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	item, err := s.unfavoriteCommunityPost(r.Context(), postID, uid)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Message: "unfavorited", Data: item})
}

func (s *appServer) handleCommunityMyPosts(w http.ResponseWriter, r *http.Request, uid int64) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	items, err := s.listCommunityMyPosts(
		r.Context(),
		uid,
		strings.TrimSpace(r.URL.Query().Get("q")),
		strings.TrimSpace(r.URL.Query().Get("status")),
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: items})
}

func (s *appServer) handleCommunityPostRoutes(w http.ResponseWriter, r *http.Request) {
	id, rest, ok := parseIDWithRest(w, r.URL.Path, "/api/community/posts/")
	if !ok {
		return
	}
	switch rest {
	case "":
		switch r.Method {
		case http.MethodGet:
			uid, _ := s.optionalUserIDFromHeader(r)
			item, err := s.getCommunityPostByID(r.Context(), id, uid, false)
			if err != nil {
				writeError(w, http.StatusNotFound, errors.New("post not found"))
				return
			}
			writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: item})
		case http.MethodDelete:
			uid, err := s.userIDFromHeader(r)
			if err != nil {
				writeError(w, http.StatusUnauthorized, err)
				return
			}
			if err := s.deleteOwnCommunityPost(r.Context(), id, uid); err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			writeJSON(w, http.StatusOK, apiResponse{Message: "deleted"})
		default:
			writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		}
	case "comments":
		switch r.Method {
		case http.MethodGet:
			uid, _ := s.optionalUserIDFromHeader(r)
			items, err := s.listCommunityComments(r.Context(), id, uid)
			if err != nil {
				writeError(w, http.StatusInternalServerError, err)
				return
			}
			writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: items})
		case http.MethodPost:
			uid, err := s.userIDFromHeader(r)
			if err != nil {
				writeError(w, http.StatusUnauthorized, err)
				return
			}
			var in communityCreateCommentRequest
			if err := readJSON(r, &in); err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			item, message, err := s.createCommunityComment(r.Context(), id, uid, in)
			if err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			writeJSON(w, http.StatusCreated, apiResponse{Message: message, Data: item})
		default:
			writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		}
	case "like":
		uid, err := s.userIDFromHeader(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err)
			return
		}
		switch r.Method {
		case http.MethodPost:
			item, err := s.likeCommunityPost(r.Context(), id, uid)
			if err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			writeJSON(w, http.StatusOK, apiResponse{Message: "liked", Data: item})
		case http.MethodDelete:
			item, err := s.unlikeCommunityPost(r.Context(), id, uid)
			if err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			writeJSON(w, http.StatusOK, apiResponse{Message: "unliked", Data: item})
		default:
			writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		}
	case "report":
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
			return
		}
		uid, err := s.userIDFromHeader(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err)
			return
		}
		var in communityReportRequest
		if err := readJSON(r, &in); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		report, err := s.reportCommunityPost(r.Context(), id, uid, in)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusCreated, apiResponse{Message: "reported", Data: report})
	default:
		writeError(w, http.StatusNotFound, errors.New("route not found"))
	}
}

func (s *appServer) handleCommunityCommentByID(w http.ResponseWriter, r *http.Request) {
	id, rest, ok := parseIDWithRest(w, r.URL.Path, "/api/community/comments/")
	if !ok {
		return
	}
	switch rest {
	case "":
		switch r.Method {
		case http.MethodDelete:
			uid, err := s.userIDFromHeader(r)
			if err != nil {
				writeError(w, http.StatusUnauthorized, err)
				return
			}
			if err := s.deleteOwnCommunityComment(r.Context(), id, uid); err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			writeJSON(w, http.StatusOK, apiResponse{Message: "deleted"})
		default:
			writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		}
	case "report":
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
			return
		}
		uid, err := s.userIDFromHeader(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err)
			return
		}
		var in communityReportRequest
		if err := readJSON(r, &in); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		report, err := s.reportCommunityComment(r.Context(), id, uid, in)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusCreated, apiResponse{Message: "reported", Data: report})
	default:
		writeError(w, http.StatusNotFound, errors.New("route not found"))
	}
}

func (s *appServer) handleCommunityUserRoutes(w http.ResponseWriter, r *http.Request) {
	id, rest, ok := parseIDWithRest(w, r.URL.Path, "/api/community/users/")
	if !ok {
		return
	}
	if rest != "block" {
		writeError(w, http.StatusNotFound, errors.New("route not found"))
		return
	}
	uid, err := s.userIDFromHeader(r)
	if err != nil {
		writeError(w, http.StatusUnauthorized, err)
		return
	}
	switch r.Method {
	case http.MethodPost:
		if err := s.blockCommunityUser(r.Context(), uid, id); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "blocked"})
	case http.MethodDelete:
		if err := s.unblockCommunityUser(r.Context(), uid, id); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "unblocked"})
	default:
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}

func (s *appServer) handleCommunityAdminPosts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	items, err := s.listCommunityPosts(r.Context(), 0, strings.TrimSpace(r.URL.Query().Get("q")), strings.TrimSpace(r.URL.Query().Get("category")), true)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	status := strings.TrimSpace(r.URL.Query().Get("status"))
	if status != "" {
		filtered := make([]communityPostItem, 0, len(items))
		for _, item := range items {
			if item.Status == status {
				filtered = append(filtered, item)
			}
		}
		items = filtered
	}
	writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: items})
}

func (s *appServer) handleCommunityAdminPostsAdmin(w http.ResponseWriter, r *http.Request, _ int64) {
	s.handleCommunityAdminPosts(w, r)
}

func (s *appServer) handleCommunityAdminPostByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r.URL.Path, "/api/community/admin/posts/")
	if !ok {
		return
	}
	switch r.Method {
	case http.MethodPatch:
		var in struct {
			Status     string `json:"status"`
			ReviewNote string `json:"reviewNote"`
		}
		if err := readJSON(r, &in); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		item, err := s.updateCommunityPostStatus(r.Context(), id, in.Status, in.ReviewNote)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "updated", Data: item})
	case http.MethodDelete:
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		if _, err := s.db.ExecContext(ctx, `DELETE FROM community_posts WHERE id=?`, id); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "deleted"})
	default:
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}

func (s *appServer) handleCommunityAdminPostByIDAdmin(w http.ResponseWriter, r *http.Request, _ int64) {
	s.handleCommunityAdminPostByID(w, r)
}

func (s *appServer) handleCommunityAdminComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	items, err := s.listCommunityCommentsForAdmin(r.Context(), strings.TrimSpace(r.URL.Query().Get("status")))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: items})
}

func (s *appServer) handleCommunityAdminCommentsAdmin(w http.ResponseWriter, r *http.Request, _ int64) {
	s.handleCommunityAdminComments(w, r)
}

func (s *appServer) handleCommunityAdminCommentByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r.URL.Path, "/api/community/admin/comments/")
	if !ok {
		return
	}
	switch r.Method {
	case http.MethodPatch:
		var in struct {
			Status     string `json:"status"`
			ReviewNote string `json:"reviewNote"`
		}
		if err := readJSON(r, &in); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		item, err := s.updateCommunityCommentStatus(r.Context(), id, in.Status, in.ReviewNote)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "updated", Data: item})
	case http.MethodDelete:
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		if _, err := s.db.ExecContext(ctx, `DELETE FROM community_comments WHERE id=?`, id); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		if err := s.refreshCommunityCounts(ctx); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "deleted"})
	default:
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}

func (s *appServer) handleCommunityAdminCommentByIDAdmin(w http.ResponseWriter, r *http.Request, _ int64) {
	s.handleCommunityAdminCommentByID(w, r)
}

func (s *appServer) handleCommunityAdminReports(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	items, err := s.listCommunityReports(r.Context(), strings.TrimSpace(r.URL.Query().Get("status")))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: items})
}

func (s *appServer) handleCommunityAdminReportsAdmin(w http.ResponseWriter, r *http.Request, _ int64) {
	s.handleCommunityAdminReports(w, r)
}

func (s *appServer) handleCommunityAdminReportByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r.URL.Path, "/api/community/admin/reports/")
	if !ok {
		return
	}
	if r.Method != http.MethodPatch {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	var in struct {
		Status     string `json:"status"`
		PostStatus string `json:"postStatus"`
		ReviewNote string `json:"reviewNote"`
	}
	if err := readJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	item, err := s.updateCommunityReportStatus(r.Context(), id, in.Status, in.PostStatus, in.ReviewNote)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Message: "updated", Data: item})
}

func (s *appServer) handleCommunityAdminReportByIDAdmin(w http.ResponseWriter, r *http.Request, _ int64) {
	s.handleCommunityAdminReportByID(w, r)
}

func (s *appServer) handleCommunityAdminCommentReports(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	items, err := s.listCommunityCommentReports(r.Context(), strings.TrimSpace(r.URL.Query().Get("status")))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: items})
}

func (s *appServer) handleCommunityAdminCommentReportsAdmin(w http.ResponseWriter, r *http.Request, _ int64) {
	s.handleCommunityAdminCommentReports(w, r)
}

func (s *appServer) handleCommunityAdminCommentReportByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r.URL.Path, "/api/community/admin/comment-reports/")
	if !ok {
		return
	}
	if r.Method != http.MethodPatch {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	var in struct {
		Status        string `json:"status"`
		CommentStatus string `json:"commentStatus"`
		ReviewNote    string `json:"reviewNote"`
	}
	if err := readJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	item, err := s.updateCommunityCommentReportStatus(r.Context(), id, in.Status, in.CommentStatus, in.ReviewNote)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Message: "updated", Data: item})
}

func (s *appServer) handleCommunityAdminCommentReportByIDAdmin(w http.ResponseWriter, r *http.Request, _ int64) {
	s.handleCommunityAdminCommentReportByID(w, r)
}

func (s *appServer) handleCommunityAdminKeywords(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items, err := s.listCommunityKeywordRules(r.Context())
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: items})
	case http.MethodPost:
		var in communityKeywordRuleItem
		if err := readJSON(r, &in); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		item, err := s.createCommunityKeywordRule(r.Context(), in)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusCreated, apiResponse{Message: "created", Data: item})
	default:
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}

func (s *appServer) handleCommunityAdminKeywordsAdmin(w http.ResponseWriter, r *http.Request, _ int64) {
	s.handleCommunityAdminKeywords(w, r)
}

func (s *appServer) handleCommunityAdminKeywordByID(w http.ResponseWriter, r *http.Request) {
	id, ok := parseID(w, r.URL.Path, "/api/community/admin/keywords/")
	if !ok {
		return
	}
	switch r.Method {
	case http.MethodPut:
		var in communityKeywordRuleItem
		if err := readJSON(r, &in); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		item, err := s.updateCommunityKeywordRule(r.Context(), id, in)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "updated", Data: item})
	case http.MethodDelete:
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		if _, err := s.db.ExecContext(ctx, `DELETE FROM community_keyword_rules WHERE id=?`, id); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "deleted"})
	default:
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}

func (s *appServer) handleCommunityAdminKeywordByIDAdmin(w http.ResponseWriter, r *http.Request, _ int64) {
	s.handleCommunityAdminKeywordByID(w, r)
}

func (s *appServer) optionalUserIDFromHeader(r *http.Request) (int64, error) {
	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	if auth == "" {
		return 0, nil
	}
	return s.userIDFromHeader(r)
}

func (s *appServer) listCommunityPosts(ctx context.Context, uid int64, q, category string, includeAll bool) ([]communityPostItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	statusFilter := "p.status='published'"
	if includeAll {
		statusFilter = "1=1"
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT
			p.id,p.user_id,p.category,p.title,p.content,p.image,p.images,p.status,p.review_note,p.like_count,p.report_count,p.comment_count,p.created_at,p.updated_at,
			u.id,u.nickname,u.avatar,
			CASE WHEN ?=0 THEN 0 ELSE EXISTS(SELECT 1 FROM community_post_likes l WHERE l.post_id=p.id AND l.user_id=?) END AS liked,
			CASE WHEN ?=0 THEN 0 ELSE EXISTS(SELECT 1 FROM community_post_favorites f WHERE f.post_id=p.id AND f.user_id=?) END AS favorited
		FROM community_posts p
		JOIN users u ON u.id=p.user_id
		WHERE `+statusFilter+`
		  AND (?='' OR p.category=?)
		  AND (?='' OR p.title LIKE CONCAT('%',?,'%') OR p.content LIKE CONCAT('%',?,'%'))
		  AND (?=0 OR NOT EXISTS (
			SELECT 1 FROM community_user_blocks b WHERE b.user_id=? AND b.blocked_user_id=p.user_id
		  ))
		ORDER BY p.created_at DESC
	`, uid, uid, uid, uid, category, category, q, q, q, uid, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]communityPostItem, 0)
	for rows.Next() {
		item, err := scanCommunityPost(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *appServer) listCommunityMyPosts(ctx context.Context, uid int64, q, status string) ([]communityPostItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, `
		SELECT
			p.id,p.user_id,p.category,p.title,p.content,p.image,p.images,p.status,p.review_note,p.like_count,p.report_count,p.comment_count,p.created_at,p.updated_at,
			u.id,u.nickname,u.avatar,
			CASE WHEN ?=0 THEN 0 ELSE EXISTS(SELECT 1 FROM community_post_likes l WHERE l.post_id=p.id AND l.user_id=?) END AS liked,
			CASE WHEN ?=0 THEN 0 ELSE EXISTS(SELECT 1 FROM community_post_favorites f WHERE f.post_id=p.id AND f.user_id=?) END AS favorited
		FROM community_posts p
		JOIN users u ON u.id=p.user_id
		WHERE p.user_id=?
		  AND (?='' OR p.status=?)
		  AND (?='' OR p.title LIKE CONCAT('%',?,'%') OR p.content LIKE CONCAT('%',?,'%'))
		ORDER BY p.created_at DESC
	`, uid, uid, uid, uid, uid, status, status, q, q, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]communityPostItem, 0)
	for rows.Next() {
		item, err := scanCommunityPost(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func scanCommunityPost(scanner interface{ Scan(dest ...any) error }) (communityPostItem, error) {
	var item communityPostItem
	var imagesRaw string
	var created time.Time
	var updated time.Time
	var liked int
	var favorited int
	if err := scanner.Scan(
		&item.ID,
		&item.UserID,
		&item.Category,
		&item.Title,
		&item.Content,
		&item.Image,
		&imagesRaw,
		&item.Status,
		&item.ReviewNote,
		&item.LikeCount,
		&item.ReportCount,
		&item.CommentCount,
		&created,
		&updated,
		&item.Author.ID,
		&item.Author.Nickname,
		&item.Author.Avatar,
		&liked,
		&favorited,
	); err != nil {
		return communityPostItem{}, err
	}
	item.CreatedAt = created.Format(time.RFC3339)
	item.UpdatedAt = updated.Format(time.RFC3339)
	item.Liked = liked == 1
	item.Favorited = favorited == 1
	item.Images = decodeStringList(imagesRaw)
	if len(item.Images) == 0 && strings.TrimSpace(item.Image) != "" {
		item.Images = []string{item.Image}
	}
	return item, nil
}

func encodeStringList(items []string) string {
	if len(items) == 0 {
		return "[]"
	}
	data, err := json.Marshal(items)
	if err != nil {
		return "[]"
	}
	return string(data)
}

func decodeStringList(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []string{}
	}
	var out []string
	if err := json.Unmarshal([]byte(raw), &out); err == nil {
		return normalizeStringList(out)
	}
	return []string{}
}

func normalizeStringList(items []string) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			out = append(out, item)
		}
	}
	return out
}

func (s *appServer) getCommunityPostByID(ctx context.Context, id, uid int64, admin bool) (communityPostItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	statusFilter := "p.status='published' OR (? > 0 AND p.user_id=?)"
	args := []any{uid, uid, uid, uid, uid, uid, id}
	if admin {
		statusFilter = "1=1"
		args = []any{uid, uid, uid, uid, id}
	}
	row := s.db.QueryRowContext(ctx, `
		SELECT
			p.id,p.user_id,p.category,p.title,p.content,p.image,p.images,p.status,p.review_note,p.like_count,p.report_count,p.comment_count,p.created_at,p.updated_at,
			u.id,u.nickname,u.avatar,
			CASE WHEN ?=0 THEN 0 ELSE EXISTS(SELECT 1 FROM community_post_likes l WHERE l.post_id=p.id AND l.user_id=?) END AS liked,
			CASE WHEN ?=0 THEN 0 ELSE EXISTS(SELECT 1 FROM community_post_favorites f WHERE f.post_id=p.id AND f.user_id=?) END AS favorited
		FROM community_posts p
		JOIN users u ON u.id=p.user_id
		WHERE (`+statusFilter+`) AND p.id=?
	`, args...)
	return scanCommunityPost(row)
}

func (s *appServer) favoriteCommunityPost(ctx context.Context, postID, uid int64) (communityPostItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	if postID <= 0 {
		return communityPostItem{}, errors.New("postId is required")
	}
	if _, err := s.db.ExecContext(ctx, `INSERT IGNORE INTO community_post_favorites(post_id,user_id) VALUES(?,?)`, postID, uid); err != nil {
		return communityPostItem{}, err
	}
	return s.getCommunityPostByID(ctx, postID, uid, false)
}

func (s *appServer) unfavoriteCommunityPost(ctx context.Context, postID, uid int64) (communityPostItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	if _, err := s.db.ExecContext(ctx, `DELETE FROM community_post_favorites WHERE post_id=? AND user_id=?`, postID, uid); err != nil {
		return communityPostItem{}, err
	}
	return s.getCommunityPostByID(ctx, postID, uid, false)
}

func (s *appServer) listCommunityFavorites(ctx context.Context, uid int64) ([]communityPostItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, `
		SELECT p.id
		FROM community_post_favorites f
		JOIN community_posts p ON p.id=f.post_id
		WHERE f.user_id=?
		ORDER BY f.created_at DESC
	`, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]communityPostItem, 0)
	for rows.Next() {
		var postID int64
		if err := rows.Scan(&postID); err != nil {
			return nil, err
		}
		item, err := s.getCommunityPostByID(ctx, postID, uid, false)
		if err == nil {
			items = append(items, item)
		}
	}
	return items, nil
}

func (s *appServer) createCommunityPost(ctx context.Context, uid int64, in communityCreatePostRequest) (communityPostItem, string, error) {
	in.Category = normalizeCommunityCategory(in.Category)
	in.Title = strings.TrimSpace(in.Title)
	in.Content = strings.TrimSpace(in.Content)
	in.Image = strings.TrimSpace(in.Image)
	in.Images = normalizeStringList(in.Images)
	if len(in.Images) > 0 {
		in.Image = in.Images[0]
	}
	if in.Title == "" || in.Content == "" {
		return communityPostItem{}, "", errors.New("title and content are required")
	}
	if len([]rune(in.Title)) > 80 {
		return communityPostItem{}, "", errors.New("title is too long")
	}
	if len([]rune(in.Content)) > 2000 {
		return communityPostItem{}, "", errors.New("content is too long")
	}

	moderation, err := s.moderateCommunityContent(ctx, in.Title, in.Content)
	if err != nil {
		return communityPostItem{}, "", err
	}
	if moderation.Status == "rejected" {
		return communityPostItem{}, "", errors.New(emptyFallback(moderation.ReviewNote, "content violates community rules"))
	}

	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	res, err := s.db.ExecContext(ctx, `
		INSERT INTO community_posts(user_id,category,title,content,image,images,status,review_note,like_count,report_count,comment_count)
		VALUES(?,?,?,?,?,?,?,?,0,0,0)
	`, uid, in.Category, in.Title, in.Content, in.Image, encodeStringList(in.Images), moderation.Status, moderation.ReviewNote)
	if err != nil {
		return communityPostItem{}, "", err
	}
	id, _ := res.LastInsertId()
	item, err := s.getCommunityPostByID(ctx, id, uid, true)
	if err != nil {
		return communityPostItem{}, "", err
	}
	if moderation.Status == "review" {
		return item, "submitted_for_review", nil
	}
	return item, "created", nil
}

func normalizeCommunityCategory(category string) string {
	category = strings.TrimSpace(category)
	valid := map[string]bool{
		"推荐":   true,
		"攻略":   true,
		"问答":   true,
		"签证经验": true,
		"材料模板": true,
	}
	if valid[category] {
		return category
	}
	return "推荐"
}

func (s *appServer) moderateCommunityContent(ctx context.Context, title, content string) (communityModerationResult, error) {
	rules, err := s.listCommunityKeywordRules(ctx)
	if err != nil {
		return communityModerationResult{}, err
	}
	text := strings.ToLower(title + "\n" + content)
	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}
		keyword := strings.ToLower(strings.TrimSpace(rule.Keyword))
		if keyword == "" || !strings.Contains(text, keyword) {
			continue
		}
		switch rule.Action {
		case "reject":
			return communityModerationResult{Status: "rejected", ReviewNote: "内容触发社区敏感词规则，请修改后再发布"}, nil
		case "review":
			return communityModerationResult{Status: "review", ReviewNote: "内容已提交审核，审核通过后会公开显示"}, nil
		}
	}
	return communityModerationResult{Status: "published"}, nil
}

func (s *appServer) likeCommunityPost(ctx context.Context, postID, uid int64) (communityPostItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	if _, err := s.db.ExecContext(ctx, `INSERT IGNORE INTO community_post_likes(post_id,user_id) VALUES(?,?)`, postID, uid); err != nil {
		return communityPostItem{}, err
	}
	if _, err := s.db.ExecContext(ctx, `UPDATE community_posts SET like_count=(SELECT COUNT(1) FROM community_post_likes WHERE post_id=?) WHERE id=?`, postID, postID); err != nil {
		return communityPostItem{}, err
	}
	return s.getCommunityPostByID(ctx, postID, uid, false)
}

func (s *appServer) unlikeCommunityPost(ctx context.Context, postID, uid int64) (communityPostItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	if _, err := s.db.ExecContext(ctx, `DELETE FROM community_post_likes WHERE post_id=? AND user_id=?`, postID, uid); err != nil {
		return communityPostItem{}, err
	}
	if _, err := s.db.ExecContext(ctx, `UPDATE community_posts SET like_count=(SELECT COUNT(1) FROM community_post_likes WHERE post_id=?) WHERE id=?`, postID, postID); err != nil {
		return communityPostItem{}, err
	}
	return s.getCommunityPostByID(ctx, postID, uid, false)
}

func (s *appServer) reportCommunityPost(ctx context.Context, postID, uid int64, in communityReportRequest) (communityReportItem, error) {
	in.Reason = strings.TrimSpace(in.Reason)
	in.Detail = strings.TrimSpace(in.Detail)
	if in.Reason == "" {
		return communityReportItem{}, errors.New("reason is required")
	}
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	res, err := s.db.ExecContext(ctx, `
		INSERT INTO community_post_reports(post_id,reporter_id,reason,detail,status)
		VALUES(?,?,?,?, 'open')
		ON DUPLICATE KEY UPDATE reason=VALUES(reason), detail=VALUES(detail), status='open', updated_at=CURRENT_TIMESTAMP
	`, postID, uid, in.Reason, in.Detail)
	if err != nil {
		return communityReportItem{}, err
	}
	if _, err := s.db.ExecContext(ctx, `UPDATE community_posts SET report_count=(SELECT COUNT(1) FROM community_post_reports WHERE post_id=? AND status='open') WHERE id=?`, postID, postID); err != nil {
		return communityReportItem{}, err
	}
	reportID, _ := res.LastInsertId()
	if reportID == 0 {
		_ = s.db.QueryRowContext(ctx, `SELECT id FROM community_post_reports WHERE post_id=? AND reporter_id=?`, postID, uid).Scan(&reportID)
	}
	return s.getCommunityReportByID(ctx, reportID)
}

func (s *appServer) reportCommunityComment(ctx context.Context, commentID, uid int64, in communityReportRequest) (communityCommentReportItem, error) {
	in.Reason = strings.TrimSpace(in.Reason)
	in.Detail = strings.TrimSpace(in.Detail)
	if in.Reason == "" {
		return communityCommentReportItem{}, errors.New("reason is required")
	}
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	var commentUserID int64
	if err := s.db.QueryRowContext(ctx, `SELECT user_id FROM community_comments WHERE id=?`, commentID).Scan(&commentUserID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return communityCommentReportItem{}, errors.New("comment not found")
		}
		return communityCommentReportItem{}, err
	}
	if commentUserID == uid {
		return communityCommentReportItem{}, errors.New("cannot report your own comment")
	}

	res, err := s.db.ExecContext(ctx, `
		INSERT INTO community_comment_reports(comment_id,reporter_id,reason,detail,status)
		VALUES(?,?,?,?, 'open')
		ON DUPLICATE KEY UPDATE reason=VALUES(reason), detail=VALUES(detail), status='open', updated_at=CURRENT_TIMESTAMP
	`, commentID, uid, in.Reason, in.Detail)
	if err != nil {
		return communityCommentReportItem{}, err
	}
	reportID, _ := res.LastInsertId()
	if reportID == 0 {
		_ = s.db.QueryRowContext(ctx, `SELECT id FROM community_comment_reports WHERE comment_id=? AND reporter_id=?`, commentID, uid).Scan(&reportID)
	}
	return s.getCommunityCommentReportByID(ctx, reportID)
}

func (s *appServer) blockCommunityUser(ctx context.Context, uid, blockedUserID int64) error {
	if uid == blockedUserID {
		return errors.New("cannot block yourself")
	}
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	_, err := s.db.ExecContext(ctx, `INSERT IGNORE INTO community_user_blocks(user_id,blocked_user_id) VALUES(?,?)`, uid, blockedUserID)
	return err
}

func (s *appServer) unblockCommunityUser(ctx context.Context, uid, blockedUserID int64) error {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	_, err := s.db.ExecContext(ctx, `DELETE FROM community_user_blocks WHERE user_id=? AND blocked_user_id=?`, uid, blockedUserID)
	return err
}

func (s *appServer) deleteOwnCommunityPost(ctx context.Context, id, uid int64) error {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	res, err := s.db.ExecContext(ctx, `DELETE FROM community_posts WHERE id=? AND user_id=?`, id, uid)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return errors.New("post not found")
	}
	return nil
}

func (s *appServer) updateCommunityPostStatus(ctx context.Context, id int64, status, reviewNote string) (communityPostItem, error) {
	status = strings.TrimSpace(status)
	valid := map[string]bool{"review": true, "published": true, "hidden": true, "rejected": true}
	if !valid[status] {
		return communityPostItem{}, errors.New("invalid status")
	}
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	if _, err := s.db.ExecContext(ctx, `UPDATE community_posts SET status=?, review_note=? WHERE id=?`, status, strings.TrimSpace(reviewNote), id); err != nil {
		return communityPostItem{}, err
	}
	return s.getCommunityPostByID(ctx, id, 0, true)
}

func (s *appServer) listCommunityReports(ctx context.Context, status string) ([]communityReportItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, `
		SELECT
			r.id,r.post_id,p.title,r.reporter_id,r.reason,r.detail,r.status,r.created_at,r.updated_at,
			u.id,u.nickname,u.avatar
		FROM community_post_reports r
		JOIN community_posts p ON p.id=r.post_id
		JOIN users u ON u.id=r.reporter_id
		WHERE (?='' OR r.status=?)
		ORDER BY r.created_at DESC
	`, status, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]communityReportItem, 0)
	for rows.Next() {
		var item communityReportItem
		var created time.Time
		var updated time.Time
		if err := rows.Scan(
			&item.ID,
			&item.PostID,
			&item.PostTitle,
			&item.ReporterID,
			&item.Reason,
			&item.Detail,
			&item.Status,
			&created,
			&updated,
			&item.Reporter.ID,
			&item.Reporter.Nickname,
			&item.Reporter.Avatar,
		); err != nil {
			return nil, err
		}
		item.CreatedAt = created.Format(time.RFC3339)
		item.UpdatedAt = updated.Format(time.RFC3339)
		items = append(items, item)
	}
	return items, nil
}

func (s *appServer) getCommunityReportByID(ctx context.Context, id int64) (communityReportItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	var item communityReportItem
	var created time.Time
	var updated time.Time
	err := s.db.QueryRowContext(ctx, `
		SELECT
			r.id,r.post_id,p.title,r.reporter_id,r.reason,r.detail,r.status,r.created_at,r.updated_at,
			u.id,u.nickname,u.avatar
		FROM community_post_reports r
		JOIN community_posts p ON p.id=r.post_id
		JOIN users u ON u.id=r.reporter_id
		WHERE r.id=?
	`, id).Scan(
		&item.ID,
		&item.PostID,
		&item.PostTitle,
		&item.ReporterID,
		&item.Reason,
		&item.Detail,
		&item.Status,
		&created,
		&updated,
		&item.Reporter.ID,
		&item.Reporter.Nickname,
		&item.Reporter.Avatar,
	)
	if err != nil {
		return communityReportItem{}, err
	}
	item.CreatedAt = created.Format(time.RFC3339)
	item.UpdatedAt = updated.Format(time.RFC3339)
	return item, nil
}

func (s *appServer) updateCommunityReportStatus(ctx context.Context, id int64, status, postStatus, reviewNote string) (communityReportItem, error) {
	status = strings.TrimSpace(status)
	valid := map[string]bool{"open": true, "resolved": true, "dismissed": true}
	if !valid[status] {
		return communityReportItem{}, errors.New("invalid report status")
	}
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	var postID int64
	if err := s.db.QueryRowContext(ctx, `SELECT post_id FROM community_post_reports WHERE id=?`, id).Scan(&postID); err != nil {
		return communityReportItem{}, err
	}
	if _, err := s.db.ExecContext(ctx, `UPDATE community_post_reports SET status=? WHERE id=?`, status, id); err != nil {
		return communityReportItem{}, err
	}
	if postStatus = strings.TrimSpace(postStatus); postStatus != "" {
		if _, err := s.db.ExecContext(ctx, `UPDATE community_posts SET status=?, review_note=? WHERE id=?`, postStatus, strings.TrimSpace(reviewNote), postID); err != nil {
			return communityReportItem{}, err
		}
	}
	if _, err := s.db.ExecContext(ctx, `UPDATE community_posts SET report_count=(SELECT COUNT(1) FROM community_post_reports WHERE post_id=? AND status='open') WHERE id=?`, postID, postID); err != nil {
		return communityReportItem{}, err
	}
	return s.getCommunityReportByID(ctx, id)
}

func (s *appServer) listCommunityCommentReports(ctx context.Context, status string) ([]communityCommentReportItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, `
		SELECT
			r.id,r.comment_id,c.post_id,p.title,c.content,r.reporter_id,r.reason,r.detail,r.status,r.created_at,r.updated_at,
			ru.id,ru.nickname,ru.avatar,
			cu.id,cu.nickname,cu.avatar
		FROM community_comment_reports r
		JOIN community_comments c ON c.id=r.comment_id
		JOIN community_posts p ON p.id=c.post_id
		JOIN users ru ON ru.id=r.reporter_id
		JOIN users cu ON cu.id=c.user_id
		WHERE (?='' OR r.status=?)
		ORDER BY r.created_at DESC
	`, status, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]communityCommentReportItem, 0)
	for rows.Next() {
		var item communityCommentReportItem
		var created time.Time
		var updated time.Time
		if err := rows.Scan(
			&item.ID,
			&item.CommentID,
			&item.PostID,
			&item.PostTitle,
			&item.CommentContent,
			&item.ReporterID,
			&item.Reason,
			&item.Detail,
			&item.Status,
			&created,
			&updated,
			&item.Reporter.ID,
			&item.Reporter.Nickname,
			&item.Reporter.Avatar,
			&item.CommentAuthor.ID,
			&item.CommentAuthor.Nickname,
			&item.CommentAuthor.Avatar,
		); err != nil {
			return nil, err
		}
		item.CreatedAt = created.Format(time.RFC3339)
		item.UpdatedAt = updated.Format(time.RFC3339)
		items = append(items, item)
	}
	return items, nil
}

func (s *appServer) getCommunityCommentReportByID(ctx context.Context, id int64) (communityCommentReportItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	var item communityCommentReportItem
	var created time.Time
	var updated time.Time
	err := s.db.QueryRowContext(ctx, `
		SELECT
			r.id,r.comment_id,c.post_id,p.title,c.content,r.reporter_id,r.reason,r.detail,r.status,r.created_at,r.updated_at,
			ru.id,ru.nickname,ru.avatar,
			cu.id,cu.nickname,cu.avatar
		FROM community_comment_reports r
		JOIN community_comments c ON c.id=r.comment_id
		JOIN community_posts p ON p.id=c.post_id
		JOIN users ru ON ru.id=r.reporter_id
		JOIN users cu ON cu.id=c.user_id
		WHERE r.id=?
	`, id).Scan(
		&item.ID,
		&item.CommentID,
		&item.PostID,
		&item.PostTitle,
		&item.CommentContent,
		&item.ReporterID,
		&item.Reason,
		&item.Detail,
		&item.Status,
		&created,
		&updated,
		&item.Reporter.ID,
		&item.Reporter.Nickname,
		&item.Reporter.Avatar,
		&item.CommentAuthor.ID,
		&item.CommentAuthor.Nickname,
		&item.CommentAuthor.Avatar,
	)
	if err != nil {
		return communityCommentReportItem{}, err
	}
	item.CreatedAt = created.Format(time.RFC3339)
	item.UpdatedAt = updated.Format(time.RFC3339)
	return item, nil
}

func (s *appServer) updateCommunityCommentReportStatus(ctx context.Context, id int64, status, commentStatus, reviewNote string) (communityCommentReportItem, error) {
	status = strings.TrimSpace(status)
	valid := map[string]bool{"open": true, "resolved": true, "dismissed": true}
	if !valid[status] {
		return communityCommentReportItem{}, errors.New("invalid report status")
	}
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	var commentID int64
	if err := s.db.QueryRowContext(ctx, `SELECT comment_id FROM community_comment_reports WHERE id=?`, id).Scan(&commentID); err != nil {
		return communityCommentReportItem{}, err
	}
	if _, err := s.db.ExecContext(ctx, `UPDATE community_comment_reports SET status=? WHERE id=?`, status, id); err != nil {
		return communityCommentReportItem{}, err
	}
	if commentStatus = strings.TrimSpace(commentStatus); commentStatus != "" {
		if _, err := s.db.ExecContext(ctx, `UPDATE community_comments SET status=?, review_note=? WHERE id=?`, commentStatus, strings.TrimSpace(reviewNote), commentID); err != nil {
			return communityCommentReportItem{}, err
		}
		if err := s.refreshCommunityCounts(ctx); err != nil {
			return communityCommentReportItem{}, err
		}
	}
	return s.getCommunityCommentReportByID(ctx, id)
}

func (s *appServer) listCommunityKeywordRules(ctx context.Context) ([]communityKeywordRuleItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, `SELECT id,keyword,action,enabled,created_at,updated_at FROM community_keyword_rules ORDER BY id DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]communityKeywordRuleItem, 0)
	for rows.Next() {
		var item communityKeywordRuleItem
		var enabled int
		var created time.Time
		var updated time.Time
		if err := rows.Scan(&item.ID, &item.Keyword, &item.Action, &enabled, &created, &updated); err != nil {
			return nil, err
		}
		item.Enabled = enabled == 1
		item.CreatedAt = created.Format(time.RFC3339)
		item.UpdatedAt = updated.Format(time.RFC3339)
		items = append(items, item)
	}
	return items, nil
}

func (s *appServer) createCommunityKeywordRule(ctx context.Context, in communityKeywordRuleItem) (communityKeywordRuleItem, error) {
	in.Keyword = strings.TrimSpace(in.Keyword)
	in.Action = strings.TrimSpace(in.Action)
	if in.Keyword == "" {
		return communityKeywordRuleItem{}, errors.New("keyword is required")
	}
	if in.Action != "review" && in.Action != "reject" {
		return communityKeywordRuleItem{}, errors.New("action must be review or reject")
	}
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	res, err := s.db.ExecContext(ctx, `INSERT INTO community_keyword_rules(keyword,action,enabled) VALUES(?,?,?)`, in.Keyword, in.Action, boolToInt(in.Enabled))
	if err != nil {
		if isDuplicateErr(err) {
			return communityKeywordRuleItem{}, errors.New("keyword already exists")
		}
		return communityKeywordRuleItem{}, err
	}
	id, _ := res.LastInsertId()
	return s.getCommunityKeywordRuleByID(ctx, id)
}

func (s *appServer) updateCommunityKeywordRule(ctx context.Context, id int64, in communityKeywordRuleItem) (communityKeywordRuleItem, error) {
	in.Keyword = strings.TrimSpace(in.Keyword)
	in.Action = strings.TrimSpace(in.Action)
	if in.Keyword == "" {
		return communityKeywordRuleItem{}, errors.New("keyword is required")
	}
	if in.Action != "review" && in.Action != "reject" {
		return communityKeywordRuleItem{}, errors.New("action must be review or reject")
	}
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	_, err := s.db.ExecContext(ctx, `UPDATE community_keyword_rules SET keyword=?, action=?, enabled=? WHERE id=?`, in.Keyword, in.Action, boolToInt(in.Enabled), id)
	if err != nil {
		if isDuplicateErr(err) {
			return communityKeywordRuleItem{}, errors.New("keyword already exists")
		}
		return communityKeywordRuleItem{}, err
	}
	return s.getCommunityKeywordRuleByID(ctx, id)
}

func (s *appServer) getCommunityKeywordRuleByID(ctx context.Context, id int64) (communityKeywordRuleItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	var item communityKeywordRuleItem
	var enabled int
	var created time.Time
	var updated time.Time
	err := s.db.QueryRowContext(ctx, `SELECT id,keyword,action,enabled,created_at,updated_at FROM community_keyword_rules WHERE id=?`, id).
		Scan(&item.ID, &item.Keyword, &item.Action, &enabled, &created, &updated)
	if err != nil {
		return communityKeywordRuleItem{}, err
	}
	item.Enabled = enabled == 1
	item.CreatedAt = created.Format(time.RFC3339)
	item.UpdatedAt = updated.Format(time.RFC3339)
	return item, nil
}

func (s *appServer) seedCommunityKeywordRules(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM community_keyword_rules`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	seeds := []communityKeywordRuleItem{
		{Keyword: "加微信", Action: "reject", Enabled: true},
		{Keyword: "vx", Action: "reject", Enabled: true},
		{Keyword: "刷单", Action: "reject", Enabled: true},
		{Keyword: "色情", Action: "reject", Enabled: true},
		{Keyword: "赌博", Action: "reject", Enabled: true},
		{Keyword: "代办", Action: "review", Enabled: true},
	}
	for _, item := range seeds {
		if _, err := s.db.ExecContext(ctx, `INSERT INTO community_keyword_rules(keyword,action,enabled) VALUES(?,?,?)`, item.Keyword, item.Action, boolToInt(item.Enabled)); err != nil {
			return err
		}
	}
	return nil
}

func (s *appServer) listCommunityComments(ctx context.Context, postID, uid int64) ([]communityCommentItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, `
		SELECT c.id,c.post_id,c.user_id,COALESCE(c.parent_id,0),COALESCE(c.root_id,0),COALESCE(pu.nickname,''),c.content,c.image,c.status,c.review_note,c.reply_count,c.created_at,c.updated_at,
			u.id,u.nickname,u.avatar
		FROM community_comments c
		JOIN users u ON u.id=c.user_id
		LEFT JOIN community_comments pc ON pc.id=c.parent_id
		LEFT JOIN users pu ON pu.id=pc.user_id
		WHERE c.post_id=?
		  AND (c.status='published' OR (? > 0 AND c.user_id=?))
		ORDER BY c.created_at ASC, c.id ASC
	`, postID, uid, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	flat := make([]communityCommentItem, 0)
	for rows.Next() {
		item, err := scanCommunityComment(rows)
		if err != nil {
			return nil, err
		}
		flat = append(flat, item)
	}

	top := make([]communityCommentItem, 0)
	indexByID := map[int64]int{}
	for i := range flat {
		flat[i].Replies = []communityCommentItem{}
		if flat[i].ParentID == 0 {
			top = append(top, flat[i])
			indexByID[flat[i].ID] = len(top) - 1
		}
	}
	for _, item := range flat {
		if item.ParentID == 0 {
			continue
		}
		threadID := item.RootID
		if threadID <= 0 {
			threadID = item.ParentID
		}
		parentIndex, ok := indexByID[threadID]
		if !ok {
			continue
		}
		top[parentIndex].Replies = append(top[parentIndex].Replies, item)
		top[parentIndex].ReplyCount = len(top[parentIndex].Replies)
	}
	return top, nil
}

func scanCommunityComment(scanner interface{ Scan(dest ...any) error }) (communityCommentItem, error) {
	var item communityCommentItem
	var created time.Time
	var updated time.Time
	if err := scanner.Scan(
		&item.ID,
		&item.PostID,
		&item.UserID,
		&item.ParentID,
		&item.RootID,
		&item.ReplyToNickname,
		&item.Content,
		&item.Image,
		&item.Status,
		&item.ReviewNote,
		&item.ReplyCount,
		&created,
		&updated,
		&item.Author.ID,
		&item.Author.Nickname,
		&item.Author.Avatar,
	); err != nil {
		return communityCommentItem{}, err
	}
	item.CreatedAt = created.Format(time.RFC3339)
	item.UpdatedAt = updated.Format(time.RFC3339)
	return item, nil
}

func (s *appServer) createCommunityComment(ctx context.Context, postID, uid int64, in communityCreateCommentRequest) (communityCommentItem, string, error) {
	in.Content = strings.TrimSpace(in.Content)
	in.Image = strings.TrimSpace(in.Image)
	if in.Content == "" && in.Image == "" {
		return communityCommentItem{}, "", errors.New("content or image is required")
	}
	if len([]rune(in.Content)) > 1000 {
		return communityCommentItem{}, "", errors.New("content is too long")
	}

	moderation, err := s.moderateCommunityContent(ctx, "", in.Content)
	if err != nil {
		return communityCommentItem{}, "", err
	}
	if moderation.Status == "rejected" {
		return communityCommentItem{}, "", errors.New(emptyFallback(moderation.ReviewNote, "content violates community rules"))
	}

	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	var parentID any = nil
	rootID := int64(0)
	if in.ParentID > 0 {
		parentID = in.ParentID
		var parentRootID sql.NullInt64
		var parentParentID sql.NullInt64
		if err := s.db.QueryRowContext(ctx, `SELECT root_id,parent_id,post_id FROM community_comments WHERE id=?`, in.ParentID).Scan(&parentRootID, &parentParentID, &postID); err != nil {
			return communityCommentItem{}, "", errors.New("parent comment not found")
		}
		if parentRootID.Valid && parentRootID.Int64 > 0 {
			rootID = parentRootID.Int64
		} else if parentParentID.Valid && parentParentID.Int64 > 0 {
			rootID = parentParentID.Int64
		} else {
			rootID = in.ParentID
		}
	}

	res, err := s.db.ExecContext(ctx, `
		INSERT INTO community_comments(post_id,user_id,parent_id,root_id,content,image,status,review_note,reply_count)
		VALUES(?,?,?,?,?,?,?,?,0)
	`, postID, uid, parentID, nullableID(rootID), in.Content, in.Image, moderation.Status, moderation.ReviewNote)
	if err != nil {
		return communityCommentItem{}, "", err
	}
	id, _ := res.LastInsertId()

	if in.ParentID > 0 {
		var total int
		if err := s.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM community_comments WHERE parent_id=? AND status='published'`, in.ParentID).Scan(&total); err != nil {
			return communityCommentItem{}, "", err
		}
		if _, err := s.db.ExecContext(ctx, `UPDATE community_comments SET reply_count=? WHERE id=?`, total, in.ParentID); err != nil {
			return communityCommentItem{}, "", err
		}
	}
	if err := s.refreshCommunityCounts(ctx); err != nil {
		return communityCommentItem{}, "", err
	}

	item, err := s.getCommunityCommentByID(ctx, id)
	if err != nil {
		return communityCommentItem{}, "", err
	}
	if moderation.Status == "review" {
		return item, "submitted_for_review", nil
	}
	return item, "created", nil
}

func (s *appServer) getCommunityCommentByID(ctx context.Context, id int64) (communityCommentItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	row := s.db.QueryRowContext(ctx, `
		SELECT c.id,c.post_id,c.user_id,COALESCE(c.parent_id,0),COALESCE(c.root_id,0),COALESCE(pu.nickname,''),c.content,c.image,c.status,c.review_note,c.reply_count,c.created_at,c.updated_at,
			u.id,u.nickname,u.avatar
		FROM community_comments c
		JOIN users u ON u.id=c.user_id
		LEFT JOIN community_comments pc ON pc.id=c.parent_id
		LEFT JOIN users pu ON pu.id=pc.user_id
		WHERE c.id=?
	`, id)
	return scanCommunityComment(row)
}

func (s *appServer) deleteOwnCommunityComment(ctx context.Context, id, uid int64) error {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	var postID int64
	var parentID sql.NullInt64
	if err := s.db.QueryRowContext(ctx, `SELECT post_id,parent_id FROM community_comments WHERE id=? AND user_id=?`, id, uid).Scan(&postID, &parentID); err != nil {
		return errors.New("comment not found")
	}
	if _, err := s.db.ExecContext(ctx, `DELETE FROM community_comments WHERE id=? AND user_id=?`, id, uid); err != nil {
		return err
	}
	if parentID.Valid {
		var total int
		if err := s.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM community_comments WHERE parent_id=? AND status='published'`, parentID.Int64).Scan(&total); err != nil {
			return err
		}
		if _, err := s.db.ExecContext(ctx, `UPDATE community_comments SET reply_count=? WHERE id=?`, total, parentID.Int64); err != nil {
			return err
		}
	}
	return s.refreshCommunityCounts(ctx)
}

func (s *appServer) refreshCommunityCounts(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx, `
		UPDATE community_posts p
		SET comment_count = (SELECT COUNT(1) FROM community_comments c WHERE c.post_id=p.id AND c.status='published')
	`); err != nil {
		return err
	}
	if _, err := s.db.ExecContext(ctx, `UPDATE community_comments SET reply_count=0`); err != nil {
		return err
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT parent_id, COUNT(1) AS total
		FROM community_comments
		WHERE parent_id IS NOT NULL AND status='published'
		GROUP BY parent_id
	`)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var parentID int64
		var total int
		if err := rows.Scan(&parentID, &total); err != nil {
			return err
		}
		if _, err := s.db.ExecContext(ctx, `UPDATE community_comments SET reply_count=? WHERE id=?`, total, parentID); err != nil {
			return err
		}
	}
	return nil
}

func (s *appServer) listCommunityCommentsForAdmin(ctx context.Context, status string) ([]communityCommentItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, `
		SELECT c.id,c.post_id,c.user_id,COALESCE(c.parent_id,0),COALESCE(c.root_id,0),COALESCE(pu.nickname,''),c.content,c.image,c.status,c.review_note,c.reply_count,c.created_at,c.updated_at,
			u.id,u.nickname,u.avatar
		FROM community_comments c
		JOIN users u ON u.id=c.user_id
		LEFT JOIN community_comments pc ON pc.id=c.parent_id
		LEFT JOIN users pu ON pu.id=pc.user_id
		WHERE (?='' OR c.status=?)
		ORDER BY c.created_at DESC
	`, status, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]communityCommentItem, 0)
	for rows.Next() {
		item, err := scanCommunityComment(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (s *appServer) updateCommunityCommentStatus(ctx context.Context, id int64, status, reviewNote string) (communityCommentItem, error) {
	status = strings.TrimSpace(status)
	valid := map[string]bool{"review": true, "published": true, "hidden": true, "rejected": true}
	if !valid[status] {
		return communityCommentItem{}, errors.New("invalid status")
	}
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	if _, err := s.db.ExecContext(ctx, `UPDATE community_comments SET status=?, review_note=? WHERE id=?`, status, strings.TrimSpace(reviewNote), id); err != nil {
		return communityCommentItem{}, err
	}
	if err := s.refreshCommunityCounts(ctx); err != nil {
		return communityCommentItem{}, err
	}
	return s.getCommunityCommentByID(ctx, id)
}
