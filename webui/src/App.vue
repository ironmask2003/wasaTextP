<script setup>
import { RouterLink, RouterView } from 'vue-router'
import Modal from './components/Modal.vue'
</script>
<script>
export default {
  data() {
    return {
      searchModalIsVisible: false,
      isLoggedIn: sessionStorage.token ? true : false,
      userID: sessionStorage.userID,
      newUsername: "",
      newProPic: null,
      username: sessionStorage.username,
      photo: sessionStorage.photo,
      updateNameModalIsVisible: false,
      updateProPicIsVisible: false,
      usernameValidation: new RegExp('^\\w{3,16}$'),
    }
  },
  methods: {
    async handleFileChange(event) {
      this.errorMsg = "";
      const file = event.target.files[0];
      if (file.type !== "image/jpeg") {
        this.errorMsg = "File type not supported, only jpg and jpeg are allowed";
        return
      }
      if (file.size > 5242880) {
        this.errorMsg = "File size is too big. Max size is 5MB";
        return
      }
      this.newProPic = file;
    },
    handleUpdateProPicToggle() {
      sessionStorage.photo = this.photo;
      this.updateProPicIsVisible = !this.updateProPicIsVisible;
      this.newProPic = "";
      this.errorMsg = "";
    },  
    handleUpdateNameToggle() {
      sessionStorage.username = this.username;
			this.updateNameModalIsVisible = !this.updateNameModalIsVisible;
      this.newUsername = "";
      this.errorMsg = "";
		},
    async updateProPic() {
      this.errorMsg = "";
      const formData = new FormData();
      formData.append('image', this.newProPic);
      this.$axios.put(`/profiles/${this.userID}/photo`, formData, { headers: { 'Authorization': `${sessionStorage.token}` } })
        .then(response => {
          this.photo = response.data.photo;
          this.handleUpdateProPicToggle();
        })
        .catch(e => {
          this.errorMsg = e.toString
        });
    },
    async updateUsername() {
    if(this.newUsername == this.username){
        this.errorMsg = "You must enter a new username";  
        return
    } 
    if(this.newUsername.length < 3 || this.newUsername.length > 16){
        this.errorMsg = "Invalid username, it must contains min 3 characters and max 16 characters";
        return
    }
    if(!this.usernameValidation.test(this.newUsername)){
        this.errorMsg = "Invalid username, it must contain only letters and numbers";
        return
    }
    try{
        let _ = await this.$axios.put(`/profiles/${sessionStorage.userID}/username`, { username: this.newUsername }, { headers: { 'Authorization': `${sessionStorage.token}` } })
        this.username = this.newUsername;
        this.handleUpdateNameToggle();
    } catch (e) {
        if(e.response.data == "Username already exist\n"){
            this.errorMsg = "This username is already taken. Please try another one.";
        }else{
            this.errorMsg = e.toString();
        }
    }
		},
    handleSearchModalToggle() {
			this.searchModalIsVisible = !this.searchModalIsVisible;
		},
    logout() {
      sessionStorage.clear();
      this.isLoggedIn = false;
      this.$router.push("/");
    },
    handleLoginSuccess(){
      this.isLoggedIn = true;
      this.username = sessionStorage.username;
      this.photo = sessionStorage.photo;
    }
  }
}
</script>

<template>

  <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
    <a class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href="#/">Wasa Text</a>
  </header>

  <div class="container-fluid">
    <div class="row">
      <nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse" v-show="isLoggedIn">
        <div class="position-sticky pt-3 sidebar-sticky">
          <Modal :show="searchModalIsVisible" @close="handleSearchModalToggle" title="search">
						<template v-slot:header>
							<h3>Users</h3>
						</template>
					</Modal>
          <Modal :show="updateNameModalIsVisible" @close="handleUpdateNameToggle" title = "username">
            <template v-slot:header>
                <h3>Update Username</h3>
            </template>
            <template v-slot:body>
                <form class="username-form">
                    <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                    <input type="text" v-model="newUsername" placeholder="New username" />
                    <button type="submit" @click.prevent="updateUsername">Update</button>
                </form>
            </template>
          </Modal>
          <Modal :show="updateProPicIsVisible" @close="handleUpdateProPicToggle" title = "photo">
            <template v-slot:header>
                <h3>Update Profile Picture</h3>
            </template>
            <template v-slot:body>
                <form class="username-form">
                    <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                    <input type="file" ref="file" accept=".jpg,.jpeg" @change="handleFileChange"/>
                    <button type="submit" @click.prevent="updateProPic">Update</button>
                </form>
            </template>
          </Modal>
          <h6
            class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
            <span>General</span>
          </h6>
          <ul class="nav flex-column">
            <li class="nav-item">
              <RouterLink to="/home" class="nav-link m-2">
                <svg class="feather">
                  <use href="/feather-sprite-v4.29.0.svg#home" />
                </svg>
                Home
              </RouterLink>
            </li>
            <li class="nav-item m-2" v-if="isLoggedIn">
							<a class="nav-link" @click="handleSearchModalToggle">
								<svg class="feather">
									<use href="/feather-sprite-v4.29.0.svg#search" />
								</svg> 
							Search
							</a>
						</li>
            <li class="nav-item m-2" v-if="isLoggedIn">
              <a class="nav-link" @click="logout">
                <svg class="feather">
                  <use href="/feather-sprite-v4.29.0.svg#log-out" />
                </svg>
                Logout
              </a>
            </li>
            <li class="nav-item" v-else>
              <RouterLink to="/session" class="nav-link m-2">
                <svg class="feather">
                  <use href="/feather-sprite-v4.29.0.svg#key" />
                </svg>
                Login
              </RouterLink>
            </li>
            <li class="nav-item m-2" v-if="isLoggedIn">
              <button @click="handleUpdateNameToggle">
                <svg class="feather">
                    <use href="/feather-sprite-v4.29.0.svg#edit" />
                </svg>
              </button>
              Set new username
						</li>
            <li class="nav-item m-2" v-if="isLoggedIn">
              <button @click="handleUpdateProPicToggle">
                <svg class="feather">
                    <use href="/feather-sprite-v4.29.0.svg#edit" />
                </svg>
              </button>
              Set new profile picture
						</li>
            <li class="nav-item m-2" v-if="isLoggedIn">
              <img :src="`data:image/jpg;base64,${photo}`" alt="Profile Picture" class="profile-picture"/>
              <span class="username">{{ username }}</span>
            </li>
          </ul>
        </div>
      </nav>

      <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
				<RouterView @login-success="handleLoginSuccess"/>
			</main>
    </div>
  </div>
</template>

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
