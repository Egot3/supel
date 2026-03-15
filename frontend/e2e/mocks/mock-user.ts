import express from 'express';

console.log('entered mocking server');

const mock = express();
mock.use(express.json());

let users: { [key: string]: { password: string; nickname: string; token: string } } = {};

const t =
	'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c';

mock.get('/', (req, res) => {
	res.send('alive');
});

console.log('set up health');
mock.post('/api/user/register', (req, res) => {
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

	res.status(200).json({
		token: t
	});
});

mock.post('/api/user/login', (req, res) => {
	const { email, password } = req.body;

	const passwordMatch = users[email] && users[email].password === password;

	if (passwordMatch) {
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

mock.get('/api/user/:token', (req, res) => {
	const reqUserToken = req.params.token;

	const found = Object.values(users).some((userData) => {
		return userData.token === reqUserToken;
	});

	console.log('user found: ', found);

	res.status(found ? 204 : 401).send('No such user');
});

mock.post('/clear', (_, res) => {
	users = {} as { [key: string]: { password: string; nickname: string; token: string } };
	res.status(204).send();
});

const PORT = 5003;
mock.listen(PORT, () => {
	console.log('mocking ', PORT);
});
