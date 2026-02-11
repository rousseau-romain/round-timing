document.addEventListener('htmx:configRequest', function(event) {
	var el = document.getElementById('csrf-data');
	if (el) {
		event.detail.headers['X-CSRF-Token'] = el.dataset.token;
	}
});
