import AuthService from '../../services/auth.js';
<template>
  <div class="modal-backdrop" style="position:fixed;top:0;left:0;width:100vw;height:100vh;background:rgba(0,0,0,0.3);z-index:2100;display:flex;align-items:center;justify-content:center;">
    <div class="modal-dialog" style="background:white;border-radius:8px;box-shadow:0 2px 16px rgba(0,0,0,0.15);max-width:320px;width:100%;">
      <div class="modal-content" style="padding:1.25rem;">
        <div class="modal-header" style="display:flex;align-items:center;justify-content:space-between;">
          <h5 class="modal-title">Reactions</h5>
          <button type="button" class="close" @click="$emit('close')" style="background:none;border:none;font-size:1.5rem;">&times;</button>
        </div>
        <div class="modal-body">
          <ul class="reaction-list" style="padding:0;list-style:none;">
            <li v-for="reaction in reactions" :key="reaction.id" class="reaction-list-item" style="display:flex;align-items:center;gap:0.75rem;padding:0.5rem 0;">
              <span class="emoji" style="font-size:1.5rem;">{{ reaction.emoticon }}</span>
              <span class="username" style="font-size:1rem;color:#212529;">{{ reaction.username }}</span>
              <button v-if="isOwnReaction(reaction)" class="btn btn-sm btn-outline-danger" style="margin-left:auto;" @click="$emit('remove', reaction)">
                Remove
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
  name: 'ReactionListModal',
  props: {
    reactions: {
      type: Array,
      required: true
    }
  },
  methods: {
    isOwnReaction(reaction) {
      // Use currentUserId from prop or computed only
      return reaction.userId === this.currentUserId;
    }
  },
  computed: {
    currentUserId() {
      // Always use AuthService for user ID
      return AuthService.getUserId();
    }
  }
}
</script>

<style scoped>
/* Add custom styles if needed */
</style>
