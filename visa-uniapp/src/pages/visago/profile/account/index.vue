<template>
  <view class="account-page">
    <VisagoTopBar :show-back="true" :show-notice="false" page-name="账号资料" />

    <scroll-view scroll-y class="account-scroll">
      <view class="account-content visago-page-width">
        <view class="hero-avatar-wrap" @tap="chooseAvatar">
          <image class="hero-avatar" :src="form.avatar || defaultAvatar" mode="aspectFill" />
          <view class="hero-camera-chip">
            <text class="material-symbols-outlined">photo_camera</text>
          </view>
        </view>

        <view class="card-block">
          <view class="info-card">
            <view class="info-row">
              <text class="row-label">UUID</text>
              <text class="row-value row-value--fixed">{{ displayUserId }}</text>
            </view>

            <view class="info-row">
              <text class="row-label">昵称</text>
              <input v-model.trim="form.nickname" class="row-input" placeholder="请输入昵称" />
            </view>

            <view class="info-row info-row--clickable" @tap="openBioEditor">
              <text class="row-label">简介</text>
              <view class="row-right">
                <text class="row-value row-value--clip">{{ bioPreview }}</text>
                <text class="material-symbols-outlined row-arrow">chevron_right</text>
              </view>
            </view>

            <view class="info-row info-row--clickable" @tap="openGenderPicker">
              <text class="row-label">性别</text>
              <view class="row-right">
                <text class="row-value">{{ genderText }}</text>
                <text class="material-symbols-outlined row-arrow">chevron_right</text>
              </view>
            </view>

            <picker
              class="row-picker"
              mode="multiSelector"
              :range="locationRange"
              :value="locationPickerValue"
              @columnchange="onLocationColumnChange"
              @change="onLocationChange"
            >
              <view class="info-row info-row--clickable">
                <text class="row-label">所在地</text>
                <view class="row-right">
                  <text class="row-value row-value--clip">{{ locationPreview }}</text>
                  <text class="material-symbols-outlined row-arrow">chevron_right</text>
                </view>
              </view>
            </picker>
          </view>
        </view>

        <view class="card-block">
          <view class="info-card">
            <view class="info-row">
              <text class="row-label">手机号</text>
              <text class="row-value row-value--fixed">{{ form.phone || '未绑定' }}</text>
            </view>
          </view>
        </view>

        <view class="action-block">
          <view class="save-btn" :class="{ 'save-btn--disabled': submitting }" @tap="submitProfile">
            {{ submitting ? '保存中...' : '保存资料' }}
          </view>
        </view>

        <view class="logout-block">
          <view class="logout-btn" @tap="handleLogout">退出登录</view>
        </view>
      </view>
    </scroll-view>

    <view v-if="genderPickerVisible" class="sheet-mask" @tap="closeGenderPicker">
      <view class="sheet-panel" @tap.stop>
        <view class="sheet-title">选择性别</view>
        <view
          v-for="item in genderOptions"
          :key="item.value"
          class="sheet-option"
          :class="{ 'sheet-option--active': form.gender === item.value }"
          @tap="selectGender(item.value)"
        >
          {{ item.label }}
        </view>
      </view>
    </view>

    <view v-if="bioEditorVisible" class="sheet-mask" @tap="closeBioEditor">
      <view class="editor-panel" @tap.stop>
        <view class="sheet-title">修改简介</view>
        <textarea
          v-model.trim="bioDraft"
          class="editor-textarea"
          maxlength="120"
          placeholder="请输入简介"
        />
        <view class="editor-actions">
          <view class="editor-btn editor-btn--cancel" @tap="closeBioEditor">取消</view>
          <view class="editor-btn editor-btn--confirm" @tap="confirmBioEditor">确定</view>
        </view>
      </view>
    </view>
  </view>
</template>

