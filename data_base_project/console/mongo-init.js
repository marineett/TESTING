// Create database and user
db = db.getSiblingDB('admin');

// Create admin user
db.createUser({
  user: 'postgres',
  pwd: 'postgres',
  roles: [
    { role: 'userAdminAnyDatabase', db: 'admin' },
    { role: 'dbAdminAnyDatabase', db: 'admin' },
    { role: 'readWriteAnyDatabase', db: 'admin' }
  ]
});

// Create application database
db = db.getSiblingDB('data_base_project_db');
db.createCollection('init');

// Initialize sequence for personal data IDs
db.counters.insertOne({
  _id: "personal_data_id",
  seq: 0
}); 