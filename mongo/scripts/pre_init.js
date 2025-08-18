db = db.getSiblingDB('fhir_hca');
db.createCollection("encounters", {
  validator: {
    $jsonSchema: {
      bsonType: "object",
      required: ["fhirId", "fullUrl", "status", "class", "period"],
      properties: {
        fhirId: { bsonType: "string", description: "Hapi Api FhirID" },
        fullUrl: { bsonType: "string", description: "FullUrl of the resource at Api Hapi" },
        status: {
          enum: ["planned", "in-progress", "on-hold", "discharged", "completed", "finished", "cancelled", "discontinued", "entered-in-error", "unknown"],
          description: "Encounter Status"
        },
        class: {
          enum: ["IMP", "AMB", "OBSENC", "EMER", "VR", "HH"],
          description: "Encounter Class"
        },
        period: {
          bsonType: "object",
          required: ["start"],
          properties: {
            start: { bsonType: "date" },
            end: { bsonType: ["date", "null"] }
          }
        },
        practitionerId: {bsonType: "objectId", description: "Internal Reference to Practitioner Resource"},
        patientId: {bsonType: "objectId", description: "Internal Reference to Patient Resource"}
      }
    }
  },
  validationLevel: "strict", 
  validationAction: "error"  
});


db.createCollection('patients');
db.createCollection('practitioners');

db = db.getSiblingDB('fhir_hcb');
db.createCollection('encounters');
db.createCollection('patients');
db.createCollection('practitioners');
