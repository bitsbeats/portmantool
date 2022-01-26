import {renderTbody} from '../helpers';

export default class Expected {
	constructor(api) {
		this.api = api;

		this.load();
	}

	async load() {
		const target = document.getElementById('target');

		try {
			const expected = await this.api.get('expected');

			renderTbody(target, expected, [
				'address',
				'port',
				'protocol',
				'state',
				'comment',
				null,
			]);
		} catch (error) {
			console.log(error);
		}

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

			try {
				await this.api.update([{
					address: address.value,
					port: Number.parseInt(port.value),
					protocol: protocol.value,
					state: state.value,
					comment: comment.value,
				}]);
			} catch (error) {
				console.log(error);
				return;
			}

			address.value = '';
			port.value = '';
			protocol.value = '';
			state.value = '';
			comment.value = '';
		});

		td.append(button);

		inputs.append(td);

		target.prepend(inputs);
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
				option.selected = value === placeholder;
				option.value = value;

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
