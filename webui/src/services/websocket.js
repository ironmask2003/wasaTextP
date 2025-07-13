class WebSocketService {
    constructor() {
        this.socket = null;
    }

    connect() {
        this.socket = new WebSocket(`ws://localhost:3000/ws?username=${sessionStorage.getItem('username')}`);

        this.socket.onopen = () => {
            console.log("WebSocket connesso");
        };

        this.socket.onmessage = (event) => {
            console.log("Messaggio ricevuto:", event.data);
            // Inoltra ogni messaggio alla pagina principale
            window.parent.postMessage({
                type: "ws-message",
                data: event.data
            }, "*");
        };
    }

    disconnect() {
        if (this.socket) {
            this.socket.close();
            console.log("WebSocket disconnesso");
        }
    }
}

export default new WebSocketService();