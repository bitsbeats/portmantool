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
