// dashboard.js - логика личного кабинета
document.addEventListener('DOMContentLoaded', function () {
    const { API_BASE, checkAuth, showError, showSuccess, logout, generateQRCode } = window.AppCommon;

    let USER_TOKEN = checkAuth();
    if (!USER_TOKEN) return;

    // Инициализация
    initDashboard();

    function initDashboard() {
        // Показываем email пользователя
        const userEmail = localStorage.getItem('userEmail');
        if (userEmail) {
            const emailElement = document.getElementById('userEmail');
            if (emailElement) {
                emailElement.textContent = userEmail;
            }
        }

        // Обработчик для кнопки выхода
        const logoutBtn = document.getElementById('logoutBtn');
        if (logoutBtn) {
            logoutBtn.addEventListener('click', function () {
                window.AppCommon.logout();
            });
        }

        const redeemBtn = document.getElementById('redeemBtn');
        if (redeemBtn) {
            redeemBtn.addEventListener('click', redeemBonus);
        }

        // Загружаем данные
        loadUserStats();
        loadBonuses();
    }

    // Загрузка статистики
    async function loadUserStats() {
        try {
            const response = await fetch(`${API_BASE}/users/stats`, {
                method: "GET",
                headers: {
                    "Authorization": "Bearer " + USER_TOKEN,
                    "Content-Type": "application/json"
                }
            });

            if (!response.ok) {
                if (response.status === 401) {
                    logout();
                    return;
                }
                throw new Error(`Ошибка: ${response.status}`);
            }

            const stats = await response.json();
            updateStatsUI(stats);

        } catch (error) {
            console.error('Ошибка загрузки статистики:', error);
            showError('Не удалось загрузить статистику');
        }
    }

    // Обновление UI статистики
    function updateStatsUI(stats) {
        document.getElementById('totalReviews').textContent = stats.total_reviews || 0;
        document.getElementById('avgRating').textContent = (stats.avg_rating || 0).toFixed(1);
        document.getElementById('points').textContent = stats.points || 0;
        document.getElementById('bonusesActive').textContent = stats.bonuses_active || 0;
        document.getElementById('currentPoints').textContent = stats.points || 0;

        // Активируем кнопку если достаточно баллов
        const REQUIRED_POINTS = 50;

        if (stats.points >= REQUIRED_POINTS) {
            redeemBtn.disabled = false;
            redeemBtn.innerHTML = '<i class="bi bi-gift me-2"></i>Получить бонус';
        } else {
            redeemBtn.disabled = true;
            redeemBtn.innerHTML = '<i class="bi bi-lock me-2"></i>Недостаточно баллов';
        }
    }

    // Загрузка бонусов
    async function loadBonuses() {
        try {
            console.log('Загрузка бонусов...');
            const response = await fetch(`${API_BASE}/bonuses`, {
                method: "GET",
                headers: {
                    "Authorization": "Bearer " + USER_TOKEN,
                    "Content-Type": "application/json"
                }
            });

            console.log('Статус ответа:', response.status);

            if (!response.ok) {
                throw new Error(`Ошибка: ${response.status}`);
            }

            const bonuses = await response.json();
            console.log('Получены бонусы:', bonuses);
            displayBonuses(bonuses);

        } catch (error) {
            console.error('Ошибка загрузки бонусов:', error);
            document.getElementById('bonusesContainer').innerHTML = `
                <div class="alert alert-danger">
                    Не удалось загрузить список бонусов
                </div>
            `;
        }
    }

    // Отображение бонусов
    function displayBonuses(bonuses) {
        const activeContainer = document.getElementById("activeBonuses");
        const usedContainer = document.getElementById("usedBonuses");

        if (!bonuses || bonuses.length === 0) {
            activeContainer.innerHTML = `
            <div class="text-center text-muted py-4">
                <i class="bi bi-gift fs-1"></i>
                <p class="mt-2">У вас пока нет бонусов</p>
            </div>
        `;
            usedContainer.innerHTML = `
            <div class="text-center text-muted py-4">
                Использованных бонусов нет
            </div>
        `;
            return;
        }

        let activeHTML = "";
        let usedHTML = "";

        bonuses.forEach(bonus => {
            const isUsed = bonus.is_used === true;

            const cardHTML = `
            <div class="col-md-6 mb-3">
                <div class="card h-100 border-${isUsed ? 'secondary' : 'success'}">
                    <div class="card-body">
                        <h5 class="card-title">${bonus.reward_type}</h5>
                        <p class="card-text">Списано баллов: <strong>${bonus.required_points}</strong></p>
                        <p class="card-text">QR токен: ${bonus.qr_token}</p>

                        <span class="badge ${isUsed ? 'bg-secondary' : 'bg-success'}">
                            ${isUsed ? 'Использован' : 'Активен'}
                        </span>
                    </div>

                    <div class="card-footer bg-transparent">
                        ${!isUsed
                    ? `<button class="btn btn-outline-primary btn-sm show-qr-btn"
                                     data-bonus-token="${bonus.qr_token}">
                                        <i class="bi bi-qr-code me-1"></i>Показать QR
                                   </button>`
                    : `<small class="text-muted">
                                        <i class="bi bi-check-circle me-1"></i>Использован
                                   </small>`
                }
                    </div>
                </div>
            </div>
        `;

            // Разделение
            if (isUsed) {
                usedHTML += cardHTML;
            } else {
                activeHTML += cardHTML;
            }
        });

        // Рендер
        activeContainer.innerHTML = `<div class="row">${activeHTML}</div>`;
        usedContainer.innerHTML = `<div class="row">${usedHTML}</div>`;

        // Навешиваем обработчики
        document.querySelectorAll(".show-qr-btn").forEach(btn => {
            btn.addEventListener("click", function () {
                showBonusQR(this.getAttribute("data-bonus-token"));
            });
        });
    }

    // Показ QR кода бонуса
    function showBonusQR(qrToken) {
        const qrUrl = generateQRCode(qrToken);

        document.getElementById('qrCodeImage').src = qrUrl;
        document.getElementById('bonusDescription').textContent = `QR токен: ${qrToken}`;

        const modal = new bootstrap.Modal(document.getElementById('qrModal'));
        modal.show();
    }

    // Обмен баллов на бонус
    async function redeemBonus() {
        const button = document.getElementById('redeemBtn');
        const messageDiv = document.getElementById('redeemMessage');

        const rewardType = document.getElementById("rewardType").value;

        button.disabled = true;
        button.innerHTML = '<div class="spinner-border spinner-border-sm me-2"></div>Обмен...';
        messageDiv.innerHTML = '';

        try {
            const response = await fetch(`${API_BASE}/bonuses/redeem`, {
                method: "POST",
                headers: {
                    "Authorization": "Bearer " + USER_TOKEN,
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    reward_type: rewardType
                })
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || `Ошибка: ${response.status}`);
            }

            showSuccess('Бонус успешно получен!', messageDiv);

            setTimeout(() => {
                loadUserStats();
                loadBonuses();
            }, 1000);

        } catch (error) {
            console.error('Ошибка обмена:', error);
            showError(error.message, messageDiv);
        } finally {
            setTimeout(() => {
                button.disabled = false;
                button.innerHTML = '<i class="bi bi-gift me-2"></i>Получить бонус';
            }, 2000);
        }
    }
});