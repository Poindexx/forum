document.addEventListener('DOMContentLoaded', function() {
    const currentUrl = window.location.href;
    const linkAdr = document.getElementById("linkAdr");
    if (linkAdr === null) {
        return
    }
    linkAdr.innerHTML = currentUrl;

    document.getElementById('vk').href = 'https://vk.com/share.php?url=' + currentUrl;
    document.getElementById('instagram').href = 'https://www.instagram.com/?url=' + currentUrl;
    document.getElementById('twitter').href = 'https://twitter.com/intent/tweet?url=' + currentUrl + '&text=Текст для твита';
    document.getElementById('github').href = 'https://github.com';
    document.getElementById('facebook').href = 'https://www.facebook.com/sharer/sharer.php?u=' + currentUrl;
    var copyButton = document.getElementById('copy-url-button');
    copyButton.addEventListener('click', function() {
        navigator.clipboard.writeText(window.location.href).then(function() {
            var toast = new bootstrap.Toast(document.getElementById('copyToast'));
            var texterr = document.getElementById("toastErrText")
            texterr.innerHTML = "Ссылка скопирована в буфер обмена!"
            toast.show();
        }).catch(function(err) {
            var toast = new bootstrap.Toast(document.getElementById('copyToast'));
            var texterr = document.getElementById("toastErrText")
            texterr.innerHTML = "Ошибка при копировании: "
            toast.show();
            console.error('Ошибка при копировании: ', err);
        });
    });
});