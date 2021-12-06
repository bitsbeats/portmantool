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
