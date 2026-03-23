<template>
  <div class="activity-page animate-fade-in">
    <nav class="page-header">
      <div class="header-left">
        <ActivityIcon :size="20" class="primary-icon" />
        <h2>Activity Feed</h2>
      </div>
      <div class="header-right">
        <div class="filter-dropdown" v-if="boards.length > 0">
          <select v-model="selectedBoardId" class="board-filter" id="board-filter-select">
            <option :value="null">All Boards</option>
            <option v-for="b in boards" :key="b.id" :value="b.id">{{ b.title }}</option>
          </select>
        </div>
        <button class="btn-action-premium" @click="$router.push('/')">
          <ArrowLeftIcon :size="16" /> Dashboard
        </button>
      </div>
    </nav>

    <div class="activity-content">
      <!-- Loading -->
      <div v-if="isLoading" class="timeline-loading">
        <div class="loading-spinner"></div>
        <p>Loading activity...</p>
      </div>

      <!-- Timeline -->
      <div v-else-if="filteredActivities.length > 0" class="timeline">
        <div v-for="(group, date) in groupedActivities" :key="date" class="timeline-group">
          <div class="timeline-date-badge">
            <CalendarIcon :size="14" />
            {{ formatDateHeader(date) }}
          </div>
          <div class="timeline-items">
            <div v-for="activity in group" :key="activity.id" class="timeline-item">
              <div class="timeline-dot"></div>
              <div class="timeline-card glass-panel">
                <div class="timeline-card-header">
                  <span class="board-badge">
                    <LayoutDashboardIcon :size="12" />
                    {{ activity.board_title }}
                  </span>
                  <span class="timeline-time">
                    <ClockIcon :size="12" />
                    {{ formatTime(activity.created_at) }}
                  </span>
                </div>
                <p class="timeline-action">{{ activity.action }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Empty -->
      <div v-else class="empty-state">
        <div class="empty-icon-box">
          <ActivityIcon :size="48" />
        </div>
        <h3>No activity yet</h3>
        <p>Activity from your boards will appear here as you create cards, move items, and collaborate.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import api from '../api';
import {
  Activity as ActivityIcon,
  ArrowLeft as ArrowLeftIcon,
  Calendar as CalendarIcon,
  Clock as ClockIcon,
  LayoutDashboard as LayoutDashboardIcon
} from 'lucide-vue-next';

const boards = ref([]);
const allActivities = ref([]);
const isLoading = ref(true);
const selectedBoardId = ref(null);

const filteredActivities = computed(() => {
  if (!selectedBoardId.value) return allActivities.value;
  return allActivities.value.filter(a => a.board_id === selectedBoardId.value);
});

const groupedActivities = computed(() => {
  const groups = {};
  filteredActivities.value.forEach(activity => {
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
  const date = new Date(dateStr);
  return date.toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' });
};

const fetchAll = async () => {
  isLoading.value = true;
  try {
    const boardsRes = await api.get('/boards');
    boards.value = boardsRes.data || [];

    const activityPromises = boards.value.map(async (board) => {
      try {
        const res = await api.get(`/boards/${board.id}/activity`);
        return (res.data || []).map(a => ({ ...a, board_title: board.title, board_id: board.id }));
      } catch {
        return [];
      }
    });

    const results = await Promise.all(activityPromises);
    allActivities.value = results
      .flat()
      .sort((a, b) => new Date(b.created_at) - new Date(a.created_at));
  } catch (err) {
    console.error('Failed to fetch activity', err);
  } finally {
    isLoading.value = false;
  }
};

onMounted(fetchAll);
</script>

<style scoped>
.activity-page {
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

.header-right { display: flex; align-items: center; gap: 12px; }

.board-filter {
  padding: 8px 16px;
  background: var(--surface-secondary);
  border: 1px solid var(--border-subtle);
  border-radius: 10px;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-main);
  cursor: pointer;
  outline: none;
}

.board-filter:focus { border-color: var(--brand-primary); }

.btn-action-premium {
  background: var(--brand-primary-light);
  color: var(--brand-primary);
  border: none;
  padding: 8px 16px;
  border-radius: 10px;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: all 0.2s;
}

.btn-action-premium:hover { background: var(--brand-primary); color: white; transform: translateY(-1px); }

.activity-content {
  flex: 1;
  padding: 40px 48px;
  max-width: 740px;
  margin: 0 auto;
  width: 100%;
}

/* Timeline */
.timeline { position: relative; }

.timeline-group { margin-bottom: 40px; }

.timeline-date-badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  background: var(--brand-primary-light);
  color: var(--brand-primary);
  padding: 6px 16px;
  border-radius: 20px;
  font-size: 13px;
  font-weight: 700;
  margin-bottom: 20px;
}

.timeline-items {
  position: relative;
  padding-left: 28px;
  border-left: 2px solid var(--border-subtle);
}

.timeline-item {
  position: relative;
  margin-bottom: 16px;
}

.timeline-dot {
  position: absolute;
  left: -35px;
  top: 20px;
  width: 12px;
  height: 12px;
  background: var(--brand-primary);
  border-radius: 50%;
  border: 3px solid var(--surface-primary);
  box-shadow: var(--shadow-sm);
}

.timeline-card {
  background: var(--surface-primary);
  border: 1px solid var(--border-subtle);
  border-radius: 14px;
  padding: 18px 22px;
  transition: all 0.2s;
}

.timeline-card:hover {
  transform: translateX(4px);
  box-shadow: var(--shadow-md);
}

.timeline-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
}

.board-badge {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  font-weight: 700;
  color: var(--brand-primary);
  background: var(--brand-primary-light);
  padding: 4px 10px;
  border-radius: 8px;
}

.timeline-time {
  font-size: 12px;
  color: var(--text-extramuted);
  display: flex;
  align-items: center;
  gap: 4px;
  font-weight: 600;
}

.timeline-action {
  font-size: 14px;
  color: var(--text-main);
  line-height: 1.5;
  margin: 0;
}

/* Loading */
.timeline-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 100px 20px;
  gap: 16px;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 4px solid var(--border-subtle);
  border-top-color: var(--brand-primary);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

.timeline-loading p { color: var(--text-muted); font-weight: 600; }

/* Empty */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 100px 20px;
  text-align: center;
}

.empty-icon-box {
  width: 100px;
  height: 100px;
  background: var(--brand-primary-light);
  border-radius: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--brand-primary);
  margin-bottom: 24px;
}

.empty-state h3 { font-size: 20px; font-weight: 800; color: var(--text-main); margin-bottom: 8px; }
.empty-state p { font-size: 15px; color: var(--text-muted); max-width: 360px; line-height: 1.6; }

.animate-fade-in { animation: fadeIn 0.6s cubic-bezier(0.16, 1, 0.3, 1); }
@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
</style>
