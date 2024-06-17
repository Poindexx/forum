document.addEventListener('DOMContentLoaded', function() {
    const usernameElement = document.querySelector('.navbar-user[href="/"]');
    const username = usernameElement.textContent.trim();
    if (username.trim() !== "") {
        const loginButton = document.querySelector('.navbar .btn-primary');
        loginButton.setAttribute('href', '/exit');
        loginButton.classList.remove('btn-primary');
        loginButton.classList.add('btn');
        loginButton.innerText = '✖️';

        const navLinks = document.querySelectorAll('.navbar-nav .nav-link.disabled');
        navLinks.forEach(link => {
            link.classList.remove('disabled');
        });
    }
});