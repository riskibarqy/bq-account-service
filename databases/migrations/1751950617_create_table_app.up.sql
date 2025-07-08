CREATE TABLE public."app" (
    "id" SERIAL PRIMARY KEY,
    "name" VARCHAR(100) NOT NULL,
    "slug" VARCHAR(50) NOT NULL UNIQUE,  -- e.g., 'budgetbuddy'
    "created_at" INT NOT NULL,
    "updated_at" INT NOT NULL,
    "deleted_at" INT
);
CREATE INDEX app_slug_idx ON public."app"("slug");