export class BasketItem {
    name: string;
    lastModified: Date;

    constructor(name: string) {
        this.name = name;
        this.lastModified = new Date();
    }

    touch() {
        this.lastModified = new Date();
    }
}