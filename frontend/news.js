(function () {
	const FETCH_LIMIT = 20;  // сколько берем с бэка
	const SHOW_COUNT = 4;    // сколько показываем вверху (3..5)
	const ROTATE_MS = 9000;  // как часто менять (мс)

	let all = [];
	let start = 0;
	let timer = null;

	function el(id) { return document.getElementById(id); }

	function formatDt(iso) {
		try {
			const d = new Date(iso);
			const dd = String(d.getDate()).padStart(2, "0");
			const mm = String(d.getMonth() + 1).padStart(2, "0");
			const hh = String(d.getHours()).padStart(2, "0");
			const mi = String(d.getMinutes()).padStart(2, "0");
			return `${dd}.${mm} ${hh}:${mi}`;
		} catch (_) {
			return "";
		}
	}

	function sliceWindow(items, startIdx, count) {
		if (!items || items.length === 0) return [];
		const res = [];
		for (let i = 0; i < count; i++) {
			res.push(items[(startIdx + i) % items.length]);
		}
		return res;
	}

	function renderWindow() {
		const row = el("newsRow");     // <-- важно: newsRow (верхний ряд)
		const empty = el("newsEmpty");
		if (!row || !empty) return;

		row.innerHTML = "";

		if (!all || all.length === 0) {
			empty.style.display = "block";
			return;
		}
		empty.style.display = "none";

		const win = sliceWindow(all, start, Math.min(SHOW_COUNT, all.length));

		for (const it of win) {
			const card = document.createElement("div");
			card.className = "newsCard";

			const a = document.createElement("a");
			a.className = "newsTitle";
			a.href = it.link;
			a.target = "_blank";
			a.rel = "noopener noreferrer";
			a.textContent = it.title || "Без названия";

			const meta = document.createElement("div");
			meta.className = "newsMeta";

			const left = document.createElement("span");
			left.textContent = (it.source || "AFISHA").toUpperCase();

			const right = document.createElement("span");
			right.textContent = it.published_at ? formatDt(it.published_at) : "";

			meta.appendChild(left);
			meta.appendChild(right);

			card.appendChild(a);
			card.appendChild(meta);

			row.appendChild(card);
		}
	}

	function startRotation() {
		stopRotation();
		if (!all || all.length <= SHOW_COUNT) return;

		timer = setInterval(() => {
			start = (start + 1) % all.length;
			renderWindow();
		}, ROTATE_MS);
	}

	function stopRotation() {
		if (timer) {
			clearInterval(timer);
			timer = null;
		}
	}

	async function load() {
		const btn = el("newsRefreshBtn");
		if (btn) btn.disabled = true;

		try {
			const res = await fetch(`${API_BASE}/news?limit=${FETCH_LIMIT}`, {
				headers: { Accept: "application/json" },
			});
			const data = await res.json().catch(() => ({}));
			if (!res.ok) throw new Error(data?.error || `news: ${res.status}`);

			all = Array.isArray(data.items) ? data.items : [];
			start = 0;
			renderWindow();
			startRotation();
		} catch (e) {
			console.error("news load error:", e);
			all = [];
			start = 0;
			stopRotation();
			renderWindow();
		} finally {
			if (btn) btn.disabled = false;
		}
	}

	document.addEventListener("DOMContentLoaded", () => {
		el("newsRefreshBtn")?.addEventListener("click", load);

		// чтобы не “крутилось” в фоне, когда вкладка не активна
		document.addEventListener("visibilitychange", () => {
			if (document.hidden) stopRotation();
			else startRotation();
		});

		load();
	});
})();
