CREATE TABLE `users` (
  `i_d` BIGINT NOT NULL,
  `created_at` DATETIME,
  `updated_at` DATETIME,
  `deleted_at` TEXT,
  `username` VARCHAR(50) NOT NULL,
  `email` TEXT,
  `password` VARCHAR(255) NOT NULL,
  `status` INT DEFAULT 1,
  PRIMARY KEY (`i_d`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;