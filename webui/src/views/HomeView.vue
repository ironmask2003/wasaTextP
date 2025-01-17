<script setup>
import Modal from '../components/Modal.vue';
</script>

<script>
export default {
  data: function () {
    return {
      errormsg: null,
      some_data: [],
      searchModalIsVisible: false,
    }
  },
  methods: {
    async getConversations() {
      this.errormsg = null;
      try {
        let response = await this.$axios.get(`/profiles/${sessionStorage.userID}/conversations`, { headers: { 'Authorization': sessionStorage.token } });
        this.some_data = response.data;
      } catch (e) {
        this.errormsg = e.toString();
      }
    },
    goToConversation(response) {
      localStorage.clear();
      localStorage.userID = response.user.userId;
      localStorage.username = response.user.username;
      localStorage.photo = response.user.photo;
      // Potrei fare un getConversation per prendere i dati dell'utente
      this.$router.push(`/conversation/${response.conversation.conversationId}`);
    },
    handleSearchModalToggle() {
      this.searchModalIsVisible = !this.searchModalIsVisible;
    },
  },
  mounted() {
    if (!sessionStorage.token) {
      this.$router.push("/");
      return;
    }
    this.getConversations();
  }
}
</script>

<template>
  <div>
    <div
      class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">Home page</h1>
      <Modal :show="searchModalIsVisible" @close="handleSearchModalToggle" title="search">
        <template v-slot:header>
          <h3>Users</h3>
        </template>
      </Modal>
      <div class="btn-toolbar mb-2 mb-md-0">
        <div class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="getConversations">
            Refresh
          </button>
        </div>
        <div class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-primary" @click="exportList">
            New Group
          </button>
        </div>
        <div class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-primary" @click="handleSearchModalToggle">
            New Chat
          </button>
        </div>
      </div>
    </div>

    <div class="conversations" v-for="response in some_data" :key="response.conversation.conversationId">
      <button v-if="response.message.photo == ''" type="button" class="btn btn-sm btn-outline-primary"
        @click="goToConversation(response)">
        {{ response.user.username }} <br> {{ response.senderUser.username }}: {{ response.message.text }}
      </button>
      <button type="button" class="btn btn-sm btn-outline-primary" @click="goToConversation(response)" v-else>
        {{ response.user.username }} <br> {{ response.senderUser.username }}: Photo
      </button>
      <hr>
    </div>

    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
  </div>
</template>

<style></style>
