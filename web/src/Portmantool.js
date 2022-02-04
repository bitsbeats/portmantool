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
