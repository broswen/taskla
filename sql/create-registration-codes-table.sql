CREATE TABLE IF NOT EXISTS public.registration_codes
(
    code text NOT NULL,
    expiration timestamp without time zone NOT NULL,
    used boolean NOT NULL,
    PRIMARY KEY (code)
);

ALTER TABLE public.registration_codes
    OWNER to taskla;