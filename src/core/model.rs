use serde::{Deserialize, Serialize};

#[derive(Deserialize, Serialize)]
pub struct Value {
    pub vector: Vec<f32>,
    pub meta: serde_json::Value,
}

pub struct Entry {
    pub collection: String,
    pub key: String,
    pub value: Value,
}