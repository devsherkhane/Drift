<template>
  <nav class="board-header">
    <div class="header-left">
      <h2 class="brand">Trello Clone</h2>
      <div class="divider"></div>
      <h2 class="board-title">{{ boardTitle }}</h2>
      <button class="btn-action" @click="$emit('open-share')">
        👥 Share
      </button>
      <button class="btn-action" @click="$emit('toggle-menu')">
        ☰ Menu
      </button>
    </div>

    <div class="header-right">
      <div class="search-wrapper">
        <input 
          v-model="localSearchQuery" 
          @input="onSearch"
          placeholder="Search cards..." 
          class="search-input" 
        />
        
        <!-- Search Results Dropdown -->
        <div v-if="localSearchQuery.trim()" class="search-results-dropdown">
          <div v-if="isSearching" class="search-status">Searching...</div>
          <div v-else-if="searchResults.length === 0" class="search-status">No cards found.</div>
          <div 
            v-else 
            v-for="card in searchResults" 
            :key="card.id" 
            class="search-result-item"
            @click="$emit('open-card', card, 'Search Result')"
          >
            <div class="search-card-title">{{ card.title }}</div>
            <div v-if="card.description" class="search-card-desc">
              {{ card.description.substring(0, 50) }}{{ card.description.length > 50 ? '...' : '' }}
            </div>
          </div>
        </div>
      </div>
      
      <router-link to="/profile" class="btn-icon-link" title="Profile Settings">
        👤
      </router-link>
      <button @click="$emit('logout')" class="btn-logout" title="Logout">Logout</button>
    </div>
  </nav>
</template>

<script setup>
import { ref } from 'vue';

const props = defineProps({
  boardTitle: { type: String, default: '' },
  searchResults: { type: Array, default: () => [] },
  isSearching: { type: Boolean, default: false }
});

const emit = defineEmits(['open-share', 'toggle-menu', 'open-card', 'logout', 'search']);

const localSearchQuery = ref('');
let searchTimeout = null;

const onSearch = () => {
  clearTimeout(searchTimeout);
  if (!localSearchQuery.value.trim()) {
    emit('search', '');
    return;
  }
  searchTimeout = setTimeout(() => {
    emit('search', localSearchQuery.value);
  }, 300);
};
</script>

<style scoped>
.board-header {
  padding: 10px 20px;
  background: var(--surface-glass);
  backdrop-filter: blur(24px);
  -webkit-backdrop-filter: blur(24px);
  border-bottom: 1px solid var(--border-subtle);
  color: var(--text-main);
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: sticky;
  top: 0;
  z-index: 1000;
  box-shadow: 0 4px 30px rgba(0, 0, 0, 0.05);
}

.header-left, .header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand {
  font-size: 18px;
  font-weight: 800;
  letter-spacing: -0.5px;
  margin: 0;
  color: var(--brand-primary);
}

.divider {
  width: 1px;
  height: 24px;
  background: var(--border-subtle);
  margin: 0 10px;
}

.board-title {
  font-size: 16px;
  font-weight: 600;
  margin: 0;
  margin-right: 15px;
}

.btn-action {
  background: var(--surface-primary);
  border: 1px solid var(--border-subtle);
  color: var(--text-main);
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  box-shadow: var(--shadow-soft);
}
.btn-action:hover {
  background: var(--surface-secondary);
  transform: translateY(-1px);
}

.search-wrapper {
  position: relative;
  margin-right: 10px;
}

.btn-icon-link {
  font-size: 18px;
  text-decoration: none;
  background: var(--surface-primary);
  width: 34px;
  height: 34px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: all 0.2s;
  color: var(--text-main);
  border: 1px solid var(--border-subtle);
}
.btn-icon-link:hover {
  background: var(--surface-secondary);
}

.btn-logout {
  background: rgba(239, 68, 68, 0.2);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #fca5a5;
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s;
}
.btn-logout:hover {
  background: rgba(239, 68, 68, 0.3);
  color: white;
}

.search-input {
  background: var(--surface-primary);
  border: 1px solid var(--border-subtle);
  border-radius: var(--border-radius-sm);
  padding: 8px 16px;
  color: var(--text-main);
  width: 250px;
  font-size: 14px;
  transition: all 0.2s ease;
}
.search-input::placeholder { color: var(--text-extramuted); }
.search-input:focus {
  background: var(--surface-primary);
  border-color: var(--brand-primary);
  outline: none;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15);
}

.search-results-dropdown {
  position: absolute;
  top: calc(100% + 10px);
  right: 0;
  width: 320px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  overflow-y: auto;
  z-index: 1100;
  color: var(--text-primary);
  border: 1px solid #e2e8f0;
}
[data-theme="dark"] .search-results-dropdown {
  background: var(--surface-color);
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.search-status {
  padding: 16px;
  text-align: center;
  color: var(--text-secondary);
  font-size: 14px;
}

.search-result-item {
  padding: 12px 16px;
  border-bottom: 1px solid #e2e8f0;
  cursor: pointer;
  transition: background 0.2s;
}
[data-theme="dark"] .search-result-item { border-bottom: 1px solid rgba(255, 255, 255, 0.1); }
.search-result-item:last-child { border-bottom: none; }
.search-result-item:hover { background: #f8fafc; }
[data-theme="dark"] .search-result-item:hover { background: rgba(255, 255, 255, 0.05); }

.search-card-title { font-weight: 600; font-size: 14px; margin-bottom: 4px; }
.search-card-desc { font-size: 12px; color: var(--text-secondary); }
</style>
