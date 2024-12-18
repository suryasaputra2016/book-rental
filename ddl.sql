--populating books
INSERT INTO books (isbn, title, author, category, rental_cost) VALUES
-- Comic books
('9781234567897', 'Comic Adventures Vol. 1', 'John Doe', 'comic', 10000.00),
('9781234567800', 'Superhero Saga', 'Jane Smith', 'comic', 15000.00),

-- Novels
('9780987654321', 'Mystery Novel', 'Arthur Conan', 'novel', 20000.00),
('9781122334455', 'Romantic Tale', 'Emily Brontë', 'novel', 18000.00),

-- Biographies
('9782233445566', 'Life of a Genius', 'Albert Einstein', 'biography', 25000.00),
('9783344556677', 'Famous Artist', 'Vincent Van Gogh', 'biography', 24000.00),

-- Art books
('9784455667788', 'Modern Art Trends', 'Pablo Picasso', 'art', 30000.00),
('9785566778899', 'Renaissance Art', 'Leonardo da Vinci', 'art', 28000.00),

-- Textbooks
('9786677889900', 'Physics 101', 'Richard Feynman', 'textbook', 29000.00),
('9787788990011', 'Chemistry Basics', 'Marie Curie', 'textbook', 27000.00),
('9788899001122', 'Mathematics Simplified', 'Carl Gauss', 'textbook', 22000.00);


-- populating book copies
INSERT INTO book_copies (book_id, copy_number, status) VALUES
-- Comic Adventures Vol. 1 (ID: 1)
(1, 1, 'available'),
(1, 2, 'available'),

-- Superhero Saga (ID: 2)
(2, 1, 'available'),

-- Mystery Novel (ID: 3)
(3, 1, 'available'),
(3, 2, 'available'),

-- Romantic Tale (ID: 4)
(4, 1, 'available'),

-- Life of a Genius (ID: 5)
(5, 1, 'available'),

-- Famous Artist (ID: 6)
(6, 1, 'available'),
(6, 2, 'available'),

-- Modern Art Trends (ID: 7)
(7, 1, 'available'),

-- Renaissance Art (ID: 8)
(8, 1, 'available'),

-- Physics 101 (ID: 9)
(9, 1, 'available'),
(9, 2, 'available'),

-- Chemistry Basics (ID: 10)
(10, 1, 'available'),

-- Mathematics Simplified (ID: 11)
(11, 1, 'available'),
(11, 2, 'available');