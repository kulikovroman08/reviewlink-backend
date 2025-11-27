document.addEventListener("DOMContentLoaded", () => {
    const { API_BASE, showError, showSuccess } = window.AppCommon;

    const urlParams = new URLSearchParams(window.location.search);
    const token = urlParams.get("token");
    const placeId = urlParams.get("place_id");

    if (!placeId) {
        document.body.innerHTML = `
        <div class="container mt-5">
            <div class="card shadow-sm">
                <div class="card-body text-center">
                    <h4 class="text-danger mb-3">Ошибка ссылки</h4>
                    <p class="text-muted">
                        В ссылке отсутствует place_id.<br>
                        Обратитесь к продавцу или перезапустите QR-код.
                    </p>
                    <a href="dashboard.html" class="btn btn-primary mt-3">В личный кабинет</a>
                </div>
            </div>
        </div>`;
        return;
    }

    // Если токена нет — показываем сообщение
    if (!token) {
        document.body.innerHTML = `
            <div class="container mt-5">
                <div class="card shadow-sm">
                    <div class="card-body text-center">
                        <h4 class="text-danger mb-3">Токен не найден</h4>
                        <p class="text-muted">
                            Страница доступна только при переходе через QR-код.
                        </p>
                        <a href="dashboard.html" class="btn btn-primary">Перейти в личный кабинет</a>
                    </div>
                </div>
            </div>
        `;
        return;
    }

    // Проверка авторизации
    const userToken = localStorage.getItem("userToken");
    if (!userToken) {
        window.location.href = `login.html?redirect=${encodeURIComponent(window.location.href)}`;
        return;
    }

    let selectedRating = 5;
    let isSubmitting = false;

    const ratingStars = document.getElementById("ratingStars");
    const submitBtn = document.getElementById("submitBtn");
    const formMessage = document.getElementById("formMessage");

    // Звёздочки
    ratingStars.addEventListener("click", (e) => {
        if (e.target.innerText === "★" || e.target.innerText === "☆") {
            const index = [...ratingStars.children].indexOf(e.target);
            selectedRating = index + 1;
            updateStars();
        }
    });

    function updateStars() {
        let html = "";
        for (let i = 1; i <= 5; i++) {
            html += `<span style="cursor:pointer;">${i <= selectedRating ? "★" : "☆"}</span>`;
        }
        ratingStars.innerHTML = html;
    }

    updateStars();

    // Отправка отзыва
    submitBtn.addEventListener("click", submitReview);

    async function submitReview() {
        if (isSubmitting) return;
        isSubmitting = true;

        submitBtn.disabled = true;
        submitBtn.innerHTML =
            `<div class="spinner-border spinner-border-sm me-2"></div>Отправка...`;
        formMessage.innerHTML = "";

        try {
            const res = await fetch(`${API_BASE}/reviews`, {
                method: "POST",
                headers: {
                    "Authorization": "Bearer " + userToken,
                    "Content-Type": "application/json"
                },
                body: JSON.stringify({
                    token: token,
                    place_id: placeId,
                    rating: selectedRating,
                    content: document.getElementById("reviewContent").value
                })
            });

            const data = await safeParseJSON(res);

            if (!res.ok) {
                if (data.error === "too many reviews today") {

                    showError("Вы уже оценили это заведение сегодня. Спасибо за отзыв!", formMessage);

                    // Разблокируем кнопку
                    submitBtn.disabled = false;
                    submitBtn.innerHTML = "Отправить";

                    // Через 3 секунды — переход в личный кабинет
                    setTimeout(() => {
                        window.location.href = "dashboard.html";
                    }, 5000);

                    return;
                }

                throw new Error(data.error || "Ошибка отправки");
            }

            showSuccess("Спасибо! Ваш отзыв отправлен.", formMessage);
            submitBtn.style.display = "none";

            setTimeout(() => {
                window.location.href = "dashboard.html";
            }, 3000);

        } catch (err) {
            showError(err.message, formMessage);
        } finally {
            isSubmitting = false;
        }
    }

    async function safeParseJSON(res) {
        try {
            return await res.json();
        } catch {
            return {};
        }
    }
});
