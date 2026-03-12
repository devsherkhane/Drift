<template>
  <div class="list-wrapper">
    <div class="list-content">
      <div class="list-header">
        <input 
          :value="list.title" 
          @change="$emit('update-title', list.id, $event.target.value)"
          class="list-title-input" 
        />
        <button @click="$emit('delete-list', list.id)" class="btn-delete-list">🗑️</button>
      </div>

      <draggable 
        :modelValue="list.cards" 
        @update:modelValue="$emit('update-cards', $event)"
        group="cards" 
        item-key="id" 
        class="cards-list" 
        ghost-class="ghost-card"
        @change="$emit('card-moved', $event, list.id)"
      >
        <template #item="{ element }">
          <div class="card" @click="$emit('open-card', element, list.title)">
            <div v-if="element.label_color" class="card-label-bar"
              :style="{ backgroundColor: element.label_color }"></div>
            {{ element.title }}
          </div>
        </template>
      </draggable>

      <div v-if="isEditing" class="add-card-form">
        <textarea 
          v-model="localNewCardTitle" 
          placeholder="Enter a title..."
          @keyup.enter="onSubmitCard"
        ></textarea>
        <div class="form-actions">
          <button class="btn-add" @click="onSubmitCard">Add card</button>
          <button class="btn-close" @click="$emit('cancel-add')">✕</button>
        </div>
      </div>
      <button v-else class="add-card-btn" @click="$emit('start-add')">+ Add a card</button>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue';
import draggable from 'vuedraggable';

const props = defineProps({
  list: { type: Object, required: true },
  isEditing: { type: Boolean, default: false }
});

const emit = defineEmits([
  'update-title', 'delete-list', 'update-cards',
  'card-moved', 'open-card', 'start-add', 'cancel-add', 'submit-card'
]);

const localNewCardTitle = ref('');

// Reset title when editing starts
watch(() => props.isEditing, (val) => {
  if (val) localNewCardTitle.value = '';
});

const onSubmitCard = () => {
  if (!localNewCardTitle.value.trim()) return;
  emit('submit-card', props.list.id, localNewCardTitle.value.trim());
  localNewCardTitle.value = '';
};
</script>

<style scoped>
.list-wrapper {
  width: 300px;
  flex-shrink: 0;
}

.list-content {
  background: var(--trello-list);
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid var(--border-subtle);
  border-radius: var(--border-radius);
  padding: 12px;
  display: flex;
  flex-direction: column;
  max-height: calc(100vh - 120px);
  box-shadow: var(--shadow-soft);
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.list-title-input {
  background: transparent;
  border: none;
  font-weight: bold;
  width: 80%;
  font-size: 16px;
  color: #172b4d;
}

.btn-delete-list {
  background: transparent;
  border: none;
  cursor: pointer;
  font-size: 14px;
}

.card {
  background: var(--trello-card);
  padding: 12px;
  margin-bottom: 8px;
  border-radius: var(--border-radius-sm);
  box-shadow: var(--card-shadow);
  cursor: grab;
  font-size: 14px;
  position: relative;
  transition: transform 0.2s cubic-bezier(0.175, 0.885, 0.32, 1.275), box-shadow 0.2s ease;
  border: 1px solid var(--border-subtle);
  color: var(--text-main);
}
.card:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
  border-color: var(--brand-primary);
}
.card:active { cursor: grabbing; }

.card-label-bar {
  height: 6px;
  width: 48px;
  border-radius: 3px;
  margin-bottom: 8px;
}

.ghost-card {
  opacity: 0.5;
}

.btn-add {
  background: var(--brand-primary);
  color: var(--text-on-brand);
  border: none;
  padding: 8px 16px;
  border-radius: var(--border-radius-sm);
  font-weight: 600;
  cursor: pointer;
  transition: background 0.2s;
}
.btn-add:hover { background: var(--brand-primary-hover); }

.btn-close {
  background: none;
  border: none;
  font-size: 20px;
  cursor: pointer;
  color: var(--text-secondary);
  padding: 4px 8px;
  border-radius: var(--border-radius-sm);
}
.btn-close:hover {
  background: var(--surface-secondary);
  color: var(--text-main);
}

.add-card-btn {
  background: transparent;
  border: none;
  padding: 10px 12px;
  border-radius: var(--border-radius-sm);
  color: var(--text-secondary);
  font-weight: 500;
  cursor: pointer;
  text-align: left;
  transition: all 0.2s ease;
  width: 100%;
}
.add-card-btn:hover {
  background: var(--surface-secondary);
  color: var(--text-main);
}

.add-card-form textarea {
  width: 100%;
  border: 1px solid rgba(99, 102, 241, 0.2);
  border-radius: var(--border-radius-sm);
  padding: 12px;
  margin-bottom: 8px;
  box-sizing: border-box;
  font-family: inherit;
  font-size: 14px;
  resize: vertical;
  min-height: 40px;
}
.add-card-form textarea:focus {
  outline: none;
  border-color: var(--brand-primary);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15);
}

.form-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>
