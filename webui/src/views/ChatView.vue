<!-- 

Pagina utilizzata per visulizzare una conversazione con un utente o con un gruppo

L'utente loggato può:
- inviare messaggi testuali o con foto
- eliminare messaggi
- inoltrare messaggi
- lasciare commenti ai messaggi e eliminarli
- visualizzare i commenti ai messaggi
- visualizzare i messaggi della conversazione

-->

<script>
import Modal from '../components/ModalConv.vue';
import Comments from '../components/ModalComments.vue';

export default {
  data() {
    return {
      errormsg: null,

      // Dati della conversazione
      userToSend: localStorage.username,  // Nome dell'utente con cui si sta conversando o del gruppo
      userIdToSend: localStorage.userID,  // ID dell'utente con cui si sta conversando o del gruppo
      proPic64: localStorage.photo, // Foto dell'utente con cui si sta conversando o del gruppo
      convId: parseInt(this.$route.params.convId),  // ID della conversazione (se è una nuova conversazione, sarà undefined)

      // Testo del messaggio da inviare
      text: null,
      // Foto del messaggio da inviare
      photo: null,

      // Lista dei messaggi della conversazione
      some_data: [],

      // Utilizzata per verificare se si sta conversando è un gruppo o no
      isGroup: localStorage.users ? true : false,

      // Utilizzato per mostrare il modale di ricerca di una conversazione a cui inoltrare un messaggio selezionato
      searchModalIsVisible: false,

      // Utilizzato per mostrare il modale per lasciare un commento al messaggio selezionato
      commentModalIsVisible: false,

      // Messaggio da inoltrare
      messageToFordward: null,

      // Messaggio da commentare
      messageToComment: null,

      // Commenti del messaggio selezionato
      comments: null,
    }
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
      this.photo = file;  // Assegna il file inserito dall'utente alla variabile photo
    },
    // Funzione che prende i messaggi della conversazione
    async getConversation() {
      this.errormsg = null;
      // Effettua la richiesta al server per ottenre i messaggi della conversazione
      this.$axios.get(`/profiles/${sessionStorage.userID}/conversations/${this.convId}`, { headers: { 'Authorization': sessionStorage.token } })
        .then(response => {
          // Salva i messaggi nella variabile some_data
          this.some_data = response.data;
        })
        .catch(e => {
          this.errormsg = e.toString();
        });
    },
    // Funzione che controlla se la convId è stata presa correttamente dai parametri
    check() {
      // Se è NaN o undefined allora crea la conversazione
      if (isNaN(this.convId) || this.convId == undefined) {
        this.createConversation();
      }
      // Altrimenti manda il messaggio
      else {
        this.sendMessage();
      }
    },
    // Funzione che crea la conversazione
    createConversation() {
      this.errormsg = null;
      // Effettua la richiesta al server per creare una nuova conversazione aggiungendo il messaggio scritto dall'utente come primo messaggio
      this.$axios.put(`/profiles/${sessionStorage.userID}/conversations/${this.userIdToSend}`, {
        text: this.text
      }, { headers: { 'Authorization': sessionStorage.token } })
        .then(response => {
          // Assegna la convId
          this.convId = response.data.conversationId;
          // Reindirizza alla pagina della conversazione appena creata
          this.$router.push(`/conversation/${this.convId}`);
        })
        .catch(e => {
          this.errormsg = e.toString();
        });

    },
    // Funzione che elimina un messaggio
    async deleteMessage(msgId) {
      this.errormsg = null;
      // Effettua la richiesta al server per eliminare il messaggio selezionato
      this.$axios.delete(`/profiles/${sessionStorage.userID}/conversations/${this.convId}/messages/${msgId}`, { headers: { 'Authorization': sessionStorage.token } })
        .then(() => {
          // Aggiorna i messaggi della conversazione
          this.getConversation();
        })
        .catch(e => {
          this.errormsg = e.toString();
        });
    },
    // Funzione che invia un messaggio
    async sendMessage() {
      this.errormsg = null;
      // Crea un form data con il testo e la foto del messaggio
      const formData = new FormData();
      formData.append('text', this.text);
      // Controlla se è stata inserita una foto e la aggiunge al from data
      if (this.photo != null) {
        formData.append('photo', this.photo);
      }
      // Effettua la richiesta al server per inviare il messaggio
      this.$axios.post(`/profiles/${sessionStorage.userID}/conversations/${this.convId}/messages`, formData, { headers: { 'Authorization': sessionStorage.token } })
        .then(() => {
          // Resetta le variabili utilizzate per mandare il messaggio
          this.text = null;
          this.photo = null;
          // Aggiorna la conversazione
          this.getConversation();
        })
        .catch(e => {
          this.errormsg = e.toString();
        });
    },
    // Funzione utilizzata per mostrare o nascondere il modale per lasciare un commento al messaggio selezionato
    handleCommentModalToggle(cmt, commentsMsg) {
      // Assegna il messaggio selezionato e i commenti del messaggio selezionato alle variabili messageToComment e comments
      this.messageToComment = cmt;
      this.comments = commentsMsg;
      // Mostra o nasconde il modale
      this.commentModalIsVisible = !this.commentModalIsVisible;
    },
    // Funzione utilizzata per mostrare o nascondere il modale per inoltrare un messaggio selezionato
    handleSearchModalToggle(msg) {
      // Assegna il messaggio selezionato alla variabile messageToFordward
      this.messageToFordward = msg;
      // Mostra o nasconde il modale
      this.searchModalIsVisible = !this.searchModalIsVisible;
    },
    // Funzione utlizzata per andare alla pagina delle infromazioni di un gruppo (utilizzata solo nel caso in cui la conversazione è con un gruppo)
    goToGroupInfo() {
      // Reindirizza alla pagina delle informazioni del gruppo
      this.$router.push(`/groups/${this.userIdToSend}`);
    }
  },
  mounted() {
    // Controlla se l'utente è loggato altrimenti reindirizza alla pagina di login
    if (!sessionStorage.token) {
      this.$router.push("/");
      return;
    }
    // Se la convId è stata presa dai parametri allora prendi i messaggi della conversazione
    if (this.convId != undefined && !isNaN(this.convId)) {
      this.getConversation()
    }
  },
  components: { Modal, Comments },
}
</script>

