function fetchCategories1() {
    fetch('/get-categories')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(categories => {
            const selectElement = document.getElementById('rendCategory');
            selectElement.innerHTML = '';
            categories.forEach(category => {
                const button = document.createElement('button');
                button.className = "btn btn-primary";
                button.setAttribute('onclick', `document.location='/Categorys/${category.id}'`);
                button.innerText = category.name;
                selectElement.appendChild(button);
            });
        })
        .catch(error => console.error('Error:', error.message));
    }

    document.addEventListener('DOMContentLoaded', fetchCategories1);

document.addEventListener('DOMContentLoaded', function() {
    const selectElement = document.getElementById('postCategoryID2');

// Проверка, если элемент уже инициализирован
if (selectElement && !selectElement.classList.contains('choices')) {
    const choices = new Choices(selectElement, {
        allowHTML: true,
        removeItemButton: true,
        searchResultLimit: 3,
        renderChoiceLimit: 10
    });

    function fetchCategories2() {
        fetch('/get-categories')
            .then(response => {
                if (!response.ok) {
                    throw new Error('Нету категории');
                }
                return response.json();
            })
            .then(categories => {
                choices.clearStore();
                categories.forEach(category => {
                    choices.setChoices([
                        { value: category.id, label: category.name, selected: false }
                    ], 'value', 'label', false);
                });
            })
            .catch(error => console.error('Error:', error.message));
    }
    fetchCategories2();
}
});