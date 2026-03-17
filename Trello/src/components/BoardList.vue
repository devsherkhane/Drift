<template>
  <div class="list-wrapper">
    <div class="list-content glass-panel">
      <div class="list-header">
        <input 
          :value="list.title" 
          @change="$emit('update-title', list.id, $event.target.value)"
          class="list-title-input" 
          placeholder="List title..."
        />
        <button @click="$emit('delete-list', list.id)" class="btn-icon-delete" title="Delete List">
          <Trash2Icon :size="14" />
        </button>
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
          <div class="card shadow-sm" @click="$emit('open-card', element, list.title)">
            <div v-if="element.label_color" class="card-label-bar"
              :style="{ backgroundColor: element.label_color }"></div>
            <div class="card-title-text">{{ element.title }}</div>
            <div class="card-footer">
              <ChevronRightIcon :size="12" class="card-arrow" />
            </div>
          </div>
        </template>
      </draggable>

      <div v-if="isEditing" class="add-card-form animate-fade-in">
        <textarea 
          v-model="localNewCardTitle" 
          placeholder="Enter a title for this card..."
          @keyup.enter="onSubmitCard"
          autofocus
        ></textarea>
        <div class="form-actions">
          <button class="btn-add-card" @click="onSubmitCard">Add Card</button>
          <button class="btn-close-form" @click="$emit('cancel-add')">Cancel</button>
        </div>
      </div>
      <button v-else class="add-card-btn" @click="$emit('start-add')">
        <PlusIcon :size="16" /> Add a card
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue';
import draggable from 'vuedraggable';
import { Trash2 as Trash2Icon, Plus as PlusIcon, ChevronRight as ChevronRightIcon } from 'lucide-vue-next';

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
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border: 1px solid var(--border-subtle);
  border-radius: 16px;
  padding: 16px;
  display: flex;
  flex-direction: column;
  max-height: calc(100vh - 120px);
  box-shadow: var(--shadow-md);
  transition: box-shadow 0.3s ease;
}

.list-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding: 0 4px;
}

.list-title-input {
  background: transparent;
  border: none;
  font-weight: 800;
  width: 85%;
  font-size: 15px;
  color: var(--text-main);
  padding: 4px;
  border-radius: 6px;
  transition: background 0.2s;
}
.list-title-input:focus {
  background: var(--surface-primary);
  outline: none;
  box-shadow: 0 0 0 2px var(--brand-primary-light);
}

.btn-icon-delete {
  background: transparent;
  border: none;
  cursor: pointer;
  color: var(--text-extramuted);
  padding: 6px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}
.btn-icon-delete:hover {
  background: #fee2e2;
  color: #ef4444;
}

.cards-list {
  min-height: 10px;
  display: flex;
  flex-direction: column;
}

.card {
  background: var(--trello-card);
  padding: 14px;
  margin-bottom: 12px;
  border-radius: 12px;
  cursor: grab;
  font-size: 14px;
  position: relative;
  transition: all 0.25s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  border: 1px solid var(--border-subtle);
  color: var(--text-main);
  display: flex;
  flex-direction: column;
}
.card:hover {
  transform: translateY(-4px) scale(1.02);
  box-shadow: var(--shadow-md);
  border-color: var(--brand-primary);
}
.card:active { cursor: grabbing; }

.card-title-text {
  font-weight: 600;
  line-height: 1.4;
}

.card-label-bar {
  height: 5px;
  width: 32px;
  border-radius: 10px;
  margin-bottom: 10px;
}

.card-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: 8px;
  opacity: 0.4;
}

.card-arrow { transition: transform 0.2s; }
.card:hover .card-arrow { transform: translateX(2px); }

.ghost-card {
  opacity: 0.3;
  background: var(--brand-primary-light);
  border: 2px dashed var(--brand-primary);
}

.btn-add-card {
  background: var(--brand-primary);
  color: white;
  border: none;
  padding: 10px 16px;
  border-radius: 10px;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.2s;
  box-shadow: 0 4px 10px rgba(99, 102, 241, 0.2);
}
.btn-add-card:hover { 
  background: var(--brand-primary-hover); 
  transform: translateY(-1px);
}

.btn-close-form {
  background: transparent;
  border: none;
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  color: var(--text-muted);
  padding: 8px 12px;
}
.btn-close-form:hover {
  color: var(--text-main);
}

.add-card-btn {
  background: transparent;
  border: none;
  padding: 12px;
  border-radius: 10px;
  color: var(--text-muted);
  font-weight: 700;
  cursor: pointer;
  text-align: left;
  transition: all 0.2s;
  width: 100%;
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 4px;
}
.add-card-btn:hover {
  background: var(--surface-secondary);
  color: var(--brand-primary);
  padding-left: 16px;
}

.add-card-form textarea {
  width: 100%;
  border: 2px solid var(--surface-secondary);
  border-radius: 12px;
  padding: 14px;
  margin-bottom: 12px;
  box-sizing: border-box;
  font-family: inherit;
  font-size: 14px;
  resize: vertical;
  min-height: 80px;
  transition: all 0.2s;
}
.add-card-form textarea:focus {
  outline: none;
  border-color: var(--brand-primary);
  background: var(--surface-primary);
  box-shadow: 0 0 0 3px var(--brand-primary-light);
}

.form-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.animate-fade-in {
  animation: fadeIn 0.3s ease-out;
}
@keyframes fadeIn {
  from { opacity: 0; transform: translateY(5px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>
