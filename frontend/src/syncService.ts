import {BasketItemManager} from "./basketItemManager";
import {BasketItem} from "./basketItem";

export class SyncService {
    private socket: WebSocket | null = null;

    private basketItemManager: BasketItemManager;

    constructor(basketItemManager: BasketItemManager) {
        this.basketItemManager = basketItemManager;
    }

    start(): SyncService {
        this.connect();
        return this;
    }

    syncItemUpdate(basketItem: BasketItem) {
        this.sendMessage("itemUpdate", basketItem);
    }

    private connect(): void {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            console.log("WebSocket is already connected.");
            return;
        }

        this.socket = new WebSocket("/ws");

        // Event handler for when the connection is established.
        this.socket.onopen = () => {
            console.log("WebSocket connection established.");
        };

        // Event handler for receiving messages from the server.
        this.socket.onmessage = (event) => {
            try {
                const message: WebSocketMessage = JSON.parse(event.data);
                console.log("Message received from server: ", message);
                this.onMessageReceived(message);
            } catch (error) {
                console.error("Error parsing message from server:", error);
            }
        };
    }

    private onMessageReceived(message: WebSocketMessage): void {
        switch (message.method) {
            case 'itemUpdate':
                this.basketItemManager.upsertBasketItemFromNetwork(message.payload)
                break;
        }
    }

    private sendMessage(method: string, payload: BasketItem): void {
        if (!this.socket || this.socket.readyState !== WebSocket.OPEN) {
            console.error("WebSocket is not connected.");
            return;
        }

        const message: WebSocketMessage = {method, payload};
        this.socket.send(JSON.stringify(message));
        console.log("Message sent to server: ", message);
    }
}

interface WebSocketMessage {
    method: string;
    payload: BasketItem;
}