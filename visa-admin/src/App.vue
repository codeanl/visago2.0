<template>
  <AdminLoginView
    v-if="!sessionReady || !adminUser"
    :loading="authLoading"
    @success="handleAdminLoginSuccess"
    @loading="authLoading = $event"
  />

  <el-container v-else class="admin-shell">
    <el-aside width="240px" class="sidebar">
      <div class="brand">
        <div class="brand-mark">V</div>
        <div>
          <div class="brand-title">Visago 管理台</div>
          <div class="brand-sub">前后端联调与签证数据管理</div>
        </div>
      </div>

      <el-menu :default-active="activeModule" class="nav-menu" @select="activeModule = $event">
        <el-menu-item index="admins">后台账号</el-menu-item>
        <el-menu-item index="users">用户模块</el-menu-item>
        <el-menu-item index="visa">签证模块</el-menu-item>
        <el-menu-item index="community">社区模块</el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="topbar">
        <div>
          <div class="page-title">{{ currentModuleTitle }}</div>
          <div class="page-sub">后端状态：{{ healthText }}</div>
        </div>
        <div class="topbar-actions">
          <div class="admin-chip">
            <span class="admin-chip__label">管理员</span>
            <span class="admin-chip__name">{{ adminUser?.nickname || adminUser?.phone || 'Admin' }}</span>
          </div>
          <el-button type="primary" plain @click="refreshCurrentModule">刷新数据</el-button>
          <el-button plain @click="handleAdminLogout">退出</el-button>
        </div>
      </el-header>

      <el-main class="main">
        <section v-if="activeModule === 'admins'" class="panel">
          <div class="panel-head">
            <div>
              <div class="panel-title">后台账号管理</div>
              <div class="panel-sub">独立维护后台管理员账号，只有登录后台账号后才能进入管理台</div>
            </div>
            <div class="toolbar">
              <el-input v-model="adminKeyword" clearable placeholder="按用户名/昵称/手机号/邮箱搜索" @keyup.enter="loadAdminAccounts" />
              <el-button type="primary" @click="openAdminAccountDialog()">新增后台账号</el-button>
            </div>
          </div>

          <el-table v-loading="loading.adminAccounts" :data="adminAccounts" border>
            <el-table-column prop="id" label="ID" width="72" />
            <el-table-column prop="username" label="用户名" min-width="130" />
            <el-table-column prop="nickname" label="昵称" min-width="120" />
            <el-table-column prop="phone" label="手机号" min-width="140" />
            <el-table-column prop="email" label="邮箱" min-width="180" />
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 'active' ? 'success' : 'info'">{{ row.status === 'active' ? '启用' : '停用' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="lastLoginAt" label="最近登录" min-width="170" />
            <el-table-column prop="createdAt" label="创建时间" min-width="170" />
            <el-table-column label="操作" width="180" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="openAdminAccountDialog(row)">编辑</el-button>
                <el-button link type="danger" @click="removeAdminAccount(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </section>

        <section v-else-if="activeModule === 'users'" class="panel">
          <div class="panel-head">
            <div>
              <div class="panel-title">用户列表</div>
              <div class="panel-sub">注册、登录、资料、会员信息都来自 MySQL</div>
            </div>
            <div class="toolbar">
              <el-input v-model="userKeyword" clearable placeholder="按用户名/姓名/手机号/邮箱搜索" @keyup.enter="loadUsers" />
              <el-button type="primary" @click="openUserDialog()">新增用户</el-button>
            </div>
          </div>

<el-table v-loading="loading.users" :data="users" border>
            <el-table-column prop="id" label="ID" width="72" />
            <el-table-column label="头像" width="88">
              <template #default="{ row }">
                <img v-if="row.avatar" :src="row.avatar" class="user-avatar-thumb" alt="avatar" />
                <span v-else>-</span>
              </template>
            </el-table-column>
            <el-table-column prop="nickname" label="昵称" min-width="120" />
            <el-table-column prop="bio" label="简介" min-width="180" show-overflow-tooltip />
            <el-table-column prop="gender" label="性别" width="90" />
            <el-table-column prop="location" label="所在地" min-width="120" />
            <el-table-column prop="phone" label="手机号" min-width="140" />
            <el-table-column prop="email" label="邮箱" min-width="180" />
            <el-table-column label="会员" min-width="230">
              <template #default="{ row }">
                <div class="membership-cell" v-if="row.membership.hasMembership">
                  <div>{{ row.membership.planName || '-' }}</div>
                  <div class="membership-meta">
                    {{ formatDate(row.membership.startedAt) }} ~ {{ formatDate(row.membership.expiresAt) }} 
                  </div>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="Actions" width="250" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="openUserDialog(row)">编辑</el-button>
                <el-button link type="warning" @click="openMembershipDialog(row)">会员</el-button>
                <el-button link type="danger" @click="removeUser(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </section>

        <CommunityAdminPanel v-else-if="activeModule === 'community'" ref="communityAdminRef" class="main-panel" />

        <section v-else class="panel">
          <el-tabs v-model="visaTab" class="visa-tabs">
            <el-tab-pane label="Countries" name="countries">
              <div class="panel-head">
                <div>
                  <div class="panel-title">国家管理</div>
                  <div class="panel-sub">维护 uniapp 国家列表与国家基础信息</div>
                </div>
                <div class="toolbar">
                  <el-input v-model="countryKeyword" clearable placeholder="按国家/代码/地区搜索" @keyup.enter="loadCountries" />
                  <el-button type="primary" @click="openCountryDialog()">新增国家</el-button>
                </div>
              </div>

              <el-table v-loading="loading.countries" :data="countries" border>
                <el-table-column prop="id" label="ID" width="72" />
                <el-table-column prop="name" label="Country" min-width="140" />
                <el-table-column prop="code" label="Code" width="90" />
                <el-table-column prop="region" label="Region" width="120" />
                <el-table-column prop="status" label="Status" width="110" />
                <el-table-column prop="note" label="Note" min-width="180" show-overflow-tooltip />
                <el-table-column label="Actions" width="240" fixed="right">
                  <template #default="{ row }">
                    <el-button link type="primary" @click="focusCountryVisas(row)">签证</el-button>
                    <el-button link type="primary" @click="openCountryDialog(row)">编辑</el-button>
                    <el-button link type="danger" @click="removeCountry(row)">删除</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </el-tab-pane>

            <el-tab-pane label="免签国家" name="freeCountries">
              <div class="panel-head">
                <div>
                  <div class="panel-title">免签/落地签国家</div>
                  <div class="panel-sub">维护免签政策国家，并可映射到已维护的具体签证</div>
                </div>
                <div class="toolbar toolbar--wide">
                  <el-input v-model="freeCountryKeyword" clearable placeholder="按国家/代码/城市搜索" @keyup.enter="loadFreeCountries" />
                  <el-button type="primary" @click="openFreeCountryDialog()">新增免签国家</el-button>
                </div>
              </div>

              <el-table v-loading="loading.freeCountries" :data="freeCountries" border>
                <el-table-column prop="id" label="ID" width="72" />
                <el-table-column prop="name" label="国家" min-width="120" />
                <el-table-column prop="code" label="代码" width="88" />
                <el-table-column prop="policyType" label="政策类型" width="130" />
                <el-table-column prop="stay" label="停留时长" width="120" />
                <el-table-column prop="city" label="城市" width="120" />
                <el-table-column label="签证映射" min-width="180">
                  <template #default="{ row }">
                    <el-tag v-if="row.supportedVisaId" type="success">{{ row.supportedCountryName ? `${row.supportedCountryName} · ${row.supportedVisaName}` : row.supportedVisaName || '已映射签证' }}</el-tag>
                    <span v-else>未映射</span>
                  </template>
                </el-table-column>
                <el-table-column label="启用" width="90">
                  <template #default="{ row }">
                    <el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? '启用' : '停用' }}</el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="note" label="说明" min-width="220" show-overflow-tooltip />
                <el-table-column label="操作" width="180" fixed="right">
                  <template #default="{ row }">
                    <el-button link type="primary" @click="openFreeCountryDialog(row)">编辑</el-button>
                    <el-button link type="danger" @click="removeFreeCountry(row)">删除</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </el-tab-pane>

            <el-tab-pane label="驻华使领馆" name="embassies">
              <div class="panel-head">
                <div>
                  <div class="panel-title">驻华使领馆名录</div>
                  <div class="panel-sub">维护 uniapp 端使用的驻华使领馆 / 签证中心联系方式与地图坐标</div>
                </div>
                <div class="toolbar toolbar--wide">
                  <el-input v-model="embassyKeyword" clearable placeholder="按国家 / 城市 / 使领馆搜索" @keyup.enter="loadEmbassies" />
                  <el-button type="primary" @click="openEmbassyDialog()">新增使领馆</el-button>
                </div>
              </div>

              <el-table v-loading="loading.embassies" :data="embassies" border>
                <el-table-column prop="id" label="ID" width="72" />
                <el-table-column prop="country" label="国家" min-width="120" />
                <el-table-column prop="city" label="城市" width="110" />
                <el-table-column prop="name" label="使领馆 / 签证中心" min-width="220" show-overflow-tooltip />
                <el-table-column label="地区" width="130">
                  <template #default="{ row }">
                    <span>{{ getEmbassyRegionLabel(row.region) }}</span>
                  </template>
                </el-table-column>
                <el-table-column prop="phone" label="电话" min-width="160" />
                <el-table-column prop="hours" label="办公时间" min-width="170" show-overflow-tooltip />
                <el-table-column label="服务标签" min-width="220">
                  <template #default="{ row }">
                    <el-tag
                      v-for="tag in row.services || []"
                      :key="tag"
                      size="small"
                      style="margin-right: 6px; margin-bottom: 6px;"
                    >
                      {{ tag }}
                    </el-tag>
                    <span v-if="!(row.services || []).length">-</span>
                  </template>
                </el-table-column>
                <el-table-column label="启用" width="90">
                  <template #default="{ row }">
                    <el-tag :type="row.enabled ? 'success' : 'info'">{{ row.enabled ? '启用' : '停用' }}</el-tag>
                  </template>
                </el-table-column>
                <el-table-column label="Actions" width="180" fixed="right">
                  <template #default="{ row }">
                    <el-button link type="primary" @click="openEmbassyDialog(row)">编辑</el-button>
                    <el-button link type="danger" @click="removeEmbassy(row)">删除</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </el-tab-pane>

            <el-tab-pane label="国家签证" name="visas">
              <div class="panel-head">
                <div>
                  <div class="panel-title">签证类型</div>
                  <div class="panel-sub">维护国家下的签证类型与基础说明</div>
                </div>
                <div class="toolbar toolbar--wide">
                  <el-select v-model="visaCountryFilter" clearable filterable placeholder="全部国家" @change="loadVisas">
                    <el-option v-for="country in countries" :key="country.id" :label="country.name" :value="country.id" />
                  </el-select>
                  <el-input v-model="visaKeyword" clearable placeholder="按签证名称/类型搜索" @keyup.enter="loadVisas" />
                  <el-button type="primary" @click="openVisaDialog()">新增签证</el-button>
                </div>
              </div>

              <el-table v-loading="loading.visas" :data="visas" border @row-click="selectVisaRow">
                <el-table-column prop="id" label="ID" width="72" />
                <el-table-column prop="countryName" label="国家" min-width="110" />
                <el-table-column prop="name" label="签证名称" min-width="170" />
                <el-table-column prop="visaType" label="签证类型" width="110" />
                <el-table-column label="是否热门" width="100">
                  <template #default="{ row }">
                    <el-tag :type="row.hot ? 'danger' : 'info'">{{ row.hot ? '是' : '否' }}</el-tag>
                  </template>
                </el-table-column>
                <el-table-column label="免签" width="90">
                  <template #default="{ row }">
                    <el-tag v-if="row.visaFree" type="success">免签</el-tag>
                    <span v-else>-</span>
                  </template>
                </el-table-column>
                <el-table-column prop="processingTime" label="办理时长" width="130" />
                <el-table-column prop="fee" label="费用" width="90" />
                <el-table-column prop="validity" label="有效期" width="110" />
                <el-table-column prop="entries" label="入境次数" width="90" />
                <el-table-column prop="status" label="状态" width="90" />
                <el-table-column label="操作" width="230" fixed="right">
                  <template #default="{ row }">
                    <el-button link type="primary" @click.stop="openVisaDialog(row)">编辑</el-button>
                    <el-button link type="warning" @click.stop="openVisaDetailEditor(row)">详情模板</el-button>
                    <el-button link type="danger" @click.stop="removeVisa(row)">删除</el-button>
                  </template>
                </el-table-column>
              </el-table>
            </el-tab-pane>

            <el-tab-pane label="签证详情模板" name="detail">
              <div class="panel-head">
                <div>
                  <div class="panel-title">签证步骤模板</div>
                  <div class="panel-sub">统一维护步骤、攻略、指南、材料、任务状态</div>
                </div>
                <div class="toolbar toolbar--wide">
                  <el-select v-model="detailVisaId" filterable placeholder="选择签证" @change="loadDetailVisa">
                    <el-option v-for="visa in visas" :key="visa.id" :label="`${visa.countryName} - ${visa.name}`" :value="visa.id" />
                  </el-select>
                  <el-button type="primary" :disabled="!detailVisaId" @click="saveVisaDetail">保存模板</el-button>
                </div>
              </div>

              <el-empty v-if="!detailVisaId" description="请先选择一个签证" />

              <div v-else class="detail-editor" v-loading="loading.detail">
                <div class="detail-header-row">
                  <div>
                    <div class="detail-name">{{ detailVisaBase.countryName }} - {{ detailVisaBase.name }}</div>
                    <div class="detail-desc">{{ detailVisaBase.description || '-' }}</div>
                  </div>
                  <el-button type="success" plain @click="addStep">新增步骤</el-button>
                </div>

                <div class="detail-overview">
                  <el-card class="overview-card" shadow="never">
                    <template #header>
                      <div class="overview-title">签证详情预览</div>
                    </template>

                    <div class="overview-grid">
                      <div class="overview-item">
                        <span class="overview-label">国家</span>
                        <span class="overview-value">{{ detailVisaBase.countryName || '-' }}</span>
                      </div>
                      <div class="overview-item">
                        <span class="overview-label">签证名称</span>
                        <span class="overview-value">{{ detailVisaBase.name || '-' }}</span>
                      </div>
                      <div class="overview-item">
                        <span class="overview-label">类型</span>
                        <span class="overview-value">{{ detailVisaBase.visaType || '-' }}</span>
                      </div>
                      <div class="overview-item">
                        <span class="overview-label">办理时长</span>
                        <span class="overview-value">{{ detailVisaBase.processingTime || '-' }}</span>
                      </div>
                      <div class="overview-item">
                        <span class="overview-label">费用</span>
                        <span class="overview-value">{{ detailVisaBase.fee || '-' }}</span>
                      </div>
                      <div class="overview-item">
                        <span class="overview-label">有效期</span>
                        <span class="overview-value">{{ detailVisaBase.validity || '-' }}</span>
                      </div>
                      <div class="overview-item">
                        <span class="overview-label">入境次数</span>
                        <span class="overview-value">{{ detailVisaBase.entries || '-' }}</span>
                      </div>
                      <div class="overview-item">
                        <span class="overview-label">状态</span>
                        <span class="overview-value">{{ detailVisaBase.status || '-' }}</span>
                      </div>
                    </div>

                    <div class="overview-intro">{{ detailVisaBase.longIntro || detailVisaBase.description || '-' }}</div>
                  </el-card>

                  <div v-if="detailSteps.length" class="preview-step-list">
                    <el-card v-for="(step, stepIndex) in detailSteps" :key="`preview-${step.localId}`" class="preview-step-card" shadow="never">
                      <div class="preview-step-head">
                        <div class="preview-step-index">第 {{ stepIndex + 1 }} 步</div>
                        <div>
                          <div class="preview-step-title">{{ step.title || `第 ${stepIndex + 1} 步` }}</div>
                          <div class="preview-step-key">{{ step.stepKey || '-' }}</div>
                        </div>
                      </div>

                      <div v-if="getStepStrategies(step).length" class="preview-block">
                        <div class="preview-block-title">办理攻略</div>
                        <div class="preview-bullet-list">
                          <div v-for="item in getStepStrategies(step)" :key="item" class="preview-bullet-item">
                            <span class="preview-bullet-dot" />
                            <span>{{ item }}</span>
                          </div>
                        </div>
                      </div>

                      <div v-if="getStepMaterials(step).length" class="preview-block">
                        <div class="preview-block-title">材料要求</div>
                        <div class="preview-bullet-list preview-bullet-list--muted">
                          <div v-for="item in getStepMaterials(step)" :key="item" class="preview-bullet-item">
                            <span class="preview-bullet-dot" />
                            <span>{{ item }}</span>
                          </div>
                        </div>
                      </div>

                      <div v-if="step.guides.length" class="preview-block">
                        <div class="preview-block-title">签证指南</div>
                        <div v-for="(guide, guideIndex) in step.guides" :key="`${guide.title}-${guideIndex}`" class="preview-guide">
                          <div class="preview-guide-title">{{ guide.title || '未命名指南' }}</div>
                          <div class="preview-guide-desc">{{ guide.desc || '-' }}</div>
                          <a v-if="guide.url" class="preview-guide-link" :href="guide.url" target="_blank" rel="noreferrer">
                            {{ guide.cta || '打开链接' }}
                          </a>
                          <img v-if="guide.image" class="preview-guide-image" :src="guide.image" alt="guide" />
                        </div>
                      </div>

                      <div v-if="step.tasks.length" class="preview-block">
                        <div class="preview-block-title">任务清单</div>
                        <div class="preview-task-list">
                          <div v-for="(task, taskIndex) in step.tasks" :key="`${task.taskKey}-${taskIndex}`" class="preview-task">
                            <span class="preview-task-title">{{ task.title || '未命名任务' }}</span>
                            <span class="preview-task-key">{{ task.taskKey || '-' }}</span>
                          </div>
                        </div>
                      </div>
                    </el-card>
                  </div>
                </div>

