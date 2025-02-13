<!--

Modale utilizzato per la ricerca di conversazioni per inoltrare un messaggio

L'utente loggato può:
- visualizzare la lista di conversazioni
- inoltrare il messaggio alla conversazione selezionata

-->


<script>
import Notification from './Notification.vue';

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
      users: [],          // Lista di tutti gli utenti
      searchText: "",     // Testo della ricerca
      filteredUsers: [],  // Lista degli utenti filtrati
      selectedUsers: [],   // Lista degli utenti selezionati

      // Utilizzato per la verifica dell'username inserito e del nome del gruppo
      usernameValidation: new RegExp('^\\w{0,16}$'),

      handleNotify: false,
    };
  },
  methods: {
    triggerNotification() {
      this.$refs.notification.show()
      setTimeout(() => {
        this.closeModal();
      }, 500);
    },
    closeModal() {
      this.convs = [];
      window.location.reload();
      this.$emit('close');
    },
    async filterUsers() {
      this.errorMsg = "";
      this.filteredUsers = this.users;

      if (this.searchText.length > 0) {
        if (this.searchText.length > 16 || !this.usernameValidation.test(this.searchText)) {
          this.errorMsg = "Invalid username, it can contain only letters and numbers for a maximum of 16 characters.";
          this.filteredUsers = [];
          return;
        }
        try {
          const url = `/profiles?username=${this.searchText}`;
          let response = await this.$axios.get(url, { headers: { 'Authorization': `${sessionStorage.token}` } });
          if (response.data == null) {
            this.filteredUsers = [];
            return;
          }
          this.filteredUsers = response.data;
        } catch (e) {
          this.errorMsg = e.toString();
          this.filteredUsers = [];
        }
      }
    },
    async forwardMessage() {
      this.errorMsg = null;
      for (let user of this.selectedUsers) {
        const dest_user = user.userId; // Assumi che l'ID della conversazione per l'utente selezionato sia disponibile in user.convId
        const url = `/profiles/${sessionStorage.userID}/conversations/${this.convId}/messages/${this.msg.messageId}?dest_user=${dest_user}`;
        await this.$axios.post(url, {}, { headers: { 'Authorization': sessionStorage.token } })
          .then(() => {
            this.triggerNotification();
          })
          .catch(e => {
            this.errorMsg = e.toString();
          });
      }
    },
    selectUser(user) {
      if (!this.selectedUsers.find(u => u.username === user.username)) {
        this.selectedUsers.push(user);
      }
    },
    removeUser(username) {
      this.selectedUsers = this.selectedUsers.filter(user => user.username !== username);
    },
  },
  watch: {
    searchText() {
      this.filterUsers();
    },
    show() {
      this.filteredUsers = this.users;
    }
  },
  components: { Notification },
}
</script>

<template>
  <Transition name="modal">
    <div v-if="show" class="modal-mask">
      <Notification ref="notification" message="Operazione completata!" :duration="3000"/>
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
              <!-- Campo di ricerca -->
              <div class="search-input">
                <input type="text" v-model="searchText" placeholder="Search" />
              </div>
              <p></p>
              <div class="btn-group me-2">
                <button class="btn btn-sm btn-outline-primary" @click="forwardMessage">Forward Message</button>
              </div>

              <!-- Risultati della ricerca -->
              <div class="search-results">
                <div v-for="user in filteredUsers" :key="user.userId" @click="selectUser(user)" class="user">
                  <p>{{ user.username }}</p>
                </div>
              </div>

              <!-- Lista di utenti selezionati -->
              <div class="selected-users">
                <h4>Selected Users:</h4>
                <div v-for="user in selectedUsers" :key="user.userId" class="selected-user">
                  <span>{{ user.username }}</span>
                  <button @click="removeUser(user.username)">Remove</button>
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
  text-decoration: none;
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
