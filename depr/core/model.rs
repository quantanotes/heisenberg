use serde::{Serialize, Deserialize};

const MAX_KEY_SIZE: usize = 64;

pub type Key = [u8; MAX_KEY_SIZE];

pub struct Entry {
    pub collection: String,
    pub key: String,
    pub value: Value,
}

#[derive(Serialize, Deserialize)]
pub struct Value {
    pub vector: Vec<f32>,
    pub meta: serde_json::Value,
}
