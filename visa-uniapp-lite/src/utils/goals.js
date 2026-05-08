const GOALS_KEY = 'visago_lite_goals'

const GOAL_STEP_TEMPLATES = [
  {
    stepKey: 'apply',
    title: '申请',
    tasks: [
      { id: 'confirm-visa', title: '确认签证方案', icon: 'travel_explore' },
    ],
  },
  {
    stepKey: 'docs',
    title: '材料',
    tasks: [
      { id: 'prepare-materials', title: '整理基础材料', icon: 'description' },
      { id: 'note-plan', title: '补充办理备注', icon: 'edit_note' },
    ],
  },
  {
    stepKey: 'book',
    title: '预约',
    tasks: [
      { id: 'schedule-time', title: '安排办理时间', icon: 'calendar_month' },
    ],
  },
  {
    stepKey: 'result',
    title: '结果',
    tasks: [
      { id: 'mark-complete', title: '记录本次结果', icon: 'task_alt' },
    ],
  },
]

const STATUS_TEXT = {
  done: '已完成',
  review: '进行中',
  todo: '待处理',
  missing: '待补充',
}

function safeArray(value) {
  return Array.isArray(value) ? value : []
}

function cloneTaskStateMap(value) {
  return value && typeof value === 'object' && !Array.isArray(value) ? { ...value } : {}
}

function normalizeGuide(guide) {
  const current = guide && typeof guide === 'object' ? guide : {}
  return {
    title: String(current.title || '').trim(),
    desc: String(current.desc || current.description || '').trim(),
    image: String(current.image || '').trim(),
    cta: String(current.cta || '').trim(),
    url: String(current.url || '').trim(),
  }
}

function normalizeVisaSteps(steps) {
  return safeArray(steps)
    .map((step) => {
      const current = step && typeof step === 'object' ? step : {}
      return {
        stepKey: String(current.stepKey || '').trim(),
        title: String(current.title || '').trim(),
        strategies: safeArray(current.strategies).map((item) => String(item || '').trim()).filter(Boolean),
        guides: safeArray(current.guides).map(normalizeGuide).filter((item) => item.title || item.desc || item.image || item.cta || item.url),
        materials: safeArray(current.materials).map((item) => String(item || '').trim()).filter(Boolean),
      }
    })
    .filter((step) => step.stepKey || step.title || step.strategies.length || step.guides.length || step.materials.length)
}

function allTaskIds() {
  return GOAL_STEP_TEMPLATES.flatMap((step) => step.tasks.map((task) => task.id))
}

function buildDetailUrl(goal) {
  return `/pages/visa-detail/index?visaId=${goal.visaId}&countryName=${encodeURIComponent(goal.countryName || '')}`
}

function seedTaskStateMap(goal) {
  const current = cloneTaskStateMap(goal.taskStateMap)
  const next = {
    'confirm-visa': true,
    ...current,
  }

  if ((goal.note || '').trim()) {
    next['note-plan'] = true
  }

  if (goal.status === 'preparing') {
    next['prepare-materials'] = true
  }

  if (goal.status === 'done') {
    allTaskIds().forEach((taskId) => {
      next[taskId] = true
    })
  }

  return next
}

function sortGoals(goals) {
  return safeArray(goals).sort((a, b) => String(b.updatedAt || '').localeCompare(String(a.updatedAt || '')))
}

function normalizeGoal(goal) {
  const now = new Date().toISOString()
  return {
    status: 'planned',
    note: '',
    resultStatus: 'pending',
    resultNote: '',
    resultAt: '',
    visaSteps: [],
    createdAt: now,
    updatedAt: now,
    ...goal,
    visaSteps: normalizeVisaSteps((goal && goal.visaSteps) || []),
    taskStateMap: seedTaskStateMap(goal || {}),
  }
}

function findVisaStep(goal, stepKey, stepIndex) {
  const steps = normalizeVisaSteps(goal && goal.visaSteps)
  return steps.find((item) => item.stepKey === stepKey) || steps[stepIndex] || null
}

