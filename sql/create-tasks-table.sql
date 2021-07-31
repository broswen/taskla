CREATE TABLE IF NOT EXISTS public.tasks
(
    id bigserial NOT NULL,
    username text NOT NULL,
    group_id bigint NOT NULL,
    name text NOT NULL,
    description text,
    status text NOT NULL,
    PRIMARY KEY (id),
    CONSTRAINT username FOREIGN KEY (username)
        REFERENCES public.users (username) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID,
    CONSTRAINT group_id FOREIGN KEY (group_id)
        REFERENCES public.groups (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE CASCADE
        NOT VALID
);

ALTER TABLE public.tasks
    OWNER to taskla;