package main

import (
	"bytes"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	facebody "github.com/alibabacloud-go/facebody-20191230/v4/client"
	teaUtil "github.com/alibabacloud-go/tea-utils/v2/service"
)

const (
	photoCheckDailyLimit    = 5
	photoCheckMaxUploadSize = 20 << 20
)

type photoCheckQuotaInfo struct {
	DailyLimit int    `json:"dailyLimit"`
	UsedCount  int    `json:"usedCount"`
	Remaining  int    `json:"remaining"`
	Date       string `json:"date"`
	Configured bool   `json:"configured"`
	Provider   string `json:"provider"`
}

type photoCheckImageInfo struct {
	Width         int    `json:"width"`
	Height        int    `json:"height"`
	FileSizeKB    int    `json:"fileSizeKb"`
	FileSizeBytes int    `json:"fileSizeBytes"`
	Format        string `json:"format"`
}

type photoCheckCheckItem struct {
	Key    string `json:"key"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
	Pass   bool   `json:"pass"`
}

type photoCheckTag struct {
	Label string `json:"label"`
	Icon  string `json:"icon"`
	Type  string `json:"type"`
}

type photoCheckAnalyzeResult struct {
	Score     int                  `json:"score"`
	Summary   string               `json:"summary"`
	Country   string               `json:"country,omitempty"`
	Visa      string               `json:"visa,omitempty"`
	Provider  string               `json:"provider"`
	RequestID string               `json:"requestId,omitempty"`
	Image     photoCheckImageInfo  `json:"image"`
	Checks    []photoCheckCheckItem `json:"checks"`
	Tags      []photoCheckTag      `json:"tags"`
	Tips      []string             `json:"tips"`
	Quota     photoCheckQuotaInfo  `json:"quota"`
}

type photoCheckFaceRect struct {
	Left   int
	Top    int
	Width  int
	Height int
}

type photoCheckPose struct {
	Yaw   float64
	Pitch float64
	Roll  float64
}

type photoCheckQuality struct {
	Score float64
	Blur  float64
	Fnf   float64
	Glass float64
	Illu  float64
	Mask  float64
	Noise float64
	Pose  float64
}

type photoCheckCloudResult struct {
	RequestID       string
	FaceCount       int
	FaceProbability float64
	FaceRect        photoCheckFaceRect
	Pose            photoCheckPose
	Quality         photoCheckQuality
	Glasses         int
	Hat             int
	Mask            int
}

func (s *appServer) handlePhotoCheckQuota(w http.ResponseWriter, r *http.Request, uid int64) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	quota, err := s.getPhotoCheckQuota(r.Context(), uid)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, apiResponse{Message: "ok", Data: quota})
}

func (s *appServer) handlePhotoCheckAnalyze(w http.ResponseWriter, r *http.Request, uid int64) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}
	if !s.photoCheckConfigured() {
		writeError(w, http.StatusServiceUnavailable, errors.New("阿里云证件照检测服务尚未配置"))
		return
	}
	quota, err := s.getPhotoCheckQuota(r.Context(), uid)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	if quota.Remaining <= 0 {
		writeError(w, http.StatusTooManyRequests, errors.New("今日检测次数已用完"))
		return
	}
	if err := r.ParseMultipartForm(25 << 20); err != nil {
		writeError(w, http.StatusBadRequest, errors.New("无效的图片上传请求"))
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		writeError(w, http.StatusBadRequest, errors.New("请先上传照片"))
		return
	}
	defer file.Close()

	imageBytes, imageInfo, err := readPhotoCheckImage(file)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	country := strings.TrimSpace(r.FormValue("country"))
	visa := strings.TrimSpace(r.FormValue("visa"))

	cloudResult, err := s.analyzePhotoCheckWithAliyun(imageBytes)
	if err != nil {
		log.Printf("photo check analyze failed for user %d: %v", uid, err)
		writeError(w, http.StatusBadGateway, errors.New("阿里云证件照检测失败，请稍后重试"))
		return
	}

	result := buildPhotoCheckAnalyzeResult(country, visa, imageInfo, cloudResult)
	updatedQuota, quotaErr := s.incrementPhotoCheckUsage(r.Context(), uid)
	if quotaErr != nil {
		log.Printf("photo check usage update failed for user %d: %v", uid, quotaErr)
		updatedQuota = quota
		if updatedQuota.UsedCount < updatedQuota.DailyLimit {
			updatedQuota.UsedCount++
		}
		if updatedQuota.Remaining > 0 {
			updatedQuota.Remaining--
		}
	}
	result.Quota = updatedQuota
	writeJSON(w, http.StatusOK, apiResponse{Message: "analyzed", Data: result})
}

func (s *appServer) getPhotoCheckQuota(ctx context.Context, uid int64) (photoCheckQuotaInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	today := time.Now().Format("2006-01-02")
	var used sql.NullInt64
	err := s.db.QueryRowContext(ctx, `
		SELECT used_count
		FROM user_photo_check_daily_usage
		WHERE user_id=? AND usage_date=?
	`, uid, today).Scan(&used)
	if err != nil && err != sql.ErrNoRows {
		return photoCheckQuotaInfo{}, err
	}

	usedCount := 0
	if used.Valid && used.Int64 > 0 {
		usedCount = int(used.Int64)
	}
	if usedCount > photoCheckDailyLimit {
		usedCount = photoCheckDailyLimit
	}
	remaining := photoCheckDailyLimit - usedCount
	if remaining < 0 {
		remaining = 0
	}

	return photoCheckQuotaInfo{
		DailyLimit: photoCheckDailyLimit,
		UsedCount:  usedCount,
		Remaining:  remaining,
		Date:       today,
		Configured: s.photoCheckConfigured(),
		Provider:   "Aliyun",
	}, nil
}

func (s *appServer) incrementPhotoCheckUsage(ctx context.Context, uid int64) (photoCheckQuotaInfo, error) {
	ctx, cancel := context.WithTimeout(ctx, 8*time.Second)
	defer cancel()

	today := time.Now().Format("2006-01-02")
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO user_photo_check_daily_usage(user_id, usage_date, used_count)
		VALUES(?, ?, 1)
		ON DUPLICATE KEY UPDATE used_count = LEAST(used_count + 1, ?)
	`, uid, today, photoCheckDailyLimit)
	if err != nil {
		return photoCheckQuotaInfo{}, err
	}
	return s.getPhotoCheckQuota(ctx, uid)
}

