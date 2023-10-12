CREATE TABLE users
(
    id            serial       not null unique,
    firstname          varchar(255) not null,
    secondname          varchar(255) not null,
    lastname          varchar(255) not null,
    email          varchar(255) not null,
    login      varchar(255) not null unique,
    pass varchar(255) not null
);
