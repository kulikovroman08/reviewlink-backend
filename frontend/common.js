// common.js - общие утилиты и константы
const API_BASE = window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1'
    ? 'http://172.27.78.199:8080'
    : 'http://localhost:8080';

// Проверка авторизации
function checkAuth() {
    const token = localStorage.getItem('userToken');
    console.log('Токен из localStorage:', token);
    if (!token) {
        window.location.href = 'login.html';
        return null;
    }
    return token;
}

// Показ ошибок
function showError(message, container = document.body) {
    const alert = document.createElement('div');
    alert.className = 'alert alert-danger alert-dismissible fade show';
    alert.innerHTML = `
        <i class="bi bi-exclamation-triangle me-2"></i>${message}
        <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
    `;

    if (container === document.body) {
        // Вставляем в начало контейнера
        const mainContainer = document.querySelector('.container');
        if (mainContainer) {
            mainContainer.insertBefore(alert, mainContainer.firstChild);
        } else {
            container.prepend(alert);
        }
    } else {
        container.innerHTML = alert.outerHTML;
    }

    return alert;
}

// Показ успеха
function showSuccess(message, container = document.body) {
    const alert = document.createElement('div');
    alert.className = 'alert alert-success alert-dismissible fade show';
    alert.innerHTML = `
        <i class="bi bi-check-circle me-2"></i>${message}
        <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
    `;

    if (container === document.body) {
        const mainContainer = document.querySelector('.container');
        if (mainContainer) {
            mainContainer.insertBefore(alert, mainContainer.firstChild);
        } else {
            container.prepend(alert);
        }
    } else {
        container.innerHTML = alert.outerHTML;
    }

    return alert;
}

// Выход из системы
function logout() {
    localStorage.removeItem('userToken');
    localStorage.removeItem('userEmail');
    const currentUrl = window.location.href;
    const newUrl = currentUrl.replace('/frontend/dashboard.html', '/frontend/login.html');
    window.location.href = newUrl;
}

// Генерация QR кода
function generateQRCode(data, size = 200) {
    return `https://api.qrserver.com/v1/create-qr-code/?size=${size}x${size}&data=${encodeURIComponent(data)}`;
}

// Экспортируем для использования в других файлах
window.AppCommon = {
    API_BASE,
    checkAuth,
    showError,
    showSuccess,
    logout,
    generateQRCode
};