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

export function matchFilter(filter, object) {
	for (const clause of filter) {
		let matched = true;

		for (const key in clause) {
			if (object[key] === undefined) {
				continue;
			}

			matched = false;

			for (let value of clause[key]) {
				const negated = /^[-!^~]/.test(value);
				value = value.replace(/^[-!^~]+/, '');

				switch (key) {
					case 'address':
						// support CIDR notation
						//break;
					case 'port':
						// support ranges
						//break;
					default:
						matched = negated ? object[key] !== value : object[key] === value;
				}

				if (matched) {
					break;
				}
			}

			if (!matched) {
				break;
			}
		}

		if (matched) {
			return true;
		}
	}

	return false;
}

export function parseFilter(value) {
	const filter = [{}];

	const words = value.trim().split(/\s+/);
	for (const w of words) {
		let key;

		if (w === '') {
			continue;
		} else if (/^(?:[+\/:|]+|or)$/.test(w)) {
			filter.push({});
			continue;
		} else if (/^addr(?:ess)?:/.test(w)) {
			key = 'address';
		} else if (w.startsWith('port:')) {
			key = 'port';
		} else if (/^proto(?:col)?:/.test(w)) {
			key = 'protocol';
		} else if (w.startsWith('state:')) {
			key = 'state';
		} else if (/^exp(?:ect|ected)?:/.test(w)) {
			key = 'expected';
		} else if (/^actual|cur(?:r|rent)?:/.test(w)) {
			key = 'actual';
		} else {
			continue;
		}

		const [_, values] = w.split(':', 2);
		if (filter[filter.length-1][key] === undefined) {
			filter[filter.length-1][key] = [];
		}
		filter[filter.length-1][key] = [...filter[filter.length-1][key], ...values.split(',')];
	}

	return filter;
}
