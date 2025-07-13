<!-- Page used for visualize profile info -->

<script>
import Modal from '../components/Modal.vue'

export default {
    data: function () {
        return {
            window: window,
            // Utilizzato per mostrare o nascondere il modal di aggiornamento dell'username
            updateNameModalIsVisible: false,
            updateProPicIsVisible: false,
            
            // Nuovo username da aggiornare
            newUsername: "",
            
            // Username dell'utente loggato
            username: sessionStorage.username,
            photo: sessionStorage.photo,
            tempPhoto: sessionStorage.photo,

            // Utilizzato per controllare se l'username inserito dall'utente è valido
            usernameValidation: new RegExp('^\\w{3,16}$'),

            // Messaggio di errore
            errorMsg: "",
        }
    },
    components: { Modal },
    methods:{
        // Funzione utilizzata per mostare o nascondere il modale di aggiornamento dell'username
        handleUpdateNameToggle() {
            this.updateNameModalIsVisible = !this.updateNameModalIsVisible; // Nasconde o mostra il modale
            this.newUsername = "";
            this.errorMsg = "";
        },
        // Funzione utilizzata per mostrare o nascondere il modale di aggiornamento immagine del profilo
        handleUpdateProPicToggle() {
            this.updateProPicIsVisible = !this.updateProPicIsVisible; // Nasconde o mostra il modale
            this.newProPic = "";
            this.errorMsg = "";
        },
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
            
            // Convert file to base64
            const reader = new FileReader();
            reader.onload = (e) => {
                const base64 = e.target.result.split(',')[1]; // Remove data:image/jpeg;base64, prefix
                this.tempPhoto = base64;
            };
            reader.readAsDataURL(file);
        },
        // Funzione utilizzata per l'aggiornamento della foto profilo dell'utente
        async updateProPic() {
            this.errorMsg = "";

            // Crea un nuovo oggetto FormData e vi aggiunge l'immagine inserita dall'utente da mandare al server
            const formData = new FormData();
            formData.append('image', this.newProPic);

            // Effettua una richiesta PUT al server per l'aggiornamento della foto profilo
            this.$axios.put(`/profiles/${sessionStorage.userID}/photo`, formData, { headers: { 'Authorization': `${sessionStorage.token}` } })
            .then(response => {
                this.photo = response.data.photo; // Assegna la nuova immagine del profilo alla variabile photo per l'aggiornamento della pagina
                this.handleUpdateProPicToggle(); // Nasconde il modale di aggiornamento dell'immagine del profilo e aggiorna l'immagine del profilo della sessione
            })
            .catch(e => {
                this.errorMsg = e.toString();
            });
        },
        // Funzione utilizzata per l'aggiornamento dell'username dell'utente
        async updateUsername() {

            // COntrolla se l'username inserito dall'utente è uguale a quello attuale
            if (this.newUsername == this.username) {
                this.errorMsg = "You must enter a new username";
                return
            }
            // Controlla se l'username inserito dall'utente ha una lunghezza valida
            if (this.newUsername.length < 3 || this.newUsername.length > 16) {
                this.errorMsg = "Invalid username, it must contains min 3 characters and max 16 characters";
                return
            }
            // Conotrolla se l'username inserito dall'utente è valido
            if (!this.usernameValidation.test(this.newUsername)) {
                this.errorMsg = "Invalid username, it must contain only letters and numbers";
                return
            }
            try {
                
                // Effettua una richiesta PUT al server per l'aggiornamento dell'username
                let _ = await this.$axios.put(`/profiles/${sessionStorage.userID}/username`, 
                    { username: this.newUsername }, 
                    { headers: { 'Authorization': `${sessionStorage.token}` } }
                );
                
                this.username = this.newUsername
                this.handleUpdateNameToggle(); // Nasconde il modale di aggiornamento dell'username
            } catch (e) {
                if (e.response && e.response.data == "Username already exist\n") {
                    this.errorMsg = "This username is already taken. Please try another one.";
                } else {
                    this.errorMsg = e.toString();
                }
            }
        },
    }
}
</script>

