import {BasketItem} from "./basketItem";

export class BasketItemManager {
    private itemNameToBasketItem: Map<string, BasketItem> = new Map();

    private readonly uiElementCreator: (BasketItem: BasketItem) => void;
    private readonly uiAddToBuyBasket: (BasketItem: BasketItem) => void;
    private readonly uiAddToBoughtBasket: (BasketItem: BasketItem) => void;

    constructor(uiElementCreator: (BasketItem: BasketItem) => void,
                uiAddToBuyBasket: (BasketItem: BasketItem) => void,
                uiAddToBoughtBasket: (BasketItem: BasketItem) => void) {
        this.uiElementCreator = uiElementCreator;
        this.uiAddToBuyBasket = uiAddToBuyBasket;
        this.uiAddToBoughtBasket = uiAddToBoughtBasket;
    }

    addItemToBuyBasket(itemText: string) {
        let basketItem = this.getBasketItemByItemName(itemText);
        if (basketItem === undefined) {
            basketItem = this.createBasketItem(itemText)
        }

        if (basketItem.toBuy) {
            this.uiAddToBuyBasket(basketItem);
        } else {
            this.uiAddToBoughtBasket(basketItem);
        }
    }

    addBasketItemToBuyBasket(basketItem: BasketItem): void {
        if (this.getBasketItemByItemName(basketItem.name) === undefined) {
            this.uiElementCreator(basketItem);
            this.itemNameToBasketItem.set(basketItem.name, basketItem);
        }

        if (basketItem.toBuy) {
            this.uiAddToBuyBasket(basketItem);
        } else {
            this.uiAddToBoughtBasket(basketItem);
        }
    }

    getBasketItemByItemName(name: string): BasketItem | undefined {
        return this.itemNameToBasketItem.get(name);
    }

    addToBuyBasket(basketItem: BasketItem): void {
        this.uiAddToBuyBasket(basketItem);
        basketItem.toBuy = true;
        basketItem.lastModified = new Date();
    }

    addToBoughtBasket(basketItem: BasketItem): void {
        this.uiAddToBoughtBasket(basketItem);
        basketItem.toBuy = false;
        basketItem.lastModified = new Date();
    }

    private createBasketItem(name: string): BasketItem {
        let basketItem = new BasketItem(name);
        this.uiElementCreator(basketItem);
        this.itemNameToBasketItem.set(name, basketItem);
        return basketItem;
    }
}