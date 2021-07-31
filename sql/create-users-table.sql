CREATE TABLE IF NOT EXISTS public.users
(
    username text NOT NULL,
    password text NOT NULL,
    role text NOT NULL,
    PRIMARY KEY (username)
);

ALTER TABLE public.users
    OWNER to taskla;