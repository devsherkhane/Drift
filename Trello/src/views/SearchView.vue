<template>
  <div class="search-page animate-fade-in">
    <nav class="page-header">
      <div class="header-left">
        <SearchIcon :size="20" class="primary-icon" />
        <h2>Search</h2>
      </div>
      <div class="header-right">
        <button class="btn-action-premium" @click="$router.push('/')">
          <ArrowLeftIcon :size="16" /> Back to Dashboard
        </button>
      </div>
    </nav>

    <div class="search-content">
      <div class="search-box-wrapper">
        <div class="search-input-container glass-panel shadow-md">
          <SearchIcon :size="20" class="search-input-icon" />
          <input
            v-model="query"
            type="text"
            placeholder="Search cards across all your boards..."
            @input="debouncedSearch"
            ref="searchInput"
            autofocus
            id="global-search-input"
          />
          <kbd v-if="!query" class="search-kbd">⌘K</kbd>
          <button v-if="query" class="search-clear" @click="clearSearch">
            <XIcon :size="16" />
          </button>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="isLoading" class="results-container">
        <div class="skeleton-card glass-panel" v-for="n in 4" :key="n">
          <div class="skeleton-line skeleton-title"></div>
          <div class="skeleton-line skeleton-short"></div>
          <div class="skeleton-line skeleton-meta"></div>
        </div>
      </div>

      <!-- Results -->
      <div v-else-if="results.length > 0" class="results-container">
        <p class="results-count">
          <FileTextIcon :size="16" />
          Found <strong>{{ results.length }}</strong> result{{ results.length !== 1 ? 's' : '' }}
        </p>
        <div class="results-grid">
          <div
            v-for="card in results"
            :key="card.id"
            class="result-card glass-panel"
            @click="navigateToCard(card)"
          >
            <div class="result-label-strip" v-if="card.label_color" :style="{ background: card.label_color }"></div>
            <div class="result-body">
              <h3 class="result-title">{{ card.title }}</h3>
              <p v-if="card.description" class="result-desc">{{ card.description }}</p>
              <div class="result-meta">
                <span class="meta-tag">
                  <LayoutDashboardIcon :size="12" /> {{ card.board_name }}
                </span>
                <span class="meta-tag">
                  <ListIcon :size="12" /> {{ card.list_name }}
                </span>
              </div>
            </div>
            <ChevronRightIcon :size="18" class="result-arrow" />
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else-if="query && hasSearched" class="empty-state">
        <div class="empty-icon-box">
          <SearchXIcon :size="48" />
        </div>
        <h3>No results found</h3>
        <p>Try a different search term or check your spelling.</p>
      </div>

      <!-- Initial State -->
      <div v-else-if="!query" class="empty-state initial-state">
        <div class="empty-icon-box">
          <SearchIcon :size="48" />
        </div>
        <h3>Search across all your boards</h3>
        <p>Find cards by title or description instantly.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import api from '../api';
import {
  Search as SearchIcon,
  ArrowLeft as ArrowLeftIcon,
  X as XIcon,
  FileText as FileTextIcon,
  LayoutDashboard as LayoutDashboardIcon,
  List as ListIcon,
  ChevronRight as ChevronRightIcon,
  SearchX as SearchXIcon
} from 'lucide-vue-next';

const query = ref('');
const results = ref([]);
const isLoading = ref(false);
const hasSearched = ref(false);
const searchInput = ref(null);

let debounceTimer = null;

const debouncedSearch = () => {
  clearTimeout(debounceTimer);
  if (!query.value.trim()) {
    results.value = [];
    hasSearched.value = false;
    return;
  }
  debounceTimer = setTimeout(() => {
    performSearch();
  }, 350);
};

const performSearch = async () => {
  if (!query.value.trim()) return;
  isLoading.value = true;
  hasSearched.value = true;
  try {
    const response = await api.get('/search', { params: { q: query.value } });
    results.value = response.data || [];
  } catch (err) {
    results.value = [];
  } finally {
    isLoading.value = false;
  }
};

const clearSearch = () => {
  query.value = '';
  results.value = [];
  hasSearched.value = false;
  searchInput.value?.focus();
};

const navigateToCard = (card) => {
  // Navigate to the board; the card ID could be used later with modal
  // For now we just go to the board
  // We need to find the board ID — it's not returned by the search API directly
  // So we just show a toast or do nothing fancy
};

