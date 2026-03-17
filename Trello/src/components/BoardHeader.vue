<template>
  <nav class="board-header glass-panel">
    <div class="header-left">
      <router-link to="/" class="brand-link" title="Dashboard">
        <h2 class="brand">Trello Clone</h2>
      </router-link>
      <div class="divider"></div>
      <h2 class="board-title">{{ boardTitle }}</h2>
      <button class="btn-action-premium" @click="$emit('open-share')">
        <UsersIcon :size="16" />
        <span>Share</span>
      </button>
      <button class="btn-action-premium" @click="$emit('toggle-menu')">
        <MenuIcon :size="16" />
        <span>Activity</span>
      </button>
    </div>

    <div class="header-right">
      <div class="search-wrapper">
        <div class="search-input-container">
          <SearchIcon :size="14" class="search-icon-inside" />
          <input 
            v-model="localSearchQuery" 
            @input="onSearch"
            placeholder="Search cards..." 
            class="search-input" 
          />
        </div>
        
        <!-- Search Results Dropdown -->
        <div v-if="localSearchQuery.trim()" class="search-results-dropdown glass-panel shadow-lg">
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
      
      <router-link to="/profile" class="btn-nav-icon" title="Profile Settings">
        <UserIcon :size="18" />
      </router-link>
      <button @click="$emit('logout')" class="btn-nav-icon btn-logout-hover" title="Logout">
        <LogOutIcon :size="18" />
      </button>
    </div>
  </nav>
</template>

<script setup>
import { ref } from 'vue';
import { 
  Users as UsersIcon, 
  Menu as MenuIcon, 
  Search as SearchIcon, 
  User as UserIcon, 
  LogOut as LogOutIcon 
} from 'lucide-vue-next';

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
  padding: 12px 24px;
  background: var(--surface-glass);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-bottom: 1px solid var(--border-subtle);
  color: var(--text-main);
  display: flex;
  justify-content: space-between;
  align-items: center;
  position: sticky;
  top: 0;
  z-index: 1000;
  border-radius: 0; /* Header spans full width */
}

.header-left, .header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand-link {
  text-decoration: none;
  transition: opacity 0.2s;
}
.brand-link:hover { opacity: 0.8; }

.brand {
  font-size: 18px;
  font-weight: 800;
  letter-spacing: -0.5px;
  margin: 0;
  color: var(--brand-primary);
}

.divider {
  width: 1px;
  height: 20px;
  background: var(--border-subtle);
  margin: 0 8px;
}

.board-title {
  font-size: 16px;
  font-weight: 700;
  margin: 0;
  margin-right: 12px;
  color: var(--text-main);
}

.btn-action-premium {
  background: var(--brand-primary-light);
  border: 1px solid transparent;
  color: var(--brand-primary);
  padding: 6px 14px;
  border-radius: 10px;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  display: flex;
  align-items: center;
  gap: 6px;
}
.btn-action-premium:hover {
  background: var(--brand-primary);
  color: white;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(99, 102, 241, 0.2);
}

.search-wrapper {
  position: relative;
  margin-right: 8px;
}

.search-input-container {
  position: relative;
  display: flex;
  align-items: center;
}

.search-icon-inside {
  position: absolute;
  left: 14px;
  color: var(--text-muted);
  pointer-events: none;
}

.btn-nav-icon {
  background: var(--surface-secondary);
  width: 38px;
  height: 38px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
  color: var(--text-muted);
  border: 1px solid var(--border-subtle);
  text-decoration: none;
}
.btn-nav-icon:hover {
  background: var(--brand-primary-light);
  color: var(--brand-primary);
  transform: translateY(-2px);
  border-color: var(--brand-primary-light);
}

.btn-logout-hover:hover {
  background: #fee2e2;
  color: #ef4444;
  border-color: #fee2e2;
}

.search-input {
  background: var(--surface-secondary);
  border: 1px solid transparent;
  border-radius: 12px;
  padding: 8px 16px 8px 38px;
  color: var(--text-main);
  width: 240px;
  font-size: 14px;
  transition: all 0.2s ease;
}
.search-input::placeholder { color: var(--text-extramuted); }
.search-input:focus {
  background: var(--surface-primary);
  border-color: var(--brand-primary);
  outline: none;
  box-shadow: 0 0 0 3px var(--brand-primary-light);
  width: 280px;
}

.search-results-dropdown {
  position: absolute;
  top: calc(100% + 12px);
  right: 0;
  width: 340px;
  background: var(--surface-primary);
  border-radius: 16px;
  overflow: hidden;
  z-index: 1100;
  border: 1px solid var(--border-subtle);
}

.search-status {
  padding: 20px;
  text-align: center;
  color: var(--text-muted);
  font-size: 14px;
}

.search-result-item {
  padding: 14px 20px;
  border-bottom: 1px solid var(--surface-secondary);
  cursor: pointer;
  transition: all 0.2s;
}
.search-result-item:last-child { border-bottom: none; }
.search-result-item:hover { background: var(--surface-secondary); }

.search-card-title { font-weight: 700; font-size: 14px; margin-bottom: 4px; color: var(--text-main); }
.search-card-desc { font-size: 12px; color: var(--text-muted); }
</style>
