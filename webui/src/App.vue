<!-- 

Pagina principale dell'applicazione, contiene la navbar con i link alle pagine principali 
dell'applicazione e i modali per l'aggiornamento dell'username e dell'immagine del profilo dell'utente loggato.

L'utente può:
- visualizzare le conversazioni con altri utenti o gruppi andando alla home 
- cercare un utente con cui aprire una nuova conversazione
- aggiornare l'username
- aggiornare l'immagine del profilo
- effettuare il logout
- visualizzare la sua immagine del profilo e il suo username

-->

<script setup>
import { RouterLink, RouterView } from 'vue-router'
import Modal from './components/Modal.vue'
</script>
<script>
import socket from './services/websocket.js'

export default {
  data() {
    return {
      errorMsg: null,

      title: "Wasa Text",
      
      // Utilizzato per mostrare o nascondere il modale di ricerca
      searchModalIsVisible: false,
      // Utilizato per mostarer determinati contenuti della pagina solo se un utente ha effettuato il login
      isLoggedIn: sessionStorage.token ? true : false,
      
      socket: socket, // Socket per la comunicazione in tempo reale
      
      // UserId dell'utente loggato
      userID: sessionStorage.userID,

      // Username dell'utente loggato
      username: sessionStorage.username,
      // Profile picture dell'utente logagto
      photo: sessionStorage.photo,

      // Utilizzato per controllare se l'username inserito dall'utente è valido
      usernameValidation: new RegExp('^\\w{3,16}$'),
    }
  },
  mounted() {

    if (this.isLoggedIn) {
      this.socket.connect(); // Connette il socket al server
    }

    window.addEventListener("message", (event) => {
      if (event.data.type === "ws-message") {
        const messaggio = JSON.parse(event.data.data);

        // Controlla se il messaggio ricevuto è relativo all'aggiornamento dell'username
        if (messaggio.message == "Username Updated" && messaggio.oldUsername == sessionStorage.username) {
          // Estrae il nuovo username dal messaggio
          sessionStorage.username = messaggio.data; // Aggiorna lo username nella sessione
          this.username = messaggio.data; // Aggiorna lo username nella componente
          console.log("Username updated: " + messaggio.data);
        }

        if (messaggio.message == "Photo Updated" && messaggio.username == sessionStorage.username) {
          // Aggiorna l'immagine del profilo nella sessione e nella componente
          sessionStorage.photo = messaggio.data;
          this.photo = sessionStorage.photo;
          console.log("Profile picture updated");
        }
      }
    });
  },
  methods: {
    // Funzione utilizzata per mostrare o nascondere il modale di ricerca di un utente per aprire una nuova conversazione
    handleSearchModalToggle() {
      this.searchModalIsVisible = !this.searchModalIsVisible;
    },
    // Funzione utilizzata per il logout dell'utente
    logout() {
      sessionStorage.clear();
      this.isLoggedIn = false;
      this.$router.push("/");
      this.socket.disconnect(); // Disconnette il socket dal server
    },
    // Funzione utilizzata per il login dell'utente
    handleLoginSuccess() {
      this.isLoggedIn = true;
      this.userID = sessionStorage.userID;
      this.username = sessionStorage.username;
      this.photo = sessionStorage.photo;
      socket.connect(); // Connette il socket al server
      sessionStorage.socket = socket; // Salva il socket nella sessione
    },
    // Funzione utilizzata per andare alla pagina del profilo dell'utente
    goToProfile(){
      this.title = "Profile";
      this.$router.push(`/profile/${this.userID}`);
    }
  }
}
</script>