func (s *appServer) photoCheckConfigured() bool {
	return strings.TrimSpace(s.cfg.AliyunAccessKeyID) != "" &&
		strings.TrimSpace(s.cfg.AliyunAccessKeySecret) != ""
}

func (s *appServer) analyzePhotoCheckWithAliyun(imageBytes []byte) (photoCheckCloudResult, error) {
	client, err := s.newAliyunFacebodyClient()
	if err != nil {
		return photoCheckCloudResult{}, err
	}
	request := &facebody.RecognizeFaceAdvanceRequest{}
	request.SetImageURLObject(bytes.NewReader(imageBytes))
	request.SetGlass(true)
	request.SetHat(true)
	request.SetMask(true)
	request.SetQuality(true)

	response, err := client.RecognizeFaceAdvance(request, &teaUtil.RuntimeOptions{})
	if err != nil {
		return photoCheckCloudResult{}, err
	}
	if response == nil || response.Body == nil || response.Body.Data == nil {
		return photoCheckCloudResult{}, errors.New("empty recognize face response")
	}
	data := response.Body.Data
	result := photoCheckCloudResult{
		RequestID:       stringValue(response.Body.RequestId),
		FaceCount:       int32Value(data.FaceCount),
		FaceProbability: firstFloat32(data.FaceProbabilityList),
		FaceRect:        firstFaceRect(data.FaceRectangles),
		Pose:            firstPose(data.PoseList),
		Quality:         firstQuality(data.Qualities),
		Glasses:         firstInt32SliceValue(data.Glasses),
		Hat:             firstInt32SliceValue(data.HatList),
		Mask:            firstInt64SliceValue(data.Masks),
	}
	return result, nil
}

