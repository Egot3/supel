import axios from 'axios';
import express from 'express';

const mock = express();
mock.use(express.json());

const posts: { caption: string; text: string; attached: string[] }[] = [];

mock.get('/', (_, res) => {
	res.send("Not just alive, but alive and kickin'");
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

	res.status(200).json({
		items: posts.slice(Number(page) * Number(size), Number(page + 1) * Number(size)),
		page: page,
		size: size,
		total: posts.length
	});
});

mock.post('/api/post', async (req, res) => {
	const token = req.get('Authorization')?.split(' ')[1];
	const response = await axios.get(`http://localhost:5003/${token}`);
	if (response.status === 204) {
		res.status(401);
	}

	const { caption, text, attachment } = req.body;

	posts.push({ caption: caption, text: text, attached: attachment });

	res.status(204);
});

mock.post('/clear', (_, res) => {
	posts.length = 0;
	res.status(204);
});

const PORT = 5004;
mock.listen(PORT, () => {
	console.log('mocking posts');
});
