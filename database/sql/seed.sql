INSERT INTO users (id, name, email, role, password_hash, created_at, updated_at) VALUES
	('ae127240-ce13-4789-aafd-d2f31e7ee487', 'Admin', 'admin@govod.com', 'ADMIN', '$2a$10$zlgi68l46JRGskhBjd2TiOyKMuNI.kbUmBLN.EDRCl.g8s0Da2Qsm', '2022-09-16 00:00:00', '2022-09-16 00:00:00'),
	('45b5fbd3-755f-4379-8f07-a58d4a30fa2f', 'User Test', 'user-test@govod.com', 'USER', '$2a$10$g0qa7gfKnfyapT5fPk7CMuYFpyzyVhZ.UXu04Bfr.rc0EKiJNyk5u', '2019-03-24 00:00:00', '2019-03-24 00:00:00')
	ON CONFLICT DO NOTHING;
