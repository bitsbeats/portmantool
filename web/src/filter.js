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

const INET_BITS = {
	4: 32n,
	6: 128n,
};

function parseInet(address) {
	const octets = address.split('.').map(n => Number.parseInt(n));
	if (octets.length === 4 && octets.every(n => n >= 0 && n < 256)) {
		return { version: 4, bits: octets.map(n => BigInt(n)).reduce((r, s) => r << 8n | s) };
	}

	const hextets = address.split(':').map(n => n === '' ? n : Number.parseInt(n, 16));
	const empty = hextets.indexOf('');
	const lastEmpty = hextets.lastIndexOf('');
	if (hextets.length >= 3 && empty === lastEmpty && octets.every(n => n === '' || n >= 0 && n < 65536)) {
		if (empty !== -1) {
			hextets.splice(empty, 1, ...Array(8-(hextets.length-1)).fill(0));
		}

		if (hextets.length !== 8) {
			return null;
		}

		return { version: 6, bits: hextets.map(n => BigInt(n)).reduce((r, s) => r << 16n | s) };
	}

	return null;
}

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
						if (value.startsWith('cidr:')) {
							const parts = value.split(':');
							if (parts.length !== 4) {
								break;
							}

							const [_, v, a, b] = parts;

							try {
								const address = BigInt(a);
								const bitmask = BigInt(b);

								const inet = parseInet(object[key]);
								if (inet !== null && `${inet.version}` === v) {
									const x = address & bitmask;
									const y = inet.bits & bitmask;
									matched = negated ? x !== y : x === y;
								}
							} catch (error) {
							}

							break;
						}
					case 'port':
						// support ranges
						const range = value.split('-', 2);
						if (range.length === 2) {
							let [min, max] = range;
							min = Math.max(0, Number.parseInt(min));
							max = Math.min(65535, max === '' ? 65535 : Number.parseInt(max));

							const n = Number.parseInt(object[key]);
							const inRange = n >= min && n <= max;
							matched = negated ? !inRange : inRange;

							break;
						}
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

	for (const clause of filter) {
		if (clause['address'] === undefined) {
			continue;
		}

		clause['address'] = clause['address'].map(value => {
			const slash = value.indexOf('/');
			if (slash === -1) {
				return value;
			}

			const negated = /^[-!^~]/.test(value);
			value = value.replace(/^[-!^~]+/, '');

			const [address, subnet] = value.split('/', 2);

			let bits = Number.parseInt(subnet);
			if (!Number.isInteger(bits) || bits < 0) {
				return null;
			}
			bits = BigInt(bits);

			const inet = parseInet(address);
			if (inet === null || bits > INET_BITS[inet.version]) {
				return null;
			}

			const bitmask = (1n << bits) - 1n << INET_BITS[inet.version] - bits;

			return `${negated ? '!' : ''}cidr:${inet.version}:0x${inet.bits.toString(16)}:0x${bitmask.toString(16)}`;
		}).filter(value => value !== null);
	}

	return filter;
}
