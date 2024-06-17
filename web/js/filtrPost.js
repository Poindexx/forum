$(document).ready(function() {
    $('#startDate').datepicker({
        format: 'dd.mm.yyyy',
        autoclose: true,
        todayHighlight: true
    });
    $('#endDate').datepicker({
        format: 'dd.mm.yyyy',
        autoclose: true,
        todayHighlight: true
    });

    $('#dateFilterForm').on('submit', function(e) {
        e.preventDefault();

        const categoryIDs = Array.from(document.getElementById('postCategoryID2').selectedOptions).map(option => option.value);
        const categoryNames = Array.from(document.getElementById('postCategoryID2').selectedOptions).map(option => option.textContent);
        const sortPost = document.getElementById('sort_post').value;
        const textDis = document.getElementById('text_dis').value;
        const startDate = $('#startDate').datepicker('getFormattedDate', 'dd.mm.yyyy');
        const endDate = $('#endDate').datepicker('getFormattedDate', 'dd.mm.yyyy');

        const data = {
            category_ids: categoryIDs,
            category_names: categoryNames,
            sort_post: sortPost,
            text_dis: textDis,
            start_date: startDate,
            end_date: endDate
        };
        console.log(data);

        fetch('/getSortedPost', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Проблема с сетью');
            }
            return response.json();
        })
        .then(data => {
            // Обработайте полученные данные здесь
            console.log('Успех:', data);
            // Обновите интерфейс с полученными данными
            updateUIWithPosts(data);
        })
        .catch(error => {
            const postsContainer = document.querySelector('.row.row-cols-md-1');
            postsContainer.innerHTML = '';
            postsContainer.innerHTML = '<h3 class="mt-5">Нету новостей</h3>';
            console.error('Ошибка:', error.message);
        });
    });
});

function updateUIWithPosts(posts) {
const postsContainer = document.querySelector('.row.row-cols-md-1');
postsContainer.innerHTML = '';

if (posts.length === 0) {
    postsContainer.innerHTML = '<h3 class="mt-5">Нету новостей</h3>';
    return;
}

posts.forEach(post => {
    const usernameElement = document.querySelector('.navbar-user[href="/"]');
    const username = usernameElement.textContent.trim();
    const postElement = document.createElement('div');
    postElement.classList.add('col-auto');
    postElement.innerHTML = `
        <div class="card mb-5">
            <div class="col-md-12">
                <div class="card-header">
                    <a href="/Id/${post.id}" class="card-title fs-2 fw-bold link-primary link-offset-2 link-underline-opacity-25 link-underline-opacity-100-hover">${post.title}</a>
                </div>
            </div>
            <div class="col-md-12">
                <div class="card-body">
                    <img src="${post.imageurl}" class="card-img-top img-fluid rounded" alt="Image">
                    <p class="card-title formatted-text mt-3">${post.anons}</p>
                    <p class="card-text">Категория:
                        ${(post.category && post.category_id) ? post.category.map((category, index) => `
                            <a href="/Categorys/${post.category_id[index]}">
                                <span class="badge rounded-pill text-bg-primary">${category}</span>
                            </a>
                        `).join(' ') : 'Без категории'}
                    </p>
                    <div class="card-text">
                        <img src="https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQNL_ZnOTpXSvhf1UaK7beHey2BX42U6solRA&usqp=CAU" alt="Автор" class="rounded-circle" style="width: 40px; height: 40px;">
                        <a href="/Author/${post.author_id}"><span class="ml-2">${post.author}</span></a>
                        <span class="ml-2"> | </span>
                        <span class="ml-auto">${post.date}</span>
                    </div>
                    <div class="card-text mt-3">
                        <button type="button" id="like-${post.id}" class="btn LikePost btn-primary mr-2" onclick="handleLikeDislike('${post.id}', 'like', '${username}')">👍 ${post.like}</button>
                        <button type="button" id="dislike-${post.id}" class="btn DislikePost btn-danger mr-2" onclick="handleLikeDislike('${post.id}', 'dislike', '${username}')">👎 ${post.doslike}</button>
                        <button type="button" onclick="window.location.href = '/Id/${post.id}';" class="btn btn-outline-secondary">${post.comment_len} Комментарии</button>
                    </div>
                </div>
            </div>
        </div>
    `;
    postsContainer.appendChild(postElement);
});
}