<div v-if="detailSteps.length" class="step-list-wrap">
                  <div class="step-list">
                  <el-card v-for="(step, stepIndex) in detailSteps" :key="step.localId" class="step-card">
                    <template #header>
                      <div class="step-card-head">
                        <div class="step-head-main" @click="toggleStep(step)">
                          <span class="material-symbols-outlined step-toggle-icon">{{ step.collapsed ? 'chevron_right' : 'expand_more' }}</span>
                          <div class="step-title-wrap">
                            <div class="step-title">第 {{ stepIndex + 1 }} 步</div>
                            <div class="step-subtitle">{{ step.title || step.stepKey || '-' }}</div>
                          </div>
                        </div>
                        <el-button link type="danger" @click.stop="removeStep(stepIndex)">删除步骤</el-button>
                      </div>
                    </template>

                    <div v-show="!step.collapsed" class="step-body">
                      <div class="step-meta-grid">
                        <section class="editor-panel editor-panel--compact">
                          <div class="editor-field">
                            <span class="editor-field__label">步骤标识</span>
                            <el-input v-model="step.stepKey" placeholder="apply/docs/book/result" />
                          </div>
                        </section>
                        <section class="editor-panel editor-panel--compact">
                          <div class="editor-field">
                            <span class="editor-field__label">步骤标题</span>
                            <el-input v-model="step.title" placeholder="请输入步骤标题" />
                          </div>
                        </section>
                        <section class="editor-panel editor-panel--compact">
                          <div class="editor-field">
                            <span class="editor-field__label">排序</span>
                            <el-input-number v-model="step.sortOrder" :min="1" />
                          </div>
                        </section>
                      </div>

                      <div class="step-content-grid">
                        <section class="editor-panel">
                          <div class="editor-panel__head">
                            <div>
                              <div class="editor-panel__title">办理攻略</div>
                              <div class="editor-panel__sub">按条维护当前步骤的关键处理建议，前端会按列表展示。</div>
                            </div>
                            <el-button plain type="primary" @click="addStepListItem(step, 'strategies')">新增攻略</el-button>
                          </div>
                          <div class="list-editor">
                            <div v-if="step.strategies.length" class="list-editor__rows">
                              <div v-for="(item, itemIndex) in step.strategies" :key="`${step.localId}-strategy-${itemIndex}`" class="list-editor__row">
                                <span class="list-editor__index">{{ itemIndex + 1 }}</span>
                                <el-input v-model="step.strategies[itemIndex]" placeholder="请输入一条办理攻略" />
                                <el-button link type="danger" @click="removeStepListItem(step, 'strategies', itemIndex)">删除</el-button>
                              </div>
                            </div>
                            <div v-else class="list-editor__empty">当前还没有办理攻略，点击右上角新增。</div>
                          </div>
                        </section>

                        <section class="editor-panel">
                          <div class="editor-panel__head">
                            <div>
                              <div class="editor-panel__title">签证指南</div>
                              <div class="editor-panel__sub">维护说明文案、外链、按钮文案和图片地址。</div>
                            </div>
                            <el-button link type="primary" @click="addGuide(step)">新增指南</el-button>
                          </div>
                          <el-table :data="step.guides" size="small" border class="editor-table">
                          <el-table-column label="标题" min-width="120">
                            <template #default="{ row }">
                              <el-input v-model="row.title" />
                            </template>
                          </el-table-column>
                          <el-table-column label="描述" min-width="180">
                            <template #default="{ row }">
                              <el-input v-model="row.desc" type="textarea" :rows="2" />
                            </template>
                          </el-table-column>
                          <el-table-column label="图片" min-width="160">
                            <template #default="{ row }">
                              <div class="guide-image-uploader">
                                <div class="guide-image-uploader__actions">
                                  <el-upload :auto-upload="false" :show-file-list="false" accept="image/*" :on-change="(file) => handleGuideImageChange(file, row)">
                                    <el-button plain type="primary" :loading="Boolean(row._uploading)">上传图片</el-button>
                                  </el-upload>
                                  <el-button v-if="row.image" link type="danger" @click="clearGuideImage(row)">移除</el-button>
                                </div>
                                <el-image v-if="row.image" :src="row.image" fit="cover" :preview-src-list="[row.image]" class="guide-image-thumb" />
                                <div v-else class="guide-image-empty">未上传图片</div>
                              </div>
                            </template>
                          </el-table-column>
                          <el-table-column label="按钮文案" min-width="120">
                            <template #default="{ row }">
                              <el-input v-model="row.cta" />
                            </template>
                          </el-table-column>
                          <el-table-column label="URL" min-width="180">
                            <template #default="{ row }">
                              <el-input v-model="row.url" />
                            </template>
                          </el-table-column>
                          <el-table-column label="" width="62">
                            <template #default="scope">
                              <el-button link type="danger" @click="removeGuide(step, scope.$index)">删除</el-button>
                            </template>
                          </el-table-column>
                          </el-table>
                        </section>

                        <section class="editor-panel">
                          <div class="editor-panel__head">
                            <div>
                              <div class="editor-panel__title">材料要求</div>
                              <div class="editor-panel__sub">逐条维护该步骤所需材料，适合在移动端按清单方式阅读。</div>
                            </div>
                            <el-button plain type="primary" @click="addStepListItem(step, 'materials')">新增材料</el-button>
                          </div>
                          <div class="list-editor">
                            <div v-if="step.materials.length" class="list-editor__rows">
                              <div v-for="(item, itemIndex) in step.materials" :key="`${step.localId}-material-${itemIndex}`" class="list-editor__row">
                                <span class="list-editor__index">{{ itemIndex + 1 }}</span>
                                <el-input v-model="step.materials[itemIndex]" placeholder="请输入一条材料要求" />
                                <el-button link type="danger" @click="removeStepListItem(step, 'materials', itemIndex)">删除</el-button>
                              </div>
                            </div>
                            <div v-else class="list-editor__empty">当前还没有材料要求，点击右上角新增。</div>
                          </div>
                        </section>

                        <section class="editor-panel">
                          <div class="editor-panel__head">
                            <div>
                              <div class="editor-panel__title">任务清单</div>
                              <div class="editor-panel__sub">维护任务标题、图标和排序，对应计划页的可勾选任务。</div>
                            </div>
                            <el-button link type="primary" @click="addTask(step)">新增任务</el-button>
                          </div>
                          <el-table :data="step.tasks" size="small" border class="editor-table">
                          <el-table-column label="标识" min-width="110">
                            <template #default="{ row }">
                              <el-input v-model="row.taskKey" />
                            </template>
                          </el-table-column>
                          <el-table-column label="标题" min-width="130">
                            <template #default="{ row }">
                              <el-input v-model="row.title" />
                            </template>
                          </el-table-column>
                          <el-table-column label="图标" min-width="100">
                            <template #default="scope">
                              <button class="task-icon-picker" type="button" @click="openTaskIconPicker(step, scope.$index)">
                                <span class="material-symbols-outlined task-icon-picker__icon">{{ scope.row.icon || 'task_alt' }}</span>
                                <span class="task-icon-picker__label">{{ getTaskIconLabel(scope.row.icon) }}</span>
                              </button>
                            </template>
                          </el-table-column>
                          <el-table-column label="" width="62">
                            <template #default="scope">
                              <el-button link type="danger" @click="removeTask(step, scope.$index)">删除</el-button>
                            </template>
                          </el-table-column>
                          </el-table>
                        </section>
                      </div>
                    </div>
                  </el-card>
                  </div>
                </div>

                <el-empty v-else description="当前没有步骤，请先新增步骤" />
              </div>
            </el-tab-pane>
          </el-tabs>
        </section>
      </el-main>
    </el-container>
  </el-container>

  <el-dialog
    v-model="adminAccountDialog.visible"
    :title="adminAccountDialog.form.id ? '编辑后台账号' : '新增后台账号'"
    width="560px"
  >
    <el-form :model="adminAccountDialog.form" label-width="100px">
      <el-form-item label="用户名">
        <el-input v-model="adminAccountDialog.form.username" placeholder="例如：admin" />
      </el-form-item>
      <el-form-item label="昵称">
        <el-input v-model="adminAccountDialog.form.nickname" placeholder="例如：系统管理员" />
      </el-form-item>
      <el-form-item label="手机号">
        <el-input v-model="adminAccountDialog.form.phone" />
      </el-form-item>
      <el-form-item label="邮箱">
        <el-input v-model="adminAccountDialog.form.email" />
      </el-form-item>
      <el-form-item :label="adminAccountDialog.form.id ? '重置密码' : '登录密码'">
        <el-input v-model="adminAccountDialog.form.password" placeholder="留空则不修改" show-password />
      </el-form-item>
      <el-form-item label="状态">
        <el-select v-model="adminAccountDialog.form.status">
          <el-option label="启用" value="active" />
          <el-option label="停用" value="disabled" />
        </el-select>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="adminAccountDialog.visible = false">取消</el-button>
      <el-button :loading="adminAccountDialog.saving" type="primary" @click="saveAdminAccount">保存</el-button>
    </template>
  </el-dialog>

  <el-dialog
    v-model="userDialog.visible"
    :title="userDialog.form.id ? '编辑用户' : '新增用户'"
    width="700px"
    @closed="resetUserDialogState"
  >
    <el-form :model="userDialog.form" label-width="96px">
      <el-form-item label="头像">
        <div class="avatar-editor">
          <el-avatar :size="72" :src="userAvatarPreview">
            {{ (userDialog.form.nickname || 'U').slice(0, 1) }}
          </el-avatar>
          <div class="avatar-editor__actions">
            <div class="avatar-editor__buttons">
              <el-upload :auto-upload="false" :show-file-list="false" accept="image/*" :on-change="handleUserAvatarChange">
                <el-button type="primary" plain>上传头像</el-button>
              </el-upload>
              <el-button
                :disabled="!userDialog.avatarFile && !userDialog.form.avatar"
                plain
                type="danger"
                @click="clearUserAvatar"
              >
                移除头像
              </el-button>
            </div>
            <div class="form-hint">支持 JPG、PNG、WEBP，点击保存后会上传到后端并更新用户头像。</div>
          </div>
        </div>
      </el-form-item>
      <el-form-item label="昵称"><el-input v-model="userDialog.form.nickname" /></el-form-item>
      <el-form-item label="简介"><el-input v-model="userDialog.form.bio" type="textarea" :rows="2" /></el-form-item>
      <el-form-item label="手机号"><el-input v-model="userDialog.form.phone" /></el-form-item>
      <el-form-item label="邮箱"><el-input v-model="userDialog.form.email" /></el-form-item>
      <el-form-item label="密码"><el-input v-model="userDialog.form.password" placeholder="留空则不修改" /></el-form-item>
      <el-form-item label="性别">
        <el-select v-model="userDialog.form.gender">
          <el-option label="未设置" value="" />
          <el-option label="男" value="male" />
          <el-option label="女" value="female" />
          <el-option label="其他" value="other" />
        </el-select>
      </el-form-item>
      <el-form-item label="所在地">
        <div class="location-editor">
          <el-cascader
            v-model="userDialog.locationParts"
            :options="regionOptions"
            class="region-cascader"
            clearable
            filterable
            placeholder="请选择省份和城市"
            @change="handleUserLocationChange"
          />
          <div class="form-hint">
            {{ userDialog.form.location ? `保存值：${userDialog.form.location}` : '将保存为：省份 城市' }}
          </div>
        </div>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="userDialog.visible = false">取消</el-button>
      <el-button :loading="userDialog.saving" type="primary" @click="saveUser">保存</el-button>
    </template>
  </el-dialog>

  <el-dialog
    v-model="membershipDialog.visible"
    :title="membershipDialog.userNickname ? `会员设置 · ${membershipDialog.userNickname}` : '会员设置'"
    width="560px"
  >
    <el-form :model="membershipDialog.form" label-width="120px">
      <el-form-item label="会员类型">
        <el-select v-model="membershipDialog.form.planKey" @change="handleMembershipPlanChange">
          <el-option v-for="item in membershipPlanOptions" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
      </el-form-item>
      <el-form-item label="会员名称">
        <el-input v-model="membershipDialog.form.planName" placeholder="例如：月度会员" />
      </el-form-item>
      <el-form-item label="开始时间">
        <el-date-picker
          v-model="membershipDialog.form.startedAt"
          type="datetime"
          format="YYYY-MM-DD HH:mm:ss"
          value-format="YYYY-MM-DDTHH:mm:ssZ"
          placeholder="请选择开始时间"
        />
      </el-form-item>
      <el-form-item label="到期时间">
        <el-date-picker
          v-model="membershipDialog.form.expiresAt"
          type="datetime"
          format="YYYY-MM-DD HH:mm:ss"
          value-format="YYYY-MM-DDTHH:mm:ssZ"
          placeholder="请选择到期时间"
        />
      </el-form-item>
      <el-form-item label="状态">
        <el-select v-model="membershipDialog.form.status">
          <el-option v-for="item in membershipStatusOptions" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="membershipDialog.visible = false">取消</el-button>
      <el-button type="primary" @click="saveMembership">保存</el-button>
    </template>
  </el-dialog>

  <el-dialog v-model="countryDialog.visible" :title="countryDialog.form.id ? 'Edit Country' : 'Create Country'" width="620px">
    <el-form :model="countryDialog.form" label-width="100px">
      <el-form-item label="Name"><el-input v-model="countryDialog.form.name" /></el-form-item>
      <el-form-item label="Code"><el-input v-model="countryDialog.form.code" /></el-form-item>
      <el-form-item label="Region"><el-input v-model="countryDialog.form.region" /></el-form-item>
      <el-form-item label="Flag"><el-input v-model="countryDialog.form.flag" /></el-form-item>
      <el-form-item label="Image"><el-input v-model="countryDialog.form.image" /></el-form-item>
      <el-form-item label="Note"><el-input v-model="countryDialog.form.note" type="textarea" :rows="2" /></el-form-item>
      <el-form-item label="Tags CSV"><el-input v-model="countryDialog.form.tagsCsv" placeholder="tourism,business" /></el-form-item>
      <el-form-item label="Keywords CSV"><el-input v-model="countryDialog.form.keywordsCsv" placeholder="japan,tokyo" /></el-form-item>
      <el-form-item label="Status">
        <el-select v-model="countryDialog.form.status">
          <el-option label="启用" value="active" />
          <el-option label="停用" value="disabled" />
        </el-select>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="countryDialog.visible = false">Cancel</el-button>
      <el-button type="primary" @click="saveCountry">Save</el-button>
    </template>
  </el-dialog>

  <el-dialog v-model="freeCountryDialog.visible" :title="freeCountryDialog.form.id ? '编辑免签国家' : '新增免签国家'" width="720px">
    <el-form :model="freeCountryDialog.form" label-width="110px">
      <el-form-item label="国家名称"><el-input v-model="freeCountryDialog.form.name" /></el-form-item>
      <el-form-item label="国家代码"><el-input v-model="freeCountryDialog.form.code" /></el-form-item>
      <el-form-item label="国旗"><el-input v-model="freeCountryDialog.form.flag" /></el-form-item>
      <el-form-item label="地区"><el-input v-model="freeCountryDialog.form.region" /></el-form-item>
      <el-form-item label="代表城市"><el-input v-model="freeCountryDialog.form.city" /></el-form-item>
      <el-form-item label="政策类型">
        <el-select v-model="freeCountryDialog.form.policyType">
          <el-option label="免签" value="免签" />
          <el-option label="落地签" value="落地签" />
          <el-option label="电子入境" value="电子入境" />
          <el-option label="免签过境" value="免签过境" />
        </el-select>
      </el-form-item>
      <el-form-item label="停留时长"><el-input v-model="freeCountryDialog.form.stay" placeholder="例如：最长30天" /></el-form-item>
      <el-form-item label="说明"><el-input v-model="freeCountryDialog.form.note" type="textarea" :rows="3" /></el-form-item>
      <el-form-item label="地图坐标 X"><el-input-number v-model="freeCountryDialog.form.mapX" :min="0" :max="100" :step="0.1" /></el-form-item>
      <el-form-item label="地图坐标 Y"><el-input-number v-model="freeCountryDialog.form.mapY" :min="0" :max="100" :step="0.1" /></el-form-item>
      <el-form-item label="映射签证">
        <el-select v-model="freeCountryDialog.form.supportedVisaId" clearable filterable placeholder="未维护可留空">
          <el-option v-for="visa in visas" :key="visa.id" :label="`${visa.countryName} - ${visa.name}`" :value="visa.id" />
        </el-select>
      </el-form-item>
      <el-form-item label="关键词 CSV"><el-input v-model="freeCountryDialog.form.keywordsCsv" placeholder="泰国,落地签,曼谷" /></el-form-item>
      <el-form-item label="启用"><el-switch v-model="freeCountryDialog.form.enabled" /></el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="freeCountryDialog.visible = false">取消</el-button>
      <el-button type="primary" @click="saveFreeCountry">保存</el-button>
    </template>
  </el-dialog>

  <el-dialog v-model="embassyDialog.visible" :title="embassyDialog.form.id ? '编辑使领馆' : '新增使领馆'" width="760px">
    <el-form :model="embassyDialog.form" label-width="120px">
      <el-form-item label="国家"><el-input v-model="embassyDialog.form.country" /></el-form-item>
      <el-form-item label="国家代码"><el-input v-model="embassyDialog.form.countryCode" placeholder="例如：FR / US / JP" /></el-form-item>
      <el-form-item label="旗帜"><el-input v-model="embassyDialog.form.flag" placeholder="可留空，按国家代码自动生成" /></el-form-item>
      <el-form-item label="地区">
        <el-select v-model="embassyDialog.form.region" filterable>
          <el-option v-for="item in embassyRegionOptions" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
      </el-form-item>
      <el-form-item label="名称"><el-input v-model="embassyDialog.form.name" /></el-form-item>
      <el-form-item label="城市"><el-input v-model="embassyDialog.form.city" /></el-form-item>
      <el-form-item label="距离"><el-input v-model="embassyDialog.form.distance" placeholder="例如：约 1.2 km" /></el-form-item>
      <el-form-item label="地址"><el-input v-model="embassyDialog.form.address" type="textarea" :rows="2" /></el-form-item>
      <el-form-item label="电话"><el-input v-model="embassyDialog.form.phone" /></el-form-item>
      <el-form-item label="办公时间"><el-input v-model="embassyDialog.form.hours" placeholder="例如：周一至周五 09:00-12:00" /></el-form-item>
      <el-form-item label="服务标签 CSV"><el-input v-model="embassyDialog.form.servicesCsv" placeholder="旅游签证,护照协助,预约咨询" /></el-form-item>
      <el-form-item label="图片"><el-input v-model="embassyDialog.form.image" /></el-form-item>
      <el-form-item label="纬度"><el-input-number v-model="embassyDialog.form.latitude" :precision="6" :step="0.000001" /></el-form-item>
      <el-form-item label="经度"><el-input-number v-model="embassyDialog.form.longitude" :precision="6" :step="0.000001" /></el-form-item>
      <el-form-item label="关键词 CSV"><el-input v-model="embassyDialog.form.keywordsCsv" placeholder="法国,北京,申根,签证" /></el-form-item>
      <el-form-item label="启用"><el-switch v-model="embassyDialog.form.enabled" /></el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="embassyDialog.visible = false">取消</el-button>
      <el-button type="primary" @click="saveEmbassy">保存</el-button>
    </template>
  </el-dialog>

  <el-dialog v-model="visaDialog.visible" :title="visaDialog.form.id ? '编辑签证' : '新增签证'" width="700px">
    <el-form :model="visaDialog.form" label-width="120px">
      <el-form-item label="所属国家">
        <el-select v-model="visaDialog.form.countryId" filterable>
          <el-option v-for="country in countries" :key="country.id" :label="country.name" :value="country.id" />
        </el-select>
      </el-form-item>
      <el-form-item label="签证名称"><el-input v-model="visaDialog.form.name" /></el-form-item>
      <el-form-item label="签证类型"><el-input v-model="visaDialog.form.visaType" /></el-form-item>
      <el-form-item label="办理时长"><el-input v-model="visaDialog.form.processingTime" /></el-form-item>
      <el-form-item label="费用"><el-input v-model="visaDialog.form.fee" /></el-form-item>
      <el-form-item label="有效期"><el-input v-model="visaDialog.form.validity" /></el-form-item>
      <el-form-item label="入境次数"><el-input v-model="visaDialog.form.entries" /></el-form-item>
      <el-form-item label="简介"><el-input v-model="visaDialog.form.description" type="textarea" :rows="2" /></el-form-item>
      <el-form-item label="详细介绍"><el-input v-model="visaDialog.form.longIntro" type="textarea" :rows="3" /></el-form-item>
      <el-form-item label="是否热门"><el-switch v-model="visaDialog.form.hot" /></el-form-item>
      <el-form-item label="免签"><el-switch v-model="visaDialog.form.visaFree" /></el-form-item>
      <el-form-item label="状态">
        <el-select v-model="visaDialog.form.status">
          <el-option label="active" value="active" />
          <el-option label="disabled" value="disabled" />
        </el-select>
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="visaDialog.visible = false">取消</el-button>
      <el-button type="primary" @click="saveVisa">保存</el-button>
    </template>
  </el-dialog>

  <el-dialog v-model="iconPicker.visible" title="选择任务图标" width="720px" @closed="resetTaskIconPicker">
    <div class="icon-picker">
      <el-input v-model="iconPicker.keyword" clearable placeholder="搜索图标名称或用途" />
      <div class="icon-picker__grid">
        <button
          v-for="item in filteredTaskIcons"
          :key="item.value"
          class="icon-picker__item"
          :class="{ 'icon-picker__item--active': item.value === iconPicker.currentIcon }"
          type="button"
          @click="selectTaskIcon(item.value)"
        >
          <span class="material-symbols-outlined icon-picker__symbol">{{ item.value }}</span>
          <span class="icon-picker__name">{{ item.label }}</span>
          <span class="icon-picker__value">{{ item.value }}</span>
        </button>
      </div>
      <div v-if="!filteredTaskIcons.length" class="icon-picker__empty">没有匹配图标，换个关键词试试。</div>
    </div>
  </el-dialog>
