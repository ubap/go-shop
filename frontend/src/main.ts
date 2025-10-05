//

const itemInput = document.querySelector<HTMLInputElement>('#item-input')!;
// const addBtn = document.getElementById('add-btn');
const list = document.getElementById('list')!;
const completedList = document.getElementById('completed-list')!;
const suggestionsContainer = document.getElementById('suggestions')!;


itemInput.addEventListener('keypress', function (e) {
    if (e.key === 'Enter') {
        addItem();
    }
});


const addItem = () => {
    const itemText = itemInput.value.trim();
    if (itemText !== '') {
        const li = document.createElement('li');

        const checkbox = document.createElement('input');
        checkbox.type = 'checkbox';
        checkbox.addEventListener('change', function() {
            li.classList.toggle('bought');
            if (this.checked) {
                completedList.appendChild(li);
            } else {
                list.appendChild(li);
            }
        });

        const textSpan = document.createElement('span');
        textSpan.className = 'item-text';
        textSpan.textContent = itemText;

        li.appendChild(checkbox);
        li.appendChild(textSpan);
        list.appendChild(li);

        itemInput.value = '';
        suggestionsContainer.innerHTML = '';
        suggestionsContainer.style.display = 'none';
    }
};
