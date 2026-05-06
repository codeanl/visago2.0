export const homeQuickMenu = [
  { key: 'visa', label: '查签证', icon: 'description' },
  { key: 'photo', label: '拍照证件', icon: 'photo_camera' },
  { key: 'translate', label: '翻译件', icon: 'translate' },
  { key: 'form', label: '填表辅助', icon: 'edit_note' },
]

export const homeHotDestinations = [
  {
    id: 'jp',
    name: '日本',
    flag: '🇯🇵',
    price: '¥299',
    type: '电子签',
    time: '5-7个工作日',
    hot: true,
    image:
      'https://images.pexels.com/photos/918275/pexels-photo-918275.jpeg?auto=compress&cs=tinysrgb&w=1200',
    keywords: ['日本', '电子签', '旅游', 'Japan', '富士山'],
  },
  {
    id: 'au',
    name: '澳大利亚',
    flag: '🇦🇺',
    price: '¥1,150',
    type: '访客签证',
    time: '15-20个工作日',
    hot: false,
    image:
      'https://images.pexels.com/photos/995765/pexels-photo-995765.jpeg?auto=compress&cs=tinysrgb&w=1200',
    keywords: ['澳大利亚', '访客签证', '探亲', 'Australia', '悉尼歌剧院'],
  },
  {
    id: 'uk',
    name: '英国',
    flag: '🇬🇧',
    price: '¥1,280',
    type: '旅游签证',
    time: '10-15个工作日',
    hot: false,
    image:
      'https://images.pexels.com/photos/460672/pexels-photo-460672.jpeg?auto=compress&cs=tinysrgb&w=1200',
    keywords: ['英国', '旅游', '工作', 'UK', '大本钟'],
  },
]

export const homeToolGroups = {
  quick: [
    { key: 'fx', icon: 'currency_exchange', title: '汇率换算' },
    { key: 'photo-check', icon: 'portrait', title: '证件照标准查询' },
    { key: 'scan', icon: 'document_scanner', title: '文件扫描仪' },
  ],
  common: [
    { key: 'translate', icon: 'g_translate', title: '翻译助手' },
    { key: 'insurance', icon: 'health_and_safety', title: '保险购买' },
    { key: 'embassy', icon: 'apartment', title: '驻华使领馆' },
  ],
}

export const visaFreeCountries = [
  {
    id: 'th',
    name: '泰国',
    city: '曼谷',
    type: '落地签',
    stay: '最长30天',
    note: '持中国护照可申请落地签，建议准备返程机票与酒店订单。',
    x: 74,
    y: 62,
  },
  {
    id: 'sg',
    name: '新加坡',
    city: '新加坡',
    type: '免签过境/电子入境',
    stay: '短停友好',
    note: '中转和短停场景常见，入境前请确认最新政策。',
    x: 77,
    y: 68,
  },
  {
    id: 'ae',
    name: '阿联酋',
    city: '迪拜',
    type: '免签',
    stay: '最长30天',
    note: '商务与旅游热度高，护照有效期需满足要求。',
    x: 62,
    y: 54,
  },
  {
    id: 'my',
    name: '马来西亚',
    city: '吉隆坡',
    type: '免签',
    stay: '最长30天',
    note: '热门自由行目的地，建议提前准备回程证明。',
    x: 76,
    y: 66,
  },
  {
    id: 'id',
    name: '印度尼西亚',
    city: '雅加达',
    type: '落地签',
    stay: '最长30天',
    note: '巴厘岛等目的地受欢迎，入境材料要齐全。',
    x: 80,
    y: 72,
  },
  {
    id: 'qa',
    name: '卡塔尔',
    city: '多哈',
    type: '免签',
    stay: '最长30天',
    note: '常见中转地，建议核对航班衔接时间。',
    x: 60,
    y: 56,
  },
  {
    id: 'ma',
    name: '摩洛哥',
    city: '马拉喀什',
    type: '免签',
    stay: '最长90天',
    note: '北非热门目的地，入境时请准备行程单。',
    x: 43,
    y: 49,
  },
  {
    id: 'mu',
    name: '毛里求斯',
    city: '路易港',
    type: '免签',
    stay: '最长60天',
    note: '海岛旅游热门，建议准备保险和酒店凭证。',
    x: 61,
    y: 78,
  },
]

export const visaRegions = ['全部地区', '亚洲', '欧洲', '美洲', '大洋洲', '非洲']