</template>

<script setup>
import { computed, onMounted, reactive, ref, watch } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { api } from './api/client'
import { clearAdminSession, getAdminToken, getAdminUser } from './auth'
import AdminLoginView from './components/AdminLoginView.vue'
import CommunityAdminPanel from './components/CommunityAdminPanel.vue'
import { formatProvinceCity, parseProvinceCity, regionOptions } from './data/chinaRegions'

const activeModule = ref('users')
const visaTab = ref('countries')
const healthText = ref('检测中...')
const sessionReady = ref(false)
const authLoading = ref(false)
const adminUser = ref(getAdminUser())

const currentModuleTitle = computed(() => {
  if (activeModule.value === 'admins') return '后台账号管理'
  if (activeModule.value === 'users') return '用户管理'
  if (activeModule.value === 'community') return '社区管理'
  return '签证管理'
})

const adminAccounts = ref([])
const users = ref([])
const countries = ref([])
const freeCountries = ref([])
const embassies = ref([])
const visas = ref([])
const communityAdminRef = ref(null)

const adminKeyword = ref('')
const userKeyword = ref('')
const countryKeyword = ref('')
const freeCountryKeyword = ref('')
const embassyKeyword = ref('')
const visaKeyword = ref('')
const visaCountryFilter = ref()

