import {renderTbody} from '../helpers';

export default class Expected {
	constructor(api) {
		this.api = api;

		this.update();
	}

	async update() {
		try {
			const expected = await this.api.get('expected');

			const target = document.getElementById('target');
			renderTbody(target, expected, [
				'address',
				'port',
				'protocol',
				'state',
				'comment',
			]);
		} catch (error) {
			console.log(error);
		}

		setTimeout(() => this.update(), 5000);
	}
}

