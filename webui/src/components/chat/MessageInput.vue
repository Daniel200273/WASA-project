<template>
  <div class="message-input">
    <div class="input-container">
      <!-- Photo upload button -->
      <button class="input-action-btn" :disabled="disabled" title="Attach photo" @click="selectPhoto">
        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#image" /></svg>
      </button>

      <!-- Text input (hidden when photo is selected) -->
      <div v-if="!selectedPhoto" class="text-input-wrapper">
        <textarea
          ref="textInput"
          v-model="messageText"
          :placeholder="placeholder"
          :disabled="disabled"
          rows="1"
          class="text-input"
          @keydown="handleKeyDown"
          @input="adjustHeight"
        />
      </div>

      <!-- Send button -->
      <button 
        class="send-btn" 
        :disabled="disabled || (!messageText.trim() && !selectedPhoto)" 
        title="Send message"
        @click="sendMessage"
      >
        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#send" /></svg>
      </button>
    </div>

    <!-- Photo preview -->
    <div v-if="selectedPhoto" class="photo-preview">
      <div class="photo-preview-container">
        <img :src="photoPreviewUrl" alt="Photo to send">
        <button class="remove-photo" title="Remove photo" @click="removePhoto">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x" /></svg>
        </button>
      </div>
      <div class="photo-caption">
        <span>Photo ready to send</span>
      </div>
    </div>

    <!-- Hidden file input -->
    <input
      ref="photoInput"
      type="file"
      accept="image/*"
      style="display: none;"
      @change="handlePhotoSelect"
    >
  </div>
</template>

<script>
export default {
  name: 'MessageInput',
  props: {
    placeholder: {
      type: String,
      default: 'Type a message...'
    },
    disabled: {
      type: Boolean,
      default: false
    },
    replyingTo: {
      type: Object,
      default: null
    }
  },
  emits: ['send-message', 'send-photo'],
  data() {
    return {
      messageText: '',
      selectedPhoto: null,
      photoPreviewUrl: null
    }
  },
  watch: {
    replyingTo() {
      // Focus input when starting a reply
      if (this.replyingTo) {
        this.$nextTick(() => {
          this.$refs.textInput.focus();
        });
      }
    }
  },
  mounted() {
    this.adjustHeight();
  },
  beforeUnmount() {
    // Clean up photo preview URL
    if (this.photoPreviewUrl) {
      URL.revokeObjectURL(this.photoPreviewUrl);
    }
  },
  methods: {
    handleKeyDown(event) {
      // Send message on Enter (but not Shift+Enter)
      if (event.key === 'Enter' && !event.shiftKey) {
        event.preventDefault();
        this.sendMessage();
      }
    },

    sendMessage() {
      if (this.disabled) return;

      if (this.selectedPhoto) {
        // Send photo message only (no text caption)
        this.$emit('send-photo', this.selectedPhoto, null);
        this.resetInput();
      } else if (this.messageText.trim()) {
        // Send text message
        this.$emit('send-message', this.messageText.trim());
        this.resetInput();
      }
    },

    selectPhoto() {
      if (this.disabled) return;
      this.$refs.photoInput.click();
    },

    handlePhotoSelect(event) {
      const file = event.target.files[0];
      if (!file) return;

      // Validate file type
      if (!file.type.startsWith('image/')) {
        alert('Please select an image file');
        return;
      }

      // Validate file size (10MB limit)
      const maxSize = 10 * 1024 * 1024;
      if (file.size > maxSize) {
        alert('Image must be smaller than 10MB');
        return;
      }

      this.selectedPhoto = file;
      this.photoPreviewUrl = URL.createObjectURL(file);

      // Clear message text since it won't be visible/editable with photo
      this.messageText = '';

      // Clear the input so the same file can be selected again
      event.target.value = '';
    },

    removePhoto() {
      if (this.photoPreviewUrl) {
        URL.revokeObjectURL(this.photoPreviewUrl);
      }
      this.selectedPhoto = null;
      this.photoPreviewUrl = null;
    },

    resetInput() {
      this.messageText = '';
      this.removePhoto();
      this.adjustHeight();
      
      // Reset file input
      if (this.$refs.photoInput) {
        this.$refs.photoInput.value = '';
      }
    },

    adjustHeight() {
      this.$nextTick(() => {
        const textarea = this.$refs.textInput;
        if (textarea) {
          textarea.style.height = 'auto';
          textarea.style.height = Math.min(textarea.scrollHeight, 120) + 'px';
        }
      });
    }
  }
}
</script>

