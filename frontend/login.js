// login.js
document.addEventListener('DOMContentLoaded', function () {
    const { API_BASE, showError, showSuccess } = window.AppCommon;
    // ===== mode toggle (login/signup) =====
    let mode = "login"; // "login" | "signup"

    const modeLoginBtn = document.getElementById("modeLoginBtn");
    const modeSignupBtn = document.getElementById("modeSignupBtn");
    const nameWrap = document.getElementById("nameWrap");
    const nameInput = document.getElementById("nameInput");

    function setMode(next) {
        mode = next;
        const isSignup = mode === "signup";

        if (nameWrap) nameWrap.style.display = isSignup ? "" : "none";

        // подсветка активной кнопки (используем твои .btn/.btn-primary)
        if (modeLoginBtn) modeLoginBtn.classList.toggle("btn-primary", !isSignup);
        if (modeSignupBtn) modeSignupBtn.classList.toggle("btn-primary", isSignup);

        if (modeLoginBtn) modeLoginBtn.setAttribute("aria-pressed", String(!isSignup));
        if (modeSignupBtn) modeSignupBtn.setAttribute("aria-pressed", String(isSignup));
    }

    if (modeLoginBtn) modeLoginBtn.addEventListener("click", () => setMode("login"));
    if (modeSignupBtn) modeSignupBtn.addEventListener("click", () => setMode("signup"));
    setMode("login");


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


    // Проверяем, если уже авторизован
    document.getElementById('loginForm').addEventListener('submit', async function (e) {
        e.preventDefault();

        const email = document.getElementById('emailInput').value.trim();
        const password = document.getElementById('passwordInput').value;
        const isSignup = mode === "signup";
        const name = (nameInput ? nameInput.value.trim() : "");
        const messageDiv = document.getElementById('loginMessage');
        const button = e.target.querySelector('button[type="submit"]');

        // Валидация
        if (!email || !password || (isSignup && !name)) {
            showError(isSignup ? 'Заполните имя, email и пароль' : 'Заполните все поля', messageDiv);
            return;
        }

        // Показываем загрузку
        button.disabled = true;
        button.innerHTML = isSignup ? 'Регистрация...' : 'Вход...';
        messageDiv.innerHTML = '';

        try {
            const endpoint = isSignup ? "signup" : "login";

            const response = await fetch(`${API_BASE}/${endpoint}`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json"
                },
                body: JSON.stringify(
                    isSignup
                        ? { email: email, name: name, password: password }
                        : { email: email, password: password }
                )
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

                showSuccess(
                    isSignup
                        ? "Аккаунт создан!"
                        : "Успешный вход!",
                    messageDiv
                );

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
                let errorData = {};
                try {
                    errorData = await response.json();
                } catch {
                    errorData = {};
                }
                let errorMessage = 'Ошибка входа';

                if (errorData.error) {
                    errorMessage = errorData.error;
                } else if (!isSignup && response.status === 401) {
                    errorMessage = 'Неверный email или пароль';
                } else if (isSignup && response.status === 409) {
                    errorMessage = 'Email уже используется';
                }

                showError(errorMessage, messageDiv);
            }
        } catch (error) {
            console.error('Ошибка:', error);
            showError('Ошибка соединения с сервером', messageDiv);
        } finally {
            // Восстанавливаем кнопку
            button.disabled = false;
            button.innerHTML = isSignup
                ? '<i class="bi bi-person-plus"></i>Зарегистрироваться'
                : '<i class="bi bi-box-arrow-in-right"></i>Войти';
        }
    });
});