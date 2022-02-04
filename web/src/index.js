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

import {Dashboard, Expected, Scans} from './pages';
import Portmantool from './Portmantool';
import './style.scss';

document.addEventListener('DOMContentLoaded', () => {
	const api = new Portmantool();

	switch (window.location.pathname) {
		case '/expected.html':
			new Expected(api);
			break;
		case '/scans.html':
			new Scans(api);
			break;
		default:
			new Dashboard(api);
	}
});
