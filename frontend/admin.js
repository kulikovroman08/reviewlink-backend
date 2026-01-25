document.addEventListener("DOMContentLoaded", () => {
    const parseJwt = window.AppCommon ? window.AppCommon.parseJwt : null;

    const { API_BASE } = window.AppCommon || {};
    const ADMIN_API_BASE = API_BASE || window.location.origin;

    // ===== helpers =====
    function getToken() {
        // админка должна брать ТОЛЬКО adminToken
        return localStorage.getItem("adminToken") || "";
    }

    function isAdminToken(token) {
        if (!token) return false;
        // Используем глобальный parseJwt из window.AppCommon
        const payload = window.AppCommon ? window.AppCommon.parseJwt(token) : null;
        return payload && payload.role === "admin";
    }

    function authHeaders(extra = {}) {
        const token = getToken();
        // бек ожидает просто токен (НЕ Bearer)
        return token ? { ...extra, Authorization: token } : { ...extra };
    }

    async function safeParseJSON(res) {
        try {
            return await res.json();
        } catch {
            return {};
        }
    }

    function escapeHtml(str) {
        return String(str || "")
            .replaceAll("&", "&amp;")
            .replaceAll("<", "&lt;")
            .replaceAll(">", "&gt;")
            .replaceAll('"', "&quot;")
            .replaceAll("'", "&#039;");
    }

    function formatDate(iso) {
        if (!iso) return "";
        const d = new Date(iso);
        if (Number.isNaN(d.getTime())) return iso;
        return d.toLocaleString();
    }

    function stars(rating) {
        const r = Number(rating) || 0;
        let out = "";
        for (let i = 1; i <= 5; i++) out += i <= r ? "★" : "☆";
        return out;
    }

    // Генерация QR (для отзывов)
    function generateQRCode(token, placeId) {
        const publicBase = window.location.origin;
        const reviewUrl = `${publicBase}/frontend/review-form.html?token=${token}&place_id=${placeId}`;
        return `https://api.qrserver.com/v1/create-qr-code/?size=150x150&data=${encodeURIComponent(
            reviewUrl
        )}`;
    }

    // ===== Guard: если токен есть, но не админ — уводим в dashboard =====
    (function guardAdminPage() {
        const token = getToken();

        // Если нет токена или токен не админский - редирект на логин
        if (!token) {
            window.location.href = "login.html";
            return;
        }

        // Используем window.AppCommon.parseJwt напрямую
        const payload = window.AppCommon ? window.AppCommon.parseJwt(token) : null;

        if (!payload || payload.role !== "admin") {
            window.location.href = "login.html";
            return;
        }

        // Если авторизован как админ - показываем статус
        const dbg = document.getElementById("loginStatus");
        if (dbg) {
            dbg.textContent = "Токен admin OK ✅";
        }
    })();

    // ===== UI refs =====
    const loginBtn = document.getElementById("loginBtn");
    const emailInput = document.getElementById("emailInput");
    const passwordInput = document.getElementById("passwordInput");
    const statusDiv = document.getElementById("loginStatus");

    const generateBtn = document.getElementById("generateBtn");
    const loadReviewsBtn = document.getElementById("loadReviewsBtn");

    // ===== Авторизация =====
    if (loginBtn) {
        loginBtn.onclick = async function () {
            const email = (emailInput?.value || "").trim();
            const password = passwordInput?.value || "";

            if (!statusDiv) return;

            if (!email || !password) {
                statusDiv.innerHTML = '<div class="text-danger">Заполните email и пароль</div>';
                return;
            }

            statusDiv.innerHTML = '<div class="text-info">Вход в систему...</div>';
            loginBtn.disabled = true;

            try {
                const res = await fetch(`${ADMIN_API_BASE}/login`, {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ email, password }),
                });

                const data = await safeParseJSON(res);
                if (!res.ok) throw new Error(data.error || "Ошибка входа");

                const token = data.token || "";
                if (!token) throw new Error("Не пришёл token от /login");

                // проверяем роль
                const payload = window.AppCommon ? window.AppCommon.parseJwt(token) : null;
                if (!payload || payload.role !== "admin") {
                    statusDiv.innerHTML = `<div class="text-danger">У вас нет прав администратора. Роль: ${payload ? payload.role : 'нет роли'}</div>`;
                    loginBtn.disabled = false;  // <-- ВОССТАНАВЛИВАЕМ
                    return;
                }

                // сохраняем как ОБЩИЙ токен
                localStorage.setItem("adminToken", token);
                localStorage.setItem("adminEmail", email);

                statusDiv.innerHTML = '<div class="text-success">Успешный вход!</div>';

                // включаем кнопки
                if (generateBtn) generateBtn.disabled = false;
                if (loadReviewsBtn) loadReviewsBtn.disabled = false;

                addLogoutButton();

                await loadPlaces();
            } catch (e) {
                statusDiv.innerHTML = `<div class="text-danger">${escapeHtml(e.message)}</div>`;
            } finally {
                loginBtn.disabled = false;  // <-- ВОССТАНАВЛИВАЕМ В ЛЮБОМ СЛУЧАЕ
            }
        };
    }

    // ===== Автологин (если уже есть userToken и он админ) =====
    (async function autoLogin() {
        const token = getToken();
        const email = localStorage.getItem("adminEmail") || localStorage.getItem("userEmail") || "";

        if (!token) return;
        if (!isAdminToken(token)) return;
        if (loginBtn) loginBtn.disabled = true;

        if (emailInput && email) emailInput.value = email;

        if (statusDiv) statusDiv.innerHTML = '<div class="text-success">Авторизован</div>';
        if (generateBtn) generateBtn.disabled = false;
        if (loadReviewsBtn) loadReviewsBtn.disabled = false;

        // Добавляем кнопку выхода
        addLogoutButton();

        await loadPlaces();
    })();

    // ===== Загрузка заведений =====
    async function loadPlaces() {
        const select = document.getElementById("placeSelect");
        if (!select) return;

        select.innerHTML = `<option>Загрузка...</option>`;

        try {
            const res = await fetch(`${ADMIN_API_BASE}/places`, {
                method: "GET",
                headers: authHeaders(),
            });

            const places = await safeParseJSON(res);
            select.innerHTML = "";

            (places || []).forEach((place) => {
                const option = document.createElement("option");
                option.value = place.id;
                option.textContent = `${place.name} (${place.address})`;
                select.appendChild(option);
            });

            // при смене place — подгружаем отзывы
            select.addEventListener("change", () => {
                if (select.value) loadReviews(select.value);
            });

            if (select.value) loadReviews(select.value);
        } catch (e) {
            select.innerHTML = `<option value="">Ошибка загрузки</option>`;
        }
    }

    // ===== Генерация токенов =====
    if (generateBtn) {
        generateBtn.onclick = async function () {
            const placeId = document.getElementById("placeSelect")?.value;
            const count = document.getElementById("countInput")?.value;

            if (!placeId) {
                alert("Выберите заведение");
                return;
            }

            const container = document.getElementById("resultsContainer");
            if (container) container.innerHTML = `<div class="alert alert-info">Генерация токенов...</div>`;

            try {
                const res = await fetch(`${ADMIN_API_BASE}/admin/tokens`, {
                    method: "POST",
                    headers: authHeaders({ "Content-Type": "application/json" }),
                    body: JSON.stringify({ place_id: placeId, count: Number(count) }),
                });

                const data = await safeParseJSON(res);
                if (!res.ok) throw new Error(data.error || "Ошибка генерации");

                showResults(data, container, placeId);
            } catch (e) {
                if (container) container.innerHTML = `<div class="alert alert-danger">${escapeHtml(e.message)}</div>`;
            }
        };
    }

    function showResults(data, container, placeId) {
        if (!container) return;

        const tokens = data.tokens || data.Tokens || [];
        let html = `
      <div class="card">
        <div class="card-header bg-success text-white">
          Успешно сгенерировано ${tokens.length} токенов
        </div>
        <div class="card-body">
    `;

        tokens.forEach((token, i) => {
            const publicBase = window.location.origin;
            const url = `${publicBase}/frontend/review-form.html?token=${token}&place_id=${placeId}`;
            const qr = generateQRCode(token, placeId);

            html += `
        <div class="mb-4 p-3 border rounded bg-white">
          <strong>Токен ${i + 1}:</strong>
          <code class="d-block mt-1">${escapeHtml(token)}</code>

          <div class="mt-2">
            <small class="text-muted">Ссылка для клиента:</small>
            <div class="bg-light p-2 rounded small mt-1">${escapeHtml(url)}</div>
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

    // ===== Reviews (business replies only) =====
    if (loadReviewsBtn) {
        loadReviewsBtn.addEventListener("click", () => {
            const placeId = document.getElementById("placeSelect")?.value;
            if (placeId) loadReviews(placeId);
        });
    }

    async function loadReviews(placeId) {
        const statusEl = document.getElementById("reviewsStatus");
        const container = document.getElementById("reviewsContainer");

        if (statusEl) statusEl.textContent = "Загрузка отзывов...";
        if (container) container.innerHTML = "";

        try {
            const res = await fetch(`${ADMIN_API_BASE}/places/${placeId}/reviews`, {
                method: "GET",
                headers: authHeaders(),
            });

            const data = await safeParseJSON(res);
            if (!res.ok) throw new Error(data.error || `Ошибка загрузки отзывов: ${res.status}`);

            const reviews = Array.isArray(data) ? data : data.reviews || [];
            if (statusEl) statusEl.textContent = `Найдено отзывов: ${reviews.length}`;

            if (!container) return;

            if (reviews.length === 0) {
                container.innerHTML = `<div class="text-muted">Отзывов пока нет.</div>`;
                return;
            }

            container.innerHTML = reviews.map(renderReviewCard).join("");
            wireReplyHandlers(placeId);
        } catch (e) {
            if (statusEl) statusEl.innerHTML = `<span class="text-danger">${escapeHtml(e.message)}</span>`;
        }
    }

    function normalizeReview(r) {
        const id = r.id || r.review_id;
        const author = r.authorName || r.author_name || r.user_name || "Пользователь";
        const content = r.content || "";
        const rating = r.rating || 0;
        const createdAt = r.createdAt || r.created_at;

        const replyObj = r.reply || r.ownerReply || r.owner_reply || null;

        const replyContent = replyObj ? (replyObj.content || replyObj.reply || "") : "";
        const replyAt = replyObj
            ? (replyObj.updated_at || replyObj.updatedAt || replyObj.created_at || replyObj.createdAt || "")
            : "";

        const hasReply = !!replyContent && replyContent.trim().length > 0;

        return { id, author, content, rating, createdAt, replyContent, replyAt, hasReply };
    }

    function renderReviewCard(r0) {
        const r = normalizeReview(r0);

        return `
      <div class="card mb-3" data-review-id="${escapeHtml(r.id)}" data-has-reply="${r.hasReply ? "1" : "0"}">
        <div class="card-body">

          <div class="d-flex justify-content-between align-items-start">
            <div>
              <div class="fw-semibold">${escapeHtml(r.author)}</div>
              <div class="text-muted small">${escapeHtml(formatDate(r.createdAt))}</div>
            </div>
            <div class="fs-5" title="Рейтинг">${stars(r.rating)}</div>
          </div>

          <div class="mt-2">${escapeHtml(r.content)}</div>

          <hr class="my-3" />

          <div class="p-3 rounded bg-light">

            <div class="d-flex align-items-center gap-2 mb-2">
              <div class="fw-semibold">Ответ заведения</div>
              <span class="badge bg-success js-replied-badge ${r.hasReply ? "" : "d-none"}">Ответ добавлен</span>
            </div>

            <div class="js-reply-view ${r.hasReply ? "" : "d-none"}">
              <div class="js-reply-content">${escapeHtml(r.replyContent)}</div>
              <div class="text-muted small mt-1 js-reply-date">${escapeHtml(formatDate(r.replyAt))}</div>
            </div>

            <div class="js-reply-empty ${r.hasReply ? "d-none" : ""} text-muted">
              Ответа пока нет.
            </div>

            <div class="mt-3">
              <textarea class="form-control form-control-sm js-reply-text" rows="2"
                placeholder="Напишите ответ...">${escapeHtml(r.replyContent)}</textarea>

              <div class="d-flex gap-2 mt-2">
                <button class="btn btn-sm btn-success js-reply-save" type="button">
                  ${r.hasReply ? "Обновить ответ" : "Ответить"}
                </button>
              </div>

              <div class="small mt-2 js-reply-status"></div>
            </div>

          </div>
        </div>
      </div>
    `;
    }

    function wireReplyHandlers(placeId) {
        document.querySelectorAll("#reviewsContainer .card[data-review-id]").forEach((card) => {
            const reviewId = card.getAttribute("data-review-id");
            const saveBtn = card.querySelector(".js-reply-save");

            saveBtn.addEventListener("click", async () => {
                const text = card.querySelector(".js-reply-text").value.trim();
                await saveReply(card, reviewId, text);
            });
        });
    }

    async function saveReply(card, reviewId, content) {
        const status = card.querySelector(".js-reply-status");
        status.innerHTML = "";

        if (!content) {
            status.innerHTML = `<span class="text-danger">Текст ответа не может быть пустым</span>`;
            return;
        }

        const btn = card.querySelector(".js-reply-save");
        btn.disabled = true;
        btn.innerText = "Сохранение...";

        try {
            let res = await fetch(`${ADMIN_API_BASE}/admin/reviews/${reviewId}/reply`, {
                method: "POST",
                headers: authHeaders({ "Content-Type": "application/json" }),
                body: JSON.stringify({ content }),
            });

            let data = await safeParseJSON(res);

            if (!res.ok) {
                res = await fetch(`${ADMIN_API_BASE}/admin/reviews/${reviewId}/reply`, {
                    method: "PUT",
                    headers: authHeaders({ "Content-Type": "application/json" }),
                    body: JSON.stringify({ content }),
                });
                data = await safeParseJSON(res);
            }

            if (!res.ok) throw new Error(data.error || `Ошибка сохранения ответа: ${res.status}`);

            const wasReply = card.getAttribute("data-has-reply") === "1";
            status.innerHTML = `<span class="text-success">${wasReply ? "Ответ обновлён" : "Ответ сохранён"}</span>`;

            card.setAttribute("data-has-reply", "1");

            const view = card.querySelector(".js-reply-view");
            const empty = card.querySelector(".js-reply-empty");
            const contentEl = card.querySelector(".js-reply-content");
            const dateEl = card.querySelector(".js-reply-date");
            const badge = card.querySelector(".js-replied-badge");

            if (contentEl) contentEl.textContent = content;
            if (view) view.classList.remove("d-none");
            if (empty) empty.classList.add("d-none");
            if (badge) badge.classList.remove("d-none");
            if (dateEl) dateEl.textContent = new Date().toLocaleString();

            btn.innerText = "Обновить ответ";

            setTimeout(() => {
                status.innerHTML = "";
            }, 2500);
        } catch (e) {
            status.innerHTML = `<span class="text-danger">${escapeHtml(e.message)}</span>`;
        } finally {
            btn.disabled = false;
            btn.innerText =
                card.getAttribute("data-has-reply") === "1" ? "Обновить ответ" : "Ответить";
        }
    }
    // ===== Кнопка выхода (как в dashboard) =====
    function addLogoutButton() {
        const navRight = document.querySelector(".nav-right");
        if (!navRight) return;

        // уже есть
        if (document.querySelector("#adminLogoutBtn")) return;

        const logoutBtn = document.createElement("button");
        logoutBtn.id = "adminLogoutBtn";

        // стиль как у "Тема" / dashboard (не danger)
        logoutBtn.className = "btn btn-sm";
        logoutBtn.type = "button";

        // иконка + текст (inline svg — без библиотек)
        logoutBtn.innerHTML = `
      <span style="display:inline-flex;align-items:center;gap:8px;">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none"
             xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
          <path d="M10 7V6a2 2 0 0 1 2-2h6a2 2 0 0 1 2 2v12a2 2 0 0 1-2 2h-6a2 2 0 0 1-2-2v-1"
                stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
          <path d="M14 12H3m0 0 3-3M3 12l3 3"
                stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
        </svg>
        <span>Выйти</span>
      </span>
    `;

        logoutBtn.onclick = function () {
            localStorage.removeItem("adminToken");
            localStorage.removeItem("adminEmail");
            window.location.href = "login.html";
        };

        // вставляем сразу после кнопки "Тема" (как на dashboard)
        const themeBtn = navRight.querySelector('button[onclick*="toggleTheme"]');
        if (themeBtn && themeBtn.nextSibling) {
            navRight.insertBefore(logoutBtn, themeBtn.nextSibling);
        } else if (themeBtn) {
            navRight.appendChild(logoutBtn);
        } else {
            navRight.insertBefore(logoutBtn, navRight.firstChild);
        }
    }

});
