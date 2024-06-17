function createPostinS() {
    const usernameElement = document.querySelector('.navbar-user[href="/"]');
    const username = usernameElement.textContent.trim();
    if (!username) {
      var toast = new bootstrap.Toast(document.getElementById('copyToast'));
      var texterr = document.getElementById("toastErrText")
      texterr.innerHTML = "Ошибка: пользователь не авторизован"
      toast.show();
      return;
    }
    const modal = new bootstrap.Modal(document.getElementById('createPostModal'));
    modal.show();
    
  }