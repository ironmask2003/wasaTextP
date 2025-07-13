<!-- 

Modale utilizzato per la ricerca di utenti

-->

<script>

export default {
  // Props passati al componente
  props: {
    show: Boolean,  // Utilizzato per mostrare o nascondere il modal
    users: Array, // Lista di utenti da filtrare
    title: String,  // Titolo del modale
  },
  data() {
    return {
      errorMsg: "",

      // Variabile utilizzate per la ricerca degli utenti
      searchText: "",

      // Utilizzato per la verifica dell'username inserito
      usernameValidation: new RegExp('^\\w{0,16}$'),

      // Username dell'utente che ha effettuato il login
      username: sessionStorage.username,

      // Lisat di utenti filtrati in base all'username inserito
      filteredUsers: [],
    };
  },
  methods: {
    // Funzione utilizzata per chiudere il modale salvando in localStorage l'utente selezionato con cui aprire un conversazione
    closeModal() {
      this.searchText = "";
      window.location.reload();
      this.$emit('close');
    },
    // Funzione utilizzata per filtrare gli utenti in base all'username inserito
    async filterUsers() {
      this.errorMsg = "";

      if (this.searchText.length >= 0) {
        if (this.searchText.length == 0){
          this.filteredUsers = [];
          return
        }
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
            // In base al risultato della GET assegna la lista di utenti filtrati
            if (response.data == null) {
              this.filteredUsers = [];
              return;
            }
            if (this.filteredUsers != response.data) {
              this.filteredUsers = response.data;            
              this.filteredUsers = this.filteredUsers.filter(user => user.username !== this.owner);
            }
          } catch (e) {
            this.errorMsg = e.toString();
            this.filteredUsers = [];
          }
        } else {
          this.filteredUsers = this.users.filter(user => user.username.toLowerCase().includes(this.searchText.toLowerCase()));
        }
      }
    },
    // Funzione utilizzara per selezionare l'utente con cui aprire una conversazione
    async selectUser(userToSend, username, photo) {
      localStorage.clear();
      localStorage.userID = userToSend;
      localStorage.username = username;
      localStorage.photo = photo;
      this.$router.push('/conversation');
      closeModal();
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
              <!-- Input per l'inserimento dell'utente da cercare -->
              <div class="search-input">
                <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                <input type="text" v-model="searchText" placeholder="Search" />
              </div>
              <!-- Risultato della ricerca -->
              <div class="search-results" v-if="filteredUsers">
                <div v-for="user in filteredUsers" :key="user.userId"
                  @click="selectUser(user.userId, user.username, user.photo)">
                    <div class="user">
                      <p v-if="user.username != username" >{{ user.username }}</p>
                    </div>
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
  background-color: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
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
  background-color: white;
  border-radius: 16px !important;
  overflow: hidden;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.2);
  transition: all 0.3s ease;
}

.modal-header {
  height: 70px;
  padding: 20px 15px 10px 15px;
  background: linear-gradient(135deg, #25d366 0%, #128c7e 100%);
  color: white;
  border-radius: 16px 16px 0 0 !important;
  position: relative;
}

.modal-header h3 {
  margin-top: 0;
  font-size: 25px;
  color: white;
  font-weight: 600;
}

.modal-header button {
  background: none;
  border: none;
  padding: 8px;
  line-height: 12px;
  font-size: 15px;
  position: absolute;
  top: 20px;
  right: 15px;
  border-radius: 50%;
  transition: all 0.2s ease;
  opacity: 0.8;
}

.modal-header button:hover {
  background: rgba(255, 255, 255, 0.2);
  opacity: 1;
  transform: scale(1.1);
}

.modal-header button svg {
  width: 20px;
  height: 20px;
  filter: brightness(0) invert(1);
}

.modal-body {
  background-color: white;
  border-radius: 0 0 16px 16px !important;
  padding: 0;
  margin: 0;
  overflow: hidden;
}

.search-input {
  padding: 15px 15px 0 15px;
  background-color: white;
  overflow: hidden;
}

.search-input input {
  height: 30px;
  width: 100%;
  outline: none;
  border-radius: 8px;
  border: 2px solid #e5e5e5;
  padding: 8px 12px;
  transition: border-color 0.2s ease;
}

.search-input input:focus {
  border-color: #25d366;
}

.search-results {
  font-size: 15px;
  padding: 10px 15px 20px 15px;
  border-bottom: none;
  cursor: pointer;
  max-height: 200px;
  overflow-y: scroll;
  background-color: white;
  border-radius: 0 0 16px 16px !important;
  box-sizing: border-box;
}

.modal-default-button {
  float: right;
}

.username-form {
  display: flex;
  flex-direction: column;
  padding: 15px 15px 20px 15px;
  background-color: white;
  border-radius: 0 0 16px 16px !important;
  margin: 0;
  box-sizing: border-box;
}

.username-form input {
  margin-bottom: 10px;
  margin-top: 5px;
  outline: none;
  border-radius: 8px;
  border: 2px solid #e5e5e5;
  padding: 8px 12px;
  transition: border-color 0.2s ease;
}

.username-form input:focus {
  border-color: #25d366;
}

.username-form button {
  margin-bottom: 0;
  border-radius: 8px;
  background: linear-gradient(135deg, #25d366 0%, #128c7e 100%);
  border: none;
  color: white;
  padding: 10px 20px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
}

.username-form button:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(37, 211, 102, 0.4);
}

/* Animazioni per l'apertura/chiusura del modale */
.modal-enter-active, .modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-active .modal-container,
.modal-leave-active .modal-container {
  transition: transform 0.3s ease-out;
}

.modal-enter-from {
  opacity: 0;
}

.modal-enter-from .modal-container {
  transform: scale(0.9) translateY(-20px);
}

.modal-leave-to {
  opacity: 0;
}

.modal-leave-to .modal-container {
  transform: scale(0.9) translateY(-20px);
}

/* Regole specifiche per forzare i border-radius */
.modal-mask .modal-wrapper .modal-container {
  border-radius: 16px !important;
}

.modal-mask .modal-wrapper .modal-container .modal-header {
  border-radius: 16px 16px 0 0 !important;
}

.modal-mask .modal-wrapper .modal-container .modal-body {
  border-radius: 0 0 16px 16px !important;
}

.modal-mask .modal-wrapper .modal-container .modal-body .search-results {
  border-radius: 0 0 16px 16px !important;
}

.modal-mask .modal-wrapper .modal-container .modal-body .username-form {
  border-radius: 0 0 16px 16px !important;
}
</style>
