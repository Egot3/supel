import axios from 'axios';
import express from 'express';

const mock = express();
mock.use(express.json());

const posts: { id: string; caption: string; text: string; attached?: string[] }[] = [];

mock.get('/v1/news', (req, res) => {
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

mock.post('/v1/news', async (req, res) => {
	console.log('posting from backend');
	if (!req.headers.authorization) {
		console.log('no auth header');
		return res.status(401).send('no auth header');
	}

	const token = req.get('Authorization')?.split(' ')[1];
	const response = await axios.get(`http://localhost/v1/public/validate/${token}`);
	if (response.status === 401) {
		console.log('Bad token');
		return res.status(401).send('Bad token');
	}

	const { caption, text, attachment } = req.body;

	posts.push({ id: crypto.randomUUID(), caption: caption, text: text, attached: attachment });

	console.log('sending 204');
	return res.status(204).send();
});

mock.post('/clearPost', (_, res) => {
	posts.length = 0;
	return res.status(204).send();
});

let users: { [key: string]: { password: string; nickname: string; token: string } } = {};

const t =
	'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c';

mock.get('/', (req, res) => {
	res.send('alive');
});

console.log('set up health');
mock.post('/v1/public/register', (req, res) => {
	const { email, password, nickname } = req.body;

	console.log(email, password, nickname, 'registering');

	if (users[email]) {
		res.status(409).json({
			error: 'EMAIL_ALREADY_IN_USE',
			errorMessage: 'Email already in use'
		});
		return;
	}

	console.log('are right');
	users[email] = {
		password: password,
		nickname: nickname,
		token: t
	};

	console.log('returning with 200');
	res.cookie('auth_token', t, {
		httpOnly: true,
		maxAge: 36000, //10h
		sameSite: 'strict'
	});
	res.status(200).json({
		token: t
	});
});

mock.post('/v1/public/login', (req, res) => {
	const { email, password } = req.body;

	const passwordMatch = users[email] && users[email].password === password;

	if (passwordMatch) {
		res.cookie('auth_token', t, {
			httpOnly: true,
			maxAge: 36000, //10h
			sameSite: 'strict'
		});
		res.status(200).json({
			token: t
		});
		return;
	} else {
		res.status(401).json({
			error: 'BAD_CREDITENTIALS',
			errorMessage: 'Session terminated, bad password or email'
		});
	}
});

mock.get('/v1/public/validate/:token', (req, res) => {
	const reqUserToken = req.params.token;

	const found = Object.values(users).some((userData) => {
		return userData.token === reqUserToken;
	});

	console.log('user found: ', found);

	res.status(found ? 204 : 401).send('No such user');
});

mock.post('/clearUser', (_, res) => {
	users = {} as { [key: string]: { password: string; nickname: string; token: string } };
	res.status(204).send();
});

const PORT = 80;
mock.listen(PORT, () => {
	console.log('mocking posts');
});
