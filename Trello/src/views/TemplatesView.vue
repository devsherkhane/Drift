<template>
  <div class="templates-page animate-fade-in">
    <nav class="page-header">
      <div class="header-left">
        <CopyIcon :size="20" class="primary-icon" />
        <h2>Board Templates</h2>
      </div>
      <div class="header-right">
        <button class="btn-action-premium" @click="$router.push('/')">
          <ArrowLeftIcon :size="16" /> Dashboard
        </button>
      </div>
    </nav>

    <div class="templates-content">
      <div class="templates-hero">
        <h1>Start from a template</h1>
        <p>Get up and running fast with a pre-built board structure. Click a template to create your board instantly.</p>
      </div>

      <div class="templates-grid">
        <div
          v-for="(template, i) in templates"
          :key="i"
          class="template-card glass-panel"
          :class="{ 'is-creating': creatingIndex === i }"
          @click="createFromTemplate(template, i)"
        >
          <div class="template-icon-box" :style="{ background: template.iconBg }">
            <component :is="template.icon" :size="28" :style="{ color: template.iconColor }" />
          </div>
          <div class="template-info">
            <h3>{{ template.name }}</h3>
            <p>{{ template.description }}</p>
          </div>
          <div class="template-lists-preview">
            <span class="preview-label">Lists:</span>
            <div class="preview-tags">
              <span v-for="(list, j) in template.lists" :key="j" class="preview-tag">{{ list }}</span>
            </div>
          </div>
          <div class="template-footer">
            <span v-if="creatingIndex === i" class="creating-label">
              <div class="mini-spinner"></div> Creating...
            </span>
            <span v-else class="use-template-label">
              <PlusIcon :size="14" /> Use Template
            </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import api from '../api';
import { useToast } from 'vue-toastification';
import {
  Copy as CopyIcon,
  ArrowLeft as ArrowLeftIcon,
  Plus as PlusIcon,
  Kanban as KanbanIcon,
  Bug as BugIcon,
  Calendar as CalendarIcon,
  User as UserIcon,
  Briefcase as BriefcaseIcon,
  Rocket as RocketIcon
} from 'lucide-vue-next';

const router = useRouter();
const toast = useToast();
const creatingIndex = ref(null);

const templates = [
  {
    name: 'Project Management',
    description: 'Track project milestones, tasks, and deliverables from start to finish.',
    icon: BriefcaseIcon,
    iconBg: '#e0e7ff',
    iconColor: '#6366f1',
    lists: ['Backlog', 'To Do', 'In Progress', 'Review', 'Done']
  },
  {
    name: 'Sprint Board',
    description: 'Agile sprint planning with clear workflow stages for development teams.',
    icon: RocketIcon,
    iconBg: '#d1fae5',
    iconColor: '#10b981',
    lists: ['Sprint Backlog', 'In Development', 'Code Review', 'QA Testing', 'Deployed']
  },
  {
    name: 'Bug Tracker',
    description: 'Capture, prioritize, and squash bugs systematically across your product.',
    icon: BugIcon,
    iconBg: '#fee2e2',
    iconColor: '#ef4444',
    lists: ['New Bugs', 'Triaged', 'In Progress', 'Fixed (Needs QA)', 'Closed']
  },
  {
    name: 'Content Calendar',
    description: 'Plan and schedule content across channels with a clear editorial pipeline.',
    icon: CalendarIcon,
    iconBg: '#fef3c7',
    iconColor: '#f59e0b',
    lists: ['Ideas', 'Drafting', 'Editing', 'Scheduled', 'Published']
  },
  {
    name: 'Personal Tasks',
    description: 'A simple board to manage your daily tasks and personal goals.',
    icon: UserIcon,
    iconBg: '#f3e8ff',
    iconColor: '#8b5cf6',
    lists: ['To Do', 'Doing', 'Done']
  },
  {
    name: 'Kanban Board',
    description: 'Classic Kanban workflow for any team or individual project.',
    icon: KanbanIcon,
    iconBg: '#fce7f3',
    iconColor: '#ec4899',
    lists: ['To Do', 'In Progress', 'Done']
  }
];

const createFromTemplate = async (template, index) => {
  if (creatingIndex.value !== null) return;
  creatingIndex.value = index;

  try {
    // Create the board
    const boardRes = await api.post('/boards', { title: template.name });
    const boardId = boardRes.data.id;

    // Create lists in order
    for (let i = 0; i < template.lists.length; i++) {
      await api.post('/lists', {
        board_id: boardId,
        title: template.lists[i],
        position: i + 1
      });
    }

    toast.success(`"${template.name}" board created!`);
    router.push(`/board/${boardId}`);
  } catch (err) {
    toast.error('Failed to create board from template');
  } finally {
    creatingIndex.value = null;
  }
};
</script>

<style scoped>
.templates-page {
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

.templates-content {
  flex: 1; padding: 48px;
  max-width: 1100px; margin: 0 auto; width: 100%;
}

.templates-hero {
  text-align: center; margin-bottom: 56px;
}

.templates-hero h1 {
  font-size: 36px; font-weight: 900; color: var(--text-main);
  letter-spacing: -2px; margin-bottom: 12px;
}

.templates-hero p {
  font-size: 17px; color: var(--text-muted);
  max-width: 500px; margin: 0 auto; line-height: 1.6;
}

.templates-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 24px;
}

.template-card {
  background: var(--surface-primary);
  border: 1px solid var(--border-subtle);
  border-radius: 20px;
  padding: 28px;
  cursor: pointer;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.template-card:hover {
  transform: translateY(-6px);
  box-shadow: var(--shadow-lg);
  border-color: var(--brand-primary);
}

.template-card.is-creating {
  opacity: 0.7;
  pointer-events: none;
}

.template-icon-box {
  width: 56px; height: 56px; border-radius: 16px;
  display: flex; align-items: center; justify-content: center;
}

.template-info h3 {
  font-size: 18px; font-weight: 800; color: var(--text-main);
  letter-spacing: -0.5px; margin: 0 0 6px;
}

.template-info p {
  font-size: 13px; color: var(--text-muted); line-height: 1.6; margin: 0;
}

.template-lists-preview {
  display: flex; flex-direction: column; gap: 8px;
}

.preview-label {
  font-size: 11px; font-weight: 700; color: var(--text-extramuted);
  text-transform: uppercase; letter-spacing: 0.5px;
}

.preview-tags {
  display: flex; flex-wrap: wrap; gap: 6px;
}

.preview-tag {
  background: var(--surface-secondary); color: var(--text-muted);
  font-size: 12px; font-weight: 600; padding: 4px 12px;
  border-radius: 8px; border: 1px solid var(--border-subtle);
}

.template-footer {
  margin-top: auto; padding-top: 8px;
  border-top: 1px solid var(--border-subtle);
}

.use-template-label {
  display: flex; align-items: center; gap: 6px;
  color: var(--brand-primary); font-size: 13px; font-weight: 700;
}

.creating-label {
  display: flex; align-items: center; gap: 8px;
  color: var(--text-muted); font-size: 13px; font-weight: 700;
}

.mini-spinner {
  width: 16px; height: 16px;
  border: 2px solid var(--border-subtle);
  border-top-color: var(--brand-primary);
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin { to { transform: rotate(360deg); } }

.animate-fade-in { animation: fadeIn 0.6s cubic-bezier(0.16, 1, 0.3, 1); }
@keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }

@media (max-width: 700px) {
  .templates-grid { grid-template-columns: 1fr; }
  .templates-content { padding: 24px; }
}
</style>
