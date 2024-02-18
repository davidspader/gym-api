INSERT INTO users (name, email, password) VALUES
    ('User1', 'user1@gmail.com', '$2a$10$0iGYlKCAYTyJV/vC6nLGgeWFwD6AhSkWLsVRO/.M4lNK8OtIkfggy'),
    ('User2', 'user2@gmail.com', '$2a$10$0iGYlKCAYTyJV/vC6nLGgeWFwD6AhSkWLsVRO/.M4lNK8OtIkfggy');

INSERT INTO workouts (user_id, name) VALUES
    (1, 'Workout1'),
    (2, 'Workout2'),
    (1, 'Workout3'),
    (2, 'Workout4');

INSERT INTO exercises (user_id, name, weight, reps) VALUES
    (1, 'Exercise1', 10, 12),
    (1, 'Exercise2', 15, 10),
    (1, 'Exercise3', 20, 8),
    (2, 'Exercise4', 12, 15),
    (2, 'Exercise5', 18, 10),
    (2, 'Exercise6', 25, 8);

INSERT INTO exercises_workout (workout_id, exercise_id) VALUES
    (1, 1),
    (1, 2),
    (1, 4),
    (2, 2),
    (2, 3),
    (2, 5),
    (3, 1),
    (3, 3),
    (3, 6),
    (4, 4),
    (4, 5),
    (4, 6);
