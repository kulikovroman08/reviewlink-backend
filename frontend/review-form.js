document.addEventListener("DOMContentLoaded", () => {
    const { showError, showSuccess, apiFetch } = window.AppCommon;

    const urlParams = new URLSearchParams(window.location.search);
    const token = urlParams.get("token");
    const placeId = urlParams.get("place_id");

    // theme
    const themeBtn = document.getElementById("themeBtn");
    if (themeBtn) themeBtn.addEventListener("click", () => window.RL?.toggleTheme?.());

    if (!placeId) {
        document.body.innerHTML = `
      <div class="container mt-4">
        <div class="card" style="max-width:640px; margin:0 auto;">
          <div class="card-header"><h5 class="section-title">Ошибка ссылки</h5></div>
          <div class="card-body center">
            <div class="muted">В ссылке отсутствует place_id. Обратитесь к заведению или перезапустите QR.</div>
            <div style="margin-top: var(--space-4);">
              <a class="btn btn-sm" href="dashboard.html">В личный кабинет</a>
            </div>
          </div>
        </div>
      </div>`;
        return;
    }

    if (!token) {
        document.body.innerHTML = `
      <div class="container mt-4">
        <div class="card" style="max-width:640px; margin:0 auto;">
          <div class="card-header"><h5 class="section-title">Токен не найден</h5></div>
          <div class="card-body center">
            <div class="muted">Страница доступна только при переходе через QR-код.</div>
            <div style="margin-top: var(--space-4);">
              <a class="btn btn-sm" href="dashboard.html">Перейти в личный кабинет</a>
            </div>
          </div>
        </div>
      </div>`;
        return;
    }

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

    function renderStars() {
        if (!ratingStars) return;

        ratingStars.innerHTML = "";
        for (let i = 1; i <= 5; i++) {
            const star = document.createElement("button");
            star.type = "button";
            star.className = "rl-star";
            star.dataset.value = String(i);
            star.setAttribute("aria-label", `Оценка ${i}`);
            star.setAttribute("aria-pressed", i === selectedRating ? "true" : "false");

            const filled = i <= selectedRating;
            if (filled) star.classList.add("is-filled");

            star.textContent = "★";

            ratingStars.appendChild(star);
        }
    }

    ratingStars?.addEventListener("click", (e) => {
        const btn = e.target.closest(".rl-star");
        if (!btn) return;
        const val = Number(btn.dataset.value);
        if (!Number.isFinite(val)) return;
        selectedRating = val;
        renderStars();
    });

    // hover подсветка
    ratingStars?.addEventListener("mouseover", (e) => {
        const btn = e.target.closest(".rl-star");
        if (!btn) return;

        const val = Number(btn.dataset.value);
        if (!Number.isFinite(val)) return;

        ratingStars.querySelectorAll(".rl-star").forEach((s) => {
            const v = Number(s.dataset.value);
            s.classList.toggle("is-hover", v <= val);
        });
    });

    ratingStars?.addEventListener("mouseout", () => {
        ratingStars.querySelectorAll(".rl-star").forEach((s) => s.classList.remove("is-hover"));
    });


    renderStars();

    submitBtn?.addEventListener("click", submitReview);

    async function submitReview() {
        if (isSubmitting) return;
        isSubmitting = true;

        if (submitBtn) {
            submitBtn.disabled = true;
            submitBtn.textContent = "Отправка...";
        }
        if (formMessage) formMessage.innerHTML = "";

        try {
            const contentEl = document.getElementById("reviewContent");
            const content = contentEl ? String(contentEl.value || "").trim() : "";

            await apiFetch("/reviews", {
                method: "POST",
                body: {
                    token,
                    place_id: placeId,
                    rating: selectedRating,
                    content,
                },
            });

            showSuccess("Спасибо! Ваш отзыв отправлен.", formMessage);

            if (submitBtn) submitBtn.style.display = "none";

            setTimeout(() => {
                window.location.href = "dashboard.html";
            }, 2000);
        } catch (err) {
            if (err && err.message === "too many reviews today") {
                showError("Вы уже оценили это заведение сегодня. Спасибо!", formMessage);
                setTimeout(() => (window.location.href = "dashboard.html"), 2500);
                return;
            }

            showError(err?.message || "Ошибка отправки", formMessage);

            if (submitBtn) {
                submitBtn.disabled = false;
                submitBtn.innerHTML = `<i class="bi bi-send"></i> Отправить отзыв`;
            }
        } finally {
            isSubmitting = false;
        }
    }
});
