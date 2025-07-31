<template>
  <div class="modal-overlay" @click.self="cancel">
    <div class="confirmation-modal">
      <div class="modal-header">
        <div class="icon-container">
          <svg class="feather icon-warning">
            <use href="/feather-sprite-v4.29.0.svg#alert-triangle" />
          </svg>
        </div>
        <h3 class="modal-title">{{ title }}</h3>
      </div>
      
      <div class="modal-body">
        <p>{{ message }}</p>
      </div>
      
      <div class="modal-footer">
        <button class="btn btn-secondary" @click="cancel">
          Cancel
        </button>
        <button class="btn btn-danger" @click="confirm">
          Confirm
        </button>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ConfirmationModal',
  props: {
    title: {
      type: String,
      default: 'Confirm Action'
    },
    message: {
      type: String,
      required: true
    }
  },
  emits: ['confirm', 'cancel'],
  mounted() {
    // Add keyboard listener for ESC key
    document.addEventListener('keydown', this.handleKeydown);
  },
  beforeUnmount() {
    document.removeEventListener('keydown', this.handleKeydown);
  },
  methods: {
    confirm() {
      this.$emit('confirm');
    },
    cancel() {
      this.$emit('cancel');
    },
    handleKeydown(event) {
      if (event.key === 'Escape') {
        this.cancel();
      } else if (event.key === 'Enter') {
        this.confirm();
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

.confirmation-modal {
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
  background-color: #fff3cd;
  color: #856404;
}

.icon-warning {
  width: 24px;
  height: 24px;
}

.modal-title {
  margin: 0;
  font-size: 1.25rem;
  font-weight: 600;
  color: #212529;
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
  gap: 0.75rem;
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

.btn-secondary {
  background-color: #6c757d;
  border-color: #6c757d;
  color: white;
}

.btn-secondary:hover {
  background-color: #545b62;
  border-color: #4e555b;
}

.btn-danger {
  background-color: #dc3545;
  border-color: #dc3545;
  color: white;
}

.btn-danger:hover {
  background-color: #c82333;
  border-color: #bd2130;
}

.btn:focus {
  outline: none;
  box-shadow: 0 0 0 3px rgba(0, 123, 255, 0.25);
}

.btn-secondary:focus {
  box-shadow: 0 0 0 3px rgba(108, 117, 125, 0.25);
}

.btn-danger:focus {
  box-shadow: 0 0 0 3px rgba(220, 53, 69, 0.25);
}

/* Responsive adjustments */
@media (max-width: 480px) {
  .confirmation-modal {
    margin: 1rem;
    width: calc(100% - 2rem);
  }
  
  .modal-header,
  .modal-body,
  .modal-footer {
    padding-left: 1rem;
    padding-right: 1rem;
  }
  
  .modal-footer {
    flex-direction: column;
  }
  
  .modal-footer .btn {
    width: 100%;
    justify-content: center;
  }
}
</style>
