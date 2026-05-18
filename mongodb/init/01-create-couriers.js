const database = db.getSiblingDB("gradin-courier");

try {
  rs.initiate({
    _id: "rs0",
    members: [{ _id: 0, host: "mongodb:27017" }]
  });
} catch (error) {
  if (!String(error).includes("already initialized")) {
    print(`Replica set initialization skipped: ${error}`);
  }
}

for (let attempt = 0; attempt < 30; attempt += 1) {
  try {
    if (rs.status().myState === 1) {
      break;
    }
  } catch (error) {
    print(`Waiting for replica set primary: ${error}`);
  }
  sleep(1000);
}

const existingCollections = database.getCollectionNames();

if (!existingCollections.includes("couriers")) {
  database.createCollection("couriers", {
    validator: {
      $jsonSchema: {
        bsonType: "object",
        required: ["name", "level", "status", "registered_at", "created_at", "updated_at"],
        properties: {
          name: {
            bsonType: "string"
          },
          email: {
            bsonType: ["string", "null"]
          },
          phone: {
            bsonType: ["string", "null"]
          },
          level: {
            bsonType: ["int", "long"],
            minimum: 1,
            maximum: 5
          },
          vehicle_type: {
            bsonType: ["string", "null"]
          },
          license_plate: {
            bsonType: ["string", "null"]
          },
          status: {
            enum: ["active", "inactive", "suspended"]
          },
          registered_at: {
            bsonType: "date"
          },
          created_at: {
            bsonType: "date"
          },
          updated_at: {
            bsonType: "date"
          },
          deleted_at: {
            bsonType: ["date", "null"]
          }
        }
      }
    },
    validationLevel: "moderate"
  });
}

database.couriers.createIndex({ name: 1 }, { name: "couriers_name_idx" });
database.couriers.createIndex({ level: 1 }, { name: "couriers_level_idx" });
database.couriers.createIndex({ registered_at: 1 }, { name: "couriers_registered_at_idx" });
database.couriers.createIndex(
  { email: 1 },
  {
    name: "couriers_email_unique_sparse_idx",
    unique: true,
    sparse: true
  }
);
