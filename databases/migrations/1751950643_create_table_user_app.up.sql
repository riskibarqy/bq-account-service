CREATE TABLE public."user_app" (
    "id" SERIAL PRIMARY KEY,
    "user_id" INT NOT NULL REFERENCES public."user"("id") ON DELETE CASCADE,
    "app_id" INT NOT NULL REFERENCES public."app"("id") ON DELETE CASCADE,
    "role" VARCHAR(50) DEFAULT 'user',
    "metadata" JSONB,
    "joined_at" INT NOT NULL,
    "created_at" INT NOT NULL,
    "updated_at" INT NOT NULL,
    "deleted_at" INT,
    UNIQUE ("user_id", "app_id")  -- prevent duplicate relations
);
CREATE INDEX user_app_user_id_idx ON public."user_app"("user_id");
CREATE INDEX user_app_app_id_idx ON public."user_app"("app_id");