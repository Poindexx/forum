
document.addEventListener('DOMContentLoaded', function() {
    const loginForm = document.getElementById('login-form');
    const tabPanes = document.querySelectorAll('.tab-pane');

    function updateFormAction(activeTabId) {
        if (activeTabId === 'pills-home') {
            loginForm.setAttribute('action', '/login-process');
        } else if (activeTabId === 'pills-profile') {
            loginForm.setAttribute('action', '/register-process');
        }
    }

    function resetFormFields() {
        loginForm.reset();
    }
    tabPanes.forEach(tabPane => {
        tabPane.addEventListener('click', function() {
            const activeTabId = this.getAttribute('id');
            updateFormAction(activeTabId);
        });

        // Добавляем обработчик нажатия клавиши Enter к каждому полю ввода формы
        loginForm.querySelectorAll('input').forEach(input => {
            input.addEventListener('keydown', function (e) {
                var key = e.which || e.keyCode;
                if (key === 13) {
                    const activeTabId = document.querySelector('.tab-pane.active').getAttribute('id');
                    if (activeTabId === 'pills-home') {
                        login();
                    } else if (activeTabId === 'pills-profile') {
                        register();
                    }
                }
            });
        });
    });

});

function login() {
    // Получение данных формы
    var username = document.getElementsByName("login_name")[0].value;
    var password = document.getElementsByName("login_password")[0].value;

    // Создание объекта XMLHttpRequest
    var xhr = new XMLHttpRequest();

    // Настройка запроса
    xhr.open("POST", "/login-process", true);
    xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");

    // Обработка ответа от сервера
    xhr.onload = function() {
        var response = JSON.parse(xhr.responseText);
        if (xhr.status == 200) {
            // Если ответ успешный, перенаправляем на главную страницу или выполняем другие действия
            window.location.href = "/";
        } else {
            // Если есть ошибка, отображаем ее на странице
            document.getElementById("error-message").innerText = response.error;
        }
    };

    // Отправка данных формы на сервер
    xhr.send("login_name=" + encodeURIComponent(username) + "&login_password=" + encodeURIComponent(password));
}

function register() {
    // Получение данных формы
    var username = document.getElementsByName("reg_name")[0].value;
    var password = document.getElementsByName("reg_password")[0].value;
    var phone = document.getElementsByName("reg_phone")[0].value;
    var email = document.getElementsByName("reg_email")[0].value;

    // Создание объекта XMLHttpRequest
    var xhr = new XMLHttpRequest();

    // Настройка запроса
    xhr.open("POST", "/register-process", true);
    xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");

    // Обработка ответа от сервера
    xhr.onload = function() {
        var response = JSON.parse(xhr.responseText);
        if (xhr.status == 200) {
            // Если ответ успешный, перенаправляем на главную страницу или выполняем другие действия
            window.location.href = "/";
        } else {
            // Если есть ошибка, отображаем ее на странице
            document.getElementById("error-message_reg").innerText = response.error;
        }
    };

    // Отправка данных формы на сервер
    xhr.send("reg_name=" + encodeURIComponent(username) + "&reg_password=" + encodeURIComponent(password) + "&reg_phone=" + encodeURIComponent(phone) + "&reg_email=" + encodeURIComponent(email) );
}