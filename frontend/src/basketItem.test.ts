import {expect, test} from "vitest";
import {BasketItem} from "./basketItem";

test('id should not repeat', () => {
    const water = new BasketItem("Water")
    const bread = new BasketItem("Bread")
    expect(water.id).not.toBe(bread.id);
});
