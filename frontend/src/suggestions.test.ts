import {expect, test} from "vitest";
import {BasketItem, suggest} from "./suggestions";

test('basic suggest, 1 word', () => {
    let basket: BasketItem[] = [
        new BasketItem("water", false, new Date("2018-03-08T08:15:16.097Z")),
        new BasketItem("apple", false, new Date("2018-03-08T09:15:16.097Z")),
        new BasketItem("potato", false, new Date("2018-03-08T10:15:16.097Z")),
        new BasketItem("juice", false, new Date("2018-03-08T11:15:16.097Z")),
        new BasketItem("Walnuts", false, new Date("2018-03-08T12:15:16.097Z")),
    ];

    let query = "WA";

    let basketItem = suggest(query, basket);

    expect(basketItem.name).toBe("Walnuts");
});