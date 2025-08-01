<template>
  <div class="modal-overlay" @click="$emit('cancel')">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h3>{{ title }}</h3>
        <button class="modal-close" @click="$emit('cancel')">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x" /></svg>
        </button>
      </div>
      
      <div class="modal-body">
        <p v-if="message">{{ message }}</p>
        <input 
          ref="input"
          v-model="inputValue"
          type="text"
          :placeholder="placeholder"
          class="modal-input"
          @keyup.enter="confirm"
          @keyup.escape="$emit('cancel')"
        >
      </div>
      
      <div class="modal-footer">
        <button class="modal-btn secondary" @click="$emit('cancel')">
          Cancel
        </button>
        <button class="modal-btn primary" :disabled="!inputValue.trim()" @click="confirm">
          Confirm
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'InputModal',
  props: {
    title: {
      type: String,
      required: true
    },
    message: {
      type: String,
      default: ''
    },
    placeholder: {
      type: String,
      default: ''
    },
    defaultValue: {
      type: String,
      default: ''
    }
  },
  emits: ['cancel', 'confirm'],
  data() {
    return {
      inputValue: this.defaultValue
    };
  },
  mounted() {
    // Focus the input when modal opens
    this.$nextTick(() => {
      if (this.$refs.input) {
        this.$refs.input.focus();
        // Select all text if there's a default value
        if (this.defaultValue) {
          this.$refs.input.select();
        }
      }
    });
  },
  methods: {
    confirm() {
      if (this.inputValue.trim()) {
        this.$emit('confirm', this.inputValue.trim());
      }
    }
  }
};
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 8px;
  min-width: 400px;
  max-width: 500px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px 24px 16px;
  border-bottom: 1px solid #e5e7eb;
}

.modal-header h3 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #111827;
}

.modal-close {
  background: none;
  border: none;
  padding: 4px;
  cursor: pointer;
  border-radius: 4px;
  color: #6b7280;
}

.modal-close:hover {
  background-color: #f3f4f6;
  color: #374151;
}

.modal-close .feather {
  width: 20px;
  height: 20px;
}

.modal-body {
  padding: 20px 24px;
}

.modal-body p {
  margin: 0 0 16px 0;
  color: #374151;
  line-height: 1.5;
}

.modal-input {
  width: 100%;
  padding: 12px 16px;
  border: 1px solid #d1d5db;
  border-radius: 6px;
  font-size: 14px;
  outline: none;
  transition: border-color 0.15s;
}

.modal-input:focus {
  border-color: #3b82f6;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.modal-footer {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  padding: 16px 24px 20px;
  border-top: 1px solid #e5e7eb;
}

.modal-btn {
  padding: 8px 16px;
  border-radius: 6px;
  border: none;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s;
}

.modal-btn.secondary {
  background-color: #f3f4f6;
  color: #374151;
}

.modal-btn.secondary:hover {
  background-color: #e5e7eb;
}

.modal-btn.primary {
  background-color: #3b82f6;
  color: white;
}

.modal-btn.primary:hover:not(:disabled) {
  background-color: #2563eb;
}

.modal-btn.primary:disabled {
  background-color: #9ca3af;
  cursor: not-allowed;
}
</style>
