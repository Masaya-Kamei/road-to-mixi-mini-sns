INSERT INTO users (user_id, name) VALUES
  (1, 'user1'), (2, 'user2'), (3, 'user3'), (4, 'user4'), (5, 'user5');

INSERT INTO friend_link (user1_id, user2_id) VALUES
  (1, 2), (1, 3), (1, 4), (2, 3), (2, 4), (3, 4);

INSERT INTO block_list (user1_id, user2_id) VALUES
  (1, 2), (1, 3), (1, 4), (2, 1);
