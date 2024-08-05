drop trigger if exists flat_create_trigger on flats;
drop function if exists insert_flat_to_outbox;

drop table if exists new_flats_outbox;
drop table if exists subscribers;
drop table if exists flats;
drop table if exists houses;
drop table if exists users;

drop type if exists user_role;
drop type if exists flat_status;
drop type if exists flat_update_msg_status;

