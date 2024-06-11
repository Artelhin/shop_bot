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
      (now(), now(), null, 'Samsung A70', '5 y.o. device', 5, null),
      (now(), now(), null, 'Samsung S24 Ultra', 'this year device' ||
                                                '+nice camera', 5, null);