export const visaCountries = [
  {
    id: 'jp',
    name: '日本',
    region: '亚洲',
    tags: ['旅游', '商务', '留学'],
    note: '东京-大阪双城热门，材料准备相对清晰',
    image:
      'https://upload.wikimedia.org/wikipedia/commons/thumb/f/f8/View_of_Mount_Fuji_from_%C5%8Cwakudani_20211202.jpg/330px-View_of_Mount_Fuji_from_%C5%8Cwakudani_20211202.jpg',
    keywords: ['日本', 'japan', '东京', '富士山'],
  },
  {
    id: 'kr',
    name: '韩国',
    region: '亚洲',
    tags: ['旅游', '探亲', '商务'],
    note: '首尔短期出行热门目的地',
    image:
      'https://upload.wikimedia.org/wikipedia/commons/thumb/6/63/%EA%B4%91%ED%99%94%EB%AC%B8_%EC%9B%94%EB%8C%80.jpg/330px-%EA%B4%91%ED%99%94%EB%AC%B8_%EC%9B%94%EB%8C%80.jpg',
    keywords: ['韩国', 'korea', '首尔', '景福宫'],
  },
  {
    id: 'sg',
    name: '新加坡',
    region: '亚洲',
    tags: ['旅游', '商务', '会展'],
    note: '电子化流程较多，出签节奏快',
    image:
      'https://upload.wikimedia.org/wikipedia/commons/thumb/c/c7/Marina_Bay_Sands_%28I%29.jpg/330px-Marina_Bay_Sands_%28I%29.jpg',
    keywords: ['新加坡', 'singapore', '滨海湾'],
  },
  {
    id: 'th',
    name: '泰国',
    region: '亚洲',
    tags: ['旅游', '商务', '探亲'],
    note: '旅游签常见，材料模板成熟',
    image:
      'https://upload.wikimedia.org/wikipedia/commons/thumb/c/c7/0005574_-_Wat_Phra_Kaew_006.jpg/330px-0005574_-_Wat_Phra_Kaew_006.jpg',
    keywords: ['泰国', 'thailand', '曼谷', '大皇宫'],
  },
  {
    id: 'ae',
    name: '阿联酋',
    region: '亚洲',
    tags: ['旅游', '商务', '中转'],
    note: '迪拜中转与会展场景较多',
    image:
      'https://upload.wikimedia.org/wikipedia/commons/thumb/9/90/Burj_Khalifa_%28worlds_tallest_building%29_and_the_Dubai_skyline_%2825781049892%29.jpg/330px-Burj_Khalifa_%28worlds_tallest_building%29_and_the_Dubai_skyline_%2825781049892%29.jpg',
    keywords: ['阿联酋', 'uae', 'dubai', '迪拜', '哈利法塔'],
  },
  {
    id: 'uk',
    name: '英国',
    region: '欧洲',
    tags: ['旅游', '工作', '高潜人才'],
    note: '需要预约线下采集指纹与面签',
    image:
      'https://upload.wikimedia.org/wikipedia/commons/thumb/4/43/Elizabeth_Tower%2C_June_2022.jpg/330px-Elizabeth_Tower%2C_June_2022.jpg',
    keywords: ['英国', 'uk', '伦敦', '大本钟'],
  },
  {
    id: 'fr',
    name: '法国',
    region: '欧洲',
    tags: ['申根', '旅游', '商务'],
    note: '申根热门，建议提前预约递签',
    image:
      'https://upload.wikimedia.org/wikipedia/commons/thumb/8/85/Tour_Eiffel_Wikimedia_Commons_%28cropped%29.jpg/330px-Tour_Eiffel_Wikimedia_Commons_%28cropped%29.jpg',
    keywords: ['法国', 'france', '巴黎', '埃菲尔铁塔', '申根'],
  },
  {
    id: 'it',
    name: '意大利',
    region: '欧洲',
    tags: ['申根', '旅游', '文化'],
    note: '申根长线旅行热门国家',
    image:
      'https://upload.wikimedia.org/wikipedia/commons/thumb/d/de/Colosseo_2020.jpg/330px-Colosseo_2020.jpg',
    keywords: ['意大利', 'italy', '罗马', '斗兽场', '申根'],
  },
  {
    id: 'de',
    name: '德国',
    region: '欧洲',
    tags: ['申根', '商务', '留学'],
    note: '商务与展会需求较高',
    image:
      'https://upload.wikimedia.org/wikipedia/commons/thumb/a/a6/Brandenburger_Tor_abends.jpg/330px-Brandenburger_Tor_abends.jpg',
    keywords: ['德国', 'germany', '柏林', '勃兰登堡门', '申根'],
  },
  {
    id: 'es',
    name: '西班牙',
    region: '欧洲',
    tags: ['申根', '旅游', '数字游民'],
    note: '西语区热门，旅游材料关注行程',
    image:
      'https://upload.wikimedia.org/wikipedia/commons/thumb/2/26/%CE%A3%CE%B1%CE%B3%CF%81%CE%AC%CE%B4%CE%B1_%CE%A6%CE%B1%CE%BC%CE%AF%CE%BB%CE%B9%CE%B1_2941.jpg/330px-%CE%A3%CE%B1%CE%B3%CF%81%CE%AC%CE%B4%CE%B1_%CE%A6%CE%B1%CE%BC%CE%AF%CE%BB%CE%B9%CE%B1_2941.jpg',
    keywords: ['西班牙', 'spain', '巴塞罗那', '圣家堂', '申根'],
  },
  {
    id: 'nl',
    name: '荷兰',
    region: '欧洲',
    tags: ['申根', '旅游', '商务'],
    note: '多城市联游常见，行程证明要完整',
    image:
      'https://upload.wikimedia.org/wikipedia/commons/thumb/8/80/South_facade_of_the_Rijksmuseum_Amsterdam_%28DSCF0528%29.jpg/330px-South_facade_of_the_Rijksmuseum_Amsterdam_%28DSCF0528%29.jpg',
    keywords: ['荷兰', 'netherlands', '阿姆斯特丹', '申根'],
  },
  {
    id: 'us',
    name: '美国',
    region: '美洲',
    tags: ['旅游', '商务', '探亲'],
    note: 'B类签证需求大，面签准备关键',
    image:
      'https://upload.wikimedia.org/wikipedia/commons/thumb/8/89/Front_view_of_Statue_of_Liberty_%28cropped%29.jpg/330px-Front_view_of_Statue_of_Liberty_%28cropped%29.jpg',
    keywords: ['美国', 'usa', '纽约', '自由女神像', 'b1', 'b2'],
  },
  {
    id: 'ca',
    name: '加拿大',
    region: '美洲',
    tags: ['旅游', '探亲', '留学'],
    note: '探亲与留学签证办理量高',
    image:
      'https://upload.wikimedia.org/wikipedia/commons/thumb/a/ab/3Falls_Niagara.jpg/330px-3Falls_Niagara.jpg',
    keywords: ['加拿大', 'canada', '尼亚加拉瀑布', '探亲', '留学'],
  },
  {
    id: 'br',
    name: '巴西',
    region: '美洲',
    tags: ['旅游', '商务', '展会'],
    note: '南美线路热门中转国家',
    image:
      'https://upload.wikimedia.org/wikipedia/commons/thumb/4/4f/Christ_the_Redeemer_-_Cristo_Redentor.jpg/330px-Christ_the_Redeemer_-_Cristo_Redentor.jpg',
    keywords: ['巴西', 'brazil', '里约', '基督像'],
  },
  {
    id: 'au',
    name: '澳大利亚',
    region: '大洋洲',
    tags: ['电子签', '打工度假', '探亲'],
    note: '线上申请流程成熟',
    image:
      'https://upload.wikimedia.org/wikipedia/commons/thumb/a/a0/Sydney_Australia._%2821339175489%29.jpg/330px-Sydney_Australia._%2821339175489%29.jpg',
    keywords: ['澳大利亚', 'australia', '悉尼', '歌剧院'],
  },
  {
    id: 'nz',
    name: '新西兰',
    region: '大洋洲',
    tags: ['旅游', '探亲', '打工度假'],
    note: '自然风光线路热门，材料要覆盖行程',
    image:
      'https://upload.wikimedia.org/wikipedia/commons/thumb/b/b6/Milford_Sound_%28New_Zealand%29.JPG/330px-Milford_Sound_%28New_Zealand%29.JPG',
    keywords: ['新西兰', 'new zealand', '米尔福德峡湾'],
  },
  {
    id: 'za',
    name: '南非',
    region: '非洲',
    tags: ['旅游', '商务', '展会'],
    note: '入境材料建议预留更充分时间',
    image:
      'https://upload.wikimedia.org/wikipedia/commons/thumb/d/dc/Table_Mountain_DanieVDM.jpg/330px-Table_Mountain_DanieVDM.jpg',
    keywords: ['南非', 'south africa', '开普敦', '桌山'],
  },
]

