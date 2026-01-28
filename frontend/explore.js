document.addEventListener("DOMContentLoaded", () => {
	const { API_BASE } = window.AppCommon;

	// ====== DOM ======
	const topView = document.getElementById("topView");
	const placeView = document.getElementById("placeView");

	const topStatus = document.getElementById("topStatus");
	const topContainer = document.getElementById("topContainer");
	const refreshTopBtn = document.getElementById("refreshTopBtn");

	const backToTopBtn = document.getElementById("backToTopBtn");
	const refreshReviewsBtn = document.getElementById("refreshReviewsBtn");
	const leaveReviewBtn = document.getElementById("leaveReviewBtn");

	const placeTitle = document.getElementById("placeTitle");
	const reviewsStatus = document.getElementById("reviewsStatus");
	const reviewsContainer = document.getElementById("reviewsContainer");

	const sortSelect = document.getElementById("sortSelect");
	const ratingSelect = document.getElementById("ratingSelect");

	const authPill = document.getElementById("authPill");
	const voteHint = document.getElementById("voteHint");

	// NEW UI
	const modeSelect = document.getElementById("modeSelect");
	const modeHint = document.getElementById("modeHint");
	const catalogFilters = document.getElementById("catalogFilters");
	const searchInput = document.getElementById("searchInput");
	const amenityInput = document.getElementById("amenityInput");
	const limitInput = document.getElementById("limitInput");

	const mapEl = document.getElementById("map");

	// ====== state ======
	let activeTab = "places"; // places | users | bonuses | catalog
	let currentPlaceId = getPlaceIdFromUrl();
	let currentPlaceName = "";

	// public catalog + map
	let publicAbort = null;
	let map = null;
	let markersLayer = null;

	// кеш последней выдачи каталога (чтобы клик по маркеру знал куда вести)
	let lastCatalogItems = [];

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

	// ===== ensure helper (POST /places/ensure_from_public) =====
	async function ensurePlaceIdFromPublic({ source, sourceId, name }) {
		const token = getUserToken();
		if (!token) throw new Error("not authorized");

		const res = await fetch(`${API_BASE}/places/ensure_from_public`, {
			method: "POST",
			headers: {
				Authorization: `Bearer ${token}`,
				"Content-Type": "application/json",
				Accept: "application/json",
			},
			body: JSON.stringify({
				source,
				source_id: sourceId,
				name,
			}),
		});

		const data = await safeParseJSON(res);
		if (!res.ok) throw new Error(data.error || `ensure place: ${res.status}`);

		return data.place_id || data.placeId || "";
	}

	function clampInt(v, def, min, max) {
		const n = Number.parseInt(v, 10);
		if (Number.isNaN(n)) return def;
		return Math.max(min, Math.min(max, n));
	}

	function currentReviewsQueryString() {
		const qs = new URLSearchParams();
		if (sortSelect && sortSelect.value) qs.set("sort", sortSelect.value);
		if (ratingSelect && ratingSelect.value) qs.set("rating", ratingSelect.value);
		return qs.toString();
	}

	function setModeHint() {
		if (!modeHint) return;
		if (activeTab === "catalog") {
			modeHint.textContent =
				"Каталог: можно открыть отзывы или перейти к оставлению отзыва (если место уже создано в ReviewLink).";
			return;
		}
		if (activeTab === "places") {
			modeHint.textContent = "Лидеры по заведениям.";
			return;
		}
		if (activeTab === "users") {
			modeHint.textContent = "Лидеры по пользователям.";
			return;
		}
		if (activeTab === "bonuses") {
			modeHint.textContent = "Лидеры по бонусам.";
			return;
		}
		modeHint.textContent = "";
	}

	// ====== mode switch (одна плашка) ======
	modeSelect?.addEventListener("change", async () => {
		activeTab = modeSelect.value || "places";
		setModeHint();

		if (catalogFilters) {
			catalogFilters.classList.toggle("d-none", activeTab !== "catalog");
		}
		if (mapEl) {
			mapEl.classList.toggle("d-none", activeTab !== "catalog");
		}

		await loadTop();
	});

	refreshTopBtn?.addEventListener("click", loadTop);

	// ====== loaders ======
	async function loadTop() {
		topStatus.textContent = "Загрузка...";
		topContainer.innerHTML = "";

		try {
			if (activeTab === "catalog") {
				await loadPublicCatalog();
				return;
			}

			// ЛИДЕРЫ: убрали фильтры UI — фиксированные параметры
			const limit = 50;

			if (activeTab === "places") {
				const qs = new URLSearchParams({
					limit: String(limit),
					sort_by: "reviews",
					min_rating: "0",
					min_reviews: "0",
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
            <tr data-place-id="${escapeHtml(p.id)}" data-place-name="${escapeHtml(
					p.name
				)}" class="js-place-row" style="cursor:pointer">
              <td>${p.rank}</td>
              <td><div class="fw-semibold">${escapeHtml(p.name)}</div></td>
              <td style="text-align:right;">${p.reviews_count}</td>
              <td style="text-align:right;">${Number(p.avg_rating || 0).toFixed(2)}</td>
              <td style="text-align:right;">
                <button class="btn btn-sm js-open-place" type="button"
                  data-place-id="${escapeHtml(p.id)}"
                  data-place-name="${escapeHtml(p.name)}">
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
              <td><div class="fw-semibold">${escapeHtml(u.name)}</div></td>
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

	// ======================
	// PUBLIC PLACES CATALOG
	// ======================
	function normalizePublicPlace(p) {
		const rl = p.reviewlink || {};
		return {
			source: p.source || "osm",
			sourceId: p.sourceId || "",
			name: p.name || "(без названия)",
			amenity: p.amenity || "",
			address: p.address?.display || "",
			lat: p.location?.lat,
			lon: p.location?.lon,
			placeId: rl.placeId || "",
			rating: rl.rating,
			reviewsCount: rl.reviewsCount,
		};
	}

	function renderPublicPlaces(items) {
		return `
      <table>
        <thead>
          <tr>
            <th>#</th>
            <th>Заведение</th>
            <th style="text-align:right;">Отзывы</th>
            <th style="text-align:right;">Рейтинг</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          ${items
			.map((p0, i) => {
				const p = normalizePublicPlace(p0);
				const rating = p.rating == null ? "—" : Number(p.rating).toFixed(2);
				const reviews = p.reviewsCount == null ? "—" : String(p.reviewsCount);
				const action = p.placeId ? "Отзывы →" : "Оставить отзыв →";

				return `
              <tr class="js-public-row"
                  data-index="${i}"
                  data-place-id="${escapeHtml(p.placeId)}"
                  data-source="${escapeHtml(p.source)}"
                  data-source-id="${escapeHtml(p.sourceId)}"
                  data-place-name="${escapeHtml(p.name)}"
                  style="cursor:pointer">
                <td>${i + 1}</td>
                <td>
                  <div class="fw-semibold">${escapeHtml(p.name)}</div>
                  <div class="small text-muted">
                    ${p.amenity ? escapeHtml(p.amenity) : ""}
                    ${p.address ? " • " + escapeHtml(p.address) : ""}
                  </div>
                </td>
                <td style="text-align:right;">${reviews}</td>
                <td style="text-align:right;">${rating}</td>
                <td style="text-align:right;">
                  <button class="btn btn-sm btn-primary js-public-open" type="button"
                    data-place-id="${escapeHtml(p.placeId)}"
                    data-source="${escapeHtml(p.source)}"
                    data-source-id="${escapeHtml(p.sourceId)}"
                    data-place-name="${escapeHtml(p.name)}">
                    ${action}
                  </button>
                </td>
              </tr>
            `;
			})
			.join("")}
        </tbody>
      </table>
    `;
	}

	function wirePublicPlaces(items) {
		topContainer.querySelectorAll(".js-public-row").forEach((row) => {
			row.addEventListener("click", () => {
				const idx = Number(row.getAttribute("data-index"));
				const p0 = items[idx];
				focusOnPublicPlace(p0);
			});
		});

		topContainer.querySelectorAll(".js-public-open").forEach((btn) => {
			btn.addEventListener("click", async (e) => {
				e.stopPropagation();
				await handleOpenPublicPlaceFromEl(btn);
			});
		});
	}

	// ✅ стало async + await ниже
	async function handleOpenPublicPlaceFromEl(el) {
		const placeId = el.getAttribute("data-place-id") || "";
		const name = el.getAttribute("data-place-name") || "";
		const source = el.getAttribute("data-source") || "osm";
		const sourceId = el.getAttribute("data-source-id") || "";

		await handleOpenPublicPlaceFromData({ placeId, name, source, sourceId });
	}

	// ✅ ensure helper
	async function ensurePlaceIdFromPublic({ source, sourceId, name }) {
		const token = getUserToken();
		if (!token) throw new Error("not authorized");

		const res = await fetch(`${API_BASE}/places/ensure_from_public`, {
			method: "POST",
			headers: {
				Authorization: `Bearer ${token}`,
				"Content-Type": "application/json",
				Accept: "application/json",
			},
			body: JSON.stringify({
				source,
				source_id: sourceId,
				name,
			}),
		});

		const data = await safeParseJSON(res);
		if (!res.ok) throw new Error(data.error || `ensure place: ${res.status}`);

		// ожидаем { place_id: "..." }
		return data.place_id || data.placeId || "";
	}

	async function handleOpenPublicPlaceFromData({ placeId, name, source, sourceId }) {
		if (placeId) {
			openPlace(placeId, name);
			return;
		}

		if (!isLoggedIn()) {
			const redirect = encodeURIComponent(window.location.href);
			window.location.href = `login.html?mode=login&redirect=${redirect}`;
			return;
		}

		// нет placeId -> создаём/получаем placeId на бэке и уходим на форму
		try {
			if (topStatus) topStatus.textContent = "Создаём заведение...";
			const newPlaceId = await ensurePlaceIdFromPublic({ source, sourceId, name });

			if (!newPlaceId) throw new Error("ensure did not return place_id");

			window.location.href = `review-form.html?place_id=${encodeURIComponent(newPlaceId)}`;
		} catch (e) {
			alert(`Не удалось создать заведение в ReviewLink: ${e?.message || e}`);
		} finally {
			if (topStatus) topStatus.textContent = "";
		}
	}

	function initMapOnce() {
		if (map) return;
		if (!mapEl) return;

		map = L.map("map").setView([48.708, 44.514], 11);

		L.tileLayer("https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png", {
			maxZoom: 19,
			attribution: "&copy; OpenStreetMap",
		}).addTo(map);

		markersLayer = L.layerGroup().addTo(map);
	}

	function renderPublicMarkers(items) {
		if (!map || !markersLayer) return;

		markersLayer.clearLayers();

		const points = [];

		items.forEach((p0) => {
			const p = normalizePublicPlace(p0);
			if (!p.lat || !p.lon) return;

			points.push([p.lat, p.lon]);

			const marker = L.marker([p.lat, p.lon]).addTo(markersLayer);

			marker.bindPopup(`
        <div style="min-width:180px">
          <b>${escapeHtml(p.name)}</b><br/>
          <span class="small">${p.address ? escapeHtml(p.address) : ""}</span><br/>
          <button type="button" class="btn btn-sm btn-primary" style="margin-top:8px"
            onclick="window.__rlOpenFromMap && window.__rlOpenFromMap(${JSON.stringify({
				placeId: p.placeId,
				name: p.name,
				source: p.source,
				sourceId: p.sourceId,
			}).replaceAll('"', "&quot;")})">
            ${p.placeId ? "Открыть отзывы →" : "Оставить отзыв →"}
          </button>
        </div>
      `);

			marker.on("click", () => {
				// просто открыть попап — уже удобно
			});
		});

		if (points.length > 1) {
			map.fitBounds(points, { padding: [30, 30] });
		} else if (points.length === 1) {
			map.setView(points[0], 14);
		}
	}

	// клики из попапа (Leaflet popup кнопка)
	window.__rlOpenFromMap = (payload) => {
		if (!payload) return;
		Promise.resolve(handleOpenPublicPlaceFromData(payload)).catch((e) => {
			alert(`Не удалось открыть/создать заведение: ${e?.message || e}`);
		});
	};

	function focusOnPublicPlace(p0) {
		if (!map) return;
		const p = normalizePublicPlace(p0);
		if (!p.lat || !p.lon) return;

		map.setView([p.lat, p.lon], 16);
	}

	async function loadPublicCatalog() {
		topStatus.textContent = "Загрузка...";
		topContainer.innerHTML = "";

		const city = "Волгоград";
		const limit = clampInt(limitInput?.value, 50, 1, 200);
		const search = (searchInput?.value || "").trim();
		const amenity = (amenityInput?.value || "").trim();

		if (publicAbort) publicAbort.abort();
		publicAbort = new AbortController();

		try {
			const qs = new URLSearchParams({
				city: city,
				limit: String(limit),
			});
			if (search) qs.set("search", search);
			if (amenity) qs.set("amenity", amenity);

			const res = await fetch(`${API_BASE}/places/public?${qs.toString()}`, {
				signal: publicAbort.signal,
			});

			const data = await safeParseJSON(res);
			if (!res.ok) throw new Error(data.error || `Ошибка places public: ${res.status}`);

			const items = Array.isArray(data.items) ? data.items : [];
			lastCatalogItems = items;

			topStatus.textContent = `Заведений: ${items.length}`;
			topContainer.innerHTML = renderPublicPlaces(items);
			wirePublicPlaces(items);

			initMapOnce();
			renderPublicMarkers(items);
		} catch (e) {
			if (e && e.name === "AbortError") return;
			topStatus.innerHTML = `<span class="text-danger">${escapeHtml(e.message)}</span>`;
		}
	}

	// авто-перезагрузка каталога по фильтрам (debounce)
	let catalogTimer = null;
	function debounceCatalogReload() {
		if (activeTab !== "catalog") return;
		if (catalogTimer) clearTimeout(catalogTimer);
		catalogTimer = setTimeout(() => loadTop(), 350);
	}
	searchInput?.addEventListener("input", debounceCatalogReload);
	amenityInput?.addEventListener("input", debounceCatalogReload);
	limitInput?.addEventListener("change", debounceCatalogReload);

	// ====== place reviews ======
	backToTopBtn?.addEventListener("click", () => {
		currentPlaceId = "";
		currentPlaceName = "";
		clearPlaceIdFromUrl();
		showTop();
	});

	refreshReviewsBtn?.addEventListener("click", () => {
		if (currentPlaceId) loadReviews(currentPlaceId);
	});

	leaveReviewBtn?.addEventListener("click", () => {
		if (!currentPlaceId) return;

		if (!isLoggedIn()) {
			const target = `review-form.html?place_id=${encodeURIComponent(currentPlaceId)}`;
			window.location.href = `login.html?mode=login&redirect=${encodeURIComponent(target)}`;
			return;
		}

		window.location.href = `review-form.html?place_id=${encodeURIComponent(currentPlaceId)}`;
	});

	sortSelect?.addEventListener("change", () => {
		if (currentPlaceId) loadReviews(currentPlaceId);
	});

	ratingSelect?.addEventListener("change", () => {
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
		if (leaveReviewBtn) leaveReviewBtn.classList.add("d-none");

		if (mapEl) {
			mapEl.classList.toggle("d-none", activeTab !== "catalog");
		}

		loadTop();
	}

	function showPlace() {
		if (mapEl) mapEl.classList.add("d-none");
		topView.classList.add("d-none");
		placeView.classList.remove("d-none");
		reviewsContainer.innerHTML = "";
		reviewsStatus.textContent = "Загрузка...";
		placeTitle.textContent = currentPlaceName ? `Отзывы: ${currentPlaceName}` : "Отзывы заведения";
		updateAuthUI();
		if (leaveReviewBtn) leaveReviewBtn.classList.remove("d-none");
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
			reply: r.reply
				? { content: r.reply.content, createdAt: r.reply.updated_at || r.reply.created_at }
				: null,
		};
	}

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

          ${
			r.reply
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

						const counts = extractCountsFromVoteResponse(data);
						if (counts && helpfulEl && unhelpfulEl) {
							helpfulEl.textContent = String(counts.helpful);
							unhelpfulEl.textContent = String(counts.unhelpful);
							status.innerHTML = `<span class="text-success">Готово ✅</span>`;
							return;
						}

						const fresh = await refreshSingleReviewCounts(currentPlaceId, reviewId);
						if (fresh && helpfulEl && unhelpfulEl) {
							helpfulEl.textContent = String(fresh.helpful);
							unhelpfulEl.textContent = String(fresh.unhelpful);
							status.innerHTML = `<span class="text-success">Готово ✅</span>`;
							return;
						}

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

	// init mode UI
	activeTab = modeSelect?.value || "places";
	setModeHint();
	if (catalogFilters) catalogFilters.classList.toggle("d-none", activeTab !== "catalog");
	if (mapEl) mapEl.classList.toggle("d-none", activeTab !== "catalog");

	if (currentPlaceId) {
		showPlace();
		loadReviews(currentPlaceId);
	} else {
		showTop();
	}
});
