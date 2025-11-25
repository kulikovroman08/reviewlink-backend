const API_BASE = (window.location.port === "5500" || window.location.port === "8000")
    ? "http://localhost:8080"
    : window.location.origin;

console.log("🌐 Порт:", window.location.port);
console.log("🌐 API_BASE:", API_BASE);

let BEARER_TOKEN = "";


function generateQRCode(token) {
    const reviewUrl = `${API_BASE}/frontend/review-form.html?token=${token}`;
    return `https://api.qrserver.com/v1/create-qr-code/?size=150x150&data=${encodeURIComponent(reviewUrl)}`;
}


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

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || `Ошибка: ${response.status}`);
        }

        const data = await response.json();

        BEARER_TOKEN = data.token;
        localStorage.setItem("adminToken", BEARER_TOKEN);
        localStorage.setItem("adminEmail", email);

        statusDiv.innerHTML = '<div class="text-success">✅ Успешный вход!</div>';
        document.getElementById("generateBtn").disabled = false;

        // Загружаем список мест
        await loadPlaces();

    } catch (error) {
        statusDiv.innerHTML = `<div class="text-danger">❌ ${error.message}</div>`;
    }
};


document.addEventListener("DOMContentLoaded", async function () {
    const savedToken = localStorage.getItem("adminToken");
    const savedEmail = localStorage.getItem("adminEmail");

    if (savedToken && savedEmail) {
        BEARER_TOKEN = savedToken;
        document.getElementById("emailInput").value = savedEmail;
        document.getElementById("loginStatus").innerHTML = '<div class="text-success">✅ Авторизован</div>';
        document.getElementById("generateBtn").disabled = false;

        await loadPlaces();
    }
});


async function loadPlaces() {
    const select = document.getElementById("placeSelect");
    select.innerHTML = `<option>Загрузка...</option>`;

    try {
        const response = await fetch(`${API_BASE}/places`, {
            method: "GET",
            headers: {
                "Authorization": "Bearer " + BEARER_TOKEN,
                "Content-Type": "application/json"
            }
        });

        if (!response.ok) {
            throw new Error("Не удалось загрузить список мест");
        }

        const places = await response.json();

        select.innerHTML = "";

        places.forEach(place => {
            const option = document.createElement("option");
            option.value = place.id;
            option.textContent = `${place.name} (${place.address})`;
            select.appendChild(option);
        });

    } catch (error) {
        console.error(error);
        select.innerHTML = `<option value="">Ошибка загрузки</option>`;
    }
}


document.getElementById("generateBtn").onclick = async function () {
    const placeId = document.getElementById("placeSelect").value;
    const count = document.getElementById("countInput").value;

    if (!placeId) {
        alert("Выберите заведение");
        return;
    }

    const container = document.getElementById("resultsContainer");
    container.innerHTML = `
        <div class="alert alert-info d-flex align-items-center">
            <div class="spinner-border spinner-border-sm me-2"></div>
            <div>Генерация токенов...</div>
        </div>
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

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || `Ошибка: ${response.status}`);
        }

        const data = await response.json();
        showResults(data, container);

    } catch (error) {
        container.innerHTML = `
            <div class="alert alert-danger">
                <strong>Ошибка:</strong> ${error.message}
            </div>
        `;
    }
};


function showResults(data, container) {
    const tokens = data.tokens || data.Tokens || [];

    if (tokens.length === 0) {
        container.innerHTML = '<div class="alert alert-warning">Токены не сгенерированы</div>';
        return;
    }

    let html = `
        <div class="card">
            <div class="card-header bg-success text-white">
                <strong>✅ Успешно сгенерировано ${tokens.length} токенов</strong>
            </div>
            <div class="card-body">
    `;

    tokens.forEach((token, index) => {
        const url = `${API_BASE}/frontend/review-form.html?token=${token}`;
        const qr = generateQRCode(token);

        html += `
            <div class="mb-4 p-3 border rounded">
                <div class="row">
                    <div class="col-md-8">
                        <strong>Токен ${index + 1}:</strong> 
                        <code class="ms-2 bg-light p-2 rounded d-block mt-1">${token}</code>
                        <div class="mt-2">
                            <small class="text-muted">Ссылка для клиента:</small>
                            <div class="bg-light p-2 rounded small mt-1">
                                ${url}
                            </div>
                        </div>
                    </div>
                    <div class="col-md-4 text-center">
                        <img src="${qr}" alt="QR Code" class="img-fluid border rounded">
                        <div class="mt-1">
                            <small class="text-muted">QR код для чека</small>
                        </div>
                    </div>
                </div>
            </div>
        `;
    });

    html += `
            </div>
        </div>
    `;

    container.innerHTML = html;
}