export const planSteps = [
  { key: 'apply', label: '申请', status: 'done' },
  { key: 'docs', label: '材料', status: 'active' },
  { key: 'book', label: '预约', status: 'todo' },
  { key: 'result', label: '结果', status: 'todo' },
]

export const planTasks = [
  { id: 'passport-copy', title: '护照复印件', statusText: '已就绪', status: 'done', icon: 'description' },
  { id: 'photo', title: '证件照片', statusText: '已就绪', status: 'done', icon: 'photo_camera' },
  { id: 'bank', title: '银行流水', statusText: '审核中', status: 'review', icon: 'account_balance' },
  { id: 'flight', title: '机票行程单', statusText: '缺失', status: 'missing', icon: 'flight_takeoff' },
]

export const planGuides = [
  {
    id: 'entry-requirement',
    title: '2024 日本入境要求',
    image:
      'https://lh3.googleusercontent.com/aida-public/AB6AXuATAvtJ7P2N8NNk4wY_0QkkzTejTMdoByAzVI89xaYsrYUHGZnZHj0eGtpRUEtdPAt8RLDr2W9ZujwdevdCYxHS5K516bSq1n7o_KioFd0IyTH2KwptKROzxrWO0ue08GTF4rvJfaK68FVhCBLn89uvENVt6pIveRnqdUJsDUXfTr4UxgX_oNPGCmDKHNBW7JBvFyRUJsGfHX2GACDZr8DYEixjszjrqChzkmFWcJtoL6GtTJD-ZHM-Ep3paC44GeCPLc7nu1dUhcw',
  },
]