func (s *appServer) newAliyunFacebodyClient() (*facebody.Client, error) {
	region := strings.TrimSpace(s.cfg.AliyunViapiRegion)
	if region == "" {
		region = "cn-shanghai"
	}
	cfg := &openapi.Config{
		AccessKeyId:     stringPtr(strings.TrimSpace(s.cfg.AliyunAccessKeyID)),
		AccessKeySecret: stringPtr(strings.TrimSpace(s.cfg.AliyunAccessKeySecret)),
		RegionId:        stringPtr(region),
		Protocol:        stringPtr("HTTPS"),
	}
	if endpoint := strings.TrimSpace(s.cfg.AliyunViapiEndpoint); endpoint != "" {
		cfg.Endpoint = stringPtr(endpoint)
	}
	return facebody.NewClient(cfg)
}

func buildPhotoCheckAnalyzeResult(country, visa string, imageInfo photoCheckImageInfo, cloud photoCheckCloudResult) photoCheckAnalyzeResult {
	faceRatio := 0.0
	centerOffsetX := 1.0
	centerOffsetY := 1.0
	if imageInfo.Height > 0 {
		faceRatio = float64(cloud.FaceRect.Height) / float64(imageInfo.Height)
	}
	if imageInfo.Width > 0 && cloud.FaceRect.Width > 0 {
		faceCenterX := float64(cloud.FaceRect.Left) + float64(cloud.FaceRect.Width)/2
		centerOffsetX = math.Abs(faceCenterX-float64(imageInfo.Width)/2) / float64(imageInfo.Width)
	}
	if imageInfo.Height > 0 && cloud.FaceRect.Height > 0 {
		faceCenterY := float64(cloud.FaceRect.Top) + float64(cloud.FaceRect.Height)/2
		centerOffsetY = math.Abs(faceCenterY-float64(imageInfo.Height)/2) / float64(imageInfo.Height)
	}

	singleFacePass := cloud.FaceCount == 1
	probabilityPass := cloud.FaceProbability >= 0.8
	centerPass := centerOffsetX <= 0.12 && centerOffsetY <= 0.14
	ratioPass := faceRatio >= 0.40 && faceRatio <= 0.75
	posePass := math.Abs(cloud.Pose.Yaw) <= 20 && math.Abs(cloud.Pose.Pitch) <= 20 && math.Abs(cloud.Pose.Roll) <= 15
	facePass := singleFacePass && probabilityPass && centerPass && ratioPass && posePass

	glassesPass := cloud.Glasses == 0
	hatPass := cloud.Hat == 0
	maskPass := cloud.Mask == 0
	accessoryPass := glassesPass && hatPass && maskPass

	sizePass := imageInfo.Width >= 300 && imageInfo.Height >= 300 && imageInfo.FileSizeBytes <= photoCheckMaxUploadSize
	qualityPass := cloud.Quality.Score >= 70 &&
		cloud.Quality.Blur >= 70 &&
		cloud.Quality.Noise >= 65 &&
		cloud.Quality.Illu >= 55 &&
		cloud.Quality.Fnf >= 60 &&
		cloud.Quality.Pose >= 60
	resolutionPass := qualityPass && posePass

	faceDetail := buildFaceCheckDetail(cloud, faceRatio, centerOffsetX, centerOffsetY, facePass)
	accessoryDetail := buildAccessoryCheckDetail(cloud, accessoryPass)
	sizeDetail := fmt.Sprintf("当前照片 %dx%d，约 %dKB", imageInfo.Width, imageInfo.Height, imageInfo.FileSizeKB)
	resolutionDetail := fmt.Sprintf("质量分 %.0f，清晰度 %.0f，光照 %.0f", cloud.Quality.Score, cloud.Quality.Blur, cloud.Quality.Illu)

	checks := []photoCheckCheckItem{
		{Key: "face", Title: "头像居中且比例合适", Detail: faceDetail, Pass: facePass},
		{Key: "glasses", Title: "不佩戴眼镜、帽子和口罩", Detail: accessoryDetail, Pass: accessoryPass},
		{Key: "size", Title: "尺寸与文件大小符合基础要求", Detail: sizeDetail, Pass: sizePass},
		{Key: "resolution", Title: "清晰度和光线状态良好", Detail: resolutionDetail, Pass: resolutionPass},
	}

	score := 48
	if singleFacePass {
		score += 12
	}
	if probabilityPass {
		score += 8
	}
	if centerPass {
		score += 8
	}
	if ratioPass {
		score += 8
	}
	if accessoryPass {
		score += 14
	}
	if sizePass {
		score += 10
	}
	if resolutionPass {
		score += 12
	}
	score += clampInt(int(math.Round(cloud.Quality.Score/10)), 0, 10)
	score = clampInt(score, 18, 99)

	tips := make([]string, 0, 4)
	if !facePass {
		tips = append(tips, "请保持单人出镜，正对镜头，并让头像占照片更稳定的比例。")
	}
	if !accessoryPass {
		tips = append(tips, "请摘下眼镜、帽子或口罩，避免面部关键区域被遮挡。")
	}
	if !sizePass {
		tips = append(tips, "建议上传更高分辨率的原图，并控制单张照片大小在 20MB 内。")
	}
	if !resolutionPass {
		tips = append(tips, "请在均匀光线下重新拍摄，避免模糊、逆光和明显偏头。")
	}
	if len(tips) == 0 {
		tips = append(tips, "检测结果良好，递交前仍建议结合目标国家签证照片规则做最后确认。")
	}

	summary := "照片已完成初步检测，建议根据提示调整后再提交。"
	if facePass && accessoryPass && sizePass && resolutionPass {
		summary = "照片整体符合当前证件照基础要求，可以继续做下一步人工确认。"
	}
	if strings.TrimSpace(country) != "" || strings.TrimSpace(visa) != "" {
		summary = strings.TrimSpace(strings.Join([]string{strings.TrimSpace(country), strings.TrimSpace(visa), summary}, " "))
	}

	tags := []photoCheckTag{
		buildPhotoCheckTag(facePass, "人脸状态正常", "请调整拍摄姿态"),
		buildPhotoCheckTag(accessoryPass, "面部无遮挡", "请移除面部遮挡"),
		buildPhotoCheckTag(sizePass && resolutionPass, "照片质量较好", "建议重新拍摄优化"),
	}

	return photoCheckAnalyzeResult{
		Score:     score,
		Summary:   summary,
		Country:   country,
		Visa:      visa,
		Provider:  "Aliyun",
		RequestID: cloud.RequestID,
		Image:     imageInfo,
		Checks:    checks,
		Tags:      tags,
		Tips:      tips,
	}
}

