function categriadd() {
    const categoryName = document.getElementById('categoryName').value;

    fetch('/create-category', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            name: categoryName
        })
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Такое имя уже существует');
        }
        const modal = document.getElementById('createCategoryModal');
        modal.classList.remove('show');
        modal.style.display = 'none';
        document.querySelector('.modal-backdrop').remove();
        window.location.href = '/';
    })
    .catch(error => {
        console.error('Error:', error.message);
        document.getElementById('categoryError').innerText = 'Ошибка при создании категории: ' + error.message;
    });
}

function createPost() {
    event.preventDefault();
    const title = document.getElementById('postTitle').value;
    const anons = document.getElementById('postAnons').value;
    const description = document.getElementById('postDescription').value;
    // const imageURL = document.getElementById('postImageURL').value;
    const image = document.getElementById('inputGroupFile02').files[0];
    const categoryIDs = Array.from(document.getElementById('postCategoryID').selectedOptions).map(option => option.value);
    const categorys = Array.from(document.getElementById('postCategoryID').selectedOptions).map(option => option.textContent);
    const usernameElement = document.querySelector('.navbar-user[href="/"]');
    const username = usernameElement.textContent.trim();

    if (!image) {
        document.getElementById('postError').innerText = 'Выберите одну фотографию.';
        return;
    }
    var allFormat = ['image/jpeg', 'image/svg', 'image/png', 'image/gif', 'image/web', 'image/webp']
    if (!allFormat.includes(image.type)) {
        document.getElementById('postError').innerText = 'Выберите фомат jpeg, svg, png, gif, web, webp.';
        return;
    }
    if (image.size > 24 * 1024 * 1024) {
        document.getElementById('postError').innerText = 'Фото привышает 20гб.';
        return;
    }
    console.log(image)

    if (categorys.length === 0) {
        document.getElementById('postError').innerText = 'Выберите хотя бы одну категорию.';
        return;
    } else if (categorys.length > 3) {
        document.getElementById('postError').innerText = 'Выберите не более трех категорий.';
        return;
    } else {
        document.getElementById('postError').innerText = '';
    }

    const reader = new FileReader();
    reader.onload = function(e) {
        const base64Image = e.target.result.split(',')[1]; // Получение Base64 части

        fetch('/get-user-data')
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(users => {
                const user = users.find(user => user.username === username);
                if (!user) {
                    throw new Error('User not found');
                }
                return user.id;
            })
            .then(id => {
                const postData = {
                    title: title,
                    description: description,
                    anons: anons,
                    imageBase64: base64Image,
                    imageName: image.name,
                    categoryIDs: categoryIDs,
                    categorys: categorys,
                    author_id: id,
                    author: username
                };

                return fetch('/create-post', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(postData)
                });
            })
            .then(response => {
                if (!response.ok) {
                    return response.text().then(errorText => {
                        throw new Error(errorText || 'Ошибка при создании поста');
                    });
                }
                const modal = document.getElementById('createPostModal');
                modal.classList.remove('show');
                modal.style.display = 'none';
                document.querySelector('.modal-backdrop').remove();

                window.location.href = '/';
            })
            .catch(error => {
                console.error('Error:', error.message);
                document.getElementById('postError').innerText = 'Ошибка при создании поста: ' + error.message;
            });
    };
    reader.readAsDataURL(image); // Чтение файла как Data URL (Base64)
}

function fetchCategories() {
    fetch('/get-categories')
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok');
            }
            return response.json();
        })
        .then(categories => {
            const selectElement = document.getElementById('postCategoryID');
            selectElement.innerHTML = '';
            categories.forEach(category => {
                const option = document.createElement('option');
                option.value = category.id;
                option.textContent = category.name;
                selectElement.appendChild(option);
            });
        })
        .catch(error => console.error('Error:', error.message));
}

document.addEventListener('DOMContentLoaded', fetchCategories);