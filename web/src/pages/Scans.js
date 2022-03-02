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

import {formatDateTime, renderTbody} from '../helpers';
import {showError} from '../notification';

export default class Scans {
	constructor(api) {
		this.api = api;

		this.keep = null;
		this.scan1 = null;
		this.scan2 = null;

		this.columnNames = [];

		this.load();
	}

	async load() {
		const createTh = text => {
			const th = document.createElement('th');
			th.append(document.createTextNode(text));
			return th;
		};

		const [thAddress, thPort, thProtocol, thState] = document.getElementById('columnNames').children;
		const [thStateA, thStateB, thPreviousState, thCurrentState, thLastScan] = [createTh('State A'), createTh('State B'), createTh('Previous state'), createTh('Current state'), createTh('Last scan')];

		this.columnNames.push([thAddress, thPort, thProtocol, thState]);
		this.columnNames.push([thAddress, thPort, thProtocol, thStateA, thStateB]);
		this.columnNames.push([thAddress, thPort, thProtocol, thPreviousState, thCurrentState, thLastScan]);

		await this.update();

		const keep = document.getElementById('keep');
		keep.addEventListener('change', event => {
			this.updateKeep(event.target);
		});
		this.updateKeep(keep);

		document.getElementById('prune').addEventListener('click', async () => {
			try {
				await this.api.prune(this.keep);

				await this.update();
			} catch (error) {
				console.log(error);
				showError(error);
			}

			window.location.hash = '';
		});
	}

	async update() {
		try {
			const scans = await this.api.get('scans');

			this.scan1 = null;
			this.scan2 = null;
			await this.show(true, true);

			document.getElementById('scans').replaceChildren(...scans.map(scan => this.createEntry(scan.id)));
		} catch (error) {
			console.log(error);
			showError(error);
		}
	}

	async show(scan1Changed, scan2Changed) {
		if (!scan1Changed && (this.scan1 === null || !scan2Changed)) {
			return;
		}

		if (this.scan1 !== null) {
			document.querySelectorAll('input[name=scan2]:disabled').forEach(e => {
				e.disabled = false;
			});
		}

		const result = document.getElementById('result');
		result.replaceChildren();

		const keys = [
			'address',
			'port',
			'protocol',
			'state',
		];
		let endpoint = this.scan1 !== null ? `scan/${this.scan1}` : '';

		const columnNames = document.getElementById('columnNames');
		if (this.scan1 !== null && this.scan2 !== null && this.scan1 === this.scan2) {
			keys.splice(keys.length-1, 1, 'state_a', 'state_b', 'scan_b');
			columnNames.replaceChildren(...this.columnNames[2]);

			endpoint = `diff/${this.scan1}`;
		} else if (this.scan1 !== null && this.scan2 !== null) {
			keys.splice(keys.length-1, 1, 'state_a', 'state_b');
			columnNames.replaceChildren(...this.columnNames[1]);

			endpoint = `diff/${this.scan1}/${this.scan2}`;
		} else {
			columnNames.replaceChildren(...this.columnNames[0]);
		}

		if (endpoint === '') {
			return;
		}

		try {
			const response = await this.api.get(endpoint);

			renderTbody(result, response, keys);
		} catch (error) {
			console.log(error);
			showError(error);
			return;
		}
	}

	createEntry(id) {
		const li = document.createElement('li');
		li.classList.add('is-flex');

		const scan1 = document.createElement('label');
		['radio', 'is-flex-grow-1'].forEach(cls => scan1.classList.add(cls));

		const scan1Input = document.createElement('input');
		scan1Input.name = 'scan1';
		scan1Input.type = 'radio';
		scan1Input.value = id;
		scan1Input.addEventListener('change', event => {
			this.scan1 = event.target.value;
			this.show(true, false);
		});

		const scan1Link = document.createElement('a');
		scan1Link.innerText = formatDateTime(id*1000);

		scan1.append(scan1Input);
		scan1.append(scan1Link);

		const scan2 = document.createElement('label');
		scan2.classList.add('radio');

		const scan2Input = document.createElement('input');
		scan2Input.disabled = true;
		scan2Input.name = 'scan2';
		scan2Input.type = 'radio';
		scan2Input.value = id;
		scan2Input.addEventListener('change', event => {
			this.scan2 = event.target.value;
			this.show(false, true);
		});
		scan2Input.addEventListener('click', event => {
			if (this.scan2 === event.target.value) {
				event.target.checked = false;
				this.scan2 = null;
				this.show(false, true);
			}
		});

		const scan2Link = document.createElement('a');
		scan2Link.innerText = '\u21c6';

		scan2.append(scan2Input);
		scan2.append(scan2Link);

		li.append(scan1);
		li.append(scan2);

		return li;
	}

	updateKeep(input) {
		this.keep = input.value === '' ? null : input.valueAsNumber / 1000;

		document.getElementById('keep-info').innerText = this.keep === null ? '' : ` older than ${formatDateTime(this.keep*1000)}`;
	}
}