func readPhotoCheckImage(file multipart.File) ([]byte, photoCheckImageInfo, error) {
	data, err := io.ReadAll(io.LimitReader(file, photoCheckMaxUploadSize+1))
	if err != nil {
		return nil, photoCheckImageInfo{}, err
	}
	if len(data) == 0 {
		return nil, photoCheckImageInfo{}, errors.New("请上传有效的照片文件")
	}
	if len(data) > photoCheckMaxUploadSize {
		return nil, photoCheckImageInfo{}, errors.New("照片文件过大，请控制在 20MB 内")
	}
	cfg, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, photoCheckImageInfo{}, errors.New("暂不支持该图片格式，请上传 JPG 或 PNG 照片")
	}
	return data, photoCheckImageInfo{
		Width:         cfg.Width,
		Height:        cfg.Height,
		FileSizeKB:    int(math.Ceil(float64(len(data)) / 1024)),
		FileSizeBytes: len(data),
		Format:        format,
	}, nil
}

func buildFaceCheckDetail(cloud photoCheckCloudResult, faceRatio, offsetX, offsetY float64, pass bool) string {
	if cloud.FaceCount <= 0 {
		return "未检测到清晰人脸，请正对镜头并确保照片中只有一人。"
	}
	if cloud.FaceCount > 1 {
		return fmt.Sprintf("检测到 %d 张人脸，请改为单人证件照。", cloud.FaceCount)
	}
	if !pass {
		return fmt.Sprintf("头像高度约占照片 %.0f%%，偏移量 X %.0f%% / Y %.0f%%。", faceRatio*100, offsetX*100, offsetY*100)
	}
	return fmt.Sprintf("头像高度约占照片 %.0f%%，面部位置居中。", faceRatio*100)
}

