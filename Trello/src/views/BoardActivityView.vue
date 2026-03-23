<template>
  <div class="board-activity-page animate-fade-in">
    <nav class="page-header">
      <div class="header-left">
        <ClockIcon :size="20" class="primary-icon" />
        <h2>{{ boardTitle }} — Activity</h2>
      </div>
      <div class="header-right">
        <button class="btn-action-premium" @click="$router.push(`/board/${boardId}`)">
          <ArrowLeftIcon :size="16" /> Back to Board
        </button>
      </div>
    </nav>

    <div class="activity-content">
      <!-- Loading -->
      <div v-if="isLoading" class="timeline-loading">
        <div class="loading-spinner"></div>
        <p>Loading board activity...</p>
      </div>

      <!-- Timeline -->
      <div v-else-if="activities.length > 0" class="timeline">
        <div v-for="(group, date) in groupedActivities" :key="date" class="timeline-group">
          <div class="timeline-date-badge">
            <CalendarIcon :size="14" />
            {{ formatDateHeader(date) }}
          </div>
          <div class="timeline-items">
            <div v-for="(activity, idx) in group" :key="activity.id || idx" class="timeline-item">
              <div class="timeline-dot"></div>
              <div class="timeline-card glass-panel">
                <p class="timeline-action">{{ activity.action }}</p>
                <span class="timeline-time">
                  <ClockIcon :size="11" />
                  {{ formatTime(activity.created_at) }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty -->
      <div v-else class="empty-state">
        <div class="empty-icon-box">
          <ClockIcon :size="48" />
        </div>
        <h3>No activity for this board yet</h3>
        <p>Start adding cards, moving items, or inviting collaborators to see activity here.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import api from '../api';
import {
  Clock as ClockIcon,
  ArrowLeft as ArrowLeftIcon,
  Calendar as CalendarIcon
} from 'lucide-vue-next';

const route = useRoute();
const boardId = ref(route.params.id);
const boardTitle = ref('Board');
const activities = ref([]);
const isLoading = ref(true);

const groupedActivities = computed(() => {
  const groups = {};
  activities.value.forEach(activity => {
    const date = activity.created_at ? activity.created_at.split('T')[0] : 'Unknown';
    if (!groups[date]) groups[date] = [];
    groups[date].push(activity);
  });
  return groups;
});

const formatDateHeader = (dateStr) => {
  if (dateStr === 'Unknown') return 'Unknown Date';
  const date = new Date(dateStr);
  const today = new Date();
  const yesterday = new Date(today);
  yesterday.setDate(yesterday.getDate() - 1);
  if (date.toDateString() === today.toDateString()) return 'Today';
  if (date.toDateString() === yesterday.toDateString()) return 'Yesterday';
  return date.toLocaleDateString('en-US', { weekday: 'long', month: 'short', day: 'numeric' });
};

const formatTime = (dateStr) => {
  if (!dateStr) return '';
  return new Date(dateStr).toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' });
};

onMounted(async () => {
  try {
    const [boardRes, activityRes] = await Promise.all([
      api.get(`/boards/${boardId.value}`),
      api.get(`/boards/${boardId.value}/activity`)
    ]);
    boardTitle.value = boardRes.data?.title || 'Board';
    activities.value = (activityRes.data || []).sort(
      (a, b) => new Date(b.created_at) - new Date(a.created_at)
    );
  } catch (err) {
    console.error('Failed to load board activity', err);
  } finally {
    isLoading.value = false;
  }
});
</script>

<style scoped>
.board-activity-page {
  min-height: 100vh;
  background: var(--bg-gradient);
  display: flex;
  flex-direction: column;
}

.page-header {
  padding: 12px 32px;
  background: var(--surface-glass);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-bottom: 1px solid var(--border-subtle);
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-left { display: flex; align-items: center; gap: 12px; }
.header-left h2 { color: var(--text-main); margin: 0; font-size: 18px; font-weight: 800; letter-spacing: -0.5px; }
.primary-icon { color: var(--brand-primary); }

.btn-action-premium {
  background: var(--brand-primary-light); color: var(--brand-primary);
  border: none; padding: 8px 16px; border-radius: 10px;
  font-size: 13px; font-weight: 700; cursor: pointer;
  display: flex; align-items: center; gap: 8px; transition: all 0.2s;
}
.btn-action-premium:hover { background: var(--brand-primary); color: white; transform: translateY(-1px); }

.activity-content {
  flex: 1; padding: 40px 48px;
  max-width: 700px; margin: 0 auto; width: 100%;
}

/* Timeline */
.timeline-group { margin-bottom: 40px; }

.timeline-date-badge {
  display: inline-flex; align-items: center; gap: 8px;
  background: var(--brand-primary-light); color: var(--brand-primary);
  padding: 6px 16px; border-radius: 20px; font-size: 13px; font-weight: 700;
  margin-bottom: 20px;
}

.timeline-items {
  position: relative; padding-left: 28px;
  border-left: 2px solid var(--border-subtle);
}

.timeline-item { position: relative; margin-bottom: 14px; }

.timeline-dot {
  position: absolute; left: -35px; top: 16px;
  width: 12px; height: 12px; background: var(--brand-primary);
  border-radius: 50%; border: 3px solid var(--surface-primary);
  box-shadow: var(--shadow-sm);
}

.timeline-card {
  background: var(--surface-primary); border: 1px solid var(--border-subtle);
  border-radius: 14px; padding: 16px 20px; transition: all 0.2s;
  display: flex; justify-content: space-between; align-items: center; gap: 16px;
}

.timeline-card:hover { transform: translateX(4px); box-shadow: var(--shadow-md); }

.timeline-action {
  font-size: 14px; color: var(--text-main); line-height: 1.5; margin: 0; flex: 1;
}

.timeline-time {
  font-size: 11px; color: var(--text-extramuted); display: flex;
  align-items: center; gap: 4px; font-weight: 600; white-space: nowrap;
}

/* Loading */
.timeline-loading {
  display: flex; flex-direction: column; align-items: center;
  justify-content: center; padding: 100px 20px; gap: 16px;
}
.loading-spinner {
  width: 40px; height: 40px; border: 4px solid var(--border-subtle);
  border-top-color: var(--brand-primary); border-radius: 50%;
  animation: spin 0.8s linear infinite;
}
@keyframes spin { to { transform: rotate(360deg); } }
.timeline-loading p { color: var(--text-muted); font-weight: 600; }

/* Empty */
.empty-state {
  display: flex; flex-direction: column; align-items: center;
  justify-content: center; padding: 100px 20px; text-align: center;
}
.empty-icon-box {
  width: 100px; height: 100px; background: var(--brand-primary-light);
  border-radius: 28px; display: flex; align-items: center;
  justify-content: center; color: var(--brand-primary); margin-bottom: 24px;
}
.empty-state h3 { font-size: 20px; font-weight: 800; color: var(--text-main); margin-bottom: 8px; }
.empty-state p { font-size: 15px; color: var(--text-muted); max-width: 360px; line-height: 1.6; }

.animate-fade-in { animation: fadeIn 0.6s cubic-bezier(0.16, 1, 0.3, 1); }
@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
</style>