<script>
import VisagoTopBar from '../../../../components/VisagoTopBar.vue'
import { areaList } from '@vant/area-data'
import { api } from '../../../../utils/api'
import { LOGIN_PAGE, logoutSession } from '../../../../utils/auth'
import { applyTheme, getStoredTheme } from '../../../../utils/theme'

const defaultAvatar = 'https://images.unsplash.com/photo-1544005313-94ddf0286df2?auto=format&fit=crop&w=300&q=80'
const { provinceOptions, cityMap } = buildProvinceCityData()

function createForm(profile = {}) {
  return {
    uuid: profile.uuid || '',
    nickname: profile.nickname || profile.name || '',
    bio: profile.bio || '',
    gender: profile.gender || '',
    location: profile.location || '',
    avatar: profile.avatar || '',
    phone: profile.phone || '',
    email: profile.email || '',
  }
}

function buildProvinceCityData() {
  const provinces = Object.entries((areaList && areaList.province_list) || {})
    .map(([code, name]) => ({ code: String(code), name: String(name) }))
    .sort((a, b) => Number(a.code) - Number(b.code))
  const citiesByProvince = {}
  for (const item of provinces) {
    citiesByProvince[item.code] = []
  }
  for (const [codeRaw, nameRaw] of Object.entries((areaList && areaList.city_list) || {})) {
    const code = String(codeRaw)
    const provinceCode = `${code.slice(0, 2)}0000`
    if (!citiesByProvince[provinceCode]) {
      citiesByProvince[provinceCode] = []
    }
    citiesByProvince[provinceCode].push({
      code,
      name: String(nameRaw),
    })
  }
  for (const province of provinces) {
    const list = citiesByProvince[province.code] || []
    list.sort((a, b) => Number(a.code) - Number(b.code))
    if (!list.length) {
      list.push({
        code: `${province.code.slice(0, 2)}0100`,
        name: province.name,
      })
    }
    citiesByProvince[province.code] = list
  }
  return {
    provinceOptions: provinces,
    cityMap: citiesByProvince,
  }
}

function normalizeRegionName(raw) {
  return String(raw || '')
    .trim()
    .replace(/特别行政区$/g, '')
    .replace(/维吾尔自治区$|壮族自治区$|回族自治区$/g, '')
    .replace(/自治区$/g, '')
    .replace(/省$|市$/g, '')
    .replace(/自治州$|地区$|盟$/g, '')
}

