function applyThemeMode(mode) {
	const root = document.documentElement;
	console.log("Applying theme mode:", mode);
	console.log("root", root);
	root.dataset.darkMode = mode;
	if (mode === "dark") {
		root.classList.add("dark");
	} else if (mode === "light") {
		root.classList.remove("dark");
	} else {
		root.classList.toggle(
			"dark",
			window.matchMedia("(prefers-color-scheme: dark)").matches,
		);
	}
}

document.addEventListener("DOMContentLoaded", () => {
	const mode = document.documentElement.dataset.darkMode || "auto";
	if (mode === "auto") {
		applyThemeMode("auto");
	}
	// Listen for system preference changes (auto mode)
	window
		.matchMedia("(prefers-color-scheme: dark)")
		.addEventListener("change", (e) => {
			if (
				(document.documentElement.dataset.darkMode || "auto") ===
				"auto"
			) {
				document.documentElement.classList.toggle("dark", e.matches);
			}
		});
});

// Apply theme mode after HTMX swap when dark mode config is updated
document.addEventListener("htmx:afterSettle", () => {
	const el = document.querySelector("[data-active-theme]");
	if (el) {
		applyThemeMode(el.dataset.activeTheme);
	}
});
