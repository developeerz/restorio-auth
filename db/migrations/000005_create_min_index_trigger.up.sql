CREATE OR REPLACE FUNCTION find_min_free_id() RETURNS TRIGGER AS $$
DECLARE
    min_free_id INT;
BEGIN
    SELECT MIN(id) + 1 INTO min_free_id
    FROM users
    WHERE id + 1 NOT IN (SELECT id FROM users);

    IF min_free_id IS NULL THEN
        SELECT COALESCE(MAX(id), 0) + 1 INTO min_free_id FROM users;
    END IF;

    NEW.id = min_free_id;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_before_insert
BEFORE INSERT ON users
FOR EACH ROW
EXECUTE FUNCTION find_min_free_id();