ALTER TABLE comments
    ADD CONSTRAINT comments_post_fk
        FOREIGN KEY (post_id)
        REFERENCES posts(id)
        ON DELETE CASCADE,
    ADD CONSTRAINT comments_user_fk
        FOREIGN KEY (user_id)
        REFERENCES users(id);