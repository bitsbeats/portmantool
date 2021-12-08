import {renderTbody} from '../helpers';

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
		this.nextUpdate = 0;

		this.update();
	}

	async update() {
		if (this.nextUpdate === 0) {
			try {
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
					setValue('scan_id', new Date(row['scan_id']).toLocaleString());

					return tr;
				}));

				focusedInput.focus();
			} catch (error) {
				console.log(error);
			}

			this.nextUpdate = 5;
		}

		--this.nextUpdate;

		const progress = document.getElementById('nextUpdate');
		progress.value = 4 - this.nextUpdate;

		setTimeout(() => this.update(), 1000);
	}
}
