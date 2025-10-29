db = db.getSiblingDB('admin');

db.createUser({
  user: 'postgres',
  pwd: 'postgres',
  roles: [
    { role: 'userAdminAnyDatabase', db: 'admin' },
    { role: 'dbAdminAnyDatabase', db: 'admin' },
    { role: 'readWriteAnyDatabase', db: 'admin' }
  ]
});

db = db.getSiblingDB('data_base_project_db');
db.createCollection('init');

db.counters.insertOne({
  _id: "personal_data_id",
  seq: 0
}); 