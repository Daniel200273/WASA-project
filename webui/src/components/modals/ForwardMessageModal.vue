<template>
  <div class="modal-backdrop" style="position:fixed;top:0;left:0;width:100vw;height:100vh;background:rgba(0,0,0,0.3);z-index:2000;display:flex;align-items:center;justify-content:center;">
    <div class="modal-dialog" style="background:white;border-radius:8px;box-shadow:0 2px 16px rgba(0,0,0,0.15);max-width:400px;width:100%;">
      <div class="modal-content" style="padding:1.5rem;">
        <div class="modal-header" style="display:flex;align-items:center;justify-content:space-between;">
          <h5 class="modal-title">Forward Message</h5>
          <button type="button" class="close" @click="$emit('cancel')" style="background:none;border:none;font-size:1.5rem;">&times;</button>
        </div>
        <div class="modal-body">
          <p>Select a conversation to forward this message to:</p>
          <ul class="list-group" style="padding:0;list-style:none;">
            <li v-for="conv in conversations" :key="conv.id" class="list-group-item" style="margin-bottom:0.5rem;">
              <button class="btn btn-sm btn-outline-primary" @click="$emit('confirm', conv.id)">
                {{ conv.name }}
              </button>
            </li>
          </ul>
        </div>
        <div class="modal-footer" style="text-align:right;">
          <button class="btn btn-secondary" @click="$emit('cancel')">Cancel</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'ForwardMessageModal',
  props: {
    conversations: {
      type: Array,
      required: true
    }
  }
}
</script>

<style scoped>
/* Add any custom styles if needed */
</style>