const detailVisaId = ref()
const detailVisaBase = ref({})
const detailSteps = ref([])

const loading = reactive({
  adminAccounts: false,
  users: false,
  countries: false,
  freeCountries: false,
  embassies: false,
  visas: false,
  detail: false,
})

const loadedModuleFlags = reactive({
  admins: false,
  users: false,
  visa: false,
})

const loadedVisaTabFlags = reactive({
  countries: false,
  freeCountries: false,
  embassies: false,
  visas: false,
  detail: false,
})

const defaultAdminAccount = () => ({
  id: 0,
  username: '',
  nickname: '',
  phone: '',
  email: '',
  password: '',
  status: 'active',
})

const defaultUser = () => ({
  username: '',
  nickname: '',
  email: '',
  phone: '',
  password: '12345678',
  bio: '',
  gender: '',
  location: '',
  avatar: '',
  role: 'user',
  status: 'active',
})

const defaultCountry = () => ({
  name: '',
  code: '',
  region: '',
  flag: '',
  image: '',
  note: '',
  status: 'active',
  tagsCsv: '',
  keywordsCsv: '',
})

const defaultFreeCountry = () => ({
  name: '',
  code: '',
  flag: '',
  region: '',
  city: '',
  policyType: '免签',
  stay: '',
  note: '',
  mapX: 0,
  mapY: 0,
  supportedVisaId: undefined,
  enabled: true,
  keywordsCsv: '',
})

const embassyRegionOptions = [
  { value: 'asia', label: '亚洲' },
  { value: 'europe', label: '欧洲' },
  { value: 'north-america', label: '北美' },
  { value: 'south-america', label: '南美' },
  { value: 'oceania', label: '大洋洲' },
  { value: 'africa', label: '非洲' },
  { value: 'middle-east', label: '中东' },
]

