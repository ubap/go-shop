import {BasketItem} from "./basketItem";

export class BasketItemManager {
    private toBuyBasket: BasketItem[] = [];
    private boughtBasket: BasketItem[] = [];

    private itemNameToBasketItem: Map<string, BasketItem> = new Map();

    private uiElementCreator: (BasketItem: BasketItem) => void;
    private uiAddToBuyBasket: (BasketItem: BasketItem) => void;
    private uiAddToBoughtBasket: (BasketItem: BasketItem) => void;

    constructor(uiElementCreator: (BasketItem: BasketItem) => void,
                uiAddToBuyBasket: (BasketItem: BasketItem) => void,
                uiAddToBoughtBasket: (BasketItem: BasketItem) => void) {
        this.uiElementCreator = uiElementCreator;
        this.uiAddToBuyBasket = uiAddToBuyBasket;
        this.uiAddToBoughtBasket = uiAddToBoughtBasket;
    }

    addNewItem(itemText: string) {
        let basketItem = this.getBasketItem(itemText);
        if (basketItem === undefined) {
            basketItem = this.createBasketItem(itemText)
        }
        this.addToBuyBasket(basketItem);
        basketItem.touch();
    }

    getBasketItem(name: string): BasketItem | undefined {
        return this.itemNameToBasketItem.get(name);
    }

    addToBuyBasket(basketItem: BasketItem): void {
        this.boughtBasket = this.boughtBasket
            .filter(bought => bought.name !== basketItem.name);
        this.toBuyBasket.push(basketItem);
        this.uiAddToBuyBasket(basketItem);
        basketItem.touch();
    }

    addToBoughtBasket(basketItem: BasketItem): void {
        this.toBuyBasket = this.toBuyBasket
            .filter(toBuy => toBuy.name !== basketItem.name);
        this.boughtBasket.push(basketItem);
        this.uiAddToBoughtBasket(basketItem);
        basketItem.touch();
    }

    private createBasketItem(name: string): BasketItem {
        let basketItem = new BasketItem(name);
        this.uiElementCreator(basketItem);
        this.itemNameToBasketItem.set(name, basketItem);
        return basketItem;
    }
}