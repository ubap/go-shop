import {BasketItem} from "./basketItem";
import {BasketItemManager} from "./basketItemManager";
import {SyncService} from "./syncService";


/**
 * The application does not have a concept of users, only baskets. A basket is secured and accessed solely by its unique URL.
 * Baskets can be shared with multiple people for collaboration by simply sharing this URL.
 *
 * Item names are treated as unique and case-insensitive; for example, "water", "WATER", and "Water" all refer to the same item.
 * However, each basket can store its own preferred casing (e.g., one basket uses "water" while another uses "WATER").
 * This is implemented on the backend with a per-basket item name mapping.
 *
 * The suggestion search should first look for matches within the current basket.
 * If no local matches are found, it should suggest items from a global usage list.
 * The spelling for these global items should be standardized (e.g., "Sparkling water") but should also be easy to change.
 *
 * To prevent duplicate items, the system must first check if an item with the same name already exists in the itemsToBuy or itemsBought lists.
 * If it does, the existing instance should be used instead of creating a new {@link BasketItem}.
 * If a user adds an item that is already in the itemsToBuy list, that item should be moved to the top of the list.
 * If the item is in the itemsBought list, it should be removed from itemsBought and added to the top of the itemsToBuy list.
 *
 */

document.addEventListener('DOMContentLoaded', function () {
    const itemInput = document.querySelector<HTMLInputElement>('#item-input')!;
    const addBtn = document.querySelector<HTMLButtonElement>('#add-btn')!;
    const list = document.querySelector<HTMLUListElement>('#list')!;
    const completedList = document.querySelector<HTMLUListElement>('#completed-list')!;
    const suggestionsContainer = document.querySelector<HTMLDivElement>('#suggestions')!;

    const basketItemElementsByText: Map<string, HTMLLIElement> = new Map();

    const basketItemManager = new BasketItemManager(createBasketItemElement,
        (basketItem) => {
            const basketItemElement = basketItemElementFromBasketItem(basketItem);
            list.prepend(basketItemElement);

            basketItemElement.classList.remove('bought');
            checkboxForBasketItemElement(basketItemElement).checked = false;
        },
        (basketItem) => {
            const basketItemElement = basketItemElementFromBasketItem(basketItem);
            completedList.prepend(basketItemElement);

            basketItemElement.classList.add('bought');
            checkboxForBasketItemElement(basketItemElement).checked = true;
        }
    );

    const syncService = new SyncService(basketItemManager).start();

    const addItemFromUI = () => {
        const itemText = itemInput.value.trim();
        if (itemText === '') {
            return;
        }
        basketItemManager.addItemToBuyBasket(itemText);
        const basketItem = basketItemManager.getBasketItemByItemName(itemText)!;
        syncService.syncItemUpdate(basketItem);

        itemInput.value = '';
        suggestionsContainer.innerHTML = '';
        suggestionsContainer.style.display = 'none';
    };

    addBtn.onclick = addItemFromUI;
    itemInput.addEventListener('keypress', function (e) {
        if (e.key === 'Enter') {
            addItemFromUI();
        }
    });

    /**
     * Create new HTML element for the basket item. Doesn't add it to the document.
     */
    function createBasketItemElement(basketItem: BasketItem): void {
        const li = document.createElement('li');

        const checkbox = document.createElement('input');
        checkbox.type = 'checkbox';
        checkbox.addEventListener('change', checkBoxClickedListener(li));

        /*checkbox.checked = !basketItem.toBuy;
        if (!basketItem.toBuy) {
            li.classList.add('bought');
        }*/

        const textSpan = document.createElement('span');
        textSpan.className = 'item-text';
        textSpan.textContent = basketItem.name;

        li.appendChild(checkbox);
        li.appendChild(textSpan);

        basketItemElementsByText.set(basketItem.name, li);
    }

    function itemNameFromBasketItemElement(htmlElement: HTMLLIElement): string {
        return htmlElement.getElementsByTagName('span')[0].textContent;
    }

    function basketItemElementFromBasketItem(basketItem: BasketItem): HTMLLIElement {
        return basketItemElementsByText.get(basketItem.name)!;
    }

    function checkboxForBasketItemElement(htmlElement: HTMLLIElement): HTMLInputElement {
        return htmlElement.getElementsByTagName('input')[0]
    }

    function checkBoxClickedListener(li: HTMLLIElement) {
        return function (this: HTMLInputElement) {
            let itemName = itemNameFromBasketItemElement(li);
            let basketItem = basketItemManager.getBasketItemByItemName(itemName)!;
            if (this.checked) {
                basketItemManager.addToBoughtBasket(basketItem);
                syncService.syncItemUpdate(basketItem);
            } else {
                basketItemManager.addToBuyBasket(basketItem);
                syncService.syncItemUpdate(basketItem);
            }
        };
    }
});



