truncate table categories;

alter sequence categories_id_seq restart;

insert into categories (
    created_at, updated_at, deleted_at, parent_id, name
) values
      (now(), now(), null, null, 'Открыть'),
      (now(), now(), null, 1, 'Смартфоны'),
      (now(), now(), null, 1, 'Телевизоры'),
      (now(), now(), null, 1, 'Ноутбуки'),
      (now(), now(), null, 2, 'Samsung'),
      (now(), now(), null, 2, 'iPhone'),
      (now(), now(), null, 2, 'Huawei'),
      (now(), now(), null, 4, 'Acer');

truncate table items;

alter sequence items_id_seq restart;

insert into items (
    created_at, updated_at, deleted_at, name, description, category_id, image
) values
      (now(), now(), null, 'Samsung A70', 'Samsung Galaxy A70 получился одним из самых красивых смартфонов на рынке, но в то же самое время его корпус непрактичный. Он имеет глянцевое покрытие, которое быстро теряет блеск из-за многочисленных отпечатков пальцев на поверхности.', 5, null),
      (now(), now(), null, 'Samsung S24 Ultra', 'Galaxy S24 Ultra с инновационной 200 МП камерой и встроенным искусственным интеллектом устанавливает новый стандарт качества съемки. Новый процессор ProVisual распознает объекты, улучшает цветовой тон, уменьшая шум и подчеркивая детали. Наслаждайтесь каждым снимком, снятым на Galaxy S24 Ultra.', 5, null),
      (now(), now(), null, 'iPhone 14 Pro Max', 'Устройство обладает модулями беспроводной передачи данных NFC, Bluetooth 5.3 и Wi-Fi. Размеры смартфона — 160,7х77,6х7,85 см, вес — 240 г. Смартфон также имеет защиту от воды и пыли по стандарту IP68, что делает его надёжным и долговечным устройством.', 6, null);

truncate table storages;

alter sequence storages_id_seq restart;

insert into storages (
    created_at, updated_at, deleted_at, name, address
) values
      (now(), now(), null, 'Авиапарк', 'м. Авиапарк'),
      (now(), now(), null, 'МТС Первомайская', 'Сиреневый б-р 62');

truncate table item_to_storage;

insert into item_to_storage (
    item_id, storage_id, count
) values
      (1, 1, 5),
      (1, 2, 4),
      (3, 1, 10);

truncate table orders;

alter sequence orders_id_seq restart;