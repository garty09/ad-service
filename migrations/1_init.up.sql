CREATE TABLE ad
(
    id          SERIAL PRIMARY KEY,
    created_at  TIMESTAMP     NOT NULL,
    title       VARCHAR(200)  NOT NULL,
    description VARCHAR(1000) NOT NULL,
    price       INTEGER       NOT NULL,
    photo_links VARCHAR[3]
);
