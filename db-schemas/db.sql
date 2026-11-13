-- Database `expenses-tracker` for MariaDB.
CREATE DATABASE IF NOT EXISTS `expenses-tracker`
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_ci;

USE `expenses-tracker`;

-- Categories table
CREATE TABLE categories (
    id INT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT NULL
) ENGINE=InnoDB;

-- Expenses table
CREATE TABLE expenses (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    item VARCHAR(255) NOT NULL,
    category_id INT UNSIGNED NULL,
    price DECIMAL(12, 2) NOT NULL,
    date DATE NOT NULL,

    CONSTRAINT fk_expenses_categories
        FOREIGN KEY (category_id)
        REFERENCES categories(id)
        ON UPDATE CASCADE
        ON DELETE SET NULL,

    INDEX idx_expenses_date (date),
    INDEX idx_expenses_category (category_id)
) ENGINE=InnoDB;
