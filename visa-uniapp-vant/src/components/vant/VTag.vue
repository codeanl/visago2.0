<template>
  <text class="v-tag" :class="[`v-tag--${type}`]">{{ text }}</text>
</template>

<script>
export default {
  name: 'VTag',
  props: {
    text: {
      type: String,
      default: '',
    },
    type: {
      type: String,
      default: 'default',
    },
  },
}
</script>

<style scoped>
.v-tag {
  display: inline-block;
  padding: 4px 10px;
  border-radius: 999px;
  font-size: 11px;
  line-height: 1.2;
  font-weight: 700;
  white-space: nowrap;
  text-align: center;
  vertical-align: middle;
}

.v-tag--default {
  color: #3f506b;
  background: #f2f4f7;
}

.v-tag--primary {
  color: #1677ff;
  background: rgba(22, 119, 255, 0.12);
}

.v-tag--warning {
  color: #a05a00;
  background: rgba(255, 187, 0, 0.22);
}

.v-tag--success {
  color: #0b8f55;
  background: rgba(7, 193, 96, 0.14);
}

.v-tag--danger {
  color: #d4380d;
  background: rgba(238, 10, 36, 0.12);
}
</style>
