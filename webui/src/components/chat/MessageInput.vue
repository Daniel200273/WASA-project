<template>
  <div class="message-input">
    <div class="input-container">
      <!-- Photo upload button -->
      <button class="input-action-btn" @click="selectPhoto" :disabled="disabled" title="Attach photo">
        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#image" /></svg>
      </button>

      <!-- Text input -->
      <div class="text-input-wrapper">
        <textarea
          ref="textInput"
          v-model="messageText"
          :placeholder="placeholder"
          :disabled="disabled"
          @keydown="handleKeyDown"
          @input="adjustHeight"
          rows="1"
          class="text-input"
        ></textarea>
      </div>

      <!-- Send button -->
      <button 
        class="send-btn" 
        @click="sendMessage" 
        :disabled="disabled || (!messageText.trim() && !selectedPhoto)"
        title="Send message"
      >
        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#send" /></svg>
      </button>
    </div>

    <!-- Photo preview -->
    <div v-if="selectedPhoto" class="photo-preview">
      <div class="photo-preview-container">
        <img :src="photoPreviewUrl" alt="Photo to send" />
        <button class="remove-photo" @click="removePhoto" title="Remove photo">
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
      @change="handlePhotoSelect"
      style="display: none;"
    />
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
        // Send photo message
        this.$emit('send-photo', this.selectedPhoto);
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
  },
  beforeUnmount() {
    // Clean up photo preview URL
    if (this.photoPreviewUrl) {
      URL.revokeObjectURL(this.photoPreviewUrl);
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
  align-items: flex-end;
  gap: 0.5rem;
  background-color: #f8f9fa;
  border-radius: 24px;
  padding: 0.5rem;
}

/* Input action buttons */
.input-action-btn {
  background: none;
  border: none;
  padding: 0.5rem;
  border-radius: 50%;
  cursor: pointer;
  color: #6c757d;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.15s ease;
  flex-shrink: 0;
}

.input-action-btn:hover:not(:disabled) {
  background-color: #e9ecef;
  color: #495057;
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
  transition: all 0.15s ease;
  flex-shrink: 0;
  width: 40px;
  height: 40px;
}

.send-btn:hover:not(:disabled) {
  background-color: #0056b3;
}

.send-btn:disabled {
  background-color: #6c757d;
  cursor: not-allowed;
  opacity: 0.5;
}

.send-btn .feather {
  width: 18px;
  height: 18px;
}

/* Photo preview */
.photo-preview {
  margin-top: 0.75rem;
  border-radius: 12px;
  background-color: #f8f9fa;
  overflow: hidden;
}

.photo-preview-container {
  position: relative;
  max-width: 200px;
}

.photo-preview-container img {
  width: 100%;
  height: auto;
  display: block;
  border-radius: 8px;
}

.remove-photo {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  background-color: rgba(0, 0, 0, 0.6);
  border: none;
  border-radius: 50%;
  width: 24px;
  height: 24px;
  cursor: pointer;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.15s ease;
}

.remove-photo:hover {
  background-color: rgba(0, 0, 0, 0.8);
}

.remove-photo .feather {
  width: 14px;
  height: 14px;
}

.photo-caption {
  padding: 0.5rem;
  font-size: 0.8rem;
  color: #6c757d;
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
