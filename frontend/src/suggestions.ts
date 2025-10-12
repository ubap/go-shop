export class BasketItem {
    name: string;
    lastModified: Date;
    htmlElement : HTMLLIElement;

    constructor(name: string, htmlElement: HTMLLIElement) {
        this.name = name;
        this.htmlElement = htmlElement;
        this.lastModified = new Date();
    }

    touch() {
        this.lastModified = new Date();
    }
}

/**
 * 1. item that  name starts with the query, if multiple then return the most recently modified
 * 2. Any word in the item name starts with the query, if multiple then return the most recently modified
 * 3. Any item where the query can be found in the name, if multiple then the most recently modified
 *
 */
// TODO: This needs more effective data structures
export function suggest(query: string, basket: BasketItem[]): BasketItem {
    let item = basket
        .filter((item: BasketItem) => item.name.toLowerCase().startsWith(query.toLowerCase()))
        .sort((a, b) => b.lastModified.getTime() - a.lastModified.getTime())
        [0];

    return item;
}