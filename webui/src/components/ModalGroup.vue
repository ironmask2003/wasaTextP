<script>
import { RouterLink } from 'vue-router';

export default {
  props: {
    show: Boolean,
    users: Array,
    title: String,
  },
  data() {
    return {
      searchText: "",
      groupName: "",
      errorMsg: "",
      usernameValidation: new RegExp('^\\w{0,16}$'),
      filteredUsers: [],
      owner: sessionStorage.username,
      selectedUsers: [], // Lista di utenti selezionati
    };
  },
  methods: {
    closeModal() {
      this.searchText = "";
      this.groupName = "";
      this.selectedUsers = [];
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

        if (this.title === "search") {
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
        } else {
          this.filteredUsers = this.users.filter(user => user.username.toLowerCase().includes(this.searchText.toLowerCase()));
        }
      }
    },
    async createGroup() {
      if(this.groupName.length < 3 || this.groupName.length > 16){
        this.errorMsg = "Invalid group name, it must contains min 3 characters and max 16 characters";
        return;
      }
      try {
        let response = await this.$axios.post(`/profiles/${sessionStorage.userID}/groups`, {
          groupName: this.groupName,
          users: this.selectedUsers,
        }, { headers: { 'Authorization': `${sessionStorage.token}` } });
        localStorage.clear();
        localStorage.userID = response.data.group.GroupId;
        localStorage.username = response.data.group.groupName;
        localStorage.photo = response.data.group.photo;
        this.closeModal();
        this.$router.push(`/conversations/${response.data.conversation.conversationId}`);
      } catch (e) {
        this.errorMsg = e.toString
      }
    },
    selectUser(user) {
      // Controlla se l'utente è già nella lista
      if (!this.selectedUsers.find(u => u.username === user.username)) {
        this.selectedUsers.push(user); // Aggiungi utente selezionato
      }
    },
    removeUser(username) {
      // Rimuove l'utente dalla lista selezionata
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
  components: { RouterLink }
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
              <!-- Selezione del nome del gruppo -->
              <div class="search-input">
                <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                <input type="text" v-model="groupName" placeholder="Select group name" />
              </div>
              <!-- Campo di ricerca -->
              <div class="search-input">
                <input type="text" v-model="searchText" placeholder="Search" />
              </div>
              <p></p>
              <div class="btn-group me-2">
                <button class="btn btn-sm btn-outline-primary" @click="createGroup">Create Group</button>
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
                <span class="selected-user">{{owner}}</span>
                <div v-for="user in selectedUsers" :key="user.userId" class="selected-user">
                  <span>{{ user.username }}</span>
                  <button v-if="user.username !== owner" @click="removeUser(user.username)">Remove</button>
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
.selected-users {
  margin-top: 20px;
  padding: 10px;
  border-top: 1px solid #ccc;
}

.selected-user {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.selected-user span {
  font-size: 14px;
  font-weight: bold;
}

.selected-user button {
  background: red;
  color: white;
  border: none;
  border-radius: 5px;
  padding: 5px 10px;
  cursor: pointer;
}
</style>
