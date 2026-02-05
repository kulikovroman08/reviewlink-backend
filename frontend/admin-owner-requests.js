document.addEventListener("DOMContentLoaded", () => {
	const { API_BASE, parseJwt } = window.AppCommon || {};
	const BASE = API_BASE || window.location.origin;

	const container = document.getElementById("requestsContainer");
	const statusLine = document.getElementById("statusLine");
	const refreshBtn = document.getElementById("refreshBtn");
	const adminBadge = document.getElementById("adminBadge");
	const qInput = document.getElementById("qInput");
	const sortSelect = document.getElementById("sortSelect");

	function getToken() {
		return localStorage.getItem("adminToken") || "";
	}

	function authHeaders(extra = {}) {
		const token = getToken();
		return token ? { ...extra, Authorization: token } : { ...extra };
	}

	async function safeParseJSON(res) {
		try { return await res.json(); } catch { return {}; }
	}

	function escapeHtml(str) {
		return String(str || "")
			.replaceAll("&", "&amp;")
			.replaceAll("<", "&lt;")
			.replaceAll(">", "&gt;")
			.replaceAll('"', "&quot;")
			.replaceAll("'", "&#039;");
	}

	// ===== guard: только platform admin =====
	(function guard() {
		const token = getToken();
		if (!token) {
			window.location.href = "login.html";
			return;
		}

		const payload = parseJwt ? parseJwt(token) : null;
		if (!payload || payload.role !== "admin") {
			window.location.href = "login.html";
			return;
		}

		if (adminBadge) {
			adminBadge.className = "badge ok";
			adminBadge.textContent = "admin";
		}
	})();

	addLogoutButton();

	// ===== загрузка заявок =====
	async function loadRequests() {
		if (!container) return;

		container.innerHTML = "";
		if (statusLine) statusLine.textContent = "Загрузка заявок...";

		try {
			const res = await fetch(`${BASE}/admin/owner_requests/pending`, {
				method: "GET",
				headers: authHeaders(),
			});

			const data = await safeParseJSON(res);
			if (!res.ok) throw new Error(data.error || `Ошибка: ${res.status}`);

			let items = data.items || [];

			// сортировка
			if (sortSelect?.value === "old") {
				items = items.slice().sort((a, b) =>
					new Date(a.created_at) - new Date(b.created_at)
				);
			} else {
				items = items.slice().sort((a, b) =>
					new Date(b.created_at) - new Date(a.created_at)
				);
			}

			// фильтр
			const q = (qInput?.value || "").toLowerCase();
			if (q) {
				items = items.filter(r =>
					`${r.name} ${r.source} ${r.source_id} ${r.comment || ""} ${r.user_id}`
						.toLowerCase()
						.includes(q)
				);
			}

			if (statusLine) {
				statusLine.textContent = `Найдено заявок: ${items.length}`;
			}

			if (items.length === 0) {
				container.innerHTML = `<div class="text-muted">Заявок нет 🎉</div>`;
				return;
			}

			container.innerHTML = items.map(renderCard).join("");
			wireActions();

		} catch (e) {
			if (statusLine) statusLine.innerHTML = `<span class="text-danger">${escapeHtml(e.message)}</span>`;
		}
	}

	function renderCard(r) {
		return `
		<div class="card" data-id="${escapeHtml(r.id)}">
			<div class="card-body">

				<div class="d-flex justify-content-between align-items-start">
					<div>
						<div class="fw-semibold">${escapeHtml(r.name)}</div>
						<div class="text-muted small">
							${escapeHtml(r.source)} / ${escapeHtml(r.source_id)}
						</div>
					</div>
					<span class="badge">${escapeHtml(r.status)}</span>
				</div>

				<div class="subtle" style="margin-top:6px;">
					user_id: <code>${escapeHtml(r.user_id)}</code>
				</div>

				<div class="subtle">
					создано: ${escapeHtml(r.created_at || "")}
				</div>

				<hr />

				<textarea
					class="form-control form-control-sm js-comment"
					placeholder="Комментарий (например: документы проверены)"
					rows="2"
				></textarea>

				<div class="row" style="margin-top:10px; gap:10px;">
					<button class="btn btn-success js-approve" type="button">Одобрить</button>
					<button class="btn btn-danger js-reject" type="button">Отклонить</button>
					<span class="muted js-status"></span>
				</div>

			</div>
		</div>
		`;
	}

	function wireActions() {
		document.querySelectorAll("#requestsContainer .card").forEach(card => {
			const id = card.getAttribute("data-id");
			const commentInput = card.querySelector(".js-comment");
			const statusEl = card.querySelector(".js-status");

			card.querySelector(".js-approve").onclick = async () => {
				await sendDecision(card, id, "approve", commentInput.value, statusEl);
			};

			card.querySelector(".js-reject").onclick = async () => {
				await sendDecision(card, id, "reject", commentInput.value, statusEl);
			};
		});
	}

	async function sendDecision(card, id, action, comment, statusEl) {
		statusEl.textContent = "Отправка...";

		try {
			const res = await fetch(`${BASE}/admin/owner_requests/${id}/${action}`, {
				method: "POST",
				headers: authHeaders({ "Content-Type": "application/json" }),
				body: JSON.stringify({ comment }),
			});

			const data = await safeParseJSON(res);
			if (!res.ok) throw new Error(data.error || "Ошибка");

			statusEl.innerHTML = `<span class="text-success">Готово</span>`;

			setTimeout(() => {
				card.remove();
			}, 400);

		} catch (e) {
			statusEl.innerHTML = `<span class="text-danger">${escapeHtml(e.message)}</span>`;
		}
	}

	// ===== UI =====
	refreshBtn?.addEventListener("click", loadRequests);
	qInput?.addEventListener("input", loadRequests);
	sortSelect?.addEventListener("change", loadRequests);

	function addLogoutButton() {
		const navRight = document.querySelector(".nav-right");
		if (!navRight) return;

		// уже есть
		if (document.querySelector("#adminLogoutBtn")) return;

		const logoutBtn = document.createElement("button");
		logoutBtn.id = "adminLogoutBtn";
		logoutBtn.className = "btn btn-sm";
		logoutBtn.type = "button";

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

		// вставляем сразу после кнопки "Тема"
		const themeBtn = navRight.querySelector('button[onclick*="toggleTheme"]');
		if (themeBtn && themeBtn.nextSibling) {
			navRight.insertBefore(logoutBtn, themeBtn.nextSibling);
		} else if (themeBtn) {
			navRight.appendChild(logoutBtn);
		} else {
			navRight.insertBefore(logoutBtn, navRight.firstChild);
		}
	}


	loadRequests();
});
