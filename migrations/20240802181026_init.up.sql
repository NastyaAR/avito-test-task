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
    moderator_id uuid references users(user_id),
    primary key (flat_id, house_id)
);

create index houses_id_on_flats
    on flats (house_id);

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
    subscriber_mail text;
begin
    select mail into subscriber_mail
    from subscribers join users u on u.user_id = subscribers.user_id
    where subscribers.user_id=NEW.user_id and house_id=NEW.house_id;

    if subscriber_mail is null then
        return NEW;
    end if;

    insert into new_flats_outbox(flat_id, house_id, mail, status) values (NEW.flat_id, NEW.house_id, subscriber_mail, 'no send');
    return NEW;
end;
$$ language plpgsql;

CREATE TRIGGER flat_create_trigger
    AFTER INSERT ON flats
    FOR EACH ROW
EXECUTE FUNCTION insert_flat_to_outbox();

create or replace  function check_exists_subscriber()
    returns trigger as $$
declare
    usr uuid;
begin
    select subscribers.user_id into usr from subscribers
    where subscribers.user_id=NEW.user_id and subscribers.house_id=NEW.house_id;

    if usr is not null then
        return NULL;
    end if;

    return NEW;
end;
$$ language plpgsql;

CREATE TRIGGER insert_subscribe_trigger
    BEFORE INSERT ON subscribers
    FOR EACH ROW
EXECUTE FUNCTION check_exists_subscriber();

create or replace  function update_status(new_status flat_status, new_flat_id int, new_house_id int, new_moderator_id uuid)
    returns setof flats as $$
declare
    mod_id uuid;
begin
    if new_status = 'on moderation' then
        select flats.moderator_id into mod_id from flats
        where flats.flat_id=new_flat_id and flats.house_id=new_house_id;

        if mod_id != new_moderator_id then
            raise exception 'flat already on moderation';
        end if;
    end if;

    return query
        update flats set status=new_status, moderator_id=new_moderator_id
            where flats.flat_id=new_flat_id and flats.house_id=new_house_id;
end;
$$ language plpgsql;