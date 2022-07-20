CREATE TABLE IF NOT EXISTS staff_department (
id bigint NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
staffId  bigint REFERENCES staff(staffId),
departmentId bigint REFERENCES department(id)
);
