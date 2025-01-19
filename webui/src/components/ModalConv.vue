<script>

export default {
  props: {
    show: Boolean,
    title: String,
    msg: Object,
  },
  data() {
    return {
      errorMsg: "",
      convs: [],
      convId: parseInt(this.$route.params.convId),
    };
  },
  methods: {
    closeModal(user) {
      // Cleare localStorae
      localStorage.clear();
      localStorage.userID = user.userId;
      localStorage.username = user.username;
      localStorage.photo = user.photo;
      this.$emit('close');
    },
    async forwardMessage(destConvId, user) {
      // Set the local storage
      this.errormsg = null;
      const url = `/profiles/${sessionStorage.userID}/conversations/${this.convId}/messages/${this.msg.messageId}?dest_conv=${destConvId}`;
      this.$axios.post(url, {}, { headers: { 'Authorization': sessionStorage.token } })
        .then(() => {
          this.closeModal(user);
        })
        .catch(e => {
          this.errormsg = e.toString();
        });
    },
    async getMyConversations() {
      this.errormsg = null;
      try {
        let response = await this.$axios.get(`/profiles/${sessionStorage.userID}/conversations`, { headers: { 'Authorization': sessionStorage.token } });
        this.convs = response.data;
      } catch (e) {
        this.errormsg = e.toString();
      }
    }
  },
  watch: {
    show() {
      this.getMyConversations();
    }
  },

}
</script>

<template>
  <Transition name="modal">
    <div v-if="show" class="modal-mask">
      <div class="modal-wrapper">
        <div class="modal-container">
          <div class="modal-header">
            <slot name="header">default header</slot>
            <button class="like-btn" @click="closeModal">
              <svg class="feather">
                <use href="/feather-sprite-v4.29.0.svg#x" />
              </svg>
            </button>
          </div>

          <div class="modal-body">
            <slot name="body">
              <div class="search-results">
                <div v-for="response in convs" :key="response.conversation.conversationId"
                  @click="forwardMessage(response.conversation.conversationId, response.user)">
                  <RouterLink :to="'/conversation/' + response.conversation.conversationId" class="custom-link" replace
                    force>
                    <div class="user" v-if="response.conversation.conversationId !== convId">
                      <p>{{ response.user.username }}</p>
                    </div>
                  </RouterLink>
                </div>
              </div>
            </slot>
          </div>
        </div>
      </div>
    </div>
  </Transition>
</template>

<style>
.custom-link {
  color: inherit;
  /* This will make the link have the same color as the surrounding text */
  text-decoration: none;
  /* This will remove the underline */
}

.modal-mask {
  position: fixed;
  z-index: 9998;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: table;
  transition: opacity 0.3s ease;
}

.modal-wrapper {
  display: table-cell;
  vertical-align: middle;
}

.modal-container {
  width: 350px;
  margin: 0px auto;
  background-color: #fff;
  border-radius: 2px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.33);
  transition: all 0.3s ease;
}

.modal-header {
  height: 70px;
  padding: 20px 15px 10px 15px;
}

.modal-header h3 {
  margin-top: 0;
  font-size: 25px;
  color: #42b983;
}

.modal-header button {
  color: rgb(86, 86, 86);
  background: none;
  border: none;
  padding: 5px;
  line-height: 12px;
  font-size: 15px;
}

.modal-header button svg {
  width: 20px;
  height: 20px;
}


.search-input {
  padding: 0 15px;
}

.search-input input {
  height: 30px;
  width: 100%;
  outline: none;
  border-radius: 3px;
  border: 1px solid rgb(179, 179, 179)
}

.search-results {
  font-size: 15px;
  padding: 10px 15px;
  border-bottom: 1px solid #eee;
  cursor: pointer;
  max-height: 200px;
  overflow-y: scroll;
}

.modal-default-button {
  float: right;
}

.username-form {
  display: flex;
  flex-direction: column;
  padding: 0 15px;
}

.username-form input {
  margin-bottom: 10px;
  margin-top: 5px;
  outline: none;
  border-radius: 3px;
  border: 1px solid rgb(179, 179, 179)
}

.username-form button {
  margin-bottom: 15px;
}
</style>