func buildAccessoryCheckDetail(cloud photoCheckCloudResult, pass bool) string {
	if pass {
		return "未检测到明显的眼镜、帽子或口罩遮挡。"
	}
	problems := make([]string, 0, 3)
	if cloud.Glasses != 0 {
		problems = append(problems, "眼镜")
	}
	if cloud.Hat != 0 {
		problems = append(problems, "帽子")
	}
	if cloud.Mask != 0 {
		problems = append(problems, "口罩")
	}
	if len(problems) == 0 {
		return "检测到面部存在遮挡，请重新拍摄。"
	}
	return "检测到可能存在" + strings.Join(problems, "、") + "遮挡。"
}

func buildPhotoCheckTag(pass bool, successLabel, failLabel string) photoCheckTag {
	if pass {
		return photoCheckTag{
			Label: successLabel,
			Icon:  "check_circle",
			Type:  "ok",
		}
	}
	return photoCheckTag{
		Label: failLabel,
		Icon:  "warning",
		Type:  "warn",
	}
}

func firstFaceRect(values []*int32) photoCheckFaceRect {
	if len(values) < 4 {
		return photoCheckFaceRect{}
	}
	return photoCheckFaceRect{
		Left:   int32PtrValue(values[0]),
		Top:    int32PtrValue(values[1]),
		Width:  int32PtrValue(values[2]),
		Height: int32PtrValue(values[3]),
	}
}

func firstPose(values []*float32) photoCheckPose {
	pose := photoCheckPose{}
	if len(values) > 0 {
		pose.Yaw = float32PtrValue(values[0])
	}
	if len(values) > 1 {
		pose.Pitch = float32PtrValue(values[1])
	}
	if len(values) > 2 {
		pose.Roll = float32PtrValue(values[2])
	}
	return pose
}

func firstQuality(values *facebody.RecognizeFaceResponseBodyDataQualities) photoCheckQuality {
	if values == nil {
		return photoCheckQuality{}
	}
	return photoCheckQuality{
		Score: firstFloat32(values.ScoreList),
		Blur:  firstFloat32(values.BlurList),
		Fnf:   firstFloat32(values.FnfList),
		Glass: firstFloat32(values.GlassList),
		Illu:  firstFloat32(values.IlluList),
		Mask:  firstFloat32(values.MaskList),
		Noise: firstFloat32(values.NoiseList),
		Pose:  firstFloat32(values.PoseList),
	}
}

func firstFloat32(values []*float32) float64 {
	if len(values) == 0 || values[0] == nil {
		return 0
	}
	return float64(*values[0])
}

func firstInt32SliceValue(values []*int32) int {
	if len(values) == 0 || values[0] == nil {
		return 0
	}
	return int(*values[0])
}

func firstInt64SliceValue(values []*int64) int {
	if len(values) == 0 || values[0] == nil {
		return 0
	}
	return int(*values[0])
}

func int32Value(value *int32) int {
	if value == nil {
		return 0
	}
	return int(*value)
}

func int32PtrValue(value *int32) int {
	if value == nil {
		return 0
	}
	return int(*value)
}

func float32PtrValue(value *float32) float64 {
	if value == nil {
		return 0
	}
	return float64(*value)
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}

func stringPtr(value string) *string {
	return &value
}

func clampInt(value, minValue, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}