<template>

  <!-- Header con stile WhatsApp -->
  <header class="profile-header">
      <h1 class="profile-title">{{ title }}</h1>
  </header>

  <div class="container-fluid">
    <div class="row">
      <!-- Modale utilizzato per la ricerca di un utente con cui aprire una conversazione -->
      <Modal :show="searchModalIsVisible" @close="handleSearchModalToggle" title="search">
        <template v-slot:header>
          <h1 class="profile-title">Users</h1>
        </template>
      </Modal>
      <!-- Navigation bar -->
      <div id="nav-bar">
        <div id="nav-header"><a id="nav-title" href="#/" target="_blank">M<i class="fab fa-codepen"></i>NU</a>
          <label for="nav-toggle"><span id="nav-toggle-burger"></span></label>
        </div>
        <div></div>
        <div></div>
        <div id="nav-content">
          <div v-if="isLoggedIn" class="nav-button" @click="title = 'Wasa Text'; $router.push('/home')">
              <svg class="feather">
                <use href="/feather-sprite-v4.29.0.svg#home" />
              </svg>
              <span class="subtitle">Home</span>
          </div>
          <hr/>
          <div v-if="isLoggedIn" class="nav-button" @click="handleSearchModalToggle">
            <svg class="feather">
                <use href="/feather-sprite-v4.29.0.svg#search" />
            </svg>
            <span class="subtitle">Search</span></div>
          <div v-if="isLoggedIn" class="nav-button" @click="logout">
            <svg class="feather">
                <use href="/feather-sprite-v4.29.0.svg#log-out" />
            </svg>
            <span class="subtitle">Logout</span></div>
          <div v-else class="nav-button" @click="$router.push('/session')">
            <svg class="feather">
                <use href="/feather-sprite-v4.29.0.svg#log-in"/>
            </svg>  
            <span class="subtitle">Login</span>
          </div>
        </div>
        <div id="nav-footer" v-if="isLoggedIn">
          <div id="nav-footer-heading" @click="goToProfile">
            <div id="nav-footer-avatar"><img :src="`data:image/jpg;base64,${photo}`" class="profile-picture"/></div>
            <div id="nav-footer-titlebox"><a id="nav-footer-title" class="subtitle">{{username}}</a></div>
            <label for="nav-footer-toggle"><i class="fas fa-caret-up"></i></label>
          </div>
        </div>
      </div>

      <!-- Contenuto principale della pagina -->
      <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
        <RouterView @login-success="handleLoginSuccess"/>
      </main>
    </div>
  </div>
</template>


<!-- Stili per l'immagine del profilo e dell'username dell'utente loggato -->
<style>

:root {
  --navbar-width: 256px;
  --navbar-width-min: 80px;
  --navbar-dark-primary: #075E54;
  --navbar-dark-secondary: #128C7E;
  --navbar-light-primary: #DCF8C6;
}

#nav-toggle:checked ~ #nav-header {
  width: calc(var(--navbar-width-min) - 16px);
}
#nav-toggle:checked ~ #nav-content, #nav-toggle:checked ~ #nav-footer {
  width: var(--navbar-width-min);
}
#nav-toggle:checked ~ #nav-header #nav-title {
  opacity: 0;
  pointer-events: none;
  transition: opacity 0.1s;
}
#nav-toggle:checked ~ #nav-header label[for=nav-toggle] {
  left: calc(50% - 8px);
  transform: translate(-50%);
}
#nav-toggle:checked ~ #nav-header #nav-toggle-burger {
  background: var(--navbar-light-primary);
}
#nav-toggle:checked ~ #nav-header #nav-toggle-burger:before, #nav-toggle:checked ~ #nav-header #nav-toggle-burger::after {
  width: 16px;
  background: var(--navbar-light-secondary);
  transform: translate(0, 0) rotate(0deg);
}
#nav-toggle:checked ~ #nav-content .nav-button span {
  opacity: 0;
  transition: opacity 0.1s;
}
#nav-toggle:checked ~ #nav-content .nav-button .fas {
  min-width: calc(100% - 16px);
}
#nav-toggle:checked ~ #nav-footer #nav-footer-avatar {
  margin-left: 0;
  left: 50%;
  transform: translate(-50%);
}
#nav-toggle:checked ~ #nav-footer #nav-footer-titlebox, #nav-toggle:checked ~ #nav-footer label[for=nav-footer-toggle] {
  opacity: 0;
  transition: opacity 0.1s;
  pointer-events: none;
}

