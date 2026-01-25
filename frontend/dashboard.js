document.addEventListener("DOMContentLoaded", () => {
    const { requireUser, showError, showSuccess, logout, generateQRCode, apiFetch } = window.AppCommon;

    const token = requireUser();
    if (!token) return;

    initDashboard();

    function initDashboard() {
        // email
        const userEmail =
            localStorage.getItem("userEmail") ||
            localStorage.getItem("adminEmail") ||
            "";

        const emailElement = document.getElementById("userEmail");
        if (emailElement) {
            emailElement.textContent = userEmail; // покажет email или пусто, но без "пропуска"
        }

        // logout
        const logoutBtn = document.getElementById("logoutBtn");
        if (logoutBtn) logoutBtn.addEventListener("click", logout);

        // theme
        const themeBtn = document.getElementById("themeBtn");
        if (themeBtn) {
            themeBtn.addEventListener("click", () => {
                window.RL?.toggleTheme?.();
            });
        }

        // redeem
        const redeemBtn = document.getElementById("redeemBtn");
        if (redeemBtn) redeemBtn.addEventListener("click", redeemBonus);

        // ===== QR MODAL (custom, NO bootstrap) =====
        function openQrModal({ qrToken, title }) {
            const modalEl = document.getElementById("qrModal");
            if (!modalEl) return;

            const img = document.getElementById("qrCodeImage");
            const desc = document.getElementById("bonusDescription");
            const tokenEl = document.getElementById("bonusToken");

            if (img) img.src = generateQRCode(qrToken, 240);
            if (desc) desc.textContent = title || "Бонус";
            if (tokenEl) tokenEl.textContent = qrToken || "";

            modalEl.classList.add("show");
            modalEl.style.display = "block";
            modalEl.removeAttribute("aria-hidden");
            document.body.style.overflow = "hidden";
        }

        function closeQrModal() {
            const modalEl = document.getElementById("qrModal");
            if (!modalEl) return;

            modalEl.classList.remove("show");
            modalEl.style.display = "none";
            modalEl.setAttribute("aria-hidden", "true");
            document.body.style.overflow = "";
        }

        // open from QR button
        document.addEventListener("click", (e) => {
            const btn = e.target.closest("button[data-qr]");
            if (!btn) return;

            openQrModal({
                qrToken: btn.getAttribute("data-qr"),
                title: btn.getAttribute("data-title") || "Бонус",
            });
        });

        // close by X or backdrop click
        document.addEventListener("click", (e) => {
            const modalEl = document.getElementById("qrModal");
            if (!modalEl || !modalEl.classList.contains("show")) return;

            if (e.target.closest(".rl-modal__close")) {
                closeQrModal();
                return;
            }

            if (e.target === modalEl) {
                closeQrModal();
            }
        });

        // close by Esc
        document.addEventListener("keydown", (e) => {
            if (e.key !== "Escape") return;
            const modalEl = document.getElementById("qrModal");
            if (modalEl && modalEl.classList.contains("show")) closeQrModal();
        });



        loadUserStats();
        loadBonuses();


        async function loadUserStats() {
            try {
                const stats = await apiFetch("/users/stats", { method: "GET" });
                updateStatsUI(stats);
            } catch (e) {
                showError(`Не удалось загрузить статистику: ${e.message}`);
            }
        }

        function updateStatsUI(stats) {
            const totalReviewsEl = document.getElementById("totalReviews");
            const avgRatingEl = document.getElementById("avgRating");
            const pointsEl = document.getElementById("points");
            const currentPointsEl = document.getElementById("currentPoints");
            const redeemBtn = document.getElementById("redeemBtn");

            if (totalReviewsEl) totalReviewsEl.textContent = stats.total_reviews || 0;
            if (avgRatingEl) avgRatingEl.textContent = (stats.avg_rating || 0).toFixed(1);

            const points = stats.points || 0;
            if (pointsEl) pointsEl.textContent = points;
            if (currentPointsEl) currentPointsEl.textContent = points;

            const REQUIRED_POINTS = 50;
            if (redeemBtn) redeemBtn.disabled = points < REQUIRED_POINTS;
        }

        async function loadBonuses() {
            try {
                const bonuses = await apiFetch("/bonuses", { method: "GET" });
                displayBonuses(bonuses);
            } catch (e) {
                console.error("loadBonuses error:", e);
                showError(`Не удалось загрузить бонусы: ${e.message}`);
            }
        }

        function displayBonuses(bonuses) {
            const activeEl = document.getElementById("activeBonuses");
            const usedEl = document.getElementById("usedBonuses");
            const activeCountEl = document.getElementById("bonusesActive");

            if (!activeEl || !usedEl) return;

            const list = Array.isArray(bonuses) ? bonuses : [];

            const active = list.filter((b) => b && b.is_used === false);
            const used = list.filter((b) => b && b.is_used === true);

            if (activeCountEl) activeCountEl.textContent = active.length;

            activeEl.innerHTML = renderBonusesList(active, true);
            usedEl.innerHTML = renderBonusesList(used, false);
        }

        function renderBonusesList(items, isActive) {
            if (!Array.isArray(items) || items.length === 0) {
                return `
        <div class="center subtle" style="padding: var(--space-5);">
          ${isActive ? "Нет активных бонусов" : "Нет использованных бонусов"}
        </div>
      `;
            }

            return `
      <div class="stack">
        ${items.map((b) => renderBonusItem(b, isActive)).join("")}
      </div>
    `;
        }

        function renderBonusItem(b, isActive) {
            const title = rewardTypeLabel(b.reward_type);
            const req = typeof b.required_points === "number" ? b.required_points : null;

            const metaLeft = req != null ? `<span class="badge">${req} баллов</span>` : "";
            const metaRight = isActive ? `<span class="badge ok">Активен</span>` : `<span class="badge">Использован</span>`;

            const usedAt =
                !isActive && b.used_at
                    ? `<div class="subtle" style="font-size:12px; margin-top:6px;">Использован: ${formatDateTime(b.used_at)}</div>`
                    : "";

            const qrBtn =
                isActive && b.qr_token
                    ? `<button class="btn btn-sm" type="button" data-qr="${escapeAttr(b.qr_token)}" data-title="${escapeAttr(title)}">
             <i class="bi bi-qr-code"></i> QR
           </button>`
                    : "";

            return `
      <div class="card">
        <div class="card-body" style="display:flex; align-items:flex-start; justify-content:space-between; gap: var(--space-4);">
          <div style="min-width:0;">
            <div class="fw-semibold" style="font-size:14px; margin-bottom:6px;">${escapeHtml(title)}</div>
            <div style="display:flex; gap:10px; align-items:center; flex-wrap:wrap;">
              ${metaLeft}
              ${metaRight}
            </div>
            ${usedAt}
          </div>

          <div style="display:flex; gap:10px; align-items:center; flex-shrink:0;">
            ${qrBtn}
          </div>
        </div>
      </div>
    `;
        }

        function rewardTypeLabel(rt) {
            switch (rt) {
                case "free_coffee":
                    return "Бесплатный кофе";
                case "free_meal":
                    return "Бесплатное блюдо";
                case "discount_10":
                    return "Скидка 10%";
                default:
                    return rt || "Бонус";
            }
        }

        function formatDateTime(iso) {
            const d = new Date(iso);
            if (Number.isNaN(d.getTime())) return iso;
            return d.toLocaleString();
        }

        function escapeHtml(s) {
            return String(s)
                .replaceAll("&", "&amp;")
                .replaceAll("<", "&lt;")
                .replaceAll(">", "&gt;")
                .replaceAll('"', "&quot;");
        }

        function escapeAttr(s) {
            return escapeHtml(s).replaceAll("'", "&#39;");
        }

        async function redeemBonus() {
            const button = document.getElementById("redeemBtn");
            const messageDiv = document.getElementById("redeemMessage");

            if (!button || !messageDiv) return;

            const rewardSelect = document.getElementById("rewardType");
            const rewardType = rewardSelect ? rewardSelect.value : "free_coffee";

            const prevHtml = button.innerHTML;

            button.disabled = true;
            button.innerHTML = `<i class="bi bi-arrow-repeat"></i> Обмен...`;
            messageDiv.innerHTML = "";

            try {
                await apiFetch("/bonuses/redeem", {
                    method: "POST",
                    body: { reward_type: rewardType },
                });

                showSuccess("Бонус успешно получен!", messageDiv);

                // обновим данные
                await loadUserStats();
                await loadBonuses();
            } catch (error) {
                console.error("Ошибка обмена:", error);
                showError(error.message, messageDiv);
            } finally {
                button.disabled = false;
                button.innerHTML = prevHtml;
            }
        }
    }
});
