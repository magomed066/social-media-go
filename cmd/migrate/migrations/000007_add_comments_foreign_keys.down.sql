ALTER TABLE comments
    DROP CONSTRAINT comments_post_fk,
    DROP CONSTRAINT comments_user_fk;