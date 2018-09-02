CREATE TABLE IF NOT EXISTS records (
    pk               CHARACTER VARYING(50) NOT NULL UNIQUE,
    score            CHARACTER VARYING(50) NOT NULL,
    created_at       TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    PRIMARY KEY (pk, score)
);
