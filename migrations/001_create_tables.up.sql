CREATE TABLE users (
	id serial PRIMARY KEY,
	reputation INT,
	views INT,
	upvotes INT,
	downvotes INT
);

CREATE TABLE questions (
	id serial PRIMARY KEY,
	creation_date TIMESTAMP WITHOUT TIME ZONE,
	score INT,
	view_count INT,
	owner_user_id INT,
	comment_count INT,
	closed_date TIMESTAMP WITHOUT TIME ZONE,
	favorite_count INT,
	CONSTRAINT fk_owner_user
		FOREIGN KEY(owner_user_id)
			REFERENCES users(id)
);

CREATE TABLE question_tags (
	id serial PRIMARY KEY,
	question_id INT,
	tag TEXT,
	CONSTRAINT fk_question_id
		FOREIGN KEY(question_id)
			REFERENCES questions(id)
);

CREATE TABLE answers (
	id serial PRIMARY KEY,
	creation_date TIMESTAMP WITHOUT TIME ZONE,
	score INT,
	owner_user_id INT,
	comment_count INT,
	closed_date TIMESTAMP WITHOUT TIME ZONE,
	CONSTRAINT fk_owner_user
		FOREIGN KEY(owner_user_id)
			REFERENCES users(id)
);

CREATE TABLE question_answers (
	id serial PRIMARY KEY,
	question_id INT,
	answer_id INT,
	accepted INT,
	CONSTRAINT fk_question_id
		FOREIGN KEY(question_id)
			REFERENCES questions(id),
	CONSTRAINT fk_answer_id
		FOREIGN KEY(answer_id)
			REFERENCES answers(id)
);
