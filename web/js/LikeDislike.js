function handleLikeDislike(postId, type, username) {
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
            updateLikeDislike(postId, type, user.id);
        })
        .catch(error => {
            console.error('Ошибка:', error);
            alert('Произошла ошибка при обработке вашего запроса');
        });
}

function updateLikeDislike(postId, type, authorId) {
    const url = `/updateLikes`;
    const data = {
        postId: postId,
        type: type,
        authorId: authorId
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
        if (data.newPostid === postId) {
            const likeButton = document.querySelector(`#like-${postId}`);
            const dislikeButton = document.querySelector(`#dislike-${postId}`);
            likeButton.innerHTML = `👍 ${data.newLikeCount}`;
            dislikeButton.innerHTML = `👎 ${data.newDislikeCount}`;
        }
    })
    .catch((error) => {
        console.error('Ошибка:', error);
        alert('Произошла ошибка при обновлении лайков/дизлайков');
    });
}