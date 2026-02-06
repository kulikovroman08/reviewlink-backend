// login.js
document.addEventListener('DOMContentLoaded', function () {
    const { API_BASE, showError, showSuccess } = window.AppCommon;

    let mode = "login"; // "login" | "signup"

    const modeLoginBtn = document.getElementById("modeLoginBtn");
    const modeSignupBtn = document.getElementById("modeSignupBtn");
    const nameWrap = document.getElementById("nameWrap");
    const nameInput = document.getElementById("nameInput");

    const ownerWrap = document.getElementById("ownerWrap");
    const ownerCheckbox = document.getElementById("ownerCheckbox");
    const submitBtn = document.getElementById("submitBtn");

    function setMode(next) {
        mode = next;
        const isSignup = mode === "signup";

        if (nameWrap) nameWrap.style.display = isSignup ? "" : "none";
        if (ownerWrap) ownerWrap.style.display = isSignup ? "" : "none";
        if (!isSignup && ownerCheckbox) ownerCheckbox.checked = false;

        if (modeLoginBtn) modeLoginBtn.classList.toggle("btn-primary", !isSignup);
        if (modeSignupBtn) modeSignupBtn.classList.toggle("btn-primary", isSignup);

        if (modeLoginBtn) modeLoginBtn.setAttribute("aria-pressed", String(!isSignup));
        if (modeSignupBtn) modeSignupBtn.setAttribute("aria-pressed", String(isSignup));

        if (submitBtn) {
            submitBtn.innerHTML = isSignup
                ? '<i class="bi bi-person-plus"></i>Зарегистрироваться'
                : '<i class="bi bi-box-arrow-in-right"></i>Войти';
        }
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

    async function redirectAfterAdminLogin(token) {
        const base = API_BASE || window.location.origin;

        try {
            const res = await fetch(`${base}/admin/owner_requests/pending`, {
                method: "GET",
                headers: { Authorization: token },
            });

            // 200 => platform admin
            if (res.ok) {
                window.location.href = "admin-owner-requests.html";
                return;
            }

            // обычный админ
            window.location.href = "admin.html";
        } catch {
            window.location.href = "admin.html";
        }
    }

    document.getElementById('loginForm').addEventListener('submit', async function (e) {
        e.preventDefault();

        const email = document.getElementById('emailInput').value.trim();
        const password = document.getElementById('passwordInput').value;
        const isSignup = mode === "signup";
        const name = (nameInput ? nameInput.value.trim() : "");
        const wantsOwnerFlow = isSignup && !!(ownerCheckbox && ownerCheckbox.checked);

        const messageDiv = document.getElementById('loginMessage');
        const button = e.target.querySelector('button[type="submit"]');

        if (!email || !password || (isSignup && !name)) {
            showError(isSignup ? 'Заполните имя, email и пароль' : 'Заполните все поля', messageDiv);
            return;
        }

        button.disabled = true;
        button.innerHTML = isSignup ? 'Регистрация...' : 'Вход...';
        messageDiv.innerHTML = '';

        try {
            const endpoint = isSignup ? "signup" : "login";

            const response = await fetch(`${API_BASE}/${endpoint}`, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
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

                showSuccess(isSignup ? "Аккаунт создан!" : "Успешный вход!", messageDiv);

                setTimeout(async () => {
                    const p = new URLSearchParams(window.location.search);
                    const redirect = p.get("redirect");
                    if (redirect) {
                        window.location.href = redirect;
                        return;
                    }

                    // Админ: всегда в админку (platform admin может улететь в owner-requests)
                    if (isAdmin) {
                        await redirectAfterAdminLogin(data.token);
                        return;
                    }

                    // Юзер: если при регистрации поставил "я владелец" — на форму заявки
                    if (wantsOwnerFlow) {
                        window.location.href = "owner-request.html";
                        return;
                    }

                    window.location.href = "dashboard.html";
                }, 400);

            } else {
                let errorData = {};
                try { errorData = await response.json(); } catch { errorData = {}; }

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
            button.disabled = false;
            button.innerHTML = isSignup
                ? '<i class="bi bi-person-plus"></i>Зарегистрироваться'
                : '<i class="bi bi-box-arrow-in-right"></i>Войти';
        }
    });
});
