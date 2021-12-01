export default class Portmantool {
	constructor() {
		this.apiUrl = `${window.origin}/v1`;
	}

	async get(endpoint) {
		const response = await fetch(`${this.apiUrl}/${endpoint}`);

		if (!response.ok) {
			const message = await response.text();
			throw message || `${response.status} ${response.statusText}`;
		}

		return response.json();
	}

	async update(endpoint, method, data) {
		const response = await fetch(`${this.apiUrl}/${endpoint}`, {
			method: method,
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(data),
		});

		if (!response.ok) {
			const message = await response.text();
			throw message || `${response.status} ${response.statusText}`;
		}
	}
}
