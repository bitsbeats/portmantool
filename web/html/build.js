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

const fs = require('fs');
const path = require('path');

try {
	const pre = fs.readFileSync(path.resolve(__dirname, 'pre.html'));
	const post = fs.readFileSync(path.resolve(__dirname, 'post.html'));

	const pages = fs.readdirSync(path.resolve(__dirname, 'pages'));
	for (const page of pages) {
		try {
			if (/\.html$/.test(page)) {
				const contents = fs.readFileSync(path.resolve(__dirname, 'pages', page));
				fs.writeFileSync(path.resolve(__dirname, '../dist', page), `${pre}${contents}${post}`);
			}
		} catch (error) {
			console.error(error);
		}
	}
} catch (error) {
	console.error(error);
}
