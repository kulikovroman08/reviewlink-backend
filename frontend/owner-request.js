document.addEventListener("DOMContentLoaded", () => {
	const { API_BASE } = window.AppCommon || {};
	const BASE = API_BASE || window.location.origin;

	const authBadge = document.getElementById("authBadge");
	const sourceSelect = document.getElementById("sourceSelect");
	const sourceIdInput = document.getElementById("sourceIdInput");
	const nameInput = document.getElementById("nameInputPlace");
	const sendBtn = document.getElementById("sendBtn");
	const statusEl = document.getElementById("status");

	function getToken() {
		// для заявки достаточно userToken; на всякий — fallback на adminToken
		return localStorage.getItem("userToken") || localStorage.getItem("adminToken") || "";
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

	function setStatus(html) {
		if (statusEl) statusEl.innerHTML = html || "";
	}

	// ===== guard: если не залогинен — на login с redirect =====
	(function guard() {
		const token = getToken();
		if (!token) {
			const p = new URLSearchParams();
			p.set("redirect", "owner-request.html");
			window.location.replace(`login.html?${p.toString()}`);
			return;
		}

		if (authBadge) {
			authBadge.className = "badge ok";
			authBadge.textContent = "OK";
		}
	})();

	// ===== prefill из query params =====
	(function prefillFromQuery() {
		const p = new URLSearchParams(window.location.search);
		const source = p.get("source");
		const sourceId = p.get("source_id") || p.get("sourceId");
		const name = p.get("name");

		if (source && sourceSelect) sourceSelect.value = source;
		if (sourceId && sourceIdInput) sourceIdInput.value = sourceId;
		if (name && nameInput) nameInput.value = name;
	})();

	if (sendBtn) {
		sendBtn.addEventListener("click", async () => {
			const source = (sourceSelect?.value || "osm").trim();
			const sourceId = (sourceIdInput?.value || "").trim();
			const name = (nameInput?.value || "").trim();

			if (!sourceId || !name) {
				setStatus(`<span class="text-danger">Заполни sourceId и название</span>`);
				return;
			}

			sendBtn.disabled = true;
			setStatus(`<span class="text-info">Отправка...</span>`);

			try {
				const res = await fetch(`${BASE}/owner_requests`, {
					method: "POST",
					headers: authHeaders({ "Content-Type": "application/json" }),
					body: JSON.stringify({
						source,
						sourceId, // swagger допускает sourceId
						name,
					}),
				});

				const data = await safeParseJSON(res);

				if (res.status === 401 || res.status === 403) {
					const p = new URLSearchParams();
					p.set("redirect", "owner-request.html");
					window.location.replace(`login.html?${p.toString()}`);
					return;
				}

				if (!res.ok) {
					throw new Error(data.error || `Ошибка отправки: ${res.status}`);
				}

				const idPart = data.id ? ` ID: ${escapeHtml(data.id)}` : "";

				setStatus(`
					<div class="text-success">
						Заявка отправлена ✅${idPart}<br/>
						После одобрения войдите ещё раз — откроется админ-панель.
						<div class="subtle" style="margin-top:6px;">Сейчас вернём на главную…</div>
					</div>
				`);

				// уводим на главную, а не в ЛК
				setTimeout(() => {
					window.location.href = "explore.html";
				}, 1200);

			} catch (e) {
				setStatus(`<span class="text-danger">${escapeHtml(e.message || String(e))}</span>`);
			} finally {
				sendBtn.disabled = false;
			}
		});
	}
});
