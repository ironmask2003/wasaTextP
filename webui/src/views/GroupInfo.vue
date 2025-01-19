<script setup>
import AddGroup from '../components/ModalAddGroup.vue';
import Modal from '../components/Modal.vue';
</script>

<script>
export default {
  data() {
    return {
      searchVisible: false,
      errorMsg: "",
      usernameValidation: new RegExp('^\\w{0,16}$'),
      filteredUsers: [],
      owner: sessionStorage.username,
      selectedUsers: [], // Lista di utenti selezionati
      addGroupVisible: false,

      newGroupname: "",
      newProPic: null,

      updateNameModalIsVisible: false,
      updateProPicIsVisible: false,

      // Group info
      groupName: localStorage.username,
      groupId: localStorage.userID,
      groupPhoto: localStorage.photo,
      groupMembers: JSON.parse(localStorage.users),
    };
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
      localStorage.photo = this.newProPic;
      this.updateProPicIsVisible = !this.updateProPicIsVisible;
      this.newProPic = "";
      this.errorMsg = "";
    },  
    handleUpdateNameToggle() {
      localStorage.username = this.groupName;
			this.updateNameModalIsVisible = !this.updateNameModalIsVisible;
      this.newUsername = "";
      this.errorMsg = "";
		},
    async updateProPic() {
      this.errorMsg = "";
      const formData = new FormData();
      formData.append('image', this.newProPic);
      this.$axios.put(`/profiles/${sessionStorage.userID}/groups/${this.groupId}/grouphoto`, formData, { headers: { 'Authorization': `${sessionStorage.token}` } })
        .then(response => {
          this.groupPhoto = response.data.photo;
          this.handleUpdateProPicToggle();
        })
        .catch(e => {
          this.errorMsg = e.toString
        });
    },
    async updateGroupName() {
    if(this.newGroupname == this.groupName){
        this.errorMsg = "You must enter a new username";  
        return
    } 
    if(this.newGroupname.length < 3 || this.newGroupname.length > 16){
        this.errorMsg = "Invalid username, it must contains min 3 characters and max 16 characters";
        return
    }
    if(!this.usernameValidation.test(this.newGroupname)){
        this.errorMsg = "Invalid username, it must contain only letters and numbers";
        return
    }
    try{
        let _ = await this.$axios.put(`/profiles/${sessionStorage.userID}/groups/${this.groupId}/groupname`, { groupName: this.newGroupname }, { headers: { 'Authorization': `${sessionStorage.token}` } })
        this.groupName = this.newGroupname;
        this.handleUpdateNameToggle();
    } catch (e) {
        if(e.response.data == "Username already exist\n"){
            this.errorMsg = "This username is already taken. Please try another one.";
        }else{
            this.errorMsg = e.toString();
        }
    }
		},
    leaveGroup(){
      this.errorMsg = "";
      this.$axios.delete(`/profiles/${sessionStorage.userID}/groups/${this.groupId}`, { headers: { 'Authorization': `${sessionStorage.token}` } })
        .then(() => {
          this.$router.push("/home");
        })
        .catch(e => {
          this.errorMsg = e.toString();
        });
    },
    handleAddGroupModalToggle() {
      this.addGroupVisible = !this.addGroupVisible;
    },
    handleSearchModalToggle() {
      this.searchVisible = !this.searchVisible;
    },
  },
}
</script>

<template>
  <div>
    <div
      class="d-flex justify-content-between flex-wrap flex-md-nowrap align-items-center pt-3 pb-2 mb-3 border-bottom">
      <div class="top-profile-container">
        <img :src="`data:image/jpg;base64,${groupPhoto}`">
      </div>
      <Modal :show="updateNameModalIsVisible" @close="handleUpdateNameToggle" title = "username">
        <template v-slot:header>
            <h3>Update Group Name</h3>
        </template>
        <template v-slot:body>
            <form class="username-form">
                <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                <input type="text" v-model="this.newGroupname" placeholder="New group name" />
                <button type="submit" @click.prevent="updateGroupName">Update</button>
            </form>
        </template>
      </Modal>
      <Modal :show="updateProPicIsVisible" @close="handleUpdateProPicToggle" title = "photo">
        <template v-slot:header>
            <h3>Update Group Picture</h3>
        </template>
        <template v-slot:body>
            <form class="username-form">
                <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                <input type="file" ref="file" accept=".jpg,.jpeg" @change="handleFileChange"/>
                <button type="submit" @click.prevent="updateProPic">Update</button>
            </form>
        </template>
      </Modal>
      <AddGroup :show="addGroupVisible" @close="handleAddGroupModalToggle" title="search">
        <template v-slot:header>
          <h3>Add to group</h3>
        </template>
      </AddGroup>
      <h1 class="h1">{{ this.groupName }}</h1>
    <div class="btn-toolbar mb-2 mb-md-0">
      <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
      <div class="btn-group me-2">
        <button type="button" class="btn btn-sm btn-outline-secondary" @click="handleUpdateNameToggle">
          Change group name
        </button>
        <button type="button" class="btn btn-sm btn-outline-secondary" @click="handleUpdateProPicToggle">
          Change group photo
        </button>
      </div>
      <button type="button" class="btn btn-sm btn-outline-primary" @click="handleAddGroupModalToggle">
        Add to group
      </button>
    </div>
    </div>

    <div v-for="user in groupMembers">
      <p class="username">
        <img :src="`data:image/jpg;base64,${user.photo}`" class="profile-picture">
        {{user.username}}
        <button type="button" v-if="user.username == owner" class="btn btn-sm btn-outline-primary" @click="leaveGroup">
          Leave Group
        </button>
      </p>
    </div>

    <ErrorMsg v-if="errormsg" :msg="errormsg"></ErrorMsg>
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
