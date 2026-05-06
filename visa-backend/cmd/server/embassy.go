package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"
)

type embassyItem struct {
	ID          int64    `json:"id"`
	Country     string   `json:"country"`
	CountryCode string   `json:"countryCode"`
	Flag        string   `json:"flag"`
	Region      string   `json:"region"`
	Name        string   `json:"name"`
	City        string   `json:"city"`
	Distance    string   `json:"distance"`
	Address     string   `json:"address"`
	Phone       string   `json:"phone"`
	Hours       string   `json:"hours"`
	Services    []string `json:"services"`
	Image       string   `json:"image"`
	Latitude    float64  `json:"latitude"`
	Longitude   float64  `json:"longitude"`
	Enabled     bool     `json:"enabled"`
	Keywords    []string `json:"keywords"`
}

func (s *appServer) handleEmbassies(w http.ResponseWriter, r *http.Request) {
	if !s.requireAdminForWrite(w, r) {
		return
	}
	switch r.Method {
	case http.MethodGet:
		q := strings.TrimSpace(r.URL.Query().Get("q"))
		region := strings.TrimSpace(r.URL.Query().Get("region"))
		enabledRaw := strings.TrimSpace(r.URL.Query().Get("enabled"))
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		rows, err := s.db.QueryContext(ctx, `
			SELECT id,country_name,country_code,flag,region,name,city,distance,address,phone,hours,services,image,latitude,longitude,enabled,keywords
			FROM embassies
			WHERE (?='' OR region=?)
			  AND (?='' OR country_name LIKE CONCAT('%',?,'%') OR country_code LIKE CONCAT('%',?,'%') OR name LIKE CONCAT('%',?,'%') OR city LIKE CONCAT('%',?,'%') OR address LIKE CONCAT('%',?,'%') OR services LIKE CONCAT('%',?,'%') OR keywords LIKE CONCAT('%',?,'%'))
			  AND (?='' OR enabled = CASE WHEN ? IN ('1','true','TRUE') THEN 1 ELSE 0 END)
			ORDER BY enabled DESC, region ASC, city ASC, id DESC
		`, region, region, q, q, q, q, q, q, q, q, enabledRaw, enabledRaw)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		defer rows.Close()

		items := make([]embassyItem, 0)
		for rows.Next() {
			item, err := scanEmbassy(rows)
			if err != nil {
				writeError(w, http.StatusInternalServerError, err)
				return
			}
			items = append(items, item)
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: items})
	case http.MethodPost:
		var in embassyItem
		if err := readJSON(r, &in); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		in = normalizeEmbassyItem(in)
		if in.Country == "" || in.Name == "" {
			writeError(w, http.StatusBadRequest, errors.New("country and name are required"))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		res, err := s.db.ExecContext(ctx, `
			INSERT INTO embassies(country_name,country_code,flag,region,name,city,distance,address,phone,hours,services,image,latitude,longitude,enabled,keywords)
			VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
		`,
			in.Country,
			in.CountryCode,
			in.Flag,
			in.Region,
			in.Name,
			in.City,
			in.Distance,
			in.Address,
			in.Phone,
			in.Hours,
			strings.Join(in.Services, ","),
			in.Image,
			in.Latitude,
			in.Longitude,
			boolToInt(in.Enabled),
			strings.Join(in.Keywords, ","),
		)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		id, _ := res.LastInsertId()
		item, err := s.getEmbassyByID(ctx, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusCreated, apiResponse{Message: "created", Data: item})
	default:
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}

func (s *appServer) handleEmbassyByID(w http.ResponseWriter, r *http.Request) {
	if !s.requireAdminForWrite(w, r) {
		return
	}
	id, ok := parseID(w, r.URL.Path, "/api/tools/embassies/")
	if !ok {
		return
	}
	switch r.Method {
	case http.MethodGet:
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		item, err := s.getEmbassyByID(ctx, id)
		if err != nil {
			writeError(w, http.StatusNotFound, errors.New("embassy not found"))
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: item})
	case http.MethodPut:
		var in embassyItem
		if err := readJSON(r, &in); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}
		in = normalizeEmbassyItem(in)
		if in.Country == "" || in.Name == "" {
			writeError(w, http.StatusBadRequest, errors.New("country and name are required"))
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		_, err := s.db.ExecContext(ctx, `
			UPDATE embassies
			SET country_name=?,country_code=?,flag=?,region=?,name=?,city=?,distance=?,address=?,phone=?,hours=?,services=?,image=?,latitude=?,longitude=?,enabled=?,keywords=?
			WHERE id=?
		`,
			in.Country,
			in.CountryCode,
			in.Flag,
			in.Region,
			in.Name,
			in.City,
			in.Distance,
			in.Address,
			in.Phone,
			in.Hours,
			strings.Join(in.Services, ","),
			in.Image,
			in.Latitude,
			in.Longitude,
			boolToInt(in.Enabled),
			strings.Join(in.Keywords, ","),
			id,
		)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		item, err := s.getEmbassyByID(ctx, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "updated", Data: item})
	case http.MethodDelete:
		ctx, cancel := context.WithTimeout(r.Context(), 8*time.Second)
		defer cancel()
		if _, err := s.db.ExecContext(ctx, `DELETE FROM embassies WHERE id=?`, id); err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, apiResponse{Message: "deleted"})
	default:
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
	}
}

