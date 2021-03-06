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
import {formatDateTime} from '../helpers';
import {showError} from '../notification';

const KEYS_DIFF = [
	'address',
	'port',
	'protocol',
	'expected_state',
	'actual_state',
	'scan_id',
	'comment',
];

export default class Dashboard {
	constructor(api) {
		this.api = api;

		this.ups = 1;
		this.interval = 5 * this.ups;
		this.nextUpdate = 0;

		const filter = document.getElementById('filter');

		this.filter = parseFilter(filter.value);
		this.filterNeedsUpdate = true;

		filter.addEventListener('input', event => {
			this.filter = parseFilter(event.target.value);
			this.filterNeedsUpdate = true;
		});

		document.getElementById('nextUpdate').innerText = `${this.interval}s`;

		this.update();
	}

	async update() {
		if (this.nextUpdate === 0) {
			try {
				const info = await this.api.get('info');
				document.getElementById('lastSuccessfulImport').innerText = formatDateTime(info['last_successful_import']*1000);

				const diff = await this.api.get('diff');

				const focusedInput = document.activeElement;

				const target = document.getElementById('target');
				target.replaceChildren(...diff.map(row => {
					const id = `${row['address']}-${row['port']}-${row['protocol']}`;

					let tr = document.getElementById(id);
					if (!tr) {
						const comment = document.createElement('input');
						['input', 'is-small'].forEach(cls => comment.classList.add(cls));
						comment.placeholder = 'This comment intentionally left blank.';
						comment.type = 'text';

						const button = document.createElement('button');
						['button', 'is-small', 'is-success'].forEach(cls => button.classList.add(cls));

						const icon = document.createElement('span');
						['icon', 'has-text-weight-bold'].forEach(cls => icon.classList.add(cls));
						icon.append(document.createTextNode('\u2713'));

						button.append(icon);
						button.addEventListener('click', async () => {
							try {
								await this.api.update([{
									address: row['address'],
									port: row['port'],
									protocol: row['protocol'],
									state: row['actual_state'],
									comment: comment.value,
								}]);
								tr.remove();
							} catch (error) {
								console.log(error);
								showError(error);
							}
						});

						tr = document.createElement('tr');
						tr.id = id;
						tr.replaceChildren(...KEYS_DIFF
							.concat(['button'])
							.map(key => {
								const td = document.createElement('td');
								switch (key) {
									case 'address':
									case 'port':
									case 'protocol':
										td.innerText = row[key];
										break;
									case 'comment':
										td.replaceChildren(comment);
										break;
									case 'button':
										td.replaceChildren(button);
								}

								return td;
							})
						);
					}

					const setValue = (key, value) => {
						const index = KEYS_DIFF.findIndex(k => k === key);
						if (index === -1) {
							return;
						}

						tr.children[index].replaceChildren(value);
					};
					setValue('expected_state', row['expected_state']);
					setValue('actual_state', row['actual_state']);
					setValue('scan_id', row['actual_state'] !== '' ? formatDateTime(row['scan_id']) : '');

					return tr;
				}));

				focusedInput.focus();
			} catch (error) {
				console.log(error);
				showError(error);
			}

			this.nextUpdate = this.interval;
			this.filterNeedsUpdate = true;
		}

		document.getElementById('nextUpdate').innerText = `${this.nextUpdate}s`;

		--this.nextUpdate;

		if (this.filterNeedsUpdate) {
			for (const row of document.getElementById('target').children) {
				const [address, port, protocol, expected, current] = row.children;
				if (!matchFilter(this.filter, {
					address: address.innerText,
					port: port.innerText,
					protocol: protocol.innerText,
					state: current.innerText,
					expected: expected.innerText,
					actual: current.innerText,
				})) {
					row.classList.add('is-hidden');
				} else {
					row.classList.remove('is-hidden');
				}
			}

			this.filterNeedsUpdate = false;
		}

		setTimeout(() => this.update(), 1000/this.ups);
	}
}
