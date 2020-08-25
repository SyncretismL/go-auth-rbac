CREATE TABLE public.users (
    id bigserial PRIMARY KEY,
    login text UNIQUE NOT NULL,
    password text NOT NULL,
    role text NOT NULL,
    created_at timestamp NOT NULL
);

CREATE TABLE public.sessions (
    id bigserial PRIMARY KEY,
    user_id bigserial UNIQUE NOT NULL,
    token text NOT NULL,
    created_at timestamp NOT NULL,
    valid_until timestamp NOT NULL,
    FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE CASCADE ON UPDATE CASCADE
);

INSERT INTO public.users (login, password, role, created_at) VALUES ('admin', '21232f297a57a5a743894a0e4a801fc3', 'admin', now());
