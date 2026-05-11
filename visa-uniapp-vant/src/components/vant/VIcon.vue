<template>
  <text class="v-icon" :class="[`v-icon--${name}`]">{{ glyph }}</text>
</template>

<script>
const ICON_MAP = {
  'arrow-left': '\u2039',
  search: '\u2315',
  close: '\u2715',
  'arrow': '\u203A',
  'home-o': '\u2302',
  'todo-list-o': '\u2316',
  'flag-o': '\u25CE',
  'browsing-history-o': '\u23F3',
  'notes-o': '\u2630',
  'bag-o': '\u25A3',
  'friends-o': '\u25C8',
  'award-o': '\u2726',
  'guide-o': '\u2708',
  passed: '\u2713',
  warning: '\u26A0',
  pending: '\u23F3',
  delete: '\u2716',
  edit: '\u270E',
  clock: '\u25F7',
  cash: '\u00A5',
  card: '\u25A4',
  user: '\u25C9',
  photo: '\u25A7',
  bookmark: '\u2311',
  success: '\u2713',
}

export default {
  name: 'VIcon',
  props: {
    name: {
      type: String,
      default: '',
    },
  },
  computed: {
    glyph() {
      const key = String(this.name || '').trim()
      return ICON_MAP[key] || '\u2022'
    },
  },
}
</script>

<style scoped>
.v-icon {
  display: inline-block;
  line-height: 1;
  font-size: inherit;
  text-align: center;
  font-style: normal;
  font-weight: 700;
  min-width: 1em;
}
</style>
