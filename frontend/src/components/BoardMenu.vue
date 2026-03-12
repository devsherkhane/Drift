<template>
  <div>
    <div v-if="isOpen" class="menu-overlay" @click="$emit('close')"></div>
    <div class="board-menu" :class="{ 'menu-open': isOpen }">
      <div class="menu-header">
        <h3>Menu</h3>
        <button class="btn-close" @click="$emit('close')">✕</button>
      </div>
      <div class="menu-content">
        <section class="menu-section">
          <button class="btn-danger w-100" @click="$emit('archive-board')">
            🗑️ Archive Board
          </button>
        </section>
        
        <section class="menu-section">
          <h4>Activity</h4>
          <div class="activity-feed">
            <div v-if="activityLogs.length === 0" class="empty-state">
              No recent activity.
            </div>
            <div v-for="(log, index) in activityLogs" :key="index" class="activity-item">
              <div class="activity-content">
                <p>{{ log.action }}</p>
                <span class="activity-time">{{ formatTimeAgo(log.created_at) }}</span>
              </div>
            </div>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  isOpen: { type: Boolean, default: false },
  activityLogs: { type: Array, default: () => [] }
});

const emit = defineEmits(['close', 'archive-board']);

const formatTimeAgo = (dateString) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  const now = new Date();
  const diffInSeconds = Math.floor((now - date) / 1000);
  
  if (diffInSeconds < 60) return `${diffInSeconds}s ago`;
  if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)}m ago`;
  if (diffInSeconds < 86400) return `${Math.floor(diffInSeconds / 3600)}h ago`;
  return date.toLocaleDateString();
};
</script>

<style scoped>
.menu-overlay {
  position: fixed; top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(15, 23, 42, 0.3); z-index: 80;
}

.board-menu {
  position: fixed; top: 0; right: 0; width: 320px; height: 100vh;
  background: white; box-shadow: -4px 0 15px rgba(0,0,0,0.1);
  z-index: 90; transform: translateX(100%); transition: transform 0.3s cubic-bezier(0.16, 1, 0.3, 1);
  display: flex; flex-direction: column;
}
.board-menu.menu-open { transform: translateX(0); }

.menu-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 20px; border-bottom: 1px solid #e2e8f0;
}
.menu-header h3 { margin: 0; color: var(--text-primary); }

.btn-close {
  background: none; border: none; font-size: 20px; cursor: pointer;
  color: var(--text-secondary); padding: 4px 8px; border-radius: var(--border-radius-sm);
}
.btn-close:hover { background: var(--surface-secondary); color: var(--text-main); }

.menu-content { flex: 1; overflow-y: auto; padding: 20px; background: #f8fafc; }
.menu-section { margin-bottom: 30px; }
.menu-section h4 { margin: 0 0 15px 0; color: var(--text-secondary); text-transform: uppercase; font-size: 13px; letter-spacing: 0.5px; }

.btn-danger {
  background: #fef2f2; color: #ef4444; border: 1px solid #fee2e2;
  padding: 8px 16px; border-radius: var(--border-radius-sm);
  cursor: pointer; font-weight: 600; transition: all 0.2s;
}
.btn-danger:hover { background: #fee2e2; border-color: #fca5a5; }
.w-100 { width: 100%; }

.activity-feed { display: flex; flex-direction: column; gap: 15px; }
.activity-item { display: flex; gap: 12px; padding-bottom: 15px; border-bottom: 1px solid #e2e8f0; }
.activity-content p { margin: 0 0 4px 0; font-size: 14px; color: var(--text-primary); line-height: 1.4; }
.activity-time { font-size: 12px; color: var(--text-secondary); }
.empty-state { text-align: center; color: var(--text-secondary); font-size: 14px; padding: 20px 0; }
</style>
