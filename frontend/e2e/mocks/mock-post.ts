import axios from 'axios';
import express from 'express';
import cors from 'cors';

const mock = express();
mock.use(express.json());
mock.use(
	cors({
		origin: 'http://localhost:5173',
		credentials: true
	})
);
const posts: { id: string; caption: string; text: string; attached?: string[] }[] = [];

mock.get('/', (_, res) => {
	return res.send("Not just alive, but alive and kickin'");
});

mock.get('/api/post', (req, res) => {
	let page;
	let size;

	page = req.query.page;
	if (!page) {
		page = 0;
	}
	page = Number(page);

	size = req.query.size;
	if (!size) {
		size = 0;
	}
	size = Number(size);

	return res.status(200).json({
		items: posts
			.slice(Number(page) * Number(size), Number(page + 1) * Number(size))
			.map(({ id, caption, text, attached }) => {
				return {
					id: id,
					caption: caption,
					text: text,
					attached: attached
				};
			}),
		page: page,
		size: size,
		total: posts.length
	});
});

mock.post('/api/post', async (req, res) => {
	if (!req.headers.authorization) {
		console.log('no auth header');
		return res.status(401).send('no auth header');
	}

	const token = req.get('Authorization')?.split(' ')[1];
	const response = await axios.get(`http://localhost:5003/api/user/${token}`);
	if (response.status === 401) {
		console.log('Bad token');
		return res.status(401).send('Bad token');
	}

	const { caption, text, attachment } = req.body;

	posts.push({ id: crypto.randomUUID(), caption: caption, text: text, attached: attachment });

	return res.status(204).send();
});

mock.post('/clear', (_, res) => {
	posts.length = 0;
	return res.status(204).send();
});

const PORT = 5004;
mock.listen(PORT, () => {
	console.log('mocking posts');
});
