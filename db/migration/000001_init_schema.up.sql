CREATE EXTENSION "uuid-ossp"
	SCHEMA "public";

CREATE TABLE public.users (
	id_user uuid NOT NULL DEFAULT uuid_generate_v4(),
	"name" varchar NOT NULL,
	CONSTRAINT users_pkey PRIMARY KEY (id_user)
);


CREATE TABLE public.wallet (
	id_wallet uuid NOT NULL DEFAULT uuid_generate_v4(),
	created_date timestamp NULL DEFAULT now(),
	balance numeric NOT NULL,
	currency varchar NULL,
	users_id uuid NULL,
	CONSTRAINT wallet_pkey PRIMARY KEY (id_wallet),
	CONSTRAINT wallet_users_id_fkey FOREIGN KEY (users_id) REFERENCES public.users(id_user)
);


CREATE TABLE public."transaction" (
	id_transaction uuid NOT NULL DEFAULT uuid_generate_v4(),
	"type" varchar NULL,
	amount numeric NOT NULL,
	"date" timestamp NULL DEFAULT now(),
	wallet_id uuid NULL,
	CONSTRAINT transaction_pkey PRIMARY KEY (id_transaction),
	CONSTRAINT transaction_wallet_id_fkey FOREIGN KEY (wallet_id) REFERENCES public.wallet(id_wallet)
);


