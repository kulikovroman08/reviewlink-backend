document.addEventListener("DOMContentLoaded", () => {
  const { API_BASE } = window.AppCommon;

  // ====== DOM ======
  const topView = document.getElementById("topView");
  const placeView = document.getElementById("placeView");

  const topStatus = document.getElementById("topStatus");
  const topContainer = document.getElementById("topContainer");
  const refreshTopBtn = document.getElementById("refreshTopBtn");

  const limitInput = document.getElementById("limitInput");
  const placesSort = document.getElementById("placesSort");
  const minRating = document.getElementById("minRating");
  const minReviews = document.getElementById("minReviews");

  const placesSortWrap = document.getElementById("placesSortWrap");
  const minRatingWrap = document.getElementById("minRatingWrap");
  const minReviewsWrap = document.getElementById("minReviewsWrap");

  const backToTopBtn = document.getElementById("backToTopBtn");
  const refreshReviewsBtn = document.getElementById("refreshReviewsBtn");

  const placeTitle = document.getElementById("placeTitle");
  const reviewsStatus = document.getElementById("reviewsStatus");
  const reviewsContainer = document.getElementById("reviewsContainer");

  const sortSelect = document.getElementById("sortSelect");
  const ratingSelect = document.getElementById("ratingSelect");

  const authPill = document.getElementById("authPill");
  const voteHint = document.getElementById("voteHint");

  // ====== state ======
  let activeTab = "places"; // places | users | bonuses
  let currentPlaceId = getPlaceIdFromUrl();
  let currentPlaceName = ""; // чтобы не показывать ID

  // ====== helpers ======
  function getUserToken() {
    return localStorage.getItem("userToken") || "";
  }
  function getUserEmail() {
    return localStorage.getItem("userEmail") || localStorage.getItem("adminEmail") || "";
  }
  function isLoggedIn() {
    return Boolean(getUserToken());
  }

  function updateAuthUI() {
    if (!authPill) return;
    if (isLoggedIn()) {
      const email = getUserEmail();
      authPill.textContent = email ? email : "Вход выполнен";
      authPill.classList.add("ok");
      authPill.classList.remove("bad", "warn");
      if (voteHint) voteHint.textContent = "Вы вошли — можно голосовать 👍 / 👎.";
    } else {
      authPill.textContent = "Гость";
      authPill.classList.add("warn");
      authPill.classList.remove("ok", "bad");
      if (voteHint) voteHint.textContent = "Чтобы голосовать, нужно войти. Клик по 👍/👎 перекинет на логин.";
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
    if (Number.isNaN(d.getTime())) return String(iso);
    return d.toLocaleString();
  }

  function stars(rating) {
    const r = Number(rating) || 0;
    let out = "";
    for (let i = 1; i <= 5; i++) out += i <= r ? "★" : "☆";
    return out;
  }

  function getPlaceIdFromUrl() {
    const p = new URLSearchParams(window.location.search);
    return p.get("id") || p.get("place_id") || "";
  }

  function setPlaceIdToUrl(placeId) {
    const url = new URL(window.location.href);
    url.searchParams.set("id", placeId);
    history.pushState({}, "", url.toString());
  }

  function clearPlaceIdFromUrl() {
    const url = new URL(window.location.href);
    url.searchParams.delete("id");
    url.searchParams.delete("place_id");
    history.pushState({}, "", url.toString());
  }

  async function safeParseJSON(res) {
    try {
      return await res.json();
    } catch {
      return {};
    }
  }

  function clampInt(v, def, min, max) {
    const n = Number.parseInt(v, 10);
    if (Number.isNaN(n)) return def;
    return Math.max(min, Math.min(max, n));
  }
  function clampFloat(v, def, min, max) {
    const n = Number.parseFloat(v);
    if (Number.isNaN(n)) return def;
    return Math.max(min, Math.min(max, n));
  }

  function currentReviewsQueryString() {
    const qs = new URLSearchParams();
    if (sortSelect && sortSelect.value) qs.set("sort", sortSelect.value);
    if (ratingSelect && ratingSelect.value) qs.set("rating", ratingSelect.value);
    return qs.toString();
  }

  // ====== tabs ======
  document.querySelectorAll("#topTabs .nav-link").forEach((btn) => {
    btn.addEventListener("click", async () => {
      document.querySelectorAll("#topTabs .nav-link").forEach((b) => b.classList.remove("active"));
      btn.classList.add("active");
      activeTab = btn.getAttribute("data-tab");

      const isPlaces = activeTab === "places";
      placesSortWrap.classList.toggle("d-none", !isPlaces);
      minRatingWrap.classList.toggle("d-none", !isPlaces);
      minReviewsWrap.classList.toggle("d-none", !isPlaces);

      await loadTop();
    });
  });

  refreshTopBtn.addEventListener("click", loadTop);

  // ====== top loaders ======
  async function loadTop() {
    topStatus.textContent = "Загрузка...";
    topContainer.innerHTML = "";

    const limit = clampInt(limitInput.value, 10, 1, 100);

    try {
      if (activeTab === "places") {
        const sortBy = placesSort.value || "reviews";
        const minR = clampFloat(minRating.value, 0, 0, 5);
        const minRev = clampInt(minReviews.value, 0, 0, 100000);

        const qs = new URLSearchParams({
          limit: String(limit),
          sort_by: sortBy,
          min_rating: String(minR),
          min_reviews: String(minRev),
        });

        const res = await fetch(`${API_BASE}/leaderboard/places?${qs.toString()}`);
        const data = await safeParseJSON(res);
        if (!res.ok) throw new Error(data.error || `Ошибка leaderboard places: ${res.status}`);

        topStatus.textContent = `Заведений: ${data.length}`;
        topContainer.innerHTML = renderPlacesTop(data);
        wirePlacesTop();
        return;
      }

      if (activeTab === "users") {
        const qs = new URLSearchParams({ limit: String(limit), sort_by: "reviews" });
        const res = await fetch(`${API_BASE}/leaderboard/users?${qs.toString()}`);
        const data = await safeParseJSON(res);
        if (!res.ok) throw new Error(data.error || `Ошибка leaderboard users: ${res.status}`);

        topStatus.textContent = `Пользователей: ${data.length}`;
        topContainer.innerHTML = renderUsersTop(data);
        return;
      }

      if (activeTab === "bonuses") {
        const res = await fetch(`${API_BASE}/leaderboard/bonuses`);
        const data = await safeParseJSON(res);
        if (!res.ok) throw new Error(data.error || `Ошибка leaderboard bonuses: ${res.status}`);

        topStatus.textContent = `Пользователей: ${data.length}`;
        topContainer.innerHTML = renderBonusesTop(data);
        return;
      }
    } catch (e) {
      topStatus.innerHTML = `<span class="text-danger">${escapeHtml(e.message)}</span>`;
    }
  }

  // ✅ Убираем ID заведения из таблицы (п.4)
  function renderPlacesTop(items) {
    return `
      <table>
        <thead>
          <tr>
            <th style="width:80px">#</th>
            <th>Заведение</th>
            <th style="width:140px; text-align:right;">Отзывы</th>
            <th style="width:140px; text-align:right;">Рейтинг</th>
            <th style="width:160px; text-align:right;"></th>
          </tr>
        </thead>
        <tbody>
          ${items
        .map(
          (p) => `
              <tr data-place-id="${escapeHtml(p.id)}" data-place-name="${escapeHtml(p.name)}" class="js-place-row" style="cursor:pointer">
                <td>${p.rank}</td>
                <td>
                  <div class="fw-semibold">${escapeHtml(p.name)}</div>
                </td>
                <td style="text-align:right;">${p.reviews_count}</td>
                <td style="text-align:right;">${Number(p.avg_rating || 0).toFixed(2)}</td>
                <td style="text-align:right;">
                  <button class="btn btn-sm js-open-place" type="button" data-place-id="${escapeHtml(p.id)}" data-place-name="${escapeHtml(p.name)}">
                    Отзывы →
                  </button>
                </td>
              </tr>
            `
        )
        .join("")}
        </tbody>
      </table>
    `;
  }

  function wirePlacesTop() {
    topContainer.querySelectorAll(".js-open-place").forEach((btn) => {
      btn.addEventListener("click", (e) => {
        e.stopPropagation();
        openPlace(btn.getAttribute("data-place-id"), btn.getAttribute("data-place-name") || "");
      });
    });

    topContainer.querySelectorAll(".js-place-row").forEach((row) => {
      row.addEventListener("click", () =>
        openPlace(row.getAttribute("data-place-id"), row.getAttribute("data-place-name") || "")
      );
    });
  }

  function renderUsersTop(items) {
    return `
      <table>
        <thead>
          <tr>
            <th style="width:80px">#</th>
            <th>Пользователь</th>
            <th style="width:140px; text-align:right;">Отзывы</th>
            <th style="width:140px; text-align:right;">Средн. рейтинг</th>
          </tr>
        </thead>
        <tbody>
          ${items
        .map(
          (u) => `
              <tr>
                <td>${u.rank}</td>
                <td>
                  <div class="fw-semibold">${escapeHtml(u.name)}</div>
                </td>
                <td style="text-align:right;">${u.reviews_count}</td>
                <td style="text-align:right;">${Number(u.avg_rating || 0).toFixed(2)}</td>
              </tr>
            `
        )
        .join("")}
        </tbody>
      </table>
    `;
  }

  function renderBonusesTop(items) {
    return `
      <table>
        <thead>
          <tr>
            <th style="width:80px">#</th>
            <th>Пользователь</th>
            <th style="width:160px; text-align:right;">Бонусов</th>
            <th style="width:160px; text-align:right;">Потрачено</th>
          </tr>
        </thead>
        <tbody>
          ${items
        .map(
          (u) => `
              <tr>
                <td>${u.rank}</td>
                <td><div class="fw-semibold">${escapeHtml(u.name)}</div></td>
                <td style="text-align:right;">${u.bonuses_count}</td>
                <td style="text-align:right;">${u.points_spent}</td>
              </tr>
            `
        )
        .join("")}
        </tbody>
      </table>
    `;
  }

  // ====== place reviews ======
  backToTopBtn.addEventListener("click", () => {
    currentPlaceId = "";
    currentPlaceName = "";
    clearPlaceIdFromUrl();
    showTop();
  });

  refreshReviewsBtn.addEventListener("click", () => {
    if (currentPlaceId) loadReviews(currentPlaceId);
  });

  sortSelect.addEventListener("change", () => {
    if (currentPlaceId) loadReviews(currentPlaceId);
  });

  ratingSelect.addEventListener("change", () => {
    if (currentPlaceId) loadReviews(currentPlaceId);
  });

  function openPlace(placeId, placeName = "") {
    currentPlaceId = placeId;
    currentPlaceName = placeName || currentPlaceName;
    setPlaceIdToUrl(placeId);
    showPlace();
    loadReviews(placeId);
  }

  function showTop() {
    placeView.classList.add("d-none");
    topView.classList.remove("d-none");
    loadTop();
  }

  function showPlace() {
    topView.classList.add("d-none");
    placeView.classList.remove("d-none");
    reviewsContainer.innerHTML = "";
    reviewsStatus.textContent = "Загрузка...";
    // ✅ не показываем ID (п.4)
    placeTitle.textContent = currentPlaceName ? `Отзывы: ${currentPlaceName}` : "Отзывы заведения";
    updateAuthUI();
  }

  async function loadReviews(placeId) {
    reviewsStatus.textContent = "Загрузка...";
    reviewsContainer.innerHTML = "";

    const qs = currentReviewsQueryString();

    try {
      const res = await fetch(`${API_BASE}/places/${placeId}/reviews?${qs}`);
      const data = await safeParseJSON(res);
      if (!res.ok) throw new Error(data.error || `Ошибка отзывов: ${res.status}`);

      const reviews = Array.isArray(data) ? data : data.reviews || [];
      reviewsStatus.textContent = `Найдено отзывов: ${reviews.length}`;

      if (reviews.length === 0) {
        reviewsContainer.innerHTML = `<div class="muted">Отзывов пока нет.</div>`;
        return;
      }

      reviewsContainer.innerHTML = reviews.map(renderReviewCard).join("");
      wireVotes();
      updateAuthUI();
    } catch (e) {
      reviewsStatus.innerHTML = `<span class="text-danger">${escapeHtml(e.message)}</span>`;
    }
  }

  function normalizeReview(r) {
    return {
      id: r.id,
      rating: r.rating,
      content: r.content || "",
      createdAt: r.created_at || r.createdAt,
      helpful: r.helpful_count ?? r.helpful ?? 0,
      unhelpful: r.unhelpful_count ?? r.unhelpful ?? 0,
      reply: r.reply ? { content: r.reply.content, createdAt: r.reply.updated_at || r.reply.created_at } : null,
    };
  }

  // ✅ (п.3) убрали слова "Полезно/Не полезно": только 👍/👎 + счётчик
  // ✅ (п.5) "👍 плохо видно" -> делаем 👍 btn-primary
  function renderReviewCard(r0) {
    const r = normalizeReview(r0);

    return `
      <div class="card" data-review-id="${escapeHtml(r.id)}">
        <div class="card-body stack">
          <div class="d-flex justify-content-between align-items-start">
            <div>
              <div class="fw-semibold">Отзыв</div>
              <div class="text-muted small">${escapeHtml(formatDate(r.createdAt))}</div>
            </div>
            <div class="fs-5" title="Рейтинг">${stars(r.rating)}</div>
          </div>

          <div>${escapeHtml(r.content)}</div>

          <div class="d-flex gap-2">
            <button class="btn btn-primary btn-sm js-vote" data-value="1" type="button" title="Полезно">
              👍 <span class="badge js-helpful">${r.helpful}</span>
            </button>
            <button class="btn btn-danger btn-sm js-vote" data-value="-1" type="button" title="Не полезно">
              👎 <span class="badge js-unhelpful">${r.unhelpful}</span>
            </button>
          </div>

          <div class="small js-vote-status"></div>

          ${r.reply
        ? `
            <hr />
            <div class="card" style="border-radius: var(--radius-md);">
              <div class="card-body">
                <div class="fw-semibold mb-1">Ответ заведения</div>
                <div>${escapeHtml(r.reply.content)}</div>
                <div class="text-muted small mt-1">${escapeHtml(formatDate(r.reply.createdAt))}</div>
              </div>
            </div>
          `
        : ""
      }
        </div>
      </div>
    `;
  }

  function extractCountsFromVoteResponse(data) {
    if (data && typeof data === "object") {
      const h = data.helpful_count ?? data.helpful ?? data.review?.helpful_count ?? data.review?.helpful;
      const u = data.unhelpful_count ?? data.unhelpful ?? data.review?.unhelpful_count ?? data.review?.unhelpful;
      const helpful = Number(h);
      const unhelpful = Number(u);
      if (!Number.isNaN(helpful) && !Number.isNaN(unhelpful)) return { helpful, unhelpful };
    }
    return null;
  }

  // ✅ (п.1) вместо перерисовки всей страницы:
  // после голосования — подтягиваем отзывы и обновляем ТОЛЬКО один отзыв в DOM
  async function refreshSingleReviewCounts(placeId, reviewId) {
    const qs = currentReviewsQueryString();
    const res = await fetch(`${API_BASE}/places/${placeId}/reviews?${qs}`);
    const data = await safeParseJSON(res);
    if (!res.ok) return null;

    const reviews = Array.isArray(data) ? data : data.reviews || [];
    const found = reviews.find((x) => x && x.id === reviewId);
    if (!found) return null;

    const r = normalizeReview(found);
    return { helpful: r.helpful, unhelpful: r.unhelpful };
  }

  function wireVotes() {
    reviewsContainer.querySelectorAll(".card[data-review-id]").forEach((card) => {
      const reviewId = card.getAttribute("data-review-id");

      const status = card.querySelector(".js-vote-status");
      const helpfulEl = card.querySelector(".js-helpful");
      const unhelpfulEl = card.querySelector(".js-unhelpful");

      const voteBtns = card.querySelectorAll(".js-vote");

      voteBtns.forEach((btn) => {
        btn.addEventListener("click", async () => {
          const token = getUserToken();
          if (!token) {
            const redirect = encodeURIComponent(window.location.href);
            window.location.href = `login.html?redirect=${redirect}`;
            return;
          }

          const value = Number(btn.getAttribute("data-value"));
          status.innerHTML = "";

          voteBtns.forEach((b) => (b.disabled = true));

          try {
            const res = await fetch(`${API_BASE}/reviews/${reviewId}/vote`, {
              method: "POST",
              headers: {
                Authorization: `Bearer ${token}`,
                "Content-Type": "application/json",
                Accept: "application/json",
              },
              body: JSON.stringify({ value }),
            });

            const data = await safeParseJSON(res);
            if (!res.ok) throw new Error(data.error || `Ошибка голоса: ${res.status}`);

            // 1) если бэк вернул счётчики — обновим сразу
            const counts = extractCountsFromVoteResponse(data);
            if (counts && helpfulEl && unhelpfulEl) {
              helpfulEl.textContent = String(counts.helpful);
              unhelpfulEl.textContent = String(counts.unhelpful);
              status.innerHTML = `<span class="text-success">Готово ✅</span>`;
              return;
            }

            // 2) иначе — обновим ТОЛЬКО этот отзыв через fetch списка (без перерисовки всей страницы)
            const fresh = await refreshSingleReviewCounts(currentPlaceId, reviewId);
            if (fresh && helpfulEl && unhelpfulEl) {
              helpfulEl.textContent = String(fresh.helpful);
              unhelpfulEl.textContent = String(fresh.unhelpful);
              status.innerHTML = `<span class="text-success">Готово ✅</span>`;
              return;
            }

            // 3) последний fallback — просто сообщение (без reload)
            status.innerHTML = `<span class="text-success">Голос учтён ✅</span>`;
          } catch (e) {
            status.innerHTML = `<span class="text-danger">${escapeHtml(e.message)}</span>`;
          } finally {
            voteBtns.forEach((b) => (b.disabled = false));
          }
        });
      });
    });
  }

  // ====== init ======
  updateAuthUI();

  if (currentPlaceId) {
    showPlace();
    loadReviews(currentPlaceId);
  } else {
    showTop();
  }
});
