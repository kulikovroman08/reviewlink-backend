const API_BASE = (window.location.port === "5500" || window.location.port === "8000")
    ? "http://localhost:8080"
    : window.location.origin;

let BEARER_TOKEN = "";

// Генерация QR без place_id
function generateQRCode(token, placeId) {
    const reviewUrl = `${API_BASE}/frontend/review-form.html?token=${token}&place_id=${placeId}`;
    return `https://api.qrserver.com/v1/create-qr-code/?size=150x150&data=${encodeURIComponent(reviewUrl)}`;
}

// Авторизация
document.getElementById("loginBtn").onclick = async function () {
    const email = document.getElementById("emailInput").value.trim();
    const password = document.getElementById("passwordInput").value;
    const statusDiv = document.getElementById("loginStatus");

    if (!email || !password) {
        statusDiv.innerHTML = '<div class="text-danger">Заполните email и пароль</div>';
        return;
    }

    statusDiv.innerHTML = '<div class="text-info">Вход в систему...</div>';

    try {
        const response = await fetch(`${API_BASE}/login`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ email, password })
        });

        const data = await response.json();

        if (!response.ok) {
            throw new Error(data.error || "Ошибка входа");
        }

        BEARER_TOKEN = data.token;
        localStorage.setItem("adminToken", BEARER_TOKEN);
        localStorage.setItem("adminEmail", email);

        statusDiv.innerHTML = '<div class="text-success">Успешный вход!</div>';
        document.getElementById("generateBtn").disabled = false;

        await loadPlaces();

    } catch (error) {
        statusDiv.innerHTML = `<div class="text-danger">${error.message}</div>`;
    }
};

// Автологин
document.addEventListener("DOMContentLoaded", async function () {
    const savedToken = localStorage.getItem("adminToken");
    const savedEmail = localStorage.getItem("adminEmail");

    if (savedToken) {
        BEARER_TOKEN = savedToken;
        if (savedEmail) document.getElementById("emailInput").value = savedEmail;

        document.getElementById("loginStatus").innerHTML = '<div class="text-success">Авторизован</div>';
        document.getElementById("generateBtn").disabled = false;

        await loadPlaces();
    }
});

// Загрузка заведений
async function loadPlaces() {
    const select = document.getElementById("placeSelect");
    select.innerHTML = `<option>Загрузка...</option>`;

    try {
        const response = await fetch(`${API_BASE}/places`, {
            method: "GET",
            headers: {
                "Authorization": "Bearer " + BEARER_TOKEN
            }
        });

        const places = await response.json();

        select.innerHTML = "";
        places.forEach(place => {
            const option = document.createElement("option");
            option.value = place.id;
            option.textContent = `${place.name} (${place.address})`;
            select.appendChild(option);
        });

    } catch (error) {
        select.innerHTML = `<option value="">Ошибка загрузки</option>`;
    }
}

// Генерация токенов
document.getElementById("generateBtn").onclick = async function () {
    const placeId = document.getElementById("placeSelect").value;
    const count = document.getElementById("countInput").value;

    if (!placeId) {
        alert("Выберите заведение");
        return;
    }

    const container = document.getElementById("resultsContainer");
    container.innerHTML = `
        <div class="alert alert-info">Генерация токенов...</div>
    `;

    try {
        const response = await fetch(`${API_BASE}/admin/tokens`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": "Bearer " + BEARER_TOKEN
            },
            body: JSON.stringify({
                place_id: placeId,
                count: Number(count)
            })
        });

        const data = await safeParseJSON(response);

        if (!response.ok) {
            throw new Error(data.error || "Ошибка генерации");
        }

        showResults(data, container, placeId);

    } catch (error) {
        container.innerHTML = `
            <div class="alert alert-danger">${error.message}</div>
        `;
    }
};

// Безопасный JSON
async function safeParseJSON(response) {
    try {
        return await response.json();
    } catch {
        return {};
    }
}

// Показ результата
function showResults(data, container, placeId) {
    const tokens = data.tokens || data.Tokens || [];

    let html = `
        <div class="card">
            <div class="card-header bg-success text-white">
                Успешно сгенерировано ${tokens.length} токенов
            </div>
            <div class="card-body">
    `;

    tokens.forEach((token, i) => {
        const url = `${API_BASE}/frontend/review-form.html?token=${token}&place_id=${placeId}`;
        const qr = generateQRCode(token);

        html += `
            <div class="mb-4 p-3 border rounded">
                <strong>Токен ${i + 1}:</strong>
                <code class="d-block mt-1">${token}</code>

                <div class="mt-2">
                    <small class="text-muted">Ссылка для клиента:</small>
                    <div class="bg-light p-2 rounded small mt-1">${url}</div>
                </div>

                <div class="text-center mt-3">
                    <img src="${qr}" class="img-fluid border rounded" />
                </div>
            </div>
        `;
    });

    html += `</div></div>`;
    container.innerHTML = html;
}
