import {BasketItemManager} from "./basketItemManager";
import {BasketItem} from "./basketItem";
import {generateUUID} from "./uuid";

export enum ConnectionStatus {
    CONNECTING = 'connecting',
    OPEN = 'open',
    RECONNECTING = 'reconnecting',
    OFFLINE = 'offline',
}

export class SyncService {
    private socket: WebSocket | null = null;
    private unacknowledgedMessages: Map<string, { message: WebSocketMessage, timeoutId: any | null }> = new Map();

    private basketItemManager: BasketItemManager;

    private onStatusChange: (status: ConnectionStatus) => void;

    constructor(basketItemManager: BasketItemManager, onStatusChange: (status: ConnectionStatus) => void) {
        this.basketItemManager = basketItemManager;
        this.onStatusChange = onStatusChange;
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
            this.onStatusChange(ConnectionStatus.RECONNECTING);
            this.connect();
        } else {
            console.log("Internet connection lost.");
            this.onStatusChange(ConnectionStatus.OFFLINE);
        }
    }

    syncItemUpdate(basketItem: BasketItem) {
        this.sendMessage("itemUpdate", [basketItem]);
    }

    private connect(): void {
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
            console.log("WebSocket is already connected.");
            this.onStatusChange(ConnectionStatus.OPEN);
            return;
        }
        this.onStatusChange(ConnectionStatus.CONNECTING);
        this.socket = new WebSocket("/ws");

        this.socket.onopen = () => {
            console.log("WebSocket connection established.");
            this.onStatusChange(ConnectionStatus.OPEN);
            // Todo: send queue of unsynced items
            // TODO: remote has to send all items only after receiving the collection of
            //  unsynced items, the collection can be empty, but the remote waits for the
            //  collection. The remote has to 'ack' it.

            // since ES2015 the values will be returned in the order they were inserted (FIFO)
            this.unacknowledgedMessages.values().forEach(message => {
                clearTimeout(message.timeoutId);
            });
            const messages = this.unacknowledgedMessages
                .values()
                .map(message => message.message)
                .toArray()
            this.sendMessage("unackedMessages", messages);
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
            this.onStatusChange(ConnectionStatus.RECONNECTING);
            setTimeout(() => {
                this.connect()
            }, 1000);
        };

        this.socket.onerror = (error) => {
            console.error("WebSocket error:", error);
        };

    }

    private onMessageReceived(message: WebSocketMessage): void {
        switch (message.method) {
            case 'itemUpdate':
                for (const item of message.payload) {
                    this.basketItemManager.upsertBasketItemFromNetwork(item as BasketItem);
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

    private sendMessage(method: string, payload: PayloadType): void {
        if (!this.socket || this.socket.readyState !== WebSocket.OPEN) {
            console.error("WebSocket is not connected.");
            // TODO add to unacked messages without timeout..
            return;
        }

        const messageId = generateUUID();
        const message: WebSocketMessage = {
            messageId: messageId, method: method, payload: payload
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
    payload: PayloadType;
}

type PayloadType = BasketItem[] | WebSocketMessage[];