func (s *appServer) getEmbassyByID(ctx context.Context, id int64) (embassyItem, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id,country_name,country_code,flag,region,name,city,distance,address,phone,hours,services,image,latitude,longitude,enabled,keywords
		FROM embassies
		WHERE id=?
	`, id)
	return scanEmbassy(row)
}

func scanEmbassy(scanner interface{ Scan(dest ...any) error }) (embassyItem, error) {
	var item embassyItem
	var services string
	var keywords string
	var enabled int
	if err := scanner.Scan(
		&item.ID,
		&item.Country,
		&item.CountryCode,
		&item.Flag,
		&item.Region,
		&item.Name,
		&item.City,
		&item.Distance,
		&item.Address,
		&item.Phone,
		&item.Hours,
		&services,
		&item.Image,
		&item.Latitude,
		&item.Longitude,
		&enabled,
		&keywords,
	); err != nil {
		return embassyItem{}, err
	}
	item.Services = splitCSV(services)
	item.Keywords = splitCSV(keywords)
	item.Enabled = enabled == 1
	if item.Flag == "" && item.CountryCode != "" {
		item.Flag = countryFlagEmoji(item.CountryCode)
	}
	return item, nil
}

func normalizeEmbassyItem(in embassyItem) embassyItem {
	in.Country = strings.TrimSpace(in.Country)
	in.CountryCode = strings.ToUpper(strings.TrimSpace(in.CountryCode))
	in.Flag = strings.TrimSpace(in.Flag)
	if in.Flag == "" && in.CountryCode != "" {
		in.Flag = countryFlagEmoji(in.CountryCode)
	}
	in.Region = strings.TrimSpace(in.Region)
	in.Name = strings.TrimSpace(in.Name)
	in.City = strings.TrimSpace(in.City)
	in.Distance = strings.TrimSpace(in.Distance)
	in.Address = strings.TrimSpace(in.Address)
	in.Phone = strings.TrimSpace(in.Phone)
	in.Hours = strings.TrimSpace(in.Hours)
	in.Image = strings.TrimSpace(in.Image)
	in.Services = normalizeCSVItems(in.Services)
	in.Keywords = normalizeCSVItems(in.Keywords)
	return in
}

func normalizeCSVItems(items []string) []string {
	out := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			out = append(out, item)
		}
	}
	return out
}

func countryFlagEmoji(code string) string {
	value := strings.ToUpper(strings.TrimSpace(code))
	if len(value) != 2 {
		return ""
	}
	runes := []rune(value)
	if runes[0] < 'A' || runes[0] > 'Z' || runes[1] < 'A' || runes[1] > 'Z' {
		return ""
	}
	return string([]rune{
		rune(127397 + runes[0]),
		rune(127397 + runes[1]),
	})
}

type embassySeed struct {
	Item          embassyItem
	LegacyCountry string
	LegacyName    string
}

func (s *appServer) seedEmbassies(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM embassies`).Scan(&count); err != nil {
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

	seeds := embassySeeds()

	for _, seed := range seeds {
		item := normalizeEmbassyItem(seed.Item)
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO embassies(country_name,country_code,flag,region,name,city,distance,address,phone,hours,services,image,latitude,longitude,enabled,keywords)
			VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
		`,
			item.Country,
			item.CountryCode,
			item.Flag,
			item.Region,
			item.Name,
			item.City,
			item.Distance,
			item.Address,
			item.Phone,
			item.Hours,
			strings.Join(item.Services, ","),
			item.Image,
			item.Latitude,
			item.Longitude,
			boolToInt(item.Enabled),
			strings.Join(item.Keywords, ","),
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (s *appServer) syncEmbassySeeds(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, seed := range embassySeeds() {
		item := normalizeEmbassyItem(seed.Item)
		var id int64
		var currentName string
		var currentCountry string
		err := tx.QueryRowContext(ctx, `
			SELECT id,name,country_name
			FROM embassies
			WHERE country_code=?
			ORDER BY id
			LIMIT 1
		`, item.CountryCode).Scan(&id, &currentName, &currentCountry)
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return err
			}
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			continue
		}

		shouldUpdate := currentName == seed.LegacyName ||
			currentCountry == seed.LegacyCountry ||
			currentName == item.Name ||
			currentCountry == item.Country ||
			strings.HasPrefix(currentName, "Embassy of")
		if !shouldUpdate {
			continue
		}

		if _, err := tx.ExecContext(ctx, `
			UPDATE embassies
			SET country_name=?,country_code=?,flag=?,region=?,name=?,city=?,distance=?,address=?,phone=?,hours=?,services=?,image=?,latitude=?,longitude=?,enabled=?,keywords=?
			WHERE id=?
		`,
			item.Country,
			item.CountryCode,
			item.Flag,
			item.Region,
			item.Name,
			item.City,
			item.Distance,
			item.Address,
			item.Phone,
			item.Hours,
			strings.Join(item.Services, ","),
			item.Image,
			item.Latitude,
			item.Longitude,
			boolToInt(item.Enabled),
			strings.Join(item.Keywords, ","),
			id,
		); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func embassySeeds() []embassySeed {
	return []embassySeed{
		{
			Item: embassyItem{
				Country:     "法国",
				CountryCode: "FR",
				Region:      "europe",
				Name:        "法国驻华大使馆",
				City:        "北京",
				Distance:    "约 1.2 km",
				Address:     "北京市朝阳区天泽路60号",
				Phone:       "+86 10 8531 2000",
				Hours:       "周一至周五 09:00-12:00",
				Services:    []string{"申根签证", "护照协助", "预约咨询"},
				Image:       "https://images.unsplash.com/photo-1502602898657-3e91760cbb34?auto=format&fit=crop&w=900&q=80",
				Latitude:    39.9486,
				Longitude:   116.4595,
				Enabled:     true,
				Keywords:    []string{"法国", "北京", "申根", "签证"},
			},
			LegacyCountry: "France",
			LegacyName:    "Embassy of France in Beijing",
		},
		{
			Item: embassyItem{
				Country:     "德国",
				CountryCode: "DE",
				Region:      "europe",
				Name:        "德国驻华大使馆",
				City:        "北京",
				Distance:    "约 0.8 km",
				Address:     "北京市朝阳区东直门外大街17号",
				Phone:       "+86 10 8532 9000",
				Hours:       "周一至周五 08:30-11:30",
				Services:    []string{"签证加急", "学生签证", "领事保护"},
				Image:       "https://images.unsplash.com/photo-1560969184-10fe8719e047?auto=format&fit=crop&w=900&q=80",
				Latitude:    39.9445,
				Longitude:   116.4477,
				Enabled:     true,
				Keywords:    []string{"德国", "北京", "学生签证", "加急"},
			},
			LegacyCountry: "Germany",
			LegacyName:    "Embassy of Germany in Beijing",
		},
		{
			Item: embassyItem{
				Country:     "意大利",
				CountryCode: "IT",
				Region:      "europe",
				Name:        "意大利驻华大使馆",
				City:        "北京",
				Distance:    "约 2.5 km",
				Address:     "北京市朝阳区三里屯东二街2号",
				Phone:       "+86 10 8532 7600",
				Hours:       "周一至周五 09:00-13:00",
				Services:    []string{"旅游签证", "商务签证"},
				Image:       "https://images.unsplash.com/photo-1529260830199-42c24126f198?auto=format&fit=crop&w=900&q=80",
				Latitude:    39.9413,
				Longitude:   116.4566,
				Enabled:     true,
				Keywords:    []string{"意大利", "北京", "旅游签证", "商务签证"},
			},
			LegacyCountry: "Italy",
			LegacyName:    "Embassy of Italy in Beijing",
		},
		{
			Item: embassyItem{
				Country:     "日本",
				CountryCode: "JP",
				Region:      "asia",
				Name:        "日本驻华大使馆",
				City:        "北京",
				Distance:    "约 1.7 km",
				Address:     "北京市朝阳区亮马桥东街1号",
				Phone:       "+86 10 8531 9800",
				Hours:       "周一至周五 09:00-12:00",
				Services:    []string{"旅游签证", "在留资格", "领事证明"},
				Image:       "https://images.unsplash.com/photo-1493976040374-85c8e12f0c0e?auto=format&fit=crop&w=900&q=80",
				Latitude:    39.9494,
				Longitude:   116.4591,
				Enabled:     true,
				Keywords:    []string{"日本", "北京", "旅游签证", "在留资格"},
			},
			LegacyCountry: "Japan",
			LegacyName:    "Embassy of Japan in Beijing",
		},
		{
			Item: embassyItem{
				Country:     "美国",
				CountryCode: "US",
				Region:      "north-america",
				Name:        "美国驻华大使馆",
				City:        "北京",
				Distance:    "约 2.1 km",
				Address:     "北京市朝阳区安家楼路55号",
				Phone:       "+86 10 8531 3000",
				Hours:       "周一至周五 08:00-17:00",
				Services:    []string{"B1/B2", "F1 学生签证", "紧急护照"},
				Image:       "https://images.unsplash.com/photo-1501466044931-62695aada8e9?auto=format&fit=crop&w=900&q=80",
				Latitude:    39.9528,
				Longitude:   116.4709,
				Enabled:     true,
				Keywords:    []string{"美国", "北京", "B1", "F1", "签证"},
			},
			LegacyCountry: "United States",
			LegacyName:    "Embassy of the United States in Beijing",
		},
		{
			Item: embassyItem{
				Country:     "新加坡",
				CountryCode: "SG",
				Region:      "asia",
				Name:        "新加坡驻华大使馆",
				City:        "北京",
				Distance:    "约 3.4 km",
				Address:     "北京市朝阳区建国门外秀水北街1号",
				Phone:       "+86 10 6532 1115",
				Hours:       "周一至周五 09:00-12:00",
				Services:    []string{"电子签", "商务通道"},
				Image:       "https://images.unsplash.com/photo-1525625293386-3f8f99389edd?auto=format&fit=crop&w=900&q=80",
				Latitude:    39.9174,
				Longitude:   116.4526,
				Enabled:     true,
				Keywords:    []string{"新加坡", "北京", "电子签", "商务"},
			},
			LegacyCountry: "Singapore",
			LegacyName:    "Embassy of Singapore in Beijing",
		},
	}
}
