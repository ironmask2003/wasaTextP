<!-- 

Modale utilizzato per aggiungere utenti al gruppo

L'utente può:
- cercare un utente specifico per aggiungerlo al gruppo
- visualizzare la lista di utenti aggiunti

-->

<script>
export default {
  props: {
    show: Boolean,
    users: Array,
    title: String,
  },
  data() {
    return {
      errorMsg: "",

      // Utilizzato per la ricerca degli utenti da aggiungere al gruppo
      searchText: "",

      // Utilizzato per la verifica dell'username inserito e del nome del gruppo
      usernameValidation: new RegExp('^\\w{0,16}$'),

      // Lista degli utenti filtrati in base alla ricerca effettuata
      filteredUsers: [],

      // Username dell'utente loggato
      owner: sessionStorage.username,

      // Id del gruppo a cui aggiungere gli utenti
      groupId: localStorage.userID,
      selectedUsers: [], // Lista di utenti selezionati
    };
  },
  methods: {
    // Chiude il modale
    closeModal() {
      this.searchText = "";
      this.selectedUsers = [];
      this.$emit('close');
    },
    // Funzione utilizzata per la ricerca degli utenti da aggiungere al gruppo
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
            // Effettua una richiesta GET al server per ottenere gli utenti in base alla ricerca effettuata
            const url = `/profiles?username=${this.searchText}`;
            let response = await this.$axios.get(url, { headers: { 'Authorization': `${sessionStorage.token}` } });
            // In base al risultato della GET assegna la lista degli utenti
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
    // Funzione utilizzata per aggiungere utenti ad un gruppo
    async addToGroup() {
      try {
        // Effettua la richiesta al server per aggiungere gli utenti selezionati al gruppo
        let response = await this.$axios.put(`/profiles/${sessionStorage.userID}/groups/${this.groupId}`, {
          users: this.selectedUsers,
        }, { headers: { 'Authorization': `${sessionStorage.token}` } });
        // Salva i dati del gruppo in localStorage e reindirizza l'utente alla pagina del gruppo aggiornata
        localStorage.clear();
        localStorage.userID = response.data.group.groupId;
        localStorage.username = response.data.group.groupName;
        localStorage.photo = response.data.group.photo;
        localStorage.users = JSON.stringify(response.data.members);
        // Chiude il modale
        this.closeModal();
        // Renderizza la pagina del gruppo aggiornata
        window.location.reload();
        this.$router.push(`/groups/${response.data.group.groupId}`);
      } catch (e) {
        this.errorMsg = e.toString
      }
    },
    // Funzione utilizzata per selezionare un utente da aggiungere al gruppo
    selectUser(user) {
      // Controlla se l'utente è già nella lista
      if (!this.selectedUsers.find(u => u.username === user.username)) {
        this.selectedUsers.push(user); // Aggiungi utente selezionato
      }
    },
    // Funzione utilizzata per rimuovere un utente dalla lista selezionata
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
              <!-- Campo di ricerca -->
              <div class="search-input">
                <input type="text" v-model="searchText" placeholder="Search" />
              </div>
              <p></p>
              <div class="btn-group me-2">
                <button class="btn btn-sm btn-outline-primary" @click="addToGroup">Add To Group</button>
              </div>

              <!-- Risultati della ricerca -->
              <div class="search-results">
                <div v-for="user in filteredUsers" :key="user.userId" @click="selectUser(user)" class="user">
                  <p v-if="user.username !== owner">{{ user.username }}</p>
                </div>
              </div>

              <!-- Lista di utenti selezionati -->
              <div class="selected-users">
                <h4>Selected Users:</h4>
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
