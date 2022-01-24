export default class Scans {
	constructor(api) {
		this.api = api;

		this.keep = null;
		this.scan1 = null;
		this.scan2 = null;

		this.load();
	}

	async load() {
		try {
			const scans = await this.api.get('scans');

			document.getElementById('scans').replaceChildren(...scans.map(scan => this.createEntry(scan.id)));
		} catch (error) {
			console.log(error);
		}

		const keep = document.getElementById('keep');
		keep.addEventListener('change', event => {
			this.updateKeep(event.target);
		});
		this.updateKeep(keep);

		document.getElementById('prune').addEventListener('click', async () => {
			try {
				await this.api.prune(this.keep);
			} catch (error) {
				console.log(error);
			}
		});
	}

	async show() {
		if (this.scan1 === null) {
			return;
		}

		console.log(this.scan1);
		console.log(this.scan2);
	}

	createEntry(id) {
		const li = document.createElement('li');
		li.classList.add('is-flex');

		const scan1 = document.createElement('label');
		['radio', 'is-flex-grow-1'].forEach(cls => scan1.classList.add(cls));

		const scan1Input = document.createElement('input');
		scan1Input.name = 'scan1';
		scan1Input.type = 'radio';
		scan1Input.value = id;
		scan1Input.addEventListener('change', event => {
			this.scan1 = event.target.value;
			this.show();
		});

		const scan1Link = document.createElement('a');
		scan1Link.innerText = new Date(id*1000).toLocaleString();

		scan1.append(scan1Input);
		scan1.append(scan1Link);

		const scan2 = document.createElement('label');
		scan2.classList.add('radio');

		const scan2Input = document.createElement('input');
		scan2Input.name = 'scan2';
		scan2Input.type = 'radio';
		scan2Input.value = id;
		scan2Input.addEventListener('change', event => {
			this.scan2 = event.target.value;
			this.show();
		});
		scan2Input.addEventListener('click', event => {
			if (this.scan2 === event.target.value) {
				event.target.checked = false;
				this.scan2 = null;
				this.show();
			}
		});

		const scan2Link = document.createElement('a');
		scan2Link.innerText = '\u21c6';

		scan2.append(scan2Input);
		scan2.append(scan2Link);

		li.append(scan1);
		li.append(scan2);

		return li;
	}

	updateKeep(input) {
		this.keep = input.value === '' ? null : input.valueAsNumber / 1000;

		document.getElementById('keep-info').innerText = this.keep === null ? '' : ` older than ${new Date(this.keep*1000).toLocaleString()}`;
	}
}
