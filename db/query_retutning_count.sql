WITH user_storage AS (
   SELECT *
   FROM  storage
   WHERE
    user_id = 1
    AND deleted_at IS NULL
   )
SELECT *
FROM  (
   TABLE  user_storage
   ORDER  BY id
   LIMIT  25
   OFFSET 0
   ) sub
RIGHT JOIN (SELECT count(*) FROM user_storage) AS c (total) ON true;

SELECT
    (SELECT COUNT(*)
   FROM  storage
   WHERE
    user_id = 1
    AND deleted_at IS NULL
    ) as count,
    (SELECT json_agg(t.*) FROM (
        SELECT * FROM storage
        WHERE
            user_id = 1
    AND deleted_at IS NULL
        ORDER BY id
        LIMIT 25
        OFFSET 0
    ) AS t) AS rows;

CREATE INDEX storage_user_id_idx ON storage(dir);

WITH user_storage AS (
   SELECT *
   FROM  storage
   WHERE
    user_id = 1
    AND deleted_at IS NULL
   )
SELECT *
FROM  (
   TABLE  user_storage
   ORDER  BY id DESC
   LIMIT  1
   OFFSET 1
   ) sub
RIGHT JOIN (SELECT count(*) FROM user_storage) AS c (total) ON true;
