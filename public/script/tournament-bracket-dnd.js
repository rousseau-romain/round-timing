// Bracket drag-and-drop: drag winner teams onto ghost match TBD slots.
// Follows the same patterns as tournament-dnd.js.

// Track partially filled ghost matches: key = "round-position", value = { 1: teamId, 2: teamId }
var _bracketSlotState = {};

document.addEventListener('dragend', function () {
	document.querySelectorAll('.bracket-drop-zone').forEach(function (el) {
		el.style.outline = '';
		el.style.backgroundColor = '';
	});
});

function onBracketTeamDragStart(event, el) {
	event.dataTransfer.setData('bracketTeamId', el.dataset.teamId);
	event.dataTransfer.setData('bracketTeamName', el.dataset.teamName);
	window._draggedBracketTeamId = el.dataset.teamId;
	window._draggedBracketTeamName = el.dataset.teamName;
}

function onBracketSlotDragOver(event, el) {
	var teamId = window._draggedBracketTeamId;
	if (!teamId) return;

	var key = el.dataset.round + '-' + el.dataset.position;
	var state = _bracketSlotState[key] || {};
	var otherSlot = el.dataset.slot === '1' ? '2' : '1';

	// Prevent dropping same team in both slots
	if (state[otherSlot] && state[otherSlot] === teamId) {
		el.style.outline = '2px solid #ef4444';
		el.style.backgroundColor = '#fef2f2';
		return;
	}

	// Prevent overwriting a filled slot
	if (state[el.dataset.slot]) {
		el.style.outline = '2px solid #ef4444';
		el.style.backgroundColor = '#fef2f2';
		return;
	}

	event.preventDefault();
	el.style.outline = '2px solid #6366f1';
	el.style.backgroundColor = '#eef2ff';
}

function onBracketSlotDragLeave(event, el) {
	if (!el.contains(event.relatedTarget)) {
		el.style.outline = '';
		el.style.backgroundColor = '';
	}
}

function handleBracketSlotDrop(event, el) {
	event.preventDefault();
	el.style.outline = '';
	el.style.backgroundColor = '';

	var teamId = event.dataTransfer.getData('bracketTeamId');
	var teamName = event.dataTransfer.getData('bracketTeamName');
	if (!teamId || !teamName) return;

	var key = el.dataset.round + '-' + el.dataset.position;
	var state = _bracketSlotState[key] || {};
	var slot = el.dataset.slot;
	var otherSlot = slot === '1' ? '2' : '1';

	// Prevent same team in both slots
	if (state[otherSlot] && state[otherSlot] === teamId) return;
	// Prevent overwriting
	if (state[slot]) return;

	// Fill the slot
	state[slot] = teamId;
	_bracketSlotState[key] = state;

	// Update DOM: replace TBD with team name
	el.textContent = teamName;
	el.classList.remove('text-gray-400');
	el.classList.add('text-indigo-600', 'dark:text-indigo-400', 'font-medium');
	el.removeAttribute('ondragover');
	el.removeAttribute('ondragleave');
	el.removeAttribute('ondrop');

	// If both slots are filled, create the match
	if (state['1'] && state['2']) {
		createMatchFromDrop(el.dataset.tournamentId, state['1'], state['2']);
		delete _bracketSlotState[key];
	}
}

function createMatchFromDrop(tournamentId, idTeam1, idTeam2) {
	var csrfToken = document.getElementById('csrf-data')?.dataset.token ?? '';
	var body = new URLSearchParams({
		id_team1: idTeam1,
		id_team2: idTeam2,
		bo_format: '1',
	});

	fetch('/tournament/' + tournamentId + '/match', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/x-www-form-urlencoded',
			'X-CSRF-Token': csrfToken,
		},
		body: body.toString(),
	})
		.then(function (res) {
			return res.ok ? res.text() : null;
		})
		.then(function (html) {
			if (!html) return;
			var target = document.getElementById('matches-section');
			if (!target) return;
			var tmp = document.createElement('div');
			tmp.innerHTML = html;
			var newEl = tmp.firstElementChild;
			target.replaceWith(newEl);
			htmx.process(newEl);
			document.body.dispatchEvent(new Event('tournamentStatusChanged'));
			// Reset state for fresh bracket
			_bracketSlotState = {};
		});
}