export const communityCategories = ['推荐', '最新', '热门', '攻略', '问答']

export const communityPosts = [
  {
    id: 'post-1',
    category: '推荐',
    title: '申根签证攻略：第一次申请就拿到3年多次！',
    author: 'ElenaTravels',
    likes: '2.4k',
    image:
      'https://lh3.googleusercontent.com/aida-public/AB6AXuC4j3paHus28MLAM6Hx6uBBdXb1pNtQKxPHTU7DQ7ln7RBoICAsKUEmFwxOR08VUUfKBoeaRMjqpFudBD0bt4uILnYwW2ErXH5Vtvscb_GW0rvc6dJUCeN3MRJF9Y31CPZJISrRD_Y7eilLkQTkyg0TWKhMnvy6gfYeYhs-Joo8TLtSM5viYNKo9V1uNwHhEmbCiSzfEhMRme_LZXFlaT8wnEGjuWOYcagcp7UUShpaxWU18UK1sFYX2SdqpA7AxunV5ewE2QmACxI',
    avatar:
      'https://lh3.googleusercontent.com/aida-public/AB6AXuCFV3WUov0kXDBe_MjS8XktMSsnQsVslLDbBOhewC8dZQ4t7ucKfP2fz38poEYQDaPtFDGm-tudnRnTlPrT7EGgnpQcp9sVGBimNjEHPDBpicXbdFCg6jeCWsBv5oxnzmPo-qNMop5paerLL6M2SRQPiae44yXT2hpkTIZ6D7phUKqwulqaPWl3HwWlbyDrrgGV-1gjmaQ8OCyphfvLseXs1bOqUGPdBag3asbgMIJTUr_CrgT30jcI4J9EjP_272oN32i8LXTZAxU',
  },
  {
    id: 'post-2',
    category: '攻略',
    title: '2024日本eVisa最新指南，官网上没说的坑',
    author: 'MarcusW',
    likes: '856',
    image:
      'https://lh3.googleusercontent.com/aida-public/AB6AXuBMi0hMIQTuZ3e5IBTq_qrVanDX_eTqo46XiOr_LnMdC7ukeGrUMA8DsjjL1Gm9_eRUKLGTd3NjNiBdIPwdwEb9YbhQPtDCR3TIYnzLwlnafdBkhD4_uPCpj4vIDSn1hXXf0C3Rczih-wbnI7RzhF_HsCTpCAgmu6As4Mvk_WcFDrjjXMYcEVAKSVBq0lN_ekd09p-8vQqStr9oDiC87SS1WfI7tlXIaGGkATOYveo3dS-SioKrEl5JJUhPfJDBXNueTaKKDV5M-LE',
    avatar:
      'https://lh3.googleusercontent.com/aida-public/AB6AXuCDMlE6CD1al4Qm_RBNmIk_VWkitiij7wF-kjcKbBeo0v00YA89aGCCApd47gpKD9ctnN4gCuFsxeiuSuaxCVwGH6ckUb0wbao9B_w88e63Tq9OzBTFDTgQOAlpZQgyvZg2m9J707aLBJWh1w8M15EZsBSfXIdPdYmavV73kR1ifEXlfNsDTWw-AAv4Zd4aJdMJTa6Cbx47eDr0K0SNkoDlYBtkCf_wEJnngzPYT4yHeQk2S_m0c3dO6fymeFFmgaLALSs5YFJcRw0',
  },
  {
    id: 'post-3',
    category: '问答',
    title: '美国B1/B2面签真的很难吗？分享今天新鲜出炉的经验！',
    author: '小陈',
    likes: '1.1k',
    image:
      'https://lh3.googleusercontent.com/aida-public/AB6AXuB1Q5ZM_9VBizOPS7tQPgdwPXUEfv06HOpZuyKdsEhHgVLCGXQNzXaF6HruG84KnvydVii0Ztaf-4WWrMs78JmADPLmkvCjzjrvURAupcngTIPn19WcOYG9EjYpyxeUzMC2-NVSB48R0eZ_yU_snyMhkKkXRaovCbS3B2hS4cAvuH4x0w1rFvlWqIzQKMOPd7-3HlUV0Dz3Guapna9r-0vbwIIfX_-ZthQShle4ZSucqN-q3OtY39Lr3opXl8Q7mdU1Ryhd8Iva-wc',
    avatar:
      'https://lh3.googleusercontent.com/aida-public/AB6AXuCbKpGCxcCXyUxyRtvq-3ZwgvUmZUbyQqlm59MRNvriTYbpR_8ofuVnmq07fd_HHGi83fDf_MdFAHIia52fvZJfneSs4nEVaymR4bdcBbFA7J0spTjojENqfM5yaKQho5rRP6gA8neULv9aac13zhkM7mGMIYMDrENuskLSKDZvqXbFWkcFtHcyhNYs4fqBQCpDkKPd-3mi88vNkm1dzSb-C_WImyI50ca0nY21ew6fA798YUFsF9DnJEXudc_F0KEOEsue6CC3Y2k',
  },
  {
    id: 'post-4',
    category: '热门',
    title: '西班牙数字游民签证：最全材料准备清单总结。',
    author: 'SarahNomad',
    likes: '3.5k',
    image:
      'https://lh3.googleusercontent.com/aida-public/AB6AXuA7DM_GziabixOFycBnmdmWideZe7CSqf0yNcHrfeV9lLKGRqtWD23WLitSv7wkvK8avJMeGa102FCdh8cETSrUpu-Yq69k6tJGV9T0XPJztUyCi53lZY91mC6gKjTYnkyNzx6G6JaFFQS_455sonl4irongBBKNiwVYFnhe17IHNvDNGWAoL35-R14ZVaqVEAEd4a0OuUTeMHjteeHjTky02obAnkeu8_-O2qKXkjHYQFLcWUhBwRJwqatQ5uvtQmRb5KbIODqHdI',
    avatar:
      'https://lh3.googleusercontent.com/aida-public/AB6AXuBdfHKXXZwCjNTXz9p8kNJOBLb79mp4niXMvn3brB76BB4sZD56KBUfKZySB-8liHbtzusDVE6DF6wVFZSGDI2eNbF3AtP_wCvRuuwK9NDLs4THZM-j401sb1DUffMPw8i8Fiexg3QwLLacuOCcFe-EmEV3H9VI2Zn0-otkX2CziRqlGidU2cB_vzC7kATsuFcwbINATPKlGY06nFAwYM7xztZyQizS77c0k1O8rarKlOgMomLoEqx-CJNOCSEYWqxrTlijNmkE7UU',
  },
]

export const profileStats = [
  { key: 'done', label: '已完成', value: '8', color: 'var(--visago-primary)' },
  { key: 'doing', label: '进行中', value: '2', color: 'var(--visago-primary)' },
  { key: 'fav', label: '收藏', value: '156', color: '#a855f7' },
]

export const profileActions = [
  { key: 'favorite', title: '我的收藏', icon: 'bookmark', tint: '#f3e8ff', color: '#9333ea' },
  { key: 'docs', title: '我的资料', icon: 'description', tint: '#fef3c7', color: '#d97706' },
  { key: 'history', title: '申请历史', icon: 'history', tint: '#dcfce7', color: '#16a34a', active: true },
]

export const profileSettings = [
  { key: 'notify', title: '通知设置', icon: 'notifications' },
  { key: 'privacy', title: '隐私与安全', icon: 'shield' },
  { key: 'theme', title: '主题切换', icon: 'dark_mode', type: 'theme' },
  { key: 'help', title: '帮助中心', icon: 'help' },
  { key: 'about', title: '关于我们', icon: 'info' },
]
