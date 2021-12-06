export function renderTbody(elem, data, keys) {
	elem.replaceChildren(...data.map(row => {
		const tr = document.createElement('tr');
		tr.replaceChildren(...keys.map(key => {
			const td = document.createElement('td');
			td.innerText = row[key];

			return td;
		}));

		return tr;
	}));
}
