import {BasketItem} from "./suggestions";

export class BasketItemManager {
    toBuyBasket: BasketItem[] = [];
    boughtBasket: BasketItem[] = [];

    itemNameToBasketItem: Map<string, BasketItem> = new Map();
    uiElementFactory: (itemName: string) => HTMLLIElement;

    uiAddToBuyBasket: (BasketItem: BasketItem) => void;
    uiAddToBoughtBasket: (BasketItem: BasketItem) => void;

    constructor(uiElementFactory: (itemName: string) => HTMLLIElement,
                uiAddToBuyBasket: (BasketItem: BasketItem) => void,
                uiAddToBoughtBasket: (BasketItem: BasketItem) => void) {
        this.uiElementFactory = uiElementFactory;
        this.uiAddToBuyBasket = uiAddToBuyBasket;
        this.uiAddToBoughtBasket = uiAddToBoughtBasket;
    }

    createBasketItem(name: string): BasketItem {
        let htmlLIElement = this.uiElementFactory(name);
        let basketItem = new BasketItem(name, htmlLIElement);
        this.itemNameToBasketItem.set(name, basketItem);

        return basketItem;
    }

    getBasketItem(name: string): BasketItem | undefined{
        return this.itemNameToBasketItem.get(name);
    }

    addToBuyBasket(basketItem: BasketItem): void {
        this.boughtBasket = this.boughtBasket.filter(bought => bought.name !== basketItem.name);
        this.toBuyBasket.push(basketItem);
        this.uiAddToBuyBasket(basketItem);
        basketItem.touch();
    }

    addToBoughtBasket(basketItem: BasketItem): void {
        this.toBuyBasket = this.toBuyBasket.filter(toBuy => toBuy.name !== basketItem.name);
        this.boughtBasket.push(basketItem);
        this.uiAddToBoughtBasket(basketItem);
        basketItem.touch();
    }
}