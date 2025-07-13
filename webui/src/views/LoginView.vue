<!-- In questa pagina l'utente effettua il login -->

<script>
export default {
  data() {
    return {
      // Username input dell'utente che si sta loggando
      username: "",
      password: "",
      errorMsg: "",

      // Verifica per il campo username
      usernameValidation: new RegExp('^\\w{3,16}$'),
    }
  },
  emits: ['login-success'],
  methods: {
    // Funzione per effettuare il login
    async doLogin() {
      try {
        // Controlla che l'username sia valido
        if (this.username.length < 3 || this.username.length > 16) throw "Invalid username, it must contains min 3 characters and max 16 characters"
        if (!this.usernameValidation.test(this.username)) throw "Invalid username, it must contain only letters and numbers"

        // Effettua la richietsa di login al server con l'username inserito (se l'username non esiste, verrà creato un nuovo utente)
        let response = await this.$axios.post('/session', {
          username: this.username,
        });

        // Salva i dati dell'utente nella sessionStorage
        sessionStorage.userID = response.data.userId;
        sessionStorage.username = response.data.username;
        sessionStorage.token = response.data.userId;
        sessionStorage.photo = response.data.photo;

        // Reindirizza l'utente alla home
        this.$router.push("/home");
        // Emette l'evento di login avvenuto con successo
        this.$emit('login-success');
      } catch (e) {
        this.errorMsg = e.toString();
        document.getElementsByTagName("input")[0].style.outline = "auto";
        document.getElementsByTagName("input")[0].style.outlineColor = "red";
      };
    }
  },
  mounted() {
    // Se l'utente è già loggato, reindirizza alla home
    if (sessionStorage.token) {
      this.$router.push("/home");
      return;
    }
    // Altrimewnti cancella i dati dell'utente dalla sessionStorage
    sessionStorage.clear();
  },
}

</script>

<template>
  <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
  <div class="login-container">
    <form class="form" @submit.prevent="doLogin">
        <h1 class="title">Welcome to WasaText</h1>
        <div class="flex-column">
          <label>Username </label></div>
          <div class="inputForm">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" viewBox="0 0 32 32" height="20"><g data-name="Layer 3" id="Layer_3"><path d="m30.853 13.87a15 15 0 0 0 -29.729 4.082 15.1 15.1 0 0 0 12.876 12.918 15.6 15.6 0 0 0 2.016.13 14.85 14.85 0 0 0 7.715-2.145 1 1 0 1 0 -1.031-1.711 13.007 13.007 0 1 1 5.458-6.529 2.149 2.149 0 0 1 -4.158-.759v-10.856a1 1 0 0 0 -2 0v1.726a8 8 0 1 0 .2 10.325 4.135 4.135 0 0 0 7.83.274 15.2 15.2 0 0 0 .823-7.455zm-14.853 8.13a6 6 0 1 1 6-6 6.006 6.006 0 0 1 -6 6z"></path></g></svg>
            <input placeholder="Enter your username" class="input" type="text" v-model="username">
          </div>
        
        <div class="flex-column">
          <label>Password </label></div>
          <div class="inputForm">
            <svg xmlns="http://www.w3.org/2000/svg" width="20" viewBox="-64 0 512 512" height="20"><path d="m336 512h-288c-26.453125 0-48-21.523438-48-48v-224c0-26.476562 21.546875-48 48-48h288c26.453125 0 48 21.523438 48 48v224c0 26.476562-21.546875 48-48 48zm-288-288c-8.8125 0-16 7.167969-16 16v224c0 8.832031 7.1875 16 16 16h288c8.8125 0 16-7.167969 16-16v-224c0-8.832031-7.1875-16-16-16zm0 0"></path><path d="m304 224c-8.832031 0-16-7.167969-16-16v-80c0-52.929688-43.070312-96-96-96s-96 43.070312-96 96v80c0 8.832031-7.167969 16-16 16s-16-7.167969-16-16v-80c0-70.59375 57.40625-128 128-128s128 57.40625 128 128v80c0 8.832031-7.167969 16-16 16zm0 0"></path></svg>        
            <input placeholder="Enter your Password" class="input" type="password" v-model="password">
          </div>
        <button class="button-submit" type="submit">Sign In</button>
    </form>
  </div>
</template>


<style>
.title {
    font-size: 1.5rem;
    font-weight: 600;
    margin: 0;
    color: #25d366;
    display: flex;
    justify-content: center;
    align-items: center;
}

.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 50vh;
}
/* From Uiverse.io by JohnnyCSilva */ 
.form {
  display: flex;
  flex-direction: column;
  gap: 10px;
  background-color: #ffffff;
  padding: 30px;
  width: 450px;
  border-radius: 20px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
}

::placeholder {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, 'Open Sans', 'Helvetica Neue', sans-serif;
}

.form button {
  align-self: flex-end;
}

.flex-column > label {
  color: #151717;
  font-weight: 600;
}

.inputForm {
  border: 1.5px solid #ecedec;
  border-radius: 10px;
  height: 50px;
  display: flex;
  align-items: center;
  padding-left: 10px;
  transition: 0.2s ease-in-out;
}

.input {
  margin-left: 10px;
  border-radius: 10px;
  border: none;
  width: 100%;
  height: 100%;
}

.input:focus {
  outline: none;
}

.inputForm:focus-within {
  border: 1.5px solid #25d366;
}

.flex-row {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 10px;
  justify-content: space-between;
}

.flex-row > div > label {
  font-size: 14px;
  color: black;
  font-weight: 400;
}

.span {
  font-size: 14px;
  margin-left: 5px;
  color: #25d366;
  font-weight: 500;
  cursor: pointer;
}

.button-submit {
  margin: 20px 0 10px 0;
  background-color: #25d366;
  border: none;
  color: white;
  font-size: 15px;
  font-weight: 500;
  border-radius: 10px;
  height: 50px;
  width: 100%;
  cursor: pointer;
}   
</style>
