//
import {BasketItem} from "./suggestions";

const itemsToBuy: BasketItem[] = [];
const itemsBought: BasketItem[] = [];

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
 * TODO: Decide on a rendering strategy for the lists. We could use an observable collection to separate state from the UI, or have each item keep a direct reference to its corresponding HTML element.
 * TODO: Implement ^, the rendering should be local, not distributed across multiple methods.
 */

document.addEventListener('DOMContentLoaded', function () {
    const itemInput = document.querySelector<HTMLInputElement>('#item-input')!;
    const addBtn = document.querySelector<HTMLButtonElement>('#add-btn')!;
    const list = document.querySelector<HTMLUListElement>('#list')!;
    const completedList = document.querySelector<HTMLUListElement>('#completed-list')!;
    const suggestionsContainer = document.querySelector<HTMLDivElement>('#suggestions')!;

    itemInput.addEventListener('keypress', function (e) {
        if (e.key === 'Enter') {
            addItem();
        }
    });

    const addItem = () => {
        const itemText = itemInput.value.trim();
        if (itemText === '') {
            return;
        }
        const li = document.createElement('li');

        const checkbox = document.createElement('input');
        checkbox.type = 'checkbox';
        checkbox.addEventListener('change', checkBoxClickedListener(li, completedList, list));

        const textSpan = document.createElement('span');
        textSpan.className = 'item-text';
        textSpan.textContent = itemText;

        li.appendChild(checkbox);
        li.appendChild(textSpan);
        list.appendChild(li);
        itemsToBuy.push(new BasketItem(itemText, new Date()))

        itemInput.value = '';
        suggestionsContainer.innerHTML = '';
        suggestionsContainer.style.display = 'none';
    };

    addBtn.onclick = addItem;
});

function checkBoxClickedListener(li: HTMLLIElement, completedList: HTMLUListElement, list: HTMLUListElement) {
    return function (this: HTMLInputElement) {
        li.classList.toggle('bought');
        if (this.checked) {
            completedList.appendChild(li);
        } else {
            list.appendChild(li);
        }
    };
}