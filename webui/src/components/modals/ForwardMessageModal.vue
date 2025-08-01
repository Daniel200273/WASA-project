<template>
  <div class="modal-backdrop" style="position:fixed;top:0;left:0;width:100vw;height:100vh;background:rgba(0,0,0,0.3);z-index:2000;display:flex;align-items:center;justify-content:center;">
    <div class="modal-dialog" style="background:white;border-radius:8px;box-shadow:0 2px 16px rgba(0,0,0,0.15);max-width:400px;width:100%;">
      <div class="modal-content" style="padding:1.5rem;">
        <div class="modal-header" style="display:flex;align-items:center;justify-content:space-between;">
          <h5 class="modal-title">Forward Message</h5>
          <button type="button" class="close" style="background:none;border:none;font-size:1.5rem;" @click="$emit('cancel')">&times;</button>
        </div>
        <div class="modal-body">
          <p>Select a conversation to forward this message to:</p>
          <ul class="list-group" style="padding:0;list-style:none;">
            <li v-for="conv in conversations" :key="conv.id" class="list-group-item" style="margin-bottom:0.5rem;display:flex;align-items:center;justify-content:space-between;padding:0.5rem;border:1px solid #dee2e6;border-radius:4px;">
              <span style="font-weight:500;">{{ conv.name }}</span>
              <button class="btn btn-sm btn-primary" style="padding:0.25rem 0.75rem;" @click="$emit('confirm', conv.id)">
                Forward
              </button>
            </li>
          </ul>
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
  },
  emits: ['cancel', 'confirm']
}
</script>

<style scoped>
/* Add any custom styles if needed */
</style>