<template>
  <div>
    <div
      class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <!-- User photo -->
      <div class="top-profile-container">
        <img :src="`data:image/jpg;base64,${proPic64}`">
      </div>
      <!-- Controlla se la conversazione è con un gruppo o con un utente -->
      <div v-if="isGroup">
        <!-- Se è un gruppo mostra il nome del gruppo -->
        <h1 class="h1 clickable" @click="goToGroupInfo">{{ this.userToSend }}</h1>
      </div>
      <div v-else>
        <!-- Se è un utente mostra il nome dell'utente -->
        <h1 class="h1">{{ this.userToSend }}</h1>
      </div>

      <!-- Modali della pagina -->

      <!-- Modale utilizzato per lasciare un commento a un messaggio -->
      <Comments :show="commentModalIsVisible" :comments="comments" :msg="messageToComment"
        @close="handleCommentModalToggle" title="comments">
        <template>
          <h3>Comments</h3>
        </template>
      </Comments>
      <!-- Modale utilizzato per selezionare una conversazione in cui inoltrare un messaggio -->
      <Modal :show="searchModalIsVisible" :msg="messageToFordward" @close="handleSearchModalToggle"
        title="conversations">
        <template v-slot:header>
          <h3>Conversations</h3>
        </template>
      </Modal>

      <!-- Body della pagina -->
      <div class="btn-toolbar mb-2 mb-md-0">
        <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
        <!-- Form per inviare una foto -->
        <div class="btn-group me-2">
          <form @submit.prevent="sendMessage">
            <input type="file" ref="file" accept=".jpg,.jpeg" @change="handleFileChange" />
            <!-- Pulsante per invaire la foto -->
            <button type="submit" class="btn btn-sm btn-outline-primary">
              Send photo
            </button>
          </form>
        </div>
      </div>
    </div>
    <!-- Lista dei messaggi della conversazione -->
    <div class="messages" v-for="response in some_data" :key="response.message.messageId">
      <!-- Mostra il contenuto del messaggio, con chi lo ha mandato, il contenuto e il timeStamp -->
      <p v-if="response.message.text !== 'null' || response.message.photo !== ''">
        {{ response.user.username }}
      </p>
      <p v-if="response.message.text !== 'null'">
        {{ response.message.text }}
      </p>
      <p></p>
      <!-- Mostra la foto contenuta nel messaggio nel caso in cui il messaggio contiene una foto -->
      <img class="msg_photo" v-if="response.message.photo !== ''"
        :src="`data:image/jpg;base64,${response.message.photo}`" alt="Message Photo">
      <p v-if="response.message.text !== 'null' || response.message.photo !== ''">
        {{ response.timeMsg }}
      </p>
      <div class="btn-group me-2">
        <!-- Pulsante per inoltrare il messaggio in un'altra conversazione -->
        <button type="button" class="btn btn-sm btn-outline-secondary"
          @click="handleSearchModalToggle(response.message)">
          Forward Message
        </button>
        <!-- Pulsante per eliminare il messaggio dalla conversazione -->
        <button type="button" class="btn btn-sm btn-outline-secondary"
          @click="deleteMessage(response.message.messageId)">
          Delete message
        </button>
        <!-- Pulsante per commentare il messaggio -->
        <button type="button" class="btn btn-sm btn-outline-secondary"
          @click="handleCommentModalToggle(response.message, response.comments)">
          Comment message
        </button>
      </div>
      <hr v-if="response.message.text !== 'null' || response.message.photo !== ''">
    </div>
    <!-- Input per invaire un messaggio testuale -->
    <div class="input-group">
      <input type="text" class="form-control" v-model="text" placeholder="Type your message here">
      <button class="btn btn-outline-primary" @click="check">Send</button>
    </div>

    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
  </div>
</template>

<style>
/* Stile utilizzato per visualizzare l'immagine profilo dell'utente o del gruppo con cui si sta conversando */
.top-profile-container {
  width: auto;
  height: auto;

  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: space-between;
}

/* Stile utilizzato per visualizzare la foto di un messaggio */
.msg_photo {
  width: 25%;
  height: 25%;
}

/* Stile utilizzato nel caso in cui la conversazione è con un gruppo */
.clickable {
  cursor: pointer;
  color: blue;
  text-decoration: underline;
}

.clickable:hover {
  color: darkblue;
}
</style>
