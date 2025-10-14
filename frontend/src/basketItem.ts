export class BasketItem {
    id: string;
    name: string;
    lastModified: Date;

    constructor(name: string) {
        this.id = crypto.randomUUID();
        this.name = name;
        this.lastModified = new Date();
    }

    touch() {
        this.lastModified = new Date();
    }
}