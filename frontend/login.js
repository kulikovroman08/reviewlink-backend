// login.js
document.addEventListener('DOMContentLoaded', function () {
    const { API_BASE, showError, showSuccess } = window.AppCommon;

    function getRoleFromJwt(token) {
        try {
            const parts = String(token || "").split(".");
            if (parts.length !== 3) return "";
            const payloadB64 = parts[1].replace(/-/g, "+").replace(/_/g, "/");
            const payloadJson = decodeURIComponent(
                atob(payloadB64)
                    .split("")
                    .map(c => "%" + c.charCodeAt(0).toString(16).padStart(2, "0"))
                    .join("")
            );
            const payload = JSON.parse(payloadJson);
            return payload.role || "";
        } catch {
            return "";
        }
    }


    // Проверяем, если уже авторизован - перенаправляем
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
        button.innerHTML = 'Вход...';
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
                const role = getRoleFromJwt(data.token);
                const isAdmin = String(role).toLowerCase() === "admin";

                if (isAdmin) {
                    localStorage.setItem("adminToken", data.token);
                    localStorage.setItem("adminEmail", email);
                } else {
                    localStorage.setItem("userToken", data.token);
                    localStorage.setItem("userEmail", email);
                }

                showSuccess("Успешный вход! Перенаправление...", messageDiv);

                setTimeout(() => {
                    const p = new URLSearchParams(window.location.search);
                    const redirect = p.get("redirect");

                    if (redirect) {
                        window.location.href = redirect;
                        return;
                    }

                    window.location.href = isAdmin ? "admin.html" : "dashboard.html";
                }, 400);

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
            button.innerHTML = '<i class="bi bi-box-arrow-in-right"></i>Войти';
        }
    });
});