onMounted(() => {
  searchInput.value?.focus();
});
</script>

<style scoped>
.search-page {
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

.header-left {
  display: flex;
  align-items: center;
  gap: 12px;
}

.header-left h2 {
  color: var(--text-main);
  margin: 0;
  font-size: 18px;
  font-weight: 800;
  letter-spacing: -0.5px;
}

.primary-icon { color: var(--brand-primary); }

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
  text-decoration: none;
}

.btn-action-premium:hover {
  background: var(--brand-primary);
  color: white;
  transform: translateY(-1px);
}

.search-content {
  flex: 1;
  padding: 40px 48px;
  max-width: 800px;
  margin: 0 auto;
  width: 100%;
}

.search-box-wrapper {
  margin-bottom: 40px;
}

.search-input-container {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px 24px;
  background: var(--surface-primary);
  border-radius: 18px;
  border: 2px solid transparent;
  transition: all 0.25s;
}

.search-input-container:focus-within {
  border-color: var(--brand-primary);
  box-shadow: 0 0 0 4px var(--brand-primary-light);
}

.search-input-icon { color: var(--text-extramuted); flex-shrink: 0; }

.search-input-container input {
  flex: 1;
  border: none;
  background: transparent;
  font-size: 16px;
  color: var(--text-main);
  outline: none;
  font-weight: 500;
}

.search-input-container input::placeholder { color: var(--text-extramuted); }

.search-kbd {
  background: var(--surface-secondary);
  color: var(--text-extramuted);
  padding: 4px 10px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 700;
  border: 1px solid var(--border-subtle);
}

.search-clear {
  background: var(--surface-secondary);
  border: none;
  color: var(--text-muted);
  width: 32px;
  height: 32px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.2s;
}

.search-clear:hover { background: #fee2e2; color: #ef4444; }

/* Results */
.results-count {
  font-size: 14px;
  color: var(--text-muted);
  font-weight: 600;
  margin-bottom: 20px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.results-grid {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.result-card {
  background: var(--surface-primary);
  border: 1px solid var(--border-subtle);
  border-radius: 16px;
  padding: 20px 24px;
  display: flex;
  align-items: center;
  gap: 16px;
  cursor: pointer;
  transition: all 0.25s;
  position: relative;
  overflow: hidden;
}

.result-card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
  border-color: var(--brand-primary-light);
}

.result-label-strip {
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 4px;
  border-radius: 4px 0 0 4px;
}

.result-body { flex: 1; min-width: 0; }

.result-title {
  font-size: 16px;
  font-weight: 700;
  color: var(--text-main);
  margin: 0 0 6px;
  letter-spacing: -0.3px;
}

.result-desc {
  font-size: 13px;
  color: var(--text-muted);
  margin: 0 0 10px;
  line-height: 1.5;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.result-meta {
  display: flex;
  gap: 12px;
}

.meta-tag {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  font-weight: 600;
  color: var(--text-extramuted);
  background: var(--surface-secondary);
  padding: 4px 10px;
  border-radius: 8px;
}

.result-arrow { color: var(--text-extramuted); flex-shrink: 0; }

/* Empty/Initial State */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 80px 20px;
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

.empty-state h3 {
  font-size: 20px;
  font-weight: 800;
  color: var(--text-main);
  margin-bottom: 8px;
}

.empty-state p {
  font-size: 15px;
  color: var(--text-muted);
}

/* Skeleton Loading */
.skeleton-card {
  background: var(--surface-primary);
  border-radius: 16px;
  padding: 24px;
  margin-bottom: 12px;
}

.skeleton-line {
  border-radius: 8px;
  background: linear-gradient(90deg, var(--surface-secondary) 25%, var(--border-subtle) 50%, var(--surface-secondary) 75%);
  background-size: 200% 100%;
  animation: shimmer 1.5s infinite;
}

.skeleton-title { height: 20px; width: 60%; margin-bottom: 10px; }
.skeleton-short { height: 14px; width: 80%; margin-bottom: 10px; }
.skeleton-meta { height: 14px; width: 40%; }

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

.animate-fade-in {
  animation: fadeIn 0.6s cubic-bezier(0.16, 1, 0.3, 1);
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
