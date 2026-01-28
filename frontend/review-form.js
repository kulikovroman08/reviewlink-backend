document.addEventListener("DOMContentLoaded", () => {
	const { showError, showSuccess, apiFetch } = window.AppCommon;

	const urlParams = new URLSearchParams(window.location.search);

	const token = urlParams.get("token") || "";

	// Вариант 1: reviewlink place
	const placeId = urlParams.get("place_id") || urlParams.get("id") || "";

	// Вариант 2: public place (OSM/2GIS/и т.п.)
	const source = urlParams.get("source") || "";
	const sourceId = urlParams.get("source_id") || "";
	const placeName = urlParams.get("name") || "";

	// theme
	const themeBtn = document.getElementById("themeBtn");
	if (themeBtn) themeBtn.addEventListener("click", () => window.RL?.toggleTheme?.());

	// ✅ Разрешаем форму, если:
	// - есть placeId
	// ИЛИ
	// - есть source + sourceId (публичный каталог)
	const isPublic = !placeId && !!(source && sourceId);

	if (!placeId && !isPublic) {
		document.body.innerHTML = `
      <div class="container mt-4">
        <div class="card" style="max-width:640px; margin:0 auto;">
          <div class="card-header"><h5 class="section-title">Ошибка ссылки</h5></div>
          <div class="card-body center">
            <div class="muted">
              В ссылке отсутствует place_id и source/source_id.
              Откройте заведение из каталога заново или перезапустите QR.
            </div>
            <div style="margin-top: var(--space-4);">
              <a class="btn btn-sm" href="explore.html">В каталог</a>
              <a class="btn btn-sm" href="dashboard.html" style="margin-left:8px;">В личный кабинет</a>
            </div>
          </div>
        </div>
      </div>`;
		return;
	}

	const userToken = localStorage.getItem("userToken");
	if (!userToken) {
		window.location.href = `login.html?mode=login&redirect=${encodeURIComponent(window.location.href)}`;
		return;
	}

	// Покажем режим
	const modeBadge = document.getElementById("modeBadge");
	if (modeBadge) {
		if (token) {
			modeBadge.textContent = "QR-отзыв: за этот отзыв будут начислены баллы";
		} else {
			modeBadge.textContent = "Публичный отзыв: баллы за отзыв не начисляются";
		}
	}

	// (опционально) показать название, если пришли из каталога
	const placeNameEl = document.getElementById("placeName");
	if (placeNameEl) {
		const title = placeId ? "" : (placeName ? `Заведение: ${placeName}` : "Заведение");
		placeNameEl.textContent = title;
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

			if (i <= selectedRating) star.classList.add("is-filled");
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

			// ✅ body теперь зависит от режима
			const body = {
				rating: selectedRating,
				content,
			};

			if (placeId) {
				body.place_id = placeId;
			} else {
				body.source = source;
				body.source_id = sourceId;
				if (placeName) body.name = placeName;
			}

			if (token) body.token = token; // QR режим (баллы), иначе публичный

			await apiFetch("/reviews", { method: "POST", body });

			showSuccess("Спасибо! Ваш отзыв отправлен.", formMessage);
			if (submitBtn) submitBtn.style.display = "none";

			setTimeout(() => {
				if (token) {
					window.location.href = "dashboard.html";
					return;
				}

				// после обычного публичного отзыва:
				if (placeId) {
					window.location.href = `explore.html?id=${encodeURIComponent(placeId)}`;
				} else {
					// public place без placeId — возвращаем в каталог
					window.location.href = `explore.html`;
				}
			}, 1500);
		} catch (err) {
			if (err && err.message === "too many reviews today") {
				showError("Вы уже оценили это заведение сегодня. Спасибо!", formMessage);
				setTimeout(() => (window.location.href = "dashboard.html"), 1500);
				return;
			}

			showError(err?.message || "Ошибка отправки", formMessage);

			if (submitBtn) {
				submitBtn.disabled = false;
				submitBtn.textContent = "Отправить отзыв";
			}
		} finally {
			isSubmitting = false;
		}
	}
});
