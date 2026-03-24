ALTER TABLE finances
DROP CONSTRAINT IF EXISTS finances_type_check;

ALTER TABLE finances
ADD CONSTRAINT finances_type_check
CHECK ("type" IN ('salary', 'freelance', 'insta_paid_collab', 'other', 'spend'));
