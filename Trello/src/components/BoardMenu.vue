<template>
  <div>
    <Transition name="fade">
      <div v-if="isOpen" class="menu-overlay" @click="$emit('close')"></div>
    </Transition>
    <div class="board-menu glass-panel" :class="{ 'menu-open': isOpen }">
      <div class="menu-header">
        <div class="header-title-group">
          <ActivityIcon :size="18" class="primary-icon" />
          <h3>Board Menu</h3>
        </div>
        <button class="menu-close-btn" @click="$emit('close')">
          <XIcon :size="18" />
        </button>
      </div>
      <div class="menu-content">
        <section class="menu-section">
          <div class="section-label">Actions</div>
          <button class="btn-danger-soft w-100" @click="$emit('archive-board')">
            <Trash2Icon :size="16" /> Archive Board
          </button>
        </section>
        
        <section class="menu-section">
          <div class="section-label">Recent Activity</div>
          <div class="activity-feed">
            <div v-if="activityLogs.length === 0" class="empty-state-logs">
              <span class="empty-icon">🔔</span>
              <p>No recent activity found.</p>
            </div>
            <div v-for="(log, index) in activityLogs" :key="index" class="activity-log-item">
              <div class="log-indicator"></div>
              <div class="log-details">
                <p class="log-text">{{ log.action }}</p>
                <span class="log-timestamp">{{ formatTimeAgo(log.created_at) }}</span>
              </div>
            </div>
          </div>
        </section>
      </div>
    </div>
  </div>
</template>

<script setup>
import { X as XIcon, Activity as ActivityIcon, Trash2 as Trash2Icon } from 'lucide-vue-next';

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
  background: rgba(15, 23, 42, 0.4); backdrop-filter: blur(4px);
  z-index: 1500;
}

.board-menu {
  position: fixed; top: 0; right: 0; width: 340px; height: 100vh;
  background: var(--surface-primary); 
  z-index: 1600; transform: translateX(100%); 
  transition: transform 0.4s cubic-bezier(0.16, 1, 0.3, 1);
  display: flex; flex-direction: column;
  box-shadow: var(--shadow-strong);
  border-left: 1px solid var(--border-subtle);
}
.board-menu.menu-open { transform: translateX(0); }

.menu-header {
  display: flex; justify-content: space-between; align-items: center;
  padding: 24px; border-bottom: 1px solid var(--border-subtle);
}

.header-title-group { display: flex; align-items: center; gap: 10px; }
.header-title-group h3 { margin: 0; font-size: 18px; font-weight: 800; color: var(--text-main); }
.primary-icon { color: var(--brand-primary); }

.menu-close-btn {
  background: var(--surface-secondary); border: none; padding: 8px;
  border-radius: 50%; cursor: pointer; color: var(--text-muted);
  display: flex; align-items: center; justify-content: center;
  transition: all 0.2s;
}
.menu-close-btn:hover { background: #fee2e2; color: #ef4444; }

.menu-content { flex: 1; overflow-y: auto; padding: 24px; background: var(--surface-primary); }
.menu-section { margin-bottom: 36px; }

.section-label { 
  font-size: 11px; font-weight: 800; text-transform: uppercase; 
  letter-spacing: 1px; color: var(--text-muted); margin-bottom: 16px; 
}

.btn-danger-soft {
  background: #fef2f2; color: #ef4444; border: 1px solid transparent;
  padding: 12px; border-radius: 12px;
  cursor: pointer; font-weight: 700; transition: all 0.2s;
  display: flex; align-items: center; justify-content: center; gap: 8px;
}
.btn-danger-soft:hover { background: #fee2e2; transform: translateY(-1px); }
.w-100 { width: 100%; }

.activity-feed { display: flex; flex-direction: column; gap: 4px; }

.activity-log-item { 
  display: flex; gap: 14px; padding: 12px 0; 
  position: relative;
}
.activity-log-item:not(:last-child)::after {
  content: ''; position: absolute; left: 4px; top: 28px; bottom: -8px; 
  width: 1px; background: var(--border-subtle);
}

.log-indicator {
  width: 9px; height: 9px; border-radius: 50%; background: var(--brand-primary);
  margin-top: 6px; z-index: 1; border: 2px solid var(--surface-primary);
}

.log-details { flex-grow: 1; }
.log-text { margin: 0 0 4px 0; font-size: 13.5px; color: var(--text-main); line-height: 1.5; font-weight: 500; }
.log-timestamp { font-size: 11px; color: var(--text-muted); font-weight: 600; }

.empty-state-logs { 
  text-align: center; padding: 40px 0; color: var(--text-extramuted); 
}
.empty-icon { font-size: 24px; display: block; margin-bottom: 12px; opacity: 0.3; }

.fade-enter-active, .fade-leave-active { transition: opacity 0.3s; }
.fade-enter-from, .fade-leave-to { opacity: 0; }
</style>