export default {
  components: {
    VisagoTopBar,
  },
  data() {
    return {
      defaultAvatar,
      loading: false,
      submitting: false,
      uploading: false,
      genderPickerVisible: false,
      bioEditorVisible: false,
      bioDraft: '',
      locationPickerValue: [0, 0],
      provinceOptions,
      cityMap,
      form: createForm(),
      genderOptions: [
        { label: '未选择', value: '' },
        { label: '男', value: 'male' },
        { label: '女', value: 'female' },
        { label: '其他', value: 'other' },
      ],
    }
  },
  computed: {
    genderText() {
      const matched = this.genderOptions.find((item) => item.value === this.form.gender)
      return matched ? matched.label : '未选择'
    },
    displayUserId() {
      const raw = String(this.form.uuid || '').trim()
      if (!raw) return '未生成'
      return raw.replace(/-/g, '').slice(0, 8).toUpperCase()
    },
    bioPreview() {
      return this.form.bio || '点击填写简介'
    },
    locationPreview() {
      return this.form.location || '请选择地区'
    },
    currentCityOptions() {
      return this.getCitiesByProvinceIndex(this.locationPickerValue[0])
    },
    locationRange() {
      return [
        this.provinceOptions.map((item) => item.name),
        this.currentCityOptions.map((item) => item.name),
      ]
    },
  },
  onShow() {
    applyTheme(getStoredTheme())
    this.loadProfile()
  },
  methods: {
    async loadProfile() {
      if (this.loading) return
      this.loading = true
      try {
        this.form = createForm(await api.me())
        this.syncLocationPickerByText(this.form.location)
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '加载失败',
          icon: 'none',
        })
      } finally {
        this.loading = false
      }
    },
    async chooseAvatar() {
      if (this.uploading) return
      try {
        const res = await this.chooseOneImage()
        const { filePath, fileObject } = this.pickAvatarFile(res)
        if (!filePath && !fileObject) {
          uni.showToast({
            title: '未读取到可上传图片',
            icon: 'none',
          })
          return
        }
        this.uploading = true
        uni.showLoading({
          title: '上传中...',
          mask: true,
        })
        const data = await api.uploadAvatar(filePath, fileObject)
        const nextAvatar =
          (data && data.profile && data.profile.avatar) ||
          (data && data.url) ||
          this.form.avatar
        this.form.avatar = nextAvatar
        uni.showToast({
          title: '头像已更新',
          icon: 'none',
        })
      } catch (error) {
        if (error && error.errMsg && String(error.errMsg).includes('cancel')) return
        uni.showToast({
          title: (error && error.message) || '上传失败',
          icon: 'none',
        })
      } finally {
        uni.hideLoading()
        this.uploading = false
      }
    },
    async chooseOneImage() {
      try {
        return await new Promise((resolve, reject) => {
          uni.chooseImage({
            count: 1,
            sizeType: ['compressed'],
            sourceType: ['album', 'camera'],
            success: resolve,
            fail: reject,
          })
        })
      } catch (error) {
        const errMsg = String((error && error.errMsg) || '')
        if (errMsg.includes('cancel')) {
          throw error
        }
        if (typeof uni.chooseMedia !== 'function') {
          throw error
        }
        return new Promise((resolve, reject) => {
          uni.chooseMedia({
            count: 1,
            mediaType: ['image'],
            sourceType: ['album', 'camera'],
            success: resolve,
            fail: reject,
          })
        })
      }
    },
    pickAvatarFile(res) {
      const tempPaths = Array.isArray(res && res.tempFilePaths) ? res.tempFilePaths : []
      const tempFiles = Array.isArray(res && res.tempFiles) ? res.tempFiles : []
      const firstFile = tempFiles.length ? tempFiles[0] : null
      const pathFromFile = firstFile && (firstFile.path || firstFile.tempFilePath || '')
      const filePath = tempPaths[0] || pathFromFile || ''
      const fileObject = firstFile && typeof firstFile === 'object' ? firstFile : null
      return { filePath, fileObject }
    },
    openGenderPicker() {
      this.genderPickerVisible = true
    },
    closeGenderPicker() {
      this.genderPickerVisible = false
    },
    selectGender(value) {
      this.form.gender = value
      this.closeGenderPicker()
    },
    normalizeIndex(value, size) {
      const n = Number(value || 0)
      if (size <= 0) return 0
      if (!Number.isFinite(n) || n < 0) return 0
      if (n >= size) return size - 1
      return Math.floor(n)
    },
    getCitiesByProvinceIndex(provinceIndex) {
      const idx = this.normalizeIndex(provinceIndex, this.provinceOptions.length)
      const province = this.provinceOptions[idx]
      if (!province) return []
      return this.cityMap[province.code] || []
    },
    onLocationColumnChange(event) {
      const column = Number(event && event.detail && event.detail.column)
      const value = Number(event && event.detail && event.detail.value)
      const next = [...this.locationPickerValue]
      if (column === 0) {
        next[0] = this.normalizeIndex(value, this.provinceOptions.length)
        next[1] = 0
      } else if (column === 1) {
        const cities = this.getCitiesByProvinceIndex(next[0])
        next[1] = this.normalizeIndex(value, cities.length)
      }
      this.locationPickerValue = next
    },
    onLocationChange(event) {
      const value = Array.isArray(event && event.detail && event.detail.value) ? event.detail.value : [0, 0]
      const provinceIndex = this.normalizeIndex(value[0], this.provinceOptions.length)
      const cities = this.getCitiesByProvinceIndex(provinceIndex)
      const cityIndex = this.normalizeIndex(value[1], cities.length)
      this.locationPickerValue = [provinceIndex, cityIndex]
      const province = this.provinceOptions[provinceIndex]
      const city = cities[cityIndex]
      if (province) {
        const cityText = city && city.name !== province.name ? ` ${city.name}` : ''
        this.form.location = `${province.name}${cityText}`
      }
    },
    syncLocationPickerByText(locationText) {
      const text = String(locationText || '').trim()
      if (!text) {
        this.locationPickerValue = [0, 0]
        return
      }
      const parts = text.split(/[\/,\s-]+/).map((item) => item.trim()).filter(Boolean)
      const provinceRaw = parts[0] || ''
      const cityRaw = parts[1] || ''
      const provinceNorm = normalizeRegionName(provinceRaw)
      let provinceIndex = this.provinceOptions.findIndex((item) => normalizeRegionName(item.name) === provinceNorm)
      if (provinceIndex < 0) {
        provinceIndex = this.provinceOptions.findIndex((item) => normalizeRegionName(item.name).includes(provinceNorm))
      }
      if (provinceIndex < 0) {
        provinceIndex = 0
      }
      const cities = this.getCitiesByProvinceIndex(provinceIndex)
      const cityNorm = normalizeRegionName(cityRaw)
      let cityIndex = cities.findIndex((item) => normalizeRegionName(item.name) === cityNorm)
      if (cityIndex < 0 && cityNorm) {
        cityIndex = cities.findIndex((item) => normalizeRegionName(item.name).includes(cityNorm))
      }
      if (cityIndex < 0) {
        cityIndex = 0
      }
      this.locationPickerValue = [provinceIndex, cityIndex]
    },
    openBioEditor() {
      this.bioDraft = this.form.bio || ''
      this.bioEditorVisible = true
    },
    closeBioEditor() {
      this.bioEditorVisible = false
    },
    confirmBioEditor() {
      this.form.bio = (this.bioDraft || '').trim()
      this.closeBioEditor()
    },
    async submitProfile() {
      if (this.submitting) return
      if (!this.form.nickname) {
        uni.showToast({
          title: '请填写昵称',
          icon: 'none',
        })
        return
      }
      this.submitting = true
      try {
        const payload = {
          name: this.form.nickname,
          nickname: this.form.nickname,
          bio: this.form.bio,
          gender: this.form.gender,
          location: this.form.location,
          avatar: this.form.avatar,
          phone: this.form.phone,
          email: this.form.email,
        }
        const profile = await api.updateProfile(payload)
        this.form = createForm(profile)
        this.syncLocationPickerByText(this.form.location)
        uni.showToast({
          title: '资料已同步更新',
          icon: 'none',
        })
      } catch (error) {
        uni.showToast({
          title: (error && error.message) || '保存失败',
          icon: 'none',
        })
      } finally {
        this.submitting = false
      }
    },
    handleLogout() {
      logoutSession()
      uni.reLaunch({ url: LOGIN_PAGE })
    },
  },
}
</script>

