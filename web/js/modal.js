document.addEventListener('DOMContentLoaded', function() {
	const createCategoryLink = document.querySelector('.nav-link.create-category');

	// Обработчик нажатия на ссылку "Создать категорию"
	createCategoryLink.addEventListener('click', function(event) {
		event.preventDefault(); // Отмена действия по умолчанию (переход по ссылке)
		
		// Открываем модальное окно
		const modal = new bootstrap.Modal(document.getElementById('createCategoryModal'));
		modal.show();
	});
});

document.addEventListener('DOMContentLoaded', function() {
	const createPostLink = document.querySelector('.nav-link.create-post');

	// Обработчик нажатия на ссылку "Создать пост"
	createPostLink.addEventListener('click', function(event) {
		event.preventDefault(); // Отмена действия по умолчанию (переход по ссылке)
		
		// Открываем модальное окно
		const modal = new bootstrap.Modal(document.getElementById('createPostModal'));
		modal.show();
	});
});
