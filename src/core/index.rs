#[derive(Debug, thiserror::Error)]
pub enum IndexError {
    #[error("index error: {0}")]
    Error(String)
}

pub trait Index {
    fn new(dimension: usize) -> Result<Self, IndexError> where Self: Sized;
    fn insert(&self, vector: Vec<f32>, index: String) -> Result<(), IndexError>;
    fn remove(&self, index: String) -> Result<(), IndexError>;
    fn search(&self, query: Vec<f32>, k: usize) -> Result<(Vec<String>, Vec<Vec<f32>>), IndexError>;
}
