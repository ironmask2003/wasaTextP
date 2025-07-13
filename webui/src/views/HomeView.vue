<!--

Pagina princiaple del sito in cui vengono mostrate le conversazioni dell'utente loggato.

In questa pagina l'utente può:
- visualizzare le conversazioni con altri utenti o gruppi
- aggiornare la lista delle conversazioni
- creare un nuovo gruppo
- creare una nuova conversazione con un utente

-->

<script>
import Modal from '../components/Modal.vue';
import Group from '../components/ModalGroup.vue';

export default {
  data: function () {
    return {
      errormsg: null,

      // Dati delle conversazioni
      some_data: [],

      // Utilizzato per mostrare o nascondere il modal di ricerca
      searchModalIsVisible: false,

      // Utilizzato per mostrare o nascondere il modal di creazione di un gruppo
      createGroupModalIsVisible: false,

      // Lista di utenti del nuovo gruppo
      users: [],
    }
  },
  emits: ['login-success', 'username-changed'],
  methods: {
    // Fiunzione per ottente la conversazioni dell'utente
    async getConversations() {
      this.errormsg = null;
      try {
        // Effettua la richiesta al server per ottenere le conversazioni dell'utente
        let response = await this.$axios.get(`/profiles/${sessionStorage.userID}/conversations`, { headers: { 'Authorization': sessionStorage.token } });
        // Salva i dati delle conversazioni
        this.some_data = response.data;
      } catch (e) {
        this.errormsg = e.toString();
      }
    },
    // Funzione utilizzata per poratre l'utente alla conversazione selezionata
    goToConversation(response) {
      localStorage.clear();
      // Controlla se la conversazione selezionata è di gruppo o no
      if (response.group.groupId !== 0) {
        // Salva i dati del gruppo nella localStorage
        localStorage.userID = response.group.groupId;
        localStorage.username = response.group.groupName;
        localStorage.photo = response.group.photo;
        localStorage.users = JSON.stringify(response.groupUsers);
      } else {
        // Altrimenti salva i dati dell'utente con cui conversare
        localStorage.userID = response.user.userId;
        localStorage.username = response.user.username;
        localStorage.photo = response.user.photo;
      }
      // Reindirizza l'utente alla conversazione
      this.$router.push(`/conversation/${response.conversation.conversationId}`);
    },
    // Funzione per mostrare o nascondere il modal di ricerca di utenti per iniziare una nuova conversazione
    handleSearchModalToggle() {
      this.searchModalIsVisible = !this.searchModalIsVisible;
    },
    // Funzione per mostare o nascondere il modale per la creazione di un nuovo gruppo
    handleCreateGroupModalToggle() {
      this.createGroupModalIsVisible = !this.createGroupModalIsVisible;
    },
  },
  mounted() {
    // Se l'utente non è loggato, reindirizza alla pagina di login
    if (!sessionStorage.token) {
      this.$router.push("/");
      return;
    }

    this.getConversations();

    window.addEventListener("message", (event) => {
      if (event.data.type === "ws-message") {
        console.log("Received message from WebSocket:", event.data.data);
        const messaggio = JSON.parse(event.data.data);

        // Controlla se il messaggio ricevuto è relativo all'aggiornamento dell'username
        if (messaggio.message == "Username Updated") {
          this.getConversations(); // Aggiorna la lista delle conversazioni
        }
      }
    });
  },
  components: { Modal, Group }
}
</script>

<template>
  <div>
    <div
      class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <h1 class="h2">Chat</h1>

      <!-- Modale utilizzato per la creazione di un nuovo gruppo -->
      <Group :show="createGroupModalIsVisible" @close="handleCreateGroupModalToggle" title="search">
        <template v-slot:header>
          <h1 class="profile-title">Select users</h1>
        </template>
      </Group>
      <!-- Modale utilzzato per la ricerca degli utenti con cui aprire una nuova conversazione -->
      <Modal :show="searchModalIsVisible" @close="handleSearchModalToggle" title="search">
        <template v-slot:header>
          <h1 class="profile-title">Users</h1>
        </template>
      </Modal>

      <!-- Pulsanti per aggiornare la lista delle conversazioni, creare un nuovo gruppo e cercare nuovi utenti -->
      <div class="btn-toolbar mb-2 mb-md-0">
        <!-- Pulsante per creare un nuovo gruppo -->
        <div class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-primary" @click="handleCreateGroupModalToggle">
            New Group
          </button>
        </div>
      </div>
    </div>

    <!-- Lista delle conversazioni -->
    <div v-if="some_data.length !== 0">
      <!-- Mostra le conversazioni dell'utente iterando all'interno di some_data dove sono salvate tutte le conversazioni -->
      <div class="conversations" v-for="response in some_data" :key="response.conversation.conversationId">
        <!-- Controlla se la conversazione non è con un gruppo -->
        <div v-if="response.group.groupName == ''">
          <!-- Mostra il nome dell'utente con cui si sta conversando, l'ultimo messaggio e chi lo ha inviato -->
          <!-- Se il messaggio è una foto, mostra "Photo" al posto del testo -->
          <button v-if="response.message.photo == ''" type="button" class="btn btn-sm btn-outline-primary"
            @click="goToConversation(response)">
            {{ response.user.username }} <br> {{ response.senderUser.username }}: {{ response.message.text }}
          </button>
          <!-- Altrimenti mostra il contenuto del messaggio -->
          <button type="button" class="btn btn-sm btn-outline-primary" @click="goToConversation(response)" v-else>
            {{ response.user.username }} <br> {{ response.senderUser.username }}: Photo
          </button>
        </div>
        <!-- Se la conversazione è con un gruppo -->
        <div v-else>
          <!-- Mostra il nome del gruppo, l'ultimo messaggio e chi lo ha inviato come per il singolo utente -->
          <div v-if="response.message.photo == ''" type="button" @click="goToConversation(response)">
            <div>
            <img :src="`data:image/jpg;base64,${response.group.photo}`" class="profile-picture position-picture" /> 
            <h6 class="subtitle subtitle-black">{{ response.group.groupName }} <br/> <span class="last-msg"> {{ response.senderUser.username }}: {{ response.message.text }} </span> </h6>
            </div>
          </div>
          <button type="button" class="btn btn-sm btn-outline-primary" @click="goToConversation(response)" v-else>
            {{ response.group.groupName }} <br> {{ response.senderUser.username }}: Photo
          </button>
        </div>
      </div>
    </div>

    <!-- Se non ci sono conversazioni, mostra un messaggio -->
    <div v-else>
      <p>Start a conversation</p>
    </div>

    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
  </div>
</template>

<style>
.subtitle-black {
  color: black;
  display: inline;
  position: relative;
}

.position-picture {
  position: relative;
  top: 10px;
}

.last-msg {
  position: relative;
  left: 50px;
  top: -10px;
  font-weight: 100;
}

.conversations {
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 8px;
  transition: background-color 0.2s ease;
}

.conversations:hover {
  cursor: pointer;
  background-color: #f5f5f5;
  border-radius: 8px;
}

.conversations:active {
  background-color: #e0e0e0;
}
</style>
