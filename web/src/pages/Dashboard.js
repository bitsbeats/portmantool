import {renderTbody} from '../helpers';

export default class Dashboard {
	constructor(api) {
		this.api = api;

		this.update();
	}

	async update() {
		try {
			const diff = await this.api.get('diff');
			for (const entry of diff) {
				entry['scan_id'] = new Date(entry['scan_id']).toLocaleString();
			}

			const target = document.getElementById('target');
			renderTbody(target, diff, [
				'address',
				'port',
				'protocol',
				'expected_state',
				'actual_state',
				'scan_id',
				'comment',
			]);
		} catch (error) {
			console.log(error);
		}

		setTimeout(() => this.update(), 5000);
	}
}
