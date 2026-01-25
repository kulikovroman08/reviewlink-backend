// common.js - общие утилиты и константы
const API_BASE = window.location.origin;

// ===== Alerts (lightweight, no bootstrap) =====
function _resolveContainer(container) {
    if (container && container !== document.body) return container;

    return (
        document.querySelector(".container.mt-4") ||
        document.querySelector("main .container") ||
        document.querySelector(".container") ||
        document.body
    );
}

function _renderAlert(type, message, container, opts = {}) {
    const {
        autoHideMs = type === "success" ? 1600 : 5000,
        replace = true,
        closable = type !== "success", // ✅ успех без крестика
    } = opts;

    const host = _resolveContainer(container);

    if (replace && host !== document.body && host.innerHTML != null) {
        host.innerHTML = "";
    }

    const el = document.createElement("div");
    el.className = `rl-alert rl-alert--${type}`;
    el.innerHTML = `
    <div class="rl-alert__icon">${type === "success" ? "✅" : "⚠️"}</div>
    <div class="rl-alert__text"></div>
    ${closable ? `<button class="rl-alert__close" type="button" aria-label="Закрыть">×</button>` : ""}
  `;

    const text = el.querySelector(".rl-alert__text");
    if (text) text.textContent = String(message ?? "");

    host.prepend(el);

    if (closable) {
        const closeBtn = el.querySelector(".rl-alert__close");
        if (closeBtn) closeBtn.addEventListener("click", () => el.remove());
    }

    if (autoHideMs && autoHideMs > 0) {
        window.setTimeout(() => {
            if (el && el.parentNode) el.remove();
        }, autoHideMs);
    }

    return el;
}

function showError(message, container = document.body) {
    return _renderAlert("error", message, container, { autoHideMs: 6000, replace: true });
}

function showSuccess(message, container = document.body) {
    return _renderAlert("success", message, container, { autoHideMs: 1600, replace: true, closable: false });
}


// Выход из системы
function logout() {
    localStorage.removeItem("userToken");
    localStorage.removeItem("adminToken");
    localStorage.removeItem("userEmail");
    localStorage.removeItem("adminEmail");
    window.location.href = "login.html";
}

function _getAnyToken() {
    // поддержка старого adminToken + нового userToken
    let t = localStorage.getItem("userToken") || localStorage.getItem("adminToken") || "";
    // если вдруг сохранили "Bearer xxx"
    if (t.startsWith("Bearer ")) t = t.slice(7).trim();
    return t;
}

// Генерация QR кода (картинка, в QR коде лежит qr_token)
function generateQRCode(data, size = 200) {
    return `https://api.qrserver.com/v1/create-qr-code/?size=${size}x${size}&data=${encodeURIComponent(data)}`;
}

async function apiFetch(path, options = {}) {
    const token = _getAnyToken();

    const headers = {
        ...(options.headers || {}),
        Accept: "application/json",
    };

    let body = options.body;
    const isFormData = body instanceof FormData;

    if (body && typeof body === "object" && !isFormData) {
        headers["Content-Type"] = "application/json";
        body = JSON.stringify(body);
    }

    if (token) headers.Authorization = token;

    const res = await fetch(`${API_BASE}${path}`, {
        ...options,
        headers,
        body,
    });

    if (res.status === 401) {
        logout();
        throw new Error("Unauthorized");
    }

    const contentType = res.headers.get("content-type") || "";
    const payload = contentType.includes("application/json")
        ? await res.json()
        : await res.text();

    if (!res.ok) {
        const msg =
            (payload && payload.error) ||
            (payload && payload.message) ||
            (typeof payload === "string" && payload) ||
            `Request failed: ${res.status}`;
        throw new Error(msg);
    }

    return payload;
}

function parseJwt(token) {
    try {
        const base64Url = token.split(".")[1];
        const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
        const jsonPayload = decodeURIComponent(
            atob(base64)
                .split("")
                .map(c => "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2))
                .join("")
        );
        return JSON.parse(jsonPayload);
    } catch {
        return null;
    }
}


// ===== AUTH HELPERS (exported) =====
function requireUser() {
    const token = _getAnyToken();
    if (!token) {
        window.location.href = "login.html";
        return null;
    }
    return token;
}

function requireAdmin() {
    const token = _getAnyToken();
    if (!token) {
        window.location.href = "login.html";
        return null;
    }

    const payload = parseJwt(token);
    if (!payload || payload.role !== "admin") {
        window.location.href = "dashboard.html";
        return null;
    }

    return token;
}

// theme init + toggle
(function () {
    const root = document.documentElement;
    const saved = localStorage.getItem("rl_theme");

    root.dataset.theme = saved === "pastel" ? "pastel" : "dark";
    localStorage.setItem("rl_theme", root.dataset.theme);

    window.RL = window.RL || {};
    window.RL.toggleTheme = function () {
        const next = root.dataset.theme === "pastel" ? "dark" : "pastel";
        root.dataset.theme = next;
        localStorage.setItem("rl_theme", next);
    };
})();

// Экспортируем для использования в других файлах
window.AppCommon = {
    API_BASE,
    showError,
    showSuccess,
    logout,
    generateQRCode,
    apiFetch,
    requireUser,
    requireAdmin,
    parseJwt,
};
