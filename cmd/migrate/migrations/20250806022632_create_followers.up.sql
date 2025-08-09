CREATE TABLE IF NOT EXISTS followers(
    user_id bigserial NOT NULL,
    follower_id bigserial NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),

    PRIMARY KEY(user_id, follower_id), --> composite key

    CONSTRAINT fk_followers_user FOREIGN KEY (user_id)
        REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_followers_follower FOREIGN KEY (follower_id)
        REFERENCES users(id) ON DELETE CASCADE
)