<template>
    <!-- Sezione principale del profilo -->
    <div class="profile-main">
        <!-- Foto profilo centrale -->
        <div class="profile-photo-section">
            <div class="profile-photo-wrapper" @click="handleUpdateProPicToggle">
                <img 
                    :src="photo ? `data:image/jpg;base64,${photo}` : '/storage/default_profile_photo.jpg'" 
                    :alt="`${username} profile picture`" 
                    class="profile-photo"
                />
                <div class="photo-edit-overlay">
                    <i class="fas fa-camera"></i>
                </div>
            </div>
        </div>

        <!-- Informazioni utente -->
        <div class="profile-info-section">
            <!-- Username -->
            <div class="info-item" @click="handleUpdateNameToggle">
                <div class="info-content">
                    <div class="info-label">Name</div>
                    <div class="info-value">{{ username || 'No username set' }}</div>
                </div>
                <div class="info-action">
                    <i class="fas fa-edit"></i>
                </div>
            </div>
        </div>
    </div>

    <!-- Modale utilizzato per l'aggiornamento della foto profilo dell'utente -->
    <Modal :show="updateProPicIsVisible" @close="handleUpdateProPicToggle">
        <template v-slot:header>
            <h3 class="profile-title">Select new picture</h3>
        </template>
        <template v-slot:body>
            <!-- Input in cui viene inserita la nuova immagine del profilo dell'utente -->
            <form class="username-form">
                <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                <div class="profile-photo-wrapper" @click="handleUpdateProPicToggle">
                        <img 
                                :src="tempPhoto ? `data:image/jpg;base64,${tempPhoto}` : '/storage/default_profile_photo.jpg'" 
                                :alt="`${username} profile picture`" 
                                class="profile-photo"
                        />
                </div>
                <div class="button-container">
                        <label class="attachment-button">
                                📎
                                <input style="display: none;" type="file" @change="handleFileChange" accept="image/*" class="image-upload" />
                        </label>
                        <button 
                                        type="submit" 
                                        @click.prevent="updateProPic"
                                        class="btn btn-primary btn-update"
                                >
                                        Update
                        </button>
                </div>
            </form>
        </template>
    </Modal>

    <!-- Modale per aggiornamento username -->
    <Modal :show="updateNameModalIsVisible" @close="handleUpdateNameToggle" title="Update Username">
        <template v-slot:header>
            <h3 class="profile-title">Update Username</h3>
        </template>
        <template v-slot:body>
            <form class="username-form">
                <ErrorMsg v-if="errorMsg" :msg="errorMsg"></ErrorMsg>
                <div class="form-group">
                    <input 
                        type="text" 
                        v-model="newUsername" 
                        placeholder="Enter new username" 
                        class="form-control"
                    />
                </div>
                <button 
                    type="submit" 
                    @click.prevent="updateUsername"
                    class="btn btn-primary btn-update"
                >
                    Update
                </button>
            </form>
        </template>
    </Modal>
</template>

<style scoped>

.profile-title {
    font-size: 1.5rem;
    font-weight: 600;
    margin: 0;
    color: white;
}

