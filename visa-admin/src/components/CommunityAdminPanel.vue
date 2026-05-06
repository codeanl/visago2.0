<template>
  <section class="panel">
    <el-tabs v-model="activeTab">
      <el-tab-pane label="帖子管理" name="posts">
        <div class="panel-head">
          <div>
            <div class="panel-title">社区帖子</div>
            <div class="panel-sub">审核、隐藏、驳回和删除社区帖子</div>
          </div>
          <div class="toolbar toolbar--wide">
            <el-select v-model="postStatusFilter" clearable placeholder="全部状态" @change="loadPosts">
              <el-option v-for="item in postStatuses" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
            <el-select v-model="postCategoryFilter" clearable placeholder="全部分类" @change="loadPosts">
              <el-option v-for="item in categories" :key="item" :label="item" :value="item" />
            </el-select>
            <el-input v-model="postKeyword" clearable placeholder="按标题/内容搜索" @keyup.enter="loadPosts" />
          </div>
        </div>

        <el-table v-loading="loading.posts" :data="posts" border>
          <el-table-column prop="id" label="ID" width="72" />
          <el-table-column prop="category" label="分类" width="110" />
          <el-table-column prop="title" label="标题" min-width="180" show-overflow-tooltip />
          <el-table-column label="作者" width="120">
            <template #default="{ row }">{{ row.author?.nickname || '-' }}</template>
          </el-table-column>
          <el-table-column label="状态" width="110">
            <template #default="{ row }">
              <el-tag :type="statusTagType(row.status)">{{ statusLabel(row.status) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="likeCount" label="点赞" width="80" />
          <el-table-column prop="reportCount" label="举报" width="80" />
          <el-table-column prop="createdAt" label="创建时间" min-width="170" />
          <el-table-column label="操作" width="280" fixed="right">
            <template #default="{ row }">
              <el-button link type="success" @click="changePostStatus(row, 'published')">发布</el-button>
              <el-button link type="warning" @click="changePostStatus(row, 'hidden')">隐藏</el-button>
              <el-button link type="danger" @click="changePostStatus(row, 'rejected')">驳回</el-button>
              <el-button link type="danger" @click="removePost(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="评论管理" name="comments">
        <div class="panel-head">
          <div>
            <div class="panel-title">评论与回复</div>
            <div class="panel-sub">审核评论和回复内容，必要时可隐藏或删除</div>
          </div>
          <div class="toolbar toolbar--wide">
            <el-select v-model="commentStatusFilter" clearable placeholder="全部状态" @change="loadComments">
              <el-option v-for="item in postStatuses" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </div>
        </div>

        <el-table v-loading="loading.comments" :data="comments" border>
          <el-table-column prop="id" label="ID" width="72" />
          <el-table-column prop="postId" label="帖子ID" width="88" />
          <el-table-column label="作者" width="120">
            <template #default="{ row }">{{ row.author?.nickname || '-' }}</template>
          </el-table-column>
          <el-table-column prop="content" label="评论内容" min-width="240" show-overflow-tooltip />
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="statusTagType(row.status)">{{ statusLabel(row.status) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="createdAt" label="创建时间" min-width="170" />
          <el-table-column label="操作" width="260" fixed="right">
            <template #default="{ row }">
              <el-button link type="success" @click="changeCommentStatus(row, 'published')">发布</el-button>
              <el-button link type="warning" @click="changeCommentStatus(row, 'hidden')">隐藏</el-button>
              <el-button link type="danger" @click="changeCommentStatus(row, 'rejected')">驳回</el-button>
              <el-button link type="danger" @click="removeComment(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="举报处理" name="reports">
        <div class="panel-head">
          <div>
            <div class="panel-title">举报记录</div>
            <div class="panel-sub">处理用户举报，并可联动隐藏帖子</div>
          </div>
          <div class="toolbar toolbar--wide">
            <el-select v-model="reportStatusFilter" clearable placeholder="全部状态" @change="loadReports">
              <el-option label="待处理" value="open" />
              <el-option label="已解决" value="resolved" />
              <el-option label="已忽略" value="dismissed" />
            </el-select>
          </div>
        </div>

        <el-table v-loading="loading.reports" :data="reports" border>
          <el-table-column prop="id" label="ID" width="72" />
          <el-table-column prop="postTitle" label="帖子标题" min-width="180" show-overflow-tooltip />
          <el-table-column label="举报人" width="120">
            <template #default="{ row }">{{ row.reporter?.nickname || '-' }}</template>
          </el-table-column>
          <el-table-column prop="reason" label="原因" width="120" />
          <el-table-column prop="detail" label="补充说明" min-width="180" show-overflow-tooltip />
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="reportStatusTagType(row.status)">{{ reportStatusLabel(row.status) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="createdAt" label="举报时间" min-width="170" />
          <el-table-column label="操作" width="240" fixed="right">
            <template #default="{ row }">
              <el-button link type="success" @click="resolveReport(row, 'resolved', '')">已解决</el-button>
              <el-button link type="warning" @click="resolveReport(row, 'resolved', 'hidden')">隐藏帖子</el-button>
              <el-button link type="info" @click="resolveReport(row, 'dismissed', '')">忽略</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="评论举报" name="commentReports">
        <div class="panel-head">
          <div>
            <div class="panel-title">评论举报记录</div>
            <div class="panel-sub">处理评论与回复举报，并可联动隐藏评论</div>
          </div>
          <div class="toolbar toolbar--wide">
            <el-select v-model="commentReportStatusFilter" clearable placeholder="全部状态" @change="loadCommentReports">
              <el-option label="待处理" value="open" />
              <el-option label="已解决" value="resolved" />
              <el-option label="已忽略" value="dismissed" />
            </el-select>
          </div>
        </div>

        <el-table v-loading="loading.commentReports" :data="commentReports" border>
          <el-table-column prop="id" label="ID" width="72" />
          <el-table-column prop="postTitle" label="所属帖子" min-width="170" show-overflow-tooltip />
          <el-table-column label="被举报评论" min-width="220" show-overflow-tooltip>
            <template #default="{ row }">{{ row.commentContent || '-' }}</template>
          </el-table-column>
          <el-table-column label="评论作者" width="120">
            <template #default="{ row }">{{ row.commentAuthor?.nickname || '-' }}</template>
          </el-table-column>
          <el-table-column label="举报人" width="120">
            <template #default="{ row }">{{ row.reporter?.nickname || '-' }}</template>
          </el-table-column>
          <el-table-column prop="reason" label="原因" width="120" />
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="reportStatusTagType(row.status)">{{ reportStatusLabel(row.status) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="createdAt" label="举报时间" min-width="170" />
          <el-table-column label="操作" width="240" fixed="right">
            <template #default="{ row }">
              <el-button link type="success" @click="resolveCommentReport(row, 'resolved', '')">已解决</el-button>
              <el-button link type="warning" @click="resolveCommentReport(row, 'resolved', 'hidden')">隐藏评论</el-button>
              <el-button link type="info" @click="resolveCommentReport(row, 'dismissed', '')">忽略</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="关键词屏蔽" name="keywords">
        <div class="panel-head">
          <div>
            <div class="panel-title">社区敏感词</div>
            <div class="panel-sub">用于发帖前拦截或送审，满足苹果对 UGC 的过滤要求</div>
          </div>
          <div class="toolbar">
            <el-button type="primary" @click="openKeywordDialog()">新增关键词</el-button>
          </div>
        </div>

        <el-table v-loading="loading.keywords" :data="keywords" border>
          <el-table-column prop="id" label="ID" width="72" />
          <el-table-column prop="keyword" label="关键词" min-width="180" />
          <el-table-column label="动作" width="100">
            <template #default="{ row }">
              <el-tag :type="row.action === 'reject' ? 'danger' : 'warning'">{{ row.action === 'reject' ? '拦截' : '送审' }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="启用" width="90">
            <template #default="{ row }">
              <el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? '启用' : '停用' }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="updatedAt" label="更新时间" min-width="170" />
          <el-table-column label="操作" width="180" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" @click="openKeywordDialog(row)">编辑</el-button>
              <el-button link type="danger" @click="removeKeyword(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <el-dialog v-model="keywordDialog.visible" :title="keywordDialog.form.id ? '编辑关键词' : '新增关键词'" width="520px">
      <el-form :model="keywordDialog.form" label-width="100px">
        <el-form-item label="关键词">
          <el-input v-model="keywordDialog.form.keyword" />
        </el-form-item>
        <el-form-item label="动作">
          <el-select v-model="keywordDialog.form.action">
            <el-option label="送审" value="review" />
            <el-option label="拦截" value="reject" />
          </el-select>
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="keywordDialog.form.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="keywordDialog.visible = false">取消</el-button>
        <el-button type="primary" @click="saveKeyword">保存</el-button>
      </template>
    </el-dialog>
  </section>
</template>

<script setup>
import { onMounted, reactive, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { api } from '../api/client'

const activeTab = ref('posts')
const categories = ['推荐', '攻略', '问答', '签证经验', '材料模板']

const posts = ref([])
const comments = ref([])
const reports = ref([])
const commentReports = ref([])
const keywords = ref([])

const postKeyword = ref('')
const postStatusFilter = ref('')
const postCategoryFilter = ref('')
const commentStatusFilter = ref('')
const reportStatusFilter = ref('')
const commentReportStatusFilter = ref('')

const loading = reactive({
  posts: false,
  comments: false,
  reports: false,
  commentReports: false,
  keywords: false,
})

const loadedTabs = reactive({
  posts: false,
  comments: false,
  reports: false,
  commentReports: false,
  keywords: false,
})

const postStatuses = [
  { value: 'review', label: '待审核' },
  { value: 'published', label: '已发布' },
  { value: 'hidden', label: '已隐藏' },
  { value: 'rejected', label: '已驳回' },
]

const keywordDialog = reactive({
  visible: false,
  form: {
    id: 0,
    keyword: '',
    action: 'review',
    enabled: true,
  },
})

function statusLabel(status) {
  switch (status) {
    case 'published':
      return '已发布'
    case 'hidden':
      return '已隐藏'
    case 'rejected':
      return '已驳回'
    default:
      return '待审核'
  }
}

function statusTagType(status) {
  switch (status) {
    case 'published':
      return 'success'
    case 'hidden':
      return 'warning'
    case 'rejected':
      return 'danger'
    default:
      return 'info'
  }
}

function reportStatusLabel(status) {
  switch (status) {
    case 'resolved':
      return '已解决'
    case 'dismissed':
      return '已忽略'
    default:
      return '待处理'
  }
}

function reportStatusTagType(status) {
  switch (status) {
    case 'resolved':
      return 'success'
    case 'dismissed':
      return 'info'
    default:
      return 'warning'
  }
}

async function loadPosts() {
  loading.posts = true
  try {
    posts.value = await api.listCommunityAdminPosts({
      q: postKeyword.value,
      status: postStatusFilter.value,
      category: postCategoryFilter.value,
    })
  } finally {
    loading.posts = false
  }
}

async function loadReports() {
  loading.reports = true
  try {
    reports.value = await api.listCommunityAdminReports({
      status: reportStatusFilter.value,
    })
  } finally {
    loading.reports = false
  }
}

async function loadCommentReports() {
  loading.commentReports = true
  try {
    commentReports.value = await api.listCommunityAdminCommentReports({
      status: commentReportStatusFilter.value,
    })
  } finally {
    loading.commentReports = false
  }
}

async function loadComments() {
  loading.comments = true
  try {
    comments.value = await api.listCommunityAdminComments({
      status: commentStatusFilter.value,
    })
  } finally {
    loading.comments = false
  }
}

async function loadKeywords() {
  loading.keywords = true
  try {
    keywords.value = await api.listCommunityKeywords()
  } finally {
    loading.keywords = false
  }
}

async function ensureTabData(tab = activeTab.value, force = false) {
  if (tab === 'posts') {
    if (force || !loadedTabs.posts) {
      await loadPosts()
      loadedTabs.posts = true
    }
    return
  }
  if (tab === 'comments') {
    if (force || !loadedTabs.comments) {
      await loadComments()
      loadedTabs.comments = true
    }
    return
  }
  if (tab === 'reports') {
    if (force || !loadedTabs.reports) {
      await loadReports()
      loadedTabs.reports = true
    }
    return
  }
  if (tab === 'commentReports') {
    if (force || !loadedTabs.commentReports) {
      await loadCommentReports()
      loadedTabs.commentReports = true
    }
    return
  }
  if (tab === 'keywords') {
    if (force || !loadedTabs.keywords) {
      await loadKeywords()
      loadedTabs.keywords = true
    }
  }
}

async function changePostStatus(row, status) {
  const note = status === 'rejected' ? '不符合社区规范' : status === 'hidden' ? '已被后台隐藏' : ''
  await api.updateCommunityAdminPost(row.id, { status, reviewNote: note })
  ElMessage.success('帖子状态已更新')
  await loadPosts()
}

async function removePost(row) {
  await ElMessageBox.confirm(`确认删除帖子《${row.title}》吗？`, '提示', { type: 'warning' })
  await api.deleteCommunityAdminPost(row.id)
  ElMessage.success('帖子已删除')
  await loadPosts()
}

async function changeCommentStatus(row, status) {
  const note = status === 'rejected' ? '不符合社区规范' : status === 'hidden' ? '已被后台隐藏' : ''
  await api.updateCommunityAdminComment(row.id, { status, reviewNote: note })
  ElMessage.success('评论状态已更新')
  await loadComments()
}

async function removeComment(row) {
  await ElMessageBox.confirm(`确认删除评论「${row.content}」吗？`, '提示', { type: 'warning' })
  await api.deleteCommunityAdminComment(row.id)
  ElMessage.success('评论已删除')
  await loadComments()
}

async function resolveReport(row, status, postStatus) {
  await api.updateCommunityAdminReport(row.id, {
    status,
    postStatus,
    reviewNote: postStatus === 'hidden' ? '因举报已隐藏' : '',
  })
  ElMessage.success('举报状态已更新')
  await Promise.all([loadReports(), loadPosts()])
}

async function resolveCommentReport(row, status, commentStatus) {
  await api.updateCommunityAdminCommentReport(row.id, {
    status,
    commentStatus,
    reviewNote: commentStatus === 'hidden' ? '因举报已隐藏' : '',
  })
  ElMessage.success('评论举报状态已更新')
  await Promise.all([loadCommentReports(), loadComments()])
}

function openKeywordDialog(row) {
  keywordDialog.form = row
    ? { ...row }
    : {
        id: 0,
        keyword: '',
        action: 'review',
        enabled: true,
      }
  keywordDialog.visible = true
}

async function saveKeyword() {
  if (!keywordDialog.form.keyword) {
    ElMessage.warning('关键词必填')
    return
  }
  if (keywordDialog.form.id) {
    await api.updateCommunityKeyword(keywordDialog.form.id, keywordDialog.form)
    ElMessage.success('关键词已更新')
  } else {
    await api.createCommunityKeyword(keywordDialog.form)
    ElMessage.success('关键词已创建')
  }
  keywordDialog.visible = false
  await loadKeywords()
}

async function removeKeyword(row) {
  await ElMessageBox.confirm(`确认删除关键词 ${row.keyword} 吗？`, '提示', { type: 'warning' })
  await api.deleteCommunityKeyword(row.id)
  ElMessage.success('关键词已删除')
  await loadKeywords()
}

watch(activeTab, async (value) => {
  await ensureTabData(value, false)
})

onMounted(async () => {
  await ensureTabData(activeTab.value, false)
})

defineExpose({
  loadAll: () => Promise.all([loadPosts(), loadComments(), loadReports(), loadCommentReports(), loadKeywords()]),
  loadCurrentTab: () => ensureTabData(activeTab.value, true),
})
</script>

<style scoped>
.panel {
  height: 100%;
  min-height: 0;
  overflow: hidden;
  padding: 20px;
  border-radius: 10px;
  background: #fff;
  border: 1px solid #e5e7eb;
  box-sizing: border-box;
}

.panel-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 20px;
  margin-bottom: 16px;
  position: sticky;
  top: 0;
  z-index: 4;
  background: #fff;
  padding-bottom: 16px;
}

.panel-title {
  font-size: 18px;
  font-weight: 900;
}

.panel-sub {
  margin-top: 4px;
  color: #6b7280;
  font-size: 13px;
}

.toolbar {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 10px;
}

.toolbar .el-input {
  width: 280px;
}

.toolbar .el-select {
  width: 160px;
}

.toolbar--wide {
  flex-wrap: wrap;
}

.panel :deep(.el-tabs) {
  height: 100%;
  display: flex;
  flex-direction: column;
  min-height: 0;
}

.panel :deep(.el-tabs__content) {
  flex: 1;
  min-height: 0;
  overflow: auto;
  padding-right: 4px;
}

.panel :deep(.el-tab-pane) {
  min-height: 0;
}
</style>