<style scoped>
.message-input {
  padding: 1rem;
}

.input-container {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  background-color: #f8f9fa;
  border-radius: 24px;
  padding: 0.75rem 1rem;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
  border: 1px solid #e9ecef;
}

/* Input action buttons */
.input-action-btn {
  background-color: #ffffff;
  border: 1px solid #e9ecef;
  padding: 0.6rem;
  border-radius: 50%;
  cursor: pointer;
  color: #6c757d;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  flex-shrink: 0;
}

.input-action-btn:hover:not(:disabled) {
  background-color: #f8f9fa;
  color: #495057;
  border-color: #dee2e6;
}

.input-action-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Text input wrapper */
.text-input-wrapper {
  flex: 1;
  min-height: 40px;
  display: flex;
  align-items: center;
}

.text-input {
  width: 100%;
  border: none;
  background: none;
  resize: none;
  outline: none;
  font-family: inherit;
  font-size: 0.9rem;
  line-height: 1.4;
  padding: 0.5rem;
  max-height: 120px;
  overflow-y: auto;
}

.text-input::placeholder {
  color: #6c757d;
}

.text-input:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Send button */
.send-btn {
  background-color: #007bff;
  border: none;
  padding: 0.5rem;
  border-radius: 50%;
  cursor: pointer;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s ease;
  flex-shrink: 0;
  width: 42px;
  height: 42px;
  margin-left: auto;
}

.send-btn:hover:not(:disabled) {
  background-color: #0056b3;
}

.send-btn:disabled {
  background-color: #bdbdbd;
  cursor: not-allowed;
  opacity: 0.6;
}

.send-btn .feather {
  width: 18px;
  height: 18px;
}

/* Photo preview */
.photo-preview {
  margin-top: 1rem;
  border-radius: 16px;
  background-color: #ffffff;
  padding: 1rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  border: 1px solid #e9ecef;
}

.photo-preview-container {
  position: relative;
  max-width: 240px;
  margin: 0 auto;
  border-radius: 12px;
  overflow: hidden;
  background-color: white;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
}

.photo-preview-container img {
  width: 100%;
  height: auto;
  max-height: 200px;
  object-fit: cover;
  display: block;
}

.remove-photo {
  position: absolute;
  top: 0.75rem;
  right: 0.75rem;
  background: rgba(255, 255, 255, 0.9);
  border: none;
  border-radius: 50%;
  width: 30px;
  height: 30px;
  cursor: pointer;
  color: #d32f2f;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s ease;
  box-shadow: 0 1px 4px rgba(0, 0, 0, 0.2);
}

.remove-photo:hover {
  background: rgba(255, 255, 255, 1);
}

.remove-photo .feather {
  width: 16px;
  height: 16px;
}

.photo-caption {
  padding: 0.75rem 0 0.25rem 0;
  text-align: center;
  font-size: 0.85rem;
  color: #6c757d;
  font-weight: 500;
}

/* Responsive design */
@media (max-width: 768px) {
  .message-input {
    padding: 0.75rem;
  }

  .input-container {
    padding: 0.375rem;
  }

  .text-input {
    font-size: 16px; /* Prevent zoom on iOS */
  }

  .send-btn,
  .input-action-btn {
    width: 36px;
    height: 36px;
    padding: 0.375rem;
  }

  .send-btn .feather {
    width: 16px;
    height: 16px;
  }
}

/* Common icon size */
.feather {
  width: 16px;
  height: 16px;
}
</style>
