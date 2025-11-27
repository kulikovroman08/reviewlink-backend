// login.js
document.addEventListener('DOMContentLoaded', function () {
    const { API_BASE, showError, showSuccess } = window.AppCommon;

    // Проверяем, если уже авторизован - перенаправляем
    const token = localStorage.getItem('userToken');
    if (token) {
        const currentUrl = window.location.href;
        const newUrl = currentUrl.replace('/frontend/login.html', '/frontend/dashboard.html');
        window.location.href = newUrl;
    }

    document.getElementById('loginForm').addEventListener('submit', async function (e) {
        e.preventDefault();

        const email = document.getElementById('emailInput').value.trim();
        const password = document.getElementById('passwordInput').value;
        const messageDiv = document.getElementById('loginMessage');
        const button = e.target.querySelector('button[type="submit"]');

        // Валидация
        if (!email || !password) {
            showError('Заполните все поля', messageDiv);
            return;
        }

        // Показываем загрузку
        button.disabled = true;
        button.innerHTML = '<div class="spinner-border spinner-border-sm me-2"></div>Вход...';
        messageDiv.innerHTML = '';

        try {
            const response = await fetch(`${API_BASE}/login`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    email: email,
                    password: password
                })
            });

            if (response.ok) {
                const data = await response.json();

                localStorage.setItem('userToken', data.token);
                localStorage.setItem('userEmail', email);

                // Показываем сообщение
                showSuccess('Успешный вход! Перенаправление...', messageDiv);

                setTimeout(() => {
                    const params = new URLSearchParams(window.location.search);
                    const redirectUrl = params.get("redirect");

                    if (redirectUrl) {
                        window.location.href = redirectUrl;
                    } else {
                        window.location.href = "dashboard.html";
                    }
                }, 500);

            } else {
                const errorData = await response.json();
                let errorMessage = 'Ошибка входа';

                if (errorData.error) {
                    errorMessage = errorData.error;
                } else if (response.status === 401) {
                    errorMessage = 'Неверный email или пароль';
                }

                showError(errorMessage, messageDiv);
            }
        } catch (error) {
            console.error('Ошибка:', error);
            showError('Ошибка соединения с сервером', messageDiv);
        } finally {
            // Восстанавливаем кнопку
            button.disabled = false;
            button.innerHTML = '<i class="bi bi-box-arrow-in-right me-2"></i>Войти';
        }
    });
});