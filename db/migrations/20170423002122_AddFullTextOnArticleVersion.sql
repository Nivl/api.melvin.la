
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE blog_article_versions ADD COLUMN content_vector tsvector;
CREATE INDEX blog_article_versions_tsv_content_idx ON blog_article_versions USING gin(content_vector);
CREATE INDEX blog_article_versions_content_idx ON blog_article_versions (content);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION blog_article_versions_vector_update() RETURNS trigger AS $blog_article_versions_vector_update$
BEGIN
    IF TG_OP = 'INSERT' THEN
        NEW.content_vector = to_tsvector('pg_catalog.english', COALESCE(NEW.content, ''));
    END IF;
    IF TG_OP = 'UPDATE' THEN
        IF NEW.content <> OLD.content THEN
            NEW.content_vector = to_tsvector('pg_catalog.english', COALESCE(NEW.content, ''));
        END IF;
    END IF;
    RETURN NEW;
END
$blog_article_versions_vector_update$ LANGUAGE plpgsql;
-- +goose StatementEnd

CREATE TRIGGER blog_article_versions_vector_update BEFORE INSERT OR UPDATE ON blog_article_versions
    FOR EACH ROW EXECUTE PROCEDURE blog_article_versions_vector_update();

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

