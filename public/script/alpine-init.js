document.addEventListener("alpine:init", () => {
	const body = document.body;
	const isAuthenticated = body.dataset.authenticated === "true";
	const containerExpandedId = body.dataset.configContainerExpandedId;
	const darkModeId = body.dataset.configDarkModeId;

	Alpine.store("config", {
		init() {
			this.isContainerExpanded = body.dataset.configContainerExpanded === "true";
			this.isDarkMode = body.dataset.configDarkMode === "true";
			if (!isAuthenticated) {
				this.isDarkMode = window.matchMedia("(prefers-color-scheme: dark)").matches;
			}
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
			if (isAuthenticated) {
				fetch(`/profile/configuration/${darkModeId}/toggle-configuration`, { method: "PATCH" });
			}
		},
	});
});