function buildStepContent(goal, stepKey, stepIndex) {
  const visaName = goal.visaName || '当前签证'
  const countryName = goal.countryName || '目标国家'
  const detailUrl = buildDetailUrl(goal)
  const visaStep = findVisaStep(goal, stepKey, stepIndex)
  const stepStrategies = visaStep && visaStep.strategies.length ? visaStep.strategies : null
  const stepMaterials = visaStep && visaStep.materials.length ? visaStep.materials : null
  const stepGuides = visaStep && visaStep.guides.length ? visaStep.guides : null

  if (stepKey === 'apply') {
    return {
      summary:
        (stepGuides && stepGuides[0] && stepGuides[0].desc) ||
        `先确认 ${countryName} 的 ${visaName} 是否适合当前出行安排。`,
      strategies: stepStrategies || [
        `确认 ${visaName} 的适用场景、费用和有效期。`,
        '结合出行时间判断是否需要尽快开始准备。',
      ],
      materials: stepMaterials || [
        goal.visaType ? `签证类型：${goal.visaType}` : '签证类型：待补充',
        goal.fee ? `费用参考：${goal.fee}` : '费用参考：待补充',
      ],
      guides: stepGuides || [
        {
          title: '确认签证方案',
          desc: `先把 ${visaName} 的基础信息看清楚，再决定后续办理节奏。`,
          cta: '查看签证详情',
          url: detailUrl,
        },
      ],
      note: (stepGuides && stepGuides[0] && stepGuides[0].cta) || '确认无误后，再开始准备材料会更顺手。',
    }
  }

  if (stepKey === 'docs') {
    return {
      summary:
        (stepGuides && stepGuides[0] && stepGuides[0].desc) ||
        '把基础材料和个人补充说明先整理出来。',
      strategies: stepStrategies || [
        '优先整理最稳定、最容易准备的基础文件。',
        (goal.note || '').trim() ? '你已经补充了个人备注，可在此基础上继续完善。' : '如果有特殊情况，建议在目标里补充备注方便回看。',
      ],
      materials: stepMaterials || [
        '护照原件与复印件',
        goal.validity ? `有效期要求参考：${goal.validity}` : '有效期要求：待补充',
        (goal.note || '').trim() ? `当前备注：${goal.note.trim()}` : '当前备注：未填写',
      ],
      guides: stepGuides || [
        {
          title: '整理申请材料',
          desc: '按照签证详情页的信息准备材料，先完成最确定的部分。',
          cta: '查看签证详情',
          url: detailUrl,
        },
      ],
      note: (stepGuides && stepGuides[0] && stepGuides[0].cta) || '材料先整理到位，后面预约和递交会更轻松。',
    }
  }

  if (stepKey === 'book') {
    return {
      summary:
        (stepGuides && stepGuides[0] && stepGuides[0].desc) ||
        '结合办理时长和出发日期安排递交节奏。',
      strategies: stepStrategies || [
        goal.processingTime ? `当前办理时长参考为 ${goal.processingTime}。` : '办理时长还没有补充，先按宽松时间预留。',
        '如果时间紧张，建议尽早安排递交或预约。',
      ],
      materials: stepMaterials || [
        goal.entries ? `入境次数：${goal.entries}` : '入境次数：待补充',
        goal.processingTime ? `办理时长：${goal.processingTime}` : '办理时长：待补充',
      ],
      guides: stepGuides || [
        {
          title: '安排办理时间',
          desc: '根据出行日期倒推准备时间，避免临近出发再集中处理。',
          cta: '查看签证详情',
          url: detailUrl,
        },
      ],
      note: (stepGuides && stepGuides[0] && stepGuides[0].cta) || '时间安排越清楚，后续的进度也越容易掌控。',
    }
  }

  return {
    summary:
      (stepGuides && stepGuides[0] && stepGuides[0].desc) ||
      '当所有事项处理完后，记得记录本次结果。',
    strategies: stepStrategies || [
      '全部任务完成后，可以把本次目标标记为完成。',
      '如果结果有变化，也可以随时重新调整。',
    ],
    materials: stepMaterials || [
      goal.fee ? `费用参考：${goal.fee}` : '费用参考：待补充',
      goal.processingTime ? `办理时长：${goal.processingTime}` : '办理时长：待补充',
    ],
    guides: stepGuides || [
      {
        title: '记录本次结果',
        desc: '完成全部准备后，在这里记录本次目标的最终状态。',
        cta: '查看签证详情',
        url: detailUrl,
      },
    ],
    note: (stepGuides && stepGuides[0] && stepGuides[0].cta) || '记录结果后，也能更方便回看之前的准备情况。',
  }
}

export function loadGoals() {
  try {
    const raw = uni.getStorageSync(GOALS_KEY)
    if (!raw) return []
    const parsed = Array.isArray(raw) ? raw : JSON.parse(raw)
    return sortGoals(parsed.map((item) => normalizeGoal(item)))
  } catch (error) {
    return []
  }
}

export function saveGoals(goals) {
  const next = sortGoals(safeArray(goals).map((item) => normalizeGoal(item)))
  uni.setStorageSync(GOALS_KEY, next)
  return next
}

export function getGoalByVisaId(visaId) {
  return loadGoals().find((item) => String(item.visaId) === String(visaId)) || null
}

export function upsertGoal(goal) {
  const now = new Date().toISOString()
  const items = loadGoals()
  const nextItem = normalizeGoal({
    createdAt: now,
    updatedAt: now,
    ...goal,
  })
  const index = items.findIndex((item) => String(item.visaId) === String(nextItem.visaId))
  if (index >= 0) {
    items[index] = normalizeGoal({
      ...items[index],
      ...nextItem,
      createdAt: items[index].createdAt || now,
      updatedAt: now,
    })
  } else {
    items.unshift(nextItem)
  }
  return saveGoals(items)
}

export function removeGoal(visaId) {
  return saveGoals(loadGoals().filter((item) => String(item.visaId) !== String(visaId)))
}

