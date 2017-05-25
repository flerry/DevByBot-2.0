create table IF NOT EXISTS bot_users
(
  id int auto_increment
    primary key,
  chat_id int not null,
  is_banned tinyint(1) default '0' null,
  constraint bot_users_chat_id_uindex
  unique (chat_id)
)
;

