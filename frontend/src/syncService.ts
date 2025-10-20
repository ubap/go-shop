import {BasketItemManager} from "./basketItemManager";
import {BasketItem} from "./basketItem";

export class SyncService {
    private socket: WebSocket | null = null;
    private unacknowledgedMessages: Map<string, { message: WebSocketMessage, timeoutId: any | null }> = new Map();

    private basketItemManager: BasketItemManager;

    constructor(basketItemManager: BasketItemManager) {
        this.basketItemManager = basketItemManager;
    }

    start(): SyncService {
        window.addEventListener('online', this.handleConnectionChange);
        window.addEventListener('offline', this.handleConnectionChange);
        this.connect();
        return this;
    }

    private handleConnectionChange = (): void => {
        if (navigator.onLine) {
            console.log("Internet connection restored. Attempting to reconnect immediately.");
            this.connect();
        } else {
            console.log("Internet connection lost.");
        }
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

        this.socket.onopen = () => {
            console.log("WebSocket connection established.");
            // Todo: send queue of unsynced items
            // TODO: remote has to send all items only after receiving the collection of
            //  unsynced items, the collection can be empty, but the remote waits for the
            //  collection. The remote has to 'ack' it.
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

        this.socket.onclose = (event) => {
            console.log("WebSocket connection closed with code:", event.code);
            setTimeout(() => {
                this.connect()
            }, 1000);
        };

        this.socket.onerror = (error) => {
            console.error("WebSocket error:", error);
            setTimeout(() => {
                this.connect()
            }, 1000);
        };

    }

    private onMessageReceived(message: WebSocketMessage): void {
        switch (message.method) {
            case 'itemUpdate':
                for (const item of message.payload) {
                    this.basketItemManager.upsertBasketItemFromNetwork(item);
                }
                break;
            case 'ack':
                // TODO: Handle ACK in separate method
                const unacknowledgedMessage = this.unacknowledgedMessages.get(message.messageId);
                if (unacknowledgedMessage) {
                    clearTimeout(unacknowledgedMessage.timeoutId);
                    this.unacknowledgedMessages.delete(message.messageId);
                }
                break;
        }
    }

    private sendMessage(method: string, basketItem: BasketItem): void {
        if (!this.socket || this.socket.readyState !== WebSocket.OPEN) {
            console.error("WebSocket is not connected.");
            return;
        }

        const messageId = crypto.randomUUID();
        const message: WebSocketMessage = {
            messageId: messageId, method: method, payload: [basketItem]
        };

        const timeoutId = setTimeout(() => {
            console.error(`ACK not received for message ${messageId}. Assuming connection is lost.`);
            this.socket?.close();
        }, 5000);
        this.unacknowledgedMessages.set(messageId, {message: message, timeoutId: timeoutId});

        this.socket.send(JSON.stringify(message));
        console.log("Message sent to server: ", message);
    }
}

interface WebSocketMessage {
    messageId: string;
    method: string;
    payload: BasketItem[];
}