#nav-bar {
  position: absolute;
  top: 1vw;
  height: calc(100% - 1vw);
  width: 28vh;
  background: var(--navbar-dark-primary);
  display: flex;
  flex-direction: column;
  color: var(--navbar-light-primary);
  font-family: Verdana, Geneva, Tahoma, sans-serif;
  overflow: hidden;
  user-select: none;
  border-radius: 0 60px 20px 0;
}
#nav-bar hr {
  margin: 0;
  position: relative;
  top: 10px;
  left: -16px;
  width: calc(100% - 32px);
  border: none;
  border-top: solid 1px var(--navbar-dark-secondary);
}
#nav-bar a {
  color: inherit;
  text-decoration: inherit;
}
#nav-bar input[type=checkbox] {
  display: none;
}

#nav-header {
  position: relative;
  width: var(--navbar-width);
  width: calc(var(--navbar-width) - 16px);
  background: var(--navbar-dark-primary);
  border-radius: 16px;
  z-index: 2;
  display: flex;
  align-items: center;
  transition: width 0.2s;
}
#nav-header hr {
  position: absolute;
  bottom: 0;
}

#nav-title {
  font-size: 1.5rem;
  transition: opacity 1s;
}

label[for=nav-toggle] {
  position: absolute;
  right: 0;
  width: 3rem;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

#nav-toggle-burger {
  position: relative;
  width: 16px;
  height: 2px;
  background: var(--navbar-dark-primary);
  border-radius: 99px;
  transition: background 0.2s;
}
#nav-toggle-burger:before, #nav-toggle-burger:after {
  content: "";
  position: absolute;
  top: -6px;
  width: 10px;
  height: 2px;
  background: var(--navbar-light-primary);
  border-radius: 99px;
  transform: translate(2px, 8px) rotate(30deg);
  transition: 0.2s;
}
#nav-toggle-burger:after {
  top: 6px;
  transform: translate(2px, -8px) rotate(-30deg);
}

#nav-content {
  margin: -16px 0;
  padding: 16px 0;
  position: relative;
  flex: 1;
  width: var(--navbar-width);
  background: var(--navbar-dark-primary);
  box-shadow: 0 0 0 16px var(--navbar-dark-primary);
  direction: rtl;
  overflow-x: hidden;
  transition: width 0.2s;
}
#nav-content::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}
#nav-content::-webkit-scrollbar-thumb {
  border-radius: 99px;
  background-color: #D62929;
}
#nav-content::-webkit-scrollbar-button {
  height: 16px;
}

#nav-content-highlight {
  position: absolute;
  left: 16px;
  top: -70px;
  width: calc(100% - 16px);
  height: 54px;
  background: var(--background);
  background-attachment: fixed;
  border-radius: 16px 0 0 16px;
  transition: top 0.2s;
}
#nav-content-highlight:before, #nav-content-highlight:after {
  content: "";
  position: absolute;
  right: 0;
  bottom: 100%;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  box-shadow: 16px 16px var(--background);
}
#nav-content-highlight:after {
  top: 100%;
  box-shadow: 16px -16px var(--background);
}

