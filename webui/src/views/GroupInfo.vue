<!-- 
-- Pagina per visualizzare le informazioni del gruppo 

In questa pagina vengono mostrare le informazioni del gruppo selezionato dall'utente (il quale partecipa ad esso)
L'utente loggato può:
- visualizzare il nome del gruppo, la foto del gruppo e i membri del gruppo.
- aggiornare il nome del gruppo
- aggiornare la foto del gruppo
- uscire dal gruppo

-->

<script>
import AddGroup from '../components/ModalAddGroup.vue';
import Modal from '../components/Modal.vue';

export default {
  data() {
    return {
      errorMsg: "",

      // Verifica per il campo newGroupname se è valido o no
      usernameValidation: new RegExp('^\\w{0,16}$'),

      // Lista di utenti filtrati dalla ricerca
      filteredUsers: [],

      // Username dell'utente loggato
      owner: sessionStorage.username,

      selectedUsers: [], // Lista di utenti selezionati
      // Utilizzato per mostrare o nascondere il modal per aggiungere utenti al gruppo
      addGroupVisible: false,

      newGroupname: "", // Nuovo nome del gruppo
      newProPic: null,  // Nuova foto del gruppo

      // Utilizzato per mostrare il modale di aggiornamento del nome del gruppo
      updateNameModalIsVisible: false,
      // Utilizzato per mostrare il modale di aggiornamento della foto del gruppo
      updateProPicIsVisible: false,

      // Group info
      groupName: localStorage.username,
      groupId: localStorage.userID,
      groupPhoto: localStorage.photo,
      groupMembers: JSON.parse(localStorage.users),
    };
  },
  emits: ['login-success'],
  methods: {
    // Funzione utilizzata per controllare se il file inserito dall'utente è del formato corretto
    async handleFileChange(event) {
      this.errorMsg = "";
      const file = event.target.files[0]; // Prende il file inserito dall'utente
      if (file.type !== "image/jpeg") {
        this.errorMsg = "File type not supported, only jpg and jpeg are allowed";
        return
      }
      if (file.size > 5242880) {
        this.errorMsg = "File size is too big. Max size is 5MB";
        return
      }
      this.newProPic = file; // Assegna il file inserito dall'utente alla variabile newProPic
    },
    // Funzione per mostrare o nascondere il modal di aggiornamento della foto del gruppo
    handleUpdateProPicToggle() {
      localStorage.photo = this.newProPic; // Assegna la nuova foto del gruppo al localStorage
      this.updateProPicIsVisible = !this.updateProPicIsVisible;  // Nasconde o mostra il modale
      this.newProPic = "";
      this.errorMsg = "";
    },
    // Funzione per mostrare o nascondere il modale di aggiornamento del nome del gruppo
    handleUpdateNameToggle() {
      localStorage.username = this.groupName; // Aseegna il nuovo nome del gruppo al localStorage
      this.updateNameModalIsVisible = !this.updateNameModalIsVisible; // Nasconde o mostra il modale
      this.newUsername = "";
      this.errorMsg = "";
    },
    // Funzione che aggiorna la foto del gruppo
    async updateProPic() {
      this.errorMsg = "";
      // Crea un oggetto FormData e vi aggiunge la foto del gruppo
      const formData = new FormData();
      formData.append('image', this.newProPic);
      // Effettua la richiesta al server per aggiornare la foto del gruppo
      this.$axios.put(`/profiles/${sessionStorage.userID}/groups/${this.groupId}/grouphoto`, formData, { headers: { 'Authorization': `${sessionStorage.token}` } })
        .then(response => {
          // Assegna la nuova foto del gruppo alla variabile groupPhoto
          this.groupPhoto = response.data.photo;
          this.handleUpdateProPicToggle(); // Nasconde il modale e aggiorna la foto del gruppo in localStorage
        })
        .catch(e => {
          this.errorMsg = e.toString
        });
    },
    // Funzione che aggiorna il nome del gruppo
    async updateGroupName() {
      // Controlla se il nome del grupop inserito è uguale a quello attuale
      if (this.newGroupname == this.groupName) {
        this.errorMsg = "You must enter a new username";
        return
      }
      // Controlla se la lunghezza del nome scelto è valido 
      if (this.newGroupname.length < 3 || this.newGroupname.length > 16) {
        this.errorMsg = "Invalid username, it must contains min 3 characters and max 16 characters";
        return
      }
      // Controlla se il nome scelto è valido
      if (!this.usernameValidation.test(this.newGroupname)) {
        this.errorMsg = "Invalid username, it must contain only letters and numbers";
        return
      }
      try {
        // Effettua la richiesta del server per aggiornare il nome del gruppo
        let _ = await this.$axios.put(`/profiles/${sessionStorage.userID}/groups/${this.groupId}/groupname`, { groupName: this.newGroupname }, { headers: { 'Authorization': `${sessionStorage.token}` } })
        // Assegna il nuovo nome del gruppo alla variabile groupName
        this.groupName = this.newGroupname;
        this.handleUpdateNameToggle(); // Nasconde il modale e aggiorna il nome del gruppo in localStorage
      } catch (e) {
        if (e.response.data == "Username already exist\n") {
          this.errorMsg = "This username is already taken. Please try another one.";
        } else {
          this.errorMsg = e.toString();
        }
      }
    },
    // Funzione utilizzata per far uscire l'utente loggato dal gruppo
    leaveGroup() {
      this.errorMsg = "";
      // Effettua la richiesta al server per far uscire l'utente dal gruppo
      this.$axios.delete(`/profiles/${sessionStorage.userID}/groups/${this.groupId}`, { headers: { 'Authorization': `${sessionStorage.token}` } })
        .then(() => {
          // Fa ritornare l'utente alla home dopo l'uscita dal gruppo
          this.$router.push("/home");
        })
        .catch(e => {
          this.errorMsg = e.toString();
        });
    },
    // Funzione per mostrare o nascondere il modale per aggiungere utenti al gruppo
    handleAddGroupModalToggle() {
      this.addGroupVisible = !this.addGroupVisible;
    },
  },
  components: { AddGroup, Modal },
}
</script>

