import {BasketItemManager} from "./basketItemManager";
import {BasketItem} from "./basketItem";

export class SyncService {

    private basketItemManager: BasketItemManager;

    constructor(basketItemManager: BasketItemManager) {
        this.basketItemManager = basketItemManager;
    }

    start(): SyncService {
        return this;
    }

    syncItemAddedToBuy(basketItem: BasketItem) {

    }

    syncItemBought(basketItem: BasketItem) {

    }
}