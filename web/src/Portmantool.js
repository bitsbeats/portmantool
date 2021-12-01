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

	async update(data) {
		const response = await fetch(`${this.apiUrl}/expected`, {
			method: 'PATCH',
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

	async prune(keep) {
		const response = await fetch(`${this.apiUrl}/scans${keep ? `/${keep}` : ''}`, {
			method: 'DELETE',
		});

		if (!response.ok) {
			const message = await response.text();
			throw message || `${response.status} ${response.statusText}`;
		}
	}
}
