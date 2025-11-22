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
        const redeemBtn = document.getElementById('redeemBtn');
        if (stats.points >= 100) {
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
        const container = document.getElementById('bonusesContainer');

        if (!bonuses || bonuses.length === 0) {
            container.innerHTML = `
                <div class="text-center text-muted py-4">
                    <i class="bi bi-gift fs-1"></i>
                    <p class="mt-2">У вас пока нет бонусов</p>
                </div>
            `;
            return;
        }

        let html = '<div class="row">';

        bonuses.forEach(bonus => {
            const isActive = bonus.status === 'active';
            const isUsed = bonus.status === 'used';
            const isExpired = bonus.status === 'expired';

            html += `
                <div class="col-md-6 mb-3">
                    <div class="card h-100 ${isUsed ? 'border-secondary' : 'border-success'}">
                        <div class="card-body">
                            <div class="d-flex justify-content-between align-items-start">
                                <div>
                                    <h5 class="card-title">${bonus.title || 'Бонус'}</h5>
                                    <p class="card-text">${bonus.description || ''}</p>
                                    <small class="text-muted">
                                        Создан: ${new Date(bonus.created_at).toLocaleDateString()}
                                    </small>
                                    ${bonus.expires_at ? `
                                        <br><small class="text-muted">
                                            Действует до: ${new Date(bonus.expires_at).toLocaleDateString()}
                                        </small>
                                    ` : ''}
                                </div>
                                <span class="badge ${isActive ? 'bg-success' : isUsed ? 'bg-secondary' : 'bg-danger'}">
                                    ${isActive ? 'Активен' : isUsed ? 'Использован' : 'Просрочен'}
                                </span>
                            </div>
                        </div>
                        <div class="card-footer bg-transparent">
                            ${isActive ? `
                                <button class="btn btn-outline-primary btn-sm show-qr-btn" data-bonus-id="${bonus.id}">
                                    <i class="bi bi-qr-code me-1"></i>Показать QR
                                </button>
                            ` : ''}
                            ${isUsed ? `
                                <small class="text-muted">
                                    <i class="bi bi-check-circle me-1"></i>Использован
                                </small>
                            ` : ''}
                        </div>
                    </div>
                </div>
            `;
        });

        html += '</div>';
        container.innerHTML = html;

        // Навешиваем обработчики для QR кнопок
        document.querySelectorAll('.show-qr-btn').forEach(btn => {
            btn.addEventListener('click', function () {
                showBonusQR(this.getAttribute('data-bonus-id'));
            });
        });
    }

    // Показ QR кода бонуса
    function showBonusQR(bonusId) {
        const qrUrl = generateQRCode(`BONUS_${bonusId}`);

        document.getElementById('qrCodeImage').src = qrUrl;
        document.getElementById('bonusDescription').textContent = `Бонус ID: ${bonusId}`;

        const modal = new bootstrap.Modal(document.getElementById('qrModal'));
        modal.show();
    }

    // Обмен баллов на бонус
    async function redeemBonus() {
        const button = document.getElementById('redeemBtn');
        const messageDiv = document.getElementById('redeemMessage');

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
                body: JSON.stringify({})  // ← ДОБАВЬ ПУСТОЙ ОБЪЕКТ
            });

            if (!response.ok) {
                const errorData = await response.json();
                throw new Error(errorData.error || `Ошибка: ${response.status}`);
            }

            const result = await response.json();

            // Показываем успех
            showSuccess('Бонус успешно получен!', messageDiv);

            // Обновляем данные
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