<template>
  <div>
    <div
      class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <!-- Group photo -->
      <div class="top-profile-container">
        <img :src="`data:image/jpg;base64,${groupPhoto}`">
      </div>

      <!-- Modali della pagina -->

      <!-- Modale utlizzato per aggiornare il nome del gruppo -->
      <Modal :show="updateNameModalIsVisible" @close="handleUpdateNameToggle" title="username">
        <template v-slot:header>
          <h3>Update Group Name</h3>
        </template>
        <template v-slot:body>
          <!-- Input per l'inserimento del nuovo nome per il gruppo -->
          <form class="username-form">
            <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
            <input type="text" v-model="this.newGroupname" placeholder="New group name" />
            <button type="submit" @click.prevent="updateGroupName">Update</button>
          </form>
        </template>
      </Modal>
      <!-- Modale utlizzato per aggiornare la foto del gruppo  -->
      <Modal :show="updateProPicIsVisible" @close="handleUpdateProPicToggle" title="photo">
        <template v-slot:header>
          <h3>Update Group Picture</h3>
        </template>
        <template v-slot:body>
          <!-- Input per l'inserimento della nuova foto per il gruppo -->
          <form class="username-form">
            <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
            <input type="file" ref="file" accept=".jpg,.jpeg" @change="handleFileChange" />
            <button type="submit" @click.prevent="updateProPic">Update</button>
          </form>
        </template>
      </Modal>
      <!-- Modale utilizzato per aggiungere utenti al gruppo -->
      <AddGroup :show="addGroupVisible" @close="handleAddGroupModalToggle" title="search">
        <template v-slot:header>
          <h3>Add to group</h3>
        </template>
      </AddGroup>

      <!-- Body della pagina -->

      <!-- Group name -->
      <h1 class="h1">{{ this.groupName }}</h1>
      <div class="btn-toolbar mb-2 mb-md-0">
        <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
        <!-- Pulsante per aggiornare il nome del gruppo -->
        <div class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="handleUpdateNameToggle">
            Change group name
          </button>
          <!-- Pulsante per l'aggiornamento della foto del gruppo -->
          <button type="button" class="btn btn-sm btn-outline-secondary" @click="handleUpdateProPicToggle">
            Change group photo
          </button>
        </div>
        <!-- Pulsante per aggiungere utenti al gruppo -->
        <button type="button" class="btn btn-sm btn-outline-primary" @click="handleAddGroupModalToggle">
          Add to group
        </button>
      </div>
    </div>
    <!-- Lista dei membri del gruppo -->
    <div v-for="user in groupMembers">
      <p class="username">
        <!-- Foto profilo dell'utente del gruppo -->
        <img :src="`data:image/jpg;base64,${user.photo}`" class="profile-picture">
        {{ user.username }}
        <!-- Se l'utente è quello loggato allora mostra il pulstante per uscire dal gruppo -->
        <button type="button" v-if="user.username == owner" class="btn btn-sm btn-outline-primary" @click="leaveGroup">
          Leave Group
        </button>
      </p>
    </div>

    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
  </div>
</template>

<!-- Stili utilizzati per mostarer il nome e la foto degli utenti membri del gruppo -->
<style>
.profile-picture {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  margin-right: 10px;
  object-fit: cover;
}

.username {
  font-size: 14px;
  font-weight: bold;
}
</style>