<style scoped>
.account-page {
  min-height: 100vh;
  background: var(--visago-bg);
  color: var(--visago-text);
}

.account-scroll {
  position: fixed;
  top: 74px;
  right: 0;
  bottom: 0;
  left: 0;
  height: auto;
}

.account-content {
  padding: 16px 14px calc(30px + var(--visago-safe-bottom));
  box-sizing: border-box;
}

.hero-avatar-wrap {
  position: relative;
  width: 108px;
  margin: 0 auto 18px;
}

.hero-avatar {
  width: 108px;
  height: 108px;
  border-radius: 50%;
  display: block;
  margin: 0 auto;
  border: 2px solid var(--visago-line);
}

.hero-camera-chip {
  position: absolute;
  right: 2px;
  bottom: 10px;
  width: 34px;
  height: 34px;
  border-radius: 50%;
  background: var(--visago-surface-soft);
  color: var(--visago-text);
  border: 1px solid var(--visago-line);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: var(--visago-shadow-card);
}

.hero-camera-chip .material-symbols-outlined {
  font-size: 18px;
}

.card-block + .card-block {
  margin-top: 18px;
}

.info-card {
  border-radius: 18px;
  background: var(--visago-surface);
  border: 1px solid var(--visago-line);
  box-shadow: var(--visago-shadow-card);
  overflow: hidden;
}

