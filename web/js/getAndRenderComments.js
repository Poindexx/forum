document.addEventListener('DOMContentLoaded', () => {
    fetchComments();
});

function fetchComments() {
    const postIdElement = document.querySelector('h2');
    if (postIdElement=== null) {
        return
    }
    const postId = postIdElement.getAttribute('id');
    const url = '/get-comments'; // URL для отправки запроса
    const data = {
        post_id: postId,
    };
    fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok');
        }
        return response.json();
    })
    .then(data => {
        renderComments(data);
    })
    .catch((error) => {
        console.error('There was a problem with the fetch operation:', error);
    });
}

function renderComments(comments) {
    if (comments === null) {
        return
    }
    const postCommentsElement = document.getElementById('post-comments');
    postCommentsElement.innerHTML = ''; // Очистить существующие комментарии
    comments.forEach(comment => {
        const li = document.createElement('li');
        li.className = 'media';

        const mediaBodyDiv = document.createElement('div');
        mediaBodyDiv.className = 'media-body';

        const img = document.createElement('img');
        img.src = "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcQNL_ZnOTpXSvhf1UaK7beHey2BX42U6solRA&usqp=CAU";
        img.alt = 'Автор';
        img.className = 'rounded-circle';
        img.style.width = '40px';
        img.style.height = '40px';
        
        const a = document.createElement('a');
        a.href = "/Author/" + comment.author_id;
        a.className = 'text-success';
        
        const aspan = document.createElement('span');
        aspan.className = "ml-2";
        aspan.textContent = " " + comment.author;
        
        a.appendChild(aspan);
        
        const span = document.createElement('span');
        span.className = 'text-muted';
        const small = document.createElement('small');
        small.className = 'text-muted';
        small.textContent = "   " + comment.date;

        const p = document.createElement('p');
        p.className = 'mt-2';
        p.textContent = comment.comment;

        mediaBodyDiv.appendChild(img);
        mediaBodyDiv.appendChild(a);
        span.appendChild(small);
        mediaBodyDiv.appendChild(span);
        mediaBodyDiv.appendChild(p);

        const mediaRightDiv = document.createElement('div');
        mediaRightDiv.className = 'form-check-reverse';

        const likeButton = document.createElement('button');
        likeButton.type = 'button';
        likeButton.className = 'btn m-1 LikePost btn-primary';
        likeButton.setAttribute('onclick', `handleComLikeDislike('${comment.id}', 'like')`)
        likeButton.setAttribute('id', `clike-${comment.id}`)
        likeButton.textContent = `👍 ${comment.like}`;

        const dislikeButton = document.createElement('button');
        dislikeButton.type = 'button';
        dislikeButton.className = 'btn m-1 DislikePost btn-danger';
        dislikeButton.setAttribute('onclick', `handleComLikeDislike('${comment.id}', 'dislike')`)
        dislikeButton.setAttribute('id', `cdislike-${comment.id}`)
        dislikeButton.textContent = `👎 ${comment.dislike}`;

        mediaRightDiv.appendChild(likeButton);
        mediaRightDiv.appendChild(dislikeButton);

        li.appendChild(mediaBodyDiv);
        li.appendChild(mediaRightDiv);
        const hr = document.createElement('hr');

        postCommentsElement.appendChild(li);
        postCommentsElement.appendChild(hr);
    });
}

function handleComLikeDislike(comId, type) {
    const usernameElement = document.querySelector('.navbar-user[href="/"]');
    const username = usernameElement.textContent.trim();
    if (!username) {
        var toast = new bootstrap.Toast(document.getElementById('copyToast'));
        var texterr = document.getElementById("toastErrText")
        texterr.innerHTML = "Ошибка: пользователь не авторизован"
        toast.show();
        return;
    }

    fetch('/get-user-data')
        .then(response => {
            if (!response.ok) {
                throw new Error('Ошибка сети');
            }
            return response.json();
        })
        .then(users => {
            const user = users.find(user => user.username === username);
            if (!user) {
                throw new Error('Пользователь не найден');
            }
            updateComLikeDislike(comId, type, user.id);
        })
        .catch(error => {
            console.error('Ошибка:', error);
            alert('Произошла ошибка при обработке вашего запроса');
        });
}

function updateComLikeDislike(comId, type, username) {
    const url = `/updateLikesCom`;
    const data = {
        comId: comId,
        type: type,
        authorId: username
    };

    fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
    })
    .then(response => response.json())
    .then(data => {
        if (String(data.newComid) === String(comId)) {
            const likeButton = document.querySelector(`#clike-${comId}`);
            const dislikeButton = document.querySelector(`#cdislike-${comId}`);
            likeButton.innerHTML = `👍 ${data.newLikeCount}`;
            dislikeButton.innerHTML = `👎 ${data.newDislikeCount}`;
        } 
    })
    .catch((error) => {
        console.error('Ошибка:', error);
        alert('Произошла ошибка при обновлении лайков/дизлайков');
    });
}