const defaultEmbassy = () => ({
  country: '',
  countryCode: '',
  flag: '',
  region: 'asia',
  name: '',
  city: '',
  distance: '',
  address: '',
  phone: '',
  hours: '',
  servicesCsv: '',
  image: '',
  latitude: 0,
  longitude: 0,
  enabled: true,
  keywordsCsv: '',
})

const defaultVisa = () => ({
  countryId: visaCountryFilter.value || countries.value[0]?.id,
  name: '',
  visaType: '',
  processingTime: '',
  fee: '',
  validity: '',
  entries: '',
  description: '',
  longIntro: '',
  hot: false,
  visaFree: false,
  status: 'active',
})

const membershipPlanOptions = [
  { value: 'month', label: '月度会员', defaultName: '月度会员' },
  { value: 'season', label: '季度会员', defaultName: '季度会员' },
  { value: 'year', label: '年度会员', defaultName: '年度会员' },
]

const membershipStatusOptions = [
  { value: 'active', label: '生效中' },
  { value: 'expired', label: '已过期' },
  { value: 'disabled', label: '已停用' },
]

const taskIconOptions = [
  { value: 'task_alt', label: '默认任务', keywords: ['默认', '任务', '通用'] },
  { value: 'description', label: '材料文件', keywords: ['材料', '文件', '文档'] },
  { value: 'article', label: '申请资料', keywords: ['申请', '资料', '表单'] },
  { value: 'folder_open', label: '文件夹', keywords: ['文件夹', '归档', '整理'] },
  { value: 'inventory_2', label: '资料清单', keywords: ['清单', '资料', '库存'] },
  { value: 'edit_document', label: '填写表单', keywords: ['填写', '表单', '录入'] },
  { value: 'fact_check', label: '信息核对', keywords: ['核对', '校验', '检查'] },
  { value: 'assignment', label: '任务单', keywords: ['任务', '工单', '申请单'] },
  { value: 'assignment_ind', label: '身份信息', keywords: ['身份', '个人', '信息'] },
  { value: 'badge', label: '证件证明', keywords: ['证件', '证明', '身份卡'] },
  { value: 'contact_page', label: '联系人信息', keywords: ['联系人', '联系', '信息'] },
  { value: 'contacts', label: '联系方式', keywords: ['联系方式', '电话', '邮箱'] },
  { value: 'person', label: '个人资料', keywords: ['个人', '用户', '人员'] },
  { value: 'groups', label: '同行人员', keywords: ['同行', '团队', '家庭'] },
  { value: 'photo_camera', label: '证件照片', keywords: ['照片', '拍照', '头像'] },
  { value: 'image', label: '图片附件', keywords: ['图片', '截图', '附件'] },
  { value: 'photo_library', label: '相册图片', keywords: ['相册', '图片', '图库'] },
  { value: 'cloud_upload', label: '云端上传', keywords: ['上传', '云端', '提交'] },
  { value: 'upload_file', label: '上传材料', keywords: ['上传', '材料', '附件'] },
  { value: 'download', label: '下载表格', keywords: ['下载', '表格', '模板'] },
  { value: 'print', label: '打印文件', keywords: ['打印', '纸质', '输出'] },
  { value: 'qr_code_2', label: '二维码', keywords: ['二维码', '扫码', '码'] },
  { value: 'fingerprint', label: '指纹采集', keywords: ['指纹', '生物信息', '录指纹'] },
  { value: 'calendar_month', label: '时间预约', keywords: ['预约', '时间', '日历'] },
  { value: 'event', label: '日程安排', keywords: ['日程', '安排', '事项'] },
  { value: 'event_available', label: '预约确认', keywords: ['确认', '预约', '成功'] },
  { value: 'schedule', label: '时效进度', keywords: ['时效', '进度', '时间'] },
  { value: 'hourglass_top', label: '等待处理', keywords: ['等待', '处理中', '排队'] },
  { value: 'history', label: '历史记录', keywords: ['历史', '记录', '进展'] },
  { value: 'payments', label: '支付费用', keywords: ['支付', '费用', '付款'] },
  { value: 'credit_card', label: '银行卡支付', keywords: ['银行卡', '信用卡', '支付'] },
  { value: 'account_balance_wallet', label: '钱包费用', keywords: ['钱包', '费用', '余额'] },
  { value: 'attach_money', label: '金额证明', keywords: ['金额', '资金', '财力'] },
  { value: 'receipt_long', label: '回执单据', keywords: ['回执', '单据', '收据'] },
  { value: 'account_balance', label: '官方机构', keywords: ['机构', '政府', '官方'] },
  { value: 'apartment', label: '使馆中心', keywords: ['使馆', '签证中心', '大楼'] },
  { value: 'business_center', label: '商务证明', keywords: ['商务', '公司', '单位'] },
  { value: 'work', label: '在职证明', keywords: ['在职', '工作', '公司'] },
  { value: 'school', label: '在学证明', keywords: ['在学', '学校', '学生'] },
  { value: 'local_police', label: '无犯罪证明', keywords: ['无犯罪', '公安', '证明'] },
  { value: 'verified_user', label: '审核通过', keywords: ['审核', '通过', '验证'] },
  { value: 'rule', label: '规则确认', keywords: ['规则', '确认', '要求'] },
  { value: 'check_circle', label: '完成确认', keywords: ['完成', '确认', '成功'] },
  { value: 'pending_actions', label: '待办事项', keywords: ['待办', '未完成', '事项'] },
  { value: 'warning', label: '风险提醒', keywords: ['提醒', '风险', '注意'] },
  { value: 'error', label: '错误异常', keywords: ['异常', '错误', '失败'] },
  { value: 'info', label: '信息说明', keywords: ['说明', '信息', '提示'] },
  { value: 'help', label: '帮助说明', keywords: ['帮助', '说明', '问号'] },
  { value: 'support_agent', label: '人工支持', keywords: ['客服', '人工', '支持'] },
  { value: 'mail', label: '邮件通知', keywords: ['邮件', '邮箱', '通知'] },
  { value: 'sms', label: '短信通知', keywords: ['短信', '通知', '消息'] },
  { value: 'notifications', label: '系统通知', keywords: ['通知', '提醒', '消息'] },
  { value: 'phone_in_talk', label: '电话联系', keywords: ['电话', '联系', '沟通'] },
  { value: 'public', label: '官网链接', keywords: ['官网', '网站', '链接'] },
  { value: 'language', label: '语言翻译', keywords: ['语言', '翻译', '外语'] },
  { value: 'translate', label: '资料翻译', keywords: ['翻译', '材料', '文书'] },
  { value: 'home', label: '住址信息', keywords: ['住址', '家庭', '地址'] },
  { value: 'hotel', label: '酒店订单', keywords: ['酒店', '住宿', '订单'] },
  { value: 'map', label: '行程地图', keywords: ['地图', '路线', '行程'] },
  { value: 'location_on', label: '地点定位', keywords: ['地点', '定位', '地址'] },
  { value: 'flight_takeoff', label: '行程出发', keywords: ['出发', '航班', '飞机'] },
  { value: 'luggage', label: '出行行李', keywords: ['行李', '出行', '准备'] },
  { value: 'travel', label: '旅行计划', keywords: ['旅行', '旅游', '计划'] },
]

const membershipDefaultNames = new Set(membershipPlanOptions.map((item) => item.defaultName))

const defaultMembershipForm = () => ({
  planKey: 'month',
  planName: '月度会员',
  startedAt: '',
  expiresAt: '',
  status: 'active',
})

const userDialog = reactive({
  visible: false,
  saving: false,
  form: defaultUser(),
  locationParts: [],
  avatarFile: null,
  avatarObjectUrl: '',
})
const membershipDialog = reactive({
  visible: false,
  userId: 0,
  userNickname: '',
  form: defaultMembershipForm(),
})

const adminAccountDialog = reactive({
  visible: false,
  saving: false,
  form: defaultAdminAccount(),
})
const iconPicker = reactive({
  visible: false,
  keyword: '',
  currentIcon: '',
  step: null,
  taskIndex: -1,
})
const countryDialog = reactive({ visible: false, form: defaultCountry() })
const freeCountryDialog = reactive({ visible: false, form: defaultFreeCountry() })
const embassyDialog = reactive({ visible: false, form: defaultEmbassy() })
const visaDialog = reactive({ visible: false, form: defaultVisa() })
const userAvatarPreview = computed(() => userDialog.avatarObjectUrl || userDialog.form.avatar || '')
const filteredTaskIcons = computed(() => {
  const keyword = String(iconPicker.keyword || '').trim().toLowerCase()
  if (!keyword) {
    return taskIconOptions
  }
  return taskIconOptions.filter((item) => {
    return (
      item.value.toLowerCase().includes(keyword) ||
      item.label.toLowerCase().includes(keyword) ||
      (item.keywords || []).some((entry) => String(entry).toLowerCase().includes(keyword))
    )
  })
})

function revokeUserAvatarObjectUrl() {
  if (userDialog.avatarObjectUrl) {
    URL.revokeObjectURL(userDialog.avatarObjectUrl)
    userDialog.avatarObjectUrl = ''
  }
}

function resetUserDialogState() {
  revokeUserAvatarObjectUrl()
  userDialog.form = defaultUser()
  userDialog.locationParts = []
  userDialog.avatarFile = null
  userDialog.saving = false
}

function handleUserLocationChange(value) {
  userDialog.locationParts = Array.isArray(value) ? value : []
  userDialog.form.location = formatProvinceCity(userDialog.locationParts)
}

function getMembershipDefaultName(planKey) {
  return membershipPlanOptions.find((item) => item.value === planKey)?.defaultName || ''
}

function handleMembershipPlanChange(planKey) {
  if (!membershipDialog.form.planName || membershipDefaultNames.has(membershipDialog.form.planName)) {
    membershipDialog.form.planName = getMembershipDefaultName(planKey)
  }
}

function handleUserAvatarChange(uploadFile) {
  const file = uploadFile?.raw
  if (!file) {
    return
  }
  if (!String(file.type || '').startsWith('image/')) {
    ElMessage.warning('请选择图片文件')
    return
  }

  revokeUserAvatarObjectUrl()
  userDialog.avatarFile = file
  userDialog.avatarObjectUrl = URL.createObjectURL(file)
}

function clearUserAvatar() {
  revokeUserAvatarObjectUrl()
  userDialog.avatarFile = null
  userDialog.form.avatar = ''
}

function formatDate(value) {
  if (!value) return '-'
  const d = new Date(value)
  if (Number.isNaN(d.getTime())) return String(value)
  return `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`
}

function splitCsv(value) {
  return String(value || '')
    .split(',')
    .map((v) => v.trim())
    .filter(Boolean)
}

function getEmbassyRegionLabel(region) {
  return embassyRegionOptions.find((item) => item.value === region)?.label || region || '-'
}

function splitLines(value) {
  return String(value || '')
    .split('\n')
    .map((v) => v.trim())
    .filter(Boolean)
}

function normalizeStepList(list = []) {
  return Array.isArray(list) ? list.map((item) => String(item || '')) : []
}

function getStepStrategies(step) {
  return normalizeStepList(step.strategies).map((item) => item.trim()).filter(Boolean)
}