.info-row {
  min-height: 58px;
  padding: 0 18px;
  box-sizing: border-box;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.row-label {
  flex-shrink: 0;
  font-size: 15px;
  color: var(--visago-text);
  font-weight: 500;
}

.row-right {
  min-width: 0;
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 6px;
}

.row-value,
.row-input,
.row-textarea {
  font-size: 14px;
  color: var(--visago-text-muted);
  text-align: right;
}

.row-input {
  min-width: 0;
  flex: 1;
}

.row-input::placeholder {
  color: var(--visago-text-soft);
}

.row-picker {
  display: block;
}

.row-arrow {
  font-size: 18px;
  color: var(--visago-text-soft);
  flex-shrink: 0;
}

.info-row--clickable:active {
  background: var(--visago-surface-soft);
}

.row-value--fixed {
  color: var(--visago-text-muted);
}

.row-value--clip {
  max-width: 220px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.action-block {
  margin-top: 26px;
}

.save-btn,
.logout-btn {
  height: 58px;
  border-radius: 18px;
  background: var(--visago-surface);
  border: 1px solid var(--visago-line);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 17px;
  box-shadow: var(--visago-shadow-card);
}

.save-btn {
  color: var(--visago-primary);
  font-weight: 600;
}

.save-btn--disabled {
  opacity: 0.72;
}

.logout-block {
  margin-top: 18px;
}

.logout-btn {
  color: var(--visago-danger);
  font-weight: 600;
}

.sheet-mask {
  position: fixed;
  inset: 0;
  background: rgba(8, 10, 21, 0.42);
  display: flex;
  align-items: flex-end;
  justify-content: center;
  z-index: 30;
}

.sheet-panel,
.editor-panel {
  width: 100%;
  background: var(--visago-surface);
  border-top: 1px solid var(--visago-line);
  border-top-left-radius: 20px;
  border-top-right-radius: 20px;
  padding: 14px 14px calc(20px + var(--visago-safe-bottom));
  box-sizing: border-box;
}

.sheet-title {
  text-align: center;
  font-size: 15px;
  color: var(--visago-text-soft);
  padding: 8px 0 12px;
}

.sheet-option {
  height: 50px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--visago-text);
  font-size: 15px;
  background: var(--visago-surface-soft);
  border: 1px solid var(--visago-line);
}

.sheet-option + .sheet-option {
  margin-top: 10px;
}

.sheet-option--active {
  color: var(--visago-primary);
  background: var(--visago-surface-softest);
  border-color: var(--visago-primary);
}

.editor-textarea {
  width: 100%;
  height: 108px;
  border-radius: 12px;
  background: var(--visago-surface-soft);
  border: 1px solid var(--visago-line);
  padding: 12px;
  box-sizing: border-box;
  font-size: 14px;
  color: var(--visago-text);
  line-height: 1.5;
}

.editor-textarea::placeholder {
  color: var(--visago-text-soft);
}

.editor-actions {
  margin-top: 12px;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.editor-btn {
  height: 46px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 15px;
  font-weight: 600;
}

.editor-btn--cancel {
  color: var(--visago-text-muted);
  background: var(--visago-surface-soft);
  border: 1px solid var(--visago-line);
}

.editor-btn--confirm {
  color: #ffffff;
  background: var(--visago-primary);
}

:global(html.theme-dark) .sheet-mask {
  background: rgba(2, 4, 9, 0.62);
}
</style>
