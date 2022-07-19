CREATE TABLE IF NOT EXISTS attendance (
    id bigint NOT NULL PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    studentId bigint REFERENCES students(id),
    loginTime TIMESTAMP,
    logoutTime TIMESTAMP
);