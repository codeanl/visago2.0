<script>
export default {
  onLaunch() {
    console.log('VisaNow uView Launch')
  },
}
</script>

<style>
page {
  min-height: 100%;
  background: #f7f8fb;
  color: #1f2329;
  font-family: -apple-system, BlinkMacSystemFont, 'PingFang SC', 'Helvetica Neue', sans-serif;
}

view,
text,
input,
textarea,
scroll-view,
image,
swiper,
swiper-item {
  box-sizing: border-box;
}

image {
  display: block;
}

.page-width {
  width: 100%;
  max-width: 430px;
  margin-left: auto;
  margin-right: auto;
}

.page-screen {
  min-height: 100vh;
  background: #f7f8fb;
}

.page-body {
  padding-left: 16px;
  padding-right: 16px;
}

.app-card {
  background: #ffffff;
  border: 1px solid #e8edf5;
  border-radius: 18px;
  box-shadow: 0 8px 22px rgba(31, 35, 41, 0.05);
}

.app-section-title {
  display: block;
  font-size: 22px;
  line-height: 1.2;
  font-weight: 800;
  color: #1f2329;
}

.app-section-sub {
  display: block;
  margin-top: 6px;
  font-size: 13px;
  line-height: 1.6;
  color: #7b8495;
}

.app-link {
  color: #1677ff;
  font-size: 14px;
  font-weight: 700;
}

.app-muted {
  color: #7b8495;
}

.app-pill-row {
  display: flex;
  flex-wrap: wrap;
}

.app-pill {
  margin-right: 8px;
  margin-bottom: 8px;
  padding: 5px 10px;
  border-radius: 999px;
  background: #edf5ff;
  color: #1677ff;
  font-size: 11px;
  line-height: 1.2;
}

.app-empty {
  padding: 24px 18px;
  text-align: center;
}

.app-empty-title {
  display: block;
  margin-top: 10px;
  font-size: 16px;
  line-height: 1.35;
  font-weight: 800;
  color: #1f2329;
}

.app-empty-desc {
  display: block;
  margin-top: 6px;
  font-size: 12px;
  line-height: 1.6;
  color: #7b8495;
}
</style>
