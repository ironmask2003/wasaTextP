<template>
    <div v-if="visible" class="notification">
      {{ message }}
    </div>
  </template>
  
  <script>
  export default {
    props: {
      message: String,
      duration: {
        type: Number,
        default: 1000, // Tempo predefinito: 3 secondi
      },
    },
    data() {
      return {
        visible: false,
      };
    },
    methods: {
      show() {
        this.visible = true;
        setTimeout(() => {
          this.visible = false;
        }, this.duration);
      },
    },
  };
  </script>
  
  <style scoped>
  .notification {
    position: fixed;
    top: 20px;
    right: 20px;
    background: #4caf50;
    color: white;
    padding: 15px;
    border-radius: 5px;
    box-shadow: 0px 4px 6px rgba(0, 0, 0, 0.1);
    transition: opacity 0.5s;
  }
  </style>
  