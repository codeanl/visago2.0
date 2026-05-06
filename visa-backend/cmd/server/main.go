package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	qstorage "visa-backend/internal/storage"
)

type apiResponse struct {
	Data    any    `json:"data,omitempty"`
	Message string `json:"message"`
}

type appConfig struct {
	Addr                  string
	RootDSN               string
	AppDSN                string
	SkipDBBootstrap       bool
	AdminDefaultUsername  string
	AdminDefaultNickname  string
	AdminDefaultEmail     string
	AdminDefaultPhone     string
	AdminDefaultPassword  string
	DBName                string
	AppUser               string
	AppPass               string
	JWTSecret             string
	AliyunAccessKeyID     string
	AliyunAccessKeySecret string
	AliyunViapiRegion     string
	AliyunViapiEndpoint   string
	QiniuAccessKey        string
	QiniuSecretKey        string
	QiniuBucket           string
	QiniuDomain           string
	QiniuUploadURL        string
	OpenAIAPIKey          string
	OpenAIBaseURL         string
	OpenAITranslateModel  string
}

type appServer struct {
	db        *sql.DB
	jwtSecret []byte
	cfg       appConfig
	storage   qstorage.Uploader
}

type jwtClaims struct {
	UID  int64  `json:"uid"`
	Kind string `json:"kind,omitempty"`
	jwt.RegisteredClaims
}

type adminAccountProfile struct {
	ID          int64  `json:"id"`
	Username    string `json:"username"`
	Nickname    string `json:"nickname"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Status      string `json:"status"`
	LastLoginAt string `json:"lastLoginAt"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type userProfile struct {
	ID         int64          `json:"id"`
	UUID       string         `json:"uuid"`
	Username   string         `json:"username"`
	Name       string         `json:"name"`
	Nickname   string         `json:"nickname"`
	Email      string         `json:"email"`
	Phone      string         `json:"phone"`
	Bio        string         `json:"bio"`
	Gender     string         `json:"gender"`
	Birthday   string         `json:"birthday"`
	Location   string         `json:"location"`
	Avatar     string         `json:"avatar"`
	Cover      string         `json:"cover"`
	Role       string         `json:"role"`
	Status     string         `json:"status"`
	CreatedAt  string         `json:"createdAt"`
	Membership membershipInfo `json:"membership"`
}

type membershipInfo struct {
	HasMembership bool   `json:"hasMembership"`
	PlanKey       string `json:"planKey"`
	PlanName      string `json:"planName"`
	StartedAt     string `json:"startedAt"`
	ExpiresAt     string `json:"expiresAt"`
	Status        string `json:"status"`
}

type countryItem struct {
	ID       int64    `json:"id"`
	Name     string   `json:"name"`
	Code     string   `json:"code"`
	Region   string   `json:"region"`
	Flag     string   `json:"flag"`
	Status   string   `json:"status"`
	Image    string   `json:"image"`
	Note     string   `json:"note"`
	Tags     []string `json:"tags"`
	Keywords []string `json:"keywords"`
}

type freeCountryItem struct {
	ID                   int64    `json:"id"`
	Name                 string   `json:"name"`
	Code                 string   `json:"code"`
	Flag                 string   `json:"flag"`
	Region               string   `json:"region"`
	City                 string   `json:"city"`
	PolicyType           string   `json:"policyType"`
	Stay                 string   `json:"stay"`
	Note                 string   `json:"note"`
	MapX                 float64  `json:"mapX"`
	MapY                 float64  `json:"mapY"`
	Enabled              bool     `json:"enabled"`
	SupportedCountryID   int64    `json:"supportedCountryId"`
	SupportedCountryName string   `json:"supportedCountryName"`
	SupportedVisaID      int64    `json:"supportedVisaId"`
	SupportedVisaName    string   `json:"supportedVisaName"`
	Keywords             []string `json:"keywords"`
}

type visaListItem struct {
	ID             int64  `json:"id"`
	CountryID      int64  `json:"countryId"`
	CountryName    string `json:"countryName"`
	Name           string `json:"name"`
	VisaType       string `json:"visaType"`
	ProcessingTime string `json:"processingTime"`
	Fee            string `json:"fee"`
	Validity       string `json:"validity"`
	Entries        string `json:"entries"`
	Status         string `json:"status"`
	Description    string `json:"description"`
	LongIntro      string `json:"longIntro"`
	Hot            bool   `json:"hot"`
	VisaFree       bool   `json:"visaFree"`
	UpdatedAt      string `json:"updatedAt"`
}

type hotDestinationItem struct {
	CountryID int64  `json:"countryId"`
	VisaID    int64  `json:"visaId"`
	Name      string `json:"name"`
	Flag      string `json:"flag"`
	Image     string `json:"image"`
	Note      string `json:"note"`
	VisaName  string `json:"visaName"`
	Type      string `json:"type"`
	Time      string `json:"time"`
	Price     string `json:"price"`
	Hot       bool   `json:"hot"`
}

type guideItem struct {
	Title       string `json:"title"`
	Description string `json:"desc"`
	Image       string `json:"image"`
	CTA         string `json:"cta"`
	URL         string `json:"url,omitempty"`
}

type taskItem struct {
	ID         int64  `json:"id"`
	TaskKey    string `json:"taskKey"`
	Title      string `json:"title"`
	Icon       string `json:"icon"`
	Status     string `json:"status"`
	StatusText string `json:"statusText"`
	SortOrder  int    `json:"sortOrder"`
}

type stepItem struct {
	ID         int64       `json:"id"`
	StepKey    string      `json:"stepKey"`
	Title      string      `json:"title"`
	SortOrder  int         `json:"sortOrder"`
	Status     string      `json:"status,omitempty"`
	Strategies []string    `json:"strategies"`
	Guides     []guideItem `json:"guides"`
	Materials  []string    `json:"materials"`
	Tasks      []taskItem  `json:"tasks"`
}

type visaDetail struct {
	visaListItem
	Steps []stepItem `json:"steps"`
}

type planSummary struct {
	ID            int64  `json:"id"`
	CountryID     int64  `json:"countryId"`
	CountryName   string `json:"countryName"`
	CountryFlag   string `json:"countryFlag"`
	VisaID        int64  `json:"visaId"`
	VisaTitle     string `json:"visaTitle"`
	Progress      int    `json:"progress"`
	ActiveStepKey string `json:"activeStepKey"`
	Status        string `json:"status"`
	ResultStatus  string `json:"resultStatus"`
	ResultNote    string `json:"resultNote"`
	ResultAt      string `json:"resultAt"`
	CreatedAt     string `json:"createdAt"`
}

type planDetail struct {
	planSummary
	Steps []stepItem `json:"steps"`
	Tips  []string   `json:"tips"`
}

func main() {
	cfg := loadConfig()
	if !cfg.SkipDBBootstrap {
		if err := ensureDatabase(cfg); err != nil {
			log.Fatal(err)
		}
	}

	db, err := sql.Open("mysql", cfg.AppDSN)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal(err)
	}

	s := &appServer{db: db, jwtSecret: []byte(cfg.JWTSecret), cfg: cfg, storage: newObjectStorage(cfg)}
	if err := s.migrate(ctx); err != nil {
		log.Fatal(err)
	}
	if err := s.seed(ctx); err != nil {
		log.Fatal(err)
	}
	if err := s.ensureAdminAccount(ctx); err != nil {
		log.Fatal(err)
	}
	if err := s.syncReferenceVisaData(ctx); err != nil {
		log.Fatal(err)
	}
	if err := s.seedFreeCountries(ctx); err != nil {
		log.Fatal(err)
	}
	if err := s.syncFreeCountryVisaMappings(ctx); err != nil {
		log.Fatal(err)
	}
	if err := s.seedEmbassies(ctx); err != nil {
		log.Fatal(err)
	}
	if err := s.syncEmbassySeeds(ctx); err != nil {
		log.Fatal(err)
	}
	if err := s.seedCommunityKeywordRules(ctx); err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))
	mux.HandleFunc("/api/health", s.handleHealth)
	mux.HandleFunc("/api/auth/register", s.handleAuthRegister)
	mux.HandleFunc("/api/auth/login", s.handleAuthLogin)
	mux.HandleFunc("/api/auth/password/reset", s.handleAuthPasswordReset)
	mux.HandleFunc("/api/admin/login", s.handleAdminLogin)
	mux.HandleFunc("/api/admin/me", s.withAdminAuth(s.handleAdminMe))
	mux.HandleFunc("/api/admin/accounts", s.withAdminAuth(s.handleAdminAccounts))
	mux.HandleFunc("/api/admin/accounts/", s.withAdminAuth(s.handleAdminAccountByID))
	mux.HandleFunc("/api/auth/me", s.withAuth(s.handleAuthMe))
	mux.HandleFunc("/api/auth/me/delete", s.withAuth(s.handleAuthMeDelete))
	mux.HandleFunc("/api/auth/me/profile", s.withAuth(s.handleAuthMeProfile))
	mux.HandleFunc("/api/auth/me/avatar", s.withAuth(s.handleAuthMeAvatarUpload))
	mux.HandleFunc("/api/auth/me/membership/subscribe", s.withAuth(s.handleAuthMeMembershipSubscribe))

	mux.HandleFunc("/api/users", s.withAdminAuth(s.handleUsersAdmin))
	mux.HandleFunc("/api/users/", s.withAdminAuth(s.handleUserRoutesAdmin))
	mux.HandleFunc("/api/uploads/image", s.handleImageUpload)
	mux.HandleFunc("/api/tools/photo-check/quota", s.withAuth(s.handlePhotoCheckQuota))
	mux.HandleFunc("/api/tools/photo-check/analyze", s.withAuth(s.handlePhotoCheckAnalyze))
	mux.HandleFunc("/api/tools/translate/text", s.handleTranslateText)
	mux.HandleFunc("/api/tools/exchange/countries", s.handleExchangeCountries)
	mux.HandleFunc("/api/tools/exchange/quote", s.handleExchangeQuote)
	mux.HandleFunc("/api/tools/exchange/trend", s.handleExchangeTrend)
	mux.HandleFunc("/api/tools/embassies", s.handleEmbassies)
	mux.HandleFunc("/api/tools/embassies/", s.handleEmbassyByID)
	mux.HandleFunc("/api/community/posts", s.handleCommunityPosts)
	mux.HandleFunc("/api/community/posts/", s.handleCommunityPostRoutes)
	mux.HandleFunc("/api/community/favorites", s.withAuth(s.handleCommunityFavorites))
	mux.HandleFunc("/api/community/favorites/", s.withAuth(s.handleCommunityFavoriteByID))
	mux.HandleFunc("/api/community/me/posts", s.withAuth(s.handleCommunityMyPosts))
	mux.HandleFunc("/api/community/comments/", s.handleCommunityCommentByID)
	mux.HandleFunc("/api/community/users/", s.handleCommunityUserRoutes)
	mux.HandleFunc("/api/community/admin/posts", s.withAdminAuth(s.handleCommunityAdminPostsAdmin))
	mux.HandleFunc("/api/community/admin/posts/", s.withAdminAuth(s.handleCommunityAdminPostByIDAdmin))
	mux.HandleFunc("/api/community/admin/comments", s.withAdminAuth(s.handleCommunityAdminCommentsAdmin))
	mux.HandleFunc("/api/community/admin/comments/", s.withAdminAuth(s.handleCommunityAdminCommentByIDAdmin))
	mux.HandleFunc("/api/community/admin/reports", s.withAdminAuth(s.handleCommunityAdminReportsAdmin))
	mux.HandleFunc("/api/community/admin/reports/", s.withAdminAuth(s.handleCommunityAdminReportByIDAdmin))
	mux.HandleFunc("/api/community/admin/comment-reports", s.withAdminAuth(s.handleCommunityAdminCommentReportsAdmin))
	mux.HandleFunc("/api/community/admin/comment-reports/", s.withAdminAuth(s.handleCommunityAdminCommentReportByIDAdmin))
	mux.HandleFunc("/api/community/admin/keywords", s.withAdminAuth(s.handleCommunityAdminKeywordsAdmin))
	mux.HandleFunc("/api/community/admin/keywords/", s.withAdminAuth(s.handleCommunityAdminKeywordByIDAdmin))

	mux.HandleFunc("/api/visa/hot-destinations", s.handleHotDestinations)
	mux.HandleFunc("/api/visa/countries", s.handleCountries)
	mux.HandleFunc("/api/visa/countries/", s.handleCountryRoutes)
	mux.HandleFunc("/api/visa/free-countries", s.handleFreeCountries)
	mux.HandleFunc("/api/visa/free-countries/", s.handleFreeCountryByID)
	mux.HandleFunc("/api/visa/country-visas", s.handleCountryVisas)
	mux.HandleFunc("/api/visa/country-visas/", s.handleCountryVisaByID)
	mux.HandleFunc("/api/visa/visas/", s.handleVisaRoutes)

	mux.HandleFunc("/api/plans", s.withAuth(s.handlePlans))
	mux.HandleFunc("/api/plans/", s.withAuth(s.handlePlanRoutes))

	log.Printf("visa backend listening on http://localhost%s", cfg.Addr)
	if err := http.ListenAndServe(cfg.Addr, withCORS(mux)); err != nil {
		log.Fatal(err)
	}
}

func loadConfig() appConfig {
	cfg := appConfig{
		Addr:                  getEnv("APP_ADDR", ":8080"),
		SkipDBBootstrap:       getEnvBool("SKIP_DB_BOOTSTRAP", false),
		AdminDefaultUsername:  getEnv("ADMIN_DEFAULT_USERNAME", "admin"),
		AdminDefaultNickname:  getEnv("ADMIN_DEFAULT_NICKNAME", "管理员"),
		AdminDefaultEmail:     getEnv("ADMIN_DEFAULT_EMAIL", "admin@visago.com"),
		AdminDefaultPhone:     getEnv("ADMIN_DEFAULT_PHONE", ""),
		AdminDefaultPassword:  getEnv("ADMIN_DEFAULT_PASSWORD", "12345678"),
		DBName:                getEnv("MYSQL_DB", "visago"),
		AppUser:               getEnv("MYSQL_APP_USER", "visago_app"),
		AppPass:               getEnv("MYSQL_APP_PASS", "visago_app_123"),
		JWTSecret:             getEnv("JWT_SECRET", "visago_dev_secret"),
		AliyunAccessKeyID:     getEnv("ALIYUN_ACCESS_KEY_ID", ""),
		AliyunAccessKeySecret: getEnv("ALIYUN_ACCESS_KEY_SECRET", ""),
		AliyunViapiRegion:     getEnv("ALIYUN_VIAPI_REGION", "cn-shanghai"),
		AliyunViapiEndpoint:   getEnv("ALIYUN_VIAPI_ENDPOINT", ""),
		QiniuAccessKey:        getEnv("QINIU_ACCESS_KEY", ""),
		QiniuSecretKey:        getEnv("QINIU_SECRET_KEY", ""),
		QiniuBucket:           getEnv("QINIU_BUCKET", ""),
		QiniuDomain:           getEnv("QINIU_DOMAIN", ""),
		QiniuUploadURL:        getEnv("QINIU_UPLOAD_URL", ""),
		OpenAIAPIKey:          getEnv("OPENAI_API_KEY", ""),
		OpenAIBaseURL:         getEnv("OPENAI_API_BASE_URL", "https://api.openai.com/v1"),
		OpenAITranslateModel:  getEnv("OPENAI_TRANSLATE_MODEL", "gpt-4o-mini"),
	}
	cfg.RootDSN = getEnv("MYSQL_ROOT_DSN", "root:123456@tcp(127.0.0.1:3306)/?charset=utf8mb4&parseTime=true&multiStatements=true")
	cfg.AppDSN = getEnv("MYSQL_DSN", fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?charset=utf8mb4&parseTime=true&loc=Local&multiStatements=true", cfg.AppUser, cfg.AppPass, cfg.DBName))
	return cfg
}

func ensureDatabase(cfg appConfig) error {
	db, err := sql.Open("mysql", cfg.RootDSN)
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return err
	}
	queries := []string{
		fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", cfg.DBName),
		fmt.Sprintf("CREATE USER IF NOT EXISTS '%s'@'%%' IDENTIFIED BY '%s'", cfg.AppUser, cfg.AppPass),
		fmt.Sprintf("GRANT ALL PRIVILEGES ON `%s`.* TO '%s'@'%%'", cfg.DBName, cfg.AppUser),
		"FLUSH PRIVILEGES",
	}
	for _, q := range queries {
		if _, err := db.ExecContext(ctx, q); err != nil {
			return err
		}
	}
	return nil
}