.nav-button {
  position: relative;
  height: 54px;
  align-items: center;
  color: var(--navbar-light-secondary);
  direction: ltr;
  cursor: pointer;
  z-index: 1;
  top: 30px;
  left: 16px;
  transition: color 0.2s;
}
.nav-button span {
  transition: opacity 1s;
}
.nav-button .fas {
  transition: min-width 0.2s;
}
.nav-button:nth-of-type(1):hover {
  color: var(--navbar-dark-primary);
}
.nav-button:nth-of-type(1):hover ~ #nav-content-highlight {
  top: 16px;
}
.nav-button:nth-of-type(2):hover {
  color: var(--navbar-dark-primary);
}
.nav-button:nth-of-type(2):hover ~ #nav-content-highlight {
  top: 70px;
}
.nav-button:nth-of-type(3):hover {
  color: var(--navbar-dark-primary);
}
.nav-button:nth-of-type(3):hover ~ #nav-content-highlight {
  top: 124px;
}
.nav-button:nth-of-type(4):hover {
  color: var(--navbar-dark-primary);
}
.nav-button:nth-of-type(4):hover ~ #nav-content-highlight {
  top: 178px;
}
.nav-button:nth-of-type(5):hover {
  color: var(--navbar-dark-primary);
}
.nav-button:nth-of-type(5):hover ~ #nav-content-highlight {
  top: 232px;
}
.nav-button:nth-of-type(6):hover {
  color: var(--navbar-dark-primary);
}
.nav-button:nth-of-type(6):hover ~ #nav-content-highlight {
  top: 286px;
}
.nav-button:nth-of-type(7):hover {
  color: var(--navbar-dark-primary);
}
.nav-button:nth-of-type(7):hover ~ #nav-content-highlight {
  top: 340px;
}
.nav-button:nth-of-type(8):hover {
  color: var(--navbar-dark-primary);
}
.nav-button:nth-of-type(8):hover ~ #nav-content-highlight {
  top: 394px;
}

#nav-bar .fas {
  min-width: 3rem;
  text-align: center;
}

#nav-footer {
  width: var(--navbar-width);
  height: 54px;
  background: var(--navbar-dark-secondary);
  border-radius: 16px;
  display: flex;
  flex-direction: column;
  z-index: 2;
  position: relative;
  top: -10px;
  transition: width 0.2s, height 0.2s;
}

#nav-footer-heading {
  position: relative;
  width: 100%;
  height: 54px;
  display: flex;
  align-items: center;
}

#nav-footer-heading:hover {
  cursor: pointer;
}

#nav-footer-avatar {
  position: relative;
  margin: 11px 0 11px 16px;
  left: 0;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  overflow: hidden;
  transform: translate(0);
  transition: 0.2s;
}
#nav-footer-avatar img {
  height: 100%;
}

#nav-footer-titlebox {
  position: relative;
  margin-left: 16px;
  width: 10px;
  display: flex;
  flex-direction: column;
  transition: opacity 1s;
}

#nav-footer-subtitle {
  color: var(--navbar-light-secondary);
  font-size: 0.6rem;
}

#nav-toggle:not(:checked) ~ #nav-footer-toggle:checked + #nav-footer {
  height: 30%;
  min-height: 54px;
}
#nav-toggle:not(:checked) ~ #nav-footer-toggle:checked + #nav-footer label[for=nav-footer-toggle] {
  transform: rotate(180deg);
}

label[for=nav-footer-toggle] {
  position: absolute;
  right: 0;
  width: 3rem;
  height: 100%;
  display: flex;
  align-items: center;
  cursor: pointer;
  transition: transform 0.2s, opacity 0.2s;
}

#nav-footer-content {
  margin: 0 16px 16px 16px;
  border-top: solid 1px var(--navbar-light-secondary);
  padding: 16px 0;
  color: var(--navbar-light-secondary);
  font-size: 0.8rem;
  overflow: auto;
}
#nav-footer-content::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}
#nav-footer-content::-webkit-scrollbar-thumb {
  border-radius: 99px;
  background-color: #D62929;
}

.subtitle {
    position: relative;
    font-weight: 600;
    margin: 0;
    color: white;
    left: 10px;
    top: 2px;
}

.profile-title {
  font-size: 1.5rem;
  font-weight: 600;
  margin: 0;
  color: white;
}

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

.profile-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 1000;
  background: #25d366;
  color: white;
  padding: 1rem;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.container-fluid {
  margin-top: 40px; /* Aggiungi spazio per l'header fisso */
  padding-top: 1rem;
}
</style>