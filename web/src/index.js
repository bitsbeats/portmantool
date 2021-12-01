import Portmantool from './Portmantool';
import './style.scss';

document.addEventListener('DOMContentLoaded', () => {
	const api = new Portmantool();
	const target = document.getElementById('target');

	const doUpdate = async () => {
		try {
			const diff = await api.get('diff');
			target.replaceChildren(...diff.map(data => {
				const tr = document.createElement('tr');
				tr.replaceChildren(...[
					'address',
					'port',
					'protocol',
					'expected_state',
					'actual_state',
					'scan_id',
					'comment',
				].map(property => {
					const td = document.createElement('td');
					if (property === 'scan_id') {
						td.innerText = new Date(data[property]).toLocaleString();
					} else {
						td.innerText = data[property];
					}

					return td;
				}));

				return tr;
			}));
		} catch (error) {
			console.log(error);
		}

		setTimeout(doUpdate, 5000);
	};

	doUpdate();
});
