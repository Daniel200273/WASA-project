<template>
  <div class="modal-overlay" @click.self="close">
    <div class="notification-modal" :class="type">
      <div class="modal-header">
        <div class="icon-container">
          <svg v-if="type === 'success'" class="feather icon-success">
            <use href="/feather-sprite-v4.29.0.svg#check-circle" />
          </svg>
          <svg v-else-if="type === 'error'" class="feather icon-error">
            <use href="/feather-sprite-v4.29.0.svg#x-circle" />
          </svg>
          <svg v-else class="feather icon-info">
            <use href="/feather-sprite-v4.29.0.svg#info" />
          </svg>
        </div>
        <h3 class="modal-title">{{ title }}</h3>
      </div>
      
      <div class="modal-body">
        <p>{{ message }}</p>
      </div>
      
      <div class="modal-footer">
        <button class="btn btn-primary" @click="close">
          OK
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'NotificationModal',
  props: {
    type: {
      type: String,
      default: 'info', // 'success', 'error', 'info'
      validator: value => ['success', 'error', 'info'].includes(value)
    },
    title: {
      type: String,
      required: true
    },
    message: {
      type: String,
      required: true
    }
  },
  emits: ['close'],
  mounted() {
    // Auto-close after 3 seconds for success messages
    if (this.type === 'success') {
      setTimeout(() => {
        this.close();
      }, 3000);
    }
    
    // Add keyboard listener for ESC key
    document.addEventListener('keydown', this.handleKeydown);
  },
  beforeUnmount() {
    document.removeEventListener('keydown', this.handleKeydown);
  },
  methods: {
    close() {
      this.$emit('close');
    },
    handleKeydown(event) {
      if (event.key === 'Escape') {
        this.close();
      }
    }
  }
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2000;
}

.notification-modal {
  background: white;
  border-radius: 8px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.2);
  max-width: 400px;
  width: 90%;
  max-height: 90vh;
  overflow: hidden;
}

.modal-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.5rem 1.5rem 1rem 1.5rem;
  border-bottom: 1px solid #e9ecef;
}

.icon-container {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border-radius: 50%;
  flex-shrink: 0;
}

.notification-modal.success .icon-container {
  background-color: #d4edda;
  color: #155724;
}

.notification-modal.error .icon-container {
  background-color: #f8d7da;
  color: #721c24;
}

.notification-modal.info .icon-container {
  background-color: #d1ecf1;
  color: #0c5460;
}

.icon-success,
.icon-error,
.icon-info {
  width: 24px;
  height: 24px;
}

.modal-title {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: #212529;
}

.notification-modal.success .modal-title {
  color: #155724;
}

.notification-modal.error .modal-title {
  color: #721c24;
}

.notification-modal.info .modal-title {
  color: #0c5460;
}

.modal-body {
  padding: 1rem 1.5rem;
}

.modal-body p {
  margin: 0;
  color: #6c757d;
  line-height: 1.5;
}

.modal-footer {
  padding: 1rem 1.5rem 1.5rem 1.5rem;
  display: flex;
  justify-content: flex-end;
}

.btn {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1.5rem;
  border: 1px solid transparent;
  border-radius: 4px;
  font-size: 0.875rem;
  font-weight: 500;
  text-decoration: none;
  cursor: pointer;
  transition: all 0.15s ease;
}

.btn-primary {
  background-color: #007bff;
  border-color: #007bff;
  color: white;
}

.btn-primary:hover {
  background-color: #0056b3;
  border-color: #0056b3;
}

.btn-primary:focus {
  outline: none;
  box-shadow: 0 0 0 3px rgba(0, 123, 255, 0.25);
}

/* Responsive adjustments */
@media (max-width: 480px) {
  .notification-modal {
    margin: 1rem;
    width: calc(100% - 2rem);
  }
  
  .modal-header,
  .modal-body,
  .modal-footer {
    padding-left: 1rem;
    padding-right: 1rem;
  }
}
</style>