.profile-main {
    background: white;
    margin: 20px;
    border-radius: 10px;
    overflow: hidden;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.profile-photo-section {
    background: linear-gradient(135deg, #25d366 0%, #128c7e 100%);
    padding: 40px 20px;
    display: flex;
    justify-content: center;
    align-items: center;
}

.profile-photo-wrapper {
    position: relative;
    cursor: pointer;
}

.profile-photo {
    width: 140px;
    height: 140px;
    border-radius: 50%;
    object-fit: cover;
    border: 4px solid white;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
    transition: all 0.3s ease;
}

.profile-photo:hover {
    transform: scale(1.05);
    box-shadow: 0 6px 20px rgba(0, 0, 0, 0.3);
}

.photo-edit-overlay {
    position: absolute;
    bottom: 8px;
    right: 8px;
    background: #25d366;
    color: white;
    width: 40px;
    height: 40px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 3px solid white;
    font-size: 16px;
    cursor: pointer;
    transition: all 0.2s ease;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.button-container {
    display: flex;
    align-items: center;
    gap: 16px;
    justify-content: center;
    margin-top: 20px;
}

.attachment-button {
    width: 50px;
    height: 50px;
    border: none;
    border-radius: 50%;
    background-color: #f0f0f0;
    font-size: 24px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: background-color 0.2s, box-shadow 0.2s;
    flex-shrink: 0; /* Non si riduce */
}

.btn-update {
    background: linear-gradient(135deg, #25d366 0%, #128c7e 100%);
    border: none;
    color: white;
    padding: 14px 32px;
    border-radius: 8px;
    font-size: 16px;
    font-weight: 600;
    cursor: pointer;
    flex: 1; /* Occupa lo spazio rimanente */
    transition: all 0.2s ease;
    box-shadow: 0 2px 8px rgba(37, 211, 102, 0.3);
    text-transform: uppercase;
    letter-spacing: 0.5px;
    min-height: 50px; /* Stessa altezza del attachment-button */
}

.attachment-button:hover {
background-color: #e0e0e0;
}

.attachment-button:active {
background-color: #d0d0d0;
box-shadow: inset 0 0 5px rgba(0,0,0,0.2);
}

.photo-edit-overlay:hover {
    background: #1da851;
    transform: scale(1.1);
}

.profile-info-section {
    padding: 0;
}

.info-item {
    padding: 20px 24px;
    border-bottom: 1px solid #e5e5e5;
    display: flex;
    justify-content: space-between;
    align-items: center;
    cursor: pointer;
    transition: all 0.2s ease;
    min-height: 70px;
}

.info-item:hover {
    background: #f5f5f5;
    transform: translateX(2px);
}

.info-item:last-child {
    border-bottom: none;
}

.info-content {
    flex: 1;
}

.info-label {
    font-size: 14px;
    color: #667781;
    margin-bottom: 6px;
    font-weight: 500;
    text-transform: uppercase;
    letter-spacing: 0.5px;
}

.info-value {
    font-size: 16px;
    color: #111b21;
    font-weight: 400;
    line-height: 1.4;
    word-break: break-word;
}

.info-secondary {
    color: #667781;
    font-style: italic;
}

.info-action {
    color: #25d366;
    font-size: 18px;
    margin-left: 15px;
    transition: all 0.2s ease;
}

.info-action:hover {
    transform: scale(1.2);
    color: #1da851;
}

.username-form {
    padding: 24px 0;
}

.form-group {
    margin-bottom: 24px;
}

.form-control {
    width: 100%;
    padding: 14px 18px;
    border: 2px solid #e5e5e5;
    border-radius: 8px;
    font-size: 16px;
    transition: all 0.2s ease;
    background: #fafafa;
}

.form-control:focus {
    outline: none;
    border-color: #25d366;
    background: white;
    box-shadow: 0 0 0 3px rgba(37, 211, 102, 0.1);
}

.btn-update:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(37, 211, 102, 0.4);
}

.btn-update:active {
    transform: translateY(0);
}

/* === RESPONSIVE BREAKPOINTS === */

/* Large Desktop: 1200px and up */
@media (min-width: 1200px) {
    .whatsapp-profile-container {
        max-width: 800px;
        margin: 0 auto;
    }
    
    .profile-main {
        margin: 40px auto;
        max-width: 800px;
    }
    
    .profile-photo {
        width: 160px;
        height: 160px;
    }
    
    .photo-edit-overlay {
        width: 45px;
        height: 45px;
        font-size: 18px;
    }
}

/* Desktop: 992px to 1199px */
@media (min-width: 992px) and (max-width: 1199px) {
    .whatsapp-profile-container {
        max-width: 700px;
    }
    
    .profile-main {
        margin: 30px;
    }
    
    .profile-photo {
        width: 150px;
        height: 150px;
    }
}

/* Tablet: 768px to 991px */
@media (min-width: 768px) and (max-width: 991px) {
    .whatsapp-profile-container {
        max-width: 100%;
        padding: 0 10px;
    }
    
    .profile-main {
        margin: 20px auto;
        max-width: 500px;
    }
    
    .profile-photo {
        width: 130px;
        height: 130px;
    }
    
    .info-item {
        padding: 18px 20px;
    }
    
    .profile-title {
        font-size: 1.4rem;
    }
}

/* Mobile Large: 576px to 767px */
@media (min-width: 576px) and (max-width: 767px) {
    .whatsapp-profile-container {
        margin: 0;
        padding: 0;
    }
    
    .profile-main {
        margin: 15px;
        border-radius: 8px;
    }
    
    .profile-photo-section {
        padding: 30px 15px;
    }
    
    .profile-photo {
        width: 110px;
        height: 110px;
    }
    
    .photo-edit-overlay {
        width: 35px;
        height: 35px;
        font-size: 14px;
        bottom: 5px;
        right: 5px;
    }
    
    .info-item {
        padding: 16px 18px;
        min-height: 65px;
    }
    
    .info-value {
        font-size: 15px;
    }
    
    .profile-title {
        font-size: 1.3rem;
    }
    
    .profile-header {
        padding: 16px;
    }
}

/* Mobile Small: up to 575px */
@media (max-width: 575px) {
    .whatsapp-profile-container {
        margin: 0;
        padding: 0;
        background: white;
    }
    
    .profile-main {
        margin: 10px;
        border-radius: 6px;
        box-shadow: 0 1px 4px rgba(0, 0, 0, 0.1);
    }
    
    .profile-photo-section {
        padding: 25px 10px;
    }
    
    .profile-photo {
        width: 100px;
        height: 100px;
        border: 3px solid white;
    }
    
    .photo-edit-overlay {
        width: 32px;
        height: 32px;
        font-size: 12px;
        bottom: 3px;
        right: 3px;
        border: 2px solid white;
    }
    
    .info-item {
        padding: 14px 16px;
        min-height: 60px;
    }
    
    .info-label {
        font-size: 12px;
        margin-bottom: 4px;
    }
    
    .info-value {
        font-size: 14px;
    }
    
    .info-action {
        font-size: 16px;
        margin-left: 10px;
    }
    
    .profile-title {
        font-size: 1.2rem;
    }
    
    .profile-header {
        padding: 14px;
    }
    
    .form-control {
        padding: 12px 14px;
        font-size: 16px; /* Prevent zoom on iOS */
    }
    
    .btn-update {
        padding: 12px 24px;
        font-size: 14px;
    }
    
    .username-form {
        padding: 16px 0;
    }
    
    .form-group {
        margin-bottom: 16px;
    }
}

/* Extra Small Mobile: up to 320px */
@media (max-width: 320px) {
    .profile-main {
        margin: 5px;
    }
    
    .profile-photo {
        width: 90px;
        height: 90px;
    }
    
    .photo-edit-overlay {
        width: 28px;
        height: 28px;
        font-size: 10px;
    }
    
    .info-item {
        padding: 12px 14px;
        min-height: 55px;
    }
    
    .profile-photo-section {
        padding: 20px 8px;
    }
    
    .profile-title {
        font-size: 1.1rem;
    }
}

/* Touch device optimizations */
@media (hover: none) and (pointer: coarse) {
    .info-item:hover {
        background: transparent;
        transform: none;
    }
    
    .info-item:active {
        background: #f0f0f0;
        transform: scale(0.98);
    }
    
    .photo-edit-overlay:hover {
        transform: none;
    }
    
    .photo-edit-overlay:active {
        transform: scale(0.95);
    }
    
    .btn-update:hover {
        transform: none;
        box-shadow: 0 2px 8px rgba(37, 211, 102, 0.3);
    }
    
    .btn-update:active {
        transform: scale(0.98);
    }
}

/* High DPI screens */
@media (-webkit-min-device-pixel-ratio: 2), (min-resolution: 192dpi) {
    .profile-photo {
        image-rendering: -webkit-optimize-contrast;
        image-rendering: crisp-edges;
    }
}

/* Landscape orientation on mobile */
@media (max-width: 767px) and (orientation: landscape) {
    .profile-photo-section {
        padding: 20px 15px;
    }
    
    .profile-photo {
        width: 80px;
        height: 80px;
    }
    
    .photo-edit-overlay {
        width: 25px;
        height: 25px;
        font-size: 10px;
    }
    
    .profile-header {
        padding: 10px;
    }
    
    .profile-title {
        font-size: 1.1rem;
    }
}
</style>