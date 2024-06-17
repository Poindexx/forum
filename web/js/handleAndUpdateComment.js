function handleComment(postId, username) {
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
            updateComment(postId, user.id, username);
        })
        .catch(error => {
            console.error('Ошибка:', error);
            alert('Произошла ошибка при обработке вашего запроса');
        });
    }

    function updateComment(postId, authorId, username) {
    const comment = document.getElementById('comment').value;
    console.log(comment)
    if (comment === "") {
        return
    }
    const url = `/updateComment`;
    const data = {
        post_id: postId,
        comment: comment,
        author_id: authorId,
        author: username
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
            return response.text().then(errorText => {
                throw new Error(errorText || 'Ошибка при создании комментария');
            });
        }
        return response.json();
    })
    .then(data => {
        document.getElementById('comment').value = ''; // Очистить текстовое поле
        fetchComments(); // Обновить список комментариев
    })
    .catch((error) => {
        console.error('Ошибка:', error);
        alert('Произошла ошибка при добавлении комментариев');
    });
}