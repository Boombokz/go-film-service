CREATE TABLE public.genres (
	id serial4 NOT NULL,
	title text NULL,
	CONSTRAINT genres_pkey PRIMARY KEY (id)
);

CREATE TABLE public.movies (
	id serial4 NOT NULL,
	title text NULL,
	description text NULL,
	release_year int4 NULL,
	director text NULL,
	rating int4 DEFAULT 0 NULL,
	is_watched bool DEFAULT false NULL,
	trailer_url text NULL,
	poster_url text NULL,
	CONSTRAINT movies_pkey PRIMARY KEY (id)
);

CREATE TABLE public.users (
	id serial4 NOT NULL,
	"name" text NOT NULL,
	email text NOT NULL,
	password_hash text NOT NULL,
	CONSTRAINT users_email_key UNIQUE (email),
	CONSTRAINT users_pkey PRIMARY KEY (id)
);

CREATE TABLE public.movies_genres (
	movie_id int4 NULL,
	genre_id int4 NULL,
	CONSTRAINT movies_genres_genre_id_fkey FOREIGN KEY (genre_id) REFERENCES public.genres(id),
	CONSTRAINT movies_genres_movies_id_fkey FOREIGN KEY (movie_id) REFERENCES public.movies(id)
);

CREATE TABLE public.watch_list (
	id int4 DEFAULT nextval('watch_queue_id_seq'::regclass) NOT NULL,
	movie_id int4 NOT NULL,
	added_at timestamp DEFAULT now() NULL,
	CONSTRAINT watch_queue_pkey PRIMARY KEY (id),
	CONSTRAINT watch_queue_movie_id_fkey FOREIGN KEY (movie_id) REFERENCES public.movies(id)
);

insert into users (name, email, password_hash)
values ('test', 'test@test.kz', '$2a$10$icUrnO4tI.v6JHXMNe4MR.TO0LPFNcq5clSbnU5RzD.o8zaeQnHCW')