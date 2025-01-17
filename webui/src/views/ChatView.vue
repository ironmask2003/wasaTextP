<script setup>
import Modal from '../components/ModalConv.vue';
</script>

<script>
export default {
  data() {
    return {
      errormsg: null,
      userToSend: localStorage.username,
      userIdToSend: localStorage.userID,
      convId: parseInt(this.$route.params.convId),
      text: null,
      photo: null,
      proPic64: localStorage.photo,
      some_data: [],
      searchModalIsVisible: false,
      messageToFordward: null,
    }
  },
  methods: {
    async handleFileChange(event) {
      this.errorMsg = "";
      const file = event.target.files[0];
      if (file.type !== "image/jpeg") {
        this.errorMsg = "File type not supported, only jpg and jpeg are allowed";
        return
      }
      if (file.size > 5242880) {
        this.errorMsg = "File size is too big. Max size is 5MB";
        return
      }
      this.photo = file;
    },
    async getConversation() {
      this.errormsg = null;
      this.$axios.get(`/profiles/${sessionStorage.userID}/conversations/${this.convId}`, { headers: { 'Authorization': sessionStorage.token } })
        .then(response => {
          this.some_data = response.data;
        })
        .catch(e => {
          this.errormsg = e.toString();
        });
    },
    check() {
      if (isNaN(this.convId) || this.convId == undefined) {
        this.createConversation();
      } else {
        this.sendMessage();
      }
    },
    createConversation() {
      this.errormsg = null;
      this.$axios.put(`/profiles/${sessionStorage.userID}/conversations/${this.userIdToSend}`, {
        text: this.text
      }, { headers: { 'Authorization': sessionStorage.token } })
        .then(response => {
          this.convId = response.data.conversationId;
          this.getConversation();
        })
        .catch(e => {
          this.errormsg = e.toString();
        });

    },
    async sendMessage() {
      this.errormsg = null;
      const formData = new FormData();
      formData.append('text', this.text);
      if (this.photo != null) {
        formData.append('photo', this.photo);
      }
      this.$axios.post(`/profiles/${sessionStorage.userID}/conversations/${this.convId}/messages`, formData, { headers: { 'Authorization': sessionStorage.token } })
        .then(() => {
          this.text = null;
          this.photo = null;
          this.getConversation();
        })
        .catch(e => {
          this.errormsg = e.toString();
        });
    },
    handleSearchModalToggle(msg) {
      this.messageToFordward = msg;
      this.searchModalIsVisible = !this.searchModalIsVisible;
    },
  },
  mounted() {
    if (!sessionStorage.token) {
      this.$router.push("/");
      return;
    }
    if (this.convId != undefined || !isNaN(this.convId)) {
      this.getConversation()
    }
  }
}
</script>

<template>
  <div>
    <div
      class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <div class="top-profile-container">
        <img :src="`data:image/jpg;base64,${proPic64}`">
      </div>
      <h1 class="h1">{{ this.userToSend }}</h1>
      <Modal :show="searchModalIsVisible" :msg="messageToFordward" @close="handleSearchModalToggle" title="search">
        <template v-slot:header>
          <h3>Conversations</h3>
        </template>
      </Modal>
      <div class="btn-toolbar mb-2 mb-md-0">
        <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
        <div class="btn-group me-2">
          <form @submit.prevent="sendMessage">
            <input type="file" ref="file" accept=".jpg,.jpeg" @change="handleFileChange" />
            <button type="submit" class="btn btn-sm btn-outline-primary">
              Send photo
            </button>
          </form>
        </div>
      </div>
    </div>
    <div class="messages" v-for="response in some_data" :key="response.message.messageId">
      <p v-if="response.message.text !== 'null' || response.message.photo !== ''">
        {{ response.user.username }}
      </p>
      <p v-if="response.message.text !== 'null'">
        {{ response.message.text }}
      </p>
      <p></p>
      <img class="msg_photo" v-if="response.message.photo !== ''"
        :src="`data:image/jpg;base64,${response.message.photo}`" alt="Message Photo">
      <p v-if="response.message.text !== 'null' || response.message.photo !== ''">
        {{ response.timeMsg }}
      </p>
      <button type="button" class="btn btn-sm btn-outline-secondary" @click="handleSearchModalToggle(response.message)">
        Forward Message
      </button>
      <button type="button" class="btn btn-sm btn-outline-secondary">
        Delete message
      </button>
      <hr v-if="response.message.text !== 'null' || response.message.photo !== ''">
    </div>
    <div class="input-group">
      <input type="text" class="form-control" v-model="text" placeholder="Type your message here">
      <button class="btn btn-outline-primary" @click="check">Send</button>
    </div>

    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
  </div>
</template>

<style>
.top-profile-container {
  width: auto;
  height: auto;

  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: space-between;
}

.msg_photo {
  width: 25%;
  height: 25%;
}
</style>
