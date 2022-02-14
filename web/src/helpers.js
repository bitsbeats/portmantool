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

export function formatDateTime(date) {
	if (!(date instanceof Date)) {
		date = new Date(date);
	}

	const i = 'yyyy-mm-dd'.length;
	const j = 'hh:mm:ss'.length;

	const str = date.toISOString();
	return `${str.substring(0, i)} ${str.substring(i+1, i+j+1)}`;
}

export function renderTbody(elem, data, keys) {
	elem.replaceChildren(...data.map(row => {
		const tr = document.createElement('tr');
		tr.replaceChildren(...keys.map(key => {
			const td = document.createElement('td');
			if (key !== null) {
				td.innerText = row[key];
			}

			return td;
		}));

		return tr;
	}));
}
