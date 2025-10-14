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

    syncItemAddedToBuy(basketItem: BasketItem) {
        this.sendMessage("itemAddedToBuy", basketItem);
    }

    syncItemBought(basketItem: BasketItem) {
        this.sendMessage("itemBought", basketItem);
    }

    connect(): void {
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

    onMessageReceived(message: WebSocketMessage): void {
        // TODO: The next improvement is should be here.
        // around the methods in basketItemManager - rethink it's interface
        /**
         * Idea: just one method: update, and set the latest state.
         * This should be enough. The backend then can deal with conflicts, if any.
         *
         */
        switch (message.method) {
            case 'itemAddedToBuy':
                this.basketItemManager.addNewItem(message.payload.name);
                break;
            case 'itemBought':
                this.basketItemManager.addToBoughtBasket(message.payload);
                break;
        }
    }

    sendMessage(method: string, payload: any): void {
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
    payload: BasketItem; // 'any' is used here for flexibility, but you can define more specific types.
}