function getStepMaterials(step) {
  return normalizeStepList(step.materials).map((item) => item.trim()).filter(Boolean)
}

function getTaskIconLabel(iconValue) {
  return taskIconOptions.find((item) => item.value === iconValue)?.label || iconValue || '选择图标'
}

function mkLocalStep(step = {}, idx = 0) {
  return {
    localId: `${Date.now()}-${Math.random()}-${idx}`,
    stepKey: step.stepKey || `step-${idx + 1}`,
    title: step.title || '',
    sortOrder: step.sortOrder || idx + 1,
    collapsed: false,
    strategies: normalizeStepList(step.strategies),
    materials: normalizeStepList(step.materials),
    guides: (step.guides || []).map((g) => ({
      title: g.title || '',
      desc: g.desc || '',
      image: g.image || '',
      cta: g.cta || '',
      url: g.url || '',
    })),
    tasks: (step.tasks || []).map((t) => ({
      taskKey: t.taskKey || '',
      title: t.title || '',
      icon: t.icon || 'task_alt',
      status: t.status || 'todo',
      statusText: t.statusText || 'todo',
      sortOrder: t.sortOrder || 0,
    })),
  }
}

function toggleStep(step) {
  step.collapsed = !step.collapsed
}

async function loadHealth() {
  try {
    await api.health()
    healthText.value = '已连接'
  } catch (error) {
    healthText.value = '未连接'
  }
}

async function ensureCurrentVisaTabData(force = false) {
  const tab = visaTab.value
  if (tab === 'countries') {
    if (force || !loadedVisaTabFlags.countries) {
      await loadCountries()
      loadedVisaTabFlags.countries = true
    }
    return
  }
  if (tab === 'freeCountries') {
    if (force || !loadedVisaTabFlags.freeCountries) {
      await loadFreeCountries()
      loadedVisaTabFlags.freeCountries = true
    }
    return
  }
  if (tab === 'embassies') {
    if (force || !loadedVisaTabFlags.embassies) {
      await loadEmbassies()
      loadedVisaTabFlags.embassies = true
    }
    return
  }
  if (tab === 'visas') {
    if (force || !loadedVisaTabFlags.visas) {
      await loadVisas()
      loadedVisaTabFlags.visas = true
    }
    return
  }
  if (tab === 'detail') {
    if (detailVisaId.value && (force || !loadedVisaTabFlags.detail)) {
      await loadDetailVisa()
      loadedVisaTabFlags.detail = true
    }
  }
}

async function ensureCurrentModuleData(force = false) {
  if (!adminUser.value) return
  if (activeModule.value === 'admins') {
    if (force || !loadedModuleFlags.admins) {
      await loadAdminAccounts()
      loadedModuleFlags.admins = true
    }
    return
  }
  if (activeModule.value === 'users') {
    if (force || !loadedModuleFlags.users) {
      await loadUsers()
      loadedModuleFlags.users = true
    }
    return
  }
  if (activeModule.value === 'visa') {
    loadedModuleFlags.visa = true
    await ensureCurrentVisaTabData(force)
  }
}

async function loadAdminAccounts() {
  loading.adminAccounts = true
  try {
    adminAccounts.value = await api.listAdminAccounts(adminKeyword.value)
  } finally {
    loading.adminAccounts = false
  }
}

async function loadUsers() {
  loading.users = true
  try {
    users.value = await api.listUsers(userKeyword.value)
  } finally {
    loading.users = false
  }
}

async function loadCountries() {
  loading.countries = true
  try {
    countries.value = await api.listCountries(countryKeyword.value)
  } finally {
    loading.countries = false
  }
}

async function loadFreeCountries() {
  loading.freeCountries = true
  try {
    freeCountries.value = await api.listFreeCountries(freeCountryKeyword.value)
  } finally {
    loading.freeCountries = false
  }
}

async function loadEmbassies() {
  loading.embassies = true
  try {
    embassies.value = await api.listEmbassies({ q: embassyKeyword.value })
  } finally {
    loading.embassies = false
  }
}

async function loadVisas() {
  loading.visas = true
  try {
    if (visaCountryFilter.value) {
      visas.value = await api.listCountryVisasByCountry(visaCountryFilter.value, visaKeyword.value)
    } else {
      visas.value = await api.listCountryVisas({ q: visaKeyword.value })
    }
  } finally {
    loading.visas = false
  }
}

async function loadDetailVisa() {
  if (!detailVisaId.value) {
    detailVisaBase.value = {}
    detailSteps.value = []
    return
  }
  loading.detail = true
  try {
    const detail = await api.getVisaDetail(detailVisaId.value)
    detailVisaBase.value = detail
    detailSteps.value = (detail.steps || []).map((step, idx) => mkLocalStep(step, idx))
  } finally {
    loading.detail = false
  }
}

async function refreshCurrentModule() {
  if (!adminUser.value) return
  await loadHealth()
  if (activeModule.value === 'admins') {
    await ensureCurrentModuleData(true)
    return
  }
  if (activeModule.value === 'users') {
    await ensureCurrentModuleData(true)
    return
  }
  if (activeModule.value === 'community') {
    if (communityAdminRef.value && typeof communityAdminRef.value.loadCurrentTab === 'function') {
      await communityAdminRef.value.loadCurrentTab()
    }
    return
  }
  await ensureCurrentModuleData(true)
}

function openAdminAccountDialog(row) {
  adminAccountDialog.saving = false
  adminAccountDialog.form = row
    ? {
        id: row.id,
        username: row.username || '',
        nickname: row.nickname || '',
        phone: row.phone || '',
        email: row.email || '',
        password: '',
        status: row.status || 'active',
      }
    : defaultAdminAccount()
  adminAccountDialog.visible = true
}

async function saveAdminAccount() {
  if (!adminAccountDialog.form.username) {
    ElMessage.warning('用户名必填')
    return
  }
  if (!adminAccountDialog.form.id && !adminAccountDialog.form.password) {
    ElMessage.warning('新增后台账号时必须设置密码')
    return
  }
  adminAccountDialog.saving = true
  try {
    const payload = {
      username: adminAccountDialog.form.username,
      nickname: adminAccountDialog.form.nickname,
      phone: adminAccountDialog.form.phone,
      email: adminAccountDialog.form.email,
      password: adminAccountDialog.form.password,
      status: adminAccountDialog.form.status,
    }
    if (adminAccountDialog.form.id) {
      await api.updateAdminAccount(adminAccountDialog.form.id, payload)
      ElMessage.success('后台账号已更新')
    } else {
      await api.createAdminAccount(payload)
      ElMessage.success('后台账号已创建')
    }
    adminAccountDialog.visible = false
    await loadAdminAccounts()
  } finally {
    adminAccountDialog.saving = false
  }
}

async function removeAdminAccount(row) {
  await ElMessageBox.confirm(`确认删除后台账号 ${row.username} 吗？`, '提示', { type: 'warning' })
  await api.deleteAdminAccount(row.id)
  ElMessage.success('后台账号已删除')
  await loadAdminAccounts()
}

async function restoreAdminSession() {
  const token = getAdminToken()
  if (!token) {
    clearAdminSession()
    adminUser.value = null
    sessionReady.value = true
    return
  }
  authLoading.value = true
  try {
    adminUser.value = await api.getAdminMe()
  } catch (error) {
    clearAdminSession()
    adminUser.value = null
  } finally {
    authLoading.value = false
    sessionReady.value = true
  }
}

async function handleAdminLoginSuccess(user) {
  adminUser.value = user || getAdminUser()
  sessionReady.value = true
  await loadHealth()
  await ensureCurrentModuleData(true)
}

function handleAdminLogout() {
  clearAdminSession()
  adminUser.value = null
  sessionReady.value = true
  healthText.value = '未登录'
}

function openUserDialog(row) {
  revokeUserAvatarObjectUrl()
  userDialog.avatarFile = null
  userDialog.saving = false
  userDialog.form = row
    ? {
        ...defaultUser(),
        id: row.id,
        username: row.username,
        nickname: row.nickname,
        email: row.email,
        phone: row.phone,
        bio: row.bio,
        gender: row.gender,
        location: row.location,
        avatar: row.avatar,
        password: '',
        role: row.role,
        status: row.status,
      }
    : defaultUser()
  userDialog.locationParts = parseProvinceCity(userDialog.form.location)
  if (userDialog.locationParts.length === 2) {
    userDialog.form.location = formatProvinceCity(userDialog.locationParts)
  }
  userDialog.visible = true
}

async function saveUser() {
  if (userDialog.saving) {
    return
  }
  if (!userDialog.form.nickname || !userDialog.form.phone) {
    ElMessage.warning('昵称和手机号必填')
    return
  }
  const isEdit = Boolean(userDialog.form.id)
  const location = userDialog.locationParts.length === 2 ? formatProvinceCity(userDialog.locationParts) : String(userDialog.form.location || '')
  const payload = {
    username: userDialog.form.username || userDialog.form.phone,
    name: userDialog.form.nickname,
    nickname: userDialog.form.nickname,
    email: userDialog.form.email,
    phone: userDialog.form.phone,
    bio: userDialog.form.bio,
    gender: userDialog.form.gender,
    location,
    avatar: userDialog.form.avatar,
    password: userDialog.form.password,
    role: userDialog.form.role,
    status: userDialog.form.status,
  }

  userDialog.saving = true
  let savedUser = null
  try {
    if (isEdit) {
      savedUser = await api.updateUser(userDialog.form.id, payload)
    } else {
      savedUser = await api.createUser(payload)
      userDialog.form.id = savedUser.id
    }

    if (userDialog.avatarFile) {
      const uploadResult = await api.uploadUserAvatar(savedUser.id, userDialog.avatarFile)
      savedUser = uploadResult.profile || savedUser
      userDialog.form.avatar = savedUser.avatar || ''
    }
  } catch (error) {
    if (savedUser?.id) {
      ElMessage.warning(`资料已保存，但头像上传失败：${error.message || '请重试'}`)
      await loadUsers()
    } else {
      ElMessage.error(error.message || '保存失败')
    }
    userDialog.saving = false
    return
  }

  userDialog.saving = false
  userDialog.visible = false
  ElMessage.success(isEdit ? '用户已更新' : '用户已创建')
  await loadUsers()
}

async function removeUser(row) {
  await ElMessageBox.confirm(`确认删除用户 ${row.name || row.username} 吗？`, '提示', { type: 'warning' })
  await api.deleteUser(row.id)
  ElMessage.success('用户已删除')
  await loadUsers()
}

function openMembershipDialog(row) {
  const membership = row.membership || {}
  const planKey = membership.planKey || 'month'
  membershipDialog.userId = row.id
  membershipDialog.userNickname = row.nickname || row.name || ''
  membershipDialog.form = {
    planKey,
    planName: membership.planName || getMembershipDefaultName(planKey),
    startedAt: membership.startedAt || '',
    expiresAt: membership.expiresAt || '',
    status: membership.status || 'active',
  }
  membershipDialog.visible = true
}

