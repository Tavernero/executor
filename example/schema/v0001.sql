
--+Migrate Down

delete from table "robot";
delete from table "endpoint";
delete from table "task";

--+Migrate Up

insert into "task" ( "version", "context", "function", "step", "status", "buffer" )
    values ( '1', 'api.com', 'toto_function', 'starting', 'TODO', '{"toto":"tata"}' );

insert into "endpoint" ( "id", "version", "name", "method", "url" )
    values (
        '0be72df3-3e61-4a87-9e60-bdf79a350834',
        '1',
        'toto starting',
        'GET',
        'http://www.mocky.io/v2/5a491783310000d609a82145'
    ), (
        '091a6422-bd2f-4b5b-99d2-26e6c8f4a246',
        '1', 
        'toto begin', 
        'GET', 
        'http://www.mocky.io/v2/5a491783310000d609a82145'
    ), (
        'c808480c-ba21-45b2-b223-e3c942d952b2',
        '1', 
        'toto next', 
        'GET', 
        'http://www.mocky.io/v2/5a491783310000d609a82145'
    );

insert into "robot" ( "id", "function", "version", "status", "definition" )
    values ( 
        '2a477996-0610-4c9b-b228-7b44a969e23a', 
        'toto_function', 
        '1', 
        'ACTIVE', 
        '{"sequence":[{"name":"starting","endpoint_id":"0be72df3-3e61-4a87-9e60-bdf79a350834"},{"name":"start","endpoint_id":"091a6422-bd2f-4b5b-99d2-26e6c8f4a246"}]}'
    );

