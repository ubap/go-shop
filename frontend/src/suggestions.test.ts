import {expect, test} from "vitest";
import {BasketItem} from "./basketItem";
import {suggest} from "./suggestions";

test('basic suggest, 1 word', () => {
    let basket: BasketItem[] = [
        basketItemLater("apple"),
        basketItemLater("potato"),
        basketItemLater("juice"),
        basketItemLater("Walnuts"),
    ];

    let query = "WA";

    let basketItem = suggest(query, basket);

    expect(basketItem.name).toBe("Walnuts");
});



let mockedTime = new Date().getTime();

function timeTick(): Date {
    mockedTime += 1;
    return new Date(mockedTime);
}

function basketItemLater(text: string) {
    let basketItem = new BasketItem(text);
    basketItem.lastModified = timeTick();
    return basketItem;
}

