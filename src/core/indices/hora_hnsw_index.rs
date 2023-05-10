use crate::core::index::{Index, IndexError};

pub struct HoraHNSWIndex {
    index: usize,
    dimension: usize,
}

impl Index for HoraHNSWIndex {
    fn new(dimension: usize) -> Result<Self, IndexError> where Self: Sized {
        todo!()
    }

    fn insert(&self, vector: Vec<f32>, index: String) -> Result<(), IndexError> {
        todo!()
    }

    fn remove(&self, index: String) -> Result<(), IndexError> {
        todo!()
    }

    fn search(&self, query: Vec<f32>, k: usize) -> Result<(Vec<String>, Vec<Vec<f32>>), IndexError> {
        todo!()
    }
}
