const KEYS_EXPECTED = [
	'address',
	'port',
	'protocol',
	'state',
	'comment',
	null, // "add" button
];

export default class Expected {
	constructor(api) {
		this.api = api;
		this.rows = [];

		this.load();
	}

	async load() {
		const target = document.getElementById('target');

		const inputs = document.createElement('tr');

		inputs.append(this.createInputCell('address', 'text', null, '192.0.2.42'));
		inputs.append(this.createInputCell('port', 'number', [1, 65535], '2342'));
		inputs.append(this.createInputCell('protocol', 'select', ['tcp', 'udp'], 'tcp'));
		inputs.append(this.createInputCell('state', 'select', ['open', 'closed', 'filtered', 'unfiltered', 'open|filtered', 'closed|filtered'], 'closed'));
		inputs.append(this.createInputCell('comment', 'text', null, 'This comment intentionally left blank.'));

		const td = document.createElement('td');

		const button = document.createElement('button');
		['button', 'is-small', 'is-success'].forEach(cls => button.classList.add(cls));

		const icon = document.createElement('span');
		['icon', 'has-text-weight-bold'].forEach(cls => icon.classList.add(cls));
		icon.append(document.createTextNode('\u2713'));

		button.append(icon);
		button.addEventListener('click', async () => {
			const [address, port, protocol, state, comment] = ['address', 'port', 'protocol', 'state', 'comment'].map(id => document.getElementById(id));
			const data = {
				address: address.value,
				port: Number.parseInt(port.value),
				protocol: protocol.value,
				state: state.value,
				comment: comment.value,
			};

			try {
				await this.api.update([data]);
			} catch (error) {
				console.log(error);
				return;
			}

			let row = this.rows.find(r => {
				const [addr, p, proto] = r.children;
				return addr.innerText === address.value && p.innerText === port.value && proto.innerText === protocol.value;
			});

			if (row === null) {
				row = document.createElement('tr');

				inputs.after(row);
			}

			row.classList.add('is-selected');
			setTimeout(() => row.classList.remove('is-selected'), 500);

			row.replaceChildren(...KEYS_EXPECTED.map(key => {
				const td = document.createElement('td');
				if (key !== null) {
					td.innerText = data[key];
				}

				return td;
			}));

			address.value = '';
			port.value = '';
			protocol.value = document.querySelector('#protocol option[selected]').value;
			state.value = document.querySelector('#state option[selected]').value;
			comment.value = '';
		});

		td.append(button);

		inputs.append(td);

		try {
			const expected = await this.api.get('expected');

			this.rows = expected.map(data => {
				const tr = document.createElement('tr');
				['is-clickable'].forEach(cls => tr.classList.add(cls));

				tr.replaceChildren(...KEYS_EXPECTED.map(key => {
					const td = document.createElement('td');
					if (key !== null) {
						td.innerText = data[key];
					}

					return td;
				}));

				tr.addEventListener('click', () => {
					const [address, port, protocol, state, comment] = ['address', 'port', 'protocol', 'state', 'comment'].map(id => document.getElementById(id));
					const [addr, p, proto, s, c] = tr.children;

					address.value = addr.innerText;
					port.value = p.innerText;
					protocol.value = proto.innerText;
					state.value = s.innerText;
					comment.value = c.innerText;
				});

				return tr;
			});
		} catch (error) {
			console.log(error);
		}

		target.replaceChildren(inputs, ...this.rows);
	}

	createInputCell(id, type, range, placeholder) {
		const td = document.createElement('td');

		if (type === 'select') {
			const select = document.createElement('div');
			['select', 'is-small'].forEach(cls => select.classList.add(cls));

			const input = document.createElement('select');
			input.id = id;
			input.replaceChildren(...range.map(value => {
				const option = document.createElement('option');
				option.innerText = value;
				option.value = value;

				if (value === placeholder) {
					option.setAttribute('selected', '');
				}

				return option;
			}));

			select.append(input);

			td.append(select);
		} else {
			const input = document.createElement('input');
			['input', 'is-small'].forEach(cls => input.classList.add(cls));
			input.id = id;
			input.type = type;
			input.placeholder = placeholder;

			if (type === 'number') {
				const [min, max] = range;
				input.max = max;
				input.min = min;
			}

			td.append(input);
		}

		return td;
	}
}
