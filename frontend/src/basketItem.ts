import {generateUUID} from "./uuid";

export class BasketItem {
    id: string;
    name: string;
    lastModified: Date;
    toBuy: boolean; // false = bought, true = still has to be bought

    constructor(name: string) {
        this.id = generateUUID();
        this.name = name;
        this.lastModified = new Date();
        this.toBuy = true;
    }
}