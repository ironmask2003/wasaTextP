<script setup>
import Modal from '../components/Modal.vue';
import Group from '../components/ModalGroup.vue';
</script>

<script>
export default {
  data: function () {
    return {
      errormsg: null,
      some_data: [],
      searchModalIsVisible: false,
      createGroupModalIsVisible: false,
      users: [],
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
      if(response.group.groupId !== 0) {
        localStorage.userID = response.group.groupId;
        localStorage.username = response.group.groupName;
        localStorage.photo = response.group.photo;
        localStorage.users = JSON.stringify(response.groupUsers);
      } else {
        localStorage.userID = response.user.userId;
        localStorage.username = response.user.username;
        localStorage.photo = response.user.photo;
      }
      // Potrei fare un getConversation per prendere i dati dell'utente
      this.$router.push(`/conversation/${response.conversation.conversationId}`);
    },
    handleSearchModalToggle() {
      this.searchModalIsVisible = !this.searchModalIsVisible;
    },
    handleCreateGroupModalToggle() {
      this.createGroupModalIsVisible = !this.createGroupModalIsVisible;
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
      <Group :show="createGroupModalIsVisible" @close="handleCreateGroupModalToggle" title="search">
        <template v-slot:header>
          <h3>Select users</h3>
        </template>
      </Group>
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
          <button type="button" class="btn btn-sm btn-outline-primary" @click="handleCreateGroupModalToggle">
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
      <div v-if="response.group.groupName == '' ">
        <button v-if="response.message.photo == ''" type="button" class="btn btn-sm btn-outline-primary"
          @click="goToConversation(response)">
          {{ response.user.username }} <br> {{ response.senderUser.username }}: {{ response.message.text }}
        </button>
        <button type="button" class="btn btn-sm btn-outline-primary" @click="goToConversation(response)" v-else>
          {{ response.user.username }} <br> {{ response.senderUser.username }}: Photo
        </button>
      </div>
      <div v-else>
        <button v-if="response.message.photo == ''" type="button" class="btn btn-sm btn-outline-primary"
          @click="goToConversation(response)">
          {{ response.group.groupName }} <br> {{ response.senderUser.username }}: {{ response.message.text }}
        </button>
        <button type="button" class="btn btn-sm btn-outline-primary" @click="goToConversation(response)" v-else>
          {{ response.group.groupName }} <br> {{ response.senderUser.username }}: Photo
        </button>
      </div>
      <hr>
    </div>

    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
  </div>
</template>

<style></style>
