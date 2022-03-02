// Copyright 2020-2022 Thomann Bits & Beats GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import {matchFilter, parseFilter} from '../filter';
import {showError} from '../notification';

const KEYS_EXPECTED = [
	'address',
	'port',
	'protocol',
	'state',
	'comment',
	'button',
];

export default class Expected {
	constructor(api) {
		this.api = api;
		this.rows = [];
		this.filter = parseFilter(document.getElementById('filter').value);
		this.filterNeedsUpdate = true;

		this.load();
	}

	async load() {
		document.getElementById('filter').addEventListener('input', event => {
			this.filter = parseFilter(event.target.value);
			this.filterNeedsUpdate = true;
		});

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
				showError(error);
				return;
			}

			let row = this.rows.find(r => {
				const [addr, p, proto] = r.children;
				return addr.innerText === address.value && p.innerText === port.value && proto.innerText === protocol.value;
			});

			if (row === undefined) {
				row = document.createElement('tr');

				this.rows.splice(0, 0, row);

				inputs.after(row);
			}

			row.classList.add('is-selected');
			setTimeout(() => row.classList.remove('is-selected'), 500);

			const button = document.createElement('button');
			['button', 'is-small', 'is-light'].forEach(cls => button.classList.add(cls));

			const icon = document.createElement('span');
			['icon', 'has-text-weight-bold'].forEach(cls => icon.classList.add(cls));
			icon.append(document.createTextNode('\u25b2'));

			button.append(icon);
			button.addEventListener('click', () => {
				const [address, port, protocol, state, comment] = ['address', 'port', 'protocol', 'state', 'comment'].map(id => document.getElementById(id));
				const [addr, p, proto, s, c] = row.children;

				address.value = addr.innerText;
				port.value = p.innerText;
				protocol.value = proto.innerText;
				state.value = s.innerText;
				comment.value = c.innerText;
			});

			row.replaceChildren(...KEYS_EXPECTED.map(key => {
				const td = document.createElement('td');
				if (key === 'button') {
					td.append(button);
				} else {
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

				const button = document.createElement('button');
				['button', 'is-small', 'is-light'].forEach(cls => button.classList.add(cls));

				const icon = document.createElement('span');
				['icon', 'has-text-weight-bold'].forEach(cls => icon.classList.add(cls));
				icon.append(document.createTextNode('\u25b2'));

				button.append(icon);
				button.addEventListener('click', () => {
					const [address, port, protocol, state, comment] = ['address', 'port', 'protocol', 'state', 'comment'].map(id => document.getElementById(id));
					const [addr, p, proto, s, c] = tr.children;

					address.value = addr.innerText;
					port.value = p.innerText;
					protocol.value = proto.innerText;
					state.value = s.innerText;
					comment.value = c.innerText;
				});

				tr.replaceChildren(...KEYS_EXPECTED.map(key => {
					const td = document.createElement('td');
					if (key === 'button') {
						td.append(button);
					} else {
						td.innerText = data[key];
					}

					return td;
				}));

				return tr;
			});
		} catch (error) {
			console.log(error);
			showError(error);
		}

		target.replaceChildren(inputs, ...this.rows);

		this.update();
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

	update() {
		if (this.filterNeedsUpdate) {
			for (const row of this.rows) {
				const [address, port, protocol, expected] = row.children;
				if (!matchFilter(this.filter, {
					address: address.innerText,
					port: port.innerText,
					protocol: protocol.innerText,
					state: expected.innerText,
					expected: expected.innerText,
				})) {
					row.classList.add('is-hidden');
				} else {
					row.classList.remove('is-hidden');
				}
			}

			this.filterNeedsUpdate = false;
		}

		setTimeout(() => this.update(), 1000);
	}
}
