CREATE TABLE personaldata (
  id SERIAL UNIQUE NOT NULL,
  name VARCHAR(10) NOT NULL,
  creditcard TEXT,
  address TEXT NOT NULL,
  postcode VARCHAR(6) NOT NULL
);

INSERT INTO personaldata (
    name, creditcard, address, postcode
)
SELECT
    left(md5(i::text), 10),
    md5(random()::text),
    md5(random()::text),
    left(md5(random()::text), 6)
FROM generate_series(1, 1000000) s(i)
