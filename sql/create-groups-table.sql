CREATE TABLE IF NOT EXISTS public.groups
(
    id bigserial NOT NULL,
    username text NOT NULL,
    name text NOT NULL,
    description text,
    PRIMARY KEY (id),
    CONSTRAINT username FOREIGN KEY (username)
        REFERENCES public.users (username) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
);

ALTER TABLE public.groups
    OWNER to taskla;