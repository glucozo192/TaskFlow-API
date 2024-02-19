CREATE TABLE IF NOT EXISTS "users"(
  id text PRIMARY KEY,
  username text UNIQUE,
  email text UNIQUE,
  password text,
  role_id text,
  token text,
  created_by text,
  created_at timestamptz DEFAULT now(),
  updated_at timestamptz DEFAULT NOW(),
  deleted_at timestamptz
);

CREATE TABLE public.roles (
    id text PRIMARY KEY,
    name text UNIQUE,
    display_name text,
    created_by text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone
);

CREATE TABLE public.role_permissions (
    id text PRIMARY KEY,
    role_id text NOT NULL,
    path text,
    name text,
    created_at timestamp with time zone DEFAULT now(),
    updated_at timestamp with time zone DEFAULT now(),
    deleted_at timestamp with time zone,
    FOREIGN KEY ("role_id") REFERENCES "roles"("id") ON DELETE CASCADE
);

