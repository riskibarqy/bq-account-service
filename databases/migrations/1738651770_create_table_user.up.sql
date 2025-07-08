CREATE TABLE public."user" (
    "id" SERIAL PRIMARY KEY,
    "clerk_id" VARCHAR(100) NOT NULL UNIQUE,
    "email" VARCHAR(100) NOT NULL UNIQUE,
    "name" VARCHAR(100) NOT NULL,
    "username" VARCHAR(50) NOT NULL UNIQUE,  -- e.g., 'john_doe'
    "phone" VARCHAR(20) NOT NULL UNIQUE,  -- e.g., '+6281234567890'
    "is_active" BOOLEAN NOT NULL DEFAULT TRUE,  -- Indicates if the user is active
    "is_verified" BOOLEAN NOT NULL DEFAULT FALSE,  -- Indicates if the user has verified
    "created_at" INT NOT NULL,
    "updated_at" INT NOT NULL,
    "deleted_at" INT
);

CREATE INDEX user_email_idx ON public."user"("email");
CREATE INDEX user_deleted_at_idx ON public."user"("deleted_at");
