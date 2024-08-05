create type user_role as enum ('moderator', 'client');
create type flat_status as enum ('created', 'approved', 'declined', 'on moderation');
create type flat_update_msg_status as enum ('send', 'no send');

create table users (
    user_id uuid primary key ,
    mail varchar(50) not null,
    password text not null,
    role user_role not null
);

create table houses (
    house_id serial primary key,
    address text not null,
    construct_year int,
    developer text,
    create_house_date timestamp without time zone,
    update_flat_date timestamp without time zone
);

create table flats (
    flat_id int not null,
    house_id int references houses(house_id),
    user_id uuid references users(user_id),
    price int not null,
    rooms int not null,
    status flat_status not null,
    moderator_id int,
    primary key (flat_id, house_id)
);

create table subscribers (
    user_id uuid references users(user_id),
    house_id int references houses(house_id)
);

create table new_flats_outbox (
    id serial primary key,
    flat_id int,
    house_id int,
    mail text not null ,
    status flat_update_msg_status not null,
    foreign key (flat_id, house_id) references flats(flat_id, house_id)
);

create or replace function insert_flat_to_outbox()
    returns trigger as $$
declare
    subscriberID uuid;
begin
    select user_id into subscriberID from subscribers where user_id=NEW.user_id and house_id=NEW.house_id;

    if subscriberID is null then
        return NEW;
    end if;

    insert into subscribers(user_id, house_id) values (NEW.user_id, NEW.house_id);
    return NEW;
end;
$$ language plpgsql;

CREATE TRIGGER flat_create_trigger
    AFTER INSERT ON flats
    FOR EACH ROW
EXECUTE FUNCTION insert_flat_to_outbox();

