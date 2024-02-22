CREATE TYPE permissions AS ENUM ('DOCTOR', 'PATIENT', 'ADMIN','MANAGER');

CREATE TYPE verificationStatus AS ENUM ('PENDING', 'VERIFIED');

CREATE TABLE IF NOT EXISTS "users" (
     "id" SERIAL NOT NULL,
     "email" VARCHAR(200),
     "phone" VARCHAR(20),
     "password" VARCHAR(200),
     "permissions" permissions NOT NULL,
     "status" verificationStatus DEFAULT 'PENDING',
     "created_at" TIMESTAMP WITHOUT TIME ZONE DEFAULT CURRENT_TIMESTAMP,
     PRIMARY KEY ("id")
);

-- Create unique index for (id, email)
CREATE UNIQUE INDEX idx_unique_id_email ON users (id, email);

-- Create unique index for (id, phone)
CREATE UNIQUE INDEX idx_unique_id_phone ON users (id, phone);