async function saveMembership() {
  if (!membershipDialog.form.planName) {
    membershipDialog.form.planName = getMembershipDefaultName(membershipDialog.form.planKey)
  }
  await api.updateUserMembership(membershipDialog.userId, membershipDialog.form)
  ElMessage.success('会员信息已更新')
  membershipDialog.visible = false
  await loadUsers()
}

function openCountryDialog(row) {
  countryDialog.form = row
    ? {
        ...row,
        tagsCsv: (row.tags || []).join(','),
        keywordsCsv: (row.keywords || []).join(','),
      }
    : defaultCountry()
  countryDialog.visible = true
}

function openFreeCountryDialog(row) {
  freeCountryDialog.form = row
    ? {
        ...row,
        keywordsCsv: (row.keywords || []).join(','),
      }
    : defaultFreeCountry()
  freeCountryDialog.visible = true
}

function openEmbassyDialog(row) {
  embassyDialog.form = row
    ? {
        ...row,
        servicesCsv: (row.services || []).join(','),
        keywordsCsv: (row.keywords || []).join(','),
      }
    : defaultEmbassy()
  embassyDialog.visible = true
}

async function saveCountry() {
  const payload = {
    ...countryDialog.form,
    tags: splitCsv(countryDialog.form.tagsCsv),
    keywords: splitCsv(countryDialog.form.keywordsCsv),
  }
  if (countryDialog.form.id) {
    await api.updateCountry(countryDialog.form.id, payload)
    ElMessage.success('国家已更新')
  } else {
    await api.createCountry(payload)
    ElMessage.success('国家已创建')
  }
  countryDialog.visible = false
  await loadCountries()
  await loadFreeCountries()
  await loadVisas()
}

async function saveFreeCountry() {
  if (!freeCountryDialog.form.name || !freeCountryDialog.form.code) {
    ElMessage.warning('国家名称和代码必填')
    return
  }
  const payload = {
    ...freeCountryDialog.form,
    keywords: splitCsv(freeCountryDialog.form.keywordsCsv),
    supportedVisaId: Number(freeCountryDialog.form.supportedVisaId) || 0,
  }
  if (freeCountryDialog.form.id) {
    await api.updateFreeCountry(freeCountryDialog.form.id, payload)
    ElMessage.success('免签国家已更新')
  } else {
    await api.createFreeCountry(payload)
    ElMessage.success('免签国家已创建')
  }
  freeCountryDialog.visible = false
  await loadFreeCountries()
}

async function saveEmbassy() {
  if (!embassyDialog.form.country || !embassyDialog.form.name) {
    ElMessage.warning('国家和使领馆名称必填')
    return
  }
  const payload = {
    ...embassyDialog.form,
    services: splitCsv(embassyDialog.form.servicesCsv),
    keywords: splitCsv(embassyDialog.form.keywordsCsv),
  }
  if (embassyDialog.form.id) {
    await api.updateEmbassy(embassyDialog.form.id, payload)
    ElMessage.success('使领馆已更新')
  } else {
    await api.createEmbassy(payload)
    ElMessage.success('使领馆已创建')
  }
  embassyDialog.visible = false
  await loadEmbassies()
}

async function removeCountry(row) {
  await ElMessageBox.confirm(`确认删除国家 ${row.name} 吗？该国家下的签证也会一起删除。`, '提示', { type: 'warning' })
  await api.deleteCountry(row.id)
  ElMessage.success('国家已删除')
  await loadCountries()
  await loadFreeCountries()
  await loadVisas()
}

async function removeFreeCountry(row) {
  await ElMessageBox.confirm(`确认删除免签国家 ${row.name} 吗？`, '提示', { type: 'warning' })
  await api.deleteFreeCountry(row.id)
  ElMessage.success('免签国家已删除')
  await loadFreeCountries()
}

async function removeEmbassy(row) {
  await ElMessageBox.confirm(`确认删除使领馆 ${row.name} 吗？`, '提示', { type: 'warning' })
  await api.deleteEmbassy(row.id)
  ElMessage.success('使领馆已删除')
  await loadEmbassies()
}

function focusCountryVisas(row) {
  visaTab.value = 'visas'
  visaCountryFilter.value = row.id
  loadVisas()
}

function openVisaDialog(row) {
  if (!row && !visaCountryFilter.value && !countries.value.length) {
    ElMessage.warning('请先创建国家')
    return
  }
  visaDialog.form = row
    ? { ...row }
    : {
        ...defaultVisa(),
        countryId: visaCountryFilter.value || countries.value[0]?.id,
      }
  visaDialog.visible = true
}

async function saveVisa() {
  if (!visaDialog.form.countryId || !visaDialog.form.name) {
    ElMessage.warning('国家和签证名称必填')
    return
  }
  if (visaDialog.form.id) {
    await api.updateCountryVisa(visaDialog.form.id, visaDialog.form)
    ElMessage.success('签证已更新')
  } else {
    await api.createCountryVisaForCountry(visaDialog.form.countryId, visaDialog.form)
    ElMessage.success('签证已创建')
  }
  visaDialog.visible = false
  await loadVisas()
}

async function removeVisa(row) {
  await ElMessageBox.confirm(`确认删除签证 ${row.name} 吗？`, '提示', { type: 'warning' })
  await api.deleteCountryVisa(row.id)
  ElMessage.success('签证已删除')
  await loadVisas()
  if (detailVisaId.value === row.id) {
    detailVisaId.value = undefined
    detailVisaBase.value = {}
    detailSteps.value = []
  }
}

function selectVisaRow(row) {
  detailVisaId.value = row.id
}

function openVisaDetailEditor(row) {
  visaTab.value = 'detail'
  detailVisaId.value = row.id
  loadDetailVisa()
}

function addStep() {
  detailSteps.value.push(
    mkLocalStep(
      {
        stepKey: `step-${detailSteps.value.length + 1}`,
        title: '',
        sortOrder: detailSteps.value.length + 1,
        strategies: [],
        materials: [],
        guides: [],
        tasks: [],
      },
      detailSteps.value.length
    )
  )
}

function removeStep(index) {
  detailSteps.value.splice(index, 1)
}

function addGuide(step) {
  step.guides.push({ title: '', desc: '', image: '', cta: '', url: '' })
}

function removeGuide(step, index) {
  step.guides.splice(index, 1)
}

async function handleGuideImageChange(uploadFile, guide) {
  const file = uploadFile?.raw
  if (!file) {
    return
  }
  if (!String(file.type || '').startsWith('image/')) {
    ElMessage.warning('请选择图片文件')
    return
  }

  guide._uploading = true
  try {
    const result = await api.uploadImage(file, 'guides')
    guide.image = result.url || ''
    ElMessage.success('图片已上传')
  } catch (error) {
    ElMessage.error(error.message || '图片上传失败')
  } finally {
    guide._uploading = false
  }
}

function clearGuideImage(guide) {
  guide.image = ''
}

function addStepListItem(step, key) {
  if (!Array.isArray(step[key])) {
    step[key] = []
  }
  step[key].push('')
}

function removeStepListItem(step, key, index) {
  if (!Array.isArray(step[key])) {
    return
  }
  step[key].splice(index, 1)
}

function addTask(step) {
  step.tasks.push({
    taskKey: `task-${step.tasks.length + 1}`,
    title: '',
    icon: 'task_alt',
    sortOrder: step.tasks.length + 1,
  })
}

function removeTask(step, index) {
  step.tasks.splice(index, 1)
}

function openTaskIconPicker(step, taskIndex) {
  iconPicker.step = step
  iconPicker.taskIndex = taskIndex
  iconPicker.currentIcon = step?.tasks?.[taskIndex]?.icon || 'task_alt'
  iconPicker.keyword = ''
  iconPicker.visible = true
}

function resetTaskIconPicker() {
  iconPicker.keyword = ''
  iconPicker.currentIcon = ''
  iconPicker.step = null
  iconPicker.taskIndex = -1
}

function selectTaskIcon(iconValue) {
  if (iconPicker.step && iconPicker.taskIndex >= 0 && iconPicker.step.tasks?.[iconPicker.taskIndex]) {
    iconPicker.step.tasks[iconPicker.taskIndex].icon = iconValue
  }
  iconPicker.currentIcon = iconValue
  iconPicker.visible = false
}

async function saveVisaDetail() {
  if (!detailVisaId.value) return
  const payloadSteps = detailSteps.value.map((step, idx) => ({
    stepKey: step.stepKey || `step-${idx + 1}`,
    title: step.title || `Step ${idx + 1}`,
    sortOrder: Number(step.sortOrder) || idx + 1,
    strategies: getStepStrategies(step),
    guides: step.guides.map((guide) => ({
      title: guide.title || '',
      desc: guide.desc || '',
      image: guide.image || '',
      cta: guide.cta || '',
      url: guide.url || '',
    })),
    materials: getStepMaterials(step),
    tasks: step.tasks.map((task, taskIndex) => ({
      taskKey: task.taskKey || `task-${taskIndex + 1}`,
      title: task.title || '',
      icon: task.icon || 'task_alt',
      sortOrder: Number(task.sortOrder) || taskIndex + 1,
    })),
  }))

  await api.updateVisaDetail(detailVisaId.value, payloadSteps)
  ElMessage.success('签证详情模板已更新')
  await loadDetailVisa()
}

watch(activeModule, async (value) => {
  if (!adminUser.value) return
  if (value === 'community') {
    return
  }
  await ensureCurrentModuleData(false)
})

watch(visaTab, async () => {
  if (!adminUser.value || activeModule.value !== 'visa') return
  await ensureCurrentVisaTabData(false)
})

onMounted(async () => {
  await loadHealth()
  await restoreAdminSession()
  if (adminUser.value) {
    await ensureCurrentModuleData(false)
  }
})
</script>

<style scoped>
.admin-shell {
  height: 100vh;
  overflow: hidden;
}

.sidebar {
  background: #0f172a;
  color: #fff;
  padding: 22px 14px;
  height: 100vh;
  overflow: auto;
  box-sizing: border-box;
}

.brand {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 28px;
  padding: 0 8px;
}

.brand-mark {
  width: 42px;
  height: 42px;
  border-radius: 10px;
  background: #1d4ed8;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 900;
}

.brand-title {
  font-size: 17px;
  font-weight: 800;
}

.brand-sub {
  margin-top: 4px;
  color: #94a3b8;
  font-size: 12px;
}

.nav-menu {
  border-right: none;
  background: transparent;
}

.nav-menu :deep(.el-menu-item) {
  color: #cbd5e1;
  border-radius: 8px;
  margin-bottom: 6px;
}

.nav-menu :deep(.el-menu-item.is-active) {
  color: #fff;
  background: #1d4ed8;
}

.topbar {
  height: 72px;
  padding: 0 28px;
  background: #fff;
  border-bottom: 1px solid #e5e7eb;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.topbar-actions {
  display: flex;
  align-items: center;
  gap: 12px;
}

.admin-chip {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  border-radius: 999px;
  background: #eff6ff;
  color: #1d4ed8;
}

.admin-chip__label {
  font-size: 12px;
  font-weight: 700;
}

.admin-chip__name {
  font-size: 13px;
  font-weight: 800;
}

.page-title {
  font-size: 20px;
  font-weight: 900;
}

.page-sub,
.panel-sub,
.membership-meta,
.detail-desc {
  margin-top: 4px;
  color: #6b7280;
  font-size: 13px;
}

.main {
  padding: 24px;
  min-height: 0;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  box-sizing: border-box;
}

.panel {
  flex: 1;
  min-height: 0;
  height: 100%;
  overflow: auto;
  padding: 20px;
  border-radius: 10px;
  background: #fff;
  border: 1px solid #e5e7eb;
  box-sizing: border-box;
}

.main-panel {
  flex: 1;
  min-height: 0;
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
  width: 220px;
}

.toolbar--wide {
  flex-wrap: wrap;
}

.membership-cell {
  line-height: 1.4;
}

.avatar-editor {
  display: flex;
  align-items: center;
  gap: 16px;
}

.avatar-editor__actions,
.location-editor {
  display: grid;
  gap: 8px;
  min-width: 0;
}

.avatar-editor__buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.region-cascader {
  width: 100%;
}

.form-hint {
  font-size: 12px;
  color: #6b7280;
  line-height: 1.5;
}

.user-avatar-thumb {
  width: 40px;
  height: 40px;
  border-radius: 999px;
  object-fit: cover;
  display: block;
  margin: 0 auto;
}

.membership-meta {
  margin-top: 2px;
  font-size: 12px;
}

.visa-tabs :deep(.el-tabs__header) {
  margin-bottom: 16px;
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

.detail-editor {
  display: grid;
  grid-template-columns: minmax(0, 1.2fr) minmax(360px, 0.8fr);
  grid-template-areas:
    'header header'
    'editor preview';
  align-items: start;
  gap: 16px;
  min-height: 0;
  overflow: visible;
  border: 1px solid #e5e7eb;
  border-radius: 10px;
  padding: 16px;
  box-sizing: border-box;
}

.detail-header-row {
  grid-area: header;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.detail-overview {
  grid-area: preview;
  display: grid;
  gap: 14px;
  align-self: start;
  position: sticky;
  top: 0;
  max-height: calc(100vh - 220px);
  overflow: auto;
  padding-right: 4px;
}

.overview-card,
.preview-step-card {
  border-radius: 12px;
  border: 1px solid #e5e7eb;
}

.overview-title {
  font-weight: 700;
}

.overview-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 10px;
}

.overview-item {
  padding: 10px 12px;
  border-radius: 10px;
  background: #f8fafc;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.overview-label {
  font-size: 12px;
  color: #64748b;
}

.overview-value {
  font-size: 14px;
  font-weight: 600;
  color: #0f172a;
}

.overview-intro {
  margin-top: 12px;
  padding: 12px;
  border-radius: 10px;
  background: #f8fafc;
  line-height: 1.6;
  color: #334155;
}

.preview-step-list {
  display: grid;
  gap: 12px;
}

.preview-step-head {
  display: flex;
  align-items: center;
  gap: 12px;
}

.preview-step-index {
  min-width: 64px;
  height: 28px;
  padding: 0 10px;
  border-radius: 999px;
  background: #dbeafe;
  color: #1d4ed8;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
}

.preview-step-title {
  font-weight: 700;
  color: #0f172a;
}

.preview-step-key {
  margin-top: 4px;
  font-size: 12px;
  color: #64748b;
}

.preview-block {
  margin-top: 14px;
}

.preview-block-title {
  margin-bottom: 8px;
  font-size: 12px;
  font-weight: 700;
  color: #334155;
  text-transform: uppercase;
  letter-spacing: 0.04em;
}

.preview-bullet-list {
  display: grid;
  gap: 8px;
}

.preview-bullet-list--muted .preview-bullet-item {
  background: #f8fafc;
}

.preview-bullet-list--muted .preview-bullet-dot {
  background: #64748b;
}

.preview-bullet-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 10px;
  background: #eff6ff;
  color: #334155;
  line-height: 1.5;
  font-size: 13px;
}

.preview-bullet-dot {
  width: 8px;
  height: 8px;
  border-radius: 999px;
  background: #1d4ed8;
  margin-top: 6px;
  flex-shrink: 0;
}

.preview-guide + .preview-guide {
  margin-top: 10px;
  padding-top: 10px;
  border-top: 1px dashed #cbd5e1;
}

.preview-guide-title {
  font-weight: 600;
  color: #0f172a;
}

.preview-guide-desc {
  margin-top: 4px;
  line-height: 1.5;
  color: #475569;
}

.preview-guide-link {
  display: inline-flex;
  margin-top: 8px;
  color: #2563eb;
  text-decoration: none;
  font-weight: 600;
}

.preview-guide-image {
  margin-top: 8px;
  width: 100%;
  max-width: 220px;
  height: 126px;
  object-fit: cover;
  border-radius: 10px;
  border: 1px solid #dbe3ef;
}

.preview-task-list {
  display: grid;
  gap: 8px;
}

.preview-task {
  padding: 10px 12px;
  border-radius: 10px;
  background: #f8fafc;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.preview-task-title {
  color: #0f172a;
}

.preview-task-key {
  font-size: 12px;
  color: #64748b;
}

.preview-task-status {
  padding: 4px 8px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
}

.preview-task-status--done {
  background: #dcfce7;
  color: #15803d;
}

.preview-task-status--review {
  background: #dbeafe;
  color: #1d4ed8;
}

.preview-task-status--missing {
  background: #fee2e2;
  color: #dc2626;
}

.preview-task-status--todo {
  background: #f1f5f9;
  color: #475569;
}

.detail-name {
  font-size: 16px;
  font-weight: 700;
}

.step-list-wrap {
  grid-area: editor;
  align-self: start;
  height: calc(100vh - 220px);
  min-height: 0;
  overflow-y: auto;
  overflow-x: hidden;
  padding-right: 6px;
}

.step-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.detail-editor > .el-empty {
  grid-area: editor;
}

.step-card :deep(.el-card__header) {
  padding: 12px 16px;
}

.step-card :deep(.el-card__body) {
  padding: 16px;
}

.step-card-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.step-head-main {
  min-width: 0;
  flex: 1;
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}

.step-toggle-icon {
  font-size: 20px;
  color: #64748b;
  flex-shrink: 0;
}

.step-title-wrap {
  min-width: 0;
}

.step-title {
  font-weight: 700;
  color: #0f172a;
}

.step-subtitle {
  margin-top: 4px;
  font-size: 12px;
  color: #64748b;
}

.step-body {
  display: grid;
  gap: 16px;
  padding-top: 4px;
}

.step-meta-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
}

.step-content-grid {
  display: grid;
  grid-template-columns: 1fr;
  gap: 16px;
  align-items: start;
}

.editor-panel {
  display: grid;
  gap: 14px;
  padding: 14px;
  border-radius: 14px;
  border: 1px solid #e2e8f0;
  background: #fbfcfe;
}

.editor-panel--compact {
  padding: 12px;
}

.editor-panel__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.editor-panel__title {
  font-size: 14px;
  font-weight: 700;
  color: #0f172a;
}

.editor-panel__sub {
  margin-top: 4px;
  font-size: 12px;
  line-height: 1.5;
  color: #64748b;
}

.editor-field {
  display: grid;
  gap: 8px;
}

.editor-field__label {
  font-size: 12px;
  font-weight: 700;
  color: #475569;
}

.editor-table :deep(.el-table__header th) {
  background: #f8fafc;
}

.guide-image-uploader {
  display: grid;
  gap: 10px;
}

.guide-image-uploader__actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.guide-image-thumb {
  width: 132px;
  height: 96px;
  border-radius: 10px;
  border: 1px solid #dbe3ef;
  overflow: hidden;
  background: #fff;
}

.guide-image-empty {
  width: 132px;
  height: 96px;
  border-radius: 10px;
  border: 1px dashed #cbd5e1;
  background: #f8fafc;
  color: #64748b;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
}

.list-editor {
  display: grid;
  gap: 10px;
  width: 100%;
}

.list-editor__rows {
  display: grid;
  gap: 8px;
}

.list-editor__row {
  display: grid;
  grid-template-columns: 28px minmax(0, 1fr) auto;
  gap: 10px;
  align-items: center;
  padding: 10px 12px;
  border-radius: 10px;
  background: #f8fafc;
  border: 1px solid #e2e8f0;
}

.list-editor__index {
  width: 28px;
  height: 28px;
  border-radius: 999px;
  background: #dbeafe;
  color: #1d4ed8;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  font-weight: 700;
}

.list-editor__empty {
  padding: 12px;
  border-radius: 10px;
  background: #f8fafc;
  border: 1px dashed #cbd5e1;
  color: #64748b;
  font-size: 13px;
}

.task-icon-picker {
  width: 100%;
  padding: 8px 10px;
  border-radius: 10px;
  border: 1px solid #dbe3ef;
  background: #f8fafc;
  display: flex;
  align-items: center;
  gap: 10px;
  color: #334155;
  cursor: pointer;
}

.task-icon-picker__icon {
  font-size: 20px;
  color: #1d4ed8;
}

.task-icon-picker__label {
  font-size: 13px;
  font-weight: 600;
}

.icon-picker {
  display: grid;
  gap: 14px;
}

.icon-picker__grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 10px;
  max-height: 420px;
  overflow: auto;
  padding-right: 4px;
}

.icon-picker__item {
  padding: 12px 10px;
  border-radius: 12px;
  border: 1px solid #dbe3ef;
  background: #fff;
  display: grid;
  justify-items: center;
  gap: 6px;
  color: #334155;
  cursor: pointer;
  text-align: center;
}

.icon-picker__item--active {
  border-color: #1d4ed8;
  background: #eff6ff;
  color: #1d4ed8;
}

.icon-picker__symbol {
  font-size: 22px;
}

.icon-picker__name {
  font-size: 13px;
  font-weight: 700;
}

.icon-picker__value {
  font-size: 11px;
  color: #64748b;
}

.icon-picker__empty {
  padding: 12px;
  border-radius: 10px;
  background: #f8fafc;
  color: #64748b;
  text-align: center;
}

@media (max-width: 1200px) {
  .detail-editor {
    grid-template-columns: 1fr;
    grid-template-areas:
      'header'
      'editor'
      'preview';
    height: auto;
    max-height: none;
    overflow: visible;
  }

  .detail-overview {
    position: static;
    max-height: none;
    height: auto;
    overflow: visible;
    padding-right: 0;
  }

  .step-list-wrap {
    height: auto;
    overflow: visible;
    padding-right: 0;
  }

  .overview-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .step-meta-grid,
  .step-content-grid {
    grid-template-columns: 1fr;
  }

  .icon-picker__grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
</style>
