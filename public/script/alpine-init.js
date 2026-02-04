document.addEventListener("alpine:init", () => {
	const body = document.body;
	const isAuthenticated = body.dataset.authenticated === "true";
	const containerExpandedId = body.dataset.configContainerExpandedId;
	const darkModeId = body.dataset.configDarkModeId;

	function applyDarkMode(enabled) {
		document.documentElement.classList.toggle("dark", enabled);
	}

	Alpine.store("config", {
		init() {
			this.isContainerExpanded = body.dataset.configContainerExpanded === "true";

			if (isAuthenticated) {
				this.isDarkMode = body.dataset.configDarkMode === "true";
			} else {
				this.isDarkMode = window.matchMedia("(prefers-color-scheme: dark)").matches;
			}
			applyDarkMode(this.isDarkMode);
		},

		isContainerExpanded: false,
		toggleContainerExpanded() {
			this.isContainerExpanded = !this.isContainerExpanded;
			if (isAuthenticated) {
				fetch(`/profile/configuration/${containerExpandedId}/toggle-configuration`, { method: "PATCH" });
			}
		},

		isDarkMode: false,
		toggleDarkMode() {
			this.isDarkMode = !this.isDarkMode;
			applyDarkMode(this.isDarkMode);
			if (isAuthenticated) {
				fetch(`/profile/configuration/${darkModeId}/toggle-configuration`, { method: "PATCH" });
			}
		},
	});
});
