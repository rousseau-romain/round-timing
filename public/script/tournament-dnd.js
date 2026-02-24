document.addEventListener('dragend', function() {
	document.querySelectorAll('[data-team-id]').forEach(function(el) {
		el.style.outline = '';
		el.style.backgroundColor = '';
	});
});

function onPlayerDragStart(event, el) {
	event.dataTransfer.setData('playerId', el.dataset.playerId);
	window._draggedPlayerId = el.dataset.playerId;
}

function onTeamDragOver(event, el) {
	const ids = el.dataset.playerIds ? el.dataset.playerIds.split(',') : [];
	if (window._draggedPlayerId && ids.includes(window._draggedPlayerId)) {
		el.style.outline = '2px solid #ef4444';
		el.style.backgroundColor = '#fef2f2';
		// Do NOT call preventDefault() → browser shows "not-allowed" cursor
	} else {
		event.preventDefault();
		el.style.outline = '2px solid #6366f1';
		el.style.backgroundColor = '#eef2ff';
	}
}

function onTeamDragLeave(event, el) {
	if (!el.contains(event.relatedTarget)) {
		el.style.outline = '';
		el.style.backgroundColor = '';
	}
}

function handlePlayerDrop(event, teamId, el) {
	event.preventDefault();
	el.style.outline = '';
	el.style.backgroundColor = '';
	const playerId = event.dataTransfer.getData('playerId');
	if (!playerId) return;
	const ids = el.dataset.playerIds ? el.dataset.playerIds.split(',') : [];
	if (ids.includes(playerId)) return;
	const csrfToken = document.getElementById('csrf-data')?.dataset.token ?? '';
	const body = new URLSearchParams({ id_player: playerId });
	fetch('/tournament/team/' + teamId + '/player', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/x-www-form-urlencoded',
			'X-CSRF-Token': csrfToken,
		},
		body: body.toString()
	})
	.then(res => res.ok ? res.text() : null)
	.then(html => {
		if (!html) return;
		const target = document.getElementById('team-composition-' + teamId);
		if (!target) return;
		const tmp = document.createElement('div');
		tmp.innerHTML = html;
		const newEl = tmp.firstElementChild;
		target.replaceWith(newEl);
		htmx.process(newEl);
	});
}