func (s *appServer) migrate(ctx context.Context) error {
	sqls := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			uuid CHAR(36) NOT NULL UNIQUE,
			nickname VARCHAR(64) NOT NULL DEFAULT '',
			email VARCHAR(128) NOT NULL DEFAULT '',
			phone VARCHAR(32) NOT NULL UNIQUE,
			password_hash VARCHAR(255) NOT NULL,
			bio VARCHAR(255) NOT NULL DEFAULT '',
			gender VARCHAR(16) NOT NULL DEFAULT '',
			birthday DATE NULL,
			location VARCHAR(64) NOT NULL DEFAULT '',
			avatar VARCHAR(512) NOT NULL DEFAULT '',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS admin_accounts (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			username VARCHAR(64) NOT NULL UNIQUE,
			nickname VARCHAR(64) NOT NULL DEFAULT '',
			phone VARCHAR(32) NOT NULL DEFAULT '',
			email VARCHAR(128) NOT NULL DEFAULT '',
			password_hash VARCHAR(255) NOT NULL,
			status VARCHAR(32) NOT NULL DEFAULT 'active',
			last_login_at DATETIME NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS user_memberships (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			user_id BIGINT NOT NULL UNIQUE,
			plan_key VARCHAR(32) NOT NULL,
			plan_name VARCHAR(64) NOT NULL,
			started_at DATETIME NOT NULL,
			expires_at DATETIME NOT NULL,
			status VARCHAR(32) NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			CONSTRAINT fk_membership_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS countries (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(64) NOT NULL,
			code VARCHAR(16) NOT NULL UNIQUE,
			region VARCHAR(32) NOT NULL,
			flag VARCHAR(16) NOT NULL DEFAULT '',
			status VARCHAR(32) NOT NULL DEFAULT 'active',
			image VARCHAR(512) NOT NULL DEFAULT '',
			note VARCHAR(255) NOT NULL DEFAULT '',
			tags TEXT NOT NULL,
			keywords TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS visas (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			country_id BIGINT NOT NULL,
			name VARCHAR(128) NOT NULL,
			visa_type VARCHAR(64) NOT NULL,
			processing_time VARCHAR(64) NOT NULL DEFAULT '',
			fee VARCHAR(64) NOT NULL DEFAULT '',
			validity VARCHAR(64) NOT NULL DEFAULT '',
			entries VARCHAR(32) NOT NULL DEFAULT '',
			status VARCHAR(32) NOT NULL DEFAULT 'active',
			description TEXT NOT NULL,
			long_intro TEXT NOT NULL,
			hot TINYINT(1) NOT NULL DEFAULT 0,
			visa_free TINYINT(1) NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			CONSTRAINT fk_visa_country FOREIGN KEY (country_id) REFERENCES countries(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS visa_free_countries (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			name VARCHAR(64) NOT NULL,
			code VARCHAR(16) NOT NULL UNIQUE,
			flag VARCHAR(16) NOT NULL DEFAULT '',
			region VARCHAR(32) NOT NULL DEFAULT '',
			city VARCHAR(64) NOT NULL DEFAULT '',
			policy_type VARCHAR(64) NOT NULL DEFAULT '',
			stay VARCHAR(64) NOT NULL DEFAULT '',
			note VARCHAR(255) NOT NULL DEFAULT '',
			map_x DOUBLE NOT NULL DEFAULT 0,
			map_y DOUBLE NOT NULL DEFAULT 0,
			enabled TINYINT(1) NOT NULL DEFAULT 1,
			supported_country_id BIGINT NULL,
			supported_visa_id BIGINT NULL,
			keywords TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			CONSTRAINT fk_free_country_supported FOREIGN KEY (supported_country_id) REFERENCES countries(id) ON DELETE SET NULL,
			CONSTRAINT fk_free_country_supported_visa FOREIGN KEY (supported_visa_id) REFERENCES visas(id) ON DELETE SET NULL
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS embassies (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			country_name VARCHAR(64) NOT NULL,
			country_code VARCHAR(16) NOT NULL DEFAULT '',
			flag VARCHAR(16) NOT NULL DEFAULT '',
			region VARCHAR(32) NOT NULL DEFAULT '',
			name VARCHAR(128) NOT NULL,
			city VARCHAR(64) NOT NULL DEFAULT '',
			distance VARCHAR(32) NOT NULL DEFAULT '',
			address VARCHAR(255) NOT NULL DEFAULT '',
			phone VARCHAR(64) NOT NULL DEFAULT '',
			hours VARCHAR(128) NOT NULL DEFAULT '',
			services TEXT NOT NULL,
			image VARCHAR(512) NOT NULL DEFAULT '',
			latitude DOUBLE NOT NULL DEFAULT 0,
			longitude DOUBLE NOT NULL DEFAULT 0,
			enabled TINYINT(1) NOT NULL DEFAULT 1,
			keywords TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS community_posts (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			user_id BIGINT NOT NULL,
			category VARCHAR(32) NOT NULL DEFAULT '推荐',
			title VARCHAR(128) NOT NULL,
			content TEXT NOT NULL,
			image VARCHAR(512) NOT NULL DEFAULT '',
			images TEXT NOT NULL,
			status VARCHAR(32) NOT NULL DEFAULT 'review',
			review_note VARCHAR(255) NOT NULL DEFAULT '',
			like_count INT NOT NULL DEFAULT 0,
			report_count INT NOT NULL DEFAULT 0,
			comment_count INT NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			CONSTRAINT fk_community_post_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS community_comments (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			post_id BIGINT NOT NULL,
			user_id BIGINT NOT NULL,
			parent_id BIGINT NULL,
			root_id BIGINT NULL,
			content TEXT NOT NULL,
			image VARCHAR(512) NOT NULL DEFAULT '',
			status VARCHAR(32) NOT NULL DEFAULT 'review',
			review_note VARCHAR(255) NOT NULL DEFAULT '',
			reply_count INT NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			CONSTRAINT fk_community_comment_post FOREIGN KEY (post_id) REFERENCES community_posts(id) ON DELETE CASCADE,
			CONSTRAINT fk_community_comment_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			CONSTRAINT fk_community_comment_parent FOREIGN KEY (parent_id) REFERENCES community_comments(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS community_post_likes (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			post_id BIGINT NOT NULL,
			user_id BIGINT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE KEY uniq_community_post_like (post_id, user_id),
			CONSTRAINT fk_community_like_post FOREIGN KEY (post_id) REFERENCES community_posts(id) ON DELETE CASCADE,
			CONSTRAINT fk_community_like_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS community_post_favorites (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			post_id BIGINT NOT NULL,
			user_id BIGINT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE KEY uniq_community_post_favorite (post_id, user_id),
			CONSTRAINT fk_community_favorite_post FOREIGN KEY (post_id) REFERENCES community_posts(id) ON DELETE CASCADE,
			CONSTRAINT fk_community_favorite_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS community_post_reports (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			post_id BIGINT NOT NULL,
			reporter_id BIGINT NOT NULL,
			reason VARCHAR(64) NOT NULL,
			detail VARCHAR(255) NOT NULL DEFAULT '',
			status VARCHAR(32) NOT NULL DEFAULT 'open',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE KEY uniq_community_post_report (post_id, reporter_id),
			CONSTRAINT fk_community_report_post FOREIGN KEY (post_id) REFERENCES community_posts(id) ON DELETE CASCADE,
			CONSTRAINT fk_community_report_user FOREIGN KEY (reporter_id) REFERENCES users(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS community_comment_reports (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			comment_id BIGINT NOT NULL,
			reporter_id BIGINT NOT NULL,
			reason VARCHAR(64) NOT NULL,
			detail VARCHAR(255) NOT NULL DEFAULT '',
			status VARCHAR(32) NOT NULL DEFAULT 'open',
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE KEY uniq_community_comment_report (comment_id, reporter_id),
			CONSTRAINT fk_community_comment_report_comment FOREIGN KEY (comment_id) REFERENCES community_comments(id) ON DELETE CASCADE,
			CONSTRAINT fk_community_comment_report_user FOREIGN KEY (reporter_id) REFERENCES users(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS community_user_blocks (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			user_id BIGINT NOT NULL,
			blocked_user_id BIGINT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			UNIQUE KEY uniq_community_user_block (user_id, blocked_user_id),
			CONSTRAINT fk_community_block_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			CONSTRAINT fk_community_blocked_user FOREIGN KEY (blocked_user_id) REFERENCES users(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS community_keyword_rules (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			keyword VARCHAR(64) NOT NULL,
			action VARCHAR(16) NOT NULL DEFAULT 'review',
			enabled TINYINT(1) NOT NULL DEFAULT 1,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE KEY uniq_community_keyword (keyword)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS user_photo_check_daily_usage (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			user_id BIGINT NOT NULL,
			usage_date DATE NOT NULL,
			used_count INT NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE KEY uk_user_photo_check_usage_date (user_id, usage_date),
			CONSTRAINT fk_user_photo_check_usage_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS visa_steps (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			visa_id BIGINT NOT NULL,
			step_key VARCHAR(32) NOT NULL,
			title VARCHAR(128) NOT NULL,
			sort_order INT NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE KEY uniq_visa_step(visa_id, step_key),
			CONSTRAINT fk_step_visa FOREIGN KEY (visa_id) REFERENCES visas(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS visa_step_strategies (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			step_id BIGINT NOT NULL,
			content TEXT NOT NULL,
			sort_order INT NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			CONSTRAINT fk_strategy_step FOREIGN KEY (step_id) REFERENCES visa_steps(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS visa_step_guides (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			step_id BIGINT NOT NULL,
			title VARCHAR(128) NOT NULL,
			description TEXT NOT NULL,
			image VARCHAR(512) NOT NULL DEFAULT '',
			cta VARCHAR(128) NOT NULL DEFAULT '',
			link_url VARCHAR(512) NOT NULL DEFAULT '',
			sort_order INT NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			CONSTRAINT fk_guide_step FOREIGN KEY (step_id) REFERENCES visa_steps(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS visa_step_materials (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			step_id BIGINT NOT NULL,
			content VARCHAR(255) NOT NULL,
			sort_order INT NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			CONSTRAINT fk_material_step FOREIGN KEY (step_id) REFERENCES visa_steps(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS visa_step_tasks (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			step_id BIGINT NOT NULL,
			task_key VARCHAR(64) NOT NULL,
			title VARCHAR(128) NOT NULL,
			icon VARCHAR(64) NOT NULL DEFAULT 'task_alt',
			default_status VARCHAR(32) NOT NULL DEFAULT 'todo',
			default_status_text VARCHAR(64) NOT NULL DEFAULT '待处理',
			sort_order INT NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			CONSTRAINT fk_task_step FOREIGN KEY (step_id) REFERENCES visa_steps(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS user_visa_plans (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			user_id BIGINT NOT NULL,
			country_id BIGINT NOT NULL,
			visa_id BIGINT NOT NULL,
			country_name VARCHAR(64) NOT NULL,
			visa_title VARCHAR(128) NOT NULL,
			source VARCHAR(32) NOT NULL DEFAULT 'visa',
			progress INT NOT NULL DEFAULT 0,
			active_step_key VARCHAR(32) NOT NULL DEFAULT 'apply',
			status VARCHAR(32) NOT NULL DEFAULT 'active',
			result_status VARCHAR(32) NOT NULL DEFAULT 'pending',
			result_note VARCHAR(255) NOT NULL DEFAULT '',
			result_at DATETIME NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			CONSTRAINT fk_plan_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
			CONSTRAINT fk_plan_country FOREIGN KEY (country_id) REFERENCES countries(id) ON DELETE RESTRICT,
			CONSTRAINT fk_plan_visa FOREIGN KEY (visa_id) REFERENCES visas(id) ON DELETE RESTRICT
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS user_plan_steps (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			plan_id BIGINT NOT NULL,
			step_key VARCHAR(32) NOT NULL,
			title VARCHAR(128) NOT NULL,
			sort_order INT NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			CONSTRAINT fk_user_plan_step_plan FOREIGN KEY (plan_id) REFERENCES user_visa_plans(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS user_plan_step_strategies (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			plan_step_id BIGINT NOT NULL,
			content TEXT NOT NULL,
			sort_order INT NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			CONSTRAINT fk_user_plan_strategy_step FOREIGN KEY (plan_step_id) REFERENCES user_plan_steps(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS user_plan_step_guides (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			plan_step_id BIGINT NOT NULL,
			title VARCHAR(128) NOT NULL,
			description TEXT NOT NULL,
			image VARCHAR(512) NOT NULL DEFAULT '',
			cta VARCHAR(128) NOT NULL DEFAULT '',
			link_url VARCHAR(512) NOT NULL DEFAULT '',
			sort_order INT NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			CONSTRAINT fk_user_plan_guide_step FOREIGN KEY (plan_step_id) REFERENCES user_plan_steps(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS user_plan_step_materials (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			plan_step_id BIGINT NOT NULL,
			content VARCHAR(255) NOT NULL,
			sort_order INT NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			CONSTRAINT fk_user_plan_material_step FOREIGN KEY (plan_step_id) REFERENCES user_plan_steps(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS user_plan_step_tasks (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			plan_step_id BIGINT NOT NULL,
			task_key VARCHAR(64) NOT NULL,
			title VARCHAR(128) NOT NULL,
			icon VARCHAR(64) NOT NULL DEFAULT 'task_alt',
			status VARCHAR(32) NOT NULL DEFAULT 'todo',
			status_text VARCHAR(64) NOT NULL DEFAULT '待处理',
			sort_order INT NOT NULL DEFAULT 0,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			CONSTRAINT fk_user_plan_task_step FOREIGN KEY (plan_step_id) REFERENCES user_plan_steps(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS user_plan_task_states (
			id BIGINT PRIMARY KEY AUTO_INCREMENT,
			plan_id BIGINT NOT NULL,
			task_id BIGINT NOT NULL,
			status VARCHAR(32) NOT NULL,
			status_text VARCHAR(64) NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			UNIQUE KEY uniq_plan_task(plan_id, task_id),
			CONSTRAINT fk_state_plan FOREIGN KEY (plan_id) REFERENCES user_visa_plans(id) ON DELETE CASCADE,
			CONSTRAINT fk_state_task FOREIGN KEY (task_id) REFERENCES visa_step_tasks(id) ON DELETE CASCADE
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
	}
	for _, q := range sqls {
		if _, err := s.db.ExecContext(ctx, q); err != nil {
			log.Printf("migrate sql failed: %s", q)
			return err
		}
	}
	if err := s.ensureGuideURLColumn(ctx); err != nil {
		return err
	}
	if err := s.ensureVisaFreeColumn(ctx); err != nil {
		return err
	}
	if err := s.ensurePlanResultColumns(ctx); err != nil {
		return err
	}
	if err := s.ensureFreeCountryVisaColumn(ctx); err != nil {
		return err
	}
	if err := s.ensureUserUUIDColumn(ctx); err != nil {
		return err
	}
	if err := s.ensureUserLeanSchema(ctx); err != nil {
		return err
	}
	if err := s.ensureCommunityPostImagesColumn(ctx); err != nil {
		return err
	}
	if err := s.ensureCommunityPostCommentCountColumn(ctx); err != nil {
		return err
	}
	if err := s.ensureCommunityCommentImageColumn(ctx); err != nil {
		return err
	}
	return nil
}

func (s *appServer) ensureGuideURLColumn(ctx context.Context) error {
	rows, err := s.db.QueryContext(ctx, `SHOW COLUMNS FROM visa_step_guides LIKE 'link_url'`)
	if err != nil {
		return err
	}
	defer rows.Close()
	if rows.Next() {
		return nil
	}
	_, err = s.db.ExecContext(ctx, `ALTER TABLE visa_step_guides ADD COLUMN link_url VARCHAR(512) NOT NULL DEFAULT ''`)
	return err
}

func (s *appServer) ensureVisaFreeColumn(ctx context.Context) error {
	exists, err := s.hasColumn(ctx, "visas", "visa_free")
	if err != nil {
		return err
	}
	if exists {
		return nil
	}
	_, err = s.db.ExecContext(ctx, `ALTER TABLE visas ADD COLUMN visa_free TINYINT(1) NOT NULL DEFAULT 0 AFTER hot`)
	return err
}

func (s *appServer) ensurePlanResultColumns(ctx context.Context) error {
	columns := []struct {
		Name string
		SQL  string
	}{
		{"result_status", `ALTER TABLE user_visa_plans ADD COLUMN result_status VARCHAR(32) NOT NULL DEFAULT 'pending' AFTER status`},
		{"result_note", `ALTER TABLE user_visa_plans ADD COLUMN result_note VARCHAR(255) NOT NULL DEFAULT '' AFTER result_status`},
		{"result_at", `ALTER TABLE user_visa_plans ADD COLUMN result_at DATETIME NULL AFTER result_note`},
	}
	for _, column := range columns {
		exists, err := s.hasColumn(ctx, "user_visa_plans", column.Name)
		if err != nil {
			return err
		}
		if exists {
			continue
		}
		if _, err := s.db.ExecContext(ctx, column.SQL); err != nil {
			return err
		}
	}
	return nil
}

func (s *appServer) ensureFreeCountryVisaColumn(ctx context.Context) error {
	exists, err := s.hasColumn(ctx, "visa_free_countries", "supported_visa_id")
	if err != nil {
		return err
	}
	if !exists {
		if _, err := s.db.ExecContext(ctx, `ALTER TABLE visa_free_countries ADD COLUMN supported_visa_id BIGINT NULL AFTER supported_country_id`); err != nil {
			return err
		}
	}
	return nil
}

func (s *appServer) ensureUserUUIDColumn(ctx context.Context) error {
	rows, err := s.db.QueryContext(ctx, `SHOW COLUMNS FROM users LIKE 'uuid'`)
	if err != nil {
		return err
	}
	hasColumn := rows.Next()
	rows.Close()
	if !hasColumn {
		if _, err := s.db.ExecContext(ctx, `ALTER TABLE users ADD COLUMN uuid CHAR(36) NOT NULL DEFAULT '' AFTER id`); err != nil {
			return err
		}
	}

	if _, err := s.db.ExecContext(ctx, `UPDATE users SET uuid=UUID() WHERE TRIM(uuid)=''`); err != nil {
		return err
	}

	idxRows, err := s.db.QueryContext(ctx, `SHOW INDEX FROM users WHERE Key_name='uniq_users_uuid'`)
	if err != nil {
		return err
	}
	hasIndex := idxRows.Next()
	idxRows.Close()
	if !hasIndex {
		if _, err := s.db.ExecContext(ctx, `ALTER TABLE users ADD UNIQUE KEY uniq_users_uuid (uuid)`); err != nil {
			return err
		}
	}
	return nil
}

func (s *appServer) ensureCommunityPostImagesColumn(ctx context.Context) error {
	exists, err := s.hasColumn(ctx, "community_posts", "images")
	if err != nil {
		return err
	}
	if !exists {
		if _, err := s.db.ExecContext(ctx, `ALTER TABLE community_posts ADD COLUMN images TEXT NOT NULL AFTER image`); err != nil {
			return err
		}
	}
	return nil
}

func (s *appServer) ensureCommunityPostCommentCountColumn(ctx context.Context) error {
	exists, err := s.hasColumn(ctx, "community_posts", "comment_count")
	if err != nil {
		return err
	}
	if !exists {
		if _, err := s.db.ExecContext(ctx, `ALTER TABLE community_posts ADD COLUMN comment_count INT NOT NULL DEFAULT 0 AFTER report_count`); err != nil {
			return err
		}
	}
	return nil
}

func (s *appServer) ensureCommunityCommentImageColumn(ctx context.Context) error {
	exists, err := s.hasColumn(ctx, "community_comments", "image")
	if err != nil {
		return err
	}
	if !exists {
		if _, err := s.db.ExecContext(ctx, `ALTER TABLE community_comments ADD COLUMN image VARCHAR(512) NOT NULL DEFAULT '' AFTER content`); err != nil {
			return err
		}
	}
	return nil
}

func (s *appServer) hasColumn(ctx context.Context, tableName, columnName string) (bool, error) {
	safeTable := strings.ReplaceAll(tableName, "`", "")
	safeColumn := strings.ReplaceAll(columnName, "'", "''")
	query := fmt.Sprintf("SHOW COLUMNS FROM `%s` LIKE '%s'", safeTable, safeColumn)
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	return rows.Next(), nil
}

func (s *appServer) ensureUserLeanSchema(ctx context.Context) error {
	hasNickname, err := s.hasColumn(ctx, "users", "nickname")
	if err != nil {
		return err
	}
	if !hasNickname {
		if _, err := s.db.ExecContext(ctx, `ALTER TABLE users ADD COLUMN nickname VARCHAR(64) NOT NULL DEFAULT '' AFTER uuid`); err != nil {
			return err
		}
	}

	hasName, err := s.hasColumn(ctx, "users", "name")
	if err != nil {
		return err
	}
	if hasName {
		if _, err := s.db.ExecContext(ctx, `UPDATE users SET nickname=COALESCE(NULLIF(TRIM(nickname),''), NULLIF(TRIM(name),''), phone) WHERE TRIM(nickname)=''`); err != nil {
			return err
		}
	}

	columnsToDrop := []string{"username", "name", "cover", "role", "status"}
	for _, col := range columnsToDrop {
		exists, err := s.hasColumn(ctx, "users", col)
		if err != nil {
			return err
		}
		if !exists {
			continue
		}
		if _, err := s.db.ExecContext(ctx, fmt.Sprintf("ALTER TABLE users DROP COLUMN %s", col)); err != nil {
			return err
		}
	}
	return nil
}

func (s *appServer) seed(ctx context.Context) error {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM countries`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	hash, _ := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
	users := []struct {
		Nickname string
		Email    string
		Phone    string
	}{
		{"体验用户A", "sample1@visago.com", "13800000011"},
		{"签证顾问", "consultant@visago.com", "13800000012"},
		{"普通用户", "user@visago.com", "13800000013"},
	}
	for _, u := range users {
		userUUID, err := newUUIDString()
		if err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `INSERT INTO users(uuid,nickname,email,phone,password_hash) VALUES(?,?,?,?,?)`,
			userUUID, u.Nickname, u.Email, u.Phone, string(hash)); err != nil {
			return err
		}
	}

	type cseed struct {
		Name, Code, Region, Flag, Image, Note string
		Tags, Keywords                        []string
	}
	seeds := []cseed{
		{"日本", "JP", "亚洲", "🇯🇵", "https://images.pexels.com/photos/918275/pexels-photo-918275.jpeg?auto=compress&cs=tinysrgb&w=1200", "东京-大阪双城热门", []string{"旅游", "商务", "留学"}, []string{"日本", "japan", "东京"}},
		{"韩国", "KR", "亚洲", "🇰🇷", "https://images.pexels.com/photos/237211/pexels-photo-237211.jpeg?auto=compress&cs=tinysrgb&w=1200", "首尔短期出行热门", []string{"旅游", "探亲", "商务"}, []string{"韩国", "korea"}},
		{"新加坡", "SG", "亚洲", "🇸🇬", "https://images.pexels.com/photos/3152126/pexels-photo-3152126.jpeg?auto=compress&cs=tinysrgb&w=1200", "电子化流程较多", []string{"旅游", "商务", "会展"}, []string{"新加坡", "singapore"}},
		{"泰国", "TH", "亚洲", "🇹🇭", "https://images.pexels.com/photos/6440428/pexels-photo-6440428.jpeg?auto=compress&cs=tinysrgb&w=1200", "旅游签常见", []string{"旅游", "探亲"}, []string{"泰国", "thailand"}},
		{"英国", "GB", "欧洲", "🇬🇧", "https://images.pexels.com/photos/460672/pexels-photo-460672.jpeg?auto=compress&cs=tinysrgb&w=1200", "需预约线下采集", []string{"旅游", "工作"}, []string{"英国", "uk"}},
		{"法国", "FR", "欧洲", "🇫🇷", "https://images.pexels.com/photos/699466/pexels-photo-699466.jpeg?auto=compress&cs=tinysrgb&w=1200", "申根热门", []string{"申根", "旅游"}, []string{"法国", "france"}},
		{"美国", "US", "美洲", "🇺🇸", "https://images.pexels.com/photos/356844/pexels-photo-356844.jpeg?auto=compress&cs=tinysrgb&w=1200", "面签准备关键", []string{"旅游", "商务"}, []string{"美国", "usa", "b1", "b2"}},
		{"加拿大", "CA", "美洲", "🇨🇦", "https://images.pexels.com/photos/417173/pexels-photo-417173.jpeg?auto=compress&cs=tinysrgb&w=1200", "探亲与留学需求较高", []string{"旅游", "探亲"}, []string{"加拿大", "canada"}},
		{"澳大利亚", "AU", "大洋洲", "🇦🇺", "https://images.pexels.com/photos/995765/pexels-photo-995765.jpeg?auto=compress&cs=tinysrgb&w=1200", "线上申请流程成熟", []string{"电子签", "访客"}, []string{"澳大利亚", "australia"}},
		{"新西兰", "NZ", "大洋洲", "🇳🇿", "https://images.pexels.com/photos/2092587/pexels-photo-2092587.jpeg?auto=compress&cs=tinysrgb&w=1200", "自然风光线热门", []string{"旅游", "探亲"}, []string{"新西兰", "new zealand"}},
	}
	countryIDs := map[string]int64{}
	for _, c := range seeds {
		res, err := tx.ExecContext(ctx, `INSERT INTO countries(name,code,region,flag,status,image,note,tags,keywords) VALUES(?,?,?,?,?,?,?,?,?)`,
			c.Name, c.Code, c.Region, c.Flag, "active", c.Image, c.Note, strings.Join(c.Tags, ","), strings.Join(c.Keywords, ","))
		if err != nil {
			return err
		}
		id, _ := res.LastInsertId()
		countryIDs[c.Code] = id
	}

	for code, cid := range countryIDs {
		if err := seedVisasForCountry(ctx, tx, code, cid); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *appServer) ensureAdminAccount(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM admin_accounts`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(s.cfg.AdminDefaultPassword)), bcrypt.DefaultCost)
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO admin_accounts(username,nickname,phone,email,password_hash,status)
		VALUES(?,?,?,?,?,'active')
	`,
		strings.TrimSpace(s.cfg.AdminDefaultUsername),
		emptyFallback(strings.TrimSpace(s.cfg.AdminDefaultNickname), "管理员"),
		strings.TrimSpace(s.cfg.AdminDefaultPhone),
		strings.TrimSpace(s.cfg.AdminDefaultEmail),
		string(hash),
	)
	if err != nil && !isDuplicateErr(err) {
		return err
	}
	return nil
}

type seedCountryInfo struct {
	ID    int64
	Code  string
	Name  string
	Image string
}

type guideSource struct {
	InfoURL   string
	ApplyURL  string
	DocsURL   string
	ResultURL string
	Image     string
}

type visaSeed struct {
	Item  visaListItem
	Steps []stepItem
}

func seedVisasForCountry(ctx context.Context, tx *sql.Tx, code string, countryID int64) error {
	country := seedCountryInfo{ID: countryID, Code: code}
	for _, seed := range referenceVisasForCountry(country) {
		res, err := tx.ExecContext(ctx, `INSERT INTO visas(country_id,name,visa_type,processing_time,fee,validity,entries,status,description,long_intro,hot,visa_free) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)`,
			countryID, seed.Item.Name, seed.Item.VisaType, seed.Item.ProcessingTime, seed.Item.Fee, seed.Item.Validity, seed.Item.Entries, "active", seed.Item.Description, seed.Item.LongIntro, boolToInt(seed.Item.Hot), boolToInt(seed.Item.VisaFree))
		if err != nil {
			return err
		}
		visaID, _ := res.LastInsertId()
		if err := replaceVisaStepsTx(ctx, tx, visaID, seed.Steps); err != nil {
			return err
		}
	}
	return nil
}

func (s *appServer) syncReferenceVisaData(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	rows, err := s.db.QueryContext(ctx, `SELECT id,code,name,image FROM countries ORDER BY id`)
	if err != nil {
		return err
	}
	defer rows.Close()

	countries := make([]seedCountryInfo, 0)
	for rows.Next() {
		var country seedCountryInfo
		if err := rows.Scan(&country.ID, &country.Code, &country.Name, &country.Image); err != nil {
			return err
		}
		countries = append(countries, country)
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, country := range countries {
		seeds := referenceVisasForCountry(country)
		if len(seeds) == 0 {
			continue
		}

		existingIDs := make([]int64, 0)
		visaRows, err := tx.QueryContext(ctx, `SELECT id FROM visas WHERE country_id=? ORDER BY id`, country.ID)
		if err != nil {
			return err
		}
		for visaRows.Next() {
			var visaID int64
			if err := visaRows.Scan(&visaID); err != nil {
				visaRows.Close()
				return err
			}
			existingIDs = append(existingIDs, visaID)
		}
		visaRows.Close()

		for idx, seed := range seeds {
			if idx < len(existingIDs) {
				visaID := existingIDs[idx]
				if _, err := tx.ExecContext(ctx, `
					UPDATE visas
					SET name=?,visa_type=?,processing_time=?,fee=?,validity=?,entries=?,status=?,description=?,long_intro=?,hot=?
					WHERE id=?
				`, seed.Item.Name, seed.Item.VisaType, seed.Item.ProcessingTime, seed.Item.Fee, seed.Item.Validity, seed.Item.Entries, "active", seed.Item.Description, seed.Item.LongIntro, boolToInt(seed.Item.Hot), visaID); err != nil {
					return err
				}
				if _, err := tx.ExecContext(ctx, `UPDATE user_visa_plans SET country_name=?,visa_title=? WHERE visa_id=?`, country.Name, seed.Item.Name, visaID); err != nil {
					return err
				}
				if err := replaceVisaStepsTx(ctx, tx, visaID, seed.Steps); err != nil {
					return err
				}
				continue
			}

			res, err := tx.ExecContext(ctx, `INSERT INTO visas(country_id,name,visa_type,processing_time,fee,validity,entries,status,description,long_intro,hot,visa_free) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)`,
				country.ID, seed.Item.Name, seed.Item.VisaType, seed.Item.ProcessingTime, seed.Item.Fee, seed.Item.Validity, seed.Item.Entries, "active", seed.Item.Description, seed.Item.LongIntro, boolToInt(seed.Item.Hot), boolToInt(seed.Item.VisaFree))
			if err != nil {
				return err
			}
			visaID, _ := res.LastInsertId()
			if err := replaceVisaStepsTx(ctx, tx, visaID, seed.Steps); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func replaceVisaStepsTx(ctx context.Context, tx *sql.Tx, visaID int64, steps []stepItem) error {
	oldRows, err := tx.QueryContext(ctx, `SELECT id FROM visa_steps WHERE visa_id=?`, visaID)
	if err != nil {
		return err
	}
	oldIDs := make([]int64, 0)
	for oldRows.Next() {
		var stepID int64
		if err := oldRows.Scan(&stepID); err != nil {
			oldRows.Close()
			return err
		}
		oldIDs = append(oldIDs, stepID)
	}
	oldRows.Close()

	for _, stepID := range oldIDs {
		if _, err := tx.ExecContext(ctx, `DELETE FROM visa_step_strategies WHERE step_id=?`, stepID); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `DELETE FROM visa_step_guides WHERE step_id=?`, stepID); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `DELETE FROM visa_step_materials WHERE step_id=?`, stepID); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `DELETE FROM visa_step_tasks WHERE step_id=?`, stepID); err != nil {
			return err
		}
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM visa_steps WHERE visa_id=?`, visaID); err != nil {
		return err
	}

	for idx, step := range steps {
		sortOrder := step.SortOrder
		if sortOrder == 0 {
			sortOrder = idx + 1
		}
		res, err := tx.ExecContext(ctx, `INSERT INTO visa_steps(visa_id,step_key,title,sort_order) VALUES(?,?,?,?)`, visaID, step.StepKey, step.Title, sortOrder)
		if err != nil {
			return err
		}
		stepID, _ := res.LastInsertId()

		for sIdx, strategy := range step.Strategies {
			if strings.TrimSpace(strategy) == "" {
				continue
			}
			if _, err := tx.ExecContext(ctx, `INSERT INTO visa_step_strategies(step_id,content,sort_order) VALUES(?,?,?)`, stepID, strategy, sIdx+1); err != nil {
				return err
			}
		}

		for gIdx, guide := range step.Guides {
			if strings.TrimSpace(guide.Title) == "" && strings.TrimSpace(guide.Description) == "" {
				continue
			}
			if _, err := tx.ExecContext(ctx, `INSERT INTO visa_step_guides(step_id,title,description,image,cta,link_url,sort_order) VALUES(?,?,?,?,?,?,?)`,
				stepID, emptyFallback(guide.Title, "缁涙崘鐦夐幐鍥у础"), guide.Description, guide.Image, guide.CTA, guide.URL, gIdx+1); err != nil {
				return err
			}
		}

		for mIdx, material := range step.Materials {
			if strings.TrimSpace(material) == "" {
				continue
			}
			if _, err := tx.ExecContext(ctx, `INSERT INTO visa_step_materials(step_id,content,sort_order) VALUES(?,?,?)`, stepID, material, mIdx+1); err != nil {
				return err
			}
		}

		for tIdx, task := range step.Tasks {
			if _, err := tx.ExecContext(ctx, `INSERT INTO visa_step_tasks(step_id,task_key,title,icon,default_status,default_status_text,sort_order) VALUES(?,?,?,?,?,?,?)`,
				stepID, emptyFallback(task.TaskKey, fmt.Sprintf("task-%d", tIdx+1)), task.Title, emptyFallback(task.Icon, "task_alt"), "todo", defaultStatusText("todo"), tIdx+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func referenceVisasForCountry(country seedCountryInfo) []visaSeed {
	source := guideSourceForCountry(country)
	switch country.Code {
	case "JP":
		return []visaSeed{
			buildVisaSeed(country, source, "单次旅游签证", "旅游", "5-7 个工作日", "¥299", "90 天", "单次", true,
				"适用于赴日旅游和探亲，重点准备行程与资金材料。",
				[]string{"护照首页扫描件", "近 6 个月银行流水", "在职/在学证明", "往返机票预订单", "酒店预订单"},
				true),
		}
	case "US":
		return []visaSeed{
			buildVisaSeed(country, source, "B1/B2 访客签证", "访客", "10-15 个工作日", "¥1280", "10 年", "多次", true,
				"适用于赴美旅游、探亲和短期商务访问。",
				[]string{"DS-160 确认页", "护照", "签证照片", "面签预约确认页", "在职/在学证明", "资金证明"},
				true),
		}
	default:
		return []visaSeed{
			buildVisaSeed(country, source, "旅游签证", "旅游", "7-15 个工作日", "¥800", "90 天", "多次", false,
				fmt.Sprintf("适用于前往 %s 的短期旅行。", country.Name),
				[]string{"护照", "照片", "资金证明", "行程证明"},
				true),
		}
	}
}

func buildVisaSeed(country seedCountryInfo, source guideSource, name, visaType, processing, fee, validity, entries string, hot bool, desc string, materials []string, appointment bool) visaSeed {
	if len(materials) == 0 {
		materials = defaultMaterialsForCustomVisa(visaType)
	}
	return visaSeed{
		Item: visaListItem{
			CountryID:      country.ID,
			CountryName:    country.Name,
			Name:           name,
			VisaType:       visaType,
			ProcessingTime: processing,
			Fee:            fee,
			Validity:       validity,
			Entries:        entries,
			Status:         "active",
			Description:    desc,
			LongIntro:      desc,
			Hot:            hot,
		},
		Steps: buildReferenceSteps(country, name, visaType, source, materials, appointment),
	}
}

func buildReferenceSteps(country seedCountryInfo, visaName, visaType string, source guideSource, materials []string, appointment bool) []stepItem {
	if len(materials) == 0 {
		materials = defaultMaterialsForCustomVisa(visaType)
	}

	applyTasks := []taskItem{
		{TaskKey: "profile", Title: "填写申请资料", Icon: "edit_document", Status: "review", StatusText: "审核中", SortOrder: 1},
		{TaskKey: "type-check", Title: "确认签证类型", Icon: "fact_check", Status: "todo", StatusText: "待处理", SortOrder: 2},
	}
	docTasks := []taskItem{
		{TaskKey: "passport", Title: "护照信息核验", Icon: "badge", Status: "todo", StatusText: "待处理", SortOrder: 1},
		{TaskKey: "photo", Title: "证件照片", Icon: "photo_camera", Status: "todo", StatusText: "待处理", SortOrder: 2},
		{TaskKey: "finance", Title: "资金证明", Icon: "account_balance", Status: "todo", StatusText: "待处理", SortOrder: 3},
	}
	bookTasks := []taskItem{
		{TaskKey: "submit", Title: "提交申请", Icon: "send", Status: "todo", StatusText: "待处理", SortOrder: 1},
	}
	if appointment {
		bookTasks = append(bookTasks,
			taskItem{TaskKey: "booking", Title: "预约递签时间", Icon: "event_available", Status: "todo", StatusText: "待处理", SortOrder: 2},
			taskItem{TaskKey: "biometric", Title: "录入生物信息", Icon: "fingerprint", Status: "todo", StatusText: "待处理", SortOrder: 3},
		)
	}
	resultTasks := []taskItem{
		{TaskKey: "waiting", Title: "跟踪审理进度", Icon: "hourglass_top", Status: "todo", StatusText: "待处理", SortOrder: 1},
		{TaskKey: "pickup", Title: "领取护照或下载电子签", Icon: "inventory_2", Status: "todo", StatusText: "待处理", SortOrder: 2},
	}

	return []stepItem{
		{
			StepKey:   "apply",
			Title:     "申请",
			SortOrder: 1,
			Strategies: []string{
				fmt.Sprintf("先在 %s 官方页面确认 %s 的受理条件。", country.Name, visaName),
				"申请表信息请与证件保持一致。",
			},
			Guides: []guideItem{
				{Title: "官方信息", Description: "查看官方签证要求和入口。", Image: source.Image, CTA: "打开官网", URL: source.InfoURL},
			},
			Materials: materials,
			Tasks:     applyTasks,
		},
		{
			StepKey:   "docs",
			Title:     "材料",
			SortOrder: 2,
			Strategies: []string{
				"按清单准备材料，优先保证护照、照片和资金证明完整。",
			},
			Guides: []guideItem{
				{Title: "材料清单", Description: "以官方要求为准。", Image: source.Image, CTA: "查看材料", URL: source.DocsURL},
			},
			Materials: materials,
			Tasks:     docTasks,
		},
		{
			StepKey:   "book",
			Title:     "预约",
			SortOrder: 3,
			Strategies: []string{
				"提交申请后，按要求完成预约或采集流程。",
			},
			Guides: []guideItem{
				{Title: "提交入口", Description: "前往官方入口提交申请。", Image: source.Image, CTA: "开始申请", URL: source.ApplyURL},
			},
			Materials: materials,
			Tasks:     bookTasks,
		},
		{
			StepKey:   "result",
			Title:     "结果",
			SortOrder: 4,
			Strategies: []string{
				"持续跟踪状态更新，及时领取或下载签证结果。",
			},
			Guides: []guideItem{
				{Title: "结果查询", Description: "通过官方渠道查询进度和结果。", Image: source.Image, CTA: "查询结果", URL: source.ResultURL},
			},
			Materials: materials,
			Tasks:     resultTasks,
		},
	}
}
func guideSourceForCountry(country seedCountryInfo) guideSource {
	source := guideSource{Image: country.Image}
	switch country.Code {
	case "JP":
		source.InfoURL = "https://www.mofa.go.jp/j_info/visit/visa/"
		source.ApplyURL = "https://www.evisa.mofa.go.jp/index"
		source.DocsURL = "https://www.mofa.go.jp/j_info/visit/visa/visaonline.html"
		source.ResultURL = "https://www.evisa.mofa.go.jp/index"
	case "KR":
		source.InfoURL = "https://www.visa.go.kr/openMain.do?LANG_TYPE=EN"
		source.ApplyURL = "https://www.visa.go.kr/openMain.do?LANG_TYPE=EN"
		source.DocsURL = "https://www.visa.go.kr/openMain.do?LANG_TYPE=EN"
		source.ResultURL = "https://www.visa.go.kr/openMain.do?LANG_TYPE=EN"
	case "SG":
		source.InfoURL = "https://www.ica.gov.sg/enter-transit-depart/entering-singapore/visa_requirements"
		source.ApplyURL = "https://www.ica.gov.sg/enter-transit-depart/entering-singapore/visa_requirements"
		source.DocsURL = "https://www.ica.gov.sg/enter-transit-depart/entering-singapore/visa_requirements"
		source.ResultURL = "https://www.ica.gov.sg/enter-transit-depart/entering-singapore/visa_requirements"
	case "TH":
		source.InfoURL = "https://www.thaievisa.go.th/"
		source.ApplyURL = "https://www.thaievisa.go.th/"
		source.DocsURL = "https://www.thaievisa.go.th/Content/Manual.pdf"
		source.ResultURL = "https://www.thaievisa.go.th/"
	case "AE":
		source.InfoURL = "https://u.ae/en/information-and-services/visa-and-emirates-id"
		source.ApplyURL = "https://u.ae/en/information-and-services/visa-and-emirates-id"
		source.DocsURL = "https://u.ae/en/information-and-services/visa-and-emirates-id"
		source.ResultURL = "https://u.ae/en/information-and-services/visa-and-emirates-id"
	case "GB":
		source.InfoURL = "https://www.gov.uk/standard-visitor"
		source.ApplyURL = "https://www.gov.uk/standard-visitor/apply-standard-visitor-visa"
		source.DocsURL = "https://www.gov.uk/standard-visitor/documents-you-must-provide"
		source.ResultURL = "https://www.gov.uk/standard-visitor"
	case "FR":
		source.InfoURL = "https://france-visas.gouv.fr/en/web/france-visas/tourism-private-stay"
		source.ApplyURL = "https://france-visas.gouv.fr/en/web/france-visas/home"
		source.DocsURL = "https://france-visas.gouv.fr/en/web/france-visas/tourism-private-stay"
		source.ResultURL = "https://france-visas.gouv.fr/en/web/france-visas/home"
	case "IT":
		source.InfoURL = "https://vistoperitalia.esteri.it/home/en"
		source.ApplyURL = "https://vistoperitalia.esteri.it/home/en"
		source.DocsURL = "https://vistoperitalia.esteri.it/home/en"
		source.ResultURL = "https://vistoperitalia.esteri.it/home/en"
	case "DE":
		source.InfoURL = "https://www.germany.info/us-en/service/visa"
		source.ApplyURL = "https://www.germany.info/us-en/service/visa"
		source.DocsURL = "https://www.germany.info/us-en/service/visa"
		source.ResultURL = "https://www.germany.info/us-en/service/visa"
	case "ES":
		source.InfoURL = "https://www.exteriores.gob.es/Consulados/chicago/en/ServiciosConsulares/Paginas/Consular/Visados-Schengen.aspx"
		source.ApplyURL = "https://www.exteriores.gob.es/Consulados/chicago/en/ServiciosConsulares/Paginas/Consular/Visados-Schengen.aspx"
		source.DocsURL = "https://www.exteriores.gob.es/Consulados/chicago/en/ServiciosConsulares/Paginas/Consular/Visados-Schengen.aspx"
		source.ResultURL = "https://www.exteriores.gob.es/Consulados/chicago/en/ServiciosConsulares/Paginas/Consular/Visados-Schengen.aspx"
	case "NL":
		source.InfoURL = "https://www.netherlandsworldwide.nl/visa-the-netherlands/schengen-visa"
		source.ApplyURL = "https://www.netherlandsworldwide.nl/visa-the-netherlands/schengen-visa"
		source.DocsURL = "https://www.netherlandsworldwide.nl/visa-the-netherlands/schengen-visa"
		source.ResultURL = "https://www.netherlandsworldwide.nl/visa-the-netherlands/schengen-visa"
	case "US":
		source.InfoURL = "https://travel.state.gov/content/travel/en/us-visas/tourism-visit/visitor.html"
		source.ApplyURL = "https://travel.state.gov/content/travel/en/us-visas/tourism-visit/visitor.html"
		source.DocsURL = "https://travel.state.gov/content/travel/en/us-visas/tourism-visit/visitor.html"
		source.ResultURL = "https://travel.state.gov/content/travel/en/us-visas/tourism-visit/visitor.html"
	case "CA":
		source.InfoURL = "https://www.canada.ca/en/immigration-refugees-citizenship/services/visit-canada/visitor-visa.html"
		source.ApplyURL = "https://www.canada.ca/en/immigration-refugees-citizenship/services/visit-canada/steps-apply-visitor-visa.html"
		source.DocsURL = "https://www.canada.ca/en/immigration-refugees-citizenship/services/visit-canada/visitor-visa.html"
		source.ResultURL = "https://www.canada.ca/en/immigration-refugees-citizenship/services/visit-canada/visitor-visa.html"
	case "BR":
		source.InfoURL = "https://www.gov.br/pt-br/servicos/obter-visto-para-viajar-ao-brasil"
		source.ApplyURL = "https://www.gov.br/pt-br/servicos/obter-visto-para-viajar-ao-brasil"
		source.DocsURL = "https://www.gov.br/pt-br/servicos/obter-visto-para-viajar-ao-brasil"
		source.ResultURL = "https://www.gov.br/pt-br/servicos/obter-visto-para-viajar-ao-brasil"
	case "AU":
		source.InfoURL = "https://immi.homeaffairs.gov.au/visas/getting-a-visa/visa-listing/visitor-600"
		source.ApplyURL = "https://immi.homeaffairs.gov.au/visas/getting-a-visa/visa-listing/visitor-600"
		source.DocsURL = "https://immi.homeaffairs.gov.au/visas/getting-a-visa/visa-listing/visitor-600"
		source.ResultURL = "https://immi.homeaffairs.gov.au/visas/getting-a-visa/visa-listing/visitor-600"
	case "NZ":
		source.InfoURL = "https://www.immigration.govt.nz/new-zealand-visas/visas/visa/visitor-visa"
		source.ApplyURL = "https://www.immigration.govt.nz/new-zealand-visas/visas/visa/visitor-visa"
		source.DocsURL = "https://www.immigration.govt.nz/new-zealand-visas/visas/visa/visitor-visa"
		source.ResultURL = "https://www.immigration.govt.nz/new-zealand-visas/visas/visa/visitor-visa"
	case "ZA":
		source.InfoURL = "https://www.gov.za/services/temporary-residence/visa"
		source.ApplyURL = "https://www.gov.za/services/temporary-residence/visa"
		source.DocsURL = "https://www.gov.za/services/temporary-residence/visa"
		source.ResultURL = "https://www.gov.za/services/temporary-residence/visa"
	default:
		source.InfoURL = "https://www.iatatravelcentre.com/"
		source.ApplyURL = "https://www.iatatravelcentre.com/"
		source.DocsURL = "https://www.iatatravelcentre.com/"
		source.ResultURL = "https://www.iatatravelcentre.com/"
	}
	return source
}

func requiresAppointmentForCountry(code string) bool {
	switch code {
	case "JP", "KR", "GB", "FR", "IT", "DE", "ES", "NL", "US", "BR", "ZA":
		return true
	default:
		return false
	}
}

func defaultMaterialsForCustomVisa(visaType string) []string {
	switch strings.TrimSpace(visaType) {
	case "商务":
		return []string{"护照首页扫描件", "近期证件照", "邀请函", "在职证明", "资金证明"}
	case "申根":
		return []string{"护照原件及旧护照", "2 张近照", "申根旅行保险", "往返机票预订单", "酒店订单与行程单", "近 3-6 个月银行流水"}
	case "访客":
		return []string{"护照扫描件", "近期证件照", "在职/在学证明", "近 6 个月银行流水", "旅行计划说明"}
	default:
		return []string{"护照首页扫描件", "近期证件照", "往返机票预订单", "酒店订单或邀请函", "资金证明"}
	}
}
func (s *appServer) withAuth(next func(http.ResponseWriter, *http.Request, int64)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid, err := s.userIDFromHeader(r)
		if err != nil {
			writeError(w, http.StatusUnauthorized, err)
			return
		}
		next(w, r, uid)
	}
}

func (s *appServer) withAdminAuth(next func(http.ResponseWriter, *http.Request, int64)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid, err := s.adminUserIDFromHeader(r)
		if err != nil {
			status := http.StatusUnauthorized
			if err.Error() == "admin access required" {
				status = http.StatusForbidden
			}
			writeError(w, status, err)
			return
		}
		next(w, r, uid)
	}
}

func (s *appServer) parseTokenClaims(r *http.Request) (*jwtClaims, error) {
	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	if auth == "" {
		return nil, errors.New("missing authorization")
	}
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
		return nil, errors.New("invalid authorization")
	}
	tokenStr := strings.TrimSpace(parts[1])
	tk, err := jwt.ParseWithClaims(tokenStr, &jwtClaims{}, func(token *jwt.Token) (any, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("invalid token method")
		}
		return s.jwtSecret, nil
	})
	if err != nil || !tk.Valid {
		return nil, errors.New("invalid token")
	}
	c, ok := tk.Claims.(*jwtClaims)
	if !ok || c.UID <= 0 {
		return nil, errors.New("invalid token claims")
	}
	return c, nil
}

func (s *appServer) userIDFromHeader(r *http.Request) (int64, error) {
	c, err := s.parseTokenClaims(r)
	if err != nil {
		return 0, err
	}
	if strings.EqualFold(strings.TrimSpace(c.Kind), "admin") {
		return 0, errors.New("invalid token claims")
	}
	return c.UID, nil
}

func (s *appServer) adminUserIDFromHeader(r *http.Request) (int64, error) {
	c, err := s.parseTokenClaims(r)
	if err != nil {
		return 0, err
	}
	if !strings.EqualFold(strings.TrimSpace(c.Kind), "admin") {
		return 0, errors.New("admin access required")
	}
	return c.UID, nil
}

func (s *appServer) requireAdminForWrite(w http.ResponseWriter, r *http.Request) bool {
	switch r.Method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	}
	_, err := s.adminUserIDFromHeader(r)
	if err != nil {
		status := http.StatusUnauthorized
		if err.Error() == "admin access required" {
			status = http.StatusForbidden
		}
		writeError(w, status, err)
		return false
	}
	return true
}

func (s *appServer) issueToken(userID int64) (string, error) {
	now := time.Now()
	claims := jwtClaims{
		UID:  userID,
		Kind: "user",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(30 * 24 * time.Hour)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
}

func (s *appServer) issueAdminToken(adminID int64) (string, error) {
	now := time.Now()
	claims := jwtClaims{
		UID:  adminID,
		Kind: "admin",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(adminID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(30 * 24 * time.Hour)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.jwtSecret)
}

func (s *appServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: map[string]string{"status": "up"}})
}

func (s *appServer) handleAuthRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	var in struct {
		Name     string `json:"name"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
	if err := readJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.Phone = strings.TrimSpace(in.Phone)
	if in.Name == "" || in.Phone == "" || len(strings.TrimSpace(in.Password)) < 8 {
		writeError(w, http.StatusBadRequest, errors.New("name, phone and password(>=8) are required"))
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	userUUID, err := newUUIDString()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()
	_, err = s.db.ExecContext(ctx, `INSERT INTO users(uuid,nickname,phone,password_hash) VALUES(?,?,?,?)`,
		userUUID, in.Name, in.Phone, string(hash))
	if err != nil {
		if isDuplicateErr(err) {
			writeError(w, http.StatusConflict, errors.New("phone already exists"))
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	u, err := s.getUserByPhone(ctx, in.Phone)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	p, err := s.getUserProfile(ctx, u.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusCreated, apiResponse{Message: "created", Data: p})
}

func (s *appServer) handleAuthLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	var in struct {
		Phone    string `json:"phone"`
		Account  string `json:"account"`
		Password string `json:"password"`
	}
	if err := readJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	in.Phone = strings.TrimSpace(in.Phone)
	in.Account = strings.TrimSpace(in.Account)
	loginPhone := in.Phone
	if loginPhone == "" {
		// Backward compatibility: older clients still send `account`.
		loginPhone = in.Account
	}
	if loginPhone == "" || in.Password == "" {
		writeError(w, http.StatusBadRequest, errors.New("phone and password are required"))
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()
	u, err := s.getUserByPhone(ctx, loginPhone)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusUnauthorized, errors.New("\u624b\u673a\u53f7\u4e0d\u5b58\u5728"))
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)) != nil {
		writeError(w, http.StatusUnauthorized, errors.New("\u5bc6\u7801\u9519\u8bef"))
		return
	}
	token, err := s.issueToken(u.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	p, err := s.getUserProfile(ctx, u.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: map[string]any{"token": token, "user": p}})
}

func (s *appServer) handleAdminLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	var in struct {
		Phone    string `json:"phone"`
		Account  string `json:"account"`
		Password string `json:"password"`
	}
	if err := readJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	loginAccount := strings.TrimSpace(in.Phone)
	if loginAccount == "" {
		loginAccount = strings.TrimSpace(in.Account)
	}
	if loginAccount == "" || strings.TrimSpace(in.Password) == "" {
		writeError(w, http.StatusBadRequest, errors.New("account and password are required"))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()

	var u userRow
	var status string
	err := s.db.QueryRowContext(ctx, `
		SELECT id,password_hash,phone,email,status
		FROM admin_accounts
		WHERE username=? OR phone=? OR email=?
		LIMIT 1
	`, loginAccount, loginAccount, loginAccount).
		Scan(&u.ID, &u.PasswordHash, &u.Phone, &u.Email, &status)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusUnauthorized, errors.New("管理员账号不存在"))
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if strings.TrimSpace(status) != "active" {
		writeError(w, http.StatusForbidden, errors.New("管理员账号已停用"))
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(in.Password)) != nil {
		writeError(w, http.StatusUnauthorized, errors.New("密码错误"))
		return
	}
	if _, err := s.db.ExecContext(ctx, `UPDATE admin_accounts SET last_login_at=? WHERE id=?`, time.Now(), u.ID); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	token, err := s.issueAdminToken(u.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	p, err := s.getAdminAccountByID(ctx, u.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: map[string]any{"token": token, "user": p}})
}

func (s *appServer) handleAdminMe(w http.ResponseWriter, r *http.Request, uid int64) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()
	p, err := s.getAdminAccountByID(ctx, uid)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: p})
}

func (s *appServer) handleAdminAccounts(w http.ResponseWriter, r *http.Request, _ int64) {
	switch r.Method {
	case http.MethodGet:
		q := strings.TrimSpace(r.URL.Query().Get("q"))
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		rows, err := s.db.QueryContext(ctx, `
			SELECT id,username,nickname,phone,email,status,last_login_at,created_at,updated_at
			FROM admin_accounts
			WHERE ?='' OR username LIKE CONCAT('%',?,'%') OR nickname LIKE CONCAT('%',?,'%') OR phone LIKE CONCAT('%',?,'%') OR email LIKE CONCAT('%',?,'%')
			ORDER BY id DESC
		`, q, q, q, q, q)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		defer rows.Close()
		out := make([]adminAccountProfile, 0)
		for rows.Next() {
			item, err := scanAdminAccount(rows)
			if err != nil {
				writeError(w, http.StatusInternalServerError, err)
				return
			}
			out = append(out, item)
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: out})
	case http.MethodPost:
		var in struct {
			Username string `json:"username"`
			Nickname string `json:"nickname"`
			Phone    string `json:"phone"`
			Email    string `json:"email"`
			Password string `json:"password"`
			Status   string `json:"status"`
		}
		if err := readJSON(r, &in); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if strings.TrimSpace(in.Username) == "" || strings.TrimSpace(in.Password) == "" {
			writeError(w, http.StatusBadRequest, errors.New("username and password are required"))
			return
		}
		if len(strings.TrimSpace(in.Password)) < 8 {
			writeError(w, http.StatusBadRequest, errors.New("password must be at least 8 chars"))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		if err := s.ensureAdminAccountUnique(ctx, 0, in.Username, in.Phone, in.Email); err != nil {
			writeError(w, http.StatusConflict, err)
			return
		}
		hash, _ := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(in.Password)), bcrypt.DefaultCost)
		res, err := s.db.ExecContext(ctx, `
			INSERT INTO admin_accounts(username,nickname,phone,email,password_hash,status)
			VALUES(?,?,?,?,?,?)
		`,
			strings.TrimSpace(in.Username),
			emptyFallback(strings.TrimSpace(in.Nickname), strings.TrimSpace(in.Username)),
			strings.TrimSpace(in.Phone),
			strings.TrimSpace(in.Email),
			string(hash),
			emptyFallback(strings.TrimSpace(in.Status), "active"),
		)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		id, _ := res.LastInsertId()
		item, err := s.getAdminAccountByID(ctx, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusCreated, apiResponse{Message: "created", Data: item})
	default:
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}

func (s *appServer) handleAdminAccountByID(w http.ResponseWriter, r *http.Request, adminID int64) {
	id, ok := parseID(w, r.URL.Path, "/api/admin/accounts/")
	if !ok {
		return
	}
	switch r.Method {
	case http.MethodGet:
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		item, err := s.getAdminAccountByID(ctx, id)
		if err != nil {
			writeError(w, http.StatusNotFound, errors.New("admin account not found"))
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: item})
	case http.MethodPut:
		var in struct {
			Username string `json:"username"`
			Nickname string `json:"nickname"`
			Phone    string `json:"phone"`
			Email    string `json:"email"`
			Password string `json:"password"`
			Status   string `json:"status"`
		}
		if err := readJSON(r, &in); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		if err := s.ensureAdminAccountUnique(ctx, id, in.Username, in.Phone, in.Email); err != nil {
			writeError(w, http.StatusConflict, err)
			return
		}
		query := `UPDATE admin_accounts SET username=?,nickname=?,phone=?,email=?,status=?`
		args := []any{
			strings.TrimSpace(in.Username),
			emptyFallback(strings.TrimSpace(in.Nickname), strings.TrimSpace(in.Username)),
			strings.TrimSpace(in.Phone),
			strings.TrimSpace(in.Email),
			emptyFallback(strings.TrimSpace(in.Status), "active"),
		}
		if strings.TrimSpace(in.Password) != "" {
			if len(strings.TrimSpace(in.Password)) < 8 {
				writeError(w, http.StatusBadRequest, errors.New("password must be at least 8 chars"))
				return
			}
			hash, _ := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(in.Password)), bcrypt.DefaultCost)
			query += `,password_hash=?`
			args = append(args, string(hash))
		}
		query += ` WHERE id=?`
		args = append(args, id)
		if _, err := s.db.ExecContext(ctx, query, args...); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		item, err := s.getAdminAccountByID(ctx, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "updated", Data: item})
	case http.MethodDelete:
		if id == adminID {
			writeError(w, http.StatusBadRequest, errors.New("cannot delete current admin account"))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		if _, err := s.db.ExecContext(ctx, `DELETE FROM admin_accounts WHERE id=?`, id); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "deleted"})
	default:
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}

func (s *appServer) handleAuthPasswordReset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	var in struct {
		Name     string `json:"name"`
		Phone    string `json:"phone"`
		Password string `json:"password"`
	}
	if err := readJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	in.Name = strings.TrimSpace(in.Name)
	in.Phone = strings.TrimSpace(in.Phone)
	in.Password = strings.TrimSpace(in.Password)
	if in.Name == "" || in.Phone == "" || len(in.Password) < 8 {
		writeError(w, http.StatusBadRequest, errors.New("name, phone and password(>=8) are required"))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()

	var userID int64
	err := s.db.QueryRowContext(ctx, `SELECT id FROM users WHERE phone=? AND nickname=?`, in.Phone, in.Name).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusBadRequest, errors.New("手机号与姓名/昵称不匹配"))
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if _, err := s.db.ExecContext(ctx, `UPDATE users SET password_hash=? WHERE id=?`, string(hash), userID); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, apiResponse{Message: "password reset"})
}

func (s *appServer) handleAuthMe(w http.ResponseWriter, r *http.Request, uid int64) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()
	p, err := s.getUserProfile(ctx, uid)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: p})
}

func (s *appServer) handleAuthMeDelete(w http.ResponseWriter, r *http.Request, uid int64) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	var in struct {
		Password string `json:"password"`
		Confirm  bool   `json:"confirm"`
	}
	if err := readJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	in.Password = strings.TrimSpace(in.Password)
	if !in.Confirm {
		writeError(w, http.StatusBadRequest, errors.New("confirmation is required"))
		return
	}
	if in.Password == "" {
		writeError(w, http.StatusBadRequest, errors.New("password is required"))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()

	var passwordHash string
	if err := s.db.QueryRowContext(ctx, `SELECT password_hash FROM users WHERE id=?`, uid).Scan(&passwordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusNotFound, errors.New("user not found"))
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(in.Password)) != nil {
		writeError(w, http.StatusBadRequest, errors.New("当前密码错误"))
		return
	}

	if _, err := s.db.ExecContext(ctx, `DELETE FROM users WHERE id=?`, uid); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{
		Message: "account deleted",
		Data: map[string]any{
			"communityPolicy": "你的社区帖子、评论、收藏、举报与屏蔽记录会随账号一起删除，且无法恢复。",
		},
	})
}

func (s *appServer) handleAuthMeProfile(w http.ResponseWriter, r *http.Request, uid int64) {
	if r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	var in struct {
		Name     string `json:"name"`
		Nickname string `json:"nickname"`
		Email    string `json:"email"`
		Phone    string `json:"phone"`
		Bio      string `json:"bio"`
		Gender   string `json:"gender"`
		Location string `json:"location"`
		Avatar   string `json:"avatar"`
		Password string `json:"password"`
	}
	if err := readJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()
	current, err := s.getUserProfile(ctx, uid)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	nickname := emptyFallback(strings.TrimSpace(in.Nickname), emptyFallback(strings.TrimSpace(in.Name), "\u7528\u6237"))
	err = s.updateUserCore(ctx, uid, userProfile{
		Nickname: nickname,
		Email:    strings.TrimSpace(in.Email),
		Phone:    strings.TrimSpace(current.Phone),
		Bio:      strings.TrimSpace(in.Bio),
		Gender:   strings.TrimSpace(in.Gender),
		Location: strings.TrimSpace(in.Location),
		Avatar:   strings.TrimSpace(in.Avatar),
	}, strings.TrimSpace(in.Password))
	if err != nil {
		if strings.Contains(err.Error(), "password must be at least 8 chars") {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if isDuplicateErr(err) {
			writeError(w, http.StatusConflict, errors.New("phone already exists"))
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	p, err := s.getUserProfile(ctx, uid)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Message: "updated", Data: p})
}

func (s *appServer) handleAuthMeAvatarUpload(w http.ResponseWriter, r *http.Request, uid int64) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	s.handleUserAvatarUpload(w, r, uid)
}

func (s *appServer) handleUserAvatarUpload(w http.ResponseWriter, r *http.Request, uid int64) {
	if err := r.ParseMultipartForm(12 << 20); err != nil {
		writeError(w, http.StatusBadRequest, errors.New("invalid upload payload"))
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, errors.New("file is required"))
		return
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
	defer cancel()
	objectInfo, err := s.uploadToObjectStorage(ctx, "avatars", fmt.Sprintf("user-%d", uid), header.Filename, file)
	if err != nil {
		mapUploadError(w, err)
		return
	}
	avatarURL := objectInfo.URL

	ctx, cancel = context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()
	if _, err := s.db.ExecContext(ctx, `UPDATE users SET avatar=? WHERE id=?`, avatarURL, uid); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	p, err := s.getUserProfile(ctx, uid)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Message: "uploaded", Data: map[string]any{
		"url":     avatarURL,
		"path":    objectInfo.Key,
		"profile": p,
	}})
}

func (s *appServer) handleImageUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	if err := r.ParseMultipartForm(12 << 20); err != nil {
		writeError(w, http.StatusBadRequest, errors.New("invalid upload payload"))
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, errors.New("file is required"))
		return
	}
	defer file.Close()

	folder := strings.TrimSpace(r.URL.Query().Get("folder"))
	if folder == "" {
		folder = "guides"
	}
	ctx, cancel := context.WithTimeout(r.Context(), 45*time.Second)
	defer cancel()
	objectInfo, err := s.uploadToObjectStorage(ctx, folder, "img", header.Filename, file)
	if err != nil {
		mapUploadError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Message: "uploaded", Data: map[string]any{
		"url":  objectInfo.URL,
		"path": objectInfo.Key,
	}})
}

func (s *appServer) handleAuthMeMembershipSubscribe(w http.ResponseWriter, r *http.Request, uid int64) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	var in struct {
		PlanKey string `json:"planKey"`
	}
	if err := readJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	planName, months, ok := membershipPlan(strings.TrimSpace(in.PlanKey))
	if !ok {
		writeError(w, http.StatusBadRequest, errors.New("invalid planKey"))
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()
	var oldExpire sql.NullTime
	_ = s.db.QueryRowContext(ctx, `SELECT expires_at FROM user_memberships WHERE user_id=?`, uid).Scan(&oldExpire)
	now := time.Now()
	start := now
	if oldExpire.Valid && oldExpire.Time.After(now) {
		start = oldExpire.Time
	}
	expire := start.AddDate(0, months, 0)
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO user_memberships(user_id,plan_key,plan_name,started_at,expires_at,status)
		VALUES(?,?,?,?,?,?)
		ON DUPLICATE KEY UPDATE
			plan_key=VALUES(plan_key),
			plan_name=VALUES(plan_name),
			started_at=VALUES(started_at),
			expires_at=VALUES(expires_at),
			status=VALUES(status)
	`, uid, in.PlanKey, planName, start, expire, "active")
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	p, err := s.getUserProfile(ctx, uid)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Message: "subscribed", Data: p})
}

type userRow struct {
	ID           int64
	PasswordHash string
	Phone        string
	Email        string
}

func scanAdminAccount(scanner interface{ Scan(dest ...any) error }) (adminAccountProfile, error) {
	var item adminAccountProfile
	var lastLogin sql.NullTime
	var created time.Time
	var updated time.Time
	if err := scanner.Scan(
		&item.ID,
		&item.Username,
		&item.Nickname,
		&item.Phone,
		&item.Email,
		&item.Status,
		&lastLogin,
		&created,
		&updated,
	); err != nil {
		return adminAccountProfile{}, err
	}
	if lastLogin.Valid {
		item.LastLoginAt = lastLogin.Time.Format(time.RFC3339)
	}
	item.CreatedAt = created.Format(time.RFC3339)
	item.UpdatedAt = updated.Format(time.RFC3339)
	return item, nil
}

func (s *appServer) getAdminAccountByID(ctx context.Context, id int64) (adminAccountProfile, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id,username,nickname,phone,email,status,last_login_at,created_at,updated_at
		FROM admin_accounts
		WHERE id=?
	`, id)
	return scanAdminAccount(row)
}

func (s *appServer) ensureAdminAccountUnique(ctx context.Context, currentID int64, username, phone, email string) error {
	username = strings.TrimSpace(username)
	phone = strings.TrimSpace(phone)
	email = strings.TrimSpace(email)
	if username == "" {
		return errors.New("username is required")
	}
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM admin_accounts WHERE username=? AND id<>?`, username, currentID).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return errors.New("username already exists")
	}
	if phone != "" {
		if err := s.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM admin_accounts WHERE phone=? AND id<>?`, phone, currentID).Scan(&count); err != nil {
			return err
		}
		if count > 0 {
			return errors.New("phone already exists")
		}
	}
	if email != "" {
		if err := s.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM admin_accounts WHERE email=? AND id<>?`, email, currentID).Scan(&count); err != nil {
			return err
		}
		if count > 0 {
			return errors.New("email already exists")
		}
	}
	return nil
}

func (s *appServer) getUserByAccount(ctx context.Context, account string) (userRow, error) {
	var u userRow
	err := s.db.QueryRowContext(ctx, `SELECT id,password_hash FROM users WHERE phone=?`, account).Scan(&u.ID, &u.PasswordHash)
	return u, err
}

func (s *appServer) getUserByPhone(ctx context.Context, phone string) (userRow, error) {
	var u userRow
	err := s.db.QueryRowContext(ctx, `SELECT id,password_hash FROM users WHERE phone=?`, phone).Scan(&u.ID, &u.PasswordHash)
	return u, err
}

func (s *appServer) getUserProfile(ctx context.Context, userID int64) (userProfile, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT u.id,u.uuid,u.phone AS username,u.nickname AS name,u.nickname,u.email,u.phone,u.bio,u.gender,u.birthday,u.location,u.avatar,'' AS cover,'user' AS role,'active' AS status,u.created_at,
		       m.plan_key,m.plan_name,m.started_at,m.expires_at,m.status
		FROM users u
		LEFT JOIN user_memberships m ON m.user_id=u.id
		WHERE u.id=?
	`, userID)
	return scanUserProfile(row)
}

func scanUserProfile(scanner interface{ Scan(dest ...any) error }) (userProfile, error) {
	u := userProfile{
		Membership: membershipInfo{},
	}
	var birthday sql.NullTime
	var created time.Time
	var mk, mn, ms sql.NullString
	var mstart, mend sql.NullTime
	err := scanner.Scan(&u.ID, &u.UUID, &u.Username, &u.Name, &u.Nickname, &u.Email, &u.Phone, &u.Bio, &u.Gender, &birthday, &u.Location, &u.Avatar, &u.Cover, &u.Role, &u.Status, &created, &mk, &mn, &mstart, &mend, &ms)
	if err != nil {
		return userProfile{}, err
	}
	u.CreatedAt = created.Format(time.RFC3339)
	if birthday.Valid {
		u.Birthday = birthday.Time.Format("2006-01-02")
	}
	u.Membership.HasMembership = mk.Valid || mn.Valid || mstart.Valid || mend.Valid || ms.Valid
	if mk.Valid {
		u.Membership.PlanKey = mk.String
	}
	if mn.Valid {
		u.Membership.PlanName = mn.String
	}
	if mstart.Valid {
		u.Membership.StartedAt = mstart.Time.Format(time.RFC3339)
	}
	if mend.Valid {
		u.Membership.ExpiresAt = mend.Time.Format(time.RFC3339)
		u.Membership.Status = membershipStatusByExpire(mend.Time)
	} else if ms.Valid {
		u.Membership.Status = ms.String
	}
	return u, nil
}

func (s *appServer) handleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		q := strings.TrimSpace(r.URL.Query().Get("q"))
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		rows, err := s.db.QueryContext(ctx, `
			SELECT u.id,u.uuid,u.phone AS username,u.nickname AS name,u.nickname,u.email,u.phone,u.bio,u.gender,u.birthday,u.location,u.avatar,'' AS cover,'user' AS role,'active' AS status,u.created_at,
			       m.plan_key,m.plan_name,m.started_at,m.expires_at,m.status
			FROM users u
			LEFT JOIN user_memberships m ON m.user_id=u.id
			WHERE ?='' OR u.nickname LIKE CONCAT('%',?,'%') OR u.phone LIKE CONCAT('%',?,'%') OR u.email LIKE CONCAT('%',?,'%')
			ORDER BY u.id DESC
		`, q, q, q, q)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		defer rows.Close()
		out := make([]userProfile, 0)
		for rows.Next() {
			item, err := scanUserProfile(rows)
			if err != nil {
				writeError(w, http.StatusInternalServerError, err)
				return
			}
			out = append(out, item)
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: out})
	case http.MethodPost:
		var in struct {
			Name     string `json:"name"`
			Nickname string `json:"nickname"`
			Email    string `json:"email"`
			Phone    string `json:"phone"`
			Bio      string `json:"bio"`
			Gender   string `json:"gender"`
			Location string `json:"location"`
			Avatar   string `json:"avatar"`
			Password string `json:"password"`
		}
		if err := readJSON(r, &in); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		nickname := emptyFallback(strings.TrimSpace(in.Nickname), emptyFallback(strings.TrimSpace(in.Name), "用户"))
		if nickname == "" || strings.TrimSpace(in.Phone) == "" {
			writeError(w, http.StatusBadRequest, errors.New("nickname and phone are required"))
			return
		}
		if strings.TrimSpace(in.Password) == "" {
			in.Password = "12345678"
		}
		hash, _ := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
		userUUID, err := newUUIDString()
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		res, err := s.db.ExecContext(ctx, `
			INSERT INTO users(uuid,nickname,email,phone,password_hash,bio,gender,location,avatar)
			VALUES(?,?,?,?,?,?,?,?,?)
		`,
			userUUID,
			nickname,
			strings.TrimSpace(in.Email),
			strings.TrimSpace(in.Phone),
			string(hash),
			strings.TrimSpace(in.Bio),
			strings.TrimSpace(in.Gender),
			strings.TrimSpace(in.Location),
			strings.TrimSpace(in.Avatar),
		)
		if err != nil {
			if isDuplicateErr(err) {
				writeError(w, http.StatusConflict, errors.New("phone already exists"))
				return
			}
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		id, _ := res.LastInsertId()
		item, err := s.getUserProfile(ctx, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusCreated, apiResponse{Message: "created", Data: item})
	default:
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}

func (s *appServer) handleUsersAdmin(w http.ResponseWriter, r *http.Request, _ int64) {
	s.handleUsers(w, r)
}

func (s *appServer) handleUserRoutes(w http.ResponseWriter, r *http.Request) {
	id, rest, ok := parseIDWithRest(w, r.URL.Path, "/api/users/")
	if !ok {
		return
	}
	if rest == "membership" {
		s.handleUserMembershipByID(w, r, id)
		return
	}
	if rest == "avatar" {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
			return
		}
		s.handleUserAvatarUpload(w, r, id)
		return
	}
	if rest != "" {
		writeError(w, http.StatusNotFound, errors.New("route not found"))
		return
	}
	switch r.Method {
	case http.MethodGet:
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		item, err := s.getUserProfile(ctx, id)
		if err != nil {
			writeError(w, http.StatusNotFound, errors.New("user not found"))
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: item})
	case http.MethodPut:
		var in struct {
			Name     string `json:"name"`
			Nickname string `json:"nickname"`
			Email    string `json:"email"`
			Phone    string `json:"phone"`
			Bio      string `json:"bio"`
			Gender   string `json:"gender"`
			Location string `json:"location"`
			Avatar   string `json:"avatar"`
			Password string `json:"password"`
		}
		if err := readJSON(r, &in); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		name := emptyFallback(strings.TrimSpace(in.Nickname), emptyFallback(strings.TrimSpace(in.Name), "\u7528\u6237"))
		err := s.updateUserCore(ctx, id, userProfile{
			Nickname: name,
			Email:    strings.TrimSpace(in.Email),
			Phone:    strings.TrimSpace(in.Phone),
			Bio:      strings.TrimSpace(in.Bio),
			Gender:   strings.TrimSpace(in.Gender),
			Location: strings.TrimSpace(in.Location),
			Avatar:   strings.TrimSpace(in.Avatar),
		}, strings.TrimSpace(in.Password))
		if err != nil {
			if strings.Contains(err.Error(), "password must be at least 8 chars") {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			if isDuplicateErr(err) {
				writeError(w, http.StatusConflict, errors.New("phone already exists"))
				return
			}
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		item, err := s.getUserProfile(ctx, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "updated", Data: item})
	case http.MethodDelete:
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		_, err := s.db.ExecContext(ctx, `DELETE FROM users WHERE id=?`, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "deleted"})
	default:
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}

func (s *appServer) handleUserRoutesAdmin(w http.ResponseWriter, r *http.Request, _ int64) {
	s.handleUserRoutes(w, r)
}

func (s *appServer) handleUserMembershipByID(w http.ResponseWriter, r *http.Request, userID int64) {
	if r.Method != http.MethodPut {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	var in struct {
		PlanKey   string `json:"planKey"`
		PlanName  string `json:"planName"`
		StartedAt string `json:"startedAt"`
		ExpiresAt string `json:"expiresAt"`
		Status    string `json:"status"`
	}
	if err := readJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if strings.TrimSpace(in.PlanKey) == "" {
		writeError(w, http.StatusBadRequest, errors.New("planKey is required"))
		return
	}
	planName := strings.TrimSpace(in.PlanName)
	if planName == "" {
		if n, _, ok := membershipPlan(in.PlanKey); ok {
			planName = n
		}
	}
	start := nullableTimeFromAny(in.StartedAt)
	expire := nullableTimeFromAny(in.ExpiresAt)
	if !start.Valid {
		start = sql.NullTime{Time: time.Now(), Valid: true}
	}
	if !expire.Valid {
		expire = sql.NullTime{Time: start.Time.AddDate(0, 1, 0), Valid: true}
	}
	status := strings.TrimSpace(in.Status)
	if status == "" {
		status = membershipStatusByExpire(expire.Time)
	}
	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO user_memberships(user_id,plan_key,plan_name,started_at,expires_at,status)
		VALUES(?,?,?,?,?,?)
		ON DUPLICATE KEY UPDATE plan_key=VALUES(plan_key),plan_name=VALUES(plan_name),started_at=VALUES(started_at),expires_at=VALUES(expires_at),status=VALUES(status)
	`, userID, in.PlanKey, planName, start.Time, expire.Time, status)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	item, err := s.getUserProfile(ctx, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Message: "updated", Data: item})
}

func (s *appServer) handleCountries(w http.ResponseWriter, r *http.Request) {
	if !s.requireAdminForWrite(w, r) {
		return
	}
	switch r.Method {
	case http.MethodGet:
		q := strings.TrimSpace(r.URL.Query().Get("q"))
		region := strings.TrimSpace(r.URL.Query().Get("region"))
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		rows, err := s.db.QueryContext(ctx, `
			SELECT id,name,code,region,flag,status,image,note,tags,keywords
			FROM countries
			WHERE (?='' OR region=?)
			  AND (?='' OR name LIKE CONCAT('%',?,'%') OR code LIKE CONCAT('%',?,'%') OR region LIKE CONCAT('%',?,'%') OR keywords LIKE CONCAT('%',?,'%'))
			ORDER BY id
		`, region, region, q, q, q, q, q)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		defer rows.Close()
		items := make([]countryItem, 0)
		for rows.Next() {
			var it countryItem
			var tags, keywords string
			if err := rows.Scan(&it.ID, &it.Name, &it.Code, &it.Region, &it.Flag, &it.Status, &it.Image, &it.Note, &tags, &keywords); err != nil {
				writeError(w, http.StatusInternalServerError, err)
				return
			}
			it.Tags = splitCSV(tags)
			it.Keywords = splitCSV(keywords)
			items = append(items, it)
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: items})
	case http.MethodPost:
		var in countryItem
		if err := readJSON(r, &in); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.Code) == "" {
			writeError(w, http.StatusBadRequest, errors.New("name and code are required"))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		res, err := s.db.ExecContext(ctx, `INSERT INTO countries(name,code,region,flag,status,image,note,tags,keywords) VALUES(?,?,?,?,?,?,?,?,?)`,
			in.Name, strings.ToUpper(strings.TrimSpace(in.Code)), strings.TrimSpace(in.Region), strings.TrimSpace(in.Flag), emptyFallback(in.Status, "active"), strings.TrimSpace(in.Image), strings.TrimSpace(in.Note), strings.Join(in.Tags, ","), strings.Join(in.Keywords, ","))
		if err != nil {
			if isDuplicateErr(err) {
				writeError(w, http.StatusConflict, errors.New("country code already exists"))
				return
			}
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		id, _ := res.LastInsertId()
		item, err := s.getCountryByID(ctx, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusCreated, apiResponse{Message: "created", Data: item})
	default:
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}

func (s *appServer) getCountryByID(ctx context.Context, id int64) (countryItem, error) {
	var it countryItem
	var tags, keywords string
	err := s.db.QueryRowContext(ctx, `SELECT id,name,code,region,flag,status,image,note,tags,keywords FROM countries WHERE id=?`, id).
		Scan(&it.ID, &it.Name, &it.Code, &it.Region, &it.Flag, &it.Status, &it.Image, &it.Note, &tags, &keywords)
	if err != nil {
		return countryItem{}, err
	}
	it.Tags = splitCSV(tags)
	it.Keywords = splitCSV(keywords)
	return it, nil
}

func (s *appServer) handleCountryRoutes(w http.ResponseWriter, r *http.Request) {
	if !s.requireAdminForWrite(w, r) {
		return
	}
	id, rest, ok := parseIDWithRest(w, r.URL.Path, "/api/visa/countries/")
	if !ok {
		return
	}
	if rest == "visas" {
		s.handleCountryVisasByCountryID(w, r, id)
		return
	}
	if rest != "" {
		writeError(w, http.StatusNotFound, errors.New("route not found"))
		return
	}
	switch r.Method {
	case http.MethodPut:
		var in countryItem
		if err := readJSON(r, &in); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		_, err := s.db.ExecContext(ctx, `UPDATE countries SET name=?,code=?,region=?,flag=?,status=?,image=?,note=?,tags=?,keywords=? WHERE id=?`,
			in.Name, strings.ToUpper(strings.TrimSpace(in.Code)), strings.TrimSpace(in.Region), strings.TrimSpace(in.Flag), emptyFallback(in.Status, "active"), strings.TrimSpace(in.Image), strings.TrimSpace(in.Note), strings.Join(in.Tags, ","), strings.Join(in.Keywords, ","), id)
		if err != nil {
			if isDuplicateErr(err) {
				writeError(w, http.StatusConflict, errors.New("country code already exists"))
				return
			}
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		item, err := s.getCountryByID(ctx, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "updated", Data: item})
	case http.MethodDelete:
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		_, err := s.db.ExecContext(ctx, `DELETE FROM countries WHERE id=?`, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "deleted"})
	default:
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}

func (s *appServer) handleHotDestinations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	limit, _ := strconv.Atoi(strings.TrimSpace(r.URL.Query().Get("limit")))
	items, err := s.listHotDestinations(r.Context(), strings.TrimSpace(r.URL.Query().Get("q")), limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: items})
}

func (s *appServer) handleFreeCountries(w http.ResponseWriter, r *http.Request) {
	if !s.requireAdminForWrite(w, r) {
		return
	}
	switch r.Method {
	case http.MethodGet:
		q := strings.TrimSpace(r.URL.Query().Get("q"))
		enabledRaw := strings.TrimSpace(r.URL.Query().Get("enabled"))
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		rows, err := s.db.QueryContext(ctx, `
			SELECT f.id,f.name,f.code,f.flag,f.region,f.city,f.policy_type,f.stay,f.note,f.map_x,f.map_y,f.enabled,c.id,c.name,v.id,v.name,f.keywords
			FROM visa_free_countries f
			LEFT JOIN visas v ON v.id=f.supported_visa_id
			LEFT JOIN countries c ON c.id=v.country_id
			WHERE (?='' OR f.name LIKE CONCAT('%',?,'%') OR f.code LIKE CONCAT('%',?,'%') OR f.region LIKE CONCAT('%',?,'%') OR f.keywords LIKE CONCAT('%',?,'%'))
			  AND (?='' OR f.enabled = CASE WHEN ? IN ('1','true','TRUE') THEN 1 ELSE 0 END)
			ORDER BY f.id DESC
		`, q, q, q, q, q, enabledRaw, enabledRaw)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		defer rows.Close()
		items := make([]freeCountryItem, 0)
		for rows.Next() {
			item, err := scanFreeCountry(rows)
			if err != nil {
				writeError(w, http.StatusInternalServerError, err)
				return
			}
			items = append(items, item)
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: items})
	case http.MethodPost:
		var in freeCountryItem
		if err := readJSON(r, &in); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if strings.TrimSpace(in.Name) == "" || strings.TrimSpace(in.Code) == "" {
			writeError(w, http.StatusBadRequest, errors.New("name and code are required"))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		supportedVisaID := nullableID(in.SupportedVisaID)
		supportedCountryID, resolveErr := s.resolveSupportedCountryIDFromVisa(ctx, in.SupportedVisaID)
		if resolveErr != nil {
			writeError(w, http.StatusBadRequest, resolveErr)
			return
		}
		res, err := s.db.ExecContext(ctx, `
			INSERT INTO visa_free_countries(name,code,flag,region,city,policy_type,stay,note,map_x,map_y,enabled,supported_country_id,supported_visa_id,keywords)
			VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)
		`,
			strings.TrimSpace(in.Name),
			strings.ToUpper(strings.TrimSpace(in.Code)),
			strings.TrimSpace(in.Flag),
			strings.TrimSpace(in.Region),
			strings.TrimSpace(in.City),
			strings.TrimSpace(in.PolicyType),
			strings.TrimSpace(in.Stay),
			strings.TrimSpace(in.Note),
			in.MapX,
			in.MapY,
			boolToInt(in.Enabled),
			supportedCountryID,
			supportedVisaID,
			strings.Join(in.Keywords, ","),
		)
		if err != nil {
			if isDuplicateErr(err) {
				writeError(w, http.StatusConflict, errors.New("free country code already exists"))
				return
			}
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		id, _ := res.LastInsertId()
		item, err := s.getFreeCountryByID(ctx, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusCreated, apiResponse{Message: "created", Data: item})
	default:
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}

func (s *appServer) handleFreeCountryByID(w http.ResponseWriter, r *http.Request) {
	if !s.requireAdminForWrite(w, r) {
		return
	}
	id, ok := parseID(w, r.URL.Path, "/api/visa/free-countries/")
	if !ok {
		return
	}
	switch r.Method {
	case http.MethodGet:
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		item, err := s.getFreeCountryByID(ctx, id)
		if err != nil {
			writeError(w, http.StatusNotFound, errors.New("free country not found"))
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: item})
	case http.MethodPut:
		var in freeCountryItem
		if err := readJSON(r, &in); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		supportedVisaID := nullableID(in.SupportedVisaID)
		supportedCountryID, resolveErr := s.resolveSupportedCountryIDFromVisa(ctx, in.SupportedVisaID)
		if resolveErr != nil {
			writeError(w, http.StatusBadRequest, resolveErr)
			return
		}
		_, err := s.db.ExecContext(ctx, `
			UPDATE visa_free_countries
			SET name=?,code=?,flag=?,region=?,city=?,policy_type=?,stay=?,note=?,map_x=?,map_y=?,enabled=?,supported_country_id=?,supported_visa_id=?,keywords=?
			WHERE id=?
		`,
			strings.TrimSpace(in.Name),
			strings.ToUpper(strings.TrimSpace(in.Code)),
			strings.TrimSpace(in.Flag),
			strings.TrimSpace(in.Region),
			strings.TrimSpace(in.City),
			strings.TrimSpace(in.PolicyType),
			strings.TrimSpace(in.Stay),
			strings.TrimSpace(in.Note),
			in.MapX,
			in.MapY,
			boolToInt(in.Enabled),
			supportedCountryID,
			supportedVisaID,
			strings.Join(in.Keywords, ","),
			id,
		)
		if err != nil {
			if isDuplicateErr(err) {
				writeError(w, http.StatusConflict, errors.New("free country code already exists"))
				return
			}
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		item, err := s.getFreeCountryByID(ctx, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "updated", Data: item})
	case http.MethodDelete:
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		if _, err := s.db.ExecContext(ctx, `DELETE FROM visa_free_countries WHERE id=?`, id); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "deleted"})
	default:
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}

func (s *appServer) getFreeCountryByID(ctx context.Context, id int64) (freeCountryItem, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT f.id,f.name,f.code,f.flag,f.region,f.city,f.policy_type,f.stay,f.note,f.map_x,f.map_y,f.enabled,c.id,c.name,v.id,v.name,f.keywords
		FROM visa_free_countries f
		LEFT JOIN visas v ON v.id=f.supported_visa_id
		LEFT JOIN countries c ON c.id=v.country_id
		WHERE f.id=?
	`, id)
	return scanFreeCountry(row)
}

func scanFreeCountry(scanner interface{ Scan(dest ...any) error }) (freeCountryItem, error) {
	var item freeCountryItem
	var enabled int
	var supportedCountryID sql.NullInt64
	var supportedCountryName sql.NullString
	var supportedVisaID sql.NullInt64
	var supportedVisaName sql.NullString
	var keywords string
	if err := scanner.Scan(
		&item.ID,
		&item.Name,
		&item.Code,
		&item.Flag,
		&item.Region,
		&item.City,
		&item.PolicyType,
		&item.Stay,
		&item.Note,
		&item.MapX,
		&item.MapY,
		&enabled,
		&supportedCountryID,
		&supportedCountryName,
		&supportedVisaID,
		&supportedVisaName,
		&keywords,
	); err != nil {
		return freeCountryItem{}, err
	}
	item.Enabled = enabled == 1
	if supportedCountryID.Valid {
		item.SupportedCountryID = supportedCountryID.Int64
	}
	if supportedCountryName.Valid {
		item.SupportedCountryName = supportedCountryName.String
	}
	if supportedVisaID.Valid {
		item.SupportedVisaID = supportedVisaID.Int64
	}
	if supportedVisaName.Valid {
		item.SupportedVisaName = supportedVisaName.String
	}
	item.Keywords = splitCSV(keywords)
	return item, nil
}

func (s *appServer) handleCountryVisasByCountryID(w http.ResponseWriter, r *http.Request, countryID int64) {
	if !s.requireAdminForWrite(w, r) {
		return
	}
	switch r.Method {
	case http.MethodGet:
		items, err := s.listVisas(r.Context(), countryID, strings.TrimSpace(r.URL.Query().Get("q")))
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: items})
	case http.MethodPost:
		var in visaListItem
		if err := readJSON(r, &in); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		in.CountryID = countryID
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		item, err := s.createVisa(ctx, in)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusCreated, apiResponse{Message: "created", Data: item})
	default:
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}

func (s *appServer) handleCountryVisas(w http.ResponseWriter, r *http.Request) {
	if !s.requireAdminForWrite(w, r) {
		return
	}
	switch r.Method {
	case http.MethodGet:
		cid, _ := strconv.ParseInt(strings.TrimSpace(r.URL.Query().Get("countryId")), 10, 64)
		items, err := s.listVisas(r.Context(), cid, strings.TrimSpace(r.URL.Query().Get("q")))
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: items})
	case http.MethodPost:
		var in visaListItem
		if err := readJSON(r, &in); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		item, err := s.createVisa(ctx, in)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusCreated, apiResponse{Message: "created", Data: item})
	default:
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}

func (s *appServer) listVisas(ctx context.Context, countryID int64, q string) ([]visaListItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, `
		SELECT v.id,v.country_id,c.name,v.name,v.visa_type,v.processing_time,v.fee,v.validity,v.entries,v.status,v.description,v.long_intro,v.hot,v.visa_free,v.updated_at
		FROM visas v
		JOIN countries c ON c.id=v.country_id
		WHERE (?=0 OR v.country_id=?) AND (?='' OR v.name LIKE CONCAT('%',?,'%') OR v.visa_type LIKE CONCAT('%',?,'%') OR c.name LIKE CONCAT('%',?,'%'))
		ORDER BY v.id DESC
	`, countryID, countryID, q, q, q, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]visaListItem, 0)
	for rows.Next() {
		var it visaListItem
		var hot, visaFree int
		var updated time.Time
		if err := rows.Scan(&it.ID, &it.CountryID, &it.CountryName, &it.Name, &it.VisaType, &it.ProcessingTime, &it.Fee, &it.Validity, &it.Entries, &it.Status, &it.Description, &it.LongIntro, &hot, &visaFree, &updated); err != nil {
			return nil, err
		}
		it.Hot = hot == 1
		it.VisaFree = visaFree == 1
		it.UpdatedAt = updated.Format(time.RFC3339)
		out = append(out, it)
	}
	return out, nil
}

func (s *appServer) createVisa(ctx context.Context, in visaListItem) (visaListItem, error) {
	if in.CountryID == 0 || strings.TrimSpace(in.Name) == "" {
		return visaListItem{}, errors.New("countryId and name are required")
	}
	res, err := s.db.ExecContext(ctx, `INSERT INTO visas(country_id,name,visa_type,processing_time,fee,validity,entries,status,description,long_intro,hot,visa_free) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)`,
		in.CountryID, strings.TrimSpace(in.Name), strings.TrimSpace(in.VisaType), strings.TrimSpace(in.ProcessingTime), strings.TrimSpace(in.Fee), strings.TrimSpace(in.Validity), strings.TrimSpace(in.Entries), emptyFallback(in.Status, "active"), strings.TrimSpace(in.Description), strings.TrimSpace(in.LongIntro), boolToInt(in.Hot), boolToInt(in.VisaFree))
	if err != nil {
		return visaListItem{}, err
	}
	id, _ := res.LastInsertId()
	tx, err := s.db.BeginTx(ctx, nil)
	if err == nil {
		country := seedCountryInfo{ID: in.CountryID}
		_ = s.db.QueryRowContext(ctx, `SELECT code,name,image FROM countries WHERE id=?`, in.CountryID).Scan(&country.Code, &country.Name, &country.Image)
		steps := buildReferenceSteps(country, strings.TrimSpace(in.Name), emptyFallback(strings.TrimSpace(in.VisaType), "\u65c5\u6e38\u7b7e"), guideSourceForCountry(country), defaultMaterialsForCustomVisa(strings.TrimSpace(in.VisaType)), requiresAppointmentForCountry(country.Code))
		_ = replaceVisaStepsTx(ctx, tx, id, steps)
		_ = tx.Commit()
	}
	return s.getVisaByID(ctx, id)
}

func (s *appServer) getVisaByID(ctx context.Context, id int64) (visaListItem, error) {
	var it visaListItem
	var hot, visaFree int
	var updated time.Time
	err := s.db.QueryRowContext(ctx, `
		SELECT v.id,v.country_id,c.name,v.name,v.visa_type,v.processing_time,v.fee,v.validity,v.entries,v.status,v.description,v.long_intro,v.hot,v.visa_free,v.updated_at
		FROM visas v
		JOIN countries c ON c.id=v.country_id
		WHERE v.id=?
	`, id).Scan(&it.ID, &it.CountryID, &it.CountryName, &it.Name, &it.VisaType, &it.ProcessingTime, &it.Fee, &it.Validity, &it.Entries, &it.Status, &it.Description, &it.LongIntro, &hot, &visaFree, &updated)
	if err != nil {
		return visaListItem{}, err
	}
	it.Hot = hot == 1
	it.VisaFree = visaFree == 1
	it.UpdatedAt = updated.Format(time.RFC3339)
	return it, nil
}

func (s *appServer) handleCountryVisaByID(w http.ResponseWriter, r *http.Request) {
	if !s.requireAdminForWrite(w, r) {
		return
	}
	id, ok := parseID(w, r.URL.Path, "/api/visa/country-visas/")
	if !ok {
		return
	}
	switch r.Method {
	case http.MethodPut:
		var in visaListItem
		if err := readJSON(r, &in); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		_, err := s.db.ExecContext(ctx, `UPDATE visas SET country_id=?,name=?,visa_type=?,processing_time=?,fee=?,validity=?,entries=?,status=?,description=?,long_intro=?,hot=?,visa_free=? WHERE id=?`,
			in.CountryID, in.Name, in.VisaType, in.ProcessingTime, in.Fee, in.Validity, in.Entries, emptyFallback(in.Status, "active"), in.Description, in.LongIntro, boolToInt(in.Hot), boolToInt(in.VisaFree), id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		item, err := s.getVisaByID(ctx, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "updated", Data: item})
	case http.MethodDelete:
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		_, err := s.db.ExecContext(ctx, `DELETE FROM visas WHERE id=?`, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "deleted"})
	default:
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}

func (s *appServer) handleVisaRoutes(w http.ResponseWriter, r *http.Request) {
	if !s.requireAdminForWrite(w, r) {
		return
	}
	id, rest, ok := parseIDWithRest(w, r.URL.Path, "/api/visa/visas/")
	if !ok {
		return
	}
	if rest == "" {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
			return
		}
		detail, err := s.getVisaDetail(r.Context(), id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: detail})
		return
	}
	if rest == "detail" {
		if r.Method != http.MethodPut {
			writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
			return
		}
		var in struct {
			Steps []stepItem `json:"steps"`
		}
		if err := readJSON(r, &in); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
		defer cancel()
		if err := s.saveVisaDetail(ctx, id, in.Steps); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		detail, err := s.getVisaDetail(ctx, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "updated", Data: detail})
		return
	}
	writeError(w, http.StatusNotFound, errors.New("route not found"))
}

func (s *appServer) getVisaDetail(ctx context.Context, visaID int64) (visaDetail, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	base, err := s.getVisaByID(ctx, visaID)
	if err != nil {
		return visaDetail{}, err
	}
	steps, err := s.loadSteps(ctx, visaID, 0, false)
	if err != nil {
		return visaDetail{}, err
	}
	return visaDetail{visaListItem: base, Steps: steps}, nil
}

func (s *appServer) saveVisaDetail(ctx context.Context, visaID int64, steps []stepItem) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var oldIDs []int64
	rows, err := tx.QueryContext(ctx, `SELECT id FROM visa_steps WHERE visa_id=?`, visaID)
	if err == nil {
		for rows.Next() {
			var id int64
			_ = rows.Scan(&id)
			oldIDs = append(oldIDs, id)
		}
		rows.Close()
	}
	for _, sid := range oldIDs {
		_, _ = tx.ExecContext(ctx, `DELETE FROM visa_step_strategies WHERE step_id=?`, sid)
		_, _ = tx.ExecContext(ctx, `DELETE FROM visa_step_guides WHERE step_id=?`, sid)
		_, _ = tx.ExecContext(ctx, `DELETE FROM visa_step_materials WHERE step_id=?`, sid)
		_, _ = tx.ExecContext(ctx, `DELETE FROM visa_step_tasks WHERE step_id=?`, sid)
	}
	_, _ = tx.ExecContext(ctx, `DELETE FROM visa_steps WHERE visa_id=?`, visaID)

	for idx, st := range steps {
		sortOrder := st.SortOrder
		if sortOrder == 0 {
			sortOrder = idx + 1
		}
		res, err := tx.ExecContext(ctx, `INSERT INTO visa_steps(visa_id,step_key,title,sort_order) VALUES(?,?,?,?)`, visaID, emptyFallback(st.StepKey, fmt.Sprintf("step-%d", idx+1)), emptyFallback(st.Title, "\u6b65\u9aa4"), sortOrder)
		if err != nil {
			return err
		}
		stepID, _ := res.LastInsertId()
		for i, c := range st.Strategies {
			if strings.TrimSpace(c) == "" {
				continue
			}
			if _, err := tx.ExecContext(ctx, `INSERT INTO visa_step_strategies(step_id,content,sort_order) VALUES(?,?,?)`, stepID, c, i+1); err != nil {
				return err
			}
		}
		for i, g := range st.Guides {
			if strings.TrimSpace(g.Title) == "" && strings.TrimSpace(g.Description) == "" {
				continue
			}
			if _, err := tx.ExecContext(ctx, `INSERT INTO visa_step_guides(step_id,title,description,image,cta,link_url,sort_order) VALUES(?,?,?,?,?,?,?)`, stepID, emptyFallback(g.Title, "\u529e\u7406\u6307\u5357"), g.Description, g.Image, g.CTA, g.URL, i+1); err != nil {
				return err
			}
		}
		for i, m := range st.Materials {
			if strings.TrimSpace(m) == "" {
				continue
			}
			if _, err := tx.ExecContext(ctx, `INSERT INTO visa_step_materials(step_id,content,sort_order) VALUES(?,?,?)`, stepID, m, i+1); err != nil {
				return err
			}
		}
		for i, t := range st.Tasks {
			if strings.TrimSpace(t.Title) == "" {
				continue
			}
			if _, err := tx.ExecContext(ctx, `INSERT INTO visa_step_tasks(step_id,task_key,title,icon,default_status,default_status_text,sort_order) VALUES(?,?,?,?,?,?,?)`,
				stepID, emptyFallback(t.TaskKey, fmt.Sprintf("task-%d", i+1)), t.Title, emptyFallback(t.Icon, "task_alt"), "todo", defaultStatusText("todo"), i+1); err != nil {
				return err
			}
		}
	}
	return tx.Commit()
}

func (s *appServer) loadSteps(ctx context.Context, visaID, planID int64, includeState bool) ([]stepItem, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id,step_key,title,sort_order FROM visa_steps WHERE visa_id=? ORDER BY sort_order,id`, visaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	steps := make([]stepItem, 0)
	for rows.Next() {
		var st stepItem
		if err := rows.Scan(&st.ID, &st.StepKey, &st.Title, &st.SortOrder); err != nil {
			return nil, err
		}
		st.Strategies, _ = s.loadStrategies(ctx, st.ID)
		st.Guides, _ = s.loadGuides(ctx, st.ID)
		st.Materials, _ = s.loadMaterials(ctx, st.ID)
		st.Tasks, _ = s.loadTasks(ctx, st.ID, planID, includeState)
		steps = append(steps, st)
	}
	if includeState {
		steps = applyStepStatuses(steps)
	}
	return steps, nil
}

func (s *appServer) loadStrategies(ctx context.Context, stepID int64) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT content FROM visa_step_strategies WHERE step_id=? ORDER BY sort_order,id`, stepID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]string, 0)
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, nil
}

func (s *appServer) loadGuides(ctx context.Context, stepID int64) ([]guideItem, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT title,description,image,cta,link_url FROM visa_step_guides WHERE step_id=? ORDER BY sort_order,id`, stepID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]guideItem, 0)
	for rows.Next() {
		var g guideItem
		if err := rows.Scan(&g.Title, &g.Description, &g.Image, &g.CTA, &g.URL); err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, nil
}

func (s *appServer) loadMaterials(ctx context.Context, stepID int64) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT content FROM visa_step_materials WHERE step_id=? ORDER BY sort_order,id`, stepID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]string, 0)
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, nil
}

func (s *appServer) loadTasks(ctx context.Context, stepID, planID int64, includeState bool) ([]taskItem, error) {
	query := `SELECT t.id,t.task_key,t.title,t.icon,t.default_status,t.default_status_text,t.sort_order FROM visa_step_tasks t WHERE t.step_id=? ORDER BY t.sort_order,t.id`
	rows, err := s.db.QueryContext(ctx, query, stepID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]taskItem, 0)
	for rows.Next() {
		var t taskItem
		if err := rows.Scan(&t.ID, &t.TaskKey, &t.Title, &t.Icon, &t.Status, &t.StatusText, &t.SortOrder); err != nil {
			return nil, err
		}
		if !includeState {
			t.Status = ""
			t.StatusText = ""
			out = append(out, t)
			continue
		}
		if includeState && planID > 0 {
			var st, stText sql.NullString
			_ = s.db.QueryRowContext(ctx, `SELECT status,status_text FROM user_plan_task_states WHERE plan_id=? AND task_id=?`, planID, t.ID).Scan(&st, &stText)
			if st.Valid {
				t.Status = st.String
			}
			if stText.Valid {
				t.StatusText = stText.String
			}
		}
		out = append(out, t)
	}
	return out, nil
}

func (s *appServer) listHotDestinations(ctx context.Context, q string, limit int) ([]hotDestinationItem, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	if limit <= 0 || limit > 20 {
		limit = 8
	}
	rows, err := s.db.QueryContext(ctx, `
		SELECT v.country_id,v.id,c.name,c.flag,c.image,c.note,v.name,v.visa_type,v.processing_time,v.fee,v.hot
		FROM visas v
		JOIN countries c ON c.id=v.country_id
		WHERE v.hot=1
		  AND (?='' OR c.name LIKE CONCAT('%',?,'%') OR v.name LIKE CONCAT('%',?,'%') OR v.visa_type LIKE CONCAT('%',?,'%'))
		ORDER BY v.updated_at DESC, v.id DESC
	`, q, q, q, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	seenCountry := make(map[int64]struct{})
	items := make([]hotDestinationItem, 0, limit)
	for rows.Next() {
		var item hotDestinationItem
		var hot int
		if err := rows.Scan(&item.CountryID, &item.VisaID, &item.Name, &item.Flag, &item.Image, &item.Note, &item.VisaName, &item.Type, &item.Time, &item.Price, &hot); err != nil {
			return nil, err
		}
		if _, exists := seenCountry[item.CountryID]; exists {
			continue
		}
		seenCountry[item.CountryID] = struct{}{}
		item.Hot = hot == 1
		if strings.TrimSpace(item.Note) == "" {
			item.Note = strings.TrimSpace(item.VisaName)
		}
		items = append(items, item)
		if len(items) >= limit {
			break
		}
	}
	return items, nil
}

func applyStepStatuses(steps []stepItem) []stepItem {
	active := -1
	for i, st := range steps {
		allDone := len(st.Tasks) > 0
		for _, t := range st.Tasks {
			if t.Status != "done" {
				allDone = false
				break
			}
		}
		if !allDone && active == -1 {
			active = i
		}
	}
	for i := range steps {
		if active == -1 {
			steps[i].Status = "done"
			continue
		}
		if i < active {
			steps[i].Status = "done"
		} else if i == active {
			steps[i].Status = "active"
		} else {
			steps[i].Status = "todo"
		}
	}
	return steps
}

func (s *appServer) hasPlanSnapshot(ctx context.Context, planID int64) (bool, error) {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM user_plan_steps WHERE plan_id=?`, planID).Scan(&count); err != nil {
		return false, err
	}
	return count > 0, nil
}

func (s *appServer) ensurePlanSnapshot(ctx context.Context, planID, visaID int64) error {
	hasSnapshot, err := s.hasPlanSnapshot(ctx, planID)
	if err != nil {
		return err
	}
	if hasSnapshot {
		return nil
	}
	steps, err := s.loadSteps(ctx, visaID, planID, true)
	if err != nil {
		return err
	}
	return s.savePlanSnapshot(ctx, planID, steps)
}

func (s *appServer) savePlanSnapshot(ctx context.Context, planID int64, steps []stepItem) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := clearPlanSnapshotTx(ctx, tx, planID); err != nil {
		return err
	}

	for stepIndex, step := range steps {
		sortOrder := step.SortOrder
		if sortOrder == 0 {
			sortOrder = stepIndex + 1
		}
		res, err := tx.ExecContext(ctx, `INSERT INTO user_plan_steps(plan_id,step_key,title,sort_order) VALUES(?,?,?,?)`,
			planID, emptyFallback(step.StepKey, fmt.Sprintf("step-%d", stepIndex+1)), emptyFallback(step.Title, "步骤"), sortOrder)
		if err != nil {
			return err
		}
		planStepID, _ := res.LastInsertId()

		for i, item := range step.Strategies {
			if strings.TrimSpace(item) == "" {
				continue
			}
			if _, err := tx.ExecContext(ctx, `INSERT INTO user_plan_step_strategies(plan_step_id,content,sort_order) VALUES(?,?,?)`, planStepID, strings.TrimSpace(item), i+1); err != nil {
				return err
			}
		}
		for i, guide := range step.Guides {
			if strings.TrimSpace(guide.Title) == "" && strings.TrimSpace(guide.Description) == "" {
				continue
			}
			if _, err := tx.ExecContext(ctx, `INSERT INTO user_plan_step_guides(plan_step_id,title,description,image,cta,link_url,sort_order) VALUES(?,?,?,?,?,?,?)`,
				planStepID, emptyFallback(guide.Title, "办理指南"), strings.TrimSpace(guide.Description), strings.TrimSpace(guide.Image), strings.TrimSpace(guide.CTA), strings.TrimSpace(guide.URL), i+1); err != nil {
				return err
			}
		}
		for i, item := range step.Materials {
			if strings.TrimSpace(item) == "" {
				continue
			}
			if _, err := tx.ExecContext(ctx, `INSERT INTO user_plan_step_materials(plan_step_id,content,sort_order) VALUES(?,?,?)`, planStepID, strings.TrimSpace(item), i+1); err != nil {
				return err
			}
		}
		for i, task := range step.Tasks {
			status := strings.TrimSpace(task.Status)
			if status == "" {
				status = "todo"
			}
			statusText := strings.TrimSpace(task.StatusText)
			if statusText == "" {
				statusText = defaultStatusText(status)
			}
			if _, err := tx.ExecContext(ctx, `INSERT INTO user_plan_step_tasks(plan_step_id,task_key,title,icon,status,status_text,sort_order) VALUES(?,?,?,?,?,?,?)`,
				planStepID, emptyFallback(task.TaskKey, fmt.Sprintf("task-%d", i+1)), emptyFallback(task.Title, "未命名任务"), emptyFallback(task.Icon, "task_alt"), status, statusText, i+1); err != nil {
				return err
			}
		}
	}

	return tx.Commit()
}

func (s *appServer) seedFreeCountries(ctx context.Context) error {
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM visa_free_countries`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	type freeSeed struct {
		Name       string
		Code       string
		Flag       string
		Region     string
		City       string
		PolicyType string
		Stay       string
		Note       string
		MapX       float64
		MapY       float64
		Keywords   []string
	}

	seeds := []freeSeed{
		{"泰国", "TH", "🇹🇭", "亚洲", "曼谷", "落地签", "最长30天", "持中国护照可申请落地签，建议准备返程机票与酒店订单。", 74, 62, []string{"泰国", "落地签", "曼谷", "旅游"}},
		{"新加坡", "SG", "🇸🇬", "亚洲", "新加坡", "免签过境/电子入境", "短停友好", "中转和短停场景常见，入境前请确认最新政策。", 77, 68, []string{"新加坡", "过境", "电子入境", "短停"}},
		{"阿联酋", "AE", "🇦🇪", "亚洲", "迪拜", "免签", "最长30天", "商务与旅游热度高，护照有效期需满足要求。", 62, 54, []string{"阿联酋", "迪拜", "免签", "中转"}},
		{"马来西亚", "MY", "🇲🇾", "亚洲", "吉隆坡", "免签", "最长30天", "热门自由行目的地，建议提前准备回程证明。", 76, 66, []string{"马来西亚", "吉隆坡", "免签", "旅游"}},
		{"印度尼西亚", "ID", "🇮🇩", "亚洲", "雅加达", "落地签", "最长30天", "巴厘岛等目的地受欢迎，入境材料要齐全。", 80, 72, []string{"印度尼西亚", "印尼", "雅加达", "落地签"}},
		{"卡塔尔", "QA", "🇶🇦", "亚洲", "多哈", "免签", "最长30天", "常见中转地，建议核对航班衔接时间。", 60, 56, []string{"卡塔尔", "多哈", "免签", "中转"}},
		{"摩洛哥", "MA", "🇲🇦", "非洲", "马拉喀什", "免签", "最长90天", "北非热门目的地，入境时请准备行程单。", 43, 49, []string{"摩洛哥", "马拉喀什", "免签", "北非"}},
		{"毛里求斯", "MU", "🇲🇺", "非洲", "路易港", "免签", "最长60天", "海岛旅游热门，建议准备保险和酒店凭证。", 61, 78, []string{"毛里求斯", "路易港", "免签", "海岛"}},
	}

	countryIDByCode := map[string]int64{}
	rows, err := s.db.QueryContext(ctx, `SELECT id,code FROM countries`)
	if err != nil {
		return err
	}
	for rows.Next() {
		var id int64
		var code string
		if err := rows.Scan(&id, &code); err != nil {
			rows.Close()
			return err
		}
		countryIDByCode[strings.ToUpper(strings.TrimSpace(code))] = id
	}
	rows.Close()

	visaIDByCode := map[string]int64{}
	visaRows, err := s.db.QueryContext(ctx, `
		SELECT c.code,v.id
		FROM visas v
		JOIN countries c ON c.id=v.country_id
		WHERE v.visa_free=1
	`)
	if err != nil {
		return err
	}
	for visaRows.Next() {
		var code string
		var visaID int64
		if err := visaRows.Scan(&code, &visaID); err != nil {
			visaRows.Close()
			return err
		}
		visaIDByCode[strings.ToUpper(strings.TrimSpace(code))] = visaID
	}
	visaRows.Close()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, item := range seeds {
		supportedCountryID := nullableID(countryIDByCode[item.Code])
		supportedVisaID := nullableID(visaIDByCode[item.Code])
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO visa_free_countries(name,code,flag,region,city,policy_type,stay,note,map_x,map_y,enabled,supported_country_id,supported_visa_id,keywords)
			VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)
		`,
			item.Name,
			item.Code,
			item.Flag,
			item.Region,
			item.City,
			item.PolicyType,
			item.Stay,
			item.Note,
			item.MapX,
			item.MapY,
			1,
			supportedCountryID,
			supportedVisaID,
			strings.Join(item.Keywords, ","),
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *appServer) syncFreeCountryVisaMappings(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()
	rows, err := s.db.QueryContext(ctx, `SELECT id,code,supported_country_id,supported_visa_id FROM visa_free_countries`)
	if err != nil {
		return err
	}
	defer rows.Close()

	type mappingRow struct {
		ID                 int64
		Code               string
		SupportedCountryID sql.NullInt64
		SupportedVisaID    sql.NullInt64
	}
	items := make([]mappingRow, 0)
	for rows.Next() {
		var item mappingRow
		if err := rows.Scan(&item.ID, &item.Code, &item.SupportedCountryID, &item.SupportedVisaID); err != nil {
			return err
		}
		items = append(items, item)
	}

	for _, item := range items {
		if item.SupportedVisaID.Valid && item.SupportedVisaID.Int64 > 0 {
			continue
		}
		var visaID int64
		query := `
			SELECT v.id
			FROM visas v
			JOIN countries c ON c.id=v.country_id
			WHERE v.visa_free=1 AND c.code=?
			ORDER BY v.id DESC
			LIMIT 1
		`
		args := []any{strings.ToUpper(strings.TrimSpace(item.Code))}
		if item.SupportedCountryID.Valid && item.SupportedCountryID.Int64 > 0 {
			query = `
				SELECT id
				FROM visas
				WHERE country_id=? AND visa_free=1
				ORDER BY id DESC
				LIMIT 1
			`
			args = []any{item.SupportedCountryID.Int64}
		}
		if err := s.db.QueryRowContext(ctx, query, args...).Scan(&visaID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			return err
		}
		supportedCountryID, err := s.resolveSupportedCountryIDFromVisa(ctx, visaID)
		if err != nil {
			return err
		}
		if _, err := s.db.ExecContext(ctx, `UPDATE visa_free_countries SET supported_country_id=?, supported_visa_id=? WHERE id=?`, supportedCountryID, visaID, item.ID); err != nil {
			return err
		}
	}
	return nil
}

func (s *appServer) resolveSupportedCountryIDFromVisa(ctx context.Context, visaID int64) (any, error) {
	if visaID <= 0 {
		return nil, nil
	}
	var countryID int64
	if err := s.db.QueryRowContext(ctx, `SELECT country_id FROM visas WHERE id=?`, visaID).Scan(&countryID); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("mapped visa not found")
		}
		return nil, err
	}
	return nullableID(countryID), nil
}

func clearPlanSnapshotTx(ctx context.Context, tx *sql.Tx, planID int64) error {
	rows, err := tx.QueryContext(ctx, `SELECT id FROM user_plan_steps WHERE plan_id=?`, planID)
	if err != nil {
		return err
	}

	planStepIDs := make([]int64, 0)
	for rows.Next() {
		var stepID int64
		if err := rows.Scan(&stepID); err != nil {
			rows.Close()
			return err
		}
		planStepIDs = append(planStepIDs, stepID)
	}
	rows.Close()

	for _, stepID := range planStepIDs {
		if _, err := tx.ExecContext(ctx, `DELETE FROM user_plan_step_strategies WHERE plan_step_id=?`, stepID); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `DELETE FROM user_plan_step_guides WHERE plan_step_id=?`, stepID); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `DELETE FROM user_plan_step_materials WHERE plan_step_id=?`, stepID); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, `DELETE FROM user_plan_step_tasks WHERE plan_step_id=?`, stepID); err != nil {
			return err
		}
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM user_plan_steps WHERE plan_id=?`, planID); err != nil {
		return err
	}
	return nil
}

func (s *appServer) loadPlanSnapshotSteps(ctx context.Context, planID int64) ([]stepItem, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id,step_key,title,sort_order FROM user_plan_steps WHERE plan_id=? ORDER BY sort_order,id`, planID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	steps := make([]stepItem, 0)
	for rows.Next() {
		var step stepItem
		if err := rows.Scan(&step.ID, &step.StepKey, &step.Title, &step.SortOrder); err != nil {
			return nil, err
		}
		step.Strategies, _ = s.loadPlanSnapshotStrategies(ctx, step.ID)
		step.Guides, _ = s.loadPlanSnapshotGuides(ctx, step.ID)
		step.Materials, _ = s.loadPlanSnapshotMaterials(ctx, step.ID)
		step.Tasks, _ = s.loadPlanSnapshotTasks(ctx, step.ID)
		steps = append(steps, step)
	}
	return applyStepStatuses(steps), nil
}

func (s *appServer) loadPlanSnapshotStrategies(ctx context.Context, planStepID int64) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT content FROM user_plan_step_strategies WHERE plan_step_id=? ORDER BY sort_order,id`, planStepID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]string, 0)
	for rows.Next() {
		var item string
		if err := rows.Scan(&item); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}

func (s *appServer) loadPlanSnapshotGuides(ctx context.Context, planStepID int64) ([]guideItem, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT title,description,image,cta,link_url FROM user_plan_step_guides WHERE plan_step_id=? ORDER BY sort_order,id`, planStepID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]guideItem, 0)
	for rows.Next() {
		var item guideItem
		if err := rows.Scan(&item.Title, &item.Description, &item.Image, &item.CTA, &item.URL); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}

func (s *appServer) loadPlanSnapshotMaterials(ctx context.Context, planStepID int64) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT content FROM user_plan_step_materials WHERE plan_step_id=? ORDER BY sort_order,id`, planStepID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]string, 0)
	for rows.Next() {
		var item string
		if err := rows.Scan(&item); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}

func (s *appServer) loadPlanSnapshotTasks(ctx context.Context, planStepID int64) ([]taskItem, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id,task_key,title,icon,status,status_text,sort_order FROM user_plan_step_tasks WHERE plan_step_id=? ORDER BY sort_order,id`, planStepID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]taskItem, 0)
	for rows.Next() {
		var item taskItem
		if err := rows.Scan(&item.ID, &item.TaskKey, &item.Title, &item.Icon, &item.Status, &item.StatusText, &item.SortOrder); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}

func (s *appServer) updatePlanLifecycleStatus(ctx context.Context, planID int64) error {
	var progress int
	var resultStatus string
	if err := s.db.QueryRowContext(ctx, `SELECT progress,result_status FROM user_visa_plans WHERE id=?`, planID).Scan(&progress, &resultStatus); err != nil {
		return err
	}
	nextStatus := "active"
	switch strings.TrimSpace(resultStatus) {
	case "approved", "rejected", "withdrawn":
		nextStatus = "closed"
	default:
		if progress >= 100 {
			nextStatus = "completed"
		}
	}
	_, err := s.db.ExecContext(ctx, `UPDATE user_visa_plans SET status=? WHERE id=?`, nextStatus, planID)
	return err
}

func (s *appServer) handlePlans(w http.ResponseWriter, r *http.Request, uid int64) {
	switch r.Method {
	case http.MethodGet:
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		rows, err := s.db.QueryContext(ctx, `
			SELECT p.id,p.country_id,p.country_name,c.flag,p.visa_id,p.visa_title,p.progress,p.active_step_key,p.status,p.result_status,p.result_note,p.result_at,p.created_at
			FROM user_visa_plans p
			LEFT JOIN countries c ON c.id=p.country_id
			WHERE p.user_id=?
			ORDER BY p.created_at DESC
		`, uid)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		defer rows.Close()
		items := make([]planSummary, 0)
		for rows.Next() {
			var it planSummary
			var created time.Time
			var resultAt sql.NullTime
			if err := rows.Scan(&it.ID, &it.CountryID, &it.CountryName, &it.CountryFlag, &it.VisaID, &it.VisaTitle, &it.Progress, &it.ActiveStepKey, &it.Status, &it.ResultStatus, &it.ResultNote, &resultAt, &created); err != nil {
				writeError(w, http.StatusInternalServerError, err)
				return
			}
			it.CreatedAt = created.Format(time.RFC3339)
			if resultAt.Valid {
				it.ResultAt = resultAt.Time.Format(time.RFC3339)
			}
			items = append(items, it)
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: items})
	case http.MethodPost:
		var in struct {
			CountryID   int64  `json:"countryId"`
			VisaID      int64  `json:"visaId"`
			CountryName string `json:"countryName"`
			VisaTitle   string `json:"visaTitle"`
			Source      string `json:"source"`
		}
		if err := readJSON(r, &in); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		if in.CountryID == 0 || in.VisaID == 0 {
			writeError(w, http.StatusBadRequest, errors.New("countryId and visaId are required"))
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		if in.CountryName == "" {
			_ = s.db.QueryRowContext(ctx, `SELECT name FROM countries WHERE id=?`, in.CountryID).Scan(&in.CountryName)
		}
		if in.VisaTitle == "" {
			_ = s.db.QueryRowContext(ctx, `SELECT name FROM visas WHERE id=?`, in.VisaID).Scan(&in.VisaTitle)
		}
		var activePlanID int64
		err := s.db.QueryRowContext(ctx, `SELECT id FROM user_visa_plans WHERE user_id=? AND visa_id=? AND status='active' ORDER BY id DESC LIMIT 1`, uid, in.VisaID).Scan(&activePlanID)
		if err == nil && activePlanID > 0 {
			writeError(w, http.StatusConflict, errors.New("an active plan for this visa already exists"))
			return
		}
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		res, err := s.db.ExecContext(ctx, `INSERT INTO user_visa_plans(user_id,country_id,visa_id,country_name,visa_title,source,progress,active_step_key,status) VALUES(?,?,?,?,?,?,0,'apply','active')`,
			uid, in.CountryID, in.VisaID, in.CountryName, in.VisaTitle, emptyFallback(in.Source, "visa"))
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		planID, _ := res.LastInsertId()
		if err := s.ensurePlanSnapshot(ctx, planID, in.VisaID); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		if err := s.refreshPlanProgress(ctx, planID); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		detail, err := s.getPlanDetail(ctx, uid, planID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusCreated, apiResponse{Message: "created", Data: detail})
	default:
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}

func (s *appServer) handlePlanRoutes(w http.ResponseWriter, r *http.Request, uid int64) {
	id, rest, ok := parseIDWithRest(w, r.URL.Path, "/api/plans/")
	if !ok {
		return
	}
	if strings.HasPrefix(rest, "tasks/") {
		taskIDRaw := strings.TrimPrefix(rest, "tasks/")
		taskID, err := strconv.ParseInt(taskIDRaw, 10, 64)
		if err != nil || taskID <= 0 {
			writeError(w, http.StatusBadRequest, errors.New("invalid task id"))
			return
		}
		s.handlePlanTaskPatch(w, r, uid, id, taskID)
		return
	}
	if rest == "result" {
		s.handlePlanResultPatch(w, r, uid, id)
		return
	}
	if rest != "" {
		writeError(w, http.StatusNotFound, errors.New("route not found"))
		return
	}
	switch r.Method {
	case http.MethodGet:
		d, err := s.getPlanDetail(r.Context(), uid, id)
		if err != nil {
			writeError(w, http.StatusNotFound, errors.New("plan not found"))
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: d})
	case http.MethodDelete:
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		_, err := s.db.ExecContext(ctx, `DELETE FROM user_visa_plans WHERE id=? AND user_id=?`, id, uid)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "deleted"})
	default:
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}

func (s *appServer) handlePlanTaskPatch(w http.ResponseWriter, r *http.Request, uid, planID, taskID int64) {
	if r.Method != http.MethodPatch {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	var in struct {
		Status     string `json:"status"`
		StatusText string `json:"statusText"`
	}
	if err := readJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	status := emptyFallback(in.Status, "done")
	statusText := strings.TrimSpace(in.StatusText)
	if statusText == "" {
		statusText = defaultStatusText(status)
	}
	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()
	var visaID int64
	if err := s.db.QueryRowContext(ctx, `SELECT visa_id FROM user_visa_plans WHERE id=? AND user_id=?`, planID, uid).Scan(&visaID); err != nil {
		writeError(w, http.StatusNotFound, errors.New("plan not found"))
		return
	}
	if err := s.ensurePlanSnapshot(ctx, planID, visaID); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	res, err := s.db.ExecContext(ctx, `
		UPDATE user_plan_step_tasks t
		JOIN user_plan_steps s ON s.id=t.plan_step_id
		SET t.status=?, t.status_text=?
		WHERE t.id=? AND s.plan_id=?
	`, status, statusText, taskID, planID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		writeError(w, http.StatusNotFound, errors.New("task not found"))
		return
	}
	if err := s.refreshPlanProgress(ctx, planID); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	d, err := s.getPlanDetail(ctx, uid, planID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Message: "updated", Data: d})
}

func (s *appServer) handlePlanResultPatch(w http.ResponseWriter, r *http.Request, uid, planID int64) {
	if r.Method != http.MethodPatch {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	var in struct {
		ResultStatus string `json:"resultStatus"`
		ResultNote   string `json:"resultNote"`
	}
	if err := readJSON(r, &in); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	resultStatus := strings.TrimSpace(in.ResultStatus)
	valid := map[string]bool{"pending": true, "approved": true, "rejected": true, "withdrawn": true}
	if !valid[resultStatus] {
		writeError(w, http.StatusBadRequest, errors.New("invalid resultStatus"))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
	defer cancel()
	var progress int
	var visaID int64
	var currentResultStatus string
	if err := s.db.QueryRowContext(ctx, `SELECT progress,visa_id,result_status FROM user_visa_plans WHERE id=? AND user_id=?`, planID, uid).Scan(&progress, &visaID, &currentResultStatus); err != nil {
		writeError(w, http.StatusNotFound, errors.New("plan not found"))
		return
	}
	if err := s.ensurePlanSnapshot(ctx, planID, visaID); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if strings.TrimSpace(currentResultStatus) != "" && strings.TrimSpace(currentResultStatus) != "pending" {
		writeError(w, http.StatusBadRequest, errors.New("plan result is already locked"))
		return
	}
	if progress < 100 && resultStatus != "pending" {
		writeError(w, http.StatusBadRequest, errors.New("plan result can only be chosen after progress reaches 100"))
		return
	}

	var resultAt any = nil
	if resultStatus != "pending" {
		resultAt = time.Now()
	}
	if _, err := s.db.ExecContext(ctx, `UPDATE user_visa_plans SET result_status=?, result_note=?, result_at=? WHERE id=? AND user_id=?`,
		resultStatus, strings.TrimSpace(in.ResultNote), resultAt, planID, uid); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if err := s.updatePlanLifecycleStatus(ctx, planID); err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	detail, err := s.getPlanDetail(ctx, uid, planID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Message: "updated", Data: detail})
}

func (s *appServer) getPlanDetail(ctx context.Context, uid, planID int64) (planDetail, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	var visaID int64
	if err := s.db.QueryRowContext(ctx, `SELECT visa_id FROM user_visa_plans WHERE id=? AND user_id=?`, planID, uid).Scan(&visaID); err != nil {
		return planDetail{}, err
	}
	if err := s.ensurePlanSnapshot(ctx, planID, visaID); err != nil {
		return planDetail{}, err
	}
	if err := s.refreshPlanProgress(ctx, planID); err != nil {
		return planDetail{}, err
	}

	var p planDetail
	var created time.Time
	var resultAt sql.NullTime
	err := s.db.QueryRowContext(ctx, `
		SELECT p.id,p.country_id,p.country_name,c.flag,p.visa_id,p.visa_title,p.progress,p.active_step_key,p.status,p.result_status,p.result_note,p.result_at,p.created_at
		FROM user_visa_plans p
		LEFT JOIN countries c ON c.id=p.country_id
		WHERE p.id=? AND p.user_id=?
	`, planID, uid).Scan(&p.ID, &p.CountryID, &p.CountryName, &p.CountryFlag, &p.VisaID, &p.VisaTitle, &p.Progress, &p.ActiveStepKey, &p.Status, &p.ResultStatus, &p.ResultNote, &resultAt, &created)
	if err != nil {
		return planDetail{}, err
	}
	p.CreatedAt = created.Format(time.RFC3339)
	if resultAt.Valid {
		p.ResultAt = resultAt.Time.Format(time.RFC3339)
	}
	steps, err := s.loadPlanSnapshotSteps(ctx, p.ID)
	if err != nil {
		return planDetail{}, err
	}
	p.Steps = steps
	p.Tips = []string{
		fmt.Sprintf("%s\u7b7e\u8bc1\u9ad8\u5cf0\u671f\u5efa\u8bae\u81f3\u5c11\u63d0\u524d4\u5468\u51c6\u5907\u3002", p.CountryName),
		"\u94f6\u884c\u6d41\u6c34\u5efa\u8bae\u8986\u76d6\u672c\u6b21\u884c\u7a0b\u9884\u7b97\u5e76\u4fdd\u6301\u7a33\u5b9a\u3002",
		"\u6240\u6709\u5173\u952e\u4fe1\u606f\u9700\u4e0e\u7533\u8bf7\u8868\u4fdd\u6301\u4e00\u81f4\u3002",
	}
	return p, nil
}

func (s *appServer) refreshPlanProgress(ctx context.Context, planID int64) error {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()
	hasSnapshot, err := s.hasPlanSnapshot(ctx, planID)
	if err != nil {
		return err
	}
	query := `
		SELECT s.step_key,t.id,t.status
		FROM user_plan_steps s
		JOIN user_plan_step_tasks t ON t.plan_step_id=s.id
		WHERE s.plan_id=?
		ORDER BY s.sort_order,t.sort_order,t.id
	`
	args := []any{planID}
	if !hasSnapshot {
		var visaID int64
		if err := s.db.QueryRowContext(ctx, `SELECT visa_id FROM user_visa_plans WHERE id=?`, planID).Scan(&visaID); err != nil {
			return err
		}
		query = `
			SELECT vs.step_key,t.id,COALESCE(st.status,t.default_status)
			FROM visa_steps vs
			JOIN visa_step_tasks t ON t.step_id=vs.id
			LEFT JOIN user_plan_task_states st ON st.plan_id=? AND st.task_id=t.id
			WHERE vs.visa_id=?
			ORDER BY vs.sort_order,t.sort_order,t.id
		`
		args = []any{planID, visaID}
	}
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()
	total := 0
	done := 0
	stepOrder := []string{}
	undone := map[string]bool{}
	for rows.Next() {
		var stepKey string
		var taskID int64
		var status string
		if err := rows.Scan(&stepKey, &taskID, &status); err != nil {
			return err
		}
		total++
		if status == "done" {
			done++
		} else {
			undone[stepKey] = true
		}
		if len(stepOrder) == 0 || stepOrder[len(stepOrder)-1] != stepKey {
			stepOrder = append(stepOrder, stepKey)
		}
	}
	active := "apply"
	for _, k := range stepOrder {
		if undone[k] {
			active = k
			break
		}
		active = k
	}
	progress := 0
	if total > 0 {
		progress = int(float64(done) / float64(total) * 100)
	}
	if _, err := s.db.ExecContext(ctx, `UPDATE user_visa_plans SET progress=?,active_step_key=? WHERE id=?`, progress, active, planID); err != nil {
		return err
	}
	return s.updatePlanLifecycleStatus(ctx, planID)
}

func membershipPlan(planKey string) (string, int, bool) {
	switch strings.TrimSpace(planKey) {
	case "month":
		return "\u6708\u5ea6\u4f1a\u5458", 1, true
	case "season":
		return "\u5b63\u5ea6\u4f1a\u5458", 3, true
	case "year":
		return "\u5e74\u5ea6\u4f1a\u5458", 12, true
	default:
		return "", 0, false
	}
}

func membershipStatusByExpire(expire time.Time) string {
	if expire.After(time.Now()) {
		return "active"
	}
	return "expired"
}

func (s *appServer) updateUserCore(ctx context.Context, userID int64, in userProfile, password string) error {
	if strings.TrimSpace(password) != "" && len(strings.TrimSpace(password)) < 8 {
		return errors.New("password must be at least 8 chars")
	}

	args := []any{
		emptyFallback(in.Nickname, "\u7528\u6237"),
		strings.TrimSpace(in.Email),
		strings.TrimSpace(in.Phone),
		strings.TrimSpace(in.Bio),
		strings.TrimSpace(in.Gender),
		strings.TrimSpace(in.Location),
		strings.TrimSpace(in.Avatar),
	}

	query := `UPDATE users SET nickname=?,email=?,phone=?,bio=?,gender=?,location=?,avatar=?`
	if strings.TrimSpace(password) != "" {
		hash, _ := bcrypt.GenerateFromPassword([]byte(strings.TrimSpace(password)), bcrypt.DefaultCost)
		query += `,password_hash=?`
		args = append(args, string(hash))
	}
	query += ` WHERE id=?`
	args = append(args, userID)

	_, err := s.db.ExecContext(ctx, query, args...)
	return err
}

func saveUploadedAvatar(uid int64, originalName string, src io.Reader) (string, error) {
	return saveUploadedImage("avatars", fmt.Sprintf("user-%d", uid), originalName, src)
}

func saveUploadedImage(folder, prefix, originalName string, src io.Reader) (string, error) {
	ext := strings.ToLower(strings.TrimSpace(filepath.Ext(originalName)))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
	default:
		ext = ".jpg"
	}
	folder = strings.Trim(strings.TrimSpace(folder), "/\\")
	if folder == "" {
		folder = "misc"
	}
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		prefix = "file"
	}
	dir := filepath.Join("uploads", folder)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	filename := fmt.Sprintf("%s-%d%s", prefix, time.Now().UnixNano(), ext)
	fullpath := filepath.Join(dir, filename)
	dst, err := os.Create(fullpath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}
	return "/uploads/" + folder + "/" + filename, nil
}

func buildPublicURL(r *http.Request, path string) string {
	host := strings.TrimSpace(r.Header.Get("X-Forwarded-Host"))
	if host == "" {
		host = strings.TrimSpace(r.Host)
	}
	if host == "" {
		host = "127.0.0.1:8080"
	}
	if strings.Contains(host, ",") {
		host = strings.TrimSpace(strings.Split(host, ",")[0])
	}

	scheme := strings.TrimSpace(r.Header.Get("X-Forwarded-Proto"))
	if strings.Contains(scheme, ",") {
		scheme = strings.TrimSpace(strings.Split(scheme, ",")[0])
	}
	if scheme == "" {
		scheme = "http"
		if r.TLS != nil {
			scheme = "https"
		}
	}

	return fmt.Sprintf("%s://%s%s", scheme, host, path)
}

func newUUIDString() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf(
		"%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		b[0], b[1], b[2], b[3],
		b[4], b[5],
		b[6], b[7],
		b[8], b[9],
		b[10], b[11], b[12], b[13], b[14], b[15],
	), nil
}

func defaultStatusText(status string) string {
	switch status {
	case "done":
		return "\u5df2\u5b8c\u6210"
	case "review":
		return "\u5ba1\u6838\u4e2d"
	case "missing":
		return "\u7f3a\u5931"
	default:
		return "\u5f85\u5904\u7406"
	}
}

func nullableDate(raw string) sql.NullTime {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return sql.NullTime{}
	}
	t, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: t, Valid: true}
}

func nullableTimeFromAny(raw string) sql.NullTime {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return sql.NullTime{}
	}
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return sql.NullTime{Time: t, Valid: true}
	}
	if t, err := time.Parse("2006-01-02", raw); err == nil {
		return sql.NullTime{Time: t, Valid: true}
	}
	return sql.NullTime{}
}

func nullableDateValue(d sql.NullTime) any {
	if !d.Valid {
		return nil
	}
	return d.Time.Format("2006-01-02")
}

func splitCSV(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return []string{}
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func emptyFallback(v, fallback string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return fallback
	}
	return v
}

func boolToInt(v bool) int {
	if v {
		return 1
	}
	return 0
}

func nullableID(v int64) any {
	if v <= 0 {
		return nil
	}
	return v
}

func isDuplicateErr(err error) bool {
	return strings.Contains(strings.ToLower(err.Error()), "duplicate")
}

func parseID(w http.ResponseWriter, path, prefix string) (int64, bool) {
	raw := strings.Trim(strings.TrimPrefix(path, prefix), "/")
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, errors.New("invalid id"))
		return 0, false
	}
	return id, true
}

func parseIDWithRest(w http.ResponseWriter, path, prefix string) (int64, string, bool) {
	raw := strings.Trim(strings.TrimPrefix(path, prefix), "/")
	parts := strings.SplitN(raw, "/", 2)
	id, err := strconv.ParseInt(parts[0], 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, errors.New("invalid id"))
		return 0, "", false
	}
	if len(parts) == 1 {
		return id, "", true
	}
	return id, strings.Trim(parts[1], "/"), true
}

func readJSON(r *http.Request, dst any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, apiResponse{Message: err.Error()})
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getEnv(key, fallback string) string {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	return v
}

func getEnvBool(key string, fallback bool) bool {
	v := strings.TrimSpace(strings.ToLower(os.Getenv(key)))
	if v == "" {
		return fallback
	}
	switch v {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return fallback
	}
}