export function updateGoalNote(visaId, note) {
  const items = loadGoals()
  const index = items.findIndex((item) => String(item.visaId) === String(visaId))
  if (index < 0) return items
  const current = items[index]
  items[index] = normalizeGoal({
    ...current,
    note,
    taskStateMap: {
      ...current.taskStateMap,
      'note-plan': Boolean(String(note || '').trim()),
    },
    updatedAt: new Date().toISOString(),
  })
  return saveGoals(items)
}

export function toggleGoalTask(visaId, taskId) {
  const items = loadGoals()
  const index = items.findIndex((item) => String(item.visaId) === String(visaId))
  if (index < 0) return items
  const current = items[index]
  const nextValue = !current.taskStateMap[taskId]
  const taskStateMap = {
    ...current.taskStateMap,
    [taskId]: nextValue,
  }

  if (taskId === 'mark-complete' && nextValue) {
    allTaskIds().forEach((id) => {
      taskStateMap[id] = true
    })
  }

  items[index] = normalizeGoal({
    ...current,
    taskStateMap,
    updatedAt: new Date().toISOString(),
  })
  return saveGoals(items)
}

export function updateGoalResult(visaId, resultStatus, resultNote = '') {
  const items = loadGoals()
  const index = items.findIndex((item) => String(item.visaId) === String(visaId))
  if (index < 0) return items
  const current = items[index]
  items[index] = normalizeGoal({
    ...current,
    resultStatus,
    resultNote,
    resultAt: resultStatus === 'pending' ? '' : new Date().toISOString(),
    updatedAt: new Date().toISOString(),
  })
  return saveGoals(items)
}

export function syncGoalVisaDetail(visaId, detail) {
  const items = loadGoals()
  const index = items.findIndex((item) => String(item.visaId) === String(visaId))
  if (index < 0 || !detail) return items
  const current = items[index]
  items[index] = normalizeGoal({
    ...current,
    visaName: detail.name || current.visaName,
    visaType: detail.visaType || current.visaType,
    processingTime: detail.processingTime || current.processingTime,
    fee: detail.fee || current.fee,
    validity: detail.validity || current.validity,
    entries: detail.entries || current.entries,
    description: detail.description || detail.longIntro || current.description,
    visaSteps: normalizeVisaSteps(detail.steps),
    updatedAt: new Date().toISOString(),
  })
  return saveGoals(items)
}

export function buildGoalView(rawGoal) {
  const goal = normalizeGoal(rawGoal)
  const baseSteps = GOAL_STEP_TEMPLATES.map((step, stepIndex) => {
    const content = buildStepContent(goal, step.stepKey, stepIndex)
    return {
      stepKey: step.stepKey,
      title: step.title,
      summary: content.summary,
      strategies: content.strategies,
      materials: content.materials,
      guides: content.guides,
      note: content.note,
      tasks: step.tasks.map((task) => ({
        ...task,
        done: Boolean(goal.taskStateMap[task.id]),
      })),
    }
  })

  const totalTasks = baseSteps.reduce((sum, step) => sum + step.tasks.length, 0)
  const doneTasks = baseSteps.reduce((sum, step) => sum + step.tasks.filter((task) => task.done).length, 0)
  const firstPendingIndex = baseSteps.findIndex((step) => step.tasks.some((task) => !task.done))
  const activeStepIndex = firstPendingIndex >= 0 ? firstPendingIndex : Math.max(baseSteps.length - 1, 0)

  const steps = baseSteps.map((step, stepIndex) => {
    const allDone = step.tasks.every((task) => task.done)
    const isActive = stepIndex === activeStepIndex && !allDone
    let reviewAssigned = false
    return {
      ...step,
      status: allDone ? 'done' : isActive ? 'active' : 'todo',
      tasks: step.tasks.map((task) => {
        let status = 'todo'
        if (task.done) {
          status = 'done'
        } else if (isActive && !reviewAssigned) {
          status = 'review'
          reviewAssigned = true
        }
        return {
          ...task,
          status,
          statusText: STATUS_TEXT[status],
        }
      }),
    }
  })

  const progress = totalTasks ? Math.round((doneTasks / totalTasks) * 100) : 0
  const tips = []
  if (goal.processingTime) tips.push(`办理时长：${goal.processingTime}`)
  if (goal.fee) tips.push(`费用参考：${goal.fee}`)
  if (goal.validity) tips.push(`有效期：${goal.validity}`)
  if (goal.entries) tips.push(`入境次数：${goal.entries}`)
  if ((goal.note || '').trim()) tips.push(`备注：${goal.note.trim()}`)

  return {
    ...goal,
    id: goal.visaId,
    visaTitle: goal.visaName || '签证计划',
    countryFlag: goal.countryFlag || (goal.countryName ? String(goal.countryName).slice(0, 1) : '签'),
    progress,
    steps,
    tips,
    activeStepKey: steps[activeStepIndex] ? steps[activeStepIndex].stepKey : 'apply',
  }
}
