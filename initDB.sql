drop table users
drop table session

create table if not exists users (
    id serial primary key,
    name varchar(64) not null unique,
    email varchar(64) not null unique,
    password varchar(64) not null,
    avatar varchar(255) default '',
);

create table if not exists session (
    id varchar(64) not null primary key ,
    userID int not null,
    expire date not null,
    foreign key (userID) references users (id) on